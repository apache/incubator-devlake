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
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/domainlayer/domaininfo"
	"github.com/apache/incubator-devlake/core/plugin"
)

type ScopePagination struct {
	Pagination   `mapstructure:",squash"`
	ConnectionId uint64 `json:"connectionId" mapstructure:"connectionId" validate:"required"`
	Blueprints   bool   `json:"blueprints" mapstructure:"blueprints"`
}

type AnyScopeDetail struct {
	Scope       any                 `json:"scope"`
	ScopeConfig any                 `json:"scopeConfig,omitempty"`
	Blueprints  []*models.Blueprint `json:"blueprints,omitempty"`
}

type ScopeModelInfo interface {
	ModelInfo
	GetScopeId(any) string
	GetConnectionId(any) uint64
	GetScopeConfigId(any) uint64
	GetScopeParams(any) interface{}
}

type AnyScopeSrvHelper struct {
	ScopeModelInfo
	ScopeConfigModelInfo
	*AnyModelSrvHelper
	pluginName string
}

func NewAnyScopeSrvHelper(
	basicRes context.BasicRes,
	scopeModelInfo ScopeModelInfo,
	scopeConfigModelInfo ScopeConfigModelInfo,
	pluginName string,
	searchColumns []string,
) *AnyScopeSrvHelper {
	return &AnyScopeSrvHelper{
		AnyModelSrvHelper:    NewAnyModelSrvHelper(basicRes, scopeModelInfo, searchColumns),
		pluginName:           pluginName,
		ScopeModelInfo:       scopeModelInfo,
		ScopeConfigModelInfo: scopeConfigModelInfo,
	}
}

func (scopeSrv *AnyScopeSrvHelper) GetPluginName() string {
	return scopeSrv.pluginName
}

func (scopeSrv *AnyScopeSrvHelper) Validate(scope any) errors.Error {
	connectionId := scopeSrv.ScopeModelInfo.GetConnectionId(scope)
	connectionCount := errors.Must1(scopeSrv.db.Count(dal.From(scopeSrv.ScopeModelInfo.TableName()), dal.Where("id = ?", connectionId)))
	if connectionCount == 0 {
		return errors.BadInput.New("connectionId is invalid")
	}
	scopeConfigId := scopeSrv.ScopeModelInfo.GetScopeConfigId(scope)
	scopeConfigCount := errors.Must1(scopeSrv.db.Count(dal.From(scopeSrv.ScopeModelInfo.TableName()), dal.Where("id = ?", scopeConfigId)))
	if scopeConfigCount == 0 {
		return errors.BadInput.New("scopeConfigId is invalid")
	}
	return nil
}

func (scopeSrv *AnyScopeSrvHelper) GetScopeDetailAny(includeBlueprints bool, pkv ...interface{}) (*AnyScopeDetail, errors.Error) {
	scope, err := scopeSrv.FindByPkAny(pkv...)
	if err != nil {
		return nil, err
	}
	scopeConfigId := scopeSrv.ScopeModelInfo.GetScopeConfigId(scope)
	scopeDetail := &AnyScopeDetail{
		Scope:       scope,
		ScopeConfig: scopeSrv.getScopeConfig(scopeConfigId),
	}
	if includeBlueprints {
		connectionId := scopeSrv.ScopeModelInfo.GetConnectionId(scope)
		scopeId := scopeSrv.ScopeModelInfo.GetScopeId(scope)
		scopeDetail.Blueprints = scopeSrv.getAllBlueprinsByScope(connectionId, scopeId)
	}
	return scopeDetail, nil
}

func (scopeSrv *AnyScopeSrvHelper) GetScopeLatestSyncStateAny(pkv ...interface{}) ([]*models.LatestSyncState, errors.Error) {
	scope, err := scopeSrv.FindByPkAny(pkv...)
	if err != nil {
		return nil, err
	}
	params := plugin.MarshalScopeParams(scopeSrv.ScopeModelInfo.GetScopeParams(scope))
	scopeSrv.log.Debug("scope: %#+v, params: %+v", scope, params)
	scopeSyncStates := []*models.LatestSyncState{}
	if err := scopeSrv.db.All(
		&scopeSyncStates,
		dal.Select("raw_data_table, latest_success_start, raw_data_params"),
		dal.From("_devlake_collector_latest_state"),
		dal.Where("raw_data_params = ?", params),
	); err != nil {
		return nil, err
	}
	scopeSrv.log.Debug("param: %+v, resp: %+v", scopeSyncStates)
	sort.Slice(scopeSyncStates, func(i, j int) bool {
		if scopeSyncStates[i].LatestSuccessStart != nil && scopeSyncStates[j].LatestSuccessStart != nil {
			return scopeSyncStates[i].LatestSuccessStart.After(*scopeSyncStates[j].LatestSuccessStart)
		}
		return false
	})
	return scopeSyncStates, nil
}

