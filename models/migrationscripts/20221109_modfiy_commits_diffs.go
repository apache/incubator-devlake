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

type CommitsStatus20221109 struct {
	NewCommitSha string `gorm:"primaryKey;type:varchar(40)"`
	OldCommitSha string `gorm:"primaryKey;type:varchar(40)"`
}

func (CommitsStatus20221109) TableName() string {
	return "commits_diffs_status"
}

type CommitsDiff20221109Before struct {
	CommitSha    string `gorm:"primaryKey;type:varchar(40)"`
	NewCommitSha string `gorm:"type:varchar(40)"`
	OldCommitSha string `gorm:"type:varchar(40)"`
	SortingIndex int
}

type CommitsDiff20221109After struct {
	CommitSha    string `gorm:"primaryKey;type:varchar(40)"`
	NewCommitSha string `gorm:"primaryKey;type:varchar(40)"`
	OldCommitSha string `gorm:"primaryKey;type:varchar(40)"`
	SortingIndex int
}

func (script *modifyCommitsDiffs) Up(basicRes core.BasicRes) errors.Error {
	db := basicRes.GetDal()
	// rename table
	err := db.RenameTable("refs_commits_diffs", "commits_diffs")
	if err != nil {
		return err
	}
	// copy data
	err = migrationhelper.TransformTable(
		basicRes,
		script,
		"commits_diffs",
		func(c *CommitsDiff20221109Before) (*CommitsDiff20221109After, errors.Error) {
			dst := CommitsDiff20221109After(*c)
			return &dst, nil
		},
	)
	if err != nil {
		return err
	}
	// add new table: refs_commits_status
	return db.AutoMigrate(&CommitsStatus20221109{})
}

func (*modifyCommitsDiffs) Version() uint64 {
	return 20221109232735
}

func (*modifyCommitsDiffs) Name() string {
	return "modify commits diffs"
}
