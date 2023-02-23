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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

type BitbucketApiParams struct {
	ConnectionId uint64
	FullName     string
}

type BitbucketInput struct {
	BitbucketId int
}

type BitbucketUuidInput struct {
	BitbucketId string
}

type BitbucketPagination struct {
	Values  []interface{} `json:"values"`
	PageLen int           `json:"pagelen"`
	Size    int           `json:"size"`
	Page    int           `json:"page"`
	Next    string        `json:"next"`
}

func CreateRawDataSubTaskArgs(taskCtx plugin.SubTaskContext, Table string) (*api.RawDataSubTaskArgs, *BitbucketTaskData) {
	data := taskCtx.GetData().(*BitbucketTaskData)
	RawDataSubTaskArgs := &api.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: BitbucketApiParams{
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

func GetQuery(reqData *api.RequestData) (url.Values, errors.Error) {
	query := url.Values{}
	query.Set("state", "all")
	query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
	query.Set("pagelen", fmt.Sprintf("%v", reqData.Pager.Size))

	return query, nil
}

// GetQueryCreatedAndUpdated is a common GeyQuery for timeFilter and incremental
func GetQueryCreatedAndUpdated(fields string, collectorWithState *api.ApiCollectorStateManager) func(reqData *api.RequestData) (url.Values, errors.Error) {
	return func(reqData *api.RequestData) (url.Values, errors.Error) {
		query, err := GetQuery(reqData)
		if err != nil {
			return nil, err
		}
		query.Set("fields", fields)
		query.Set("sort", "created_on")
		if collectorWithState.IsIncremental() {
			latestSuccessStart := collectorWithState.LatestState.LatestSuccessStart.Format(time.RFC3339)
			query.Set("q", fmt.Sprintf("updated_on>=%s", latestSuccessStart))
		} else if collectorWithState.TimeAfter != nil {
			timeAfter := collectorWithState.TimeAfter.Format(time.RFC3339)
			query.Set("q", fmt.Sprintf("updated_on>=%s", timeAfter))
		}

		return query, nil
	}
}

func GetQueryFields(fields string) func(reqData *api.RequestData) (url.Values, errors.Error) {
	return func(reqData *api.RequestData) (url.Values, errors.Error) {
		query, err := GetQuery(reqData)
		if err != nil {
			return nil, err
		}
		query.Set("fields", fields)

		return query, nil
	}
}

func GetNextPageCustomData(_ *api.RequestData, prevPageResponse *http.Response) (interface{}, errors.Error) {
	var rawMessages struct {
		Next string `json:"next"`
	}
	err := decodeResponse(prevPageResponse, &rawMessages)
	if err != nil {
		return nil, err
	}
	if rawMessages.Next == `` {
		return ``, api.ErrFinishCollect
	}
	u, err := errors.Convert01(url.Parse(rawMessages.Next))
	if err != nil {
		return nil, err
	}
	return u.Query()[`page`][0], nil
}

func GetTotalPagesFromResponse(res *http.Response, args *api.ApiCollectorArgs) (int, errors.Error) {
	body := &BitbucketPagination{}
	err := api.UnmarshalResponse(res, body)
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
	err := decodeResponse(res, &rawMessages)
	if err != nil {
		return nil, err
	}

	return rawMessages.Values, nil
}

func GetPullRequestsIterator(taskCtx plugin.SubTaskContext, collectorWithState *api.ApiCollectorStateManager) (*api.DalCursorIterator, errors.Error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*BitbucketTaskData)
	clauses := []dal.Clause{
		dal.Select("bpr.bitbucket_id"),
		dal.From("_tool_bitbucket_pull_requests bpr"),
		dal.Where(
			`bpr.repo_id = ? and bpr.connection_id = ?`,
			data.Options.FullName, data.Options.ConnectionId,
		),
	}
	if collectorWithState.IsIncremental() {
		clauses = append(clauses, dal.Where("bitbucket_updated_at > ?", *collectorWithState.LatestState.LatestSuccessStart))
	}
	// construct the input iterator
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return nil, err
	}

	return api.NewDalCursorIterator(db, cursor, reflect.TypeOf(BitbucketInput{}))
}

func GetIssuesIterator(taskCtx plugin.SubTaskContext, collectorWithState *api.ApiCollectorStateManager) (*api.DalCursorIterator, errors.Error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*BitbucketTaskData)
	clauses := []dal.Clause{
		dal.Select("bpr.bitbucket_id"),
		dal.From("_tool_bitbucket_issues bpr"),
		dal.Where(
			`bpr.repo_id = ? and bpr.connection_id = ?`,
			data.Options.FullName, data.Options.ConnectionId,
		),
	}
	if collectorWithState.IsIncremental() {
		clauses = append(clauses, dal.Where("bitbucket_updated_at > ?", *collectorWithState.LatestState.LatestSuccessStart))
	}
	// construct the input iterator
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return nil, err
	}

	return api.NewDalCursorIterator(db, cursor, reflect.TypeOf(BitbucketInput{}))
}

func GetPipelinesIterator(taskCtx plugin.SubTaskContext, collectorWithState *api.ApiCollectorStateManager) (*api.DalCursorIterator, errors.Error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*BitbucketTaskData)
	clauses := []dal.Clause{
		dal.Select("bpr.bitbucket_id"),
		dal.From("_tool_bitbucket_pipelines bpr"),
		dal.Where(
			`bpr.repo_id = ? and bpr.connection_id = ?`,
			data.Options.FullName, data.Options.ConnectionId,
		),
	}
	if collectorWithState.IsIncremental() {
		clauses = append(clauses, dal.Where("bitbucket_complete_on > ?", *collectorWithState.LatestState.LatestSuccessStart))
	}
	// construct the input iterator
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return nil, err
	}

	return api.NewDalCursorIterator(db, cursor, reflect.TypeOf(BitbucketUuidInput{}))
}

func ignoreHTTPStatus404(res *http.Response) errors.Error {
	if res.StatusCode == http.StatusUnauthorized {
		return errors.Unauthorized.New("authentication failed, please check your AccessToken")
	}
	if res.StatusCode == http.StatusNotFound {
		return api.ErrIgnoreAndContinue
	}
	return nil
}
