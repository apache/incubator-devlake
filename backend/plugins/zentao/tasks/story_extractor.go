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
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	helpers "github.com/apache/incubator-devlake/helpers/utils"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

var _ plugin.SubTaskEntryPoint = ExtractStory

var ExtractStoryMeta = plugin.SubTaskMeta{
	Name:             "extractStory",
	EntryPoint:       ExtractStory,
	EnabledByDefault: true,
	Description:      "extract Zentao story",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ExtractStory(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	statusMappings := getStoryStatusMapping(data)
	dueDateField := ""
	if data.Options.ScopeConfig != nil && data.Options.ScopeConfig.StoryDueDateField != "" {
		dueDateField = data.Options.ScopeConfig.StoryDueDateField
	}

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_STORY_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			var inputParams storyInput
			err := json.Unmarshal(row.Input, &inputParams)
			if err != nil {
				return nil, errors.Default.WrapRaw(err)
			}
			res := &models.ZentaoStoryRes{}
			err = json.Unmarshal(row.Data, res)
			if err != nil {
				return nil, errors.Default.WrapRaw(err)
			}
			data.Stories[res.ID] = struct{}{}
			var results []interface{}
			projectStory := &models.ZentaoProjectStory{
				ConnectionId: data.Options.ConnectionId,
				ProjectId:    data.Options.ProjectId,
				StoryId:      res.ID,
			}
			results = append(results, projectStory)
			story := &models.ZentaoStory{
				ConnectionId:     data.Options.ConnectionId,
				ID:               res.ID,
				Product:          res.Product,
				Branch:           res.Branch,
				Version:          res.Version,
				OrderIn:          0,
				Vision:           res.Vision,
				Parent:           res.Parent,
				Module:           res.Module,
				Plan:             res.Plan,
				Source:           res.Source,
				SourceNote:       res.SourceNote,
				FromBug:          res.FromBug,
				Feedback:         res.Feedback,
				Title:            res.Title,
				Keywords:         res.Keywords,
				Type:             res.Type,
				Category:         res.Category,
				Pri:              res.Pri,
				Estimate:         res.Estimate,
				Status:           res.Status,
				SubStatus:        res.SubStatus,
				Color:            res.Color,
				Stage:            res.Stage,
				Lib:              res.Lib,
				FromStory:        res.FromStory,
				FromVersion:      res.FromVersion,
				OpenedById:       data.AccountCache.getAccountIDFromApiAccount(res.OpenedBy),
				OpenedByName:     data.AccountCache.getAccountNameFromApiAccount(res.OpenedBy),
				OpenedDate:       res.OpenedDate,
				AssignedToId:     data.AccountCache.getAccountIDFromApiAccount(res.AssignedTo),
				AssignedToName:   data.AccountCache.getAccountNameFromApiAccount(res.AssignedTo),
				AssignedDate:     res.AssignedDate,
				ApprovedDate:     res.ApprovedDate,
				LastEditedId:     data.AccountCache.getAccountIDFromApiAccount(res.LastEditedBy),
				LastEditedDate:   res.LastEditedDate,
				ChangedDate:      res.ChangedDate,
				ReviewedById:     data.AccountCache.getAccountIDFromApiAccount(res.ReviewedBy),
				ReviewedDate:     res.ReviewedDate,
				ClosedId:         data.AccountCache.getAccountIDFromApiAccount(res.ClosedBy),
				ClosedDate:       res.ClosedDate,
				ClosedReason:     res.ClosedReason,
				ActivatedDate:    res.ActivatedDate,
				ToBug:            res.ToBug,
				ChildStories:     res.ChildStories,
				LinkStories:      res.LinkStories,
				LinkRequirements: res.LinkRequirements,
				DuplicateStory:   res.DuplicateStory,
				StoryChanged:     res.StoryChanged,
				FeedbackBy:       res.FeedbackBy,
				NotifyEmail:      res.NotifyEmail,
				URChanged:        res.URChanged,
				Deleted:          res.Deleted,
				PriOrder:         res.PriOrder.String(),
				PlanTitle:        res.PlanTitle,
				Url:              row.Url,
			}
			if dueDateField != "" {
				err = res.SetAllFeilds(row.Data)
				if err != nil {
					return nil, errors.Default.WrapRaw(err)
				}
				loc, _ := time.LoadLocation("Asia/Shanghai")
				story.DueDate, _ = helpers.GetTimeFieldFromMap(res.AllFields, dueDateField, loc)
			}
			if story.StdType == "" {
				story.StdType = ticket.REQUIREMENT
			}
			switch story.Status {
			case "active", "closed", "draft", "changing", "reviewing":
			default:
				story.Status = "active"
			}
			if len(statusMappings) != 0 {
				if stdStatus, ok := statusMappings[story.Status]; ok {
					story.StdStatus = stdStatus
				} else {
					story.StdStatus = story.Status
				}
			} else {
				story.StdStatus = ticket.GetStatus(&ticket.StatusRule{
					Done:    []string{"closed"},
					Todo:    []string{"wait"},
					Default: ticket.IN_PROGRESS,
				}, story.Stage)
			}

			results = append(results, story)
			if inputParams.ExecutionId != 0 {
				executionStory := &models.ZentaoExecutionStory{
					ConnectionId: data.Options.ConnectionId,
					ProjectId:    data.Options.ProjectId,
					ExecutionId:  inputParams.ExecutionId,
					StoryId:      story.ID,
				}
				results = append(results, executionStory)
			}
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
