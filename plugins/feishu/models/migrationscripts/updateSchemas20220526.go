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
	"context"
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/feishu/models/migrationscripts/archived"
	"gorm.io/gorm/clause"
	"time"

	"gorm.io/gorm"
)

type FeishuMeetingTopUserItem20220524Temp struct {
	common.NoPKModel `json:"-"`
	StartTime        time.Time `gorm:"primaryKey"`
	Name             string    `json:"name" gorm:"primaryKey;type:varchar(255)"`
	MeetingCount     string    `json:"meeting_count" gorm:"type:varchar(255)"`
	MeetingDuration  string    `json:"meeting_duration" gorm:"type:varchar(255)"`
	UserType         int64     `json:"user_type"`
}

func (FeishuMeetingTopUserItem20220524Temp) TableName() string {
	return "_tool_feishu_meeting_top_user_items_tmp"
}

type FeishuMeetingTopUserItem20220524 struct {
}

func (FeishuMeetingTopUserItem20220524) TableName() string {
	return "_tool_feishu_meeting_top_user_items"
}

type UpdateSchemas20220524 struct{}

func (*UpdateSchemas20220524) Up(ctx context.Context, db *gorm.DB) error {
	cursor, err := db.Model(archived.FeishuMeetingTopUserItem{}).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	// 1. create a temporary table to store unique records
	err = db.Migrator().CreateTable(FeishuMeetingTopUserItem20220524Temp{})
	if err != nil {
		return err
	}
	// 2. dedupe records and insert into the temporary table
	for cursor.Next() {
		inputRow := FeishuMeetingTopUserItem20220524Temp{}
		err := db.ScanRows(cursor, &inputRow)
		if err != nil {
			return err
		}
		err = db.Clauses(clause.OnConflict{UpdateAll: true}).Create(inputRow).Error
		if err != nil {
			return err
		}
	}
	// 3. drop old table
	err = db.Migrator().DropTable(archived.FeishuMeetingTopUserItem{})
	if err != nil {
		return err
	}
	// 4. rename the temporary table to the old table
	return db.Migrator().RenameTable(FeishuMeetingTopUserItem20220524Temp{}, FeishuMeetingTopUserItem20220524{})
}

func (*UpdateSchemas20220524) Version() uint64 {
	return 20220524000001
}

func (*UpdateSchemas20220524) Name() string {
	return "change primary column `id` to start_time+name"
}
