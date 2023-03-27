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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

// CreateTransformationRule create transformation rule for Bamboo
// @Summary create transformation rule for Bamboo
// @Description create transformation rule for Bamboo
// @Tags plugins/bamboo
// @Accept application/json
// @Param connectionId path int true "connectionId"
// @Param transformationRule body models.BambooTransformationRule true "transformation rule"
// @Success 200  {object} models.BambooTransformationRule
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/transformation_rules [POST]
func CreateTransformationRule(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return trHelper.Create(input)
}

// UpdateTransformationRule update transformation rule for Bamboo
// @Summary update transformation rule for Bamboo
// @Description update transformation rule for Bamboo
// @Tags plugins/bamboo
// @Accept application/json
// @Param id path int true "id"
// @Param transformationRule body models.BambooTransformationRule true "transformation rule"
// @Param connectionId path int true "connectionId"
// @Success 200  {object} models.BambooTransformationRule
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/transformation_rules/{id} [PATCH]
func UpdateTransformationRule(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return trHelper.Update(input)
}

// GetTransformationRule return one transformation rule
// @Summary return one transformation rule
// @Description return one transformation rule
// @Tags plugins/bamboo
// @Param id path int true "id"
// @Param connectionId path int true "connectionId"
// @Success 200  {object} models.BambooTransformationRule
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/transformation_rules/{id} [GET]
func GetTransformationRule(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return trHelper.Get(input)
}

// GetTransformationRuleList return all transformation rules
// @Summary return all transformation rules
// @Description return all transformation rules
// @Tags plugins/bamboo
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Param connectionId path int true "connectionId"
// @Success 200  {object} []models.BambooTransformationRule
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/transformation_rules [GET]
func GetTransformationRuleList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return trHelper.List(input)
}
