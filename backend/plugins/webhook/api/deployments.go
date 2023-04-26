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
	"crypto/md5"
	"fmt"
	"net/http"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/webhook/models"

	"github.com/go-playground/validator/v10"
)

type WebhookDeployTaskRequest struct {
	PipelineId string `mapstructure:"pipeline_id"`
	// RepoUrl should be unique string, fill url or other unique data
	RepoId    string `mapstructure:"repo_id"`
	RepoUrl   string `mapstructure:"repo_url" validate:"required"`
	CommitSha string `mapstructure:"commit_sha" validate:"required"`
	RefName   string `mapstructure:"ref_name"`
	Result    string `mapstructure:"result"`
	// start_time and end_time is more readable for users,
	// StartedDate and FinishedDate is same as columns in db.
	// So they all keep.
	CreatedDate  *time.Time `mapstructure:"create_time"`
	StartedDate  *time.Time `mapstructure:"start_time" validate:"required"`
	FinishedDate *time.Time `mapstructure:"end_time"`
	Environment  string     `validate:"omitempty,oneof=PRODUCTION STAGING TESTING DEVELOPMENT"`
}

// PostDeploymentCicdTask
// @Summary create deployment by webhook
// @Description Create deployment pipeline by webhook.<br/>
// @Description example1: {"repo_url":"devlake","commit_sha":"015e3d3b480e417aede5a1293bd61de9b0fd051d","start_time":"2020-01-01T12:00:00+00:00","end_time":"2020-01-01T12:59:59+00:00","environment":"PRODUCTION"}<br/>
// @Description So we suggest request before task after deployment pipeline finish.
// @Description Both cicd_pipeline and cicd_task will be created
// @Tags plugins/webhook
// @Param body body WebhookDeployTaskRequest true "json body"
// @Success 200
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 403  {string} errcode.Error "Forbidden"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/webhook/:connectionId/deployments [POST]
func PostDeploymentCicdTask(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	// get request
	request := &WebhookDeployTaskRequest{}
	err = api.DecodeMapStruct(input.Body, request, true)
	if err != nil {
		return &plugin.ApiResourceOutput{Body: err.Error(), Status: http.StatusBadRequest}, nil
	}
	// validate
	vld = validator.New()
	err = errors.Convert(vld.Struct(request))
	if err != nil {
		return nil, errors.BadInput.Wrap(vld.Struct(request), `input json error`)
	}
	db := basicRes.GetDal()
	urlHash16 := fmt.Sprintf("%x", md5.Sum([]byte(request.RepoUrl)))[:16]
	scopeId := fmt.Sprintf("%s:%d", "webhook", connection.ID)
	deploymentCommitId := fmt.Sprintf("%s:%d:%s:%s", "webhook", connection.ID, urlHash16, request.CommitSha)
	pipelineId := request.PipelineId
	if pipelineId == "" {
		pipelineId = deploymentCommitId
	}
	if request.CreatedDate == nil {
		request.CreatedDate = request.StartedDate
	}
	if request.Environment == "" {
		request.Environment = devops.PRODUCTION
	}
	if request.FinishedDate == nil {
		now := time.Now()
		request.FinishedDate = &now
	}
	if request.Result == "" {
		request.Result = devops.SUCCESS
	}
	duration := uint64(request.FinishedDate.Sub(*request.StartedDate).Seconds())

	// create a deployment_commit record
	deploymentCommit := &devops.CicdDeploymentCommit{
		DomainEntity: domainlayer.DomainEntity{
			Id: deploymentCommitId,
		},
		CicdDeploymentId: pipelineId,
		CicdScopeId:      scopeId,
		Name:             fmt.Sprintf(`deployment for %s`, request.CommitSha),
		Result:           request.Result,
		Status:           devops.DONE,
		Environment:      request.Environment,
		CreatedDate:      *request.CreatedDate,
		StartedDate:      request.StartedDate,
		FinishedDate:     request.FinishedDate,
		DurationSec:      &duration,
		CommitSha:        request.CommitSha,
		RefName:          request.RefName,
		RepoId:           request.RepoId,
		RepoUrl:          request.RepoUrl,
	}
	err = db.CreateOrUpdate(deploymentCommit)
	if err != nil {
		return nil, err
	}

	// TODO: create a deployment record when the table is ready

	return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
}
