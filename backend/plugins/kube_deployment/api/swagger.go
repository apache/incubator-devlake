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

// @Summary blueprints setting for KubeDeployment
// @Description blueprint setting for KubeDeployment
// @Tags plugins/crypto-asset
// @Accept application/json
// @Param blueprint body KubeDeploymentBlueprintSetting true "json"
// @Router /blueprints/crypto-asset/blueprint-setting [post]
func _() {}

//nolint:unused,deadcode
type KubeDeploymentBlueprintSetting []struct {
	Version     string `json:"version" example:"1.0.0"`
	Connections []struct {
		Plugin string `json:"plugin" example:"KubeDeployment"`
	} `json:"connections"`
}

// @Summary pipelines plan for KubeDeployment
// @Description pipelines plan for KubeDeployment
// @Tags plugins/crypto-asset
// @Accept application/json
// @Param pipeline body KubeDeploymentPipelinePlan true "json"
// @Router /pipelines/crypto-asset/pipeline-plan [post]
func _() {}

//nolint:unused,deadcode
type KubeDeploymentPipelinePlan [][]struct {
	Plugin string `json:"plugin" example:"KubeDeployment"`
}
