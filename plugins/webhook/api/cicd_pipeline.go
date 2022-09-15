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
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/webhook/models"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

type WebhookTaskRequest struct {
	// PipelineName can be filled by any string unique in one pipeline
	PipelineName string `mapstructure:"pipeline_name" validate:"required"`

	Name         string     `validate:"required"` // Name should be unique in one pipeline
	Result       string     `validate:"oneof=SUCCESS FAILURE ABORT IN_PROGRESS"`
	Status       string     `validate:"oneof=IN_PROGRESS DONE"`
	Type         string     `validate:"oneof=CI CD CI/CD"`
	StartedDate  time.Time  `mapstructure:"created_date" validate:"required"`
	FinishedDate *time.Time `mapstructure:"finished_date"`

	RepoId    string `mapstructure:"repo_id" validate:"required"` // RepoId should be unique string
	Branch    string
	CommitSha string `mapstructure:"commit_sha"`
}

// PostCicdTask
// @Summary create pipeline by webhook
// @Description Create pipeline by webhook.<br/>
// @Description example1: {"pipeline_name":"A123","name":"unit-test","result":"IN_PROGRESS","status":"IN_PROGRESS","type":"CI","created_date":"2020-01-01T12:00:00+00:00","finished_date":"2020-01-01T12:59:59+00:00","repo_id":"devlake","branch":"main","commit_sha":"015e3d3b480e417aede5a1293bd61de9b0fd051d"}<br/>
// @Description example2: {"pipeline_name":"A123","name":"unit-test","result":"SUCCESS","status":"DONE","type":"CI/CD","created_date":"2020-01-01T12:00:00+00:00","finished_date":"2020-01-01T12:59:59+00:00","repo_id":"devlake","branch":"main","commit_sha":"015e3d3b480e417aede5a1293bd61de9b0fd051d"}<br/>
// @Description When request webhook first time for each pipeline, it will be created.
// @Description So we suggest request before task start and after pipeline finish.
// @Description Remember fill all data to request after pipeline finish.
// @Tags plugins/webhook
// @Param body body WebhookTaskRequest true "json body"
// @Success 200
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 403  {string} errcode.Error "Forbidden"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/webhook/:connectionId/cicd_tasks [POST]
func PostCicdTask(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	// get request
	request := &WebhookTaskRequest{}
	err = helper.DecodeMapStruct(input.Body, request)
	if err != nil {
		return &core.ApiResourceOutput{Body: err.Error(), Status: http.StatusBadRequest}, nil
	}
	// validate
	vld = validator.New()
	err = errors.BadInput.Wrap(vld.Struct(request), `input json error`)
	if err != nil {
		return &core.ApiResourceOutput{Body: err.Error(), Status: http.StatusBadRequest}, nil
	}

	db := basicRes.GetDal()
	pipelineId := fmt.Sprintf("%s:%d:%s", "webhook", connection.ID, request.PipelineName)
	domainCicdTask := &devops.CICDTask{
		DomainEntity: domainlayer.DomainEntity{
			Id: fmt.Sprintf("%s:%d:%s:%s", "webhook", connection.ID, request.PipelineName, request.Name),
		},
		PipelineId:   pipelineId,
		Name:         request.Name,
		Result:       request.Result,
		Status:       request.Status,
		Type:         request.Type,
		StartedDate:  request.StartedDate,
		FinishedDate: request.FinishedDate,
	}
	if domainCicdTask.FinishedDate != nil {
		domainCicdTask.DurationSec = uint64(domainCicdTask.FinishedDate.Sub(domainCicdTask.StartedDate).Seconds())
	}

	domainPipeline := &devops.CICDPipeline{}
	err = db.First(domainPipeline, dal.Where("id = ?", pipelineId))
	if err != nil {
		domainPipeline = &devops.CICDPipeline{
			DomainEntity: domainlayer.DomainEntity{
				Id: pipelineId,
			},
			Name:         request.PipelineName,
			Result:       ``,
			Status:       `IN_PROGRESS`,
			Type:         ``,
			CreatedDate:  request.StartedDate,
			FinishedDate: nil,
		}
	} else if domainPipeline.Status == `DONE` {
		return nil, errors.Forbidden.New(`can not receive this task because pipeline has already been done.`)
	}

	domainPipelineRepo := &devops.CiCDPipelineCommit{
		PipelineId: pipelineId,
		CommitSha:  request.CommitSha,
		Branch:     request.Branch,
		RepoId:     request.RepoId,
	}

	// save
	err = db.CreateOrUpdate(domainCicdTask)
	if err != nil {
		return nil, err
	}
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

// PostPipelineFinish
// @Summary set pipeline's status to DONE
// @Description set pipeline's status to DONE and cal duration
// @Tags plugins/webhook
// @Success 200
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/webhook/:connectionId/cicd_pipeline/:pipelineName/finish [POST]
func PostPipelineFinish(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}

	db := basicRes.GetDal()
	pipelineId := fmt.Sprintf("%s:%d:%s", "webhook", connection.ID, input.Params[`pipelineName`])
	println(pipelineId)
	domainPipeline := &devops.CICDPipeline{}
	err = db.First(domainPipeline, dal.Where("id = ?", pipelineId))
	if err != nil {
		return nil, errors.NotFound.Wrap(err, `pipeline not found`)
	}

	domainTasks := []devops.CICDTask{}
	err = db.All(&domainTasks, dal.Where("pipeline_id = ?", pipelineId))
	if err != nil {
		return nil, errors.NotFound.Wrap(err, `tasks not found`)
	}
	typeHasCi, typeHasCd, result := getTypeAndResultFromTasks(domainTasks)
	if typeHasCi && typeHasCd {
		domainPipeline.Type = `CI/CD`
	} else if typeHasCi {
		domainPipeline.Type = `CI`
	} else if typeHasCd {
		domainPipeline.Type = `CD`
	}
	domainPipeline.Result = result
	domainPipeline.Status = ticket.DONE
	now := time.Now()
	domainPipeline.FinishedDate = &now
	domainPipeline.DurationSec = uint64(domainPipeline.FinishedDate.Sub(domainPipeline.CreatedDate).Seconds())

	// save
	err = db.Update(domainPipeline)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
}

func getTypeAndResultFromTasks(domainTasks []devops.CICDTask) (typeHasCi bool, typeHasCd bool, result string) {
	typeHasCi = false
	typeHasCd = false
	result = `SUCCESS`
	for _, domainTask := range domainTasks {
		if domainTask.Type == `CI/CD` {
			typeHasCi = true
			typeHasCd = true
		} else if domainTask.Type == `CI` {
			typeHasCi = true
		} else if domainTask.Type == `CD` {
			typeHasCd = true
		}
		if domainTask.Result == `ABORT` {
			result = `ABORT`
		} else if domainTask.Result == `FAILURE` {
			if result == `SUCCESS` {
				result = `FAILURE`
			}
		}
	}
	return
}
