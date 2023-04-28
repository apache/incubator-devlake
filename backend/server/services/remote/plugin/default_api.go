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

package plugin

import (
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/impls/dalgorm"
	"github.com/apache/incubator-devlake/server/services/remote/bridge"
	"gorm.io/gorm"
	"reflect"
)

type pluginAPI struct {
	invoker    bridge.Invoker
	connType   *models.DynamicTabler
	txRuleType *models.DynamicTabler
	scopeType  *models.DynamicTabler
	helper     *api.ConnectionApiHelper
}

func GetDefaultAPI(
	invoker bridge.Invoker,
	connType *models.DynamicTabler,
	txRuleType *models.DynamicTabler,
	scopeType *models.DynamicTabler,
	helper *api.ConnectionApiHelper,
) map[string]map[string]plugin.ApiResourceHandler {
	papi := &pluginAPI{
		invoker:    invoker,
		connType:   connType,
		txRuleType: txRuleType,
		scopeType:  scopeType,
		helper:     helper,
	}

	resources := map[string]map[string]plugin.ApiResourceHandler{
		"test": {
			"POST": papi.TestConnection,
		},
		"connections": {
			"POST": papi.PostConnections,
			"GET":  papi.ListConnections,
		},
		"connections/:connectionId": {
			"GET":    papi.GetConnection,
			"PATCH":  papi.PatchConnection,
			"DELETE": papi.DeleteConnection,
		},
		"connections/:connectionId/scopes": {
			"PUT": papi.PutScope,
			"GET": papi.ListScopes,
		},
		"connections/:connectionId/scopes/:scopeId": {
			"GET":    papi.GetScope,
			"PATCH":  papi.PatchScope,
			"DELETE": papi.DeleteScope,
		},
		"connections/:connectionId/remote-scopes": {
			"GET": papi.GetRemoteScopes,
		},
		"connections/:connectionId/search-remote-scopes": {
			"GET": papi.SearchRemoteScopes,
		},
	}

	if txRuleType != nil {
		resources["connections/:connectionId/transformation_rules"] = map[string]plugin.ApiResourceHandler{
			"POST": papi.PostTransformationRules,
			"GET":  papi.ListTransformationRules,
		}
		resources["connections/:connectionId/transformation_rules/:id"] = map[string]plugin.ApiResourceHandler{
			"GET":   papi.GetTransformationRule,
			"PATCH": papi.PatchTransformationRule,
		}
	}
	scopeHelper = createScopeHelper(papi)
	return resources
}

func createScopeHelper(pa *pluginAPI) *api.GenericScopeHelper[any, any] {
	db := basicRes.GetDal()
	params := &api.ReflectionParameters{
		ScopeIdFieldName:  "Id",
		ScopeIdColumnName: "id",
		RawScopeParamName: "scope_id",
	}
	return api.NewGenericScopeHelper[any, any](
		basicRes,
		params,
		nil,
		func(connectionId uint64) errors.Error {
			connection := pa.connType.New()
			err := connectionHelper.FirstById(connection, connectionId)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.BadInput.New("Invalid Connection Id")
				}
				return err
			}
			return nil
		},
		func(scopes []*any) errors.Error {
			var targets []map[string]any
			for _, x := range scopes {
				ifc := reflect.ValueOf(*x).Elem().Interface()
				j, err := errors.Convert01(json.Marshal(ifc))
				if err != nil {
					return err
				}
				m := map[string]any{}
				err = errors.Convert(json.Unmarshal(j, &m))
				if err != nil {
					return err
				}
				m = dalgorm.ToDatabaseMap(pa.scopeType.TableName(), m) //or use api.DecodeMapStruct?
				targets = append(targets, m)
			}
			err := api.CallDB(db.Create, &targets, dal.From(pa.scopeType.TableName()))
			if err != nil {
				if db.IsDuplicationError(err) {
					return errors.BadInput.Wrap(err, "the scope already exists")
				}
				return err
			}
			return nil
		},
		func(connectionId uint64, scopeId string) (any, errors.Error) {
			query := dal.Where(fmt.Sprintf("connection_id = ? AND %s = ?", params.ScopeIdColumnName), connectionId, scopeId)
			scope := pa.scopeType.New()
			err := api.CallDB(db.First, scope, query)
			if basicRes.GetDal().IsErrorNotFound(err) {
				return scope, errors.NotFound.New("Scope not found")
			}
			return scope.Unwrap(), nil
		},
		func(input *plugin.ApiResourceInput, connectionId uint64) ([]*any, errors.Error) {
			limit, offset := api.GetLimitOffset(input.Query, "pageSize", "page")
			scopes := pa.scopeType.NewSlice()
			err := api.CallDB(db.All, scopes, dal.Where("connection_id = ?", connectionId), dal.Limit(limit), dal.Offset(offset))
			if err != nil {
				return nil, err
			}
			var result []*any
			for _, scope := range scopes.UnwrapSlice() {
				result = append(result, &scope)
			}
			return result, nil
		},
		func(connectionId uint64, scopeId string) errors.Error {
			rawScope := pa.scopeType.New()
			return api.CallDB(db.Delete, rawScope, dal.Where("connection_id = ? AND id = ?", connectionId, scopeId))
		},
		func(ruleId uint64) (any, errors.Error) {
			rule := pa.txRuleType.New()
			err := api.CallDB(db.First, rule, dal.Where("id = ?", ruleId))
			if err != nil {
				return rule, errors.NotFound.New("transformationRule not found")
			}
			return rule.Unwrap(), nil
		},
		func(ruleIds []uint64) ([]*any, errors.Error) {
			rules := pa.txRuleType.NewSlice()
			err := api.CallDB(db.All, rules, dal.Where("id IN (?)", ruleIds))
			if err != nil {
				return nil, err
			}
			var result []*any
			for _, scope := range rules.UnwrapSlice() {
				result = append(result, &scope)
			}
			return result, nil
		},
	)
}
