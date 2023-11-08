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
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

var _ plugin.SubTaskEntryPoint = DBGetActionHistory

type actionHistoryHandler struct {
	rawDataParams           string
	changelogBachSave       *api.BatchSave
	changelogDetailBachSave *api.BatchSave
	stories                 map[int64]struct{}
	tasks                   map[int64]struct{}
	bugs                    map[int64]struct{}
}

func newActionHistoryHandler(taskCtx plugin.SubTaskContext, divider *api.BatchSaveDivider) (*actionHistoryHandler, errors.Error) {
	data := taskCtx.GetData().(*ZentaoTaskData)
	changelogBachSave, err := divider.ForType(reflect.TypeOf(&models.ZentaoChangelog{}))
	if err != nil {
		return nil, err
	}
	changelogDetailBachSave, err := divider.ForType(reflect.TypeOf(&models.ZentaoChangelogDetail{}))
	if err != nil {
		return nil, err
	}
	blob, _ := json.Marshal(data.Options.GetParams())
	rawDataParams := string(blob)
	db := taskCtx.GetDal()
	err = db.Delete(&models.ZentaoChangelog{}, dal.Where("_raw_data_params = ?", rawDataParams))
	if err != nil {
		return nil, err
	}
	err = db.Delete(&models.ZentaoChangelogDetail{}, dal.Where("_raw_data_params = ?", rawDataParams))
	if err != nil {
		return nil, err
	}
	return &actionHistoryHandler{
		rawDataParams:           rawDataParams,
		changelogBachSave:       changelogBachSave,
		changelogDetailBachSave: changelogDetailBachSave,
		stories:                 data.Stories,
		tasks:                   data.Tasks,
		bugs:                    data.Bugs,
	}, nil
}

func (h actionHistoryHandler) collectActionHistory(rdb dal.Dal, connectionId uint64) errors.Error {
	clause := []dal.Clause{
		dal.Select("*,zt_action.id aid,zt_history.id hid "),
		dal.From("zt_action"),
		dal.Join("LEFT JOIN zt_history on zt_history.action = zt_action.id"),
		dal.Where("? IN ?", dal.ClauseColumn{Table: "zt_action", Name: "objectType"}, []string{"story", "task", "bug"}),
	}
	cursor, err := rdb.Cursor(clause...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	for cursor.Next() {
		var ah models.ZentaoRemoteDbActionHistory
		err = rdb.Fetch(cursor, &ah)
		if err != nil {
			return err
		}
		switch ah.ObjectType {
		case "story":
			if _, ok := h.stories[ah.ObjectId]; !ok {
				continue
			}
		case "task":
			if _, ok := h.tasks[ah.ObjectId]; !ok {
				continue
			}
		case "bug":
			if _, ok := h.bugs[ah.ObjectId]; !ok {
				continue
			}
		default:
			continue
		}

		zcc := ah.Convert(connectionId)
		zcc.Changelog.NoPKModel.RawDataParams = h.rawDataParams
		zcc.ChangelogDetail.NoPKModel.RawDataParams = h.rawDataParams
		zcc.Changelog.NoPKModel.RawDataTable = "zt_action"
		zcc.ChangelogDetail.NoPKModel.RawDataTable = "zt_history"
		err = h.changelogBachSave.Add(zcc.Changelog)
		if err != nil {
			return err
		}
		if zcc.ChangelogDetail.Id != 0 {
			err = h.changelogDetailBachSave.Add(zcc.ChangelogDetail)
			if err != nil {
				return err
			}
		}
	}
	err = h.changelogBachSave.Flush()
	if err != nil {
		return err
	}
	return h.changelogDetailBachSave.Flush()
}

func DBGetActionHistory(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)

	// skip if no RemoteDb
	if data.RemoteDb == nil {
		return nil
	}

	divider := api.NewBatchSaveDivider(taskCtx, 500, "", "")
	defer func() {
		err1 := divider.Close()
		if err1 != nil {
			panic(err1)
		}
	}()
	handler, err := newActionHistoryHandler(taskCtx, divider)
	if err != nil {
		return err
	}
	return handler.collectActionHistory(data.RemoteDb, data.Options.ConnectionId)
}

var DBGetChangelogMeta = plugin.SubTaskMeta{
	Name:             "collectChangelog",
	EntryPoint:       DBGetActionHistory,
	EnabledByDefault: true,
	Description:      "get action and history data to be changelog from Zentao databases",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}
