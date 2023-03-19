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
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/teambition/models"
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

func getTaskWorktimeIdGen() *didgen.DomainIdGenerator {
	if taskWorktimeIdGen == nil {
		taskWorktimeIdGen = didgen.NewDomainIdGenerator(&models.TeambitionTaskWorktime{})
	}
	return taskWorktimeIdGen
}

func CreateRawDataSubTaskArgs(taskCtx plugin.SubTaskContext, rawTable string) (*api.RawDataSubTaskArgs, *TeambitionTaskData) {
	data := taskCtx.GetData().(*TeambitionTaskData)
	filteredData := *data
	filteredData.Options = &TeambitionOptions{}
	*filteredData.Options = *data.Options
	params := TeambitionApiParams{
		ConnectionId:   data.Options.ConnectionId,
		OrganizationId: data.Options.OrganizationId,
		ProjectId:      data.Options.ProjectId,
	}
	rawDataSubTaskArgs := &api.RawDataSubTaskArgs{
		Ctx:    taskCtx,
		Params: params,
		Table:  rawTable,
	}
	return rawDataSubTaskArgs, &filteredData
}
