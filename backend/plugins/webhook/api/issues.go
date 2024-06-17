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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/webhook/models"

	"github.com/go-playground/validator/v10"
)

type WebhookIssueRequest struct {
	Url                     string     `mapstructure:"url"`
	IssueKey                string     `mapstructure:"issueKey" validate:"required"`
	Title                   string     `mapstructure:"title" validate:"required"`
	Description             string     `mapstructure:"description"`
	EpicKey                 string     `mapstructure:"epicKey"`
	Type                    string     `mapstructure:"type"`
	Status                  string     `mapstructure:"status" validate:"oneof=TODO DONE IN_PROGRESS"`
	OriginalStatus          string     `mapstructure:"originalStatus" validate:"required"`
	StoryPoint              float64    `mapstructure:"storyPoint"`
	ResolutionDate          *time.Time `mapstructure:"resolutionDate"`
	CreatedDate             *time.Time `mapstructure:"createdDate" validate:"required"`
	UpdatedDate             *time.Time `mapstructure:"updatedDate"`
	LeadTimeMinutes         uint       `mapstructure:"leadTimeMinutes"`
	ParentIssueKey          string     `mapstructure:"parentIssueKey"`
	Priority                string     `mapstructure:"priority"`
	OriginalEstimateMinutes int64      `mapstructure:"originalEstimateMinutes"`
	TimeSpentMinutes        int64      `mapstructure:"timeSpentMinutes"`
	TimeRemainingMinutes    int64      `mapstructure:"timeRemainingMinutes"`
	CreatorId               string     `mapstructure:"creatorId"`
	CreatorName             string     `mapstructure:"creatorName"`
	AssigneeId              string     `mapstructure:"assigneeId"`
	AssigneeName            string     `mapstructure:"assigneeName"`
	Severity                string     `mapstructure:"severity"`
	Component               string     `mapstructure:"component"`
	//IconURL               string
	//DeploymentId          string
}

// PostIssue
// @Summary receive a record as defined and save it
// @Description receive a record as follow and save it, example: {"url":"","issue_key":"DLK-1234","title":"a feature from DLK","description":"","epic_key":"","type":"BUG","status":"TODO","original_status":"created","story_point":0,"resolution_date":null,"created_date":"2020-01-01T12:00:00+00:00","updated_date":null,"lead_time_minutes":0,"parent_issue_key":"DLK-1200","priority":"","original_estimate_minutes":0,"time_spent_minutes":0,"time_remaining_minutes":0,"creator_id":"user1131","creator_name":"Nick name 1","assignee_id":"user1132","assignee_name":"Nick name 2","severity":"","component":""}
// @Tags plugins/webhook
// @Param body body WebhookIssueRequest true "json body"
// @Success 200  {string} noResponse ""
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/webhook/:connectionId/issues [POST]
func PostIssue(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	// get request
	request := &WebhookIssueRequest{}
	err = helper.DecodeMapStruct(input.Body, request, true)
	if err != nil {
		return &plugin.ApiResourceOutput{Body: err.Error(), Status: http.StatusBadRequest}, nil
	}
	// validate
	vld = validator.New()
	err = errors.Convert(vld.Struct(request))
	if err != nil {
		return &plugin.ApiResourceOutput{Body: err.Error(), Status: http.StatusBadRequest}, nil
	}
	db := basicRes.GetDal()
	domainIssue := &ticket.Issue{
		DomainEntity: domainlayer.DomainEntity{
			Id: fmt.Sprintf("%s:%d:%s", "webhook", connection.ID, request.IssueKey),
		},
		Url:                     request.Url,
		IssueKey:                request.IssueKey,
		Title:                   request.Title,
		Description:             request.Description,
		EpicKey:                 request.EpicKey,
		Type:                    request.Type,
		Status:                  request.Status,
		OriginalStatus:          request.OriginalStatus,
		StoryPoint:              &request.StoryPoint,
		ResolutionDate:          request.ResolutionDate,
		CreatedDate:             request.CreatedDate,
		UpdatedDate:             request.UpdatedDate,
		LeadTimeMinutes:         &request.LeadTimeMinutes,
		Priority:                request.Priority,
		OriginalEstimateMinutes: &request.OriginalEstimateMinutes,
		TimeSpentMinutes:        &request.TimeSpentMinutes,
		TimeRemainingMinutes:    &request.TimeRemainingMinutes,
		CreatorName:             request.CreatorName,
		AssigneeName:            request.AssigneeName,
		Severity:                request.Severity,
		Component:               request.Component,
	}
	if request.CreatorId != "" {
		domainIssue.CreatorId = fmt.Sprintf("%s:%d:%s", "webhook", connection.ID, request.CreatorId)
	}
	if request.AssigneeId != "" {
		domainIssue.AssigneeId = fmt.Sprintf("%s:%d:%s", "webhook", connection.ID, request.AssigneeId)
	}
	if request.ParentIssueKey != "" {
		domainIssue.ParentIssueId = fmt.Sprintf("%s:%d:%s", "webhook", connection.ID, request.ParentIssueKey)
	}

	domainBoardId := fmt.Sprintf("%s:%d", "webhook", connection.ID)

	boardIssue := &ticket.BoardIssue{
		BoardId: domainBoardId,
		IssueId: domainIssue.Id,
	}

	// check if board exists
	count, err := db.Count(dal.From(&ticket.Board{}), dal.Where("id = ?", domainBoardId))
	if err != nil {
		return nil, err
	}

	// only create board with domainBoard non-existent
	if count == 0 {
		domainBoard := &ticket.Board{
			DomainEntity: domainlayer.DomainEntity{
				Id: domainBoardId,
			},
		}
		err = db.Create(domainBoard)
		if err != nil {
			return nil, err
		}
	}

	// save
	err = db.CreateOrUpdate(domainIssue)
	if err != nil {
		return nil, err
	}

	err = db.CreateOrUpdate(boardIssue)
	if err != nil {
		return nil, err
	}

	return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
}

// CloseIssue
// @Summary set issue's status to DONE
// @Description set issue's status to DONE
// @Tags plugins/webhook
// @Success 200  {string} noResponse ""
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/webhook/:connectionId/issue/:issueKey/close [POST]
func CloseIssue(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}

	db := basicRes.GetDal()
	domainIssue := &ticket.Issue{}
	err = db.First(domainIssue, dal.Where("id = ?", fmt.Sprintf("%s:%d:%s", "webhook", connection.ID, input.Params[`issueKey`])))
	if err != nil {
		return nil, errors.NotFound.Wrap(err, `issue not found`)
	}
	domainIssue.Status = ticket.DONE
	domainIssue.OriginalStatus = ``

	// save
	err = db.Update(domainIssue)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
}
