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
	"fmt"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/server/services/remote/models"
)

type ScopeDatabaseHelperImpl struct {
	api.ScopeDatabaseHelper[models.RemoteConnection, models.RemoteScope, models.RemoteScopeConfig]
	pa         *pluginAPI
	db         dal.Dal
	params     *api.ReflectionParameters
	connHelper *api.ConnectionApiHelper
}

func NewScopeDatabaseHelperImpl(pa *pluginAPI, basicRes context.BasicRes, params *api.ReflectionParameters) *ScopeDatabaseHelperImpl {
	return &ScopeDatabaseHelperImpl{
		pa:         pa,
		db:         basicRes.GetDal(),
		params:     params,
		connHelper: pa.connhelper,
	}
}

func (s *ScopeDatabaseHelperImpl) VerifyConnection(connectionId uint64) errors.Error {
	conn := s.pa.connType.New()
	err := s.connHelper.FirstById(conn, connectionId)
	if err != nil {
		if s.db.IsErrorNotFound(err) {
			return errors.BadInput.New("Invalid Connection Id")
		}
		return err
	}
	return nil
}

func (s *ScopeDatabaseHelperImpl) SaveScope(scopes []*models.RemoteScope) errors.Error {
	now := time.Now()
	return s.save(scopes, &now, &now)
}

func (s *ScopeDatabaseHelperImpl) UpdateScope(connectionId uint64, scopeId string, scope *models.RemoteScope) errors.Error {
	// Update API on Gorm doesn't work with dynamic models. Need to do delete + create instead, unfortunately.
	if err := s.DeleteScope(connectionId, scopeId); err != nil {
		if !s.db.IsErrorNotFound(err) {
			return err
		}
	}
	now := time.Now()
	return s.save([]*models.RemoteScope{scope}, nil, &now)
}

func (s *ScopeDatabaseHelperImpl) GetScope(connectionId uint64, scopeId string) (*models.RemoteScope, errors.Error) {
	query := dal.Where(fmt.Sprintf("connection_id = ? AND %s = ?", s.params.ScopeIdColumnName), connectionId, scopeId)
	scope := s.pa.scopeType.New()
	err := api.CallDB(s.db.First, scope, query)
	if err != nil {
		return nil, errors.Default.Wrap(err, "could not get scope")
	}
	// @keon @camille: not sure if this is correct
	return (*models.RemoteScope)(scope.UnwrapPtr()), nil
}

func (s *ScopeDatabaseHelperImpl) ListScopes(input *plugin.ApiResourceInput, connectionId uint64) ([]*models.RemoteScope, errors.Error) {
	limit, offset := api.GetLimitOffset(input.Query, "pageSize", "page")
	scopes := s.pa.scopeType.NewSlice()
	err := api.CallDB(s.db.All, scopes, dal.Where("connection_id = ?", connectionId), dal.Limit(limit), dal.Offset(offset))
	if err != nil {
		return nil, err
	}
	var result []*models.RemoteScope
	for _, scope := range scopes.UnwrapSlice() {
		scope := scope.(models.RemoteScope)
		result = append(result, &scope)
	}
	return result, nil
}

func (s *ScopeDatabaseHelperImpl) DeleteScope(connectionId uint64, scopeId string) errors.Error {
	rawScope := s.pa.scopeType.New()
	return api.CallDB(s.db.Delete, rawScope, dal.Where("connection_id = ? AND id = ?", connectionId, scopeId))
}

func (s *ScopeDatabaseHelperImpl) GetScopeConfig(configId uint64) (*models.RemoteScopeConfig, errors.Error) {
	config := s.pa.scopeConfigType.New()
	err := api.CallDB(s.db.First, config, dal.Where("id = ?", configId))
	if err != nil {
		return nil, err
	}
	unwrapped := config.Unwrap().(models.RemoteScopeConfig)
	return &unwrapped, nil
}

func (s *ScopeDatabaseHelperImpl) ListScopeConfigs(configIds []uint64) ([]*models.RemoteScopeConfig, errors.Error) {
	configs := s.pa.scopeConfigType.NewSlice()
	err := api.CallDB(s.db.All, configs, dal.Where("id IN (?)", configIds))
	if err != nil {
		return nil, err
	}
	var result []*models.RemoteScopeConfig
	for _, config := range configs.UnwrapSlice() {
		config := config.(models.RemoteScopeConfig)
		result = append(result, &config)
	}
	return result, nil
}

func (s *ScopeDatabaseHelperImpl) save(scopes []*models.RemoteScope, createdAt *time.Time, updatedAt *time.Time) errors.Error {
	var targets []map[string]any
	for _, x := range scopes {
		ifc := reflect.ValueOf(*x).Elem().Interface()
		m, err := models.ToDatabaseMap(s.pa.scopeType.TableName(), ifc, createdAt, updatedAt)
		if err != nil {
			return err
		}
		targets = append(targets, m)
	}
	err := api.CallDB(s.db.Create, &targets, dal.From(s.pa.scopeType.TableName()))
	if err != nil {
		return errors.Default.Wrap(err, "could not save scope")
	}
	return nil
}

var _ api.ScopeDatabaseHelper[models.RemoteConnection, models.RemoteScope, models.RemoteScopeConfig] = &ScopeDatabaseHelperImpl{}
