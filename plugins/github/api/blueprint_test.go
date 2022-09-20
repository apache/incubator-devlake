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
	"fmt"
	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/logger"
	"github.com/apache/incubator-devlake/mocks"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/runner"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestMakePipelinePlan(t *testing.T) {
	cmd := &cobra.Command{Use: "github"}
	cfg := config.GetConfig()
	log := logger.Global.Nested("test")
	db, err := runner.NewGormDb(cfg, log)
	Init(cfg, log, db)
	mockMeta := mocks.NewPluginMeta(t)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/github")
	err = core.RegisterPlugin(cmd.Use, mockMeta)
	if err != nil {
		panic(err)
	}
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
	plan, err := MakePipelinePlan(nil, 1, scopes)
	assert.Nil(t, err)
	planJson, err1 := json.Marshal(plan)
	assert.Nil(t, err1)
	githubConn := &models.GithubConnection{}
	err = connectionHelper.FirstById(githubConn, 1)
	assert.Nil(t, err)
	token := strings.Split(githubConn.Token, ",")[0]
	expectPlan := fmt.Sprintf(`[[{"plugin":"github","subtasks":[],"options":{"connectionId":1,"owner":"merico-dev","repo":"lake","transformationRules":{"prType":"hey,man,wasup"}}},{"plugin":"gitextractor","subtasks":null,"options":{"proxy":"","repoId":"github:GithubRepo:1:491450511","url":"https://git:%s@github.com/merico-dev/lake.git"}}],[{"plugin":"refdiff","subtasks":null,"options":{"tagsLimit":10,"tagsOrder":"reverse semver","tagsPattern":"pattern"}}],[{"plugin":"dora","subtasks":null,"options":{"repoId":"github:GithubRepo:1:491450511","tasks":["EnrichTaskEnv"],"transformation":{"environment":"pattern","environmentRegex":"xxxx"}}}]]`,
		token,
	)
	assert.Equal(t, expectPlan, string(planJson))

}
