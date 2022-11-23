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
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

var _ core.SubTaskEntryPoint = ExtractTask

var ExtractTaskMeta = core.SubTaskMeta{
	Name:             "extractTask",
	EntryPoint:       ExtractTask,
	EnabledByDefault: true,
	Description:      "extract Zentao task",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ExtractTask(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ZentaoApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProductId:    data.Options.ProductId,
				ExecutionId:  data.Options.ExecutionId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_TASK_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			res := &models.ZentaoTaskRes{}
			err := json.Unmarshal(row.Data, res)
			if err != nil {
				return nil, errors.Default.WrapRaw(err)
			}
			task := &models.ZentaoTask{
				ConnectionId:       data.Options.ConnectionId,
				ID:                 res.Id,
				Project:            res.Project,
				Parent:             res.Parent,
				Execution:          res.Execution,
				Module:             res.Module,
				Design:             res.Design,
				Story:              res.Story,
				StoryVersion:       res.StoryVersion,
				DesignVersion:      res.DesignVersion,
				FromBug:            res.FromBug,
				Feedback:           res.Feedback,
				FromIssue:          res.FromIssue,
				Name:               res.Name,
				Type:               res.Type,
				Mode:               res.Mode,
				Pri:                res.Pri,
				Estimate:           res.Estimate,
				Consumed:           res.Consumed,
				Deadline:           res.Deadline,
				Status:             res.Status,
				SubStatus:          res.SubStatus,
				Color:              res.Color,
				Description:        res.Description,
				Version:            res.Version,
				OpenedById:         getAccountId(res.OpenedBy),
				OpenedByName:       getAccountName(res.OpenedBy),
				OpenedDate:         res.OpenedDate,
				AssignedToId:       getAccountId(res.AssignedTo),
				AssignedToName:     getAccountName(res.AssignedTo),
				AssignedDate:       res.AssignedDate,
				EstStarted:         res.EstStarted,
				RealStarted:        res.RealStarted,
				FinishedId:         getAccountId(res.FinishedBy),
				FinishedDate:       res.FinishedDate,
				FinishedList:       res.FinishedList,
				CanceledId:         getAccountId(res.CanceledBy),
				CanceledDate:       res.CanceledDate,
				ClosedById:         getAccountId(res.ClosedBy),
				ClosedDate:         res.ClosedDate,
				PlanDuration:       res.PlanDuration,
				RealDuration:       res.RealDuration,
				ClosedReason:       res.ClosedReason,
				LastEditedId:       getAccountId(res.LastEditedBy),
				LastEditedDate:     res.LastEditedDate,
				ActivatedDate:      res.ActivatedDate,
				OrderIn:            res.OrderIn,
				Repo:               res.Repo,
				Mr:                 res.Mr,
				Entry:              res.Entry,
				NumOfLine:          res.NumOfLine,
				V1:                 res.V1,
				V2:                 res.V2,
				Deleted:            res.Deleted,
				Vision:             res.Vision,
				StoryID:            res.Story,
				StoryTitle:         res.StoryTitle,
				LatestStoryVersion: 0,
				//Product:            getAccountId(res.Product),
				//Branch:             res.Branch,
				//LatestStoryVersion: res.LatestStoryVersion,
				//StoryStatus:        res.StoryStatus,
				AssignedToRealName: res.AssignedToRealName,
				PriOrder:           res.PriOrder,
				NeedConfirm:        res.NeedConfirm,
				//ProductType:        res.ProductType,
				Progress: res.Progress,
			}
			results := make([]interface{}, 0)
			results = append(results, task)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
