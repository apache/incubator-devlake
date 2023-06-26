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
	"strconv"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/go-playground/validator/v10"
)

var vld *validator.Validate
var connectionHelper *api.ConnectionApiHelper
var scopeHelper *api.ScopeApiHelper[models.GithubConnection, models.GithubRepo, models.GithubScopeConfig]
var basicRes context.BasicRes
var scHelper *api.ScopeConfigHelper[models.GithubScopeConfig]
var remoteHelper *api.RemoteApiHelper[models.GithubConnection, models.GithubRepo, repo, plugin.ApiGroup]

func Init(br context.BasicRes, p plugin.PluginMeta) {

	basicRes = br
	vld = validator.New()
	connectionHelper = api.NewConnectionHelper(
		basicRes,
		vld,
		p.Name(),
	)
	params := &api.ReflectionParameters{
		ScopeIdFieldName:  "GithubId",
		ScopeIdColumnName: "github_id",
		RawScopeParamName: "Name",
	}
	scopeHelper = api.NewScopeHelper[models.GithubConnection, models.GithubRepo, models.GithubScopeConfig](
		basicRes,
		vld,
		connectionHelper,
		api.NewScopeDatabaseHelperImpl[models.GithubConnection, models.GithubRepo, models.GithubScopeConfig](
			basicRes, connectionHelper, params),
		params,
		&api.ScopeHelperOptions{
			GetScopeParamValue: func(db dal.Dal, scopeId string) (string, errors.Error) {
				id, err := errors.Convert01(strconv.ParseInt(scopeId, 10, 64))
				if err != nil {
					return "", err
				}
				repo := models.GithubRepo{
					GithubId: int(id),
				}
				err = db.First(&repo)
				if err != nil {
					return "", err
				}
				return repo.FullName, nil
			},
		},
	)
	scHelper = api.NewScopeConfigHelper[models.GithubScopeConfig](
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
