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
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/go-playground/validator/v10"
)

var vld *validator.Validate
var connectionHelper *api.ConnectionApiHelper
var scopeHelper *api.ScopeApiHelper[models.GitlabConnection, models.GitlabProject, models.GitlabScopeConfig]
var remoteHelper *api.RemoteApiHelper[models.GitlabConnection, models.GitlabProject, models.GitlabApiProject, models.GroupResponse]
var basicRes context.BasicRes
var scHelper *api.ScopeConfigHelper[models.GitlabScopeConfig]

func Init(br context.BasicRes, p plugin.PluginMeta) {

	basicRes = br
	vld = validator.New()
	connectionHelper = api.NewConnectionHelper(
		basicRes,
		vld,
		p.Name(),
	)
	params := &api.ReflectionParameters{
		ScopeIdFieldName:  "GitlabId",
		ScopeIdColumnName: "gitlab_id",
		RawScopeParamName: "Name",
	}
	scopeHelper = api.NewScopeHelper[models.GitlabConnection, models.GitlabProject, models.GitlabScopeConfig](
		basicRes,
		vld,
		connectionHelper,
		api.NewScopeDatabaseHelperImpl[models.GitlabConnection, models.GitlabProject, models.GitlabScopeConfig](
			basicRes, connectionHelper, params),
		params,
		nil,
	)
	remoteHelper = api.NewRemoteHelper[models.GitlabConnection, models.GitlabProject, models.GitlabApiProject, models.GroupResponse](
		basicRes,
		vld,
		connectionHelper,
	)
	scHelper = api.NewScopeConfigHelper[models.GitlabScopeConfig](
		basicRes,
		vld,
		p.Name(),
	)
}
