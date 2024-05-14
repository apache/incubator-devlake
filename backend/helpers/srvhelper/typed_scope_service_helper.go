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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
)

type GenericScopeModelInfo[S plugin.ToolLayerScope] struct {
	*GenericModelInfo[S]
}

func (*GenericScopeModelInfo[S]) GetScopeId(scope any) string {
	return scope.(plugin.ToolLayerScope).ScopeId()
}

func (*GenericScopeModelInfo[S]) GetConnectionId(scope any) uint64 {
	return scope.(plugin.ToolLayerScope).ScopeConnectionId()
}

func (*GenericScopeModelInfo[S]) GetScopeConfigId(scope any) uint64 {
	return scope.(plugin.ToolLayerScope).ScopeScopeConfigId()
}

func (*GenericScopeModelInfo[S]) GetScopeParams(scope any) any {
	return scope.(plugin.ToolLayerScope).ScopeScopeConfigId()
}

func NewScopeModelInfo[S plugin.ToolLayerScope]() *GenericScopeModelInfo[S] {
	return &GenericScopeModelInfo[S]{NewGenericModelInfo[S]()}
}

type ScopeDetail[S plugin.ToolLayerScope, SC plugin.ToolLayerScopeConfig] struct {
	Scope       S                   `json:"scope"`
	ScopeConfig *SC                 `json:"scopeConfig,omitempty"`
	Blueprints  []*models.Blueprint `json:"blueprints,omitempty"`
}

type ScopeSrvHelper[C plugin.ToolLayerConnection, S plugin.ToolLayerScope, SC plugin.ToolLayerScopeConfig] struct {
	*AnyScopeSrvHelper
}

// NewScopeSrvHelper creates a ScopeDalHelper for scope management
func NewScopeSrvHelper[
	C plugin.ToolLayerConnection,
	S plugin.ToolLayerScope,
	SC plugin.ToolLayerScopeConfig,
](
	anyScopeSrv *AnyScopeSrvHelper,
) *ScopeSrvHelper[C, S, SC] {
	return &ScopeSrvHelper[C, S, SC]{
		AnyScopeSrvHelper: anyScopeSrv,
	}
}

func (scopeSrv *ScopeSrvHelper[C, S, SC]) GetScopeDetail(includeBlueprints bool, pkv ...interface{}) (*ScopeDetail[S, SC], errors.Error) {
	return scopeSrv.mapScopeDetail(scopeSrv.GetScopeDetailAny(includeBlueprints, pkv...))
}

// MapScopeDetails returns scope details (scope and scopeConfig) for the given blueprint scopes
func (scopeSrv *ScopeSrvHelper[C, S, SC]) MapScopeDetails(connectionId uint64, bpScopes []*models.BlueprintScope) ([]*ScopeDetail[S, SC], errors.Error) {
	return scopeSrv.mapScopeDetails(scopeSrv.MapScopeDetailsAny(connectionId, bpScopes))
}

func (scopeSrv *ScopeSrvHelper[C, S, SC]) GetScopesPage(pagination *ScopePagination) ([]*ScopeDetail[S, SC], int64, errors.Error) {
	scopes, count, err := scopeSrv.GetScopesPageAny(pagination)
	scopeDetails, err := scopeSrv.mapScopeDetails(scopes, err)
	return scopeDetails, count, err
}

func (scopeSrv *ScopeSrvHelper[C, S, SC]) mapScopeDetail(detail *AnyScopeDetail, err errors.Error) (*ScopeDetail[S, SC], errors.Error) {
	if err != nil {
		return nil, err
	}
	return &ScopeDetail[S, SC]{
		Scope:       detail.Scope.(S),
		ScopeConfig: detail.ScopeConfig.(*SC),
	}, nil
}

func (scopeSrv *ScopeSrvHelper[C, S, SC]) mapScopeDetails(details []*AnyScopeDetail, err errors.Error) ([]*ScopeDetail[S, SC], errors.Error) {
	if err != nil {
		return nil, err
	}
	scopeDetails := make([]*ScopeDetail[S, SC], len(details))
	for i, detail := range details {
		scopeDetails[i] = &ScopeDetail[S, SC]{
			Scope:       detail.Scope.(S),
			ScopeConfig: detail.ScopeConfig.(*SC),
		}
	}
	return scopeDetails, err
}

func (scopeSrv *ScopeSrvHelper[C, S, SC]) DeleteScope(scope *S, dataOnly bool) (refs *DsRefs, err errors.Error) {
	return scopeSrv.DeleteScopeAny(scope, dataOnly)
}
