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
	"regexp"
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiIssuesMeta)
}

var ExtractApiIssuesMeta = plugin.SubTaskMeta{
	Name:             "extractApiIssues",
	EntryPoint:       ExtractApiIssues,
	EnabledByDefault: true,
	Description:      "Extract raw Issues data into tool layer table github_issues",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	DependencyTables: []string{RAW_ISSUE_TABLE},
	ProductTables: []string{
		models.GithubIssue{}.TableName(),
		models.GithubIssueLabel{}.TableName(),
		models.GithubRepoAccount{}.TableName(),
		models.GithubIssueAssignee{}.TableName()},
}

type IssuesResponse struct {
	GithubId    int `json:"id"`
	Number      int
	State       string
	Title       string
	Body        json.RawMessage
	HtmlUrl     string `json:"html_url"`
	PullRequest struct {
		Url     string `json:"url"`
		HtmlUrl string `json:"html_url"`
	} `json:"pull_request"`
	Labels []struct {
		Name string `json:"name"`
	} `json:"labels"`
	Assignee  *GithubAccountResponse  `json:"assignee"`
	Assignees []GithubAccountResponse `json:"assignees"`
	User      *GithubAccountResponse
	Milestone *struct {
		Id int
	}
	ClosedAt        *api.Iso8601Time `json:"closed_at"`
	GithubCreatedAt api.Iso8601Time  `json:"created_at"`
	GithubUpdatedAt api.Iso8601Time  `json:"updated_at"`
}

type IssueRegexes struct {
	SeverityRegex        *regexp.Regexp
	ComponentRegex       *regexp.Regexp
	PriorityRegex        *regexp.Regexp
	TypeBugRegex         *regexp.Regexp
	TypeRequirementRegex *regexp.Regexp
	TypeIncidentRegex    *regexp.Regexp
}

func ExtractApiIssues(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*GithubTaskData)

	config := data.Options.ScopeConfig
	issueRegexes, err := NewIssueRegexes(config)
	if err != nil {
		return nil
	}
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.Name,
			},
			/*
				Table store raw data
			*/
			Table: RAW_ISSUE_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			body := &IssuesResponse{}
			err := errors.Convert(json.Unmarshal(row.Data, body))
			if err != nil {
				return nil, err
			}
			// need to extract 2 kinds of entities here
			if body.GithubId == 0 {
				return nil, nil
			}
			//If this is a pr, ignore
			if body.PullRequest.Url != "" {
				return nil, nil
			}
			results := make([]interface{}, 0, 2)

			githubIssue, err := convertGithubIssue(body, data.Options.ConnectionId, data.Options.GithubId)
			if err != nil {
				return nil, err
			}
			githubLabels, err := convertGithubLabels(issueRegexes, body, githubIssue)
			if err != nil {
				return nil, err
			}
			results = append(results, githubLabels...)
			results = append(results, githubIssue)
			if body.Assignee != nil {
				githubIssue.AssigneeId = body.Assignee.Id
				githubIssue.AssigneeName = body.Assignee.Login
				relatedUser, err := convertAccount(body.Assignee, data.Options.GithubId, data.Options.ConnectionId)
				if err != nil {
					return nil, err
				}
				results = append(results, relatedUser)
			}
			for _, assignee := range body.Assignees {
				issueAssignee := &models.GithubIssueAssignee{
					ConnectionId: githubIssue.ConnectionId,
					IssueId:      githubIssue.GithubId,
					RepoId:       githubIssue.RepoId,
					AssigneeId:   assignee.Id,
					AssigneeName: assignee.Login,
				}
				results = append(results, issueAssignee)
			}
			if body.User != nil {
				githubIssue.AuthorId = body.User.Id
				githubIssue.AuthorName = body.User.Login
				relatedUser, err := convertAccount(body.User, data.Options.GithubId, data.Options.ConnectionId)
				if err != nil {
					return nil, err
				}
				results = append(results, relatedUser)
			}
			return results, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

func convertGithubIssue(issue *IssuesResponse, connectionId uint64, repositoryId int) (*models.GithubIssue, errors.Error) {
	githubIssue := &models.GithubIssue{
		ConnectionId:    connectionId,
		GithubId:        issue.GithubId,
		RepoId:          repositoryId,
		Number:          issue.Number,
		State:           issue.State,
		Title:           issue.Title,
		Body:            string(issue.Body),
		Url:             issue.HtmlUrl,
		ClosedAt:        api.Iso8601TimeToTime(issue.ClosedAt),
		GithubCreatedAt: issue.GithubCreatedAt.ToTime(),
		GithubUpdatedAt: issue.GithubUpdatedAt.ToTime(),
	}
	if issue.Assignee != nil {
		githubIssue.AssigneeId = issue.Assignee.Id
		githubIssue.AssigneeName = issue.Assignee.Login
	}
	if issue.User != nil {
		githubIssue.AuthorId = issue.User.Id
		githubIssue.AuthorName = issue.User.Login
	}
	if issue.ClosedAt != nil {
		githubIssue.LeadTimeMinutes = uint(issue.ClosedAt.ToTime().Sub(issue.GithubCreatedAt.ToTime()).Minutes())
	}
	if issue.Milestone != nil {
		githubIssue.MilestoneId = issue.Milestone.Id
	}
	return githubIssue, nil
}

