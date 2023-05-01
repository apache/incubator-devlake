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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
	"github.com/go-playground/validator/v10"
)

var vld *validator.Validate
var connectionHelper *api.ConnectionApiHelper
var scopeHelper *api.ScopeApiHelper[models.BambooConnection, models.BambooProject, models.BambooTransformationRule]
var remoteHelper *api.RemoteApiHelper[models.BambooConnection, models.BambooProject, models.ApiBambooProject, api.NoRemoteGroupResponse]
var trHelper *api.TransformationRuleHelper[models.BambooTransformationRule]

var basicRes context.BasicRes

func Init(br context.BasicRes) {
	basicRes = br
	vld = validator.New()
	connectionHelper = api.NewConnectionHelper(
		basicRes,
		vld,
	)
	params := &api.ReflectionParameters{
		ScopeIdFieldName:  "ProjectKey",
		ScopeIdColumnName: "project_key",
	}
	scopeHelper = api.NewScopeHelper[models.BambooConnection, models.BambooProject, models.BambooTransformationRule](
		basicRes,
		vld,
		connectionHelper,
		api.NewScopeDatabaseHelperImpl[models.BambooConnection, models.BambooProject, models.BambooTransformationRule](
			basicRes, connectionHelper, params),
		params,
	)
	remoteHelper = api.NewRemoteHelper[models.BambooConnection, models.BambooProject, models.ApiBambooProject, api.NoRemoteGroupResponse](
		basicRes,
		vld,
		connectionHelper,
	)
	trHelper = api.NewTransformationRuleHelper[models.BambooTransformationRule](
		basicRes,
		vld,
	)
}
