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
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type BitbucketServerApiParams struct {
	ConnectionId uint64
	FullName     string
}

type BitbucketServerInput struct {
	BitbucketId int
}

type BitbucketServerBranchInput struct {
	Branch string
}

type BitbucketServerCommitInput struct {
	CommitSha string
}

func CreateRawDataSubTaskArgs(taskCtx plugin.SubTaskContext, Table string) (*api.RawDataSubTaskArgs, *BitbucketServerTaskData) {
	data := taskCtx.GetData().(*BitbucketServerTaskData)
	RawDataSubTaskArgs := &api.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: BitbucketServerApiParams{
			ConnectionId: data.Options.ConnectionId,
			FullName:     data.Options.FullName,
		},
		Table: Table,
	}
	return RawDataSubTaskArgs, data
}

func decodeResponse(res *http.Response, message interface{}) errors.Error {
	if res == nil {
		return errors.Default.New("res is nil")
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error reading response from %s", res.Request.URL.String()))
	}

	err = errors.Convert(json.Unmarshal(resBody, &message))
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error decoding response from %s: raw response: %s", res.Request.URL.String(), string(resBody)))
	}
	return nil
}

func GetNextPageCustomData(_ *helper.RequestData, prevPageResponse *http.Response) (interface{}, errors.Error) {
	var rawMessages struct {
		NextPageStart *int `json:"nextPageStart"`
		IsLastPage    bool `json:"isLastPage"`
	}
	err := decodeResponse(prevPageResponse, &rawMessages)
	if err != nil {
		return nil, err
	}

	if rawMessages.IsLastPage || rawMessages.NextPageStart == nil {
		return nil, api.ErrFinishCollect
	}

	return strconv.Itoa(*rawMessages.NextPageStart), nil
}

func GetQueryForNextPage(reqData *helper.RequestData) (url.Values, errors.Error) {
	query := url.Values{}
	query.Set("state", "all")
	query.Set("limit", fmt.Sprintf("%v", reqData.Pager.Size))

	if reqData.CustomData != nil {
		query.Set("start", reqData.CustomData.(string))
	}
	return query, nil
}

func GetRawMessageFromResponse(res *http.Response) ([]json.RawMessage, errors.Error) {
	var rawMessages struct {
		Values []json.RawMessage `json:"values"`
	}
	err := decodeResponse(res, &rawMessages)
	if err != nil {
		return nil, err
	}

	return rawMessages.Values, nil
}

func GetBranchesIterator(taskCtx plugin.SubTaskContext) (*api.DalCursorIterator, errors.Error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*BitbucketServerTaskData)
	clauses := []dal.Clause{
		dal.Select("bb.branch"),
		dal.From("_tool_bitbucket_server_branches bb"),
		dal.Where(
			`bb.repo_id = ? and bb.connection_id = ?`,
			data.Options.FullName, data.Options.ConnectionId,
		),
	}

	// construct the input iterator
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return nil, err
	}

	return api.NewDalCursorIterator(db, cursor, reflect.TypeOf(BitbucketServerBranchInput{}))
}

func GetPullRequestsIterator(taskCtx plugin.SubTaskContext, collectorWithState *api.ApiCollectorStateManager) (*api.DalCursorIterator, errors.Error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*BitbucketServerTaskData)
	clauses := []dal.Clause{
		dal.Select("bpr.bitbucket_id"),
		dal.From("_tool_bitbucket_server_pull_requests bpr"),
		dal.Where(
			`bpr.repo_id = ? and bpr.connection_id = ?`,
			data.Options.FullName, data.Options.ConnectionId,
		),
	}
	if collectorWithState.IsIncremental && collectorWithState.Since != nil {
		clauses = append(clauses, dal.Where("bitbucket_server_updated_at > ?", *collectorWithState.Since))
	}

	// construct the input iterator
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return nil, err
	}

	return api.NewDalCursorIterator(db, cursor, reflect.TypeOf(BitbucketServerInput{}))
}
