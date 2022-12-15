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
	"strconv"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/mitchellh/mapstructure"
)

// CreateTransformationRule create transformation rule for Gitlab
// @Summary create transformation rule for Gitlab
// @Description create transformation rule for Gitlab
// @Tags plugins/gitlab
// @Accept application/json
// @Param transformationRule body models.GitlabTransformationRule true "transformation rule"
// @Success 200  {object} models.GitlabTransformationRule
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/gitlab/transformation_rules [POST]
func CreateTransformationRule(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	var rule models.GitlabTransformationRule
	err := mapstructure.Decode(input.Body, &rule)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error in decoding transformation rule")
	}
	err = basicRes.GetDal().Create(&rule)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving TransformationRule")
	}
	return &core.ApiResourceOutput{Body: rule, Status: http.StatusOK}, nil
}

// UpdateTransformationRule update transformation rule for Gitlab
// @Summary update transformation rule for Gitlab
// @Description update transformation rule for Gitlab
// @Tags plugins/gitlab
// @Accept application/json
// @Param id path int true "id"
// @Param transformationRule body models.GitlabTransformationRule true "transformation rule"
// @Success 200  {object} models.GitlabTransformationRule
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/gitlab/transformation_rules/{id} [PATCH]
func UpdateTransformationRule(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	transformationRuleId, err := strconv.ParseUint(input.Params["id"], 10, 64)
	if err != nil {
		return nil, errors.Default.Wrap(err, "the transformation rule ID should be an integer")
	}
	var old models.GitlabTransformationRule
	err = basicRes.GetDal().First(&old, dal.Where("id = ?", transformationRuleId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving TransformationRule")
	}
	err = helper.DecodeMapStruct(input.Body, &old)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error decoding map into transformationRule")
	}
	old.ID = transformationRuleId
	err = basicRes.GetDal().Update(&old, dal.Where("id = ?", transformationRuleId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving TransformationRule")
	}
	return &core.ApiResourceOutput{Body: old, Status: http.StatusOK}, nil
}

// GetTransformationRule return one transformation rule
// @Summary return one transformation rule
// @Description return one transformation rule
// @Tags plugins/gitlab
// @Param id path int true "id"
// @Success 200  {object} models.GitlabTransformationRule
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/gitlab/transformation_rules/{id} [GET]
func GetTransformationRule(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	transformationRuleId, err := strconv.ParseUint(input.Params["id"], 10, 64)
	if err != nil {
		return nil, errors.Default.Wrap(err, "the transformation rule ID should be an integer")
	}
	var rule models.GitlabTransformationRule
	err = basicRes.GetDal().First(&rule, dal.Where("id = ?", transformationRuleId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on get TransformationRule")
	}
	return &core.ApiResourceOutput{Body: rule, Status: http.StatusOK}, nil
}

// GetTransformationRuleList return all transformation rules
// @Summary return all transformation rules
// @Description return all transformation rules
// @Tags plugins/gitlab
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Success 200  {object} []models.GitlabTransformationRule
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/gitlab/transformation_rules [GET]
func GetTransformationRuleList(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	var rules []models.GitlabTransformationRule
	limit, offset := helper.GetLimitOffset(input.Query, "pageSize", "page")
	err := basicRes.GetDal().All(&rules, dal.Limit(limit), dal.Offset(offset))
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on get TransformationRule list")
	}
	return &core.ApiResourceOutput{Body: rules, Status: http.StatusOK}, nil
}
