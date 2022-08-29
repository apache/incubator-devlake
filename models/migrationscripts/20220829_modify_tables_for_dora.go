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
	"gorm.io/gorm"
)

type modifyTablesForDora struct{}

func (*modifyTablesForDora) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().AutoMigrate(
		&CICDPipeline0829{},
		&PullRequest0829{},
		&Issue0829{},
	)
	if err != nil {
		return err
	}
	return nil
}

func (*modifyTablesForDora) Version() uint64 {
	return 20220829232735
}

func (*modifyTablesForDora) Name() string {
	return "modify tables for dora"
}

type PullRequest0829 struct {
	CodingTimespan uint64
	ReviewLag      uint64
	ReviewTimespan uint64
	DeployTimespan uint64
	ChangeTimespan uint64
}

func (PullRequest0829) TableName() string {
	return "pull_requests"
}

type Issue0829 struct {
	DeploymentId string `gorm:"type:varchar(255)"`
}

func (Issue0829) TableName() string {
	return "issues"
}

type CICDPipeline0829 struct {
	Environment string `gorm:"type:varchar(255)"`
}

func (CICDPipeline0829) TableName() string {
	return "cicd_pipelines"
}
