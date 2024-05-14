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
	"reflect"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
)

// DsAnyHelper is a helper struct for implementing APIs for data source plugin like github/gitlab, etc.
// Normally, the plugin will have a connection model, a scope model, and a scope config model.
//   - connection: holds the data source url and credential etc
//   - scope: a scope is a collection of data source objects, like a github repo, a gitlab project, etc
//   - scope config: configuration of what to collect, how to transform, etc
//
// The helper provides APIs for CRUD operations on connection, scope, and scope config without type information
type DsAnyHelper struct {
	ConnSrv        *srvhelper.AnyConnectionSrvHelper
	ConnApi        *DsAnyConnectionApiHelper
	ScopeSrv       *srvhelper.AnyScopeSrvHelper
	ScopeApi       *DsAnyScopeApiHelper
	ScopeConfigSrv *srvhelper.AnyScopeConfigSrvHelper
	ScopeConfigApi *DsAnyScopeConfigApiHelper
}

// NewDataSourceAnyHelper creates a new DsAnyHelper
func NewDataSourceAnyHelper(
	basicRes context.BasicRes,
	pluginName string,
	scopeSearchColumns []string,
	connectionSterilizer func(c any) any,
	connModelInfo srvhelper.ConnectionModelInfo,
	scopeModelInfo srvhelper.ScopeModelInfo,
	scopeConfigModelInfo srvhelper.ScopeConfigModelInfo,
) *DsAnyHelper {
	connSrv := srvhelper.NewAnyConnectionSrvHelper(basicRes, connModelInfo, scopeModelInfo, scopeConfigModelInfo, pluginName)
	connApi := NewDsAnyConnectionApiHelper(basicRes, connSrv, connectionSterilizer)
	scopeSrv := srvhelper.NewAnyScopeSrvHelper(basicRes, scopeModelInfo, scopeConfigModelInfo, pluginName, scopeSearchColumns)
	scopeApi := NewDsAnyScopeApiHelper(basicRes, scopeSrv)
	var scopeConfigSrv *srvhelper.AnyScopeConfigSrvHelper
	var scopeConfigApi *DsAnyScopeConfigApiHelper
	if scopeConfigModelInfo != nil {
		scopeConfigSrv = srvhelper.NewAnyScopeConfigSrvHelper(basicRes, scopeConfigModelInfo, scopeModelInfo, pluginName)
		scopeConfigApi = NewDsAnyScopeConfigApiHelper(basicRes, scopeConfigSrv)
	}
	return &DsAnyHelper{
		ConnSrv:        connSrv,
		ConnApi:        connApi,
		ScopeSrv:       scopeSrv,
		ScopeApi:       scopeApi,
		ScopeConfigSrv: scopeConfigSrv,
		ScopeConfigApi: scopeConfigApi,
	}
}

var noScopeConfig = reflect.TypeOf(new(srvhelper.NoScopeConfig))

// DsHelper is the Typed version of DsAnyHelper
type DsHelper[
	C plugin.ToolLayerConnection,
	S plugin.ToolLayerScope,
	SC plugin.ToolLayerScopeConfig,
] struct {
	ConnSrv        *srvhelper.ConnectionSrvHelper[C, S, SC]
	ConnApi        *DsConnectionApiHelper[C]
	ScopeSrv       *srvhelper.ScopeSrvHelper[C, S, SC]
	ScopeApi       *DsScopeApiHelper[S]
	ScopeConfigSrv *srvhelper.ScopeConfigSrvHelper[C, S, SC]
	ScopeConfigApi *DsScopeConfigApiHelper[SC]
}

func NewDataSourceHelper[
	C plugin.ToolLayerConnection,
	S plugin.ToolLayerScope,
	SC plugin.ToolLayerScopeConfig,
](
	basicRes context.BasicRes,
	pluginName string,
	scopeSearchColumns []string,
	connectionSterilizer func(c C) C,
	scopeSterilizer func(s S) S,
	scopeConfigSterilizer func(s SC) SC,
) *DsHelper[C, S, SC] {
	var connectionModelInfo = &srvhelper.GenericConnectionModelInfo[C]{}
	var scopeModelInfo = &srvhelper.GenericScopeModelInfo[S]{}
	var scopeConfigModelInfo *srvhelper.GenericScopeConfigModelInfo[SC]
	scType := reflect.TypeOf(new(SC))
	if scType != noScopeConfig {
		scopeConfigModelInfo = &srvhelper.GenericScopeConfigModelInfo[SC]{}
	}
	anyDsHelper := NewDataSourceAnyHelper(
		basicRes,
		pluginName,
		scopeSearchColumns,
		func(c any) any {
			return connectionSterilizer(c.(C))
		},
		connectionModelInfo,
		scopeModelInfo,
		scopeConfigModelInfo,
	)
	dshelper := &DsHelper[C, S, SC]{
		ConnSrv:  srvhelper.NewConnectionSrvHelper[C, S, SC](anyDsHelper.ConnSrv),
		ConnApi:  NewDsConnectionApiHelper[C](anyDsHelper.ConnApi),
		ScopeSrv: srvhelper.NewScopeSrvHelper[C, S, SC](anyDsHelper.ScopeSrv),
		ScopeApi: NewDsScopeApiHelper[S](anyDsHelper.ScopeApi),
	}
	if anyDsHelper.ScopeConfigSrv != nil {
		dshelper.ScopeConfigSrv = srvhelper.NewScopeConfigSrvHelper[C, S, SC](anyDsHelper.ScopeConfigSrv)
	}
	if anyDsHelper.ScopeConfigApi != nil {
		dshelper.ScopeConfigApi = NewDsScopeConfigApiHelper[SC](anyDsHelper.ScopeConfigApi)
	}
	return dshelper
}
