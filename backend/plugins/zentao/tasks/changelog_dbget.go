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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

var _ plugin.SubTaskEntryPoint = DBGetActionHistory

func DBGetActionHistory(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)

	// skip if no RemoteDb
	if data.Options.RemoteDb == nil {
		return nil
	}

	divider := api.NewBatchSaveDivider(taskCtx, 500, "", "")
	defer func() {
		err1 := divider.Close()
		if err1 != nil {
			panic(err1)
		}
	}()

	return dBGetActionHistory(data.Options, func(zcc *models.ZentaoChangelogCom) errors.Error {
		batch, err := divider.ForType(reflect.TypeOf(zcc.Changelog))
		if err != nil {
			return err
		}
		zcc.Changelog.ConnectionId = data.Options.ConnectionId
		batch.Add(zcc.Changelog)

		if zcc.ChangelogDetail.Id != 0 {
			batch, err = divider.ForType(reflect.TypeOf(zcc.ChangelogDetail))
			if err != nil {
				return err
			}
			zcc.ChangelogDetail.ConnectionId = data.Options.ConnectionId
			batch.Add(zcc.ChangelogDetail)
		}
		return nil
	})
}

var DBGetChangelogMeta = plugin.SubTaskMeta{
	Name:             "DBGetChangelog",
	EntryPoint:       DBGetActionHistory,
	EnabledByDefault: true,
	Description:      "get action and history data to be changelog from Zentao databases",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

// it is work for zentao version 18.3
func dBGetActionHistory(op *ZentaoOptions, callback func(*models.ZentaoChangelogCom) errors.Error) errors.Error {
	rdb := op.RemoteDb
	atn := (models.ZentaoRemoteDbAction{}).TableName()
	htn := (models.ZentaoRemoteDbHistory{}).TableName()

	clause := []dal.Clause{
		dal.Select(fmt.Sprintf("*,%s.id aid,%s.id hid ", atn, htn)),
		dal.From(atn),
	}

	if op.ProductId != 0 {
		clause = append(clause, dal.Where(fmt.Sprintf("%s.product = ?", atn), fmt.Sprintf(",%d,", op.ProductId)))
	}
	if op.ProjectId != 0 {
		clause = append(clause, dal.Where(fmt.Sprintf("%s.project = ?", atn), op.ProjectId))
	}
	clause = append(clause, dal.Join(fmt.Sprintf("LEFT JOIN %s on %s.action = %s.id", htn, htn, atn)))

	cursor, err := rdb.Cursor(clause...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	for cursor.Next() {
		actionHistory := &models.ZentaoRemoteDbActionHistory{}

		err = rdb.Fetch(cursor, actionHistory)
		if err != nil {
			return err
		}

		err = callback(actionHistory.Convert())
		if err != nil {
			return err
		}
	}

	return nil
}
