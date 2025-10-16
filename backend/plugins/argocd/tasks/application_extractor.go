/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tasks

import (
	"encoding/json"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/argocd/models"
)

var _ plugin.SubTaskEntryPoint = ExtractApplications

var ExtractApplicationsMeta = plugin.SubTaskMeta{
	Name:             "extractApplications",
	EntryPoint:       ExtractApplications,
	EnabledByDefault: true,
	Description:      "Extract applications from raw data",
	DependencyTables: []string{RAW_APPLICATION_TABLE},
	ProductTables:    []string{models.ArgocdApplication{}.TableName()},
}

type ArgocdApiApplication struct {
	Metadata struct {
		Name              string    `json:"name"`
		Namespace         string    `json:"namespace"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
	} `json:"metadata"`
	Spec struct {
		Project string `json:"project"`
		Source  struct {
			RepoURL        string `json:"repoURL"`
			Path           string `json:"path"`
			TargetRevision string `json:"targetRevision"`
		} `json:"source"`
		Destination struct {
			Server    string `json:"server"`
			Namespace string `json:"namespace"`
		} `json:"destination"`
	} `json:"spec"`
	Status struct {
		Sync struct {
			Status string `json:"status"` // Synced, OutOfSync, Unknown
		} `json:"sync"`
		Health struct {
			Status string `json:"status"` // Healthy, Progressing, Degraded, etc.
		} `json:"health"`
	} `json:"status"`
}

func ExtractApplications(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ArgocdTaskData)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:   taskCtx,
			Table: RAW_APPLICATION_TABLE,
			Params: models.ArgocdApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.ApplicationName,
			},
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			var apiApp ArgocdApiApplication
			err := json.Unmarshal(row.Data, &apiApp)
			if err != nil {
				return nil, errors.Default.Wrap(err, "error unmarshaling application")
			}

			application := &models.ArgocdApplication{
				Name:           apiApp.Metadata.Name,
				Namespace:      apiApp.Metadata.Namespace,
				Project:        apiApp.Spec.Project,
				RepoURL:        apiApp.Spec.Source.RepoURL,
				Path:           apiApp.Spec.Source.Path,
				TargetRevision: apiApp.Spec.Source.TargetRevision,
				DestServer:     apiApp.Spec.Destination.Server,
				DestNamespace:  apiApp.Spec.Destination.Namespace,
				SyncStatus:     apiApp.Status.Sync.Status,
				HealthStatus:   apiApp.Status.Health.Status,
				CreatedDate:    &apiApp.Metadata.CreationTimestamp,
			}
			application.ConnectionId = data.Options.ConnectionId
			if data.Options.ScopeConfig != nil {
				application.ScopeConfigId = data.Options.ScopeConfig.ID
			}

			return []interface{}{application}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
