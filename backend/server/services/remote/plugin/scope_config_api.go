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
	"github.com/apache/incubator-devlake/server/services/remote/models"
	"net/http"
	"strconv"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

func (pa *pluginAPI) PostScopeConfigs(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, _ := strconv.ParseUint(input.Params["connectionId"], 10, 64)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}
	scopeConfig := pa.scopeConfigType.New()
	input.Body[`connectionId`] = connectionId
	err := api.DecodeMapStruct(input.Body, scopeConfig, false)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "error in decoding scope config")
	}
	db := basicRes.GetDal()
	err = api.CallDB(db.Create, scopeConfig)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: scopeConfig.Unwrap(), Status: http.StatusOK}, nil
}

func (pa *pluginAPI) PatchScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, trId, err := extractConfigParam(input.Params)
	if err != nil {
		return nil, err
	}

	scopeConfig := pa.scopeConfigType.New()
	db := basicRes.GetDal()
	err = api.CallDB(db.First, scopeConfig, dal.Where("connection_id = ? AND id = ?", connectionId, trId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "no scope config with given id")
	}

	input.Body[`connectionId`] = connectionId
	input.Body[`id`] = trId
	err = api.DecodeMapStruct(input.Body, scopeConfig, false)
	if err != nil {
		return nil, errors.Default.Wrap(err, "decoding error")
	}

	err = api.CallDB(db.Update, scopeConfig)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: scopeConfig.Unwrap(), Status: http.StatusOK}, nil
}

func (pa *pluginAPI) GetScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	scopeConfig := pa.scopeConfigType.New()
	db := basicRes.GetDal()
	connectionId, configId, err := extractConfigParam(input.Params)
	if err != nil {
		return nil, err
	}
	err = api.CallDB(db.First, scopeConfig, dal.Where("connection_id = ? AND id = ?", connectionId, configId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "no scope config with given id")
	}

	return &plugin.ApiResourceOutput{Body: scopeConfig.Unwrap()}, nil
}

func (pa *pluginAPI) ListScopeConfigs(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	scopeConfigs := pa.scopeConfigType.NewSlice()
	limit, offset := api.GetLimitOffset(input.Query, "pageSize", "page")
	if limit > 100 {
		return nil, errors.BadInput.New("pageSize cannot exceed 100")
	}

	db := basicRes.GetDal()
	err := api.CallDB(db.All, scopeConfigs, dal.Limit(limit), dal.Offset(offset))
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: scopeConfigs.Unwrap()}, nil
}

func (pa *pluginAPI) DeleteScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, configId, err := extractConfigParam(input.Params)
	if err != nil {
		return nil, err
	}
	scopeConfig := pa.scopeConfigType.New()
	db := basicRes.GetDal()
	err = api.CallDB(db.Delete, scopeConfig, dal.Where("connection_id = ? AND id = ?", connectionId, configId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "no scope config with given id")
	}
	err = pa.nullOutScopeReferences(configId)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: scopeConfig.Unwrap()}, nil
}

func (pa *pluginAPI) nullOutScopeReferences(scopeConfigId uint64) errors.Error {
	scopeModel := models.NewDynamicScopeModel(pa.scopeType)
	return basicRes.GetDal().Exec(fmt.Sprintf("UPDATE %s SET scope_config_id = NULL WHERE scope_config_id = %d", scopeModel.TableName(), scopeConfigId))
}

func extractConfigParam(params map[string]string) (connectionId uint64, configId uint64, err errors.Error) {
	connectionId, _ = strconv.ParseUint(params["connectionId"], 10, 64)
	configId, _ = strconv.ParseUint(params["id"], 10, 64)
	if connectionId == 0 {
		return 0, 0, errors.BadInput.New("invalid connectionId")
	}
	if configId == 0 {
		return 0, 0, errors.BadInput.New("invalid configId")
	}

	return connectionId, configId, nil
}
