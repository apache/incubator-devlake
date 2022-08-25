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

// @Summary blueprints plan for gitextractor
// @Description blueprints plan for gitextractor
// @Tags plugins/gitextractor
// @Accept application/json
// @Param blueprint body GitextractorBlueprintPlan true "json"
// @Router /blueprints/gitextractor/blueprint-plan [post]
func _() {}

type GitextractorBlueprintPlan [][]struct {
	Plugin  string `json:"plugin"`
	Options struct {
		URL    string `json:"url"`
		RepoID string `json:"repoId"`
	} `json:"options"`
}

// @Summary pipelines plan for gitextractor
// @Description pipelines plan for gitextractor
// @Tags plugins/gitextractor
// @Accept application/json
// @Param pipeline body GitextractorPipelinePlan true "json"
// @Router /pipelines/gitextractor/pipeline-plan [post]
func _() {}

type GitextractorPipelinePlan [][]struct {
	Plugin  string `json:"plugin"`
	Options struct {
		URL    string `json:"url"`
		RepoID string `json:"repoId"`
	} `json:"options"`
}
