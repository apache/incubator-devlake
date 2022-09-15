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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/webhook/models"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

type WebhookPipelineRequest struct {
	Id           string     `validate:"required"`
	Result       string     `validate:"oneof=SUCCESS FAILURE ABORT"`
	Status       string     `validate:"oneof=IN_PROGRESS DONE"`
	Type         string     `validate:"oneof=CI CD CI/CD"`
	CreatedDate  time.Time  `mapstructure:"created_date" validate:"required"`
	FinishedDate *time.Time `mapstructure:"finished_date"`

	Repo      string `validate:"required"`
	Branch    string
	CommitSha string `mapstructure:"commit_sha"`
}

// PostCicdPipeline
// @Summary create pipeline by webhook
// @Description Create pipeline by webhook, example: {"id":"A123123","result":"one of SUCCESS/FAILURE/ABORT","status":"one of IN_PROGRESS/DONE","type":"CI/CD","created_date":"2020-01-01T12:00:00+00:00","finished_date":"2020-01-01T12:59:59+00:00","repo":"devlake","branch":"main","commit_sha":"015e3d3b480e417aede5a1293bd61de9b0fd051d"}
// @Tags plugins/webhook
// @Param body body WebhookPipelineRequest true "json body"
// @Success 200
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/webhook/:connectionId/cicd_pipelines [POST]
func PostCicdPipeline(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	// get request
	request := &WebhookPipelineRequest{}
	err = helper.DecodeMapStruct(input.Body, request)
	if err != nil {
		return &core.ApiResourceOutput{Body: err.Error(), Status: http.StatusBadRequest}, nil
	}
	// validate
	vld = validator.New()
	err = errors.Convert(vld.Struct(request))
	if err != nil {
		return &core.ApiResourceOutput{Body: err.Error(), Status: http.StatusBadRequest}, nil
	}

	db := basicRes.GetDal()
	domainPipeline := &devops.CICDPipeline{
		DomainEntity: domainlayer.DomainEntity{
			Id: fmt.Sprintf("%s:%d:%s:%s:%s:%s", "webhook", connection.ID, request.Repo, request.Branch, request.CommitSha, request.Id),
		},
		Name:         request.Id,
		Result:       request.Result,
		Status:       request.Status,
		Type:         request.Type,
		CreatedDate:  request.CreatedDate,
		FinishedDate: request.FinishedDate,
	}
	if domainPipeline.FinishedDate != nil {
		domainPipeline.DurationSec = uint64(domainPipeline.FinishedDate.Sub(domainPipeline.CreatedDate).Seconds())
	}
	domainPipelineRepo := &devops.CiCDPipelineRepo{
		DomainEntity: domainlayer.DomainEntity{
			Id: fmt.Sprintf("%s:%d:%s:%s:%s", "webhook", connection.ID, request.Repo, request.Branch, request.CommitSha),
		},
		CommitSha: request.CommitSha,
		Branch:    request.Branch,
		Repo:      request.Repo,
	}
	// save
	err = db.CreateOrUpdate(domainPipeline)
	if err != nil {
		return nil, err
	}
	err = db.CreateOrUpdate(domainPipelineRepo)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
}
