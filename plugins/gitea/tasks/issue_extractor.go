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
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitea/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractApiIssuesMeta = core.SubTaskMeta{
	Name:             "extractApiIssues",
	EntryPoint:       ExtractApiIssues,
	EnabledByDefault: true,
	Description:      "Extract raw Issues data into tool layer table gitea_issues",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

type IssuesResponse struct {
	GiteaId       int    `json:"id"`
	Url           string `json:"url"`
	RepositoryUrl string `json:"repository_url"`
	Number        int    `json:"number"`
	State         string `json:"state"`
	Title         string
	Body          string
	HtmlUrl       string `json:"html_url"`
	CommentsUrl   string `json:"comments_url"`
	PullRequest   struct {
		Url     string `json:"url"`
		HtmlUrl string `json:"html_url"`
	} `json:"pull_request"`
	Labels []struct {
		Id           int
		RepositoryId int                `json:"repository_id"`
		Name         string             `json:"name"`
		CreatedAt    helper.Iso8601Time `json:"created_at"`
		UpdatedAt    helper.Iso8601Time `json:"updated_at"`
	} `json:"labels"`
	Repository struct {
		Id       int
		FullName string `json:"full_name"`
		Url      string `json:"url"`
	} `json:"repository"`
	Assignee *struct {
		Login string
		Id    int
	}
	User *struct {
		Login string
		Id    int
		Name  string
	}
	Comments        int                 `json:"comments"`
	Priority        int                 `json:"priority"`
	IssueType       string              `json:"issue_type"`
	SecurityHole    bool                `json:"security_hole"`
	IssueState      string              `json:"issue_state"`
	Branch          string              `json:"branch"`
	FinishAt        *helper.Iso8601Time `json:"finished_at"`
	GiteaCreatedAt  helper.Iso8601Time  `json:"created_at"`
	GiteaUpdatedAt  helper.Iso8601Time  `json:"updated_at"`
	IssueTypeDetail struct {
		Id        int
		Title     string
		Ident     string
		CreatedAt helper.Iso8601Time `json:"created_at"`
		UpdatedAt helper.Iso8601Time `json:"updated_at"`
	}
	IssueStateDetail struct {
		Id        int
		Title     string
		Serial    string
		CreatedAt helper.Iso8601Time `json:"created_at"`
		UpdatedAt helper.Iso8601Time `json:"updated_at"`
	}
}

func ExtractApiIssues(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_ISSUE_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			body := &IssuesResponse{}
			err := json.Unmarshal(row.Data, body)
			if err != nil {
				return nil, err
			}
			// need to extract 2 kinds of entities here
			if body.GiteaId == 0 {
				return nil, nil
			}
			//If this is a pr, ignore
			if body.PullRequest.Url != "" {
				return nil, nil
			}
			results := make([]interface{}, 0, 2)
			giteaIssue, err := convertGiteaIssue(body, data.Options.ConnectionId, data.Repo.GiteaId)
			if err != nil {
				return nil, err
			}
			for _, label := range body.Labels {
				results = append(results, &models.GiteaIssueLabel{
					ConnectionId: data.Options.ConnectionId,
					IssueId:      giteaIssue.GiteaId,
					LabelName:    label.Name,
				})

			}
			results = append(results, giteaIssue)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
func convertGiteaIssue(issue *IssuesResponse, connectionId uint64, repositoryId int) (*models.GiteaIssue, error) {
	giteaIssue := &models.GiteaIssue{
		ConnectionId:   connectionId,
		GiteaId:        issue.GiteaId,
		RepoId:         repositoryId,
		Number:         issue.Number,
		State:          issue.State,
		Title:          issue.Title,
		Body:           issue.Body,
		Url:            issue.HtmlUrl,
		ClosedAt:       helper.Iso8601TimeToTime(issue.FinishAt),
		GiteaCreatedAt: issue.GiteaCreatedAt.ToTime(),
		GiteaUpdatedAt: issue.GiteaUpdatedAt.ToTime(),
	}

	if issue.Assignee != nil {
		giteaIssue.AssigneeId = issue.Assignee.Id
		giteaIssue.AssigneeName = issue.Assignee.Login
	}
	if issue.User != nil {
		giteaIssue.AuthorId = issue.User.Id
		giteaIssue.AuthorName = issue.User.Login
	}
	if issue.FinishAt != nil {
		giteaIssue.LeadTimeMinutes = uint(issue.FinishAt.ToTime().Sub(issue.GiteaCreatedAt.ToTime()).Minutes())
	}

	return giteaIssue, nil
}
