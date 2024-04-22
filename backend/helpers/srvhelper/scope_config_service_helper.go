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

package srvhelper

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

type GenericScopeConfigModelInfo[SC plugin.ToolLayerScopeConfig] struct {
	*GenericModelInfo[SC]
}

func (*GenericScopeConfigModelInfo[SC]) GetConnectionId(scopeConfig any) uint64 {
	return scopeConfig.(plugin.ToolLayerScopeConfig).ScopeConfigConnectionId()
}

func (*GenericScopeConfigModelInfo[SC]) GetScopeConfigId(scopeConfig any) uint64 {
	return scopeConfig.(plugin.ToolLayerScopeConfig).ScopeConfigId()
}

// ScopeConfigSrvHelper
type ScopeConfigSrvHelper[C plugin.ToolLayerConnection, S plugin.ToolLayerScope, SC plugin.ToolLayerScopeConfig] struct {
	*AnyScopeConfigSrvHelper
}

func NewScopeConfigSrvHelper[
	C plugin.ToolLayerConnection,
	S plugin.ToolLayerScope,
	SC plugin.ToolLayerScopeConfig,
](basicRes context.BasicRes, searchColumns []string) *ScopeConfigSrvHelper[C, S, SC] {
	return &ScopeConfigSrvHelper[C, S, SC]{
		AnyScopeConfigSrvHelper: NewAnyScopeConfigSrvHelper(
			basicRes,
			&GenericScopeConfigModelInfo[SC]{},
			&GenericScopeModelInfo[S]{},
		),
	}
}

func (scopeConfigSrv *ScopeConfigSrvHelper[C, S, SC]) GetAllByConnectionId(connectionId uint64) ([]*SC, errors.Error) {
	all, err := scopeConfigSrv.GetAllByConnectionIdAny(connectionId)
	return all.([]*SC), err
}

func (scopeConfigSrv *ScopeConfigSrvHelper[C, S, SC]) DeleteScopeConfig(scopeConfig *SC) (refs []*S, err errors.Error) {
	all, err := scopeConfigSrv.DeleteScopeConfigAny(scopeConfig)
	return all.([]*S), err
}
