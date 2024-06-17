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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
)

var ExtractApiStagesMeta = plugin.SubTaskMeta{
	Name:             "extractApiStages",
	EntryPoint:       ExtractApiStages,
	EnabledByDefault: true,
	Description:      "Extract raw stages data into tool layer table jenkins_stages",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func ExtractApiStages(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JenkinsTaskData)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Params: JenkinsApiParams{
				ConnectionId: data.Options.ConnectionId,
				FullName:     data.Options.JobFullName,
			},
			Ctx:   taskCtx,
			Table: RAW_STAGE_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			body := &models.Stage{}
			err := errors.Convert(json.Unmarshal(row.Data, body))
			if err != nil {
				return nil, err
			}
			input := &SimpleBuild{}
			err = errors.Convert(json.Unmarshal(row.Input, input))
			if err != nil {
				return nil, err
			}

			results := make([]interface{}, 0)

			stage := &models.JenkinsStage{
				ConnectionId:        data.Options.ConnectionId,
				ID:                  body.ID,
				Name:                body.Name,
				ExecNode:            body.ExecNode,
				Status:              body.Status,
				StartTimeMillis:     body.StartTimeMillis,
				DurationMillis:      body.DurationMillis,
				PauseDurationMillis: body.PauseDurationMillis,
				BuildName:           input.FullName,
			}

			results = append(results, stage)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
