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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks"
	"github.com/mitchellh/mapstructure"
)

// CreateTransformationRule create transformation rule for Jira
// @Summary create transformation rule for Jira
// @Description create transformation rule for Jira
// @Tags plugins/jira
// @Accept application/json
// @Param connectionId path int true "connectionId"
// @Param transformationRule body tasks.JiraTransformationRule true "transformation rule"
// @Success 200  {object} tasks.JiraTransformationRule
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/connections/{connectionId}/transformation_rules [POST]
func CreateTransformationRule(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	rule, err := makeDbTransformationRuleFromInput(input)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "error in makeJiraTransformationRule")
	}
	newRule := map[string]interface{}{}
	err = errors.Convert(mapstructure.Decode(rule, &newRule))
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "error in makeJiraTransformationRule")
	}
	input.Body = newRule
	return trHelper.Create(input)
}

// UpdateTransformationRule update transformation rule for Jira
// @Summary update transformation rule for Jira
// @Description update transformation rule for Jira
// @Tags plugins/jira
// @Accept application/json
// @Param id path int true "id"
// @Param connectionId path int true "connectionId"
// @Param transformationRule body tasks.JiraTransformationRule true "transformation rule"
// @Success 200  {object} tasks.JiraTransformationRule
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/connections/{connectionId}/transformation_rules/{id} [PATCH]
func UpdateTransformationRule(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, e := strconv.ParseUint(input.Params["connectionId"], 10, 64)
	if e != nil || connectionId == 0 {
		return nil, errors.Default.Wrap(e, "the connection ID should be an non-zero integer")
	}
	transformationRuleId, e := strconv.ParseUint(input.Params["id"], 10, 64)
	if e != nil {
		return nil, errors.Default.Wrap(e, "the transformation rule ID should be an integer")
	}
	var req tasks.JiraTransformationRule
	err := api.Decode(input.Body, &req, vld)
	if err != nil {
		return nil, err
	}
	var oldDB models.JiraTransformationRule
	err = basicRes.GetDal().First(&oldDB, dal.Where("id = ?", transformationRuleId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on getting TransformationRule")
	}
	oldTr, err := tasks.MakeTransformationRules(oldDB)
	if err != nil {
		return nil, err
	}
	err = api.DecodeMapStruct(input.Body, oldTr, true)
	if err != nil {
		return nil, err
	}

	newDB, err := oldTr.ToDb()
	if err != nil {
		return nil, err
	}
	newDB.ID = transformationRuleId
	newDB.ConnectionId = connectionId
	newDB.CreatedAt = oldDB.CreatedAt
	err = basicRes.GetDal().Update(newDB)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: newDB, Status: http.StatusOK}, err
}

func makeDbTransformationRuleFromInput(input *plugin.ApiResourceInput) (*models.JiraTransformationRule, errors.Error) {
	connectionId, e := strconv.ParseUint(input.Params["connectionId"], 10, 64)
	if e != nil || connectionId == 0 {
		return nil, errors.Default.Wrap(e, "the connection ID should be an non-zero integer")
	}
	var req tasks.JiraTransformationRule
	err := api.Decode(input.Body, &req, vld)
	if err != nil {
		return nil, err
	}
	req.ConnectionId = connectionId
	return req.ToDb()
}

// GetTransformationRule return one transformation rule
// @Summary return one transformation rule
// @Description return one transformation rule
// @Tags plugins/jira
// @Param id path int true "id"
// @Param connectionId path int true "connectionId"
// @Success 200  {object} tasks.JiraTransformationRule
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/connections/{connectionId}/transformation_rules/{id} [GET]
func GetTransformationRule(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return trHelper.Get(input)
}

// GetTransformationRuleList return all transformation rules
// @Summary return all transformation rules
// @Description return all transformation rules
// @Tags plugins/jira
// @Param connectionId path int true "connectionId"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Success 200  {object} []tasks.JiraTransformationRule
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/connections/{connectionId}/transformation_rules [GET]
func GetTransformationRuleList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return trHelper.List(input)
}
