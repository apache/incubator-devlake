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

package tasks

import (
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/argocd/models"
)

func CollectDataTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		CollectApplicationsMeta,
		ExtractApplicationsMeta,
		ConvertApplicationsMeta,
		CollectSyncOperationsMeta,
		ExtractSyncOperationsMeta,
		ConvertSyncOperationsMeta,
	}
}

func PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	op, err := DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, err
	}

	connectionHelper := api.NewConnectionHelper(
		taskCtx,
		nil,
		"argocd",
	)

	connection := &models.ArgocdConnection{}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, err
	}

	apiClient, err := NewArgocdApiClient(taskCtx, connection)
	if err != nil {
		return nil, err
	}

	if op.ScopeConfig == nil {
		op.ScopeConfig = &models.ArgocdScopeConfig{
			EnvNamePattern: "(?i)prod(.*)",
		}
	}
	if op.ScopeConfigId != 0 {
		scopeConfig := &models.ArgocdScopeConfig{}
		if err := taskCtx.GetDal().First(scopeConfig, dal.Where("id = ?", op.ScopeConfigId)); err != nil {
			return nil, errors.BadInput.Wrap(err, "fail to load scopeConfig")
		}
		op.ScopeConfig = scopeConfig
		op.ScopeConfigId = scopeConfig.ID
	} else if op.ScopeConfig != nil && op.ScopeConfig.ID != 0 {
		op.ScopeConfigId = op.ScopeConfig.ID
	}
	if op.ScopeConfig.EnvNamePattern == "" {
		op.ScopeConfig.EnvNamePattern = "(?i)prod(.*)"
	}

	regexEnricher := api.NewRegexEnricher()
	if err := regexEnricher.TryAdd(devops.ENV_NAME_PATTERN, op.ScopeConfig.EnvNamePattern); err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid envNamePattern")
	}
	if err := regexEnricher.TryAdd(devops.DEPLOYMENT, op.ScopeConfig.DeploymentPattern); err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid deploymentPattern")
	}
	if err := regexEnricher.TryAdd(devops.PRODUCTION, op.ScopeConfig.ProductionPattern); err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid productionPattern")
	}

	return &ArgocdTaskData{
		Options:       op,
		ApiClient:     apiClient,
		RegexEnricher: regexEnricher,
	}, nil
}
