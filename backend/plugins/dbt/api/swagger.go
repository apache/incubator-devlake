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

// @Summary blueprints plan for dbt
// @Description blueprints plan for dbt
// @Tags plugins/dbt
// @Accept application/json
// @Param blueprint body DbtBlueprintPlan true "json"
// @Router /blueprints/dbt/blueprint-plan [post]
func _() {}

type Options struct {
	ProjectPath    string   `json:"projectPath"`
	ProjectGitURL  string   `json:"projectGitURL"`
	ProjectName    string   `json:"projectName"`
	ProjectTarget  string   `json:"projectTarget"`
	SelectedModels []string `json:"selectedModels"`
	Args           []string `json:"args"`
	FailFast       bool     `json:"failFast"`
	ProfilesPath   string   `json:"profilesPath"`
	Profile        string   `json:"profile"`
	Threads        int      `json:"threads"`
	NoVersionCheck bool     `json:"noVersionCheck"`
	ExcludeModels  []string `json:"excludeModels"`
	Selector       string   `json:"selector"`
	State          string   `json:"state"`
	Defer          bool     `json:"defer"`
	NoDefer        bool     `json:"noDefer"`
	FullRefresh    bool     `json:"fullRefresh"`
	ProjectVars    struct {
		Demokey1 string `json:"demokey1"`
		Demokey2 string `json:"demokey2"`
	} `json:"projectVars"`
}
type Plan struct {
	Plugin  string  `json:"plugin"`
	Options Options `json:"options"`
}
type DbtBlueprintPlan [][]Plan

// @Summary pipelines plan for dbt
// @Description pipelines plan for dbt
// @Tags plugins/dbt
// @Accept application/json
// @Param pipeline body DbtPipelinePlan true "json"
// @Router /pipelines/dbt/pipeline-plan [post]
func _() {}

type DbtPipelinePlan [][]Plan
