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
	"fmt"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/bitbucket/tasks"
	"github.com/go-playground/validator/v10"
	"net/url"
	"time"
)

func MakeDataSourcePipelinePlanV200(subtaskMetas []plugin.SubTaskMeta, connectionId uint64, bpScopes []*plugin.BlueprintScopeV200, syncPolicy *plugin.BlueprintSyncPolicy) (plugin.PipelinePlan, []plugin.Scope, errors.Error) {
	connectionHelper := helper.NewConnectionHelper(basicRes, validator.New())
	// get the connection info for url
	connection := &models.BitbucketConnection{}
	err := connectionHelper.FirstById(connection, connectionId)
	if err != nil {
		return nil, nil, err
	}

	plan := make(plugin.PipelinePlan, len(bpScopes))
	plan, err = makeDataSourcePipelinePlanV200(subtaskMetas, plan, bpScopes, connection, syncPolicy)
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
	subtaskMetas []plugin.SubTaskMeta,
	plan plugin.PipelinePlan,
	bpScopes []*plugin.BlueprintScopeV200,
	connection *models.BitbucketConnection,
	syncPolicy *plugin.BlueprintSyncPolicy,
) (plugin.PipelinePlan, errors.Error) {
	var err errors.Error
	for i, bpScope := range bpScopes {
		stage := plan[i]
		if stage == nil {
			stage = plugin.PipelineStage{}
		}
		repo := &models.BitbucketRepo{}
		// get repo from db
		err = basicRes.GetDal().First(repo, dal.Where(`connection_id = ? AND bitbucket_id = ?`, connection.ID, bpScope.Id))
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find repo %s", bpScope.Id))
		}
		transformationRule := &models.BitbucketTransformationRule{}
		// get transformation rules from db
		db := basicRes.GetDal()
		err = db.First(transformationRule, dal.Where(`id = ?`, repo.TransformationRuleId))
		if err != nil && !db.IsErrorNotFound(err) {
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
			refdiffOp["repoId"] = didgen.NewDomainIdGenerator(&models.BitbucketRepo{}).Generate(connection.ID, repo.BitbucketId)
			plan[j] = plugin.PipelineStage{
				{
					Plugin:  "refdiff",
					Options: refdiffOp,
				},
			}
			transformationRule.Refdiff = nil
		}

		// construct task options for bitbucket
		op := &tasks.BitbucketOptions{
			ConnectionId: repo.ConnectionId,
			FullName:     repo.BitbucketId,
		}
		if syncPolicy.TimeAfter != nil {
			op.TimeAfter = syncPolicy.TimeAfter.Format(time.RFC3339)
		}
		options, err := tasks.EncodeTaskOptions(op)
		if err != nil {
			return nil, err
		}

		subtasks, err := helper.MakePipelinePlanSubtasks(subtaskMetas, bpScope.Entities)
		if err != nil {
			return nil, err
		}
		stage = append(stage, &plugin.PipelineTask{
			Plugin:   "bitbucket",
			Subtasks: subtasks,
			Options:  options,
		})
		if err != nil {
			return nil, err
		}

		// add gitex stage
		if utils.StringsContains(bpScope.Entities, plugin.DOMAIN_TYPE_CODE) {
			cloneUrl, err := errors.Convert01(url.Parse(repo.CloneUrl))
			if err != nil {
				return nil, err
			}
			cloneUrl.User = url.UserPassword(connection.Username, connection.Password)
			stage = append(stage, &plugin.PipelineTask{
				Plugin: "gitextractor",
				Options: map[string]interface{}{
					"url":    cloneUrl.String(),
					"repoId": didgen.NewDomainIdGenerator(&models.BitbucketRepo{}).Generate(connection.ID, repo.BitbucketId),
					"proxy":  connection.Proxy,
				},
			})

		}
		plan[i] = stage
	}
	return plan, nil
}

func makeScopesV200(bpScopes []*plugin.BlueprintScopeV200, connection *models.BitbucketConnection) ([]plugin.Scope, errors.Error) {
	scopes := make([]plugin.Scope, 0)
	for _, bpScope := range bpScopes {
		repo := &models.BitbucketRepo{}
		// get repo from db
		err := basicRes.GetDal().First(repo, dal.Where(`connection_id = ? AND bitbucket_id = ?`, connection.ID, bpScope.Id))
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find repo%s", bpScope.Id))
		}
		if utils.StringsContains(bpScope.Entities, plugin.DOMAIN_TYPE_CODE_REVIEW) ||
			utils.StringsContains(bpScope.Entities, plugin.DOMAIN_TYPE_CODE) ||
			utils.StringsContains(bpScope.Entities, plugin.DOMAIN_TYPE_CROSS) {
			// if we don't need to collect gitex, we need to add repo to scopes here
			scopeRepo := &code.Repo{
				DomainEntity: domainlayer.DomainEntity{
					Id: didgen.NewDomainIdGenerator(&models.BitbucketRepo{}).Generate(connection.ID, repo.BitbucketId),
				},
				Name: repo.Name,
			}
			scopes = append(scopes, scopeRepo)
		}
		// add cicd_scope to scopes
		if utils.StringsContains(bpScope.Entities, plugin.DOMAIN_TYPE_CICD) {
			scopeCICD := &devops.CicdScope{
				DomainEntity: domainlayer.DomainEntity{
					Id: didgen.NewDomainIdGenerator(&models.BitbucketRepo{}).Generate(connection.ID, repo.BitbucketId),
				},
				Name: repo.Name,
			}
			scopes = append(scopes, scopeCICD)
		}
		// add board to scopes
		if utils.StringsContains(bpScope.Entities, plugin.DOMAIN_TYPE_TICKET) {
			scopeTicket := &ticket.Board{
				DomainEntity: domainlayer.DomainEntity{
					Id: didgen.NewDomainIdGenerator(&models.BitbucketRepo{}).Generate(connection.ID, repo.BitbucketId),
				},
				Name: repo.Name,
			}
			scopes = append(scopes, scopeTicket)
		}
	}
	return scopes, nil
}
