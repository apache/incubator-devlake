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
	"gorm.io/gorm"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

func MakeDataSourcePipelinePlanV200(subtaskMetas []core.SubTaskMeta, connectionId uint64, bpScopes []*core.BlueprintScopeV200, syncPolicy *core.BlueprintSyncPolicy) (core.PipelinePlan, []core.Scope, errors.Error) {
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

	plan := make(core.PipelinePlan, len(bpScopes))
	plan, err = makeDataSourcePipelinePlanV200(subtaskMetas, plan, bpScopes, connection, apiClient, syncPolicy)
	if err != nil {
		return nil, nil, err
	}
	scopes, err := makeScopesV200(bpScopes, connection)
	if err != nil {
		return nil, nil, err
	}

	return plan, scopes, nil
}

func makeDataSourcePipelinePlanV200(
	subtaskMetas []core.SubTaskMeta,
	plan core.PipelinePlan,
	bpScopes []*core.BlueprintScopeV200,
	connection *models.GithubConnection,
	apiClient helper.ApiClientGetter,
	syncPolicy *core.BlueprintSyncPolicy,
) (core.PipelinePlan, errors.Error) {
	var err errors.Error
	var repoRes *tasks.GithubApiRepo
	for i, bpScope := range bpScopes {
		stage := plan[i]
		if stage == nil {
			stage = core.PipelineStage{}
		}
		githubRepo := &models.GithubRepo{}
		// get repo from db
		err = basicRes.GetDal().First(githubRepo, dal.Where(`connection_id = ? AND github_id = ?`, connection.ID, bpScope.Id))
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find repo %s", bpScope.Id))
		}
		transformationRule := &models.GithubTransformationRule{}
		// get transformation rules from db
		err = basicRes.GetDal().First(transformationRule, dal.Where(`id = ?`, githubRepo.TransformationRuleId))
		if err != nil && !goerror.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		// refdiff
		if transformationRule != nil && transformationRule.Refdiff != nil {
			// add a new task to next stage
			j := i + 1
			if j == len(plan) {
				plan = append(plan, nil)
			}
			refdiffOp := transformationRule.Refdiff
			if err != nil {
				return nil, err
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
		op := &tasks.GithubOptions{
			ConnectionId: githubRepo.ConnectionId,
			ScopeId:      bpScope.Id,
		}
		if syncPolicy.CreatedDateAfter != nil {
			op.CreatedDateAfter = syncPolicy.CreatedDateAfter.Format(time.RFC3339)
		}
		options, err := tasks.EncodeTaskOptions(op)
		if err != nil {
			return nil, err
		}
		stage, err = addGithub(subtaskMetas, connection, bpScope.Entities, stage, options)
		if err != nil {
			return nil, err
		}
		ownerRepo := strings.Split(githubRepo.Name, "/")
		if len(ownerRepo) != 2 {
			return nil, errors.Default.New("Fail to parse githubRepo.Name")
		}
		op.Owner = ownerRepo[0]
		op.Repo = ownerRepo[1]
		// add gitex stage and add repo to scopes
		if utils.StringsContains(bpScope.Entities, core.DOMAIN_TYPE_CODE) {
			repoRes, err = memorizedGetApiRepo(repoRes, op, apiClient)
			if err != nil {
				return nil, err
			}
			stage, err = addGitex(bpScope.Entities, connection, repoRes, stage)
			if err != nil {
				return nil, err
			}
		}
		plan[i] = stage
	}
	return plan, nil
}

func makeScopesV200(bpScopes []*core.BlueprintScopeV200, connection *models.GithubConnection) ([]core.Scope, errors.Error) {
	scopes := make([]core.Scope, 0)
	for _, bpScope := range bpScopes {
		githubRepo := &models.GithubRepo{}
		// get repo from db
		err := basicRes.GetDal().First(githubRepo, dal.Where(`connection_id = ? AND github_id = ?`, connection.ID, bpScope.Id))
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find repo%s", bpScope.Id))
		}
		if utils.StringsContains(bpScope.Entities, core.DOMAIN_TYPE_CODE_REVIEW) ||
			utils.StringsContains(bpScope.Entities, core.DOMAIN_TYPE_CODE) ||
			utils.StringsContains(bpScope.Entities, core.DOMAIN_TYPE_CROSS) {
			// if we don't need to collect gitex, we need to add repo to scopes here
			scopeRepo := &code.Repo{
				DomainEntity: domainlayer.DomainEntity{
					Id: didgen.NewDomainIdGenerator(&models.GithubRepo{}).Generate(connection.ID, githubRepo.GithubId),
				},
				Name: githubRepo.Name,
			}
			if githubRepo.ParentHTMLUrl != "" {
				scopeRepo.ForkedFrom = githubRepo.ParentHTMLUrl
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
				Name: githubRepo.Name,
			}
			scopes = append(scopes, scopeTicket)
		}
	}
	return scopes, nil
}
