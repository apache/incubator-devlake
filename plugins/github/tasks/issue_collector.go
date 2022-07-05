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
	"fmt"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"net/http"
	"net/url"
	"regexp"
	"runtime/debug"
	"time"

	"github.com/apache/incubator-devlake/plugins/helper"

	"github.com/apache/incubator-devlake/plugins/core"
)

const RAW_ISSUE_TABLE = "github_api_issues"

// this struct should be moved to `gitub_api_common.go`
type GithubApiParams struct {
	ConnectionId uint64
	Owner        string
	Repo         string
}

var CollectApiIssuesMeta = core.SubTaskMeta{
	Name:             "collectApiIssues",
	EntryPoint:       CollectApiIssues,
	EnabledByDefault: true,
	Description:      "Collect issues data from Github api",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func CollectApiIssues(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GithubTaskData)
	config := data.Options.TransformationRules
	var issueSeverityRegex *regexp.Regexp
	var issueComponentRegex *regexp.Regexp
	var issuePriorityRegex *regexp.Regexp
	var issueTypeBugRegex *regexp.Regexp
	var issueTypeRequirementRegex *regexp.Regexp
	var issueTypeIncidentRegex *regexp.Regexp
	var issueSeverity = config.IssueSeverity
	var err error
	if len(issueSeverity) > 0 {
		issueSeverityRegex, err = regexp.Compile(issueSeverity)
		if err != nil {
			return fmt.Errorf("regexp Compile issueSeverity failed:[%s] stack:[%s]", err.Error(), debug.Stack())
		}
	}
	var issueComponent = config.IssueComponent
	if len(issueComponent) > 0 {
		issueComponentRegex, err = regexp.Compile(issueComponent)
		if err != nil {
			return fmt.Errorf("regexp Compile issueComponent failed:[%s] stack:[%s]", err.Error(), debug.Stack())
		}
	}
	var issuePriority = config.IssuePriority
	if len(issuePriority) > 0 {
		issuePriorityRegex, err = regexp.Compile(issuePriority)
		if err != nil {
			return fmt.Errorf("regexp Compile issuePriority failed:[%s] stack:[%s]", err.Error(), debug.Stack())
		}
	}
	var issueTypeBug = config.IssueTypeBug
	if len(issueTypeBug) > 0 {
		issueTypeBugRegex, err = regexp.Compile(issueTypeBug)
		if err != nil {
			return fmt.Errorf("regexp Compile issueTypeBug failed:[%s] stack:[%s]", err.Error(), debug.Stack())
		}
	}
	var issueTypeRequirement = config.IssueTypeRequirement
	if len(issueTypeRequirement) > 0 {
		issueTypeRequirementRegex, err = regexp.Compile(issueTypeRequirement)
		if err != nil {
			return fmt.Errorf("regexp Compile issueTypeRequirement failed:[%s] stack:[%s]", err.Error(), debug.Stack())
		}
	}
	var issueTypeIncident = config.IssueTypeIncident
	if len(issueTypeIncident) > 0 {
		issueTypeIncidentRegex, err = regexp.Compile(issueTypeIncident)
		if err != nil {
			return fmt.Errorf("regexp Compile issueTypeIncident failed:[%s] stack:[%s]", err.Error(), debug.Stack())
		}
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			/*
				Table store raw data
			*/
			Table: RAW_ISSUE_TABLE,
		},
		ApiClient: data.ApiClient,
		PageSize:  100,
		/*
			url may use arbitrary variables from different source in any order, we need GoTemplate to allow more
			flexible for all kinds of possibility.
			Pager contains information for a particular page, calculated by ApiCollector, and will be passed into
			GoTemplate to generate a url for that page.
			We want to do page-fetching in ApiCollector, because the logic are highly similar, by doing so, we can
			avoid duplicate logic for every tasks, and when we have a better idea like improving performance, we can
			do it in one place
		*/
		UrlTemplate: "repos/{{ .Params.Owner }}/{{ .Params.Repo }}/issues",
		/*
			(Optional) Return query string for request, or you can plug them into UrlTemplate directly
		*/
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("state", "all")
			if reqData.Since != nil {
				query.Set("since", reqData.Since.String())
			}
			query.Set("direction", "asc")
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))

			return query, nil
		},
		/*
			Some api might do pagination by http headers
		*/
		//Header: func(pager *core.Pager) http.Header {
		//},
		/*
			Sometimes, we need to collect data based on previous collected data, like jira changelog, it requires
			issue_id as part of the url.
			We can mimic `stdin` design, to accept a `Input` function which produces a `Iterator`, collector
			should iterate all records, and do data-fetching for each on, either in parallel or sequential order
			UrlTemplate: "api/3/issue/{{ Input.ID }}/changelog"
		*/
		//Input: databaseIssuesIterator,
		/*
			For api endpoint that returns number of total pages, ApiCollector can collect pages in parallel with ease,
			or other techniques are required if this information was missing.
		*/
		GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var items []json.RawMessage
			err := helper.UnmarshalResponse(res, &items)
			if err != nil {
				return nil, err
			}
			return items, nil
		},
		Extract: func(row *helper.RawData) ([]interface{}, *time.Time, error) {
			body := &IssuesResponse{}
			err := json.Unmarshal(row.Data, body)
			if err != nil {
				return nil, nil, err
			}
			// need to extract 2 kinds of entities here
			if body.GithubId == 0 {
				return nil, nil, nil
			}
			//If this is a pr, ignore
			if body.PullRequest.Url != "" {
				return nil, nil, nil
			}
			results := make([]interface{}, 0, 2)
			githubIssue, err := convertGithubIssue(body, data.Options.ConnectionId, data.Repo.GithubId)
			if err != nil {
				return nil, nil, err
			}
			if body.Assignee != nil {
				githubIssue.AssigneeId = body.Assignee.Id
				githubIssue.AssigneeName = body.Assignee.Login
				relatedUser, err := convertUser(body.Assignee, data.Options.ConnectionId)
				if err != nil {
					return nil, nil, err
				}
				results = append(results, relatedUser)
			}
			if body.User != nil {
				githubIssue.AuthorId = body.User.Id
				githubIssue.AuthorName = body.User.Login
				relatedUser, err := convertUser(body.User, data.Options.ConnectionId)
				if err != nil {
					return nil, nil, err
				}
				results = append(results, relatedUser)
			}
			for _, label := range body.Labels {
				results = append(results, &models.GithubIssueLabel{
					ConnectionId: data.Options.ConnectionId,
					IssueId:      githubIssue.GithubId,
					LabelName:    label.Name,
				})
				if issueSeverityRegex != nil {
					groups := issueSeverityRegex.FindStringSubmatch(label.Name)
					if len(groups) > 0 {
						githubIssue.Severity = groups[1]
					}
				}

				if issueComponentRegex != nil {
					groups := issueComponentRegex.FindStringSubmatch(label.Name)
					if len(groups) > 0 {
						githubIssue.Component = groups[1]
					}
				}

				if issuePriorityRegex != nil {
					groups := issuePriorityRegex.FindStringSubmatch(label.Name)
					if len(groups) > 0 {
						githubIssue.Priority = groups[1]
					}
				}

				if issueTypeBugRegex != nil {
					if ok := issueTypeBugRegex.MatchString(label.Name); ok {
						githubIssue.Type = ticket.BUG
					}
				}

				if issueTypeRequirementRegex != nil {
					if ok := issueTypeRequirementRegex.MatchString(label.Name); ok {
						githubIssue.Type = ticket.REQUIREMENT
					}
				}

				if issueTypeIncidentRegex != nil {
					if ok := issueTypeIncidentRegex.MatchString(label.Name); ok {
						githubIssue.Type = ticket.INCIDENT
					}
				}
			}
			results = append(results, githubIssue)

			return results, &githubIssue.GithubUpdatedAt, nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
