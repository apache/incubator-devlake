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
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"strings"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractStories

var ExtractStoryMeta = core.SubTaskMeta{
	Name:             "extractStories",
	EntryPoint:       ExtractStories,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ExtractStories(taskCtx core.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_STORY_TABLE, true)
	db := taskCtx.GetDal()
	statusList := make([]*models.TapdStoryStatus, 0)
	clauses := []dal.Clause{
		dal.Where("connection_id = ? and workspace_id = ?", data.Options.ConnectionId, data.Options.WorkspaceId),
	}
	err := db.All(&statusList, clauses...)
	if err != nil {
		return err
	}

	statusMap := make(map[string]string, len(statusList))
	for _, v := range statusList {
		statusMap[v.EnglishName] = v.ChineseName
	}
	mappings, err := getTypeMappings(data, db, "story")

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		BatchSize:          100,
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			var storyBody struct {
				Story models.TapdStory
			}
			err = errors.Convert(json.Unmarshal(row.Data, &storyBody))
			if err != nil {
				return nil, err
			}
			toolL := storyBody.Story
			toolL.Status = statusMap[toolL.Status]
			toolL.StdStatus = getStdStatus(toolL.Status)

			toolL.ConnectionId = data.Options.ConnectionId

			toolL.Type = mappings.typeIdMappings[toolL.WorkitemTypeId]
			toolL.StdType = mappings.stdTypeMappings[toolL.Type]
			if toolL.StdType == "" {
				toolL.StdType = "REQUIREMENT"
			}

			toolL.Url = fmt.Sprintf("https://www.tapd.cn/%d/prong/stories/view/%d", toolL.WorkspaceId, toolL.Id)
			if strings.Contains(toolL.Owner, ";") {
				toolL.Owner = strings.Split(toolL.Owner, ";")[0]
			}
			workSpaceStory := &models.TapdWorkSpaceStory{
				ConnectionId: data.Options.ConnectionId,
				WorkspaceId:  toolL.WorkspaceId,
				StoryId:      toolL.Id,
			}
			results := make([]interface{}, 0, 3)
			results = append(results, &toolL, workSpaceStory)
			if toolL.IterationId != 0 {
				iterationStory := &models.TapdIterationStory{
					ConnectionId:     data.Options.ConnectionId,
					IterationId:      toolL.IterationId,
					StoryId:          toolL.Id,
					WorkspaceId:      toolL.WorkspaceId,
					ResolutionDate:   toolL.Completed,
					StoryCreatedDate: toolL.Created,
				}
				results = append(results, iterationStory)
			}
			if toolL.Label != "" {
				labelList := strings.Split(toolL.Label, "|")
				for _, v := range labelList {
					toolLIssueLabel := &models.TapdStoryLabel{
						StoryId:   toolL.Id,
						LabelName: v,
					}
					results = append(results, toolLIssueLabel)
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

type typeMappings struct {
	typeIdMappings  map[uint64]string
	stdTypeMappings map[string]string
}
