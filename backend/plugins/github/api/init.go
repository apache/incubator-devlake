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
	"github.com/apache/incubator-devlake/core/plugin"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/go-playground/validator/v10"
)

var vld *validator.Validate

// var connSrv *srvhelper.ConnectionSrvHelper[models.GithubConnection, models.GithubRepo, models.GithubScopeConfig]
var connApi *api.DsConnectionApiHelper[models.GithubConnection, models.GithubRepo, models.GithubScopeConfig]
var scopeSrv *srvhelper.ScopeSrvHelper[models.GithubConnection, models.GithubRepo, models.GithubScopeConfig]
var scopeApi *api.DsScopeApiHelper[models.GithubConnection, models.GithubRepo, models.GithubScopeConfig]
var scSrv *srvhelper.ScopeConfigSrvHelper[models.GithubConnection, models.GithubRepo, models.GithubScopeConfig]
var scApi *api.DsScopeConfigApiHelper[models.GithubConnection, models.GithubRepo, models.GithubScopeConfig]
var connectionHelper *api.ConnectionApiHelper
var basicRes context.BasicRes
var remoteHelper *api.RemoteApiHelper[models.GithubConnection, models.GithubRepo, repo, plugin.ApiGroup]

func Init(br context.BasicRes, p plugin.PluginMeta) {
	basicRes = br
	_, connApi, scopeSrv, scopeApi, scSrv, scApi = api.NewDataSourceHelpers[
		models.GithubConnection,
		models.GithubRepo, models.GithubScopeConfig,
	](
		br,
		p.Name(),
		[]string{"full_name"},
	)
	// TODO: refactor remoteHelper
	vld = validator.New()
	connectionHelper = api.NewConnectionHelper(
		basicRes,
		vld,
		p.Name(),
	)
	remoteHelper = api.NewRemoteHelper[models.GithubConnection, models.GithubRepo, repo, plugin.ApiGroup](
		basicRes,
		vld,
		connectionHelper,
	)
}
