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
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
	"github.com/mitchellh/mapstructure"
)

type req struct {
	Data []*models.JenkinsJob `json:"data"`
}

// PutScope create or update jenkins job
// @Summary create or update jenkins job
// @Description Create or update jenkins job
// @Tags plugins/jenkins
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scope body req true "json"
// @Success 200  {object} []models.JenkinsJob
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jenkins/connections/{connectionId}/scopes [PUT]
func PutScope(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connectionId, _ := strconv.ParseUint(input.Params["connectionId"], 10, 64)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}
	var jobs req
	err := errors.Convert(mapstructure.Decode(input.Body, &jobs))
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "decoding Jenkins job error")
	}
	keeper := make(map[string]struct{})
	for _, job := range jobs.Data {
		if _, ok := keeper[job.FullName]; ok {
			return nil, errors.BadInput.New("duplicated item")
		} else {
			keeper[job.FullName] = struct{}{}
		}
		job.ConnectionId = connectionId

	}
	err = BasicRes.GetDal().CreateOrUpdate(jobs.Data)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving JenkinsJob")
	}
	return &core.ApiResourceOutput{Body: jobs.Data, Status: http.StatusOK}, nil
}

// UpdateScope patch to jenkins job
// @Summary patch to jenkins job
// @Description patch to jenkins job
// @Tags plugins/jenkins
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param fullName path string false "job's full name"
// @Param scope body models.JenkinsJob true "json"
// @Success 200  {object} models.JenkinsJob
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jenkins/connections/{connectionId}/scopes/{fullName} [PATCH]
func UpdateScope(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connectionId, fullName, err := extractParam(input.Params)
	if err != nil {
		return nil, err
	}
	var job models.JenkinsJob
	job.ConnectionId = connectionId
	job.FullName = fullName
	err = BasicRes.GetDal().First(&job, dal.Where("connection_id = ? AND full_name = ?", connectionId, fullName))
	if err != nil {
		return nil, errors.Default.Wrap(err, "getting JenkinsJob error")
	}
	err = helper.DecodeMapStruct(input.Body, &job)
	if err != nil {
		return nil, errors.Default.Wrap(err, "patch jenkins job error")
	}
	err = BasicRes.GetDal().Update(&job)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving JenkinsJob")
	}
	return &core.ApiResourceOutput{Body: job, Status: http.StatusOK}, nil
}

// GetScopeList get Jenkins jobs
// @Summary get Jenkins jobs
// @Description get Jenkins jobs
// @Tags plugins/jenkins
// @Param connectionId path int false "connection ID"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Success 200  {object} []models.JenkinsJob
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jenkins/connections/{connectionId}/scopes/ [GET]
func GetScopeList(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	var jobs []models.JenkinsJob
	connectionId, _ := strconv.ParseUint(input.Params["connectionId"], 10, 64)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid path params")
	}
	limit, offset := helper.GetLimitOffset(input.Query, "pageSize", "page")
	err := BasicRes.GetDal().All(&jobs, dal.Where("connection_id = ?", connectionId), dal.Limit(limit), dal.Offset(offset))
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jobs, Status: http.StatusOK}, nil
}

// GetScope get one Jenkins job
// @Summary get one Jenkins job
// @Description get one Jenkins job
// @Tags plugins/jenkins
// @Param connectionId path int false "connection ID"
// @Param fullName path string false "job's full name"
// @Success 200  {object} models.JenkinsJob
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jenkins/connections/{connectionId}/scopes/{fullName} [GET]
func GetScope(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	var job models.JenkinsJob
	connectionId, fullName, err := extractParam(input.Params)
	if err != nil {
		return nil, err
	}
	err = BasicRes.GetDal().First(&job, dal.Where("connection_id = ? AND full_name = ?", connectionId, fullName))
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: job, Status: http.StatusOK}, nil
}

func extractParam(params map[string]string) (uint64, string, errors.Error) {
	connectionId, _ := strconv.ParseUint(params["connectionId"], 10, 64)
	if connectionId == 0 {
		return 0, "", errors.BadInput.New("invalid connectionId")
	}
	if params["fullName"] == "" {
		return 0, "", errors.BadInput.New("invalid fullName")
	}
	return connectionId, params["fullName"], nil
}
