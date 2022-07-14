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

type GithubAccountOrg20220713 struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	AccountId    int    `gorm:"primaryKey;autoIncrement:false"`
	OrgId        int    `gorm:"primaryKey;autoIncrement:false"`
	OrgLogin     string `json:"org_login" gorm:"type:varchar(255)"`
	archived.NoPKModel
}

func (GithubAccountOrg20220713) TableName() string {
	return "_tool_github_account_orgs"
}

type updateSchemas20220713000004 struct{}

func (*updateSchemas20220713000004) Up(ctx context.Context, db *gorm.DB) error {
	return db.Migrator().AutoMigrate(&GithubAccountOrg20220713{})
}

func (*updateSchemas20220713000004) Version() uint64 {
	return 20220713000004
}

func (*updateSchemas20220713000004) Name() string {
	return "UpdateSchemas for add org in 20220713"
}
