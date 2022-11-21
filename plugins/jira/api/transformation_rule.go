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
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks"
	"github.com/mitchellh/mapstructure"
)

// CreateTransformationRule create transformation rule for Jira
// @Summary create transformation rule for Jira
// @Description create transformation rule for Jira
// @Tags plugins/jira
// @Accept application/json
// @Param transformationRule body tasks.TransformationRules true "transformation rule"
// @Success 200  {object} core.ApiResourceOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/transformation_rules [POST]
func CreateTransformationRule(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	rule, err := makeJiraTransformationRule(input)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error in makeJiraTransformationRule")
	}
	err = basicRes.GetDal().Create(&rule)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving TransformationRule")
	}
	return &core.ApiResourceOutput{Status: http.StatusOK}, nil
}

// UpdateTransformationRule update transformation rule for Jira
// @Summary update transformation rule for Jira
// @Description update transformation rule for Jira
// @Tags plugins/jira
// @Accept application/json
// @Param id path int true "id"
// @Param transformationRule body tasks.TransformationRules true "transformation rule"
// @Success 200  {object} core.ApiResourceOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/transformation_rules/{id} [PUT]
func UpdateTransformationRule(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	transformationRuleId, err := strconv.ParseUint(input.Params["id"], 10, 64)
	if err != nil {
		return nil, errors.Default.Wrap(err, "the transformation rule ID should be an integer")
	}
	rule, err := makeJiraTransformationRule(input)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error in makeJiraTransformationRule")
	}
	err = basicRes.GetDal().Update(&rule, dal.Where("id = ?", transformationRuleId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving TransformationRule")
	}
	return &core.ApiResourceOutput{Status: http.StatusOK}, nil
}

func makeJiraTransformationRule(input *core.ApiResourceInput) (*models.JiraTransformationRule, errors.Error) {
	var req tasks.TransformationRules
	err := mapstructure.Decode(input.Body, &req)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error decoding map into putBoardRequest")
	}
	blob, err := json.Marshal(req.TypeMappings)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error marshaling TypeMappings")
	}
	return &models.JiraTransformationRule{
		EpicKeyField:               req.EpicKeyField,
		StoryPointField:            req.StoryPointField,
		RemotelinkCommitShaPattern: req.RemotelinkCommitShaPattern,
		TypeMappings:               blob,
	}, nil
}

// GetTransformationRule return one transformation rule
// @Summary return one transformation rule
// @Description return one transformation rule
// @Tags plugins/jira
// @Param id path int true "id"
// @Success 200  {object} tasks.TransformationRules
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/transformation_rules/{id} [GET]
func GetTransformationRule(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	transformationRuleId, err := strconv.ParseUint(input.Params["id"], 10, 64)
	if err != nil {
		return nil, errors.Default.Wrap(err, "the transformation rule ID should be an integer")
	}
	var rule models.JiraTransformationRule
	err = basicRes.GetDal().First(&rule, dal.Where("id = ?", transformationRuleId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on get TransformationRule")
	}
	return &core.ApiResourceOutput{Body: rule, Status: http.StatusOK}, nil
}

// GetTransformationRuleList return all transformation rules
// @Summary return all transformation rules
// @Description return all transformation rules
// @Tags plugins/jira
// @Success 200  {object} []tasks.TransformationRules
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/transformation_rules [GET]
func GetTransformationRuleList(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	var rules []models.JiraTransformationRule
	err := basicRes.GetDal().All(&rules)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on get TransformationRule list")
	}
	return &core.ApiResourceOutput{Body: rules, Status: http.StatusOK}, nil
}
