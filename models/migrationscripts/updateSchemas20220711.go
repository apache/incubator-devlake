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
	"gorm.io/gorm"
)

type Component struct {
	RepoId    string `gorm:"type:varchar(255)"`
	Name      string `gorm:"primaryKey;type:varchar(255)"`
	PathRegex string `gorm:"type:varchar(255)"`
}

func (Component) TableName() string {
	return "components"
}

type CommitFile struct {
	common.NoPKModel
	ID        string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha string `gorm:"type:varchar(40)"`
	FilePath  string `gorm:"type:varchar(255)"`
	Additions int
	Deletions int
}

func (CommitFile) TableName() string {
	return "commit_files"
}

type CommitFileComponent struct {
	common.NoPKModel
	CommitFileID string `gorm:"primaryKey;type:varchar(255)"`
	Component    string `gorm:"type:varchar(255)"`
}

func (CommitFileComponent) TableName() string {
	return "commit_file_components"
}

type updateSchemas20220711 struct{}

func (*updateSchemas20220711) Up(ctx context.Context, db *gorm.DB) error {

	err := db.Migrator().AutoMigrate(Component{}, CommitFile{}, CommitFileComponent{})
	if err != nil {
		return err
	}
	return nil

}

func (*updateSchemas20220711) Version() uint64 {
	return 202207151644
}

func (*updateSchemas20220711) Name() string {
	return "add commit_file_components components table,update commit_files table"
}
