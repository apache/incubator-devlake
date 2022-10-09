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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = EnrichStoryStatusLastStep

var EnrichStoryStatusLastStepMeta = core.SubTaskMeta{
	Name:             "enrichStoryStatusLastStep",
	EntryPoint:       EnrichStoryStatusLastStep,
	EnabledByDefault: true,
	Description:      "Enrich raw data into tool layer table _tool_tapd_story_status",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func EnrichStoryStatusLastStep(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_STORY_STATUS_LAST_STEP_TABLE, false)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			var storyStatusLastStepRes struct {
				Data map[string]string
			}
			err := errors.Convert(json.Unmarshal(row.Data, &storyStatusLastStepRes))
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0)
			statusList := make([]*models.TapdStoryStatus, 0)
			clauses := []dal.Clause{
				dal.Where("connection_id = ? and workspace_id = ?", data.Options.ConnectionId, data.Options.WorkspaceId),
			}
			err = db.All(&statusList, clauses...)
			if err != nil {
				return nil, err
			}

			for _, status := range statusList {
				if storyStatusLastStepRes.Data[status.EnglishName] != "" {
					status.IsLastStep = true
					results = append(results, status)
				}
			}

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
