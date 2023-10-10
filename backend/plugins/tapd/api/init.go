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
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"github.com/go-playground/validator/v10"
)

var vld *validator.Validate
var connectionHelper *api.ConnectionApiHelper
var basicRes context.BasicRes
var scopeHelper *api.ScopeApiHelper[models.TapdConnection, models.TapdWorkspace, models.TapdScopeConfig]
var remoteHelper *api.RemoteApiHelper[models.TapdConnection, models.TapdWorkspace, models.TapdWorkspace, api.BaseRemoteGroupResponse]
var scHelper *api.ScopeConfigHelper[models.TapdScopeConfig, *models.TapdScopeConfig]

func Init(br context.BasicRes, p plugin.PluginMeta) {

	basicRes = br
	vld = validator.New()
	connectionHelper = api.NewConnectionHelper(
		basicRes,
		vld,
		p.Name(),
	)
	params := &api.ReflectionParameters{
		ScopeIdFieldName:     "Id",
		ScopeIdColumnName:    "id",
		RawScopeParamName:    "WorkSpaceId",
		SearchScopeParamName: "name",
	}
	scopeHelper = api.NewScopeHelper[models.TapdConnection, models.TapdWorkspace, models.TapdScopeConfig](
		basicRes,
		vld,
		connectionHelper,
		api.NewScopeDatabaseHelperImpl[models.TapdConnection, models.TapdWorkspace, models.TapdScopeConfig](
			basicRes, connectionHelper, params),
		params,
		nil,
	)
	remoteHelper = api.NewRemoteHelper[models.TapdConnection, models.TapdWorkspace, models.TapdWorkspace, api.BaseRemoteGroupResponse](
		basicRes,
		vld,
		connectionHelper,
	)
	scHelper = api.NewScopeConfigHelper[models.TapdScopeConfig, *models.TapdScopeConfig](
		basicRes,
		vld,
		p.Name(),
	)
}
