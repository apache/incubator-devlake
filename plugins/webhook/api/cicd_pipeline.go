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
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/webhook/models"
	"github.com/go-playground/validator/v10"
)

type WebhookTaskRequest struct {
	// PipelineName can be filled by any string unique in one pipeline
	PipelineName string `mapstructure:"pipeline_name" validate:"required"`

	Name         string     `validate:"required"` // Name should be unique in one pipeline
	Result       string     `validate:"oneof=SUCCESS FAILURE ABORT IN_PROGRESS"`
	Status       string     `validate:"oneof=IN_PROGRESS DONE"`
	Type         string     `validate:"oneof=TEST LINT BUILD DEPLOYMENT"`
	Environment  string     `validate:"oneof=PRODUCTION STAGING TESTING"`
	StartedDate  time.Time  `mapstructure:"created_date" validate:"required"`
	FinishedDate *time.Time `mapstructure:"finished_date"`

	RepoId    string `mapstructure:"repo_id" validate:"required"` // RepoId should be unique string
	Branch    string
	CommitSha string `mapstructure:"commit_sha"`
}

// PostCicdTask
// @Summary create pipeline by webhook
// @Description Create pipeline by webhook.<br/>
// @Description example1: {"pipeline_name":"A123","name":"unit-test","result":"IN_PROGRESS","status":"IN_PROGRESS","type":"TEST","environment":"PRODUCTION","created_date":"2020-01-01T12:00:00+00:00","finished_date":"2020-01-01T12:59:59+00:00","repo_id":"devlake","branch":"main","commit_sha":"015e3d3b480e417aede5a1293bd61de9b0fd051d"}<br/>
// @Description example2: {"pipeline_name":"A123","name":"unit-test","result":"SUCCESS","status":"DONE","type":"DEPLOYMENT","environment":"PRODUCTION","created_date":"2020-01-01T12:00:00+00:00","finished_date":"2020-01-01T12:59:59+00:00","repo_id":"devlake","branch":"main","commit_sha":"015e3d3b480e417aede5a1293bd61de9b0fd051d"}<br/>
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
	err = errors.Convert(vld.Struct(request))
	if err != nil {
		return nil, errors.BadInput.Wrap(vld.Struct(request), `input json error`)
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
		Environment:  request.Environment,
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

	domainPipelineCommit := &devops.CiCDPipelineCommit{
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
	err = db.CreateOrUpdate(domainPipelineCommit)
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

	now := time.Now()

	// finished all CICDTask
	cursor, err := db.Cursor(
		dal.From(&devops.CICDTask{}),
		dal.Where("pipeline_id = ?", pipelineId),
	)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on select CICDTask")
	}
	batch, err := helper.NewBatchSave(basicRes, reflect.TypeOf(&devops.CICDTask{}), 500)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error getting batch from CICDTask")
	}
	defer batch.Close()

	domainTasks := []devops.CICDTask{}
	for cursor.Next() {
		task := &devops.CICDTask{}
		err = db.Fetch(cursor, task)
		if err != nil {
			return nil, errors.Default.Wrap(err, "error on Fetch CICDTask data")
		}
		// set the IN_PROGRESS task to be ABORT
		if task.Result == `IN_PROGRESS` {
			task.Result = `ABORT`
			task.FinishedDate = &now
		}
		task.Status = ticket.DONE
		domainTasks = append(domainTasks, *task)
		err = batch.Add(task)
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("error on CICDTask batch add %v", task))
		}
	}

	// finished CICDPipeline
	domainPipeline := &devops.CICDPipeline{}
	err = db.First(domainPipeline, dal.Where("id = ?", pipelineId))
	if err != nil {
		return nil, errors.NotFound.Wrap(err, `pipeline not found`)
	}

	err = db.All(&domainTasks, dal.Where("pipeline_id = ?", pipelineId))
	if err != nil {
		return nil, errors.NotFound.Wrap(err, `tasks not found`)
	}
	pipelineType, result := getTypeAndResultFromTasks(domainTasks)
	domainPipeline.Type = pipelineType
	domainPipeline.Result = result
	domainPipeline.Status = ticket.DONE
	domainPipeline.FinishedDate = &now
	domainPipeline.DurationSec = uint64(domainPipeline.FinishedDate.Sub(domainPipeline.CreatedDate).Seconds())

	// save
	err = db.Update(domainPipeline)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
}

