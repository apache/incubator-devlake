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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/common"
	"gorm.io/gorm"
)

type RepoSnapshot struct {
	common.NoPKModel
	RepoId    string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha string `gorm:"primaryKey;type:varchar(40);"`
	FilePath  string `gorm:"primaryKey;type:varchar(255);"`
	LineNo    int    `gorm:"primaryKey;type:int;"`
}

func (RepoSnapshot) TableName() string {
	return "repo_snapshot"
}

type addRepoSnapshot struct{}

func (*addRepoSnapshot) Up(ctx context.Context, db *gorm.DB) errors.Error {
	err := db.Migrator().AutoMigrate(RepoSnapshot{})
	if err != nil {
		return errors.Convert(err)
	}
	return nil
}

func (*addRepoSnapshot) Version() uint64 {
	return 20221009111241
}

func (*addRepoSnapshot) Name() string {
	return "add snapshot table"
}
