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
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/utils"

	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
)

func MakePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	connectionId uint64,
	bpScopes []*coreModels.BlueprintScope,
	skipCollectors bool,
) (coreModels.PipelinePlan, []plugin.Scope, errors.Error) {
	// load connection, scope and scopeConfig from the db
	connection, err := dsHelper.ConnSrv.FindByPk(connectionId)
	if err != nil {
		return nil, nil, err
	}
	scopeDetails, err := dsHelper.ScopeApi.MapScopeDetails(connectionId, bpScopes)
	if err != nil {
		return nil, nil, err
	}

	sc, err := makeScopeV200(connectionId, scopeDetails)
	if err != nil {
		return nil, nil, err
	}

	pp, err := makePipelinePlanV200(subtaskMetas, connection, scopeDetails)
	if err != nil {
		return nil, nil, err
	}

	return pp, sc, nil
}

func makeScopeV200(
	connectionId uint64,
	scopeDetails []*srvhelper.ScopeDetail[models.AzuredevopsRepo, models.AzuredevopsScopeConfig],
) ([]plugin.Scope, errors.Error) {
	sc := make([]plugin.Scope, 0, 3*len(scopeDetails))

	for _, scope := range scopeDetails {
		repo, scopeConfig := scope.Scope, scope.ScopeConfig
		entities := scopeConfig.Entities

		// We are treating empty entities as 'Selected All' since collecting a scope without any entity is pointless.
		if len(entities) == 0 {
			entities = plugin.DOMAIN_TYPES
		}

		isDomainCode := utils.StringsContains(entities, plugin.DOMAIN_TYPE_CODE_REVIEW) ||
			utils.StringsContains(entities, plugin.DOMAIN_TYPE_CODE)
		isDomainCICD := utils.StringsContains(entities, plugin.DOMAIN_TYPE_CICD)
		isDomainTicket := utils.StringsContains(entities, plugin.DOMAIN_TYPE_TICKET)

		id := didgen.NewDomainIdGenerator(&models.AzuredevopsRepo{}).Generate(connectionId, repo.Id)

		// DOMAIN_TYPE_CODE (i.e. gitextractor, rediff) only works if the repository is public and not disabled
		if isDomainCode && !repo.IsDisabled && !repo.IsPrivate {
			scopeRepo := code.NewRepo(id, repo.Name)
			sc = append(sc, scopeRepo)
		}

		// add cicd_scope to scopes
		if isDomainCICD {
			scopeCICD := devops.NewCicdScope(id, repo.Name)
			sc = append(sc, scopeCICD)
		}

		// add board to scopes
		if isDomainTicket {
			scopeTicket := ticket.NewBoard(id, repo.Name)
			sc = append(sc, scopeTicket)
		}
	}

	return sc, nil
}

func makePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	connection *models.AzuredevopsConnection,
	scopeDetails []*srvhelper.ScopeDetail[models.AzuredevopsRepo, models.AzuredevopsScopeConfig],
) (coreModels.PipelinePlan, errors.Error) {
	plans := make(coreModels.PipelinePlan, 0, 3*len(scopeDetails))

	for _, scope := range scopeDetails {
		repo, scopeConfig := scope.Scope, scope.ScopeConfig

		options := make(map[string]interface{})
		options["name"] = repo.Name // this is solely for the FE to display the repo name of a task

		options["connectionId"] = connection.ID
		options["organizationId"] = repo.OrganizationId
		options["projectId"] = repo.ProjectId
		options["externalId"] = repo.ExternalId
		options["repositoryId"] = repo.Id
		options["repositoryType"] = repo.Type

		if repo.Type == "" {
			options["repositoryType"] = models.RepositoryTypeADO
			logger.Warn(nil, "repository type for repoId: %v not found. falling back to TfsGit", repo.Id)
		}

		// We are treating empty entities as 'Selected All' since collecting a scope without any entity is pointless.
		entities := scopeConfig.Entities
		if len(entities) == 0 {
			entities = plugin.DOMAIN_TYPES
		}

		var selectedEntities []string
		var blockedEntities []string

		// We are unable to check out the code or gather pull requests for repositories that are disabled (DevOps)
		// or private (GitHub)
		if repo.IsDisabled || repo.IsPrivate {
			blockedEntities = append(blockedEntities, []string{
				plugin.DOMAIN_TYPE_CODE,
				plugin.DOMAIN_TYPE_CODE_REVIEW,
			}...)
		}

		// We are unable to gather pull requests from repositories not hosted on DevOps.
		// However, we can still check out the code if the repository is publicly available
		if repo.Type != models.RepositoryTypeADO {
			blockedEntities = append(blockedEntities, []string{
				plugin.DOMAIN_TYPE_CODE_REVIEW,
			}...)
		}

		for _, v := range entities {
			if !utils.StringsContains(blockedEntities, v) {
				selectedEntities = append(selectedEntities, v)
			}
		}

		var subtasks []string
		var err errors.Error
		if len(selectedEntities) > 0 {
			// if selectedEntities is empty MakePipelinePlanSubtasks assumes that we want to
			// enable all entity types
			subtasks, err = helper.MakePipelinePlanSubtasks(subtaskMetas, selectedEntities)
		}
		if err != nil {
			return nil, err
		}

		var stage []*coreModels.PipelineTask
		if len(subtasks) > 0 {
			stage = append(stage, &coreModels.PipelineTask{
				Plugin:   "azuredevops_go",
				Subtasks: subtasks,
				Options:  options,
			})
		} else {
			logger.Printf("Skipping azuredevops_go plugin due to empty subtasks. Please check your scope config")
		}

		// collect git data by gitextractor if CODE was requested
		if !repo.IsPrivate && !repo.IsDisabled && utils.StringsContains(entities, plugin.DOMAIN_TYPE_CODE) {
			cloneUrl, err := errors.Convert01(url.Parse(repo.RemoteUrl))
			if err != nil {
				return nil, err
			}

			if repo.Type == models.RepositoryTypeADO {
				cloneUrl.User = url.UserPassword("git", connection.Token)
			}
			stage = append(stage, &coreModels.PipelineTask{
				Plugin: "gitextractor",
				Options: map[string]interface{}{
					"url":            cloneUrl.String(),
					"name":           repo.Name,
					"repoId":         didgen.NewDomainIdGenerator(&models.AzuredevopsRepo{}).Generate(connection.ID, repo.Id),
					"proxy":          connection.Proxy,
					"noShallowClone": true,
				},
			})
		}

		plans = append(plans, stage)

		// refdiff part
		if scopeConfig.Refdiff != nil {
			task := &coreModels.PipelineTask{
				Plugin:  "refdiff",
				Options: scopeConfig.Refdiff,
			}
			plans = append(plans, coreModels.PipelineStage{task})
		}
	}
	return plans, nil
}
