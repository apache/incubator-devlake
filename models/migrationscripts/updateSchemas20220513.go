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

	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RefsIssuesDiffs20220513 struct {
	NewRefId        string `gorm:"primaryKey;type:varchar(255)"`
	OldRefId        string `gorm:"primaryKey;type:varchar(255)"`
	NewRefCommitSha string `gorm:"type:varchar(40)"`
	OldRefCommitSha string `gorm:"type:varchar(40)"`
	IssueNumber     string `gorm:"type:varchar(255)"`
	IssueId         string `gorm:"primaryKey;type:varchar(255)"`
	archived.NoPKModel
}

func (RefsIssuesDiffs20220513) TableName() string {
	return "refs_issues_diffs_20220513"
}

type RefsIssuesDiffsNew struct {
	NewRefId        string `gorm:"primaryKey;type:varchar(255)"`
	OldRefId        string `gorm:"primaryKey;type:varchar(255)"`
	NewRefCommitSha string `gorm:"type:varchar(40)"`
	OldRefCommitSha string `gorm:"type:varchar(40)"`
	IssueNumber     string `gorm:"type:varchar(255)"`
	IssueId         string `gorm:"primaryKey;type:varchar(255)"`
	archived.NoPKModel
}

func (RefsIssuesDiffsNew) TableName() string {
	return "refs_issues_diffs"
}

type updateSchemas20220513 struct{}

func (*updateSchemas20220513) Up(ctx context.Context, db *gorm.DB) error {
	cursor, err := db.Model(archived.RefsIssuesDiffs{}).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	// 1. create a temporary table to store unique records
	err = db.Migrator().CreateTable(RefsIssuesDiffs20220513{})
	if err != nil {
		return err
	}
	// 2. dedupe records and insert into the temporary table
	for cursor.Next() {
		//inputRow := archived.RefsIssuesDiffs{}
		inputRow := RefsIssuesDiffs20220513{}
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
	err = db.Migrator().DropTable(archived.RefsIssuesDiffs{})
	if err != nil {
		return err
	}
	// 4. rename the temporary table to the old table
	err = db.Migrator().RenameTable(RefsIssuesDiffs20220513{}, RefsIssuesDiffsNew{})
	if err != nil {
		return err
	}

	return nil
}

func (*updateSchemas20220513) Version() uint64 {
	return 20220513212319
}

func (*updateSchemas20220513) Name() string {
	return "refs_issues_diffs add primary key"
}
