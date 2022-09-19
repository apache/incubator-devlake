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
	"encoding/json"
	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/logger"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/runner"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakePipelinePlan(t *testing.T) {
	cfg := config.GetConfig()
	log := logger.Global.Nested("test")
	db, err := runner.NewGormDb(cfg, log)
	Init(cfg, log, db)
	bs := &core.BlueprintScopeV100{
		Entities: []string{"CODE"},
		Options: json.RawMessage(`{
              "repo": "lake",
              "owner": "merico-dev"
            }`),
		Transformation: json.RawMessage(`{
              "prType": "hey,man,wasup",
              "refdiff": {
                "tagsPattern": "pattern",
                "tagsLimit": 10,
                "tagsOrder": "reverse semver"
              },
              "dora": {
                "environment": "pattern",
                "environmentRegex": "xxxx"
              }
            }`),
	}
	scopes := make([]*core.BlueprintScopeV100, 0)
	scopes = append(scopes, bs)
	_, err = MakePipelinePlan(nil, 1, scopes)
	assert.Nil(t, err)
	//assert.Equal(t, plan, result3)
}
