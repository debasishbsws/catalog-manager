/*
 Copyright 2022 Napptive

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package apps

import (
	"github.com/napptive/catalog-manager/internal/pkg/config"
	"github.com/napptive/catalog-manager/internal/pkg/connection"
	catalog_manager "github.com/napptive/catalog-manager/internal/pkg/server/catalog-manager"
	grpc_catalog_common_go "github.com/napptive/grpc-catalog-common-go"
	grpc_playground_apps_go "github.com/napptive/grpc-playground-apps-go"
	"github.com/napptive/nerrors/pkg/nerrors"
	"github.com/rs/zerolog/log"
)

type Manager interface {
	// Deploy an application on a target Playground platform. This endpoint
	// will gather the application information and send it to the target
	// playground platform.
	Deploy(userToken string, applicationID string, targetEnvironmentQualifiedName string, targetPlaygroundApiURL string) (*grpc_catalog_common_go.OpResponse, error)
}

// Manager for apps operations.
type manager struct {
	catalogManager catalog_manager.Manager
	contextHelper  *connection.ContextHelper
	cfg            *config.Config
}

// NewManager creates a new instance of the manager.
func NewManager(cfg *config.Config, catalogManager catalog_manager.Manager) Manager {
	contextHelper := connection.NewContextHelper(cfg)
	return &manager{
		catalogManager: catalogManager,
		contextHelper:  contextHelper,
		cfg:            cfg,
	}
}

// Deploy an application on a target Playground platform. This endpoint
// will gather the application information and send it to the target
// playground platform.
func (m *manager) Deploy(userToken string, applicationID string, targetEnvironmentQualifiedName string, targetPlaygroundApiURL string) (*grpc_catalog_common_go.OpResponse, error) {
	log.Debug().Str("application_id", applicationID).Str("eqn", targetEnvironmentQualifiedName).Str("target_playground_api_url", targetPlaygroundApiURL).Msg("deploying application")
	// Retrieve the target application
	app, err := m.catalogManager.Download(applicationID, true)
	if err != nil {
		return nil, nerrors.FromGRPC(err)
	}

	// GetConnection
	conn, err := connection.GetConnectionToPlayground(&m.cfg.PlaygroundConnection, targetPlaygroundApiURL)
	if err != nil {
		return nil, nerrors.NewInternalErrorFrom(err, "cannot establish connection with Playground server on %s", targetPlaygroundApiURL)
	}
	defer conn.Close()
	client := grpc_playground_apps_go.NewAppsClient(conn)

	// Create a connection with the same JWT token
	ctx, cancel := m.contextHelper.GetContext(userToken)
	defer cancel()

	response, err := client.Deploy(ctx, &grpc_playground_apps_go.DeployApplicationRequest{
		ApplicationData:                app[0].Data,
		TargetEnvironmentQualifiedName: targetEnvironmentQualifiedName,
	})
	if err != nil {
		return nil, nerrors.FromGRPC(err)
	}
	return &grpc_catalog_common_go.OpResponse{
		Status:     grpc_catalog_common_go.OpStatus(grpc_catalog_common_go.OpStatus_value[response.StatusName]),
		StatusName: response.StatusName,
		UserInfo:   response.UserInfo,
	}, nil
}