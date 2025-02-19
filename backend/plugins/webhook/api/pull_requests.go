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
	"net/http"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/log"

	"github.com/apache/incubator-devlake/helpers/dbhelper"
	"github.com/go-playground/validator/v10"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/webhook/models"
)

type WebhookPullRequestReq struct {
	Id             string     `mapstructure:"id" validate:"required"`
	BaseRepoId     string     `mapstructure:"baseRepoId"`
	HeadRepoId     string     `mapstructure:"headRepoId"`
	Status         string     `mapstructure:"status" validate:"omitempty,oneof=OPEN CLOSED MERGED"`
	OriginalStatus string     `mapstructure:"originalStatus"`
	Title          string     `mapstructure:"displayTitle" validate:"required"`
	Description    string     `mapstructure:"description"`
	Url            string     `mapstructure:"url"`
	AuthorName     string     `mapstructure:"authorName"`
	AuthorId       string     `mapstructure:"authorId"`
	MergedByName   string     `mapstructure:"mergedByName"`
	MergedById     string     `mapstructure:"mergedById"`
	ParentPrId     string     `mapstructure:"parentPrId"`
	PullRequestKey int        `mapstructure:"pullRequestKey" validate:"required"`
	CreatedDate    *time.Time `mapstructure:"createdDate" validate:"required"`
	MergedDate     *time.Time `mapstructure:"mergedDate"`
	ClosedDate     *time.Time `mapstructure:"closedDate"`
	Type           string     `mapstructure:"type"`
	Component      string     `mapstructure:"component"`
	MergeCommitSha string     `mapstructure:"mergeCommitSha"`
	HeadRef        string     `mapstructure:"headRef"`
	BaseRef        string     `mapstructure:"baseRef"`
	BaseCommitSha  string     `mapstructure:"baseCommitSha"`
	HeadCommitSha  string     `mapstructure:"headCommitSha"`
	Additions      int        `mapstructure:"additions"`
	Deletions      int        `mapstructure:"deletions"`
	IsDraft        bool       `mapstructure:"isDraft"`
	// PullRequestCommits is used for multiple commits in one pull request
	//PullRequestCommits []WebhookPullRequestCommitReq `mapstructure:"pullRequestCommits" validate:"omitempty,dive"`
}

//type WebhookPullRequestCommitReq struct {
//	DisplayTitle string     `mapstructure:"displayTitle"`
//	RepoId       string     `mapstructure:"repoId"`
//	RepoUrl      string     `mapstructure:"repoUrl" validate:"required"`
//	Name         string     `mapstructure:"name"`
//	RefName      string     `mapstructure:"refName"`
//	CommitSha    string     `mapstructure:"commitSha" validate:"required"`
//	CommitMsg    string     `mapstructure:"commitMsg"`
//	Result       string     `mapstructure:"result"`
//	Status       string     `mapstructure:"status"`
//	CreatedDate  *time.Time `mapstructure:"createdDate"`
//	// QueuedDate   *time.Time `mapstructure:"queue_time"`
//	StartedDate  *time.Time `mapstructure:"startedDate" validate:"required"`
//	FinishedDate *time.Time `mapstructure:"finishedDate" validate:"required"`
//}

// PostPullRequests
// @Summary create pull requests by webhook
// @Description Create pull request by webhook.<br/>
// @Description example1: {"repo_url":"devlake","commit_sha":"015e3d3b480e417aede5a1293bd61de9b0fd051d","start_time":"2020-01-01T12:00:00+00:00","end_time":"2020-01-01T12:59:59+00:00","environment":"PRODUCTION"}<br/>
// @Description So we suggest request before task after deployment pipeline finish.
// @Description Both cicd_pipeline and cicd_task will be created
// @Tags plugins/webhook
// @Param body body WebhookPullRequestReq true "json body"
// @Success 200
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 403  {string} errcode.Error "Forbidden"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/webhook/connections/:connectionId/pullrequests [POST]
func PostPullRequests(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.First(connection, input.Params)

	return postPullRequests(input, connection, err)
}

// PostPullRequestsByName
// @Summary create pull requests by webhook name
// @Description Create pull request by webhook name.<br/>
// @Description example1: {"repo_url":"devlake","commit_sha":"015e3d3b480e417aede5a1293bd61de9b0fd051d","start_time":"2020-01-01T12:00:00+00:00","end_time":"2020-01-01T12:59:59+00:00","environment":"PRODUCTION"}<br/>
// @Description So we suggest request before task after deployment pipeline finish.
// @Description Both cicd_pipeline and cicd_task will be created
// @Tags plugins/webhook
// @Param body body WebhookPullRequestReq true "json body"
// @Success 200
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 403  {string} errcode.Error "Forbidden"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/webhook/connections/by-name/:connectionName/pullrequests [POST]
func PostPullRequestsByName(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.FirstByName(connection, input.Params)

	return postPullRequests(input, connection, err)
}

func postPullRequests(input *plugin.ApiResourceInput, connection *models.WebhookConnection, err errors.Error) (*plugin.ApiResourceOutput, errors.Error) {
	if err != nil {
		return nil, err
	}
	// get request
	request := &WebhookPullRequestReq{}
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
	if err := CreatePullRequest(connection, request, tx, logger); err != nil {
		logger.Error(err, "create pull requests")
		return nil, err
	}

	return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
}

func CreatePullRequest(connection *models.WebhookConnection, request *WebhookPullRequestReq, tx dal.Transaction, logger log.Logger) errors.Error {
	// validation
	if request == nil {
		return errors.BadInput.New("request body is nil")
	}
	//if len(request.PullRequestCommits) == 0 {
	//	return errors.BadInput.New("pull_request_commits is empty")
	//}
	// set default values for optional fields
	// prepare pull request commits and pull request records
	// queuedDuration := dateInfo.CalculateQueueDuration()
	//prCommits := make([]*code.PullRequestCommit, len(request.PullRequestCommits))
	//for i, commit := range request.PullRequestCommits {
	//	if commit.Name == "" {
	//		commit.Name = fmt.Sprintf(`commit for %s`, commit.CommitSha)
	//	}
	//createdDate := time.Now()
	//if request.CreatedDate == nil {
	//	request.CreatedDate = &createdDate
	//}
	//	if commit.StartedDate == nil {
	//		commit.StartedDate = request.StartedDate
	//	}
	//	if commit.FinishedDate == nil {
	//		commit.FinishedDate = request.FinishedDate
	//	}
	//	// create a pull_request_commits record
	//	prCommits[i] = &code.PullRequestCommit{
	//		CommitSha:          commit.CommitSha,
	//		PullRequestId:      request.Id,
	//		CommitAuthorName:   commit.AuthorName,
	//		CommitAuthorEmail:  commit.AuthorEmail,
	//		CommitAuthoredDate: *commit.CreatedDate,
	//		NoPKModel: common.NoPKModel{
	//			CreatedAt: *commit.CreatedDate,
	//			UpdatedAt: commit.UpdatedAt,
	//		},
	//	}
	//}

	//if err := tx.CreateOrUpdate(prCommits); err != nil {
	//	logger.Error(err, "failed to save pull request commits")
	//	return err
	//}

	// create a pull_request record
	pullRequest := code.PullRequest{
		DomainEntity: domainlayer.DomainEntity{
			Id: fmt.Sprintf("%s:%d:%d", "webhook", connection.ID, request.PullRequestKey),
		},
		BaseRepoId:     request.BaseRepoId,
		HeadRepoId:     request.HeadRepoId,
		Status:         request.Status,
		OriginalStatus: request.OriginalStatus,
		Title:          request.Title,
		Description:    request.Description,
		Url:            request.Url,
		AuthorName:     request.AuthorName,
		AuthorId:       request.AuthorId,
		MergedByName:   request.MergedByName,
		MergedById:     request.MergedById,
		ParentPrId:     request.ParentPrId,
		PullRequestKey: request.PullRequestKey,
		CreatedDate:    *request.CreatedDate,
		MergedDate:     request.MergedDate,
		ClosedDate:     request.ClosedDate,
		Type:           request.Type,
		Component:      request.Component,
		MergeCommitSha: request.MergeCommitSha,
		HeadRef:        request.HeadRef,
		BaseRef:        request.BaseRef,
		BaseCommitSha:  request.BaseCommitSha,
		HeadCommitSha:  request.HeadCommitSha,
		Additions:      request.Additions,
		Deletions:      request.Deletions,
		IsDraft:        request.IsDraft,
	}
	if err := tx.CreateOrUpdate(pullRequest); err != nil {
		logger.Error(err, "failed to save pull request")
		return err
	}
	return nil
}