func convertGithubLabels(issueRegexes *IssueRegexes, issue *IssuesResponse, githubIssue *models.GithubIssue) ([]interface{}, errors.Error) {
	var results []interface{}
	var joinedLabels []string
	for _, label := range issue.Labels {
		results = append(results, &models.GithubIssueLabel{
			ConnectionId: githubIssue.ConnectionId,
			IssueId:      githubIssue.GithubId,
			LabelName:    label.Name,
		})
		joinedLabels = append(joinedLabels, label.Name)

		if issueRegexes.SeverityRegex != nil {
			groups := issueRegexes.SeverityRegex.FindStringSubmatch(label.Name)
			if len(groups) > 1 {
				githubIssue.Severity = groups[1]
			}
		}
		if issueRegexes.ComponentRegex != nil {
			groups := issueRegexes.ComponentRegex.FindStringSubmatch(label.Name)
			if len(groups) > 1 {
				githubIssue.Component = groups[1]
			}
		}
		if issueRegexes.PriorityRegex != nil {
			groups := issueRegexes.PriorityRegex.FindStringSubmatch(label.Name)
			if len(groups) > 1 {
				githubIssue.Priority = groups[1]
			}
		}
		if issueRegexes.TypeRequirementRegex != nil && issueRegexes.TypeRequirementRegex.MatchString(label.Name) {
			githubIssue.StdType = ticket.REQUIREMENT
		} else if issueRegexes.TypeBugRegex != nil && issueRegexes.TypeBugRegex.MatchString(label.Name) {
			githubIssue.StdType = ticket.BUG
		} else if issueRegexes.TypeIncidentRegex != nil && issueRegexes.TypeIncidentRegex.MatchString(label.Name) {
			githubIssue.StdType = ticket.INCIDENT
		}
	}
	if len(joinedLabels) > 0 {
		githubIssue.Type = strings.Join(joinedLabels, ",")
	}
	return results, nil
}

func NewIssueRegexes(config *models.GithubScopeConfig) (*IssueRegexes, errors.Error) {
	var issueRegexes IssueRegexes
	if config == nil {
		return &issueRegexes, nil
	}
	var issueSeverity = config.IssueSeverity
	var err error
	if len(issueSeverity) > 0 {
		issueRegexes.SeverityRegex, err = regexp.Compile(issueSeverity)
		if err != nil {
			return nil, errors.Default.Wrap(err, "regexp Compile issueSeverity failed")
		}
	}
	var issueComponent = config.IssueComponent
	if len(issueComponent) > 0 {
		issueRegexes.ComponentRegex, err = regexp.Compile(issueComponent)
		if err != nil {
			return nil, errors.Default.Wrap(err, "regexp Compile issueComponent failed")
		}
	}
	var issuePriority = config.IssuePriority
	if len(issuePriority) > 0 {
		issueRegexes.PriorityRegex, err = regexp.Compile(issuePriority)
		if err != nil {
			return nil, errors.Default.Wrap(err, "regexp Compile issuePriority failed")
		}
	}
	var issueTypeBug = config.IssueTypeBug
	if len(issueTypeBug) > 0 {
		issueRegexes.TypeBugRegex, err = regexp.Compile(issueTypeBug)
		if err != nil {
			return nil, errors.Default.Wrap(err, "regexp Compile issueTypeBug failed")
		}
	}
	var issueTypeRequirement = config.IssueTypeRequirement
	if len(issueTypeRequirement) > 0 {
		issueRegexes.TypeRequirementRegex, err = regexp.Compile(issueTypeRequirement)
		if err != nil {
			return nil, errors.Default.Wrap(err, "regexp Compile issueTypeRequirement failed")
		}
	}
	var issueTypeIncident = config.IssueTypeIncident
	if len(issueTypeIncident) > 0 {
		issueRegexes.TypeIncidentRegex, err = regexp.Compile(issueTypeIncident)
		if err != nil {
			return nil, errors.Default.Wrap(err, "regexp Compile issueTypeIncident failed")
		}
	}
	return &issueRegexes, nil
}
