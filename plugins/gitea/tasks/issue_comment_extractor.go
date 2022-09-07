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

	"github.com/apache/incubator-devlake/plugins/core/dal"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitea/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractApiIssueCommentsMeta = core.SubTaskMeta{
	Name:             "extractApiIssueComments",
	EntryPoint:       ExtractApiIssueComments,
	EnabledByDefault: true,
	Description: "Extract raw comment data  into tool layer table gitea_pull_request_comments" +
		"and gitea_issue_comments",
	DomainTypes: []string{core.DOMAIN_TYPE_TICKET},
}

type IssueComment struct {
	GiteaId int `json:"id"`
	Body    string

	User struct {
		Login string
		Id    int
	}

	Target struct {
		Issue struct {
			Id     int    `json:"id"`
			Title  string `json:"title"`
			Number string `json:"number"`
		}
		PullRequest string `json:"pull_request"`
	}

	GiteaCreatedAt helper.Iso8601Time `json:"created_at"`
	GiteaUpdatedAt helper.Iso8601Time `json:"updated_at"`
}

func ExtractApiIssueComments(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GiteaTaskData)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GiteaApiParams{
				Owner: data.Options.Owner,
				Repo:  data.Options.Repo,
			},
			Table: RAW_COMMENTS_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			apiComment := &IssueComment{}
			err := json.Unmarshal(row.Data, apiComment)
			if err != nil {
				return nil, err
			}
			// need to extract 2 kinds of entities here
			results := make([]interface{}, 0, 2)
			if apiComment.GiteaId == 0 {
				return nil, nil
			}
			//If this is a pr, ignore
			issueINumber := apiComment.Target.Issue.Number
			if err != nil {
				return nil, err
			}
			issue := &models.GiteaIssue{}
			err = taskCtx.GetDal().All(issue, dal.Where("connection_id = ? and number = ? and repo_id = ?", data.Options.ConnectionId, issueINumber, data.Repo.GiteaId))
			if err != nil {
				return nil, err
			}
			giteaIssueComment := &models.GiteaIssueComment{
				ConnectionId:   data.Options.ConnectionId,
				GiteaId:        apiComment.GiteaId,
				IssueId:        issue.GiteaId,
				Body:           apiComment.Body,
				AuthorName:     apiComment.User.Login,
				AuthorId:       apiComment.User.Id,
				GiteaCreatedAt: apiComment.GiteaCreatedAt.ToTime(),
				GiteaUpdatedAt: apiComment.GiteaUpdatedAt.ToTime(),
			}
			results = append(results, giteaIssueComment)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
