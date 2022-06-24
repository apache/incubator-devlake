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
	"io/ioutil"
	"net/http"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const RAW_PROJECT_TABLE = "ae_project"

func CollectProject(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*AeTaskData)
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: AeApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_PROJECT_TABLE,
		},
		ApiClient:   data.ApiClient,
		UrlTemplate: "projects/{{ .Params.ProjectId }}",
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}
			res.Body.Close()
			return []json.RawMessage{
				json.RawMessage(body),
			}, nil
		},
	})

	if err != nil {
		return err
	}
	return collector.Execute()
}

var CollectProjectMeta = core.SubTaskMeta{
	Name:             "collectProject",
	EntryPoint:       CollectProject,
	EnabledByDefault: true,
	Description:      "Collect analysis project data from AE api",
}
