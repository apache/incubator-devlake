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
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type modifyAllEntities struct{}

// JenkinsBuild db entity for jenkins build
type JenkinsBuild0729 struct {
	Type string `gorm:"index;type:varchar(255)" `
}

func (JenkinsBuild0729) TableName() string {
	return "_tool_jenkins_builds"
}

type JenkinsJob0729 struct {
	HasUpstreamProjects bool
}

func (JenkinsJob0729) TableName() string {
	return "_tool_jenkins_jobs"
}

type JenkinsUpDownJob0729 struct {
	ConnetionId   uint64 `gorm:"primaryKey"`
	UpstreamJob   string `gorm:"primaryKey;type:varchar(255)"`
	DownstreamJob string `gorm:"primaryKey;type:varchar(255)"`
	archived.NoPKModel
}

func (JenkinsUpDownJob0729) TableName() string {
	return "_tool_jenkins_up_down_jobs"
}

type JenkinsBuildCommitRepoUrl0729 struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	BuildName    string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha    string `gorm:"primaryKey;type:varchar(255)"`
	RemoteUrl    string `gorm:"primaryKey;type:varchar(255)"`
	common.NoPKModel
}

func (JenkinsBuildCommitRepoUrl0729) TableName() string {
	return "_tool_jenkins_build_commit_repo_urls"
}

type JenkinsBuildTriggeredBuilds0729 struct {
	ConnectionId       uint64 `gorm:"primaryKey"`
	BuildName          string `gorm:"primaryKey;type:varchar(255)"`
	TriggeredBuildName string `gorm:"primaryKey;type:varchar(255)"`
	archived.NoPKModel
}

func (JenkinsBuildTriggeredBuilds0729) TableName() string {
	return "_tool_jenkins_build_triggered_builds"
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
}

func (JenkinsStage0729) TableName() string {
	return "_tool_jenkins_stages"
}

func (*modifyAllEntities) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().AddColumn(JenkinsBuild0729{}, "type")
	if err != nil {
		return err
	}
	err = db.Migrator().AddColumn(JenkinsJob0729{}, "has_upstream_projects")
	if err != nil {
		return err
	}
	err = db.Migrator().AutoMigrate(
		JenkinsUpDownJob0729{},
		JenkinsBuildCommitRepoUrl0729{},
		JenkinsBuildTriggeredBuilds0729{},
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
