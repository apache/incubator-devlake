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
)

type GithubAccount20220713 struct {
	Name    string `json:"name" gorm:"type:varchar(255)"`
	Company string `json:"company" gorm:"type:varchar(255)"`
	Email   string `json:"Email" gorm:"type:varchar(255)"`
}

func (GithubAccount20220713) TableName() string {
	return "_tool_github_accounts"
}

type GithubRepoAccount20220713 struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	AccountId    int    `gorm:"primaryKey;autoIncrement:false"`
	RepoGithubId int    `gorm:"primaryKey;autoIncrement:false"`
	Login        string `json:"login" gorm:"type:varchar(255)"`
	archived.NoPKModel
}

func (GithubRepoAccount20220713) TableName() string {
	return "_tool_github_repo_accounts"
}

type updateSchemas20220713000003 struct{}

func (*updateSchemas20220713000003) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().AutoMigrate(&GithubAccount20220713{})
	if err != nil {
		return err
	}
	err = db.Migrator().AutoMigrate(&GithubRepoAccount20220713{})
	if err != nil {
		return err
	}
	return nil
}

func (*updateSchemas20220713000003) Version() uint64 {
	return 20220713000003
}

func (*updateSchemas20220713000003) Name() string {
	return "UpdateSchemas for extend account in 20220713"
}
