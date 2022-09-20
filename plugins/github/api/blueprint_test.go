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
	"github.com/apache/incubator-devlake/mocks"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakePipelinePlan(t *testing.T) {
	cmd := &cobra.Command{Use: "github"}
	option := &tasks.GithubOptions{
		ConnectionId: 1,
		Tasks:        nil,
		Since:        "",
		Owner:        "merico-dev",
		Repo:         "lake",
		TransformationRules: models.TransformationRules{
			PrType:               "hey,man,wasup",
			PrComponent:          "component/(.*)$",
			PrBodyClosePattern:   "(?mi)(fix|close|resolve|fixes|closes|resolves|fixed|closed|resolved)[\\s]*.*(((and )?(#|https:\\/\\/github.com\\/%s\\/%s\\/issues\\/)\\d+[ ]*)+)",
			IssueSeverity:        "severity/(.*)$",
			IssuePriority:        "^(highest|high|medium|low)$",
			IssueComponent:       "component/(.*)$",
			IssueTypeBug:         "^(bug|failure|error)$",
			IssueTypeIncident:    "",
			IssueTypeRequirement: "^(feat|feature|proposal|requirement)$",
			DeployTagPattern:     "(?i)deploy",
		},
	}
	mockConnHelper := mocks.NewBpHelper(t)
	mockConnHelper.On("GetApiRepo", uint64(1), option).Return(&tasks.GithubApiRepo{
		Name:     "test",
		GithubId: 1,
		CloneUrl: "CloneUrl",
	}, "", "", nil)
	mockMeta := mocks.NewPluginMeta(t)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/github")

	err := core.RegisterPlugin(cmd.Use, mockMeta)
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
	plan, err := DoMakePipeline(nil, 1, scopes, mockConnHelper)

	assert.Nil(t, err)
	planJson, err1 := json.Marshal(plan)
	assert.Nil(t, err1)
	fmt.Println(string(planJson))
	expectPlan := `[[{"plugin":"github","subtasks":[],"options":{"connectionId":1,"owner":"merico-dev","repo":"lake","transformationRules":{"prType":"hey,man,wasup"}}},{"plugin":"gitextractor","subtasks":null,"options":{"proxy":"","repoId":"github:GithubRepo:1:1","url":"//git:@CloneUrl"}}],[{"plugin":"refdiff","subtasks":null,"options":{"tagsLimit":10,"tagsOrder":"reverse semver","tagsPattern":"pattern"}}],[{"plugin":"dora","subtasks":null,"options":{"repoId":"github:GithubRepo:1:1","tasks":["EnrichTaskEnv"],"transformation":{"environment":"pattern","environmentRegex":"xxxx"}}}]]`
	assert.Equal(t, expectPlan, string(planJson))

}
