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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"strings"
)

var _ plugin.SubTaskEntryPoint = ExtractStories

var ExtractStoryMeta = plugin.SubTaskMeta{
	Name:             "extractStories",
	EntryPoint:       ExtractStories,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ExtractStories(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_STORY_TABLE)
	db := taskCtx.GetDal()
	statusList := make([]models.TapdStoryStatus, 0)
	statusLanguageMap, getStdStatus, err := getDefaultStdStatusMapping(data, db, statusList)
	if err != nil {
		return err
	}
	customStatusMap := getStatusMapping(data)
	stdTypeMappings := getStdTypeMappings(data)
	typeIdMapping, err := getTapdTypeMappings(data, db, "story")
	if err != nil {
		return err
	}
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		BatchSize:          100,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			var storyBody struct {
				Story models.TapdStory
			}
			err = errors.Convert(json.Unmarshal(row.Data, &storyBody))
			if err != nil {
				return nil, err
			}
			toolL := storyBody.Story
			toolL.Status = statusLanguageMap[toolL.Status]
			if len(customStatusMap) != 0 {
				toolL.StdStatus = customStatusMap[toolL.Status]
			} else {
				toolL.StdStatus = getStdStatus(toolL.Status)
			}

			toolL.ConnectionId = data.Options.ConnectionId
			toolL.Priority = priorityMap[toolL.Priority]
			toolL.Type = typeIdMapping[toolL.WorkitemTypeId]
			toolL.StdType = stdTypeMappings[toolL.Type]
			if toolL.StdType == "" {
				toolL.StdType = ticket.REQUIREMENT
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
