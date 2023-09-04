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

package plugin

import (
	"encoding/json"

	"github.com/apache/incubator-devlake/core/errors"
	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/server/services/remote/bridge"
	"github.com/apache/incubator-devlake/server/services/remote/models"
)

type (
	remoteMetricPlugin struct {
		*remotePluginImpl
	}
	remoteDatasourcePlugin struct {
		*remotePluginImpl
	}
)

func (p remoteMetricPlugin) MakeMetricPluginPipelinePlanV200(projectName string, options json.RawMessage) (coreModels.PipelinePlan, errors.Error) {
	return nil, errors.Internal.New("Remote metric coreModels not supported")
}

func (p remoteDatasourcePlugin) MakeDataSourcePipelinePlanV200(
	connectionId uint64,
	bpScopes []*coreModels.BlueprintScope,
	syncPolicy coreModels.BlueprintSyncPolicy,
) (coreModels.PipelinePlan, []plugin.Scope, errors.Error) {
	connection := p.connectionTabler.New()
	err := p.connHelper.FirstById(connection, connectionId)
	if err != nil {
		return nil, nil, err
	}

	db := basicRes.GetDal()
	var toolScopeConfigPairs = make([]interface{}, len(bpScopes))
	for i, bpScope := range bpScopes {
		toolScope, scopeConfig, err := p.getScopeAndConfig(db, connectionId, bpScope.ScopeId)
		if err != nil {
			return nil, nil, err
		}
		toolScopeConfigPairs[i] = []interface{}{toolScope, scopeConfig}
	}

	planData := models.PipelineData{}
	err = p.invoker.Call("make-pipeline", bridge.DefaultContext, toolScopeConfigPairs, connection.Unwrap()).Get(&planData)
	if err != nil {
		return nil, nil, err
	}
	scopes, err := toDomainScopes(planData.Scopes)
	if err != nil {
		return nil, nil, err
	}
	// store these domain scopes in the DB (remote plugins will not explicitly do this via standalone extractor/convertor pairs)
	for _, scope := range scopes {
		err = db.CreateOrUpdate(scope)
		if err != nil {
			return nil, nil, err
		}
	}
	return planData.Plan, scopes, nil
}

func toDomainScopes(dynamicScopes []models.DynamicDomainScope) ([]plugin.Scope, errors.Error) {
	var scopes []plugin.Scope
	for _, dynamicScope := range dynamicScopes {
		scope, err := dynamicScope.Load()
		if err != nil {
			return nil, err
		}
		scopes = append(scopes, scope)
	}
	return scopes, nil
}

var _ models.RemotePlugin = (*remoteMetricPlugin)(nil)
var _ plugin.MetricPluginBlueprintV200 = (*remoteMetricPlugin)(nil)
var _ models.RemotePlugin = (*remoteDatasourcePlugin)(nil)
var _ plugin.DataSourcePluginBlueprintV200 = (*remoteDatasourcePlugin)(nil)
