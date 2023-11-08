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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
	"github.com/go-playground/validator/v10"
)

type MixScopes struct {
	ZentaoProduct *models.ZentaoProduct `json:"product"`
	ZentaoProject *models.ZentaoProject `json:"project"`
}

var logger log.Logger
var vld *validator.Validate
var connectionHelper *api.ConnectionApiHelper
var projectScopeHelper *api.ScopeApiHelper[models.ZentaoConnection, models.ZentaoProject, models.ZentaoScopeConfig]

var projectRemoteHelper *api.RemoteApiHelper[models.ZentaoConnection, models.ZentaoProject, models.ZentaoProject, api.BaseRemoteGroupResponse]
var basicRes context.BasicRes
var scHelper *api.ScopeConfigHelper[models.ZentaoScopeConfig, *models.ZentaoScopeConfig]

func Init(br context.BasicRes, p plugin.PluginMeta) {

	basicRes = br
	logger = basicRes.GetLogger()
	vld = validator.New()
	connectionHelper = api.NewConnectionHelper(basicRes, vld, p.Name())

	projectParams := &api.ReflectionParameters{
		ScopeIdFieldName:     "Id",
		ScopeIdColumnName:    "id",
		RawScopeParamName:    "ProjectId",
		SearchScopeParamName: "name",
	}
	projectScopeHelper = api.NewScopeHelper[models.ZentaoConnection, models.ZentaoProject, models.ZentaoScopeConfig](
		basicRes,
		vld,
		connectionHelper,
		api.NewScopeDatabaseHelperImpl[models.ZentaoConnection, models.ZentaoProject, models.ZentaoScopeConfig](
			basicRes, connectionHelper, projectParams),
		projectParams,
		nil,
	)

	projectRemoteHelper = api.NewRemoteHelper[models.ZentaoConnection, models.ZentaoProject, models.ZentaoProject, api.BaseRemoteGroupResponse](basicRes, vld, connectionHelper)
	scHelper = api.NewScopeConfigHelper[models.ZentaoScopeConfig, *models.ZentaoScopeConfig](basicRes, vld, p.Name())
}
