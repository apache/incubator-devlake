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
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/plugins/core"
)

var _ core.MigrationScript = (*addCommitFileComponent)(nil)

type component20220722 struct {
	RepoId    string `gorm:"type:varchar(255)"`
	Name      string `gorm:"primaryKey;type:varchar(255)"`
	PathRegex string `gorm:"type:varchar(255)"`
}

func (component20220722) TableName() string {
	return "components"
}

type commitFile20220722 struct {
	archived.DomainEntity
	CommitSha string `gorm:"type:varchar(40)"`
	FilePath  string `gorm:"type:varchar(255)"`
	Additions int
	Deletions int
}

func (commitFile20220722) TableName() string {
	return "commit_files"
}

type commitFileComponent20220722 struct {
	archived.NoPKModel
	CommitFileId  string `gorm:"primaryKey;type:varchar(255)"`
	ComponentName string `gorm:"type:varchar(255)"`
}

func (commitFileComponent20220722) TableName() string {
	return "commit_file_components"
}

type addCommitFileComponent struct{}

func (addCommitFileComponent) Up(basicRes core.BasicRes) errors.Error {
	db := basicRes.GetDal()
	err := db.DropTables(&archived.CommitFile{})
	if err != nil {
		return err
	}

	return migrationhelper.AutoMigrateTables(
		basicRes,
		component20220722{},
		commitFile20220722{},
		commitFileComponent20220722{},
	)

}

func (*addCommitFileComponent) Version() uint64 {
	return 20220722165805
}

func (*addCommitFileComponent) Name() string {

	return "add commit_file_components components table,update commit_files table"
}
