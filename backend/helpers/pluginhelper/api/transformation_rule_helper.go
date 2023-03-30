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

// TransformationRuleHelper is used to write the CURD of transformation rule
type TransformationRuleHelper[Tr dal.Tabler] struct {
	log       log.Logger
	db        dal.Dal
	validator *validator.Validate
}

// NewTransformationRuleHelper creates a TransformationRuleHelper for transformation rule management
func NewTransformationRuleHelper[Tr dal.Tabler](
	basicRes context.BasicRes,
	vld *validator.Validate,
) *TransformationRuleHelper[Tr] {
	if vld == nil {
		vld = validator.New()
	}
	return &TransformationRuleHelper[Tr]{
		log:       basicRes.GetLogger(),
		db:        basicRes.GetDal(),
		validator: vld,
	}
}

func (t TransformationRuleHelper[Tr]) Create(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, e := strconv.ParseUint(input.Params["connectionId"], 10, 64)
	if e != nil || connectionId == 0 {
		return nil, errors.Default.Wrap(e, "the connection ID should be an non-zero integer")
	}
	var rule Tr
	if err := DecodeMapStruct(input.Body, &rule, false); err != nil {
		return nil, errors.Default.Wrap(err, "error in decoding transformation rule")
	}
	if t.validator != nil {
		if err := t.validator.Struct(rule); err != nil {
			return nil, errors.Default.Wrap(err, "error validating transformation rule")
		}
	}
	valueConnectionId := reflect.ValueOf(&rule).Elem().FieldByName("ConnectionId")
	if valueConnectionId.IsValid() {
		valueConnectionId.SetUint(connectionId)
	}

	if err := t.db.Create(&rule); err != nil {
		if t.db.IsDuplicationError(err) {
			return nil, errors.BadInput.New("there was a transformation rule with the same name, please choose another name")
		}
		return nil, errors.BadInput.Wrap(err, "error on saving TransformationRule")
	}
	return &plugin.ApiResourceOutput{Body: rule, Status: http.StatusOK}, nil
}

func (t TransformationRuleHelper[Tr]) Update(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	transformationRuleId, e := strconv.ParseUint(input.Params["id"], 10, 64)
	if e != nil {
		return nil, errors.Default.Wrap(e, "the transformation rule ID should be an integer")
	}
	var old Tr
	err := t.db.First(&old, dal.Where("id = ?", transformationRuleId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving TransformationRule")
	}
	err = DecodeMapStruct(input.Body, &old, false)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error decoding map into transformationRule")
	}
	err = t.db.Update(&old, dal.Where("id = ?", transformationRuleId))
	if err != nil {
		if t.db.IsDuplicationError(err) {
			return nil, errors.BadInput.New("there was a transformation rule with the same name, please choose another name")
		}
		return nil, errors.BadInput.Wrap(err, "error on saving TransformationRule")
	}
	return &plugin.ApiResourceOutput{Body: old, Status: http.StatusOK}, nil
}

func (t TransformationRuleHelper[Tr]) Get(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	transformationRuleId, err := strconv.ParseUint(input.Params["id"], 10, 64)
	if err != nil {
		return nil, errors.Default.Wrap(err, "the transformation rule ID should be an integer")
	}
	var rule Tr
	err = t.db.First(&rule, dal.Where("id = ?", transformationRuleId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on get TransformationRule")
	}
	return &plugin.ApiResourceOutput{Body: rule, Status: http.StatusOK}, nil
}

func (t TransformationRuleHelper[Tr]) List(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, e := strconv.ParseUint(input.Params["connectionId"], 10, 64)
	if e != nil || connectionId == 0 {
		return nil, errors.Default.Wrap(e, "the connection ID should be an non-zero integer")
	}
	var rules []Tr
	limit, offset := GetLimitOffset(input.Query, "pageSize", "page")
	err := t.db.All(&rules, dal.Where("connection_id = ?", connectionId), dal.Limit(limit), dal.Offset(offset))
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on get TransformationRule list")
	}
	return &plugin.ApiResourceOutput{Body: rules, Status: http.StatusOK}, nil
}
