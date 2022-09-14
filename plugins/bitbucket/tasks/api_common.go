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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"io"
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

type BitbucketPagination struct {
	Values  []interface{} `json:"values"`
	PageLen int           `json:"pagelen"`
	Size    int           `json:"size"`
	Page    int           `json:"page"`
	Next    string        `json:"next"`
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

func GetQuery(reqData *helper.RequestData) (url.Values, errors.Error) {
	query := url.Values{}
	query.Set("state", "all")
	query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
	query.Set("pagelen", fmt.Sprintf("%v", reqData.Pager.Size))

	return query, nil
}

func GetTotalPagesFromResponse(res *http.Response, args *helper.ApiCollectorArgs) (int, errors.Error) {
	body := &BitbucketPagination{}
	err := helper.UnmarshalResponse(res, body)
	if err != nil {
		return 0, err
	}
	pages := body.Size / args.PageSize
	if body.Size%args.PageSize > 0 {
		pages++
	}
	return pages, nil
}

func GetRawMessageFromResponse(res *http.Response) ([]json.RawMessage, errors.Error) {
	var rawMessages struct {
		Values []json.RawMessage `json:"values"`
	}
	if res == nil {
		return nil, errors.Default.New("res is nil")
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("error reading response from %s", res.Request.URL.String()))
	}

	err = errors.Convert(json.Unmarshal(resBody, &rawMessages))
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("error decoding response from %s: raw response: %s", res.Request.URL.String(), string(resBody)))
	}

	return rawMessages.Values, nil
}

func GetPullRequestsIterator(taskCtx core.SubTaskContext) (*helper.DalCursorIterator, errors.Error) {
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

func GetIssuesIterator(taskCtx core.SubTaskContext) (*helper.DalCursorIterator, errors.Error) {
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

func ignoreHTTPStatus404(res *http.Response) errors.Error {
	if res.StatusCode == http.StatusUnauthorized {
		return errors.Unauthorized.New("authentication failed, please check your AccessToken")
	}
	if res.StatusCode == http.StatusNotFound {
		return helper.ErrIgnoreAndContinue
	}
	return nil
}
