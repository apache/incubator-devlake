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
	"fmt"
	"reflect"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

type ScopeDatabaseHelper[Conn any, Scope plugin.ToolLayerScope, Tr any] interface {
	VerifyConnection(connectionId uint64) errors.Error
	SaveScope(scopes []*Scope) errors.Error
	UpdateScope(scope *Scope) errors.Error
	GetScope(connectionId uint64, scopeId string) (*Scope, errors.Error)
	ListScopes(input *plugin.ApiResourceInput, connectionId uint64) ([]*Scope, errors.Error)

	DeleteScope(scope *Scope) errors.Error
	GetScopeConfig(ruleId uint64) (*Tr, errors.Error)
	ListScopeConfigs(ruleIds []uint64) ([]*Tr, errors.Error)
	GetScopeAndConfig(connectionId uint64, scopeId string) (*Scope, *Tr, errors.Error)
}

type ScopeDatabaseHelperImpl[Conn any, Scope plugin.ToolLayerScope, Tr any] struct {
	ScopeDatabaseHelper[Conn, Scope, Tr]
	db         dal.Dal
	connHelper *ConnectionApiHelper
	params     *ReflectionParameters
}

func NewScopeDatabaseHelperImpl[Conn any, Scope plugin.ToolLayerScope, Tr any](
	basicRes context.BasicRes, connHelper *ConnectionApiHelper, params *ReflectionParameters) *ScopeDatabaseHelperImpl[Conn, Scope, Tr] {
	return &ScopeDatabaseHelperImpl[Conn, Scope, Tr]{
		db:         basicRes.GetDal(),
		connHelper: connHelper,
		params:     params,
	}
}

func (s *ScopeDatabaseHelperImpl[Conn, Scope, Tr]) VerifyConnection(connectionId uint64) errors.Error {
	var conn Conn
	err := s.connHelper.FirstById(&conn, connectionId)
	if err != nil {
		if s.db.IsErrorNotFound(err) {
			return errors.BadInput.New("Invalid Connection Id")
		}
		return err
	}
	return nil
}

func (s *ScopeDatabaseHelperImpl[Conn, Scope, Tr]) SaveScope(scopes []*Scope) errors.Error {
	err := s.db.CreateOrUpdate(&scopes)
	if err != nil {
		if s.db.IsDuplicationError(err) {
			return errors.BadInput.New("the scope already exists")
		}
		return err
	}
	return nil
}

func (s *ScopeDatabaseHelperImpl[Conn, Scope, Tr]) UpdateScope(scope *Scope) errors.Error {
	return s.db.Update(&scope)
}

func (s *ScopeDatabaseHelperImpl[Conn, Scope, Tr]) GetScope(connectionId uint64, scopeId string) (*Scope, errors.Error) {
	query := dal.Where(fmt.Sprintf("connection_id = ? AND %s = ?", s.params.ScopeIdColumnName), connectionId, scopeId)
	scope := new(Scope)
	err := s.db.First(scope, query)
	return scope, err
}

func (s *ScopeDatabaseHelperImpl[Conn, Scope, Tr]) GetScopeAndConfig(connectionId uint64, scopeId string) (*Scope, *Tr, errors.Error) {
	scope, err := s.GetScope(connectionId, scopeId)
	if err != nil {
		return nil, nil, err
	}
	scopeConfig := new(Tr)
	scIdField := reflectField(scope, "ScopeConfigId")
	if scIdField.IsValid() {
		if scIdField.Uint() != 0 {
			err = s.db.First(scopeConfig, dal.Where("id = ?", scIdField.Uint()))
			if err != nil {
				return nil, nil, err
			}
		}
	}
	entitiesField := reflectField(scopeConfig, "Entities")
	if entitiesField.IsValid() && entitiesField.IsNil() {
		entitiesField.Set(reflect.ValueOf(plugin.DOMAIN_TYPES))
	}
	return scope, scopeConfig, nil
}

func (s *ScopeDatabaseHelperImpl[Conn, Scope, Tr]) ListScopes(input *plugin.ApiResourceInput, connectionId uint64) ([]*Scope, errors.Error) {
	searchTerm := input.Query.Get("searchTerm")
	query := dal.Where("connection_id = ?", connectionId)
	if searchTerm != "" && s.params.SearchScopeParamName != "" {
		query = dal.Where(fmt.Sprintf("connection_id = ? AND %s LIKE ?", s.params.SearchScopeParamName), connectionId, "%"+searchTerm+"%")
	}
	limit, offset := GetLimitOffset(input.Query, "pageSize", "page")
	var scopes []*Scope
	err := s.db.All(&scopes, query, dal.Limit(limit), dal.Offset(offset))
	return scopes, err
}

func (s *ScopeDatabaseHelperImpl[Conn, Scope, Tr]) DeleteScope(scope *Scope) errors.Error {
	err := s.db.Delete(&scope)
	return err
}

func (s *ScopeDatabaseHelperImpl[Conn, Scope, Tr]) GetScopeConfig(ruleId uint64) (*Tr, errors.Error) {
	var rule Tr
	err := s.db.First(&rule, dal.Where("id = ?", ruleId))
	return &rule, err
}

func (s *ScopeDatabaseHelperImpl[Conn, Scope, Tr]) ListScopeConfigs(ruleIds []uint64) ([]*Tr, errors.Error) {
	var rules []*Tr
	err := s.db.All(&rules, dal.Where("id IN (?)", ruleIds))
	return rules, err
}

var _ ScopeDatabaseHelper[any, plugin.ToolLayerScope, any] = &ScopeDatabaseHelperImpl[any, plugin.ToolLayerScope, any]{}
