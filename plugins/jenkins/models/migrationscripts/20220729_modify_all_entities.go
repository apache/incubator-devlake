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

type modifyAllEntities struct{}

type JenkinsJobDag0729 struct {
	ConnetionId   uint64 `gorm:"primaryKey"`
	UpstreamJob   string `gorm:"primaryKey;type:varchar(255)"`
	DownstreamJob string `gorm:"primaryKey;type:varchar(255)"`
	archived.NoPKModel
}

func (JenkinsJobDag0729) TableName() string {
	return "_tool_jenkins_job_dags"
}

// JenkinsBuild db entity for jenkins build
type JenkinsBuild0729 struct {
	TriggeredBy string `gorm:"type:varchar(255)"`
	Type        string `gorm:"index;type:varchar(255)" `
	Class       string `gorm:"index;type:varchar(255)" `
	HasStages   bool
	Building    bool
}

func (JenkinsBuild0729) TableName() string {
	return "_tool_jenkins_builds"
}

type JenkinsBuildRepo0729 struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	BuildName    string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha    string `gorm:"primaryKey;type:varchar(255)"`
	Branch       string `gorm:"type:varchar(255)"`
	RepoUrl      string `gorm:"type:varchar(255)"`
	archived.NoPKModel
}

func (JenkinsBuildRepo0729) TableName() string {
	return "_tool_jenkins_build_repos"
}

type JenkinsStage0729 struct {
	archived.NoPKModel
	ConnectionId        uint64 `gorm:"primaryKey"`
	ID                  string `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Name                string `json:"name" gorm:"type:varchar(255)"`
	ExecNode            string `json:"execNode" gorm:"type:varchar(255)"`
	Status              string `json:"status" gorm:"type:varchar(255)"`
	StartTimeMillis     int64  `json:"startTimeMillis"`
	DurationMillis      int    `json:"durationMillis"`
	PauseDurationMillis int    `json:"pauseDurationMillis"`
	BuildName           string `gorm:"primaryKey;type:varchar(255)"`
	Type                string `gorm:"index;type:varchar(255)"`
}

func (JenkinsStage0729) TableName() string {
	return "_tool_jenkins_stages"
}

func (*modifyAllEntities) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().AddColumn(JenkinsBuild0729{}, "type")
	if err != nil {
		return err
	}
	err = db.Migrator().AddColumn(JenkinsBuild0729{}, "class")
	if err != nil {
		return err
	}
	err = db.Migrator().AddColumn(JenkinsBuild0729{}, "triggered_by")
	if err != nil {
		return err
	}
	err = db.Migrator().AddColumn(JenkinsBuild0729{}, "building")
	if err != nil {
		return err
	}
	err = db.Migrator().AddColumn(JenkinsBuild0729{}, "has_stages")
	if err != nil {
		return err
	}
	err = db.Migrator().AutoMigrate(
		JenkinsJobDag0729{},
		JenkinsBuildRepo0729{},
		JenkinsStage0729{},
	)
	if err != nil {
		return err
	}

	return nil
}

func (*modifyAllEntities) Version() uint64 {
	return 20220729231237
}

func (*modifyAllEntities) Name() string {
	return "Jenkins modify build and job"
}
