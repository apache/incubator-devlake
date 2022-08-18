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
	"time"

	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type addCICD struct{}

func (*addCICD) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().AutoMigrate(
		&CICDPipelineRepo{},
		&CICDPipeline{},
		&CICDTask{},
	)
	if err != nil {
		return err
	}

	return nil
}

func (*addCICD) Version() uint64 {
	return 20220818232735
}

func (*addCICD) Name() string {
	return "add cicd models"
}

type CICDPipeline struct {
	archived.DomainEntity
	Name         string `gorm:"type:varchar(255)"`
	CommitSha    string `gorm:"type:varchar(255);index"`
	Branch       string `gorm:"type:varchar(255);index"`
	Repo         string `gorm:"type:varchar(255);index"`
	Result       string `gorm:"type:varchar(100)"`
	Status       string `gorm:"type:varchar(100)"`
	Type         string `gorm:"type:varchar(100);comment: to indicate this is CI or CD"`
	DurationSec  uint64
	CreatedDate  time.Time
	FinishedDate *time.Time
}

func (CICDPipeline) TableName() string {
	return "cicd_pipelines"
}

type CICDTask struct {
	archived.DomainEntity
	Name         string `gorm:"type:varchar(255)"`
	PipelineId   string `gorm:"index;type:varchar(255)"`
	Result       string `gorm:"type:varchar(100)"`
	Status       string `gorm:"type:varchar(100)"`
	Type         string `gorm:"type:varchar(100);comment: to indicate this is CI or CD"`
	DurationSec  uint64
	StartedDate  time.Time
	FinishedDate *time.Time
}

func (CICDTask) TableName() string {
	return "cicd_tasks"
}

type CICDPipelineRepo struct {
	archived.DomainEntity
	CommitSha string `gorm:"primaryKey;type:varchar(255)"`
	Branch    string `gorm:"type:varchar(255)"`
	RepoUrl   string `gorm:"type:varchar(255)"`
}

func (CICDPipelineRepo) TableName() string {
	return "cicd_pipeline_repos"
}
