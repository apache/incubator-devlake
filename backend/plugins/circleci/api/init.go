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
	"github.com/apache/incubator-devlake/plugins/circleci/models"
	"github.com/go-playground/validator/v10"
)

var basicRes context.BasicRes
var vld *validator.Validate

var dsHelper *api.DsHelper[models.CircleciConnection, models.CircleciProject, models.CircleciScopeConfig]
var raProxy *api.DsRemoteApiProxyHelper[models.CircleciConnection]
var raScopeList *api.DsRemoteApiScopeListHelper[models.CircleciConnection, models.CircleciProject, srvhelper.NoPagintation]

// var raScopeSearch *api.DsRemoteApiScopeSearchHelper[models.CircleciConnection, models.CircleciProject]

func Init(br context.BasicRes, p plugin.PluginMeta) {
	basicRes = br
	vld = validator.New()
	dsHelper = api.NewDataSourceHelper[
		models.CircleciConnection, models.CircleciProject, models.CircleciScopeConfig,
	](
		br,
		p.Name(),
		[]string{"name"},
		func(c models.CircleciConnection) models.CircleciConnection {
			return c.Sanitize()
		},
		nil,
		nil,
	)
	raProxy = api.NewDsRemoteApiProxyHelper[models.CircleciConnection](dsHelper.ConnApi.ModelApiHelper)
	raScopeList = api.NewDsRemoteApiScopeListHelper[models.CircleciConnection, models.CircleciProject, srvhelper.NoPagintation](raProxy, listCircleciRemoteScopes)
	// raScopeSearch = api.NewDsRemoteApiScopeSearchHelper[models.CircleciConnection, models.CircleciProject](raProxy, searchCircleciProjects)
}
