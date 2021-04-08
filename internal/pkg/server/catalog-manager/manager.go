/**
 * Copyright 2021 Napptive
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package catalog_manager

import (
	"github.com/napptive/catalog-manager/internal/pkg/entities"
	"github.com/napptive/catalog-manager/internal/pkg/provider/metadata"
	"github.com/napptive/catalog-manager/internal/pkg/storage"
	"github.com/napptive/catalog-manager/internal/pkg/utils"
	"github.com/napptive/nerrors/pkg/nerrors"
	"github.com/rs/zerolog/log"
	"strings"
)

const (
	// defaultVersion contains the default version, when a user does not fill the version, defaultVersion is used
	defaultVersion = "latest"
	// readmeFile with the name of the readme file
	readmeFile = "readme.md"
)

type Manager interface {
	// Add Adds a new application in the repository.
	Add(name string, files []*entities.FileInfo) error
	// Download returns the files of an application
	Download(appName string) ([]*entities.FileInfo, error)
	// Remove removes an application from the repository
	Remove(appName string) error
	// Get returns a given application metadata
	Get(appName string) (*entities.ExtendedApplicationMetadata, error)
	// List returns a list of applications (without metadata and readme content)
	List() ([]*entities.ApplicationMetadata, error)
}

type manager struct {
	stManager storage.StorageManager
	provider  metadata.MetadataProvider
	// catalogURL is the URL of the repository managed by this catalog
	catalogURL string
}

// NewManager returns a new object of manager
func NewManager(stManager storage.StorageManager, provider metadata.MetadataProvider, catalogURL string) Manager {
	return &manager{
		stManager:  stManager,
		provider:   provider,
		catalogURL: catalogURL,
	}
}

// decomposeRepositoryName gets the url, repo, application and version from repository name
// and returns the url, the applicationID and an error it something fails
// [repoURL/]repoName/appName[:tag]
func (m *manager) decomposeRepositoryName(name string) (string, *entities.ApplicationID, error) {
	var version string
	var applicationName string
	var repoName string
	urlName := ""

	names := strings.Split(name, "/")
	if len(names) != 2 && len(names) != 3 {
		return "", nil, nerrors.NewFailedPreconditionError(
			"incorrect format for application name. [repoURL/]repoName/appName[:tag]")
	}

	// if len == 2 -> no url informed.
	if len(names) == 3 {
		urlName = names[0]
	}
	repoName = names[len(names)-2]

	// get the version -> appName[:tag]
	sp := strings.Split(names[len(names)-1], ":")
	if len(sp) == 1 {
		applicationName = sp[0]
		version = defaultVersion
	} else if len(sp) == 2 {
		applicationName = sp[0]
		version = sp[1]
	} else {
		return "", nil, nerrors.NewFailedPreconditionError(
			"incorrect format for application name. [repoURL/]repoName/appName[:tag]")
	}

	return urlName, &entities.ApplicationID{
		Repository:      repoName,
		ApplicationName: applicationName,
		Tag:             version,
	}, nil
}

// getApplicationMetadataFile looks for the application metadata yaml file
func (m *manager) getApplicationMetadataFile(files []*entities.FileInfo) ([]byte, *entities.CatalogMetadata, error) {

	for _, file := range files {
		// 1.- the files must have .yaml extension
		if utils.IsYamlFile(strings.ToLower(file.Path)) {
			// 2.- Get Metadata
			isMetadata, metadataObj, err := utils.IsMetadata(file.Data)
			if err != nil {
				return nil, nil, err
			}
			if isMetadata {
				return file.Data, metadataObj, nil
			}
		}
	}
	return nil, nil, nil
}

// Add Adds a new application in the repository.
func (m *manager) Add(name string, files []*entities.FileInfo) error {

	// TODO: here, validate the application

	// 1.- Store metadata into the provider
	// Locate README and metadata
	// Store in the provider
	url, appID, err := m.decomposeRepositoryName(name)
	if err != nil {
		log.Err(err).Str("name", name).Msg("Error decomposing the application name")
		return err
	}

	// check that the url of the application matches the url of the catalog
	// we avoid including applications in catalogs that do not correspond
	if url != "" &&  url != m.catalogURL {
		log.Err(err).Str("name", name).Msg("Error adding application. The application url does not match the one in the catalog")
		return nerrors.NewInternalError("The application url does not match the one in the catalog")
	}

	readme := utils.GetFile(readmeFile, files)
	appMetadata, header, err := m.getApplicationMetadataFile(files)
	if err != nil {
		return nerrors.NewInternalErrorFrom(err, "Unable to add the application. Error getting metadata")
	}
	// the metadata file is required, if is not in the Files -> return an error
	if appMetadata == nil {
		return nerrors.NewNotFoundError("Unable to add the application. Metadata file is required.")
	}
	// Metadata Name is required too
	if header == nil || header.Name == "" {
		return nerrors.NewFailedPreconditionError("Unable to add the application. Metadata name is required.")
	}

	if _, err := m.provider.Add(&entities.ApplicationMetadata{
		Repository:      appID.Repository,
		ApplicationName: appID.ApplicationName,
		Tag:             appID.Tag,
		Readme:          string(readme),
		Metadata:        string(appMetadata),
		MetadataName:    header.Name,
	}); err != nil {
		log.Err(err).Str("name", name).Msg("Error storing application metadata")
		return err
	}

	// 2.- store the files into the repository storage
	if err = m.stManager.StoreApplication(appID.Repository, appID.ApplicationName, appID.Tag, files); err != nil {
		log.Err(err).Str("name", name).Msg("Error storing application")
		// rollback operation
		if rErr := m.provider.Remove(appID); rErr != nil {
			log.Err(err).Interface("appID", appID).Msg("Error in rollback operation, metadata can not be removed")
		}
		return err
	}

	return nil
}

// Download returns the files of an application
func (m *manager) Download(appName string) ([]*entities.FileInfo, error) {

	_, appID, err := m.decomposeRepositoryName(appName)
	if err != nil {
		return nil, nerrors.NewFailedPreconditionErrorFrom(err, "unable to download the application")
	}

	return m.stManager.GetApplication(appID.Repository, appID.ApplicationName, appID.Tag)
}

// Remove removes an application from the repository
func (m *manager) Remove(appName string) error {

	// - Validate the appName
	_, appID, err := m.decomposeRepositoryName(appName)
	if err != nil {
		return nerrors.NewFailedPreconditionErrorFrom(err, "unable to remove application, wrong name")
	}

	// - Remove from metadata provider
	err = m.provider.Remove(appID)
	if err != nil {
		log.Err(err).Str("appName", appName).Msg("Unable to remove application metadata")
		return err
	}

	// - Remove from storage
	err = m.stManager.RemoveApplication(appID.Repository, appID.ApplicationName, appID.Tag)
	if err != nil {
		log.Err(err).Str("appName", appName).Msg("Unable to remove application")
		return err
	}

	return nil
}

// Get returns the application metadata for a given application
func (m *manager) Get(appName string) (*entities.ExtendedApplicationMetadata, error) {

	_, appID, err := m.decomposeRepositoryName(appName)
	if err != nil {
		return nil, err
	}

	app, err := m.provider.Get(*appID)
	if err != nil {
		return nil, err
	}

	_, metadata, err := utils.IsMetadata([]byte(app.Metadata))
	if err != nil {
		return nil, err
	}
	var obj entities.CatalogMetadata
	if metadata != nil {
		obj = *metadata
	}

	return &entities.ExtendedApplicationMetadata{
		CatalogID:       app.CatalogID,
		Repository:      app.Repository,
		ApplicationName: app.ApplicationName,
		Tag:             app.Tag,
		Readme:          app.Readme,
		Metadata:        app.Metadata,
		MetadataObj:     obj,
	}, nil
}

// List returns a list of applications (without metadata and readme content)
func (m *manager) List() ([]*entities.ApplicationMetadata, error) {
	return m.provider.List()
}
