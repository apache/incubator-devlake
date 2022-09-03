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

package core

import (
	"encoding/json"
	"github.com/apache/incubator-devlake/errors"
)

// PluginBlueprint is used to support Blueprint Normal model
type PluginBlueprintV100 interface {
	// MakePipelinePlan generates `pipeline.tasks` based on `version` and `scope`
	//
	// `version` semver from `blueprint.settings.version`
	// `scope` arbitrary json.RawMessage, depends on `version`, for v0.0.1, it is an Array of Objects
	MakePipelinePlan(connectionId uint64, scope []*BlueprintScopeV100) (PipelinePlan, errors.Error)
}

// BlueprintConnectionV100 is the connection definition for protocol v1.0.0
type BlueprintConnectionV100 struct {
	Plugin       string                `json:"plugin" validate:"required"`
	ConnectionId uint64                `json:"connectionId" validate:"required"`
	Scope        []*BlueprintScopeV100 `json:"scope" validate:"required"`
}

// BlueprintScopeV100 is the scope definition for protocol v1.0.0
type BlueprintScopeV100 struct {
	Entities       []string        `json:"entities"`
	Options        json.RawMessage `json:"options"`
	Transformation json.RawMessage `json:"transformation"`
}

// PipelineTask represents a smallest unit of execution inside a PipelinePlan
type PipelineTask struct {
	// Plugin name
	Plugin   string                 `json:"plugin" binding:"required"`
	Subtasks []string               `json:"subtasks"`
	Options  map[string]interface{} `json:"options"`
}

// PipelineStage consist of multiple PipelineTasks, they will be executed in parallel
type PipelineStage []*PipelineTask

// PipelinePlan consist of multiple PipelineStages, they will be executed in sequential order
type PipelinePlan []PipelineStage
