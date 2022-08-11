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

package api

type plan struct {
	Plugin   string   `json:"plugin" example:"github"`
	Subtasks []string `json:"subtasks" example:"collectApiIssues,extractApiIssues,collectApiPullRequests,extractApiPullRequests,collectApiComments,extractApiComments,collectApiEvents,extractApiEvents,collectApiPullRequestCommits,extractApiPullRequestCommits,collectApiPullRequestReviews,extractApiPullRequestReviews,CollectApiPrReviewCommentsMeta,extractApiPrReviewComments,collectApiMilestones,extractMilestones,collectAccounts,extractAccounts,collectAccountOrg,ExtractAccountOrg,enrichPullRequestIssues,convertRepo,convertIssues,convertIssueLabels,convertPullRequestCommits,convertPullRequests,convertPullRequestReviews,convertPullRequestLabels,convertPullRequestIssues,convertIssueComments,convertPullRequestComments,convertMilestones,convertAccounts"`
	Options  struct {
		ConnectionID        int    `json:"connectionId" example:"1"`
		Owner               string `json:"owner" example:"apache"`
		Repo                string `json:"repo" example:"incubator-devlake"`
		TransformationRules struct {
			PrType               string `json:"prType,omitempty" example:"type/(.*)$"`
			PrComponent          string `json:"prComponent,omitempty" example:"component/(.*)$"`
			IssueSeverity        string `json:"issueSeverity,omitempty" example:"severity/(.*)$"`
			IssueComponent       string `json:"issueComponent,omitempty" example:"component/(.*)$"`
			IssuePriority        string `json:"issuePriority,omitempty" example:"(highest|high|medium|low)$"`
			IssueTypeRequirement string `json:"issueTypeRequirement,omitempty" example:"(feat|feature|proposal|requirement)$"`
			IssueTypeBug         string `json:"issueTypeBug,omitempty" example:"(bug|broken)$"`
			IssueTypeIncident    string `json:"issueTypeIncident,omitempty" example:"(incident|p0|p1|p2)$"`
		} `json:"transformationRules,omitempty"`
	} `json:"options,omitempty"`
}

type blueprintOutput struct {
	ID   int      `json:"id" example:"14"`
	Plan [][]plan `json:"plan"`
	blueprintInput
}
type blueprintInput struct {
	Name     string `json:"name"`
	Settings struct {
		Version     string `json:"version" example:"1.0.0"`
		Connections []struct {
			Plugin       string `json:"plugin" example:"github"`
			ConnectionID int    `json:"connectionId" example:"1"`
			Scope        []struct {
				Transformation struct {
					PrType               string `json:"prType" example:"type/(.*)$"`
					PrComponent          string `json:"prComponent" example:"component/(.*)$"`
					IssueSeverity        string `json:"issueSeverity" example:"severity/(.*)$"`
					IssueComponent       string `json:"issueComponent" example:"component/(.*)$"`
					IssuePriority        string `json:"issuePriority" example:"(highest|high|medium|low)$"`
					IssueTypeRequirement string `json:"issueTypeRequirement" example:"(feat|feature|proposal|requirement)$"`
					IssueTypeBug         string `json:"issueTypeBug" example:"(bug|broken)$"`
					IssueTypeIncident    string `json:"issueTypeIncident" example:"(incident|p0|p1|p2)$"`
					Refdiff              struct {
						TagsOrder   string `json:"tagsOrder" example:"reverse semver"`
						TagsPattern string `json:"tagsPattern" example:"(regex)$"`
						TagsLimit   int    `json:"tagsLimit" example:"10"`
					} `json:"refdiff"`
				} `json:"transformation"`
				Options struct {
					Owner string `json:"owner" example:"apache"`
					Repo  string `json:"repo" example:"incubator-devlake"`
				} `json:"options"`
				Entities []string `json:"entities" example:"CODE,TICKET,CODEREVIEW,CROSS"`
			} `json:"scope"`
		} `json:"connections"`
	} `json:"settings"`
	CronConfig string `json:"cronConfig" example:"0 0 * * *"`
	Enable     bool   `json:"enable" example:"true"`
	Mode       string `json:"mode" example:"NORMAL"`
	IsManual   bool   `json:"isManual" example:"true"`
}
