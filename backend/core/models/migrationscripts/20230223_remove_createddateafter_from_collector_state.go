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

package migrationscripts

import (
	"encoding/json"
	"fmt"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
)

type blueprint20230223 struct {
	ID       int64
	Settings json.RawMessage `gorm:"serializer:encdec"`
}

func (blueprint20230223) TableName() string {
	return "_devlake_blueprints"
}

type removeCreatedDateAfterFromCollectorMeta20230223 struct{}

func (script *removeCreatedDateAfterFromCollectorMeta20230223) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	// step 1: rename bp.settings.createdDateAfter to timeAfter
	bp := &blueprint20230223{}
	cursor, err := db.Cursor(dal.From(bp), dal.Where("mode = ?", "NORMAL"))
	if err != nil {
		return err
	}
	defer cursor.Close()
	for cursor.Next() {
		err = db.Fetch(cursor, bp)
		if err != nil {
			return err
		}
		settingsMap := make(map[string]interface{})
		if e := json.Unmarshal(bp.Settings, &settingsMap); e != nil {
			return errors.Default.Wrap(e, fmt.Sprintf("failed to unmarshal settings for blueprint #%v", bp.ID))
		}
		if v, ok := settingsMap["createdDateAfter"]; ok {
			settingsMap["timeAfter"] = v
			delete(settingsMap, "createdDateAfter")
		} else {
			continue
		}
		if s, e := json.Marshal(settingsMap); e == nil {
			bp.Settings = s
			err = db.Update(bp)
			if err != nil {
				return err
			}
		} else {
			return errors.Default.Wrap(e, fmt.Sprintf("failed to update settings for blueprint #%v", bp.ID))
		}
	}

	// step 2: update collector_latest_state.time_after with values from created_date_after
	table := "_devlake_collector_latest_state"
	err = db.UpdateColumn(
		table,
		"time_after", dal.Expr("created_date_after"),
		dal.Where("time_after IS NULL"),
	)
	if err != nil {
		return err
	}

	// step 3: drop collector_latest_state.created_date_after
	return db.DropColumns(table, "created_date_after")
}

func (*removeCreatedDateAfterFromCollectorMeta20230223) Version() uint64 {
	return 20230223200040
}

func (*removeCreatedDateAfterFromCollectorMeta20230223) Name() string {
	return "remove created_date_after from _devlake_collector_latest_state"
}
