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

type CodeTransformationRules struct {
	PrType               string `mapstructure:"prType" json:"prType"`
	PrComponent          string `mapstructure:"prComponent" json:"prComponent"`
	PrBodyClosePattern   string `mapstructure:"prBodyClosePattern" json:"prBodyClosePattern"`
	IssueSeverity        string `mapstructure:"issueSeverity" json:"issueSeverity"`
	IssuePriority        string `mapstructure:"issuePriority" json:"issuePriority"`
	IssueComponent       string `mapstructure:"issueComponent" json:"issueComponent"`
	IssueTypeBug         string `mapstructure:"issueTypeBug" json:"issueTypeBug"`
	IssueTypeIncident    string `mapstructure:"issueTypeIncident" json:"issueTypeIncident"`
	IssueTypeRequirement string `mapstructure:"issueTypeRequirement" json:"issueTypeRequirement"`
}

type TicketTransformationRules struct {
	EpicKeyField               string `json:"epicKeyField"`
	StoryPointField            string `json:"storyPointField"`
	RemotelinkCommitShaPattern string `json:"remotelinkCommitShaPattern"`
	TypeMappings               map[string]struct {
		StandardType string `json:"standardType"`
	} `json:"typeMappings"`
}

// @Summary blueprints setting for tapd
// @Description blueprint setting for tapd
// @Tags plugins/tapd
// @Accept application/json
// @Param blueprint body TapdBlueprintSetting true "json"
// @Router /blueprints/tapd/blueprint-setting [post]
func _() {}

type TapdBlueprintSetting []struct {
	Version     string `json:"version"`
	Connections []struct {
		Plugin       string `json:"plugin"`
		ConnectionID int    `json:"connectionId"`
		Scope        []struct {
			Options struct {
				WorkspaceId uint64   `mapstruct:"workspaceId"`
				CompanyId   uint64   `mapstruct:"companyId"`
				Tasks       []string `mapstruct:"tasks,omitempty"`
				Since       string
			} `json:"options"`
			Entities []string `json:"entities"`
		} `json:"scope"`
	} `json:"connections"`
}

// @Summary pipelines plan for tapd
// @Description pipelines plan for tapd
// @Tags plugins/tapd
// @Accept application/json
// @Param pipeline body TapdPipelinePlan true "json"
// @Router /pipelines/tapd/pipeline-plan [post]
func _() {}

type TapdPipelinePlan [][]struct {
	Plugin   string   `json:"plugin"`
	Subtasks []string `json:"subtasks"`
	Options  struct {
		WorkspaceId uint64   `mapstruct:"workspaceId"`
		CompanyId   uint64   `mapstruct:"companyId"`
		Tasks       []string `mapstruct:"tasks,omitempty"`
		Since       string
	} `json:"options"`
}
