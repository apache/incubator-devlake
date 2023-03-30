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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/trello/models"
	"net/http"
	"strconv"
)

// CreateTransformationRule create transformation rule for Trello
// @Summary create transformation rule for Trello
// @Description create transformation rule for Trello
// @Tags plugins/trello
// @Accept application/json
// @Param transformationRule body models.TrelloTransformationRule true "transformation rule"
// @Success 200  {object} models.TrelloTransformationRule
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/trello/transformation_rules [POST]
func CreateTransformationRule(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var rule models.TrelloTransformationRule
	err := api.Decode(input.Body, &rule, vld)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "error in decoding transformation rule")
	}
	err = basicRes.GetDal().Create(&rule)
	if err != nil {
		if basicRes.GetDal().IsDuplicationError(err) {
			return nil, errors.BadInput.New("there was a transformation rule with the same name, please choose another name")
		}
		return nil, errors.BadInput.Wrap(err, "error on saving TransformationRule")
	}
	return &plugin.ApiResourceOutput{Body: rule, Status: http.StatusOK}, nil
}

// UpdateTransformationRule update transformation rule for Trello
// @Summary update transformation rule for Trello
// @Description update transformation rule for Trello
// @Tags plugins/trello
// @Accept application/json
// @Param id path int true "id"
// @Param transformationRule body models.TrelloTransformationRule true "transformation rule"
// @Success 200  {object} models.TrelloTransformationRule
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/trello/transformation_rules/{id} [PATCH]
func UpdateTransformationRule(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	transformationRuleId, e := strconv.ParseUint(input.Params["id"], 10, 64)
	if e != nil {
		return nil, errors.Default.Wrap(e, "the transformation rule ID should be an integer")
	}
	var old models.TrelloConnection
	err := basicRes.GetDal().First(&old, dal.Where("id = ?", transformationRuleId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving TransformationRule")
	}
	err = api.DecodeMapStruct(input.Body, &old, true)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error decoding map into transformationRule")
	}
	old.ID = transformationRuleId
	err = basicRes.GetDal().Update(&old, dal.Where("id = ?", transformationRuleId))
	if err != nil {
		if basicRes.GetDal().IsDuplicationError(err) {
			return nil, errors.BadInput.New("there was a transformation rule with the same name, please choose another name")
		}
		return nil, errors.BadInput.Wrap(err, "error on saving TransformationRule")
	}
	return &plugin.ApiResourceOutput{Body: old, Status: http.StatusOK}, nil
}

// GetTransformationRule return one transformation rule
// @Summary return one transformation rule
// @Description return one transformation rule
// @Tags plugins/trello
// @Param id path int true "id"
// @Success 200  {object} models.TrelloTransformationRule
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/trello/transformation_rules/{id} [GET]
func GetTransformationRule(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	transformationRuleId, err := strconv.ParseUint(input.Params["id"], 10, 64)
	if err != nil {
		return nil, errors.Default.Wrap(err, "the transformation rule ID should be an integer")
	}
	var rule models.TrelloTransformationRule
	err = basicRes.GetDal().First(&rule, dal.Where("id = ?", transformationRuleId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on get TransformationRule")
	}
	return &plugin.ApiResourceOutput{Body: rule, Status: http.StatusOK}, nil
}

// GetTransformationRuleList return all transformation rules
// @Summary return all transformation rules
// @Description return all transformation rules
// @Tags plugins/trello
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Success 200  {object} []models.TrelloTransformationRule
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/trello/transformation_rules [GET]
func GetTransformationRuleList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var rules []models.TrelloTransformationRule
	limit, offset := api.GetLimitOffset(input.Query, "pageSize", "page")
	err := basicRes.GetDal().All(&rules, dal.Limit(limit), dal.Offset(offset))
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on get TransformationRule list")
	}
	return &plugin.ApiResourceOutput{Body: rules, Status: http.StatusOK}, nil
}
