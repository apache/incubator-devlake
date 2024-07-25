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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	coreArchived "github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/github/models/migrationscripts/archived"
)

var _ plugin.MigrationScript = (*restructReviewer)(nil)

type reviewer20240711 struct {
	ConnectionId  uint64 `gorm:"primaryKey"`
	ReviewerId    int    `gorm:"primaryKey"`
	PullRequestId int    `gorm:"primaryKey"`
	Name          string `gorm:"type:varchar(255)"`
	Username      string `gorm:"type:varchar(255)"`
	State         string `gorm:"type:varchar(255)"`
	AvatarUrl     string `gorm:"type:varchar(255)"`
	WebUrl        string `gorm:"type:varchar(255)"`
	coreArchived.NoPKModel
}

func (reviewer20240711) TableName() string {
	return "_tool_github_reviewers"
}

type restructReviewer struct{}

func (*restructReviewer) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if err := db.DropTables(&archived.GithubReviewer{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&reviewer20240711{}); err != nil {
		return err
	}
	return nil
}

func (*restructReviewer) Version() uint64 {
	return 20240710142104
}

func (*restructReviewer) Name() string {
	return "restruct reviewer table"
}
