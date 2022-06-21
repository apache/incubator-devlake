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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

type Page struct {
	Data Data `json:"data"`
}
type Data struct {
	Count int `json:"count"`
}

var UserIdGen *didgen.DomainIdGenerator
var WorkspaceIdGen *didgen.DomainIdGenerator
var IssueIdGen *didgen.DomainIdGenerator
var IterIdGen *didgen.DomainIdGenerator

// res will not be used
func GetTotalPagesFromResponse(r *http.Response, args *helper.ApiCollectorArgs) (int, error) {
	data := args.Ctx.GetData().(*TapdTaskData)
	apiClient, err := NewTapdApiClient(args.Ctx.TaskContext(), data.Connection)
	if err != nil {
		return 0, err
	}
	query := url.Values{}
	query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
	res, err := apiClient.Get(fmt.Sprintf("%s/count", r.Request.URL.Path), query, nil)
	if err != nil {
		return 0, err
	}
	var page Page
	err = helper.UnmarshalResponse(res, &page)

	count := page.Data.Count
	totalPage := count/args.PageSize + 1

	return totalPage, err
}

func parseIterationChangelog(taskCtx core.SubTaskContext, old string, new string) (uint64, uint64, error) {
	data := taskCtx.GetData().(*TapdTaskData)
	db := taskCtx.GetDal()
	iterationFrom := &models.TapdIteration{}
	clauses := []dal.Clause{
		dal.From(&models.TapdIteration{}),
		dal.Where("connection_id = ? and workspace_id = ? and name = ?",
			data.Connection.ID, data.Options.WorkspaceId, old),
	}
	err := db.First(iterationFrom, clauses...)
	if err != nil && err.Error() != "record not found" {
		return 0, 0, err
	}

	iterationTo := &models.TapdIteration{}
	clauses = []dal.Clause{
		dal.From(&models.TapdIteration{}),
		dal.Where("connection_id = ? and workspace_id = ? and name = ?",
			data.Connection.ID, data.Options.WorkspaceId, new),
	}
	err = db.First(iterationTo, clauses...)
	if err != nil && err.Error() != "record not found" {
		return 0, 0, err
	}
	return iterationFrom.Id, iterationTo.Id, nil
}
func GetRawMessageDirectFromResponse(res *http.Response) ([]json.RawMessage, error) {
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	return []json.RawMessage{body}, nil
}

func GetRawMessageArrayFromResponse(res *http.Response) ([]json.RawMessage, error) {
	var data struct {
		Data []json.RawMessage `json:"data"`
	}
	err := helper.UnmarshalResponse(res, &data)
	return data.Data, err
}

type TapdApiParams struct {
	ConnectionId uint64
	CompanyId    uint64
	WorkspaceId  uint64
}

func CreateRawDataSubTaskArgs(taskCtx core.SubTaskContext, rawTable string) (*helper.RawDataSubTaskArgs, *TapdTaskData) {
	data := taskCtx.GetData().(*TapdTaskData)
	var params = TapdApiParams{
		ConnectionId: data.Options.ConnectionId,
		WorkspaceId:  data.Options.WorkspaceId,
	}
	if data.Options.CompanyId != 0 {
		params.CompanyId = data.Options.CompanyId
	}
	RawDataSubTaskArgs := &helper.RawDataSubTaskArgs{
		Ctx:    taskCtx,
		Params: params,
		Table:  rawTable,
	}
	return RawDataSubTaskArgs, data
}
