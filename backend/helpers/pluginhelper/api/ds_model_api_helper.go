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
	"net/http"
	"reflect"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/helpers/utils"
	"github.com/go-playground/validator/v10"
)

var vld = validator.New()

type AnyModelApiHelper struct {
	*srvhelper.AnyModelSrvHelper
	basicRes       context.BasicRes
	log            log.Logger
	modelName      string
	pkPathVarNames []string
	sterilizers    []func(m any) any
}

func NewAnyModelApiHelper(
	basicRes context.BasicRes,
	srvHelper *srvhelper.AnyModelSrvHelper,
	pkPathVarNames []string, // path variable names of primary key
	sterilizer func(m any) any,
) *AnyModelApiHelper {
	m := srvHelper.New()
	modelName := fmt.Sprintf("%T", m)
	modelApiHelper := &AnyModelApiHelper{
		AnyModelSrvHelper: srvHelper,
		basicRes:          basicRes,
		log:               basicRes.GetLogger().Nested(fmt.Sprintf("%s_dal", modelName)),
		modelName:         modelName,
		pkPathVarNames:    pkPathVarNames,
	}
	if sterilizer != nil {
		modelApiHelper.sterilizers = []func(m any) any{sterilizer}
	}
	return modelApiHelper
}

func (modelApi *AnyModelApiHelper) Post(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	model := modelApi.New()
	err := utils.DecodeMapStruct(input.Body, model, false)
	if err != nil {
		return nil, err
	}
	err = modelApi.CreateAny(model)
	if err != nil {
		return nil, err
	}
	model = modelApi.Sanitize(model)
	return &plugin.ApiResourceOutput{
		Status: http.StatusCreated,
		Body:   model,
	}, nil
}

func (modelApi *AnyModelApiHelper) ExtractPkValues(input *plugin.ApiResourceInput) ([]interface{}, errors.Error) {
	pkv := make([]interface{}, len(modelApi.pkPathVarNames))
	for i, pkn := range modelApi.pkPathVarNames {
		var ok bool
		pkv[i], ok = input.Params[pkn]
		if !ok {
			return nil, errors.BadInput.New(fmt.Sprintf("missing path variable %s", pkn))
		}
	}
	return pkv, nil
}

func (modelApi *AnyModelApiHelper) FindByPkAny(input *plugin.ApiResourceInput) (any, errors.Error) {
	pkv, err := modelApi.ExtractPkValues(input)
	if err != nil {
		return nil, err
	}
	return modelApi.AnyModelSrvHelper.FindByPkAny(pkv...)
}

func (modelApi *AnyModelApiHelper) GetDetail(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	model, err := modelApi.FindByPkAny(input)
	if err != nil {
		return nil, err
	}
	model = modelApi.Sanitize(model)
	return &plugin.ApiResourceOutput{
		Body: model,
	}, nil
}

func (modelApi *AnyModelApiHelper) Sanitize(model any) any {
	if modelApi.sterilizers != nil {
		for _, sterilizer := range modelApi.sterilizers {
			sanitizedModel := sterilizer(model)
			model = sanitizedModel
		}
	}
	return model
}

func (modelApi *AnyModelApiHelper) BatchSanitize(models any) any {
	array := reflect.ValueOf(models)
	for i := 0; i < array.Len(); i++ {
		elem := array.Index(i)
		model := elem.Interface()
		sanitized := modelApi.Sanitize(model)
		elem.Set(reflect.ValueOf(sanitized))
	}
	return models
}

type CustomMerge interface {
	MergeFromRequest(target any, body map[string]interface{}) error
}

// PatchModelAny will get an "M" from database and try to merge update from request body
// zeroFields decides whether "M" will be zeroed if "M" doesn't implement CustomMerge.
func (modelApi *AnyModelApiHelper) PatchModel(input *plugin.ApiResourceInput, zeroFields bool) (any, error) {
	model, err := modelApi.FindByPkAny(input)
	if err != nil {
		return nil, err
	}
	if v, ok := (model).(CustomMerge); ok {
		if err := v.MergeFromRequest(model, input.Body); err != nil {
			return nil, err
		}
	} else {
		err = utils.DecodeMapStruct(input.Body, model, zeroFields)
		if err != nil {
			return nil, errors.BadInput.Wrap(err, fmt.Sprintf("faled to patch %s", modelApi.modelName))
		}
	}
	return model, nil
}

func (modelApi *AnyModelApiHelper) Patch(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	model, err := modelApi.PatchModel(input, true)
	if err != nil {
		return nil, errors.Convert(err)
	}
	if err := modelApi.UpdateAny(model); err != nil {
		return nil, err
	}
	model = modelApi.Sanitize(model)
	return &plugin.ApiResourceOutput{
		Body: model,
	}, nil
}

func (modelApi *AnyModelApiHelper) Delete(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	model, err := modelApi.FindByPkAny(input)
	if err != nil {
		return nil, err
	}
	err = modelApi.DeleteModelAny(model)
	if err != nil {
		return nil, err
	}
	model = modelApi.Sanitize(model)
	return &plugin.ApiResourceOutput{
		Body: model,
	}, nil
}

func (modelApi *AnyModelApiHelper) GetAll(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	all, err := modelApi.QueryAllAny()
	all = modelApi.BatchSanitize(all)
	return &plugin.ApiResourceOutput{
		Body: all,
	}, err
}

func (modelApi *AnyModelApiHelper) PutMultipleCb(input *plugin.ApiResourceInput, beforeSave func(any) errors.Error) (*plugin.ApiResourceOutput, errors.Error) {
	var req struct {
		Data []any `json:"data"`
	}
	err := utils.DecodeMapStruct(input.Body, &req, false)
	if err != nil {
		return nil, err
	}
	for i, itemDict := range req.Data {
		item := modelApi.New()
		err := utils.DecodeMapStruct(itemDict, item, false)
		if err != nil {
			return nil, err
		}
		if beforeSave != nil {
			err := beforeSave(item)
			if err != nil {
				return nil, err
			}
		}
		err = modelApi.CreateOrUpdateAny(item)
		if err != nil {
			return nil, errors.BadInput.Wrap(err, fmt.Sprintf("failed to save item %d", i))
		}
	}
	// TODO
	// req.Data = modelApi.BatchSanitize(req.Data)
	return &plugin.ApiResourceOutput{
		Body: req.Data,
	}, nil
}

func parsePagination[P any](input *plugin.ApiResourceInput) (*P, errors.Error) {
	if !input.Query.Has("page") {
		input.Query.Set("page", "1")
	}
	if !input.Query.Has("pageSize") {
		input.Query.Set("pageSize", "100")
	}
	pagination := new(P)
	err := utils.DecodeMapStruct(input.Query, pagination, false)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "faild to decode pagination from query string")
	}
	err = utils.DecodeMapStruct(input.Params, pagination, false)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "faild to decode pagination from path variables")
	}
	if e := vld.Struct(pagination); e != nil {
		return nil, errors.BadInput.Wrap(e, "invalid pagination parameters")
	}
	return pagination, nil
}

type ModelApiHelper[M dal.Tabler] struct {
	apiHelper *AnyModelApiHelper
}

func NewModelApiHelper[M dal.Tabler](
	anyModelApiHelper *AnyModelApiHelper,
) *ModelApiHelper[M] {
	return &ModelApiHelper[M]{
		apiHelper: anyModelApiHelper,
	}
}

func (modelApi *ModelApiHelper[M]) FindByPk(input *plugin.ApiResourceInput) (*M, errors.Error) {
	model, err := modelApi.apiHelper.FindByPkAny(input)
	return model.(*M), err
}
