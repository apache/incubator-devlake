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
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"strings"

	"encoding/json"
	"github.com/apache/incubator-devlake/plugins/helper"
	"io/ioutil"
	"net/http"
	"net/url"
)

const RAW_EPIC_TABLE = "jira_api_epics"

type JiraEpicParams struct {
	ConnectionId uint64
	BoardId      uint64
}

var _ core.SubTaskEntryPoint = CollectIssues

var CollectEpicsMeta = core.SubTaskMeta{
	Name:             "collectEpics",
	EntryPoint:       CollectEpics,
	EnabledByDefault: true,
	Description:      "collect Jira epics from all boards",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET, core.DOMAIN_TYPE_CROSS},
}

func CollectEpics(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*JiraTaskData)
	externalEpicKeys, err := GetEpicKeys(db, data)
	if err != nil {
		return err
	}
	if len(externalEpicKeys) == 0 {
		taskCtx.GetLogger().Info("no external epic keys found for Jira board %d", data.Options.BoardId)
		return nil
	}
	since := data.Since
	jql := "ORDER BY created ASC"
	if since != nil {
		// prepend a time range criteria if `since` was specified, either by user or from database
		jql = fmt.Sprintf("updated >= '%s' %s", since.Format("2006/01/02 15:04"), jql)
	}
	jql = fmt.Sprintf("issue in (%s) %s", strings.Join(externalEpicKeys, ","), jql)

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraEpicParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_EPIC_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Incremental: false,
		UrlTemplate: "api/2/search",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("jql", jql)
			query.Set("issue in", fmt.Sprintf("(%s)", strings.Join(externalEpicKeys, ",")))
			query.Set("startAt", fmt.Sprintf("%v", reqData.Pager.Skip))
			query.Set("maxResults", fmt.Sprintf("%v", reqData.Pager.Size))
			query.Set("expand", "changelog")
			return query, nil
		},
		GetTotalPages: GetTotalPagesFromResponse,
		Concurrency:   10,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var data struct {
				Issues []json.RawMessage `json:"issues"`
			}
			blob, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(blob, &data)
			if err != nil {
				return nil, err
			}
			return data.Issues, nil
		},
	})
	if err != nil {
		return err
	}
	return collector.Execute()
}

func GetEpicKeys(db dal.Dal, data *JiraTaskData) ([]string, error) {
	cursor, err := db.RawCursor(`
		select
			distinct epic_key
		from
			_tool_jira_issues i
		left join _tool_jira_board_issues bi on (
			i.connection_id = bi.connection_id
			and 
			i.issue_id = bi.issue_id
		)
		where
		i.connection_id = ?
		and 
		bi.board_id = ?
		and
		i.epic_key != ''
		`, data.Options.ConnectionId, data.Options.BoardId)
	if err != nil {
		return nil, fmt.Errorf("unable to query for external epics: %v", err)
	}
	var externalEpicKeys []string
	for cursor.Next() {
		epicKey := ""
		err = cursor.Scan(&epicKey)
		if err != nil {
			return nil, fmt.Errorf("couldn't read returned epic key: %v", err)
		}
		externalEpicKeys = append(externalEpicKeys, epicKey)
	}
	return externalEpicKeys, nil
}
