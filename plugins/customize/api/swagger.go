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

import "github.com/apache/incubator-devlake/plugins/customize/tasks"

// @Summary blueprints setting for customize
// @Description blueprint setting for customize
// @Tags plugins/customize
// @Accept application/json
// @Param blueprint-setting body blueprintSetting true "json"
// @Router /blueprints/customize/blueprint-setting [post]
func _() {}

type blueprintSetting []struct {
	Version     string `json:"version" example:"1.0.0"`
	Connections []struct {
		Plugin string `json:"plugin" example:"customize"`
		Scope  []struct {
			Options struct {
				Issues []tasks.MappingRules `json:"issues"`
			} `json:"options"`
		} `json:"scope"`
	} `json:"connections"`
}

// @Summary pipelines plan for customize
// @Description pipelines plan for customize
// @Tags plugins/customize
// @Accept application/json
// @Param pipeline-plan body pipelinePlan true "json"
// @Router /pipelines/customize/pipeline-plan [post]
func _() {}

type pipelinePlan [][]struct {
	Plugin   string   `json:"plugin"`
	Subtasks []string `json:"subtasks"`
	Options  struct {
		Issues []tasks.MappingRules `json:"issues"`
	} `json:"options"`
}
