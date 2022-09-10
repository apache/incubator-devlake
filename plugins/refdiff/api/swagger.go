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

// @Summary blueprints plan for refdiff
// @Description blueprints plan for refdiff
// @Tags plugins/refdiff
// @Accept application/json
// @Param blueprint body RefdiffBlueprintPlan true "json"
// @Router /blueprints/refdiff/blueprint-plan [post]
func _() {}

//nolint:unused
type RefdiffBlueprintPlan [][]struct {
	Plugin  string `json:"plugin"`
	Options struct {
		RepoID string `json:"repoId"`
		Pairs  []struct {
			NewRef string `json:"newRef"`
			OldRef string `json:"oldRef"`
		} `json:"pairs"`
	} `json:"options"`
}
