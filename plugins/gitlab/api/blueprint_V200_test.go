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
	const TestConnectionID uint64 = 1
	const TestTransformationRuleId uint64 = 2
	const TestID int = 37
	const TestGitlabEndPoint string = "https://gitlab.com/api/v4/"
	const TestHHttpUrlToRepo string = "https://this_is_cloneUrl"
	const TestToken string = "nddtf"
	const TestName string = "gitlab-test"
	const TestTransformationRuleName string = "github transformation rule"
	const TestProxy string = ""

	var TestGitlabProject *models.GitlabProject = &models.GitlabProject{
		ConnectionId: TestConnectionID,
		GitlabId:     TestID,
		Name:         TestName,

		TransformationRuleId: TestTransformationRuleId,
		CreatedDate:          time.Time{},
		HttpUrlToRepo:        TestHHttpUrlToRepo,
	}

	var TestTransformationRule *models.GitlabTransformationRule = &models.GitlabTransformationRule{
		Model: common.Model{
			ID: TestTransformationRuleId,
		},
		Name:   TestTransformationRuleName,
		PrType: "hey,man,wasup",
		RefdiffRule: map[string]interface{}{
			"tagsPattern": "pattern",
			"tagsLimit":   10,
			"tagsOrder":   "reverse semver",
		},
	}

	var TestGitlabConnection *models.GitlabConnection = &models.GitlabConnection{
		RestConnection: helper.RestConnection{
			BaseConnection: helper.BaseConnection{
				Name: TestName,
				Model: common.Model{
					ID: TestConnectionID,
				},
			},
			Endpoint:         TestGitlabEndPoint,
			Proxy:            TestProxy,
			RateLimitPerHour: 0,
		},
		AccessToken: helper.AccessToken{
			Token: TestToken,
		},
	}

	var ExpectRepoId = "gitlab:GitlabProject:1:37"

	var TestSubTaskMeta = []core.SubTaskMeta{
		tasks.CollectProjectMeta,
		tasks.ExtractProjectMeta,
		tasks.ConvertProjectMeta,
		tasks.CollectApiIssuesMeta,
		tasks.ExtractApiIssuesMeta,
		tasks.ConvertIssuesMeta,
		tasks.ConvertIssueLabelsMeta,
	}

	var ExpectPlans core.PipelinePlan = core.PipelinePlan{
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
				},
				SkipOnFail: false,
				Options: map[string]interface{}{
					"connectionId":         uint64(1),
					"projectId":            TestID,
					"transformationRuleId": TestTransformationRuleId,
					"transformationRules":  TestTransformationRule,
				},
			},
			{
				Plugin:     "gitextractor",
				SkipOnFail: false,
				Options: map[string]interface{}{
					"proxy":  "",
					"repoId": ExpectRepoId,
					"url":    "https://git:nddtf@this_is_cloneUrl",
				},
			},
		},
	}

	var ExpectScopes []core.Scope = []core.Scope{
		&code.Repo{
			DomainEntity: domainlayer.DomainEntity{
				Id: ExpectRepoId,
			},
			Name: TestName,
		},
		&ticket.Board{
			DomainEntity: domainlayer.DomainEntity{
				Id: ExpectRepoId,
			},
			Name:        TestName,
			Description: "",
			Url:         "",
			CreatedDate: nil,
			Type:        "",
		},
	}

	var err errors.Error

	bpScopes := []*core.BlueprintScopeV200{
		{
			Entities: []string{"CODE", "TICKET"},
			Id:       strconv.Itoa(TestID),
			Name:     TestName,
		},
	}

	// register gitlab plugin for NewDomainIdGenerator
	mockMeta := mocks.NewPluginMeta(t)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/gitlab")
	err = core.RegisterPlugin("gitlab", mockMeta)
	assert.Equal(t, err, nil)

	// Refresh Global Variables and set the sql mock
	BasicRes = unithelper.NewMockBasicRes(func(mockDal *mocks.Dal) {
		mockDal.On("First", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			dst := args.Get(0).(*models.GitlabConnection)
			*dst = *TestGitlabConnection
		}).Return(nil).Once()

		mockDal.On("First", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			dst := args.Get(0).(*models.GitlabProject)
			*dst = *TestGitlabProject
		}).Return(nil).Once()

		mockDal.On("First", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			dst := args.Get(0).(*models.GitlabTransformationRule)
			*dst = *TestTransformationRule
		}).Return(nil).Once()
	})
	connectionHelper = helper.NewConnectionHelper(
		BasicRes,
		validator.New(),
	)

	plans, scopes, err := MakePipelinePlanV200(TestSubTaskMeta, TestConnectionID, bpScopes)
	assert.Equal(t, err, nil)

	assert.Equal(t, ExpectPlans, plans)
	assert.Equal(t, ExpectScopes, scopes)
}
