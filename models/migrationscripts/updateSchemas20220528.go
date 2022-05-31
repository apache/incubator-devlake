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
	"time"

	"gorm.io/gorm"

	"github.com/apache/incubator-devlake/models/common"
)

type SprintIssue20220528 struct {
	common.NoPKModel
	SprintId      string `gorm:"primaryKey;type:varchar(255)"`
	IssueId       string `gorm:"primaryKey;type:varchar(255)"`
	IsRemoved     bool
	AddedDate     *time.Time
	RemovedDate   *time.Time
	AddedStage    *string `gorm:"type:varchar(255)"`
	ResolvedStage *string `gorm:"type:varchar(255)"`
}

func (SprintIssue20220528) TableName() string {
	return "sprint_issues"
}

type updateSchemas20220528 struct{}

func (*updateSchemas20220528) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().DropColumn(&SprintIssue20220528{}, "is_removed")
	if err != nil {
		return err
	}
	err = db.Migrator().DropColumn(&SprintIssue20220528{}, "added_date")
	if err != nil {
		return err
	}
	err = db.Migrator().DropColumn(&SprintIssue20220528{}, "removed_date")
	if err != nil {
		return err
	}
	err = db.Migrator().DropColumn(&SprintIssue20220528{}, "added_stage")
	if err != nil {
		return err
	}
	return db.Migrator().DropColumn(&SprintIssue20220528{}, "resolved_stage")
}

func (*updateSchemas20220528) Version() uint64 {
	return 20220528110537
}

func (*updateSchemas20220528) Name() string {
	return "remove columns: is_removed, added_date, removed_date, added_stage, resolved_stage"
}
