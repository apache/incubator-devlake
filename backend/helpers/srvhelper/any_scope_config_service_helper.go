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
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
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
	pluginName string
}

func NewAnyScopeConfigSrvHelper(
	basicRes context.BasicRes,
	scopeConfigModelInfo ScopeConfigModelInfo,
	scopeModelInfo ScopeModelInfo,
	pluginName string,
) *AnyScopeConfigSrvHelper {
	return &AnyScopeConfigSrvHelper{
		ScopeConfigModelInfo: scopeConfigModelInfo,
		ScopeModelInfo:       scopeModelInfo,
		AnyModelSrvHelper:    NewAnyModelSrvHelper(basicRes, scopeConfigModelInfo, nil),
		pluginName:           pluginName,
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

func (scopeConfigSrv *AnyScopeConfigSrvHelper) GetProjectsByScopeConfig(scopeConfig any) (*models.ProjectScopeOutput, errors.Error) {
	ps := &models.ProjectScopeOutput{}
	projectMap := make(map[string]*models.ProjectScope)
	// 1. get all scopes that are using the scopeConfigId
	scopes := scopeConfigSrv.ScopeModelInfo.NewSlice()
	sc := scopeConfig.(plugin.ToolLayerScopeConfig)
	err := scopeConfigSrv.db.All(&scopes,
		dal.Where("scope_config_id = ?", sc.ScopeConfigId()),
	)
	if err != nil {
		return nil, err
	}
	slice := reflect.ValueOf(scopes)
	for i := 0; i < slice.Len(); i++ {
		// 2. get blueprint id by connection id and scope id
		bpScope := []*models.BlueprintScope{}
		s := slice.Index(i).Interface().(plugin.ToolLayerScope)
		err = scopeConfigSrv.db.All(&bpScope,
			dal.Where("plugin_name = ? and connection_id = ? and scope_id = ?", scopeConfigSrv.pluginName, s.ScopeConnectionId(), s.ScopeId()),
		)
		if err != nil {
			return nil, err
		}

		for _, bs := range bpScope {
			// 3. get project details by blueprint id
			bp := models.Blueprint{}
			err = scopeConfigSrv.db.All(&bp,
				dal.Where("id = ?", bs.BlueprintId),
			)
			if err != nil {
				return nil, err
			}
			if project, exists := projectMap[bp.ProjectName]; exists {
				project.Scopes = append(project.Scopes, struct {
					ScopeID   string `json:"scopeId"`
					ScopeName string `json:"scopeName"`
				}{
					ScopeID:   bs.ScopeId,
					ScopeName: s.ScopeName(),
				})
			} else {
				projectMap[bp.ProjectName] = &models.ProjectScope{
					Name:        bp.ProjectName,
					BlueprintId: bp.ID,
					Scopes: []struct {
						ScopeID   string `json:"scopeId"`
						ScopeName string `json:"scopeName"`
					}{
						{
							ScopeID:   bs.ScopeId,
							ScopeName: s.ScopeName(),
						},
					},
				}
			}
		}
	}
	// 4. combine all projects
	for _, project := range projectMap {
		ps.Projects = append(ps.Projects, *project)
	}
	ps.Count = len(ps.Projects)

	return ps, err
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
