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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"net/http"
)

// @Summary pipelines plan for starrocks
// @Description pipelines plan for starrocks
// @Tags plugins/starrocks
// @Accept application/json
// @Param blueprint body StarRocksPipelinePlan true "json"
// @Router /pipelines/starrocks/pipeline-plan [post]
func PostStarRocksPipeline(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	blueprint := &StarRocksPipelinePlan{}
	return &core.ApiResourceOutput{Body: blueprint, Status: http.StatusOK}, nil
}

type StarRocksPipelinePlan [][]struct {
	Plugin  string `json:"plugin"`
	Options struct {
		SourceType  string   `json:"source_type"`
		SourceDsn   string   `json:"source_dsn"`
		Host        string   `json:"host"`
		Port        int      `json:"port"`
		User        string   `json:"user"`
		Password    string   `json:"password"`
		Database    string   `json:"database"`
		BePort      int      `json:"be_port"`
		Tables      []string `json:"tables"`
		BatchSize   int      `json:"batch_size"`
		Extra       string   `json:"extra"`
		DomainLayer string   `json:"domain_layer"`
	} `json:"options"`
}
