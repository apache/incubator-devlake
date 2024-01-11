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

var CollectBoardFilterEndMeta = plugin.SubTaskMeta{
	Name:             "collectBoardFilterEnd",
	EntryPoint:       CollectBoardFilterEnd,
	EnabledByDefault: true,
	Description:      "Jira board filter jql checker after runnig",
	DomainTypes:      plugin.DOMAIN_TYPES,
}

func CollectBoardFilterEnd(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDal()
	logger.Info("collect board in collectBoardFilterEnd: %d", data.Options.BoardId)

	// get board filter id
	url := fmt.Sprintf("agile/1.0/board/%d/configuration", data.Options.BoardId)
	boardConfiguration, err := data.ApiClient.Get(url, nil, nil)
	if err != nil {
		return err
	}
	bc := &BoardConfiguration{}
	err = helper.UnmarshalResponse(boardConfiguration, bc)
	if err != nil {
		return err
	}
	filterId := bc.Filter.ID
	logger.Info("collect board filter:%s", filterId)

	// get board filter jql
	url = fmt.Sprintf("api/2/filter/%s", filterId)
	filterInfo, err := data.ApiClient.Get(url, nil, nil)
	if err != nil {
		return err
	}
	fi := &FilterInfo{}
	err = helper.UnmarshalResponse(filterInfo, fi)
	if err != nil {
		return err
	}
	jql := fi.Jql
	logger.Info("collect board filter jql:%s", jql)

	// should not change
	var record models.JiraBoard
	err = db.First(&record, dal.Where("connection_id = ? AND board_id = ? ", data.Options.ConnectionId, data.Options.BoardId))
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error finding record in _tool_jira_boards table for connection_id:%d board_id:%d", data.Options.ConnectionId, data.Options.BoardId))
	}

	if record.Jql != jql {
		return errors.Default.New(fmt.Sprintf("board filter jql has changed for connection_id:%d board_id:%d, please use fullSync mode!!!", data.Options.ConnectionId, data.Options.BoardId))
	}

	return nil
}
