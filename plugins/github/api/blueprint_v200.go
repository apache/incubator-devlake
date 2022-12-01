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
	"context"
	goerror "errors"
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
	"github.com/apache/incubator-devlake/utils"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

func MakeDataSourcePipelinePlanV200(subtaskMetas []core.SubTaskMeta, connectionId uint64, bpScopes []*core.BlueprintScopeV200) (core.PipelinePlan, []core.Scope, errors.Error) {
	db := basicRes.GetDal()
	connectionHelper := helper.NewConnectionHelper(basicRes, validator.New())
	// get the connection info for url
	connection := &models.GithubConnection{}
	err := connectionHelper.FirstById(connection, connectionId)
	if err != nil {
		return nil, nil, err
	}

	token := strings.Split(connection.Token, ",")[0]
	apiClient, err := helper.NewApiClient(
		context.TODO(),
		connection.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", token),
		},
		10*time.Second,
		connection.Proxy,
		basicRes,
	)
	if err != nil {
		return nil, nil, err
	}

	plan := make(core.PipelinePlan, 0, len(bpScopes))
	scopes := make([]core.Scope, 0, len(bpScopes))
	for i, bpScope := range bpScopes {
		var githubRepo *models.GithubRepo
		// get repo from db
		err = db.First(githubRepo, dal.Where(`id = ?`, bpScope.Id))
		if err != nil {
			return nil, nil, err
		}
		var transformationRule *models.GithubTransformationRule
		// get transformation rules from db
		err = db.First(transformationRule, dal.Where(`id = ?`, githubRepo.TransformationRuleId))
		if err != nil && goerror.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, err
		}
		var scope []core.Scope
		// make pipeline for each bpScope
		plan, scope, err = makeDataSourcePipelinePlanV200(subtaskMetas, i, plan, bpScope, connection, apiClient, githubRepo, transformationRule)
		if err != nil {
			return nil, nil, err
		}
		if len(scope) > 0 {
			scopes = append(scopes, scope...)
		}

	}

	return plan, scopes, nil
}

func makeDataSourcePipelinePlanV200(
	subtaskMetas []core.SubTaskMeta,
	i int,
	plan core.PipelinePlan,
	bpScope *core.BlueprintScopeV200,
	connection *models.GithubConnection,
	apiClient helper.ApiClientGetter,
	githubRepo *models.GithubRepo,
	transformationRule *models.GithubTransformationRule,
) (core.PipelinePlan, []core.Scope, errors.Error) {
	var err errors.Error
	var stage core.PipelineStage
	var repo *tasks.GithubApiRepo
	scopes := make([]core.Scope, 0)
	// refdiff
	if transformationRule != nil && transformationRule.Refdiff != nil {
		// add a new task to next stage
		j := i + 1
		if j == len(plan) {
			plan = append(plan, nil)
		}
		refdiffOp := transformationRule.Refdiff
		if err != nil {
			return nil, nil, err
		}
		plan[j] = core.PipelineStage{
			{
				Plugin:  "refdiff",
				Options: refdiffOp,
			},
		}
		transformationRule.Refdiff = nil
	}

	// construct task options for github
	var options map[string]interface{}
	err = errors.Convert(mapstructure.Decode(githubRepo, &options))
	if err != nil {
		return nil, nil, err
	}
	// make sure task options is valid
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, nil, err
	}

	var transformationRuleMap map[string]interface{}
	err = errors.Convert(mapstructure.Decode(transformationRule, &transformationRuleMap))
	if err != nil {
		return nil, nil, err
	}
	options["transformationRules"] = transformationRuleMap
	stage, err = addGithub(subtaskMetas, connection, bpScope.Entities, stage, options)
	if err != nil {
		return nil, nil, err
	}
	// add gitex stage and add repo to scopes
	if utils.StringsContains(bpScope.Entities, core.DOMAIN_TYPE_CODE) {
		repo, err = memorizedGetApiRepo(repo, op, apiClient)
		if err != nil {
			return nil, nil, err
		}
		stage, err = addGitex(bpScope.Entities, connection, repo, stage)
		if err != nil {
			return nil, nil, err
		}
		scopeRepo := &code.Repo{
			DomainEntity: domainlayer.DomainEntity{
				Id: didgen.NewDomainIdGenerator(&models.GithubRepo{}).Generate(connection.ID, githubRepo.GithubId),
			},
			Name: fmt.Sprintf("%s/%s", githubRepo.OwnerLogin, githubRepo.Name),
		}
		if repo.Parent != nil {
			scopeRepo.ForkedFrom = repo.Parent.HTMLUrl
		}
		scopes = append(scopes, scopeRepo)
	} else if utils.StringsContains(bpScope.Entities, core.DOMAIN_TYPE_CODE_REVIEW) {
		// if we don't need to collect gitex, we need to add repo to scopes here
		scopeRepo := &code.Repo{
			DomainEntity: domainlayer.DomainEntity{
				Id: didgen.NewDomainIdGenerator(&models.GithubRepo{}).Generate(connection.ID, githubRepo.GithubId),
			},
			Name: fmt.Sprintf("%s/%s", githubRepo.OwnerLogin, githubRepo.Name),
		}
		scopes = append(scopes, scopeRepo)
	}
	// add cicd_scope to scopes
	if utils.StringsContains(bpScope.Entities, core.DOMAIN_TYPE_CICD) {
		scopeCICD := &devops.CicdScope{
			DomainEntity: domainlayer.DomainEntity{
				Id: didgen.NewDomainIdGenerator(&models.GithubRepo{}).Generate(connection.ID, githubRepo.GithubId),
			},
			Name: fmt.Sprintf("%s/%s", githubRepo.OwnerLogin, githubRepo.Name),
		}
		scopes = append(scopes, scopeCICD)
	}
	// add board to scopes
	if utils.StringsContains(bpScope.Entities, core.DOMAIN_TYPE_TICKET) {
		scopeTicket := &ticket.Board{
			DomainEntity: domainlayer.DomainEntity{
				Id: didgen.NewDomainIdGenerator(&models.GithubRepo{}).Generate(connection.ID, githubRepo.GithubId),
			},
			Name: fmt.Sprintf("%s/%s", githubRepo.OwnerLogin, githubRepo.Name),
		}
		scopes = append(scopes, scopeTicket)
	}

	plan[i] = stage

	return plan, scopes, nil
}
