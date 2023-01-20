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
	"github.com/apache/incubator-devlake/plugins/pagerduty/models"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"strconv"
)

type apiService struct {
	models.Service
	TransformationRuleName string `json:"transformationRuleName,omitempty"`
}

type req struct {
	Data []*models.Service `json:"data"`
}

// PutScope create or update pagerduty repo
// @Summary create or update pagerduty repo
// @Description Create or update pagerduty repo
// @Tags plugins/pagerduty
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param scope body req true "json"
// @Success 200  {object} []models.Service
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/pagerduty/connections/{connectionId}/scopes [PUT]
func PutScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, _ := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}
	var services req
	err := errors.Convert(mapstructure.Decode(input.Body, &services))
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "decoding PagerDuty service error")
	}
	keeper := make(map[string]struct{})
	for _, service := range services.Data {
		if _, ok := keeper[service.Id]; ok {
			return nil, errors.BadInput.New("duplicated item")
		} else {
			keeper[service.Id] = struct{}{}
		}
		service.ConnectionId = connectionId
		err = verifyService(service)
		if err != nil {
			return nil, err
		}
	}
	err = basicRes.GetDal().CreateOrUpdate(services.Data)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving Service")
	}
	return &plugin.ApiResourceOutput{Body: services.Data, Status: http.StatusOK}, nil
}

// UpdateScope patch to pagerduty repo
// @Summary patch to pagerduty repo
// @Description patch to pagerduty repo
// @Tags plugins/pagerduty
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param repoId path int true "repo ID"
// @Param scope body models.Service true "json"
// @Success 200  {object} models.Service
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/pagerduty/connections/{connectionId}/scopes/{repoId} [PATCH]
func UpdateScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, serviceId := extractParam(input.Params)
	if connectionId == 0 || serviceId == 0 {
		return nil, errors.BadInput.New("invalid connectionId or serviceId")
	}
	var service models.Service
	err := basicRes.GetDal().First(&service, dal.Where("connection_id = ? AND id = ?", connectionId, serviceId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "getting Service error")
	}
	err = api.DecodeMapStruct(input.Body, &service)
	if err != nil {
		return nil, errors.Default.Wrap(err, "patch service error")
	}
	err = verifyService(&service)
	if err != nil {
		return nil, err
	}
	err = basicRes.GetDal().Update(service)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving PagerDutyRepo")
	}
	return &plugin.ApiResourceOutput{Body: service, Status: http.StatusOK}, nil
}

// GetScopeList get PagerDuty repos
// @Summary get PagerDuty repos
// @Description get PagerDuty repos
// @Tags plugins/pagerduty
// @Param connectionId path int true "connection ID"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Success 200  {object} []apiService
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/pagerduty/connections/{connectionId}/scopes/ [GET]
func GetScopeList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var services []models.Service
	connectionId, _ := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid path params")
	}
	limit, offset := api.GetLimitOffset(input.Query, "pageSize", "page")
	err := basicRes.GetDal().All(&services, dal.Where("connection_id = ?", connectionId), dal.Limit(limit), dal.Offset(offset))
	if err != nil {
		return nil, err
	}
	var ruleIds []uint64
	for _, service := range services {
		if service.TransformationRuleId > 0 {
			ruleIds = append(ruleIds, service.TransformationRuleId)
		}
	}
	var rules []models.PagerdutyTransformationRule
	if len(ruleIds) > 0 {
		err = basicRes.GetDal().All(&rules, dal.Where("id IN (?)", ruleIds))
		if err != nil {
			return nil, err
		}
	}
	names := make(map[uint64]string)
	for _, rule := range rules {
		names[rule.ID] = rule.Name
	}
	var apiRepos []apiService
	for _, service := range services {
		apiRepos = append(apiRepos, apiService{service, names[service.TransformationRuleId]})
	}
	return &plugin.ApiResourceOutput{Body: apiRepos, Status: http.StatusOK}, nil
}

// GetScope get one PagerDuty repo
// @Summary get one PagerDuty repo
// @Description get one PagerDuty repo
// @Tags plugins/pagerduty
// @Param connectionId path int true "connection ID"
// @Param repoId path int true "repo ID"
// @Success 200  {object} apiService
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/pagerduty/connections/{connectionId}/scopes/{repoId} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var service models.Service
	connectionId, serviceId := extractParam(input.Params)
	if connectionId == 0 || serviceId == 0 {
		return nil, errors.BadInput.New("invalid path params")
	}
	db := basicRes.GetDal()
	err := db.First(&service, dal.Where("connection_id = ? AND id = ?", connectionId, serviceId))
	if db.IsErrorNotFound(err) {
		return nil, errors.NotFound.New("record not found")
	}
	if err != nil {
		return nil, err
	}
	var rule models.PagerdutyTransformationRule
	if service.TransformationRuleId > 0 {
		err = basicRes.GetDal().First(&rule, dal.Where("id = ?", service.TransformationRuleId))
		if err != nil {
			return nil, err
		}
	}
	return &plugin.ApiResourceOutput{Body: apiService{service, rule.Name}, Status: http.StatusOK}, nil
}

func extractParam(params map[string]string) (uint64, uint64) {
	connectionId, _ := strconv.ParseUint(params["connectionId"], 10, 64)
	serviceId, _ := strconv.ParseUint(params["serviceId"], 10, 64)
	return connectionId, serviceId
}

func verifyService(service *models.Service) errors.Error {
	if service.ConnectionId == 0 {
		return errors.BadInput.New("invalid connectionId")
	}
	if service.Id == "" {
		return errors.BadInput.New("invalid service ID")
	}
	return nil
}
