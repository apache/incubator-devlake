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

//type CodeComponent20220711 struct {
//	ComponentId string `gorm:"primaryKey;type:varchar(255)"`
//	PathRegex   string `gorm:"type:varchar(255)"`
//}
//
//func (CodeComponent20220711) TableName() string {
//	return "code_component_20220711"
//}

type Component struct {
	RepoId    string `gorm:"primaryKey;type:varchar(255)"`
	Component string `gorm:"primaryKey;type:varchar(255)"`
	PathRegex string `gorm:"type:varchar(255)"`
}

func (Component) TableName() string {
	return "component"
}

type CommitFile struct {
	common.NoPKModel
	CommitFileID string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha    string `gorm:"type:varchar(40)"`
	FilePath     string `gorm:"type:varchar(255)"`
	Additions    int
	Deletions    int
	Component    string `gorm:"type:varchar(255)"`
}

func (CommitFile) TableName() string {
	return "commit_files"
}

type FileComponent struct {
	common.NoPKModel
	CommitFileID string `gorm:"primaryKey;type:varchar(255)"`
	RepoId       string `gorm:"primaryKey;type:varchar(255)"`
	Component    string `gorm:"type:varchar(255)"`
}

func (FileComponent) TableName() string {
	return "file_component"
}

type updateSchemas20220711 struct{}

func (*updateSchemas20220711) Up(ctx context.Context, db *gorm.DB) error {

	err := db.Migrator().AutoMigrate(Component{}, CommitFile{}, FileComponent{})
	if err != nil {
		return err
	}
	return nil

}

func (*updateSchemas20220711) Version() uint64 {
	return 202207151420
}

func (*updateSchemas20220711) Name() string {
	return "file_component table"
}
