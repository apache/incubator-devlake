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
	"regexp"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiIssuesMeta)
}

var ExtractApiIssuesMeta = plugin.SubTaskMeta{
	Name:             "Extract Issues",
	EntryPoint:       ExtractApiIssues,
	EnabledByDefault: true,
	Description:      "Extract raw Issues data into tool layer table gitlab_issues",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	Dependencies:     []*plugin.SubTaskMeta{&CollectApiIssuesMeta},
}

type IssuesResponse struct {
	ProjectId int `json:"project_id"`
	Milestone struct {
		Due_date    string
		Project_id  int
		State       string
		Description string
		Iid         int
		Id          int
		Title       string
		CreatedAt   common.Iso8601Time
		UpdatedAt   common.Iso8601Time
	}
	Author *struct {
		State     string
		WebUrl    string `json:"web_url"`
		AvatarUrl string `json:"avatar_url"`
		Username  string
		Id        int    `json:"id"`
		Name      string `json:"name"`
	}
	Description string
	State       string
	Iid         int
	Assignees   []struct {
		AvatarUrl string `json:"avatar_url"`
		WebUrl    string `json:"web_url"`
		State     string
		Username  string
		Id        int
		Name      string
	}
	Assignee *struct {
		AvatarUrl string
		WebUrl    string
		State     string
		Username  string
		Id        int
		Name      string
	}
	Type               string
	Labels             []string `json:"labels"`
	UpVotes            int
	DownVotes          int
	MergeRequestsCount int
	Id                 int `json:"id"`
	Title              string
	GitlabUpdatedAt    common.Iso8601Time  `json:"updated_at"`
	GitlabCreatedAt    common.Iso8601Time  `json:"created_at"`
	GitlabClosedAt     *common.Iso8601Time `json:"closed_at"`
	ClosedBy           struct {
		State     string
		WebUrl    string
		AvatarUrl string
		Username  string
		Id        int
		Name      string
	}
	UserNotesCount int
	DueDate        common.Iso8601Time
	WebUrl         string `json:"web_url"`
	References     struct {
		Short    string
		Relative string
		Full     string
	}
	TimeStats struct {
		TimeEstimate        *int64
		TotalTimeSpent      *int64
		HumanTimeEstimate   string
		HumanTotalTimeSpent string
	}
	HasTasks         bool
	TaskStatus       string
	Confidential     bool
	DiscussionLocked bool
	IssueType        string
	Severity         string
	Component        string
	Priority         string
	Links            struct {
		Self       string `json:"self"`
		Notes      string
		AwardEmoji string
		Project    string
	} `json:"_links"`
	TaskCompletionStatus struct {
		Count          int
		CompletedCount int
	}
}

