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

// @Summary pipelines plan for github
// @Description pipelines plan for github
// @Tags plugins/github
// @Accept application/json
// @Param pipeline body GithubPipelinePlan true "json"
// @Router /pipelines/github/pipeline-plan [post]
func _() {}

type GithubPipelinePlan [][]struct {
	Plugin   string   `json:"plugin"`
	Subtasks []string `json:"subtasks"`
	Options  struct {
		ConnectionID   int    `json:"connectionId"`
		Owner          string `json:"owner"`
		Repo           string `json:"repo"`
		Since          string
		Transformation CodeTransformationRules `json:"transformation"`
	} `json:"options"`
}

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

// @Summary blueprints setting for github
// @Description blueprint setting for github
// @Tags plugins/github
// @Accept application/json
// @Param blueprint body GithubBlueprintSetting true "json"
// @Router /blueprints/github/blueprint-setting [post]
func _() {}

type GithubBlueprintSetting []struct {
	Version     string `json:"version"`
	Connections []struct {
		Plugin       string `json:"plugin"`
		ConnectionID int    `json:"connectionId"`
		Scope        []struct {
			Transformation CodeTransformationRules `json:"transformation"`
			Options        struct {
				Owner string `json:"owner"`
				Repo  string `json:"repo"`
				Since string
			} `json:"options"`
			Entities []string `json:"entities"`
		} `json:"scope"`
	} `json:"connections"`
}
