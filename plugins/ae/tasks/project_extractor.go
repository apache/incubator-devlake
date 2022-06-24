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
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/plugins/ae/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type ApiProjectResponse struct {
	Id           int        `json:"id"`
	GitUrl       string     `json:"git_url"`
	Priority     int        `json:"priority"`
	AECreateTime *time.Time `json:"create_time"`
	AEUpdateTime *time.Time `json:"update_time"`
}

func ExtractProject(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*AeTaskData)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: AeApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_PROJECT_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			body := &ApiProjectResponse{}
			err := json.Unmarshal(row.Data, body)
			if err != nil {
				return nil, err
			}
			aeProject := &models.AEProject{
				ConnectionId: data.Options.ConnectionId,

				Id:           strconv.Itoa(body.Id),
				GitUrl:       body.GitUrl,
				Priority:     body.Priority,
				AECreateTime: body.AECreateTime,
				AEUpdateTime: body.AEUpdateTime,
			}
			return []interface{}{aeProject}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractProjectMeta = core.SubTaskMeta{
	Name:             "extractProject",
	EntryPoint:       ExtractProject,
	EnabledByDefault: true,
	Description:      "Extract raw project data into tool layer table ae_projects",
}
