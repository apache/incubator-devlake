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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

var CollectBoardFilterBeginMeta = plugin.SubTaskMeta{
	Name:             "collectBoardFilterBegin",
	EntryPoint:       CollectBoardFilterBegin,
	EnabledByDefault: true,
	Description:      "Jira board filter jql checker before running",
	DomainTypes:      plugin.DOMAIN_TYPES,
}

func CollectBoardFilterBegin(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDal()
	logger.Info("collect board in collectBoardFilterBegin: %d", data.Options.BoardId)
	// get board filter id
	filterId, err := getBoardFilterId(data)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error getting board filter id for connection_id:%d board_id:%d", data.Options.ConnectionId, data.Options.BoardId))
	}
	logger.Info("collect board filter:%s", filterId)

	// get board filter jql
	filterInfo, err := getBoardFilterJql(data, filterId)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error getting board filter jql for connection_id:%d board_id:%d", data.Options.ConnectionId, data.Options.BoardId))
	}
	logger.Info("collect board filter jql:%s", filterInfo.Jql)

	jql := filterInfo.Jql
	var record models.JiraBoard
	err = db.First(&record, dal.Where("connection_id = ? AND board_id = ? ", data.Options.ConnectionId, data.Options.BoardId))
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error finding record in _tool_jira_boards table for connection_id:%d board_id:%d", data.Options.ConnectionId, data.Options.BoardId))
	}

	// full sync
	syncPolicy := taskCtx.TaskContext().SyncPolicy()
	if syncPolicy != nil && syncPolicy.FullSync {
		if record.Jql != jql {
			record.Jql = jql
			err = db.Update(&record, dal.Where("connection_id = ? AND board_id = ? ", data.Options.ConnectionId, data.Options.BoardId))
			if err != nil {
				return errors.Default.Wrap(err, fmt.Sprintf("error updating record in _tool_jira_boards table for connection_id:%d board_id:%d", data.Options.ConnectionId, data.Options.BoardId))
			}
			logger.Info("full sync mode, update jql to %s", record.Jql)
		}
		return nil
	}

	// first run
	if record.Jql == "" {
		record.Jql = jql
		err = db.Update(&record, dal.Where("connection_id = ? AND board_id = ? ", data.Options.ConnectionId, data.Options.BoardId))
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("error updating record in _tool_jira_boards table for connection_id:%d board_id:%d", data.Options.ConnectionId, data.Options.BoardId))
		}
		logger.Info("first run, update jql to %s", record.Jql)
		return nil
	}
	// change
	if record.Jql != jql {
		return errors.Default.New(fmt.Sprintf("connection_id:%d board_id:%d filter jql has changed, please use fullSync mode. And the previous jql is %s, now jql is %s", data.Options.ConnectionId, data.Options.BoardId, record.Jql, jql))
	}
	// no change
	return nil
}

func getBoardFilterId(data *JiraTaskData) (string, error) {
	url := fmt.Sprintf("agile/1.0/board/%d/configuration", data.Options.BoardId)
	boardConfiguration, err := data.ApiClient.Get(url, nil, nil)
	if err != nil {
		return "", err
	}
	bc := &BoardConfiguration{}
	err = helper.UnmarshalResponse(boardConfiguration, bc)
	if err != nil {
		return "", err
	}
	filterId := bc.Filter.ID
	return filterId, nil
}

func getBoardFilterJql(data *JiraTaskData, filterId string) (*FilterInfo, error) {
	url := fmt.Sprintf("api/2/filter/%s", filterId)
	filterInfo, err := data.ApiClient.Get(url, nil, nil)
	if err != nil {
		return nil, err
	}
	fi := &FilterInfo{}
	err = helper.UnmarshalResponse(filterInfo, fi)
	if err != nil {
		return nil, err
	}
	return fi, nil
}

type BoardConfiguration struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Self     string `json:"self"`
	Location struct {
		Type string `json:"type"`
		Key  string `json:"key"`
		ID   string `json:"id"`
		Self string `json:"self"`
		Name string `json:"name"`
	} `json:"location"`
	Filter struct {
		ID   string `json:"id"`
		Self string `json:"self"`
	} `json:"filter"`
	ColumnConfig struct {
		Columns []struct {
			Name     string `json:"name"`
			Statuses []struct {
				ID   string `json:"id"`
				Self string `json:"self"`
			} `json:"statuses"`
		} `json:"columns"`
		ConstraintType string `json:"constraintType"`
	} `json:"columnConfig"`
	Estimation struct {
		Type  string `json:"type"`
		Field struct {
			FieldID     string `json:"fieldId"`
			DisplayName string `json:"displayName"`
		} `json:"field"`
	} `json:"estimation"`
	Ranking struct {
		RankCustomFieldID int `json:"rankCustomFieldId"`
	} `json:"ranking"`
}

type FilterInfo struct {
	Self        string `json:"self"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Owner       struct {
		Self       string `json:"self"`
		AccountID  string `json:"accountId"`
		AvatarUrls struct {
			Four8X48  string `json:"48x48"`
			Two4X24   string `json:"24x24"`
			One6X16   string `json:"16x16"`
			Three2X32 string `json:"32x32"`
		} `json:"avatarUrls"`
		DisplayName string `json:"displayName"`
		Active      bool   `json:"active"`
	} `json:"owner"`
	Jql              string `json:"jql"`
	ViewURL          string `json:"viewUrl"`
	SearchURL        string `json:"searchUrl"`
	Favourite        bool   `json:"favourite"`
	FavouritedCount  int    `json:"favouritedCount"`
	SharePermissions []struct {
		ID      int    `json:"id"`
		Type    string `json:"type"`
		Project struct {
			Self         string `json:"self"`
			ID           string `json:"id"`
			Key          string `json:"key"`
			AssigneeType string `json:"assigneeType"`
			Name         string `json:"name"`
			Roles        struct {
			} `json:"roles"`
			AvatarUrls struct {
				Four8X48  string `json:"48x48"`
				Two4X24   string `json:"24x24"`
				One6X16   string `json:"16x16"`
				Three2X32 string `json:"32x32"`
			} `json:"avatarUrls"`
			ProjectTypeKey string `json:"projectTypeKey"`
			Simplified     bool   `json:"simplified"`
			Style          string `json:"style"`
			Properties     struct {
			} `json:"properties"`
		} `json:"project"`
	} `json:"sharePermissions"`
	EditPermissions []any `json:"editPermissions"`
	IsWritable      bool  `json:"isWritable"`
	SharedUsers     struct {
		Size       int   `json:"size"`
		Items      []any `json:"items"`
		MaxResults int   `json:"max-results"`
		StartIndex int   `json:"start-index"`
		EndIndex   int   `json:"end-index"`
	} `json:"sharedUsers"`
	Subscriptions struct {
		Size       int   `json:"size"`
		Items      []any `json:"items"`
		MaxResults int   `json:"max-results"`
		StartIndex int   `json:"start-index"`
		EndIndex   int   `json:"end-index"`
	} `json:"subscriptions"`
}
