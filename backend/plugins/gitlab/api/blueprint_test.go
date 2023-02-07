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
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	mockplugin "github.com/apache/incubator-devlake/mocks/core/plugin"
	mockaha "github.com/apache/incubator-devlake/mocks/helpers/pluginhelper/api/apihelperabstract"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/gitlab/tasks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProcessScope(t *testing.T) {
	connection := &models.GitlabConnection{
		BaseConnection: helper.BaseConnection{
			Name: "gitlab-test",
			Model: common.Model{
				ID: 1,
			},
		},
		GitlabConn: models.GitlabConn{
			RestConnection: helper.RestConnection{
				Endpoint:         "https://gitlab.com/api/v4/",
				Proxy:            "",
				RateLimitPerHour: 0,
			},
			AccessToken: helper.AccessToken{
				Token: "123",
			},
		},
	}
	mockApiCLient := mockaha.NewApiClientAbstract(t)
	repo := &tasks.GitlabApiProject{
		GitlabId:      12345,
		HttpUrlToRepo: "https://this_is_HttpUrlToRepo",
	}
	js, err := json.Marshal(repo)
	assert.Nil(t, err)
	res := &http.Response{}
	res.Body = io.NopCloser(bytes.NewBuffer(js))
	res.StatusCode = http.StatusOK
	mockApiCLient.On("Get", "projects/12345", mock.Anything, mock.Anything).Return(res, nil)

	mockMeta := mockplugin.NewPluginMeta(t)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/gitlab")
	err = plugin.RegisterPlugin("gitlab", mockMeta)
	assert.Nil(t, err)
	bs := &plugin.BlueprintScopeV100{
		Entities: []string{"CODE"},
		Options: json.RawMessage(`{
              "projectId": 12345
            }`),
		Transformation: json.RawMessage(`{
              "prType": "hey,man,wasup",
              "refdiff": {
                "tagsPattern": "pattern",
                "tagsLimit": 10,
                "tagsOrder": "reverse semver"
              },
              "productionPattern": "xxxx"
            }`),
	}
	scopes := make([]*plugin.BlueprintScopeV100, 0)
	scopes = append(scopes, bs)
	plan, err := makePipelinePlan(nil, scopes, mockApiCLient, connection)
	assert.Nil(t, err)

	expectPlan := plugin.PipelinePlan{
		plugin.PipelineStage{
			{
				Plugin:   "gitlab",
				Subtasks: []string{},
				Options: map[string]interface{}{
					"connectionId": uint64(1),
					"projectId":    float64(12345),
					"transformationRules": map[string]interface{}{
						"prType":            "hey,man,wasup",
						"productionPattern": "xxxx",
					},
				},
			},
			{
				Plugin: "gitextractor",
				Options: map[string]interface{}{
					"proxy":  "",
					"repoId": "gitlab:GitlabProject:1:12345",
					"url":    "https://git:123@this_is_HttpUrlToRepo",
				},
			},
		},
		plugin.PipelineStage{
			{
				Plugin: "refdiff",
				Options: map[string]interface{}{
					"repoId":      "gitlab:GitlabProject:1:12345",
					"tagsLimit":   float64(10),
					"tagsOrder":   "reverse semver",
					"tagsPattern": "pattern",
				},
			},
		},
		plugin.PipelineStage{
			{
				Plugin:   "dora",
				Subtasks: []string{"EnrichTaskEnv"},
				Options:  map[string]interface{}{},
			},
		},
	}
	assert.Equal(t, expectPlan, plan)
}
