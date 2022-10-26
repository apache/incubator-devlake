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
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type IssuesResponse struct {
	Type        string            `json:"type"`
	BitbucketId int               `json:"id"`
	Repository  *BitbucketApiRepo `json:"repository"`
	Links       struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Html struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
	Title   string `json:"title"`
	Content struct {
		Type string `json:"type"`
		Raw  string `json:"raw"`
	} `json:"content"`
	Reporter  *BitbucketAccountResponse `json:"reporter"`
	Assignee  *BitbucketAccountResponse `json:"assignee"`
	State     string                    `json:"state"`
	Kind      string                    `json:"kind"`
	Milestone *struct {
		Id int `json:"id"`
	} `json:"milestone"`
	Component          string    `json:"component"`
	Priority           string    `json:"priority"`
	Votes              int       `json:"votes"`
	Watches            int       `json:"watches"`
	BitbucketCreatedAt time.Time `json:"created_on"`
	BitbucketUpdatedAt time.Time `json:"updated_on"`
}

var ExtractApiIssuesMeta = core.SubTaskMeta{
	Name:             "extractApiIssues",
	EntryPoint:       ExtractApiIssues,
	EnabledByDefault: true,
	Description:      "Extract raw Issues data into tool layer table bitbucket_issues",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ExtractApiIssues(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*BitbucketTaskData)
	config := data.Options.TransformationRules
	issueStatusMap, err := newIssueStatusMap(config)
	if err != nil {
		return nil
	}
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			Params: BitbucketApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			/*
				Table store raw data
			*/
			Table: RAW_ISSUE_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			body := &IssuesResponse{}
			err := errors.Convert(json.Unmarshal(row.Data, body))
			if err != nil {
				return nil, err
			}
			// need to extract 2 kinds of entities here
			if body.BitbucketId == 0 {
				return nil, nil
			}
			//If this is not an issue, ignore
			if body.Type != "issue" {
				return nil, nil
			}
			results := make([]interface{}, 0, 2)

			bitbucketIssue, err := convertBitbucketIssue(body, data.Options.ConnectionId, data.Repo.BitbucketId)
			if err != nil {
				return nil, err
			}

			results = append(results, bitbucketIssue)
			if body.Assignee != nil {
				bitbucketIssue.AssigneeId = body.Assignee.AccountId
				bitbucketIssue.AssigneeName = body.Assignee.DisplayName
				relatedUser, err := convertAccount(body.Assignee, data.Options.ConnectionId)
				if err != nil {
					return nil, err
				}
				results = append(results, relatedUser)
			}
			if body.Reporter != nil {
				bitbucketIssue.AuthorId = body.Reporter.AccountId
				bitbucketIssue.AuthorName = body.Reporter.DisplayName
				relatedUser, err := convertAccount(body.Reporter, data.Options.ConnectionId)
				if err != nil {
					return nil, err
				}
				results = append(results, relatedUser)
			}
			if status, ok := issueStatusMap[bitbucketIssue.State]; ok {
				bitbucketIssue.State = status
			}
			return results, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

func convertBitbucketIssue(issue *IssuesResponse, connectionId uint64, repositoryId string) (*models.BitbucketIssue, errors.Error) {
	bitbucketIssue := &models.BitbucketIssue{
		ConnectionId:       connectionId,
		BitbucketId:        issue.BitbucketId,
		RepoId:             repositoryId,
		Number:             issue.BitbucketId,
		State:              issue.State,
		Title:              issue.Title,
		Type:               issue.Type,
		Body:               issue.Content.Raw,
		Url:                issue.Links.Self.Href,
		Priority:           issue.Priority,
		Component:          issue.Component,
		BitbucketCreatedAt: issue.BitbucketCreatedAt,
		BitbucketUpdatedAt: issue.BitbucketUpdatedAt,
	}

	if issue.Milestone != nil {
		bitbucketIssue.MilestoneId = issue.Milestone.Id
	}
	if issue.Assignee != nil {
		bitbucketIssue.AssigneeId = issue.Assignee.AccountId
		bitbucketIssue.AssigneeName = issue.Assignee.DisplayName
	}
	if issue.Reporter != nil {
		bitbucketIssue.AuthorId = issue.Reporter.AccountId
		bitbucketIssue.AuthorName = issue.Reporter.DisplayName
	}

	return bitbucketIssue, nil
}

func newIssueStatusMap(config models.TransformationRules) (map[string]string, errors.Error) {
	issueStatusMap := make(map[string]string, 3)
	for _, state := range config.IssueStatusTODO {
		issueStatusMap[state] = ticket.TODO
	}
	for _, state := range config.IssueStatusINPROGRESS {
		issueStatusMap[state] = ticket.IN_PROGRESS
	}
	for _, state := range config.IssueStatusDONE {
		issueStatusMap[state] = ticket.DONE
	}
	return issueStatusMap, nil
}
