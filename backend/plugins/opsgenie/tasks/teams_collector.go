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
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_TEAMS_TABLE = "opsgenie_teams"

type (
	collectedTeams struct {
		TotalCount int               `json:"totalCount"`
		Data       []json.RawMessage `json:"data"`
	}
)

var CollectTeamsMeta = plugin.SubTaskMeta{
	Name:             "collectTeams",
	EntryPoint:       CollectTeams,
	EnabledByDefault: true,
	Description:      "collect Opsgenie teams.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

func CollectTeams(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*OpsgenieTaskData)
	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_TEAMS_TABLE,
		},
		ApiClient:   data.Client,
		UrlTemplate: "v2/teams",
		Query:       nil,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			rawResult := collectedTeams{}
			err := api.UnmarshalResponse(res, &rawResult)

			return rawResult.Data, err
		},
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}
