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
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/testmo/models"
)

var ExtractMilestonesMeta = plugin.SubTaskMeta{
	Name:             "extractMilestones",
	EntryPoint:       ExtractMilestones,
	EnabledByDefault: true,
	Description:      "Extract raw milestones data into tool layer table testmo_milestones",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}

func ExtractMilestones(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*TestmoTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TestmoApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_MILESTONE_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			var apiMilestone struct {
				Id          uint64     `json:"id"`
				ProjectId   uint64     `json:"project_id"`
				Name        string     `json:"name"`
				Description string     `json:"description"`
				IsCompleted bool       `json:"is_completed"`
				CreatedAt   *time.Time `json:"created_at"`
				UpdatedAt   *time.Time `json:"updated_at"`
				CompletedAt *time.Time `json:"completed_at"`
			}

			err := json.Unmarshal(row.Data, &apiMilestone)
			if err != nil {
				return nil, errors.Default.Wrap(err, "error unmarshaling milestone")
			}

			milestone := &models.TestmoMilestone{
				ConnectionId:    data.Options.ConnectionId,
				Id:              apiMilestone.Id,
				ProjectId:       apiMilestone.ProjectId,
				Name:            apiMilestone.Name,
				Description:     apiMilestone.Description,
				IsCompleted:     apiMilestone.IsCompleted,
				TestmoCreatedAt: apiMilestone.CreatedAt,
				TestmoUpdatedAt: apiMilestone.UpdatedAt,
				CompletedAt:     apiMilestone.CompletedAt,
			}

			return []interface{}{milestone}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
