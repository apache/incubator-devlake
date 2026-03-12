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
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/argocd/models"
)

const RAW_SYNC_OPERATION_TABLE = "argocd_api_sync_operations"

var _ plugin.SubTaskEntryPoint = CollectSyncOperations

var CollectSyncOperationsMeta = plugin.SubTaskMeta{
	Name:             "collectSyncOperations",
	EntryPoint:       CollectSyncOperations,
	EnabledByDefault: true,
	Description:      "Collect sync operations (deployment history) from ArgoCD API",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	DependencyTables: []string{models.ArgocdApplication{}.TableName()},
}

func CollectSyncOperations(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ArgocdTaskData)

	collector, err := api.NewStatefulApiCollector(api.RawDataSubTaskArgs{
		Ctx:   taskCtx,
		Table: RAW_SYNC_OPERATION_TABLE,
		Params: models.ArgocdApiParams{
			ConnectionId: data.Options.ConnectionId,
			Name:         data.Options.ApplicationName,
		},
	})

	if err != nil {
		return err
	}

	err = collector.InitCollector(api.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		UrlTemplate: "/applications/{{.Params.Name}}",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("refresh", "false")
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var appResponse struct {
				Status struct {
					History        []json.RawMessage `json:"history"`
					OperationState json.RawMessage   `json:"operationState"`
				} `json:"status"`
			}

			err := api.UnmarshalResponse(res, &appResponse)
			if err != nil {
				return nil, err
			}

			return sanitizeOperationEntries(appResponse.Status.OperationState, appResponse.Status.History), nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}

func sanitizeOperationEntries(operationState json.RawMessage, history []json.RawMessage) []json.RawMessage {
	results := make([]json.RawMessage, 0, len(history)+1)

	if op := bytes.TrimSpace(operationState); len(op) > 0 && !bytes.Equal(op, []byte("null")) {
		results = append(results, operationState)
	}

	for _, entry := range history {
		if h := bytes.TrimSpace(entry); len(h) == 0 || bytes.Equal(h, []byte("null")) {
			continue
		}
		results = append(results, entry)
	}

	return results
}
