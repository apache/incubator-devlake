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
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/argocd/models"
)

var _ plugin.SubTaskEntryPoint = ExtractSyncOperations

var ExtractSyncOperationsMeta = plugin.SubTaskMeta{
	Name:             "extractSyncOperations",
	EntryPoint:       ExtractSyncOperations,
	EnabledByDefault: true,
	Description:      "Extract sync operations from raw data",
	DependencyTables: []string{RAW_SYNC_OPERATION_TABLE},
	ProductTables:    []string{models.ArgocdSyncOperation{}.TableName()},
}

type ArgocdApiSyncOperation struct {
	// For history entries
	ID              int64      `json:"id"`
	Revision        string     `json:"revision"`
	DeployedAt      time.Time  `json:"deployedAt"`
	DeployStartedAt *time.Time `json:"deployStartedAt"`
	Source          struct {
		RepoURL string `json:"repoURL"`
	} `json:"source"`
	InitiatedBy struct {
		Username  string `json:"username"`
		Automated bool   `json:"automated"`
	} `json:"initiatedBy"`

	// For operationState (current operation)
	Phase      string     `json:"phase"` // Succeeded, Failed, Error, Running, Terminating
	Message    string     `json:"message"`
	StartedAt  time.Time  `json:"startedAt"`
	FinishedAt *time.Time `json:"finishedAt"`
	SyncResult struct {
		Revision  string                      `json:"revision"`
		Resources []ArgocdApiSyncResourceItem `json:"resources"`
	} `json:"syncResult"`
}

type ArgocdApiSyncResourceItem struct {
	Group     string `json:"group"`
	Version   string `json:"version"`
	Kind      string `json:"kind"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

func ExtractSyncOperations(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ArgocdTaskData)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:   taskCtx,
			Table: RAW_SYNC_OPERATION_TABLE,
			Params: models.ArgocdApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.ApplicationName,
			},
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			var apiOp ArgocdApiSyncOperation
			err := json.Unmarshal(row.Data, &apiOp)
			if err != nil {
				return nil, errors.Default.Wrap(err, "error unmarshaling sync operation")
			}

			syncOp := &models.ArgocdSyncOperation{
				ConnectionId:    data.Options.ConnectionId,
				ApplicationName: data.Options.ApplicationName,
				NoPKModel:       common.NewNoPKModel(),
			}

			normalize := func(t time.Time) *time.Time {
				if t.IsZero() {
					return nil
				}
				return &t
			}
			normalizePtr := func(t *time.Time) *time.Time {
				if t == nil {
					return nil
				}
				if t.IsZero() {
					return nil
				}
				return t
			}

			isOperationState := apiOp.Phase != ""

			if !isOperationState && apiOp.DeployedAt.IsZero() && apiOp.Revision == "" {
				return nil, nil
			}

			if isOperationState {
				start := normalize(apiOp.StartedAt)
				if start != nil {
					syncOp.DeploymentId = start.Unix()
				} else {
					syncOp.DeploymentId = time.Now().Unix()
				}
				syncOp.Revision = apiOp.SyncResult.Revision
				syncOp.StartedAt = start
				syncOp.FinishedAt = normalizePtr(apiOp.FinishedAt)
				syncOp.Phase = apiOp.Phase
				syncOp.Message = apiOp.Message
				syncOp.ResourcesCount = len(apiOp.SyncResult.Resources)
			} else {
				deployedAt := normalize(apiOp.DeployedAt)
				if deployedAt != nil {
					syncOp.DeploymentId = deployedAt.Unix()
				} else {
					syncOp.DeploymentId = time.Now().Unix()
				}

				syncOp.Revision = apiOp.Revision
				syncOp.FinishedAt = deployedAt
				start := normalizePtr(apiOp.DeployStartedAt)
				if start == nil {
					start = deployedAt
				}
				syncOp.StartedAt = start
				syncOp.Phase = "Succeeded"

				if apiOp.InitiatedBy.Automated {
					syncOp.InitiatedBy = "automated"
				} else {
					syncOp.InitiatedBy = apiOp.InitiatedBy.Username
				}
			}

			syncOp.Kind = extractPrimaryDeploymentKind(apiOp.SyncResult.Resources)

			return []interface{}{syncOp}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

// extractPrimaryDeploymentKind identifies the primary deployment resource kind from the resources list.
// It prioritizes actual deployment resources over supporting resources like Services or Ingresses.
func extractPrimaryDeploymentKind(resources []ArgocdApiSyncResourceItem) string {
	// Priority order for deployment resources
	priorityKinds := []string{
		"Rollout",     // Argo Rollouts
		"Deployment",  // Standard K8s Deployment
		"StatefulSet", // Stateful applications
		"DaemonSet",   // Node-level deployments
		"ReplicaSet",  // Direct ReplicaSet management
		"Job",         // Batch jobs
		"CronJob",     // Scheduled jobs
	}

	for _, priorityKind := range priorityKinds {
		for _, resource := range resources {
			if resource.Kind == priorityKind {
				return priorityKind
			}
		}
	}

	for _, resource := range resources {
		if resource.Kind != "" {
			return resource.Kind
		}
	}

	return ""
}
