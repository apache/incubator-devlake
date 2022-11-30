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

import "github.com/apache/incubator-devlake/plugins/jira/tasks"

// @Summary blueprints setting for jira
// @Description blueprint setting for jira
// @Tags plugins/jira
// @Accept application/json
// @Param blueprint-setting body JiraBlueprintSetting true "json"
// @Router /blueprints/jira/blueprint-setting [post]
func _() {}

type JiraBlueprintSetting []struct {
	Version     string `json:"version"`
	Connections []struct {
		Plugin       string `json:"plugin"`
		ConnectionID int    `json:"connectionId"`
		Scope        []struct {
			Transformation tasks.JiraTransformationRule `json:"transformation"`
			Options        struct {
				BoardId uint64 `json:"boardId"`
				Since   string `json:"since"`
			} `json:"options"`
			Entities []string `json:"entities"`
		} `json:"scope"`
	} `json:"connections"`
}

// @Summary pipelines plan for jira
// @Description pipelines plan for jira
// @Tags plugins/jira
// @Accept application/json
// @Param pipeline-plan body JiraPipelinePlan true "json"
// @Router /pipelines/jira/pipeline-plan [post]
func _() {}

type JiraPipelinePlan [][]struct {
	Plugin   string   `json:"plugin"`
	Subtasks []string `json:"subtasks"`
	Options  struct {
		BoardID             int                          `json:"boardId"`
		ConnectionID        int                          `json:"connectionId"`
		TransformationRules tasks.JiraTransformationRule `json:"transformationRules"`
	} `json:"options"`
}
