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
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/plugins/jira/models"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
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
		/*
			This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
			set of data to be process, for example, we process JiraIssues by Board
		*/
		Params: JiraApiParams{
			ConnectionId: data.Options.ConnectionId,
			BoardId:      data.Options.BoardId,
		},
		/*
			Table store raw data
		*/
		Table: RAW_ISSUE_TABLE,
	})
	if err != nil {
		return err
	}

	// build jql
	// IMPORTANT: we have to keep paginated data in a consistence order to avoid data-missing, if we sort issues by
	//  `updated`, issue will be jumping between pages if it got updated during the collection process
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

	err = apiCollector.InitCollector(api.ApiCollectorArgs{
		ApiClient: data.ApiClient,
		PageSize:  data.Options.PageSize,
		/*
			url may use arbitrary variables from different connection in any order, we need GoTemplate to allow more
			flexible for all kinds of possibility.
			Pager contains information for a particular page, calculated by ApiCollector, and will be passed into
			GoTemplate to generate a url for that page.
			We want to do page-fetching in ApiCollector, because the logic are highly similar, by doing so, we can
			avoid duplicate logic for every tasks, and when we have a better idea like improving performance, we can
			do it in one place
		*/
		UrlTemplate: "agile/1.0/board/{{ .Params.BoardId }}/issue",
		/*
			(Optional) Return query string for request, or you can plug them into UrlTemplate directly
		*/
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("jql", jql)
			query.Set("startAt", fmt.Sprintf("%v", reqData.Pager.Skip))
			query.Set("maxResults", fmt.Sprintf("%v", reqData.Pager.Size))
			query.Set("expand", "changelog")
			return query, nil
		},
		/*
			Some api might do pagination by http headers
		*/
		//Header: func(pager *plugin.Pager) http.Header {
		//},
		/*
			Sometimes, we need to collect data based on previous collected data, like jira changelog, it requires
			issue_id as part of the url.
			We can mimic `stdin` design, to accept a `Input` function which produces a `Iterator`, collector
			should iterate all records, and do data-fetching for each on, either in parallel or sequential order
			UrlTemplate: "api/3/issue/{{ Input.ID }}/changelog"
		*/
		//Input: databaseIssuesIterator,
		/*
			For api endpoint that returns number of total pages, ApiCollector can collect pages in parallel with ease,
			or other techniques are required if this information was missing.
		*/
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
	if err != nil {
		return err
	}

	return apiCollector.Execute()
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
	if data.JiraServerInfo.DeploymentType == models.DeploymentServer {
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
