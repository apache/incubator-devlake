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
	"net/http"
	"reflect"
	"strconv"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/go-playground/validator/v10"
)

// ScopeConfigHelper is used to write the CURD of scope config
type ScopeConfigHelper[ScopeConfig dal.Tabler] struct {
	log        log.Logger
	db         dal.Dal
	validator  *validator.Validate
	pluginName string
}

// NewScopeConfigHelper creates a ScopeConfigHelper for scope config management
func NewScopeConfigHelper[Tr dal.Tabler](
	basicRes context.BasicRes,
	vld *validator.Validate,
	pluginName string,
) *ScopeConfigHelper[Tr] {
	if vld == nil {
		vld = validator.New()
	}
	return &ScopeConfigHelper[Tr]{
		log:        basicRes.GetLogger(),
		db:         basicRes.GetDal(),
		validator:  vld,
		pluginName: pluginName,
	}
}

func (t ScopeConfigHelper[ScopeConfig]) Create(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, e := strconv.ParseUint(input.Params["connectionId"], 10, 64)
	if e != nil || connectionId == 0 {
		return nil, errors.Default.Wrap(e, "the connection ID should be an non-zero integer")
	}
	var config ScopeConfig
	if err := DecodeMapStruct(input.Body, &config, false); err != nil {
		return nil, errors.Default.Wrap(err, "error in decoding scope config")
	}
	if t.validator != nil {
		if err := t.validator.Struct(config); err != nil {
			return nil, errors.Default.Wrap(err, "error validating scope config")
		}
	}
	valueConnectionId := reflect.ValueOf(&config).Elem().FieldByName("ConnectionId")
	if valueConnectionId.IsValid() {
		valueConnectionId.SetUint(connectionId)
	}

	if err := t.db.Create(&config); err != nil {
		if t.db.IsDuplicationError(err) {
			return nil, errors.BadInput.New("there was a scope config with the same name, please choose another name")
		}
		return nil, errors.BadInput.Wrap(err, "error on saving ScopeConfig")
	}
	return &plugin.ApiResourceOutput{Body: config, Status: http.StatusOK}, nil
}

func (t ScopeConfigHelper[ScopeConfig]) Update(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	scopeConfigId, e := strconv.ParseUint(input.Params["id"], 10, 64)
	if e != nil {
		return nil, errors.Default.Wrap(e, "the scope config ID should be an integer")
	}
	var old ScopeConfig
	err := t.db.First(&old, dal.Where("id = ?", scopeConfigId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving ScopeConfig")
	}
	err = DecodeMapStruct(input.Body, &old, true)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error decoding map into scopeConfig")
	}
	err = t.db.Update(&old, dal.Where("id = ?", scopeConfigId))
	if err != nil {
		if t.db.IsDuplicationError(err) {
			return nil, errors.BadInput.New("there was a scope config with the same name, please choose another name")
		}
		return nil, errors.BadInput.Wrap(err, "error on saving ScopeConfig")
	}
	return &plugin.ApiResourceOutput{Body: old, Status: http.StatusOK}, nil
}

func (t ScopeConfigHelper[ScopeConfig]) Get(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	scopeConfigId, err := strconv.ParseUint(input.Params["id"], 10, 64)
	if err != nil {
		return nil, errors.Default.Wrap(err, "the scope config ID should be an integer")
	}
	var config ScopeConfig
	err = t.db.First(&config, dal.Where("id = ?", scopeConfigId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on get ScopeConfig")
	}
	return &plugin.ApiResourceOutput{Body: config, Status: http.StatusOK}, nil
}

func (t ScopeConfigHelper[ScopeConfig]) List(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, e := strconv.ParseUint(input.Params["connectionId"], 10, 64)
	if e != nil || connectionId == 0 {
		return nil, errors.Default.Wrap(e, "the connection ID should be an non-zero integer")
	}
	var configs []ScopeConfig
	limit, offset := GetLimitOffset(input.Query, "pageSize", "page")
	err := t.db.All(&configs, dal.Where("connection_id = ?", connectionId), dal.Limit(limit), dal.Offset(offset))
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on get ScopeConfig list")
	}
	return &plugin.ApiResourceOutput{Body: configs, Status: http.StatusOK}, nil
}

func (t ScopeConfigHelper[ScopeConfig]) Delete(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	scopeConfigId, err := errors.Convert01(strconv.ParseUint(input.Params["id"], 10, 64))
	if err != nil {
		return nil, errors.Default.Wrap(err, "the scope config ID should be an integer")
	}
	connectionId, err := errors.Convert01(strconv.ParseUint(input.Params["connectionId"], 10, 64))
	if err != nil {
		return nil, errors.Default.Wrap(err, "the scope connection ID should be an integer")
	}
	var config ScopeConfig
	err = t.db.Delete(&config, dal.Where("id = ? AND connection_id = ?", scopeConfigId, connectionId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "error deleting ScopeConfig")
	}
	err = t.nullOutScopeReferences(scopeConfigId)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
}
func (t ScopeConfigHelper[ScopeConfig]) nullOutScopeReferences(scopeConfigId uint64) errors.Error {
	p, _ := plugin.GetPlugin(t.pluginName)
	pluginSrc, ok := p.(plugin.PluginSource)
	if !ok {
		return errors.Default.New("plugin doesn't implement PluginSource")
	}
	scopeModel := pluginSrc.Scope()
	if scopeModel == nil {
		return nil
	}
	return t.db.UpdateColumn(scopeModel.TableName(), "scope_config_id", nil, dal.Where("scope_config_id = ?", scopeConfigId))
}
