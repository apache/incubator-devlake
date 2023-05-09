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

package tasks

import (
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/teambition/models"
	"strings"
)

type TeambitionApiParams struct {
	ConnectionId   uint64
	OrganizationId string
	ProjectId      string
}

type TeambitionComRes[T any] struct {
	NextPageToken string `json:"nextPageToken"`
	TotalSize     int    `json:"totalSize"`
	Result        T      `json:"result"`
	Code          int    `json:"code"`
	ErrorMessage  string `json:"errorMessage"`
	RequestId     string `json:"requestId"`
}

var accountIdGen *didgen.DomainIdGenerator
var taskIdGen *didgen.DomainIdGenerator
var taskActivityIdGen *didgen.DomainIdGenerator
var taskWorktimeIdGen *didgen.DomainIdGenerator
var projectIdGen *didgen.DomainIdGenerator
var sprintIdGen *didgen.DomainIdGenerator

func getAccountIdGen() *didgen.DomainIdGenerator {
	if accountIdGen == nil {
		accountIdGen = didgen.NewDomainIdGenerator(&models.TeambitionAccount{})
	}
	return accountIdGen
}

func getTaskIdGen() *didgen.DomainIdGenerator {
	if taskIdGen == nil {
		taskIdGen = didgen.NewDomainIdGenerator(&models.TeambitionTask{})
	}
	return taskIdGen
}

func getTaskActivityIdGen() *didgen.DomainIdGenerator {
	if taskActivityIdGen == nil {
		taskActivityIdGen = didgen.NewDomainIdGenerator(&models.TeambitionTaskActivity{})
	}
	return taskActivityIdGen
}

func getProjectIdGen() *didgen.DomainIdGenerator {
	if projectIdGen == nil {
		projectIdGen = didgen.NewDomainIdGenerator(&models.TeambitionProject{})
	}
	return projectIdGen
}

func getTaskWorktimeIdGen() *didgen.DomainIdGenerator {
	if taskWorktimeIdGen == nil {
		taskWorktimeIdGen = didgen.NewDomainIdGenerator(&models.TeambitionTaskWorktime{})
	}
	return taskWorktimeIdGen
}

func getSprintIdGen() *didgen.DomainIdGenerator {
	if sprintIdGen == nil {
		sprintIdGen = didgen.NewDomainIdGenerator(&models.TeambitionSprint{})
	}
	return sprintIdGen
}

func CreateRawDataSubTaskArgs(taskCtx plugin.SubTaskContext, rawTable string) (*api.RawDataSubTaskArgs, *TeambitionTaskData) {
	data := taskCtx.GetData().(*TeambitionTaskData)
	filteredData := *data
	filteredData.Options = &TeambitionOptions{}
	*filteredData.Options = *data.Options
	params := TeambitionApiParams{
		ConnectionId: data.Options.ConnectionId,
		ProjectId:    data.Options.ProjectId,
	}
	rawDataSubTaskArgs := &api.RawDataSubTaskArgs{
		Ctx:    taskCtx,
		Params: params,
		Table:  rawTable,
	}
	return rawDataSubTaskArgs, &filteredData
}

func getStdTypeMappings(data *TeambitionTaskData) map[string]string {
	stdTypeMappings := make(map[string]string)
	for userType, stdType := range data.Options.TransformationRules.TypeMappings {
		stdTypeMappings[userType] = strings.ToUpper(stdType.StandardType)
	}
	return stdTypeMappings
}

func getStatusMapping(data *TeambitionTaskData) map[string]string {
	statusMapping := make(map[string]string)
	mapping := data.Options.TransformationRules.StatusMappings
	for std, orig := range mapping {
		for _, v := range orig {
			statusMapping[v] = std
		}
	}
	return statusMapping
}

func FindAccountById(db dal.Dal, accountId string) (*models.TeambitionAccount, errors.Error) {
	if accountId == "" {
		return nil, errors.Default.New("account id must not empty")
	}
	account := &models.TeambitionAccount{}
	if err := db.First(account, dal.Where("user_id = ?", accountId)); err != nil {
		return nil, err
	}
	return account, nil
}

func FindProjectById(db dal.Dal, projectId string) (*models.TeambitionProject, errors.Error) {
	if projectId == "" {
		return nil, errors.Default.New("project id must not empty")
	}
	project := &models.TeambitionProject{}
	if err := db.First(project, dal.Where("id = ?", projectId)); err != nil {
		return nil, err
	}
	return project, nil
}

func FindTaskScenarioById(db dal.Dal, scenarioId string) (*models.TeambitionTaskScenario, errors.Error) {
	if scenarioId == "" {
		return nil, errors.Default.New("id must not empty")
	}
	scenario := &models.TeambitionTaskScenario{}
	if err := db.First(scenario, dal.Where("id = ?", scenarioId)); err != nil {
		return nil, err
	}
	return scenario, nil
}

func FindTaskFlowStatusById(db dal.Dal, id string) (*models.TeambitionTaskFlowStatus, errors.Error) {
	if id == "" {
		return nil, errors.Default.New("id must not empty")
	}
	status := &models.TeambitionTaskFlowStatus{}
	if err := db.First(status, dal.Where("id = ?", id)); err != nil {
		return nil, err
	}
	return status, nil
}
