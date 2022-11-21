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

var _ core.SubTaskEntryPoint = ExtractBug

var ExtractBugMeta = core.SubTaskMeta{
	Name:             "extractBug",
	EntryPoint:       ExtractBug,
	EnabledByDefault: true,
	Description:      "extract Zentao bug",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ExtractBug(taskCtx core.SubTaskContext) errors.Error {
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
			Table: RAW_BUG_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			res := &models.ZentaoBugRes{}
			err := json.Unmarshal(row.Data, res)
			if err != nil {
				return nil, errors.Default.WrapRaw(err)
			}
			bug := &models.ZentaoBug{
				ConnectionId:   data.Options.ConnectionId,
				ID:             res.ID,
				Project:        res.Project,
				Product:        res.Product,
				Injection:      res.Injection,
				Identify:       res.Identify,
				Branch:         res.Branch,
				Module:         res.Module,
				Execution:      res.Execution,
				Plan:           res.Plan,
				Story:          res.Story,
				StoryVersion:   res.StoryVersion,
				Task:           res.Task,
				ToTask:         res.ToTask,
				ToStory:        res.ToStory,
				Title:          res.Title,
				Keywords:       res.Keywords,
				Severity:       res.Severity,
				Pri:            res.Pri,
				Type:           res.Type,
				Os:             res.Os,
				Browser:        res.Browser,
				Hardware:       res.Hardware,
				Found:          res.Found,
				Steps:          res.Steps,
				Status:         res.Status,
				SubStatus:      res.SubStatus,
				Color:          res.Color,
				Confirmed:      res.Confirmed,
				ActivatedCount: res.ActivatedCount,
				ActivatedDate:  res.ActivatedDate,
				FeedbackBy:     res.FeedbackBy,
				NotifyEmail:    res.NotifyEmail,
				OpenedById:     res.OpenedBy.ID,
				OpenedByName:   res.OpenedBy.Realname,
				OpenedDate:     res.OpenedDate,
				OpenedBuild:    res.OpenedBuild,
				AssignedToId:   res.AssignedTo.ID,
				AssignedToName: res.AssignedTo.Realname,
				AssignedDate:   res.AssignedDate,
				Deadline:       res.Deadline,
				ResolvedById:   res.ResolvedBy.ID,
				Resolution:     res.Resolution,
				ResolvedBuild:  res.ResolvedBuild,
				ResolvedDate:   res.ResolvedDate,
				ClosedById:     res.ClosedBy.ID,
				ClosedDate:     res.ClosedDate,
				DuplicateBug:   res.DuplicateBug,
				LinkBug:        res.LinkBug,
				Feedback:       res.Feedback,
				Result:         res.Result,
				Repo:           res.Repo,
				Mr:             res.Mr,
				Entry:          res.Entry,
				NumOfLine:      res.NumOfLine,
				V1:             res.V1,
				V2:             res.V2,
				RepoType:       res.RepoType,
				IssueKey:       res.IssueKey,
				Testtask:       res.Testtask,
				LastEditedById: res.LastEditedBy.ID,
				LastEditedDate: res.LastEditedDate,
				Deleted:        res.Deleted,
				PriOrder:       res.PriOrder,
				SeverityOrder:  res.SeverityOrder,
				Needconfirm:    res.Needconfirm,
				StatusName:     res.StatusName,
				ProductStatus:  res.ProductStatus,
			}
			results := make([]interface{}, 0)
			results = append(results, bug)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
