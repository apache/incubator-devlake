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
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/plugins/slack/models"
	"github.com/go-playground/validator/v10"
)

var vld *validator.Validate
var connectionHelper *api.ConnectionApiHelper
var basicRes context.BasicRes
var dsHelper *api.DsHelper[models.SlackConnection, models.SlackChannel, srvhelper.NoScopeConfig]
var raProxy *api.DsRemoteApiProxyHelper[models.SlackConnection]
var raScopeList *api.DsRemoteApiScopeListHelper[models.SlackConnection, models.SlackChannel, SlackRemotePagination]
var raScopeSearch *api.DsRemoteApiScopeSearchHelper[models.SlackConnection, models.SlackChannel]

func Init(br context.BasicRes, p plugin.PluginMeta) {

	basicRes = br
	vld = validator.New()
	connectionHelper = api.NewConnectionHelper(
		basicRes,
		vld,
		p.Name(),
	)

	dsHelper = api.NewDataSourceHelper[
		models.SlackConnection, models.SlackChannel, srvhelper.NoScopeConfig,
	](
		br,
		p.Name(),
		[]string{"name"},
		func(c models.SlackConnection) models.SlackConnection { return c.Sanitize() },
		func(s models.SlackChannel) models.SlackChannel { return s },
		nil,
	)
	raProxy = api.NewDsRemoteApiProxyHelper[models.SlackConnection](dsHelper.ConnApi.ModelApiHelper)
	raScopeList = api.NewDsRemoteApiScopeListHelper[models.SlackConnection, models.SlackChannel, SlackRemotePagination](raProxy, listSlackRemoteScopes)
	raScopeSearch = api.NewDsRemoteApiScopeSearchHelper[models.SlackConnection, models.SlackChannel](raProxy, searchSlackRemoteScopes)
}
