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

// @Summary blueprints setting for sonarqube
// @Description blueprint setting for sonarqube
// @Tags plugins/sonarqube
// @Accept application/json
// @Param blueprint-setting body SonarqubeBlueprintSetting true "json"
// @Router /blueprints/sonarqube/blueprint-setting [post]
func _() {}

type SonarqubeBlueprintSetting []struct {
	Version     string `json:"version"`
	Connections []struct {
		Plugin       string `json:"plugin"`
		ConnectionID int    `json:"connectionId"`
		Scope        []struct {
			Options struct {
				ProjectKey string `json:"projectKey"`
			} `json:"options"`
			Entities []string `json:"entities"`
		} `json:"scopes"`
	} `json:"connections"`
}

// @Summary pipelines plan for sonarqube
// @Description pipelines plan for sonarqube
// @Tags plugins/sonarqube
// @Accept application/json
// @Param pipeline-plan body SonarqubePipelinePlan true "json"
// @Router /pipelines/sonarqube/pipeline-plan [post]
func _() {}

type SonarqubePipelinePlan [][]struct {
	Plugin   string   `json:"plugin"`
	Subtasks []string `json:"subtasks"`
	Options  struct {
		ProjectKey   string `json:"projectKey"`
		ConnectionID int    `json:"connectionId"`
	} `json:"options"`
}