// MapScopeDetails returns scope details (scope and scopeConfig) for the given blueprint scopes
func (scopeSrv *AnyScopeSrvHelper) MapScopeDetailsAny(connectionId uint64, bpScopes []*models.BlueprintScope) ([]*AnyScopeDetail, errors.Error) {
	var err errors.Error
	scopeDetails := make([]*AnyScopeDetail, len(bpScopes))
	for i, bpScope := range bpScopes {
		scopeDetails[i], err = scopeSrv.GetScopeDetailAny(false, connectionId, bpScope.ScopeId)
		if err != nil {
			return nil, err
		}
		if scopeDetails[i].ScopeConfig == nil {
			scopeDetails[i].ScopeConfig = scopeSrv.ScopeConfigModelInfo.New()
		}
		setDefaultEntities(scopeDetails[i].ScopeConfig)
	}
	return scopeDetails, nil
}

func (scopeSrv *AnyScopeSrvHelper) GetScopesPageAny(pagination *ScopePagination) ([]*AnyScopeDetail, int64, errors.Error) {
	if pagination.ConnectionId < 1 {
		return nil, 0, errors.BadInput.New("connectionId is required")
	}
	scopes, count, err := scopeSrv.QueryPageAny(
		&pagination.Pagination,
		dal.Where("connection_id = ?", pagination.ConnectionId),
	)
	if err != nil {
		return nil, 0, err
	}
	slice := reflect.ValueOf(scopes)
	data := make([]*AnyScopeDetail, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		// load blueprints
		scope := slice.Index(i).Interface()
		scopeConfigId := scopeSrv.ScopeModelInfo.GetScopeConfigId(scope)
		scopeDetail := &AnyScopeDetail{
			Scope:       scope,
			ScopeConfig: scopeSrv.getScopeConfig(scopeConfigId),
		}
		if pagination.Blueprints {
			connectionId := scopeSrv.ScopeModelInfo.GetConnectionId(scope)
			scopeId := scopeSrv.ScopeModelInfo.GetScopeId(scope)
			scopeDetail.Blueprints = scopeSrv.getAllBlueprinsByScope(connectionId, scopeId)
		}
		data[i] = scopeDetail
	}
	return data, count, nil
}

func (scopeSrv *AnyScopeSrvHelper) DeleteScopeAny(scope any, dataOnly bool) (refs *DsRefs, err errors.Error) {
	err = scopeSrv.NoRunningPipeline(func(tx dal.Transaction) errors.Error {
		// check referencing blueprints
		if !dataOnly {
			connectionId := scopeSrv.ScopeModelInfo.GetConnectionId(scope)
			scopeId := scopeSrv.ScopeModelInfo.GetScopeId(scope)
			refs = toDsRefs(scopeSrv.getAllBlueprinsByScope(connectionId, scopeId))
			if refs != nil {
				return errors.Conflict.New("Cannot delete the scope because it is referenced by blueprints")
			}
			errors.Must(tx.Delete(scope))
		}
		// delete data
		scopeSrv.deleteScopeData(scope, tx)
		return nil
	})
	return
}

func (scopeSrv *AnyScopeSrvHelper) getScopeConfig(scopeConfigId uint64) any {
	if scopeConfigId < 1 {
		return nil
	}
	scopeConfig := scopeSrv.ScopeConfigModelInfo.New()
	err := scopeSrv.db.First(
		&scopeConfig,
		dal.Where(
			"id = ?",
			scopeConfigId,
		),
	)
	if err != nil {
		return nil
	}
	return &scopeConfig
}

func (scopeSrv *AnyScopeSrvHelper) getAllBlueprinsByScope(connectionId uint64, scopeId string) []*models.Blueprint {
	blueprints := make([]*models.Blueprint, 0)
	errors.Must(scopeSrv.db.All(
		&blueprints,
		dal.From("_devlake_blueprints bp"),
		dal.Join("JOIN _devlake_blueprint_scopes sc ON sc.blueprint_id = bp.id"),
		dal.Where(
			"mode = ? AND sc.connection_id = ? AND sc.plugin_name = ? AND sc.scope_id = ?",
			"NORMAL",
			connectionId,
			scopeSrv.pluginName,
			scopeId,
		),
	))
	return blueprints
}

