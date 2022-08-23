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

import (
	"github.com/apache/incubator-devlake/plugins/core"
	"net/http"
)

// @Summary pipelines plan for refdiff
// @Description pipelines plan for refdiff
// @Tags plugins/refdiff
// @Accept application/json
// @Param blueprint body RefdiffPipelinePlan true "json"
// @Router /pipelines/refdiff/pipeline-plan [post]
func PostRefdiffPipeline(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	blueprint := &RefdiffPipelinePlan{}
	return &core.ApiResourceOutput{Body: blueprint, Status: http.StatusOK}, nil
}

type RefdiffPipelinePlan [][]struct {
	Plugin  string `json:"plugin"`
	Options struct {
		RepoID string `json:"repoId"`
		Pairs  []struct {
			NewRef string `json:"newRef"`
			OldRef string `json:"oldRef"`
		} `json:"pairs"`
	} `json:"options"`
}
