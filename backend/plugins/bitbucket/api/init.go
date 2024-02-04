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
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/go-playground/validator/v10"
)

var basicRes context.BasicRes
var vld *validator.Validate
var dsHelper *api.DsHelper[models.BitbucketConnection, models.BitbucketRepo, models.BitbucketScopeConfig]
var raProxy *api.DsRemoteApiProxyHelper[models.BitbucketConnection]
var raScopeList *api.DsRemoteApiScopeListHelper[models.BitbucketConnection, models.BitbucketRepo, BitbucketRemotePagination]
var raScopeSearch *api.DsRemoteApiScopeSearchHelper[models.BitbucketConnection, models.BitbucketRepo]

func Init(br context.BasicRes, p plugin.PluginMeta) {

	basicRes = br

	dsHelper = api.NewDataSourceHelper[
		models.BitbucketConnection, models.BitbucketRepo, models.BitbucketScopeConfig,
	](
		br,
		p.Name(),
		[]string{"name"},
		func(c models.BitbucketConnection) models.BitbucketConnection {
			return c.Sanitize()
		},
		nil,
		nil,
	)
	raProxy = api.NewDsRemoteApiProxyHelper[models.BitbucketConnection](dsHelper.ConnApi.ModelApiHelper)
	raScopeList = api.NewDsRemoteApiScopeListHelper[models.BitbucketConnection, models.BitbucketRepo, BitbucketRemotePagination](raProxy, listBitbucketRemoteScopes)
	raScopeSearch = api.NewDsRemoteApiScopeSearchHelper[models.BitbucketConnection, models.BitbucketRepo](raProxy, searchBitbucketRepos)

}
