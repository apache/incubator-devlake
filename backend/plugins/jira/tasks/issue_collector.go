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
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

const RAW_ISSUE_TABLE = "jira_api_issues"

var _ plugin.SubTaskEntryPoint = CollectIssues

var CollectIssuesMeta = plugin.SubTaskMeta{
	Name:             "collectIssues",
	EntryPoint:       CollectIssues,
	EnabledByDefault: true,
	Description:      "collect Jira issues, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET, plugin.DOMAIN_TYPE_CROSS},
}

func CollectIssues(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	logger := taskCtx.GetLogger()
	apiCollector, err := api.NewStatefulApiCollector(api.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: JiraApiParams{
			ConnectionId: data.Options.ConnectionId,
			BoardId:      data.Options.BoardId,
		},
		Table: RAW_ISSUE_TABLE,
	})
	if err != nil {
		return err
	}

	// IMPORTANT: we sort by `created ASC` to keep paginated data in a consistent order.
	// Sorting by `updated` would cause issues to jump between pages during collection.
	loc, err := getTimeZone(taskCtx)
	if err != nil {
		logger.Info("failed to get timezone, err: %v", err)
	} else {
		logger.Info("got user's timezone: %v", loc.String())
	}
	incrementalJql := "ORDER BY created ASC"
	if apiCollector.GetSince() != nil {
		incrementalJql = buildJQL(*apiCollector.GetSince(), loc)
	}

	// Use the search API with `filter = {id}` JQL instead of the board Agile API.
	// The board Agile API applies kanban sub-filters server-side, which silently
	// excludes resolved issues (e.g. those with a released fixVersion).
	// The search API with the saved filter JQL returns all matching issues.
	filterJql := buildFilterJQL(data.FilterId, incrementalJql)
	logger.Info("collecting issues via search API with JQL: %s", filterJql)

	pageSize := data.Options.PageSize
	if pageSize == 0 {
		pageSize = 100
	}

	if strings.EqualFold(string(data.JiraServerInfo.DeploymentType), string(models.DeploymentServer)) {
		logger.Info("Using api/2/search for JIRA Server issue collection")
		err = setupIssueV2Collector(apiCollector, data, filterJql, pageSize)
	} else {
		logger.Info("Using api/3/search/jql for JIRA Cloud issue collection")
		err = setupIssueV3Collector(apiCollector, data, filterJql, pageSize)
	}
	if err != nil {
		return err
	}

	return apiCollector.Execute()
}

func buildFilterJQL(filterId string, incrementalJql string) string {
	if filterId == "" {
		return incrementalJql
	}
	// Use Jira's `filter = {id}` syntax to reference the saved filter.
	// This avoids parenthesization bugs when composing raw JQL strings
	// that may contain OR/AND operators.
	if incrementalJql == "ORDER BY created ASC" {
		return fmt.Sprintf("filter = %s ORDER BY created ASC", filterId)
	}
	// incrementalJql contains "updated >= '...' ORDER BY created ASC"
	// We need to insert the filter reference before the incremental clause
	return fmt.Sprintf("filter = %s AND %s", filterId, incrementalJql)
}

func setupIssueV2Collector(apiCollector *api.StatefulApiCollector, data *JiraTaskData, filterJql string, pageSize int) errors.Error {
	return apiCollector.InitCollector(api.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		PageSize:    pageSize,
		UrlTemplate: "api/2/search",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("jql", filterJql)
			query.Set("startAt", fmt.Sprintf("%v", reqData.Pager.Skip))
			query.Set("maxResults", fmt.Sprintf("%v", reqData.Pager.Size))
			query.Set("expand", "changelog")
			return query, nil
		},
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
	})
}

func setupIssueV3Collector(apiCollector *api.StatefulApiCollector, data *JiraTaskData, filterJql string, pageSize int) errors.Error {
	return apiCollector.InitCollector(api.ApiCollectorArgs{
		ApiClient:             data.ApiClient,
		PageSize:              pageSize,
		UrlTemplate:           "api/3/search/jql",
		GetNextPageCustomData: getNextPageCustomDataForV3,
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("jql", filterJql)
			query.Set("maxResults", fmt.Sprintf("%v", reqData.Pager.Size))
			query.Set("expand", "changelog")
			query.Set("fields", "*all")
			if reqData.CustomData != nil {
				query.Set("nextPageToken", reqData.CustomData.(string))
			}
			return query, nil
		},
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
	})
}

// buildJQL build jql based on timeAfter and incremental mode
func buildJQL(since time.Time, location *time.Location) string {
	jql := "ORDER BY created ASC"
	if !since.IsZero() {
		if location != nil {
			since = since.In(location)
		} else {
			since = since.In(time.UTC).Add(-24 * time.Hour)
		}
		jql = fmt.Sprintf("updated >= '%s' %s", since.Format("2006/01/02 15:04"), jql)
	}
	return jql
}

// getTimeZone get user's timezone from jira API
func getTimeZone(taskCtx plugin.SubTaskContext) (*time.Location, errors.Error) {
	data := taskCtx.GetData().(*JiraTaskData)
	connectionId := data.Options.ConnectionId
	var conn models.JiraConnection
	err := taskCtx.GetDal().First(&conn, dal.Where("id = ?", connectionId))
	if err != nil {
		return nil, err
	}
	var resp *http.Response
	var path string
	var query url.Values
	if strings.EqualFold(string(data.JiraServerInfo.DeploymentType), string(models.DeploymentServer)) {
		path = "api/2/user"
		query = url.Values{"username": []string{conn.Username}}
	} else {
		path = "api/3/user"
		var accountId string
		accountId, err = getAccountId(data.ApiClient, conn.Username)
		if err != nil {
			return nil, err
		}
		query = url.Values{"accountId": []string{accountId}}
	}
	resp, err = data.ApiClient.Get(path, query, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var timeZone struct {
		TimeZone string `json:"timeZone"`
	}
	err = errors.Convert(json.NewDecoder(resp.Body).Decode(&timeZone))
	if err != nil {
		return nil, err
	}
	tz, err := errors.Convert01(time.LoadLocation(timeZone.TimeZone))
	if err != nil {
		return nil, err
	}
	if tz == nil {
		return nil, errors.Default.New(fmt.Sprintf("invalid time zone: %s", timeZone.TimeZone))
	}
	return tz, nil
}

func getAccountId(client *api.ApiAsyncClient, username string) (string, errors.Error) {
	resp, err := client.Get("api/3/user/picker", url.Values{"query": []string{username}}, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var accounts struct {
		Users []struct {
			AccountID   string `json:"accountId"`
			AccountType string `json:"accountType"`
			HTML        string `json:"html"`
			DisplayName string `json:"displayName"`
		} `json:"users"`
		Total  int    `json:"total"`
		Header string `json:"header"`
	}
	err = errors.Convert(json.NewDecoder(resp.Body).Decode(&accounts))
	if err != nil {
		return "", err
	}
	if len(accounts.Users) == 0 {
		return "", errors.Default.New("no user found")
	}
	return accounts.Users[0].AccountID, nil
}