// getTypeAndResultFromTasks will extract pipeline type and result from tasks
// type = tasks' type if all tasks have the same type, or empty string
// result = ABORT if any tasks' type is ABORT,
// or result = FAILURE if any tasks' type is ABORT and others are SUCCESS
// or result = SUCCESS if all tasks' type is SUCCESS
func getTypeAndResultFromTasks(domainTasks []devops.CICDTask) (pipelineType, result string) {
	result = `SUCCESS`
	if len(domainTasks) > 0 {
		pipelineType = domainTasks[0].Type
	}
	for _, domainTask := range domainTasks {
		if domainTask.Type != pipelineType {
			pipelineType = ``
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

type WebhookDeployTaskRequest struct {
	// RepoUrl should be unique string, fill url or other unique data
	RepoUrl   string `mapstructure:"repo_url" validate:"required"`
	CommitSha string `mapstructure:"commit_sha" validate:"required"`
	// start_time and end_time is more readable for users,
	// StartedDate and FinishedDate is same as columns in db.
	// So they all keep.
	StartedDate  *time.Time `mapstructure:"start_time" validate:"required_with=FinishedDate"`
	FinishedDate *time.Time `mapstructure:"end_time"`
	Environment  string     `validate:"omitempty,oneof=PRODUCTION STAGING TESTING DEVELOPMENT"`
}

// PostCicdTask
// @Summary create deployment pipeline by webhook
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
func PostDeploymentCicdTask(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	// get request
	request := &WebhookDeployTaskRequest{}
	err = helper.DecodeMapStruct(input.Body, request)
	if err != nil {
		return &core.ApiResourceOutput{Body: err.Error(), Status: http.StatusBadRequest}, nil
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
	pipelineId := fmt.Sprintf("%s:%d:%s:%s:%s", "webhook", connection.ID, `pipeline`, urlHash16, request.CommitSha)

	taskId := fmt.Sprintf("%s:%d:%s:%s", "webhook", connection.ID, urlHash16, request.CommitSha)
	domainCicdTask := &devops.CICDTask{
		DomainEntity: domainlayer.DomainEntity{
			Id: taskId,
		},
		PipelineId:  pipelineId,
		Name:        fmt.Sprintf(`deployment for %s`, request.CommitSha),
		Result:      devops.SUCCESS,
		Status:      devops.DONE,
		Type:        devops.DEPLOYMENT,
		Environment: request.Environment,
		CicdScopeId: scopeId,
	}
	now := time.Now()
	if request.StartedDate != nil {
		domainCicdTask.StartedDate = *request.StartedDate
		if request.FinishedDate != nil {
			domainCicdTask.FinishedDate = request.FinishedDate
		} else {
			domainCicdTask.FinishedDate = &now
		}
		domainCicdTask.DurationSec = uint64(domainCicdTask.FinishedDate.Sub(domainCicdTask.StartedDate).Seconds())
	} else {
		domainCicdTask.StartedDate = now
	}
	if domainCicdTask.Environment == `` {
		domainCicdTask.Environment = devops.PRODUCTION
	}

	domainPipeline := &devops.CICDPipeline{
		DomainEntity: domainlayer.DomainEntity{
			Id: pipelineId,
		},
		Name:         fmt.Sprintf(`pipeline for %s`, request.CommitSha),
		Result:       devops.SUCCESS,
		Status:       devops.DONE,
		Type:         devops.DEPLOYMENT,
		CreatedDate:  domainCicdTask.StartedDate,
		FinishedDate: domainCicdTask.FinishedDate,
		DurationSec:  domainCicdTask.DurationSec,
		Environment:  domainCicdTask.Environment,
		CicdScopeId:  scopeId,
	}

	domainPipelineCommit := &devops.CiCDPipelineCommit{
		PipelineId: pipelineId,
		CommitSha:  request.CommitSha,
		Branch:     ``,
		RepoId:     request.RepoUrl,
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
	err = db.CreateOrUpdate(domainPipelineCommit)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
}