func (scopeSrv *AnyScopeSrvHelper) deleteScopeData(scope any, tx dal.Transaction) {
	rawDataParams := plugin.MarshalScopeParams(scopeSrv.ScopeModelInfo.GetScopeParams(scope))
	generateWhereClause := func(table string) (string, []any) {
		var where string
		var params []interface{}
		if strings.HasPrefix(table, "_raw_") {
			// raw table: should check connection and scope
			where = "params = ?"
			params = []interface{}{rawDataParams}
		} else if strings.HasPrefix(table, "_tool_") {
			// tool layer table: should check connection and scope
			where = "_raw_data_params = ?"
			params = []interface{}{rawDataParams}
		} else {
			// framework tables: should check plugin, connection and scope
			if table == (models.CollectorLatestState{}.TableName()) {
				// diff sync state
				where = "raw_data_table LIKE ? AND raw_data_params = ?"
			} else {
				// domain layer table
				where = "_raw_data_table LIKE ? AND _raw_data_params = ?"
			}
			rawDataTablePrefix := fmt.Sprintf("_raw_%s%%", scopeSrv.pluginName)
			params = []interface{}{rawDataTablePrefix, rawDataParams}
		}
		return where, params
	}
	tables := errors.Must1(scopeSrv.getAffectedTables())
	for _, table := range tables {
		where, params := generateWhereClause(table)
		scopeSrv.log.Info("deleting data from table %s with WHERE \"%s\" and params: \"%v\"", table, where, params)
		sql := fmt.Sprintf("DELETE FROM %s WHERE %s", table, where)
		errors.Must(tx.Exec(sql, params...))
	}
}

func (scopeSrv *AnyScopeSrvHelper) getAffectedTables() ([]string, errors.Error) {
	var tables []string
	meta, err := plugin.GetPlugin(scopeSrv.pluginName)
	if err != nil {
		return nil, err
	}
	if pluginModel, ok := meta.(plugin.PluginModel); !ok {
		panic(errors.Default.New(fmt.Sprintf("plugin \"%s\" does not implement listing its tables", scopeSrv.pluginName)))
	} else {
		// Unfortunately, can't cache the tables because Python creates some tables on a per-demand basis, so such a cache would possibly get outdated.
		// It's a rare scenario in practice, but might as well play it safe and sacrifice some performance here
		var allTables []string
		if allTables, err = scopeSrv.db.AllTables(); err != nil {
			return nil, err
		}
		// collect raw tables
		for _, table := range allTables {
			if strings.HasPrefix(table, "_raw_"+scopeSrv.pluginName) {
				tables = append(tables, table)
			}
		}
		// collect tool tables
		toolModels := pluginModel.GetTablesInfo()
		for _, toolModel := range toolModels {
			if !isScopeModel(toolModel) && hasField(toolModel, "RawDataParams") {
				tables = append(tables, toolModel.TableName())
			}
		}
		// collect domain tables
		for _, domainModel := range domaininfo.GetDomainTablesInfo() {
			// we only care about tables with RawOrigin
			ok = hasField(domainModel, "RawDataParams")
			if ok {
				tables = append(tables, domainModel.TableName())
			}
		}
		// additional tables
		tables = append(tables, models.CollectorLatestState{}.TableName())
	}
	scopeSrv.log.Debug("Discovered %d tables used by plugin \"%s\": %v", len(tables), scopeSrv.pluginName, tables)
	return tables, nil
}

// TODO: sort out the follow functions
func isScopeModel(obj dal.Tabler) bool {
	_, ok := obj.(plugin.ToolLayerScope)
	return ok
}

func hasField(obj any, fieldName string) bool {
	obj = models.UnwrapObject(obj)
	_, ok := reflectType(obj).FieldByName(fieldName)
	return ok
}

func reflectType(obj any) reflect.Type {
	obj = models.UnwrapObject(obj)
	typ := reflect.TypeOf(obj)
	kind := typ.Kind()
	for kind == reflect.Ptr {
		typ = typ.Elem()
		kind = typ.Kind()
	}
	return typ
}

func setDefaultEntities(sc interface{}) {
	v := reflect.ValueOf(sc)
	if v.Kind() != reflect.Pointer {
		panic(fmt.Errorf("sc must be a pointer"))
	}
	entities := v.Elem().FieldByName("Entities")
	if !entities.IsValid() ||
		!(entities.Kind() == reflect.Array || entities.Kind() == reflect.Slice) ||
		entities.Type().Elem().Kind() != reflect.String {
		return
	}
	if entities.IsNil() || entities.Len() == 0 {
		entities.Set(reflect.ValueOf(plugin.DOMAIN_TYPES))
	}
}
