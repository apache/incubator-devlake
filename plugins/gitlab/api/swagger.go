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
	Plugin   string   `json:"plugin" example:"gitlab"`
	Subtasks []string `json:"subtasks" example:"collectApiProject,extractApiProject,collectApiIssues,extractApiIssues,collectApiMergeRequests,extractApiMergeRequests,collectApiMergeRequestsNotes,extractApiMergeRequestsNotes,collectApiMergeRequestsCommits,extractApiMergeRequestsCommits,enrichMrs,collectAccounts,extractAccounts,convertAccounts,convertApiProject,convertApiMergeRequests,convertMergeRequestComment,convertApiMergeRequestsCommits,convertIssues,convertIssueLabels,convertMrLabels,convertPipelines,convertJobs"`
	Options  struct {
		ConnectionID        int      `json:"connectionId" example:"1"`
		ProjectID           string   `json:"projectId" example:"18895622"`
		TransformationRules struct{} `json:"transformationRules"`
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
			Plugin       string `json:"plugin" example:"gitlab"`
			ConnectionID int    `json:"connectionId" example:"1"`
			Scope        []struct {
				Transformation struct{} `json:"transformation"`
				Options        struct {
					ProjectID string `json:"projectId" example:"18895622"`
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
