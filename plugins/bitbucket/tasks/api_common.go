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
	"github.com/apache/incubator-devlake/plugins/bitbucket/utils"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

type BitbucketApiParams struct {
	ConnectionId uint64
	Owner        string
	Repo         string
}

type BitbucketInput struct {
	BitbucketId int
}

func CreateRawDataSubTaskArgs(taskCtx core.SubTaskContext, Table string) (*helper.RawDataSubTaskArgs, *BitbucketTaskData) {
	data := taskCtx.GetData().(*BitbucketTaskData)
	RawDataSubTaskArgs := &helper.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: BitbucketApiParams{
			ConnectionId: data.Options.ConnectionId,
			Owner:        data.Options.Owner,
			Repo:         data.Options.Repo,
		},
		Table: Table,
	}
	return RawDataSubTaskArgs, data
}

func GetQuery(reqData *helper.RequestData) (url.Values, error) {
	query := url.Values{}
	query.Set("state", "all")
	query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
	query.Set("pagelen", fmt.Sprintf("%v", reqData.Pager.Size))

	return query, nil
}

func GetTotalPagesFromResponse(res *http.Response, args *helper.ApiCollectorArgs) (int, error) {
	link := res.Header.Get("link")
	pageInfo, err := utils.GetPagingFromLinkHeader(link)
	if err != nil {
		return 0, nil
	}
	return pageInfo.Last, nil
}

func GetRawMessageFromResponse(res *http.Response) ([]json.RawMessage, error) {
	var rawMessages struct {
		Values []json.RawMessage `json:"values"`
	}
	if res == nil {
		return nil, fmt.Errorf("res is nil")
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("%w %s", err, res.Request.URL.String())
	}

	err = json.Unmarshal(resBody, &rawMessages)
	if err != nil {
		return nil, fmt.Errorf("%w %s %s", err, res.Request.URL.String(), string(resBody))
	}

	return rawMessages.Values, nil
}

func GetPullRequestsIterator(taskCtx core.SubTaskContext) (*helper.DalCursorIterator, error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*BitbucketTaskData)
	clauses := []dal.Clause{
		dal.Select("bpr.bitbucket_id"),
		dal.From("_tool_bitbucket_pull_requests bpr"),
		dal.Where(
			`bpr.repo_id = ? and bpr.connection_id = ?`,
			"repositories/"+data.Options.Owner+"/"+data.Options.Repo, data.Options.ConnectionId,
		),
	}
	// construct the input iterator
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return nil, err
	}

	return helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(BitbucketInput{}))
}

func GetIssuesIterator(taskCtx core.SubTaskContext) (*helper.DalCursorIterator, error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*BitbucketTaskData)
	clauses := []dal.Clause{
		dal.Select("bpr.bitbucket_id"),
		dal.From("_tool_bitbucket_issues bpr"),
		dal.Where(
			`bpr.repo_id = ? and bpr.connection_id = ?`,
			"repositories/"+data.Options.Owner+"/"+data.Options.Repo, data.Options.ConnectionId,
		),
	}
	// construct the input iterator
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return nil, err
	}

	return helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(BitbucketInput{}))
}
