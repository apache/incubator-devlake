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
	"fmt"
	"reflect"
	"strings"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"

	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

const RAW_EPIC_TABLE = "jira_api_epics"

var _ plugin.SubTaskEntryPoint = CollectEpics

var CollectEpicsMeta = plugin.SubTaskMeta{
	Name:             "collectEpics",
	EntryPoint:       CollectEpics,
	EnabledByDefault: true,
	Description:      "collect Jira epics from all boards, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET, plugin.DOMAIN_TYPE_CROSS},
}

func CollectEpics(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*JiraTaskData)
	logger := taskCtx.GetLogger()
	batchSize := 100
	if strings.EqualFold(string(data.JiraServerInfo.DeploymentType), string(models.DeploymentServer)) && len(data.JiraServerInfo.VersionNumbers) == 3 && data.JiraServerInfo.VersionNumbers[0] <= 8 {
		batchSize = 1
	}
	epicIterator, err := GetEpicKeysIterator(db, data, batchSize)
	if err != nil {
		return err
	}

	apiCollector, err := api.NewStatefulApiCollector(api.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: JiraApiParams{
			ConnectionId: data.Options.ConnectionId,
			BoardId:      data.Options.BoardId,
		},
		Table: RAW_EPIC_TABLE,
	})
	if err != nil {
		return err
	}

	loc, err := getTimeZone(taskCtx)
	if err != nil {
		logger.Info("failed to get timezone, err: %v", err)
	} else {
		logger.Info("got user's timezone: %v", loc.String())
	}
	jql := "ORDER BY created ASC"
	if apiCollector.GetSince() != nil {
		jql = buildJQL(*apiCollector.GetSince(), loc)
	}

	// Choose API endpoint based on JIRA deployment type
	if strings.EqualFold(string(data.JiraServerInfo.DeploymentType), string(models.DeploymentServer)) {
		logger.Info("Using api/2/search for JIRA Server")
		err = setupApiV2Collector(apiCollector, data, epicIterator, jql)
	} else {
		logger.Info("Using api/3/search/jql for JIRA Cloud")
		err = setupApiV3Collector(apiCollector, data, epicIterator, jql)
	}
	if err != nil {
		return err
	}
	return apiCollector.Execute()
}

// JIRA Server API v2 collector
func setupApiV2Collector(apiCollector *api.StatefulApiCollector, data *JiraTaskData, epicIterator api.Iterator, jql string) errors.Error {
	return apiCollector.InitCollector(api.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Incremental: false,
		UrlTemplate: "api/2/search",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			epicKeys := []string{}

			input, ok := reqData.Input.([]interface{})
			if !ok {
				return nil, errors.Default.New("invalid input type, expected []interface{}")
			}

			for _, e := range input {
				if epicKey, ok := e.(*string); ok && epicKey != nil {
					epicKeys = append(epicKeys, *epicKey)
				}
			}

			localJQL := fmt.Sprintf("issue in (%s) and %s", strings.Join(epicKeys, ","), jql)
			query.Set("jql", localJQL)
			query.Set("startAt", fmt.Sprintf("%v", reqData.Pager.Skip))
			query.Set("maxResults", fmt.Sprintf("%v", reqData.Pager.Size))
			query.Set("expand", "changelog")
			return query, nil
		},
		Input:         epicIterator,
		GetTotalPages: GetTotalPagesFromResponse,
		Concurrency:   10,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data struct {
				Issues []json.RawMessage `json:"issues"`
			}
			blob, err := io.ReadAll(res.Body)
			if err != nil {
				return nil, errors.Convert(err)
			}
			err = json.Unmarshal(blob, &data)
			if err != nil {
				return nil, errors.Convert(err)
			}
			return data.Issues, nil
		},
		AfterResponse: ignoreHTTPStatus400,
	})
}

// JIRA Cloud API v3 collector
func setupApiV3Collector(apiCollector *api.StatefulApiCollector, data *JiraTaskData, epicIterator api.Iterator, jql string) errors.Error {
	return apiCollector.InitCollector(api.ApiCollectorArgs{
		ApiClient:             data.ApiClient,
		PageSize:              100,
		Incremental:           false,
		UrlTemplate:           "api/3/search/jql",
		GetNextPageCustomData: getNextPageCustomDataForV3,
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			epicKeys := []string{}
			for _, e := range reqData.Input.([]interface{}) {
				epicKeys = append(epicKeys, *e.(*string))
			}
			localJQL := fmt.Sprintf("issue in (%s) and %s", strings.Join(epicKeys, ","), jql)
			query.Set("jql", localJQL)
			query.Set("maxResults", fmt.Sprintf("%v", reqData.Pager.Size))
			query.Set("expand", "changelog")
			query.Set("fields", "*all")

			if reqData.CustomData != nil {
				query.Set("nextPageToken", reqData.CustomData.(string))
			}

			return query, nil
		},
		Input: epicIterator,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data struct {
				Issues []json.RawMessage `json:"issues"`
			}
			blob, err := io.ReadAll(res.Body)
			if err != nil {
				return nil, errors.Convert(err)
			}
			err = json.Unmarshal(blob, &data)
			if err != nil {
				return nil, errors.Convert(err)
			}
			return data.Issues, nil
		},
		AfterResponse: ignoreHTTPStatus400,
	})
}

// Get next page token for API v3
func getNextPageCustomDataForV3(_ *api.RequestData, prevPageResponse *http.Response) (interface{}, errors.Error) {
	var response struct {
		NextPageToken string `json:"nextPageToken"`
	}

	blob, err := io.ReadAll(prevPageResponse.Body)
	if err != nil {
		return nil, errors.Convert(err)
	}

	prevPageResponse.Body = io.NopCloser(strings.NewReader(string(blob)))

	err = json.Unmarshal(blob, &response)
	if err != nil {
		return nil, errors.Convert(err)
	}

	if response.NextPageToken == "" {
		return nil, api.ErrFinishCollect
	}

	return response.NextPageToken, nil
}

func GetEpicKeysIterator(db dal.Dal, data *JiraTaskData, batchSize int) (api.Iterator, errors.Error) {
	cursor, err := db.Cursor(
		dal.Select("DISTINCT epic_key"),
		dal.From("_tool_jira_issues i"),
		dal.Join(`
			LEFT JOIN _tool_jira_board_issues bi ON (
			i.connection_id = bi.connection_id
			AND
			i.issue_id = bi.issue_id
		)`),
		dal.Where(`
			i.connection_id = ?
			AND
			bi.board_id = ?
			AND
			i.epic_key != ''
		`, data.Options.ConnectionId, data.Options.BoardId,
		),
	)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to query for external epics")
	}
	iter, err := api.NewBatchedDalCursorIterator(db, cursor, reflect.TypeOf(""), batchSize)
	if err != nil {
		return nil, err
	}
	return iter, nil
}