func ExtractApiIssues(subtaskCtx plugin.SubTaskContext) errors.Error {
	subtaskCommonArgs, data := CreateSubtaskCommonArgs(subtaskCtx, RAW_ISSUE_TABLE)

	db := subtaskCtx.GetDal()
	config := data.Options.ScopeConfig
	var issueSeverityRegex *regexp.Regexp
	var issueComponentRegex *regexp.Regexp
	var issuePriorityRegex *regexp.Regexp

	var issueSeverity = config.IssueSeverity
	var err error
	if len(issueSeverity) > 0 {
		issueSeverityRegex, err = regexp.Compile(issueSeverity)
		if err != nil {
			return errors.Default.Wrap(err, "regexp Compile issueSeverity failed")
		}
	}
	var issueComponent = config.IssueComponent
	if len(issueComponent) > 0 {
		issueComponentRegex, err = regexp.Compile(issueComponent)
		if err != nil {
			return errors.Default.Wrap(err, "regexp Compile issueComponent failed")
		}
	}
	var issuePriority = config.IssuePriority
	if len(issuePriority) > 0 {
		issuePriorityRegex, err = regexp.Compile(issuePriority)
		if err != nil {
			return errors.Default.Wrap(err, "regexp Compile issuePriority failed")
		}
	}
	subtaskCommonArgs.SubtaskConfig = map[string]interface{}{
		"issueSeverity":      issueSeverity,
		"issueComponent":     issueComponent,
		"issuePriorityRegex": issuePriorityRegex,
	}
	extractor, err := api.NewStatefulApiExtractor(&api.StatefulApiExtractorArgs[IssuesResponse]{
		SubtaskCommonArgs: subtaskCommonArgs,
		BeforeExtract: func(body *IssuesResponse, stateManager *api.SubtaskStateManager) errors.Error {
			if stateManager.IsIncremental() {
				err := db.Delete(
					&models.GitlabIssueLabel{},
					dal.Where("connection_id = ? AND issue_id = ?", data.Options.ConnectionId, body.Id),
				)
				if err != nil {
					return err
				}
				err = db.Delete(
					&models.GitlabIssueAssignee{},
					dal.Where("connection_id = ? AND gitlab_id = ?", data.Options.ConnectionId, body.Id),
				)
				if err != nil {
					return err
				}
			}
			return nil
		},
		Extract: func(body *IssuesResponse, row *api.RawData) ([]interface{}, errors.Error) {
			if body.ProjectId == 0 {
				return nil, nil
			}

			results := make([]interface{}, 0, 2)
			gitlabIssue, err := convertGitlabIssue(body, data.Options.ProjectId)
			if err != nil {
				return nil, err
			}
			for _, label := range body.Labels {
				results = append(results, &models.GitlabIssueLabel{
					ConnectionId: data.Options.ConnectionId,
					IssueId:      gitlabIssue.GitlabId,
					LabelName:    label,
				})
				if issueSeverityRegex != nil && issueSeverityRegex.MatchString(label) {
					gitlabIssue.Severity = label
				}
				if issueComponentRegex != nil && issueComponentRegex.MatchString(label) {
					gitlabIssue.Component = label
				}
				if issuePriorityRegex != nil && issuePriorityRegex.MatchString(label) {
					gitlabIssue.Priority = label
				}
			}

			gitlabIssue.ConnectionId = data.Options.ConnectionId
			if body.Author != nil {
				gitlabAuthor, err := convertGitlabAuthor(body, data.Options.ConnectionId)
				if err != nil {
					return nil, err
				}
				results = append(results, gitlabAuthor)
			}
			results = append(results, gitlabIssue)

			for _, v := range body.Assignees {
				assignee := &models.GitlabAccount{
					ConnectionId: data.Options.ConnectionId,
					Username:     v.Username,
					Name:         v.Name,
					State:        v.State,
					AvatarUrl:    v.AvatarUrl,
					WebUrl:       v.WebUrl,
				}
				issueAssignee := &models.GitlabIssueAssignee{
					ConnectionId: data.Options.ConnectionId,
					GitlabId:     gitlabIssue.GitlabId,
					ProjectId:    gitlabIssue.ProjectId,
					AssigneeId:   v.Id,
					AssigneeName: v.Username,
				}
				results = append(results, assignee, issueAssignee)
			}

			return results, nil
		},
	})

	if err != nil {
		return errors.Convert(err)
	}

	return extractor.Execute()
}

func convertGitlabIssue(issue *IssuesResponse, projectId int) (*models.GitlabIssue, errors.Error) {
	gitlabIssue := &models.GitlabIssue{
		GitlabId:        issue.Id,
		ProjectId:       projectId,
		Number:          issue.Iid,
		State:           issue.State,
		Type:            issue.Type,
		StdType:         issue.Type,
		Severity:        issue.Severity,
		Component:       issue.Component,
		Priority:        issue.Priority,
		Title:           issue.Title,
		Body:            issue.Description,
		Url:             issue.WebUrl,
		ClosedAt:        common.Iso8601TimeToTime(issue.GitlabClosedAt),
		GitlabCreatedAt: issue.GitlabCreatedAt.ToTime(),
		GitlabUpdatedAt: issue.GitlabUpdatedAt.ToTime(),
		TimeEstimate:    issue.TimeStats.TimeEstimate,
		TotalTimeSpent:  issue.TimeStats.TotalTimeSpent,
		CreatorId:       issue.Author.Id,
		CreatorName:     issue.Author.Username,
	}

	if issue.Assignee != nil {
		gitlabIssue.AssigneeId = issue.Assignee.Id
		gitlabIssue.AssigneeName = issue.Assignee.Username
	}
	if issue.Author != nil {
		gitlabIssue.CreatorId = issue.Author.Id
		gitlabIssue.CreatorName = issue.Author.Username
	}
	if issue.GitlabClosedAt != nil {
		temp := uint(issue.GitlabClosedAt.ToTime().Sub(issue.GitlabCreatedAt.ToTime()).Minutes())
		gitlabIssue.LeadTimeMinutes = &temp
	}

	return gitlabIssue, nil
}

func convertGitlabAuthor(issue *IssuesResponse, connectionId uint64) (*models.GitlabAccount, errors.Error) {
	gitlabAuthor := &models.GitlabAccount{
		ConnectionId: connectionId,
		GitlabId:     issue.Author.Id,
		Username:     issue.Author.Username,
		Name:         issue.Author.Name,
		State:        issue.Author.State,
		AvatarUrl:    issue.Author.AvatarUrl,
		WebUrl:       issue.Author.WebUrl,
	}

	return gitlabAuthor, nil
}
