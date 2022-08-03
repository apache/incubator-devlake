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

	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/gitlab/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type modifyGitlabCI struct{}

type GitlabPipeline20220729 struct {
	ConnectionId uint64 `gorm:"primaryKey"`

	GitlabId  int    `gorm:"primaryKey"`
	ProjectId int    `gorm:"index"`
	Status    string `gorm:"type:varchar(100)"`
	Ref       string `gorm:"type:varchar(255)"`
	Sha       string `gorm:"type:varchar(255)"`
	WebUrl    string `gorm:"type:varchar(255)"`
	Duration  int

	GitlabCreatedAt *time.Time
	GitlabUpdatedAt *time.Time
	StartedAt       *time.Time
	FinishedAt      *time.Time
	Coverage        string

	common.NoPKModel
}

func (GitlabPipeline20220729) TableName() string {
	return "_tool_gitlab_pipelines"
}

type GitlabJob20220729 struct {
	ConnectionId uint64 `gorm:"primaryKey"`

	GitlabId     int     `gorm:"primaryKey"`
	ProjectId    int     `gorm:"index"`
	Status       string  `gorm:"type:varchar(255)"`
	Stage        string  `gorm:"type:varchar(255)"`
	Name         string  `gorm:"type:varchar(255)"`
	Ref          string  `gorm:"type:varchar(255)"`
	Tag          bool    `gorm:"type:boolean"`
	AllowFailure bool    `json:"allow_failure"`
	Duration     float64 `gorm:"type:text"`
	WebUrl       string  `gorm:"type:varchar(255)"`

	GitlabCreatedAt *time.Time
	StartedAt       *time.Time
	FinishedAt      *time.Time

	common.NoPKModel
}

func (GitlabJob20220729) TableName() string {
	return "_tool_gitlab_jobs"
}

func (*modifyGitlabCI) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().DropColumn(&archived.GitlabPipeline{}, "started_at")
	if err != nil {
		return err
	}
	err = db.Migrator().DropColumn(&archived.GitlabPipeline{}, "finished_at")
	if err != nil {
		return err
	}

	err = db.Migrator().AddColumn(&GitlabPipeline20220729{}, "gitlab_updated_at")
	if err != nil {
		return err
	}

	err = db.Migrator().AddColumn(&GitlabPipeline20220729{}, "started_at")
	if err != nil {
		return err
	}

	err = db.Migrator().AddColumn(&GitlabPipeline20220729{}, "finished_at")
	if err != nil {
		return err
	}

	err = db.Migrator().AutoMigrate(&GitlabJob20220729{})
	if err != nil {
		return err
	}

	return nil
}

func (*modifyGitlabCI) Version() uint64 {
	return 20220729231236
}

func (*modifyGitlabCI) Name() string {
	return "pipeline and job"
}
