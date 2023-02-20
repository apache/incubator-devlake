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
	"net/http"
	"strconv"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type TransformationRuleAPI struct {
	txRuleType *models.DynamicTabler
}

func (t *TransformationRuleAPI) PostTransformationRules(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	txRule := t.txRuleType.New()
	err := api.Decode(input.Body, txRule, vld)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "error in decoding transformation rule")
	}
	db := basicRes.GetDal()
	err = api.CallDB(db.Create, txRule)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: txRule.Unwrap(), Status: http.StatusOK}, nil
}

func (t *TransformationRuleAPI) PatchTransformationRule(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	id, err := strconv.ParseUint(input.Params["id"], 10, 64)
	if err != nil {
		return nil, errors.Default.Wrap(err, "id should be an integer")
	}

	txRule := t.txRuleType.New()
	db := basicRes.GetDal()
	err = api.CallDB(db.First, txRule, dal.Where("id = ?", id))
	if err != nil {
		return nil, errors.Default.Wrap(err, "no transformation rule with given id")
	}

	err = api.Decode(input.Body, txRule, vld)
	if err != nil {
		return nil, errors.Default.Wrap(err, "decoding error")
	}

	return &plugin.ApiResourceOutput{Body: txRule.Unwrap(), Status: http.StatusOK}, nil
}

func (t *TransformationRuleAPI) GetTransformationRule(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	txRule := t.txRuleType.New()
	db := basicRes.GetDal()
	err := api.CallDB(db.First, txRule, dal.Where("id = ?", input.Params))
	if err != nil {
		return nil, errors.Default.Wrap(err, "no transformation rule with given id")
	}

	return &plugin.ApiResourceOutput{Body: txRule.Unwrap()}, nil
}

func (t *TransformationRuleAPI) ListTransformationRules(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	txRules := t.txRuleType.NewSlice()
	limit, offset := api.GetLimitOffset(input.Query, "pageSize", "page")
	if limit > 100 {
		return nil, errors.BadInput.New("pageSize cannot exceed 100")
	}

	db := basicRes.GetDal()
	err := api.CallDB(db.All, txRules, dal.Limit(limit), dal.Offset(offset))
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: txRules.Unwrap()}, nil
}
