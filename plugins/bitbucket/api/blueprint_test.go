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
	"github.com/apache/incubator-devlake/mocks"
	"testing"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/stretchr/testify/assert"
)

func TestMakePipelinePlan(t *testing.T) {
	var mockGetter repoGetter
	mockGetter = func(connectionId uint64, owner, repo string) (string, string, errors.Error) {
		return "https://thenicetgp@bitbucket.org/thenicetgp/lake.git", "secret", nil
	}
	scope := &core.BlueprintScopeV100{
		Entities: []string{core.DOMAIN_TYPE_CODE, core.DOMAIN_TYPE_TICKET, core.DOMAIN_TYPE_CODE_REVIEW, core.DOMAIN_TYPE_CROSS},
		Options: []byte(`{
                            "owner": "thenicetgp",
                            "repo": "lake"
                        }`),
		Transformation: nil,
	}
	mockMeta := mocks.NewPluginMeta(t)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/bitbucket")
	err := core.RegisterPlugin("bitbucket", mockMeta)
	assert.Nil(t, err)
	plan, err := makePipelinePlan(nil, 1, mockGetter, []*core.BlueprintScopeV100{scope})
	assert.Nil(t, err)
	for _, stage := range plan {
		for _, task := range stage {
			if task.Plugin == "gitextractor" {
				assert.Equal(t, task.Options["url"], "https://thenicetgp:secret@bitbucket.org/thenicetgp/lake.git")
				return
			}
		}
	}
	t.Fatal("no gitextractor plugin")
}
