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
	"github.com/apache/incubator-devlake/helpers/utils"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

var _ plugin.SubTaskEntryPoint = ExtractBug

var ExtractBugMeta = plugin.SubTaskMeta{
	Name:             "extractBug",
	EntryPoint:       ExtractBug,
	EnabledByDefault: true,
	Description:      "extract Zentao bug",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ExtractBug(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	statusMappings := getBugStatusMapping(data)
	dueDateField := "deadline"
	if data.Options.ScopeConfig != nil && data.Options.ScopeConfig.BugDueDateField != "" {
		dueDateField = data.Options.ScopeConfig.BugDueDateField
	}
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_BUG_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			res := &models.ZentaoBugRes{}
			err := json.Unmarshal(row.Data, res)
			if err != nil {
				return nil, errors.Default.WrapRaw(err)
			}
			err = res.SetAllFeilds(row.Data)
			if err != nil {
				return nil, errors.Default.WrapRaw(err)
			}
			data.Bugs[res.ID] = struct{}{}
			bug := &models.ZentaoBug{
				ConnectionId:   data.Options.ConnectionId,
				ID:             res.ID,
				Project:        data.Options.ProjectId,
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
				OpenedById:     data.AccountCache.getAccountIDFromApiAccount(res.OpenedBy),
				OpenedByName:   data.AccountCache.getAccountNameFromApiAccount(res.OpenedBy),
				OpenedDate:     res.OpenedDate,
				OpenedBuild:    res.OpenedBuild,
				AssignedToId:   data.AccountCache.getAccountIDFromApiAccount(res.AssignedTo),
				AssignedToName: data.AccountCache.getAccountNameFromApiAccount(res.AssignedTo),
				AssignedDate:   res.AssignedDate,
				Deadline:       res.Deadline,
				ResolvedById:   data.AccountCache.getAccountIDFromApiAccount(res.ResolvedBy),
				Resolution:     res.Resolution,
				ResolvedBuild:  res.ResolvedBuild,
				ResolvedDate:   res.ResolvedDate,
				ClosedById:     data.AccountCache.getAccountIDFromApiAccount(res.ClosedBy),
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
				LastEditedById: data.AccountCache.getAccountIDFromApiAccount(res.LastEditedBy),
				LastEditedDate: res.LastEditedDate,
				Deleted:        res.Deleted,
				PriOrder:       res.PriOrder.String(),
				SeverityOrder:  res.SeverityOrder,
				Needconfirm:    res.Needconfirm,
				StatusName:     res.StatusName,
				ProductStatus:  res.ProductStatus,
				Url:            row.Url,
			}
			loc, _ := time.LoadLocation("Asia/Shanghai")
			bug.DueDate, _ = utils.GetTimeFeildFromMap(res.AllFeilds, dueDateField, loc)
			switch bug.Status {
			case "active", "closed", "resolved":
			default:
				bug.Status = "active"
			}
			if bug.StdType == "" {
				bug.StdType = ticket.BUG
			}
			if len(statusMappings) != 0 {
				if stdStatus, ok := statusMappings[bug.Status]; ok {
					bug.StdStatus = stdStatus
				} else {
					bug.StdStatus = bug.Status
				}
			} else {
				bug.StdStatus = ticket.GetStatus(&ticket.StatusRule{
					Done:    []string{"resolved"},
					Default: ticket.IN_PROGRESS,
				}, bug.Status)
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
