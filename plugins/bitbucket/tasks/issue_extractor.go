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
	"time"

	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type IssuesResponse struct {
	Type        string `json:"type"`
	BitbucketId int    `json:"id"`
	Number      int
	Repository  *BitbucketApiRepo
	Links       struct {
		Self struct {
			Href string
		} `json:"self"`
		Html struct {
			Href string
		} `json:"html"`
	} `json:"links"`
	Title     string `json:"title"`
	Reporter  *BitbucketAccountResponse
	Assignee  *BitbucketAccountResponse
	State     string
	Kind      string
	Milestone *struct {
		Id int
	} `json:"milestone"`
	Votes              int       `json:"votes"`
	BitbucketCreatedAt time.Time `json:"created_on"`
	BitbucketUpdatedAt time.Time `json:"updated_on"`
}

type IssueRegexes struct {
	SeverityRegex        *regexp.Regexp
	ComponentRegex       *regexp.Regexp
	PriorityRegex        *regexp.Regexp
	TypeBugRegex         *regexp.Regexp
	TypeRequirementRegex *regexp.Regexp
	TypeIncidentRegex    *regexp.Regexp
}

var ExtractApiIssuesMeta = core.SubTaskMeta{
	Name:             "extractApiIssues",
	EntryPoint:       ExtractApiIssues,
	EnabledByDefault: true,
	Description:      "Extract raw Issues data into tool layer table bitbucket_issues",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ExtractApiIssues(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*BitbucketTaskData)

	//config := data.Options.TransformationRules
	//issueRegexes, err := NewIssueRegexes(config)
	//if err != nil {
	//	return nil
	//}
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
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			body := &IssuesResponse{}
			err := json.Unmarshal(row.Data, body)
			if err != nil {
				return nil, err
			}
			// need to extract 2 kinds of entities here
			if body.BitbucketId == 0 {
				return nil, nil
			}
			//If this is a pr, ignore
			if body.Type != "issue" {
				return nil, nil
			}
			results := make([]interface{}, 0, 2)

			bitbucketIssue, err := convertBitbucketIssue(body, data.Options.ConnectionId, data.Repo.BitbucketId)
			if err != nil {
				return nil, err
			}
			//bitbucketLabels, err := convertBitbucketLabels(issueRegexes, body, bitbucketIssue)
			//if err != nil {
			//	return nil, err
			//}
			//results = append(results, bitbucketLabels...)

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
			return results, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

func convertBitbucketIssue(issue *IssuesResponse, connectionId uint64, repositoryId string) (*models.BitbucketIssue, error) {
	bitbucketIssue := &models.BitbucketIssue{
		ConnectionId:       connectionId,
		BitbucketId:        issue.BitbucketId,
		RepoId:             repositoryId,
		Number:             issue.Number,
		State:              issue.State,
		Title:              issue.Title,
		Url:                issue.Links.Self.Href,
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

//func convertBitbucketLabels(issueRegexes *IssueRegexes, issue *IssuesResponse, bitbucketIssue *models.BitbucketIssue) ([]interface{}, error) {
//	var results []interface{}
//	for _, label := range issue {
//		results = append(results, &models.BitbucketIssueLabel{
//			ConnectionId: bitbucketIssue.ConnectionId,
//			IssueId:      bitbucketIssue.BitbucketId,
//			LabelName:    label,
//		})
//		if issueRegexes.SeverityRegex != nil {
//			groups := issueRegexes.SeverityRegex.FindStringSubmatch(label.Name)
//			if len(groups) > 0 {
//				bitbucketIssue.Severity = groups[1]
//			}
//		}
//		if issueRegexes.ComponentRegex != nil {
//			groups := issueRegexes.ComponentRegex.FindStringSubmatch(label.Name)
//			if len(groups) > 0 {
//				bitbucketIssue.Component = groups[1]
//			}
//		}
//		if issueRegexes.PriorityRegex != nil {
//			groups := issueRegexes.PriorityRegex.FindStringSubmatch(label.Name)
//			if len(groups) > 0 {
//				bitbucketIssue.Priority = groups[1]
//			}
//		}
//		if issueRegexes.TypeBugRegex != nil {
//			if ok := issueRegexes.TypeBugRegex.MatchString(label.Name); ok {
//				bitbucketIssue.Type = ticket.BUG
//			}
//		}
//		if issueRegexes.TypeRequirementRegex != nil {
//			if ok := issueRegexes.TypeRequirementRegex.MatchString(label.Name); ok {
//				bitbucketIssue.Type = ticket.REQUIREMENT
//			}
//		}
//		if issueRegexes.TypeIncidentRegex != nil {
//			if ok := issueRegexes.TypeIncidentRegex.MatchString(label.Name); ok {
//				bitbucketIssue.Type = ticket.INCIDENT
//			}
//		}
//	}
//	return results, nil
//}

//func NewIssueRegexes(config models.TransformationRules) (*IssueRegexes, error) {
//	var issueRegexes IssueRegexes
//	var issueSeverity = config.IssueSeverity
//	var err error
//	if len(issueSeverity) > 0 {
//		issueRegexes.SeverityRegex, err = regexp.Compile(issueSeverity)
//		if err != nil {
//			return nil, fmt.Errorf("regexp Compile issueSeverity failed:[%s] stack:[%s]", err.Error(), debug.Stack())
//		}
//	}
//	var issueComponent = config.IssueComponent
//	if len(issueComponent) > 0 {
//		issueRegexes.ComponentRegex, err = regexp.Compile(issueComponent)
//		if err != nil {
//			return nil, fmt.Errorf("regexp Compile issueComponent failed:[%s] stack:[%s]", err.Error(), debug.Stack())
//		}
//	}
//	var issuePriority = config.IssuePriority
//	if len(issuePriority) > 0 {
//		issueRegexes.PriorityRegex, err = regexp.Compile(issuePriority)
//		if err != nil {
//			return nil, fmt.Errorf("regexp Compile issuePriority failed:[%s] stack:[%s]", err.Error(), debug.Stack())
//		}
//	}
//	var issueTypeBug = config.IssueTypeBug
//	if len(issueTypeBug) > 0 {
//		issueRegexes.TypeBugRegex, err = regexp.Compile(issueTypeBug)
//		if err != nil {
//			return nil, fmt.Errorf("regexp Compile issueTypeBug failed:[%s] stack:[%s]", err.Error(), debug.Stack())
//		}
//	}
//	var issueTypeRequirement = config.IssueTypeRequirement
//	if len(issueTypeRequirement) > 0 {
//		issueRegexes.TypeRequirementRegex, err = regexp.Compile(issueTypeRequirement)
//		if err != nil {
//			return nil, fmt.Errorf("regexp Compile issueTypeRequirement failed:[%s] stack:[%s]", err.Error(), debug.Stack())
//		}
//	}
//	var issueTypeIncident = config.IssueTypeIncident
//	if len(issueTypeIncident) > 0 {
//		issueRegexes.TypeIncidentRegex, err = regexp.Compile(issueTypeIncident)
//		if err != nil {
//			return nil, fmt.Errorf("regexp Compile issueTypeIncident failed:[%s] stack:[%s]", err.Error(), debug.Stack())
//		}
//	}
//	return &issueRegexes, nil
//}
