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
	"github.com/apache/incubator-devlake/mocks"
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/bitbucket/tasks"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"path"
	"testing"
)

func TestMakePipelinePlan(t *testing.T) {
	connection := &models.BitbucketConnection{
		RestConnection: helper.RestConnection{
			BaseConnection: helper.BaseConnection{
				Name: "github-test",
				Model: common.Model{
					ID: 1,
				},
			},
			Endpoint:         "https://TestBitBucket/",
			Proxy:            "",
			RateLimitPerHour: 0,
		},
		BasicAuth: helper.BasicAuth{
			Username: "Username",
			Password: "Password",
		},
	}

	mockApiCLient := mocks.NewApiClientGetter(t)
	repo := &tasks.BitbucketApiRepo{
		Links: struct {
			Clone []struct {
				Href string `json:"href"`
				Name string `json:"name"`
			} `json:"clone"`
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			Html struct {
				Href string `json:"href"`
			} `json:"html"`
		}{
			Clone: []struct {
				Href string `json:"href"`
				Name string `json:"name"`
			}{
				{
					Href: "https://bitbucket.org/thenicetgp/lake.git",
					Name: "https",
				},
			},
			Self: struct {
				Href string `json:"href"`
			}{},
			Html: struct {
				Href string `json:"href"`
			}{},
		},
	}
	js, err := json.Marshal(repo)
	assert.Nil(t, err)
	res := &http.Response{}
	res.Body = io.NopCloser(bytes.NewBuffer(js))
	res.StatusCode = http.StatusOK
	mockApiCLient.On("Get", path.Join("repositories", "thenicetgp", "lake"), mock.Anything, mock.Anything).Return(res, nil)
	mockMeta := mocks.NewPluginMeta(t)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/bitbucket")
	err = core.RegisterPlugin("bitbucket", mockMeta)
	scope := &core.BlueprintScopeV100{
		Entities: []string{core.DOMAIN_TYPE_CODE, core.DOMAIN_TYPE_TICKET, core.DOMAIN_TYPE_CODE_REVIEW, core.DOMAIN_TYPE_CROSS},
		Options: []byte(`{
                            "owner": "thenicetgp",
                            "repo": "lake"
                        }`),
		Transformation: nil,
	}
	scopes := make([]*core.BlueprintScopeV100, 0)
	scopes = append(scopes, scope)
	assert.Nil(t, err)
	plan, err := makePipelinePlan(nil, scopes, mockApiCLient, connection)
	assert.Nil(t, err)

	expectPlan := core.PipelinePlan{
		core.PipelineStage{
			{
				Plugin:   "bitbucket",
				Subtasks: []string{},
				Options: map[string]interface{}{
					"connectionId":        uint64(1),
					"owner":               "thenicetgp",
					"repo":                "lake",
					"transformationRules": map[string]interface{}{},
				},
			},
			{
				Plugin: "gitextractor",
				Options: map[string]interface{}{
					"repoId": "bitbucket:BitbucketRepo:1:thenicetgp/lake",
					"url":    "https://thenicetgp:Password@bitbucket.org/thenicetgp/lake.git",
				},
			},
		},
	}
	assert.Equal(t, expectPlan, plan)
}
