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
	"reflect"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
)

// NoScopeConfig is a placeholder for plugins that don't have any scope configuration yet
type NoScopeConfig struct{}

func (NoScopeConfig) TableName() string               { return "" }
func (NoScopeConfig) ScopeConfigId() uint64           { return 0 }
func (NoScopeConfig) ScopeConfigConnectionId() uint64 { return 0 }

type ScopeConfigModelInfo interface {
	ModelInfo
	GetConnectionId(any) uint64
	GetScopeConfigId(any) uint64
}

// ScopeConfigSrvHelper
type AnyScopeConfigSrvHelper struct {
	ScopeConfigModelInfo
	ScopeModelInfo
	*AnyModelSrvHelper
}

func NewAnyScopeConfigSrvHelper(
	basicRes context.BasicRes,
	scopeConfigModelInfo ScopeConfigModelInfo,
	scopeModelInfo ScopeModelInfo,
) *AnyScopeConfigSrvHelper {
	return &AnyScopeConfigSrvHelper{
		ScopeConfigModelInfo: scopeConfigModelInfo,
		ScopeModelInfo:       scopeModelInfo,
		AnyModelSrvHelper:    NewAnyModelSrvHelper(basicRes, scopeConfigModelInfo, nil),
	}
}

func (scopeConfigSrv *AnyScopeConfigSrvHelper) GetAllByConnectionIdAny(connectionId uint64) (any, errors.Error) {
	scopeConfigs := scopeConfigSrv.ScopeConfigModelInfo.NewSlice()
	err := scopeConfigSrv.db.All(&scopeConfigs,
		dal.Where("connection_id = ?", connectionId),
		dal.Orderby("id DESC"),
	)
	return scopeConfigs, err
}

func (scopeConfigSrv *AnyScopeConfigSrvHelper) DeleteScopeConfigAny(scopeConfig any) (refs any, err errors.Error) {
	err = scopeConfigSrv.NoRunningPipeline(func(tx dal.Transaction) errors.Error {
		// make sure no scope is using the scopeConfig
		connectionId := scopeConfigSrv.ScopeConfigModelInfo.GetConnectionId(scopeConfig)
		scopeConfigId := scopeConfigSrv.ScopeConfigModelInfo.GetScopeConfigId(scopeConfig)
		refs = scopeConfigSrv.ScopeModelInfo.NewSlice()
		errors.Must(tx.All(
			&refs,
			dal.Where("connection_id = ? AND scope_config_id = ?", connectionId, scopeConfigId),
		))
		if reflect.ValueOf(refs).Len() > 0 {
			return errors.Conflict.New("Please delete all data scope(s) before you delete this ScopeConfig.")
		}
		errors.Must(tx.Delete(scopeConfig))
		return nil
	})
	return
}
