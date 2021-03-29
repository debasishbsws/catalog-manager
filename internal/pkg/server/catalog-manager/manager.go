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
	"github.com/napptive/catalog-manager/internal/pkg/provider"
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
	readmeFile     = "readme.md"
	// apiMetadataVersion with the version of the metadata entity
	apiMetadataVersion     = "core.oam.dev/v1alpha1"
	// appMetadataKind with the kind of the metadata entity
	appMetadataKind           = "ApplicationMetadata"
)

type Manager struct {
	stManager *storage.StorageManager
	provider  provider.MetadataProvider
	// repositoryURL is the URL of the repository managed by this catalog
	repositoryURL string
}

// NewManager returns a new object of manager
func NewManager(manager *storage.StorageManager, provider provider.MetadataProvider) *Manager {
	return &Manager{
		stManager: manager,
		provider:  provider,
	}
}

// decomposeRepositoryName gets the url, repo, application and version from repository name
// and returns the url, the applicationID and an error it something fails
// [repoURL/]repoName/appName[:tag]
func (m *Manager) decomposeRepositoryName(name string) (string, *entities.ApplicationID, error) {
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
func (m *Manager) getApplicationMetadataFile(files []*entities.FileInfo) []byte {

	for _, file := range files {
		// 1.- the files must have .yaml extension
		if utils.IsYamlFile(strings.ToLower(file.Path)) {
			// 2.- Get Metadata
			isMetadata := utils.CheckKindAndVersion(file.Data, apiMetadataVersion, appMetadataKind)
			if isMetadata {
				log.Debug().Str("name", file.Path).Msg("Metadata found")
				return file.Data
			}
		}
	}
	return nil
}

// Add Adds a new application in the repository.
func (m *Manager) Add(name string, files []*entities.FileInfo) error {

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
	if url != m.repositoryURL {
		log.Err(err).Str("name", name).Msg("Error adding application. The application url does not match the one in the catalog")
		return nerrors.NewInternalError("The application url does not match the one in the catalog")
	}

	readme := utils.GetFile(readmeFile, files)
	appMetadata := m.getApplicationMetadataFile(files)
	// the metadata file is required, if is not in the Files -> return an error
	if appMetadata == nil {
		return nerrors.NewNotFoundError("Unable to add the application. Metadata file is required.")
	}

	if _, err := m.provider.Add(&entities.ApplicationMetadata{
		Repository:      appID.Repository,
		ApplicationName: appID.ApplicationName,
		Tag:             appID.Tag,
		Readme:          string(readme),
		Metadata:        string(appMetadata),
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
