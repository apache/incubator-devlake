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
	"strconv"
	"testing"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/helpers/unithelper"
	"github.com/apache/incubator-devlake/mocks"
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/gitlab/tasks"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMakeDataSourcePipelinePlanV200(t *testing.T) {
	const testConnectionID uint64 = 1
	const testTransformationRuleId uint64 = 2
	const testID int = 37
	const testGitlabEndPoint string = "https://gitlab.com/api/v4/"
	const testHttpUrlToRepo string = "https://this_is_cloneUrl"
	const testToken string = "nddtf"
	const testName string = "gitlab-test"
	const testTransformationRuleName string = "github transformation rule"
	const testProxy string = ""

	bpScopes := []*core.BlueprintScopeV200{
		{
			Entities: []string{core.DOMAIN_TYPE_CODE, core.DOMAIN_TYPE_TICKET, core.DOMAIN_TYPE_CICD},
			Id:       strconv.Itoa(testID),
			Name:     testName,
		},
	}

	var testGitlabProject *models.GitlabProject = &models.GitlabProject{
		ConnectionId: testConnectionID,
		GitlabId:     testID,
		Name:         testName,

		TransformationRuleId: testTransformationRuleId,
		CreatedDate:          time.Time{},
		HttpUrlToRepo:        testHttpUrlToRepo,
	}

	var testTransformationRule *models.GitlabTransformationRule = &models.GitlabTransformationRule{
		Model: common.Model{
			ID: testTransformationRuleId,
		},
		Name:   testTransformationRuleName,
		PrType: "hey,man,wasup",
		Refdiff: map[string]interface{}{
			"tagsPattern": "pattern",
			"tagsLimit":   10,
			"tagsOrder":   "reverse semver",
		},
	}

	var testGitlabConnection *models.GitlabConnection = &models.GitlabConnection{
		RestConnection: helper.RestConnection{
			BaseConnection: helper.BaseConnection{
				Name: testName,
				Model: common.Model{
					ID: testConnectionID,
				},
			},
			Endpoint:         testGitlabEndPoint,
			Proxy:            testProxy,
			RateLimitPerHour: 0,
		},
		AccessToken: helper.AccessToken{
			Token: testToken,
		},
	}

	var expectRepoId = "gitlab:GitlabProject:1:37"

	var testSubTaskMeta = []core.SubTaskMeta{
		tasks.CollectProjectMeta,
		tasks.ExtractProjectMeta,
		tasks.ConvertProjectMeta,
		tasks.CollectApiIssuesMeta,
		tasks.ExtractApiIssuesMeta,
		tasks.ConvertIssuesMeta,
		tasks.ConvertIssueLabelsMeta,
		tasks.CollectApiJobsMeta,
		tasks.ExtractApiJobsMeta,
		tasks.CollectApiPipelinesMeta,
		tasks.ExtractApiPipelinesMeta,
	}

	var expectPlans core.PipelinePlan = core.PipelinePlan{
		{
			{
				Plugin:     "refdiff",
				SkipOnFail: false,
				Options: map[string]interface{}{
					"tagsLimit":   10,
					"tagsOrder":   "reverse semver",
					"tagsPattern": "pattern",
				},
			},
			{
				Plugin: "gitlab",
				Subtasks: []string{
					tasks.CollectProjectMeta.Name,
					tasks.ExtractProjectMeta.Name,
					tasks.ConvertProjectMeta.Name,
					tasks.CollectApiIssuesMeta.Name,
					tasks.ExtractApiIssuesMeta.Name,
					tasks.ConvertIssuesMeta.Name,
					tasks.ConvertIssueLabelsMeta.Name,
					tasks.CollectApiJobsMeta.Name,
					tasks.ExtractApiJobsMeta.Name,
					tasks.CollectApiPipelinesMeta.Name,
					tasks.ExtractApiPipelinesMeta.Name,
				},
				SkipOnFail: false,
				Options: map[string]interface{}{
					"connectionId": uint64(1),
					"scopeId":      testID,
					"entities":     bpScopes[0].Entities,
				},
			},
			{
				Plugin:     "gitextractor",
				SkipOnFail: false,
				Options: map[string]interface{}{
					"proxy":  "",
					"repoId": expectRepoId,
					"url":    "https://git:nddtf@this_is_cloneUrl",
				},
			},
		},
	}

	var expectScopes []core.Scope = []core.Scope{
		&code.Repo{
			DomainEntity: domainlayer.DomainEntity{
				Id: expectRepoId,
			},
			Name:       testName,
			ForkedFrom: testGitlabProject.ForkedFromProjectWebUrl,
		},
		&devops.CicdScope{
			DomainEntity: domainlayer.DomainEntity{
				Id: expectRepoId,
			},
			Name:        testName,
			Description: "",
			Url:         "",
		},
		&ticket.Board{
			DomainEntity: domainlayer.DomainEntity{
				Id: expectRepoId,
			},
			Name:        testName,
			Description: "",
			Url:         "",
			CreatedDate: nil,
			Type:        "",
		},
	}

	var err errors.Error

	// register gitlab plugin for NewDomainIdGenerator
	mockMeta := mocks.NewPluginMeta(t)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/gitlab")
	err = core.RegisterPlugin("gitlab", mockMeta)
	assert.Equal(t, err, nil)

	// Refresh Global Variables and set the sql mock
	BasicRes = unithelper.DummyBasicRes(func(mockDal *mocks.Dal) {
		mockDal.On("First", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			dst := args.Get(0).(*models.GitlabConnection)
			*dst = *testGitlabConnection
		}).Return(nil).Once()

		mockDal.On("First", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			dst := args.Get(0).(*models.GitlabProject)
			*dst = *testGitlabProject
		}).Return(nil).Twice()

		mockDal.On("First", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			dst := args.Get(0).(*models.GitlabTransformationRule)
			*dst = *testTransformationRule
		}).Return(nil).Once()
	})
	connectionHelper = helper.NewConnectionHelper(
		BasicRes,
		validator.New(),
	)

	plans, scopes, err := MakePipelinePlanV200(testSubTaskMeta, testConnectionID, bpScopes)
	assert.Equal(t, err, nil)

	assert.Equal(t, expectPlans, plans)
	assert.Equal(t, expectScopes, scopes)
}
