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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/core"
)

var _ core.MigrationScript = (*modifyCommitsDiffs)(nil)

type modifyCommitsDiffs struct{}

type FinishedCommitsDiffs20221109 struct {
	NewCommitSha string `gorm:"primaryKey;type:varchar(40)"`
	OldCommitSha string `gorm:"primaryKey;type:varchar(40)"`
}

func (FinishedCommitsDiffs20221109) TableName() string {
	return "finished_commits_diffs"
}

type RefsCommitsDiff20221109 struct {
	NewRefId        string `gorm:"primaryKey;type:varchar(255)"`
	OldRefId        string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha       string `gorm:"primaryKey;type:varchar(40)"`
	NewRefCommitSha string `gorm:"type:varchar(40)"`
	OldRefCommitSha string `gorm:"type:varchar(40)"`
	SortingIndex    int
}

func (RefsCommitsDiff20221109) TableName() string {
	return "refs_commits_diffs"
}

type CommitsDiff20221109 struct {
	CommitSha    string `gorm:"primaryKey;type:varchar(40)"`
	NewCommitSha string `gorm:"primaryKey;type:varchar(40)"`
	OldCommitSha string `gorm:"primaryKey;type:varchar(40)"`
	SortingIndex int
}

func (CommitsDiff20221109) TableName() string {
	return "commits_diffs"
}

func (script *modifyCommitsDiffs) Up(basicRes core.BasicRes) errors.Error {
	db := basicRes.GetDal()
	// create table
	err := db.AutoMigrate(&CommitsDiff20221109{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&FinishedCommitsDiffs20221109{})
	if err != nil {
		return err
	}

	srcColumn := []string{"CommitSha", "NewRefCommitSha", "OldRefCommitSha", "SortingIndex"}
	dstColumn := []string{"CommitSha", "NewCommitSha", "OldCommitSha", "SortingIndex"}

	// copy data
	err = migrationhelper.CopyTableColumn(basicRes, &RefsCommitsDiff20221109{}, srcColumn, &CommitsDiff20221109{}, dstColumn)
	if err != nil {
		return err
	}

	return db.DropTables(&RefsCommitsDiff20221109{})
}

func (*modifyCommitsDiffs) Version() uint64 {
	return 20221109232735
}

func (*modifyCommitsDiffs) Name() string {
	return "modify commits diffs"
}
