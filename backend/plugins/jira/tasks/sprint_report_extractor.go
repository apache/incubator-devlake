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

	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/plugin"
	"github.com/apache/devlake/helpers/pluginhelper/api"
	"github.com/apache/devlake/plugins/jira/models"
	"github.com/apache/devlake/plugins/jira/tasks/apiv2models"
)

var _ plugin.SubTaskEntryPoint = ExtractSprintReport

var ExtractSprintReportMeta = plugin.SubTaskMeta{
	Name:             "extractSprintReport",
	EntryPoint:       ExtractSprintReport,
	EnabledByDefault: true,
	Description:      "extract Jira Sprint Report",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ExtractSprintReport(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	connectionId := data.Options.ConnectionId

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_SPRINT_REPORT_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			var report apiv2models.SprintReport
			if err := errors.Convert(json.Unmarshal(row.Data, &report)); err != nil {
				return nil, err
			}
			var input apiv2models.SprintReportInput
			if err := errors.Convert(json.Unmarshal(row.Input, &input)); err != nil {
				return nil, err
			}

			var result []interface{}
			addBucket := func(issues []apiv2models.SprintReportIssue, bucket string) {
				for _, issue := range issues {
					result = append(result, &models.JiraSprintReport{
						ConnectionId:             connectionId,
						BoardId:                  input.BoardId,
						SprintId:                 input.SprintId,
						IssueId:                  issue.Id,
						IssueKey:                 issue.Key,
						Bucket:                   bucket,
						Done:                     issue.Done,
						StoryPointsAtSprintStart: issue.EstimateStatistic.StatFieldValue.Value,
						StoryPointsAtSprintEnd:   issue.CurrentEstimateStatistic.StatFieldValue.Value,
					})
				}
			}
			addBucket(report.Contents.CompletedIssues, models.SprintReportBucketCompleted)
			addBucket(report.Contents.IssuesNotCompletedInCurrentSprint, models.SprintReportBucketNotCompleted)
			addBucket(report.Contents.PuntedIssues, models.SprintReportBucketPunted)
			addBucket(report.Contents.IssuesCompletedInAnotherSprint, models.SprintReportBucketCompletedInOtherSprint)

			return result, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}
