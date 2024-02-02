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
	"strings"
	"time"

	"github.com/apache/incubator-devlake/helpers/dbhelper"
	"github.com/go-playground/validator/v10"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/webhook/models"
)

type WebhookDeployTaskRequest struct {
	PipelineId string `mapstructure:"pipeline_id"`
	// RepoUrl should be unique string, fill url or other unique data
	RepoId string `mapstructure:"repo_id"`
	Result string `mapstructure:"result"`
	// start_time and end_time is more readable for users,
	// StartedDate and FinishedDate is same as columns in db.
	// So they all keep.
	CreatedDate *time.Time `mapstructure:"create_time"`
	// QueuedDate   *time.Time `mapstructure:"queue_time"`
	StartedDate  *time.Time `mapstructure:"start_time" validate:"required"`
	FinishedDate *time.Time `mapstructure:"end_time"`
	RepoUrl      string     `mapstructure:"repo_url"`
	Environment  string     `validate:"omitempty,oneof=PRODUCTION STAGING TESTING DEVELOPMENT"`
	Name         string     `mapstructure:"name"`
	RefName      string     `mapstructure:"ref_name"`
	CommitSha    string     `mapstructure:"commit_sha"`
	CommitMsg    string     `mapstructure:"commit_msg"`
	// DeploymentCommits is used for multiple commits in one deployment
	DeploymentCommits []DeploymentCommit `mapstructure:"deploymentCommits" validate:"omitempty,dive"`
}

type DeploymentCommit struct {
	RepoUrl   string `mapstructure:"repo_url" validate:"required"`
	Name      string `mapstructure:"name"`
	RefName   string `mapstructure:"ref_name"`
	CommitSha string `mapstructure:"commit_sha" validate:"required"`
	CommitMsg string `mapstructure:"commit_msg"`
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
// @Router /plugins/webhook/connections/:connectionId/deployments [POST]
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
	txHelper := dbhelper.NewTxHelper(basicRes, &err)
	defer txHelper.End()
	tx := txHelper.Begin()

	pipelineId := request.PipelineId
	scopeId := fmt.Sprintf("%s:%d", "webhook", connection.ID)
	if request.CreatedDate == nil {
		request.CreatedDate = request.StartedDate
	}
	if request.FinishedDate == nil {
		now := time.Now()
		request.FinishedDate = &now
	}
	if request.Result == "" {
		request.Result = devops.RESULT_SUCCESS
	}
	if request.Environment == "" {
		request.Environment = devops.PRODUCTION
	}
	duration := float64(request.FinishedDate.Sub(*request.StartedDate).Milliseconds() / 1e3)
	name := request.Name
	if name == "" {
		if request.DeploymentCommits == nil {
			name = fmt.Sprintf(`deployment for %s`, request.CommitSha)
		} else {
			var commitShaList []string
			for _, commit := range request.DeploymentCommits {
				commitShaList = append(commitShaList, commit.CommitSha)
			}
			name = fmt.Sprintf(`deployment for %s`, strings.Join(commitShaList, ","))
		}
	}
	createdDate := time.Now()
	if request.CreatedDate != nil {
		createdDate = *request.CreatedDate
	} else if request.StartedDate != nil {
		createdDate = *request.StartedDate
	}
	dateInfo := devops.TaskDatesInfo{
		CreatedDate: createdDate,
		// QueuedDate:   request.QueuedDate,
		StartedDate:  request.StartedDate,
		FinishedDate: request.FinishedDate,
	}
	// queuedDuration := dateInfo.CalculateQueueDuration()
	if request.DeploymentCommits == nil {
		if request.CommitSha == "" || request.RepoUrl == "" {
			return nil, errors.Convert(fmt.Errorf("commit_sha or repo_url is required"))
		}
		urlHash16 := fmt.Sprintf("%x", md5.Sum([]byte(request.RepoUrl)))[:16]
		deploymentCommitId := fmt.Sprintf("%s:%d:%s:%s", "webhook", connection.ID, urlHash16, request.CommitSha)
		if pipelineId == "" {
			pipelineId = deploymentCommitId
		}
		// create a deployment_commit record
		deploymentCommit := &devops.CicdDeploymentCommit{
			DomainEntity: domainlayer.DomainEntity{
				Id: deploymentCommitId,
			},
			CicdDeploymentId: pipelineId,
			CicdScopeId:      scopeId,
			Name:             name,
			Result:           request.Result,
			Status:           devops.STATUS_DONE,
			OriginalResult:   request.Result,
			OriginalStatus:   devops.STATUS_DONE,
			TaskDatesInfo:    dateInfo,
			DurationSec:      &duration,
			//QueuedDurationSec: queuedDuration,
			RepoId:              request.RepoId,
			RepoUrl:             request.RepoUrl,
			Environment:         request.Environment,
			OriginalEnvironment: request.Environment,
			RefName:             request.RefName,
			CommitSha:           request.CommitSha,
			CommitMsg:           request.CommitMsg,
		}
		err = tx.CreateOrUpdate(deploymentCommit)
		if err != nil {
			logger.Error(err, "create deployment commit")
			return nil, err
		}

		// create a deployment record
		if err = tx.CreateOrUpdate(deploymentCommit.ToDeployment()); err != nil {
			logger.Error(err, "create deployment")
			return nil, err
		}
	} else {
		for _, commit := range request.DeploymentCommits {
			urlHash16 := fmt.Sprintf("%x", md5.Sum([]byte(commit.RepoUrl)))[:16]
			deploymentCommitId := fmt.Sprintf("%s:%d:%s:%s", "webhook", connection.ID, urlHash16, commit.CommitSha)
			if pipelineId == "" {
				pipelineId = deploymentCommitId
			}
			// create a deployment_commit record
			deploymentCommit := &devops.CicdDeploymentCommit{
				DomainEntity: domainlayer.DomainEntity{
					Id: deploymentCommitId,
				},
				CicdDeploymentId: pipelineId,
				CicdScopeId:      scopeId,
				Result:           request.Result,
				Status:           devops.STATUS_DONE,
				OriginalResult:   request.Result,
				OriginalStatus:   devops.STATUS_DONE,
				TaskDatesInfo:    dateInfo,
				DurationSec:      &duration,
				//QueuedDurationSec: queuedDuration,
				RepoId:              request.RepoId,
				Name:                fmt.Sprintf(`deployment for %s`, commit.CommitSha),
				RepoUrl:             commit.RepoUrl,
				Environment:         request.Environment,
				OriginalEnvironment: request.Environment,
				RefName:             commit.RefName,
				CommitSha:           commit.CommitSha,
				CommitMsg:           commit.CommitMsg,
			}
			err = tx.CreateOrUpdate(deploymentCommit)
			if err != nil {
				logger.Error(err, "create deployment commit")
				return nil, err
			}

			// create a deployment record
			deploymentCommit.Name = name
			if err = tx.CreateOrUpdate(deploymentCommit.ToDeployment()); err != nil {
				logger.Error(err, "create deployment")
				return nil, err
			}
		}
	}

	return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
}
