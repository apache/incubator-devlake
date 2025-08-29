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
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	dsmodels "github.com/apache/incubator-devlake/helpers/pluginhelper/api/models"
	"github.com/apache/incubator-devlake/plugins/slack/models"
)

type SlackRemotePagination struct {
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit"`
}

type slackConvListResp struct {
	Ok               bool              `json:"ok"`
	Error            string            `json:"error"`
	Needed           string            `json:"needed"`
	Provided         string            `json:"provided"`
	Channels         []json.RawMessage `json:"channels"`
	ResponseMetadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
}

func listSlackRemoteScopes(
	_ *models.SlackConnection,
	apiClient plugin.ApiClient,
	_ string,
	page SlackRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.SlackChannel],
	nextPage *SlackRemotePagination,
	err errors.Error,
) {
	if page.Limit == 0 {
		page.Limit = 100
	}
	// helper to perform API call with given query
	call := func(q url.Values) (*slackConvListResp, errors.Error) {
		res, e := apiClient.Get("conversations.list", q, nil)
		if e != nil {
			return nil, e
		}
		resp := &slackConvListResp{}
		if e = helper.UnmarshalResponse(res, resp); e != nil {
			return nil, e
		}
		return resp, nil
	}

	q := url.Values{}
	q.Set("limit", strconv.Itoa(page.Limit))
	if page.Cursor != "" {
		q.Set("cursor", page.Cursor)
	}

	q.Set("types", "public_channel,private_channel")
	resp, e := call(q)
	if e != nil {
		err = e
		return
	}
	// handle missing_scope gracefully by retrying with private_channel only if channels:read is missing
	if !resp.Ok && resp.Error == "missing_scope" {
		if strings.Contains(resp.Needed, "channels:read") {
			// retry with private channels only (requires groups:read)
			q.Set("types", "private_channel")
			resp, e = call(q)
			if e != nil {
				err = e
				return
			}
		}
	}
	if !resp.Ok {
		err = errors.BadInput.New("slack conversations.list error: " + resp.Error)
		return
	}

	for _, raw := range resp.Channels {
		var ch models.SlackChannel
		if e := errors.Convert(json.Unmarshal(raw, &ch)); e != nil {
			err = e
			return
		}
		children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.SlackChannel]{
			Type:     helper.RAS_ENTRY_TYPE_SCOPE,
			Id:       ch.Id,
			Name:     ch.Name,
			FullName: ch.Name,
			Data:     &ch,
		})
	}
	if resp.ResponseMetadata.NextCursor != "" {
		nextPage = &SlackRemotePagination{Cursor: resp.ResponseMetadata.NextCursor, Limit: page.Limit}
	}
	return
}

func searchSlackRemoteScopes(
	apiClient plugin.ApiClient,
	params *dsmodels.DsRemoteApiScopeSearchParams,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.SlackChannel],
	err errors.Error,
) {
	cursor := ""
	remaining := params.PageSize
	for remaining > 0 {
		list, next, e := listSlackRemoteScopes(nil, apiClient, "", SlackRemotePagination{Cursor: cursor, Limit: 200})
		if e != nil {
			err = e
			return
		}
		for _, it := range list {
			if params.Search == "" || (it.Name != "" && strings.Contains(strings.ToLower(it.Name), strings.ToLower(params.Search))) {
				children = append(children, it)
				remaining--
				if remaining == 0 {
					return
				}
			}
		}
		if next == nil || next.Cursor == "" {
			break
		}
		cursor = next.Cursor
	}
	return
}

func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeList.Get(input)
}

func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeSearch.Get(input)
}
