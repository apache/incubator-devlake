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
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"time"
)

type gitlabDeployment20230926 struct {
	archived.NoPKModel `swaggerignore:"true" json:"-" mapstructure:"-"`

	ConnectionId uint64 `gorm:"primaryKey"`

	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
	Status      string    `json:"status"`

	DeploymentId int    `json:"id" gorm:"primaryKey"`
	Iid          int    `json:"iid"`
	Ref          string `json:"ref"`
	Sha          string `json:"sha"`
	Environment  string `json:"environment" gorm:"type:varchar(255)"`
	Name         string `json:"name" gorm:"type:varchar(255)"`

	DeployableCommitAuthorEmail string    `json:"deployable_commit_author_email" gorm:"type:varchar(255)"`
	DeployableCommitAuthorName  string    `json:"deployable_commit_author_name" gorm:"type:varchar(255)"`
	DeployableCommitCreatedAt   time.Time `json:"deployable_commit_created_at"`
	DeployableCommitID          string    `json:"deployable_commit_id" gorm:"type:varchar(255)"`
	DeployableCommitMessage     string    `json:"deployable_commit_message" gorm:"type:varchar(255)"`
	DeployableCommitShortID     string    `json:"deployable_commit_short_id" gorm:"type:varchar(255)"`
	DeployableCommitTitle       string    `json:"deployable_commit_title" gorm:"type:varchar(255)"`

	//DeployableCoverage   any       `json:"deployable_coverage"`
	DeployableCreatedAt  *time.Time `json:"deployable_created_at"`
	DeployableFinishedAt *time.Time `json:"deployable_finished_at"`
	DeployableID         int        `json:"deployable_id"`
	DeployableName       string     `json:"deployable_name" gorm:"type:varchar(255)"`
	DeployableRef        string     `json:"deployable_ref" gorm:"type:varchar(255)"`
	//DeployableRunner     any       `json:"deployable_runner"`
	DeployableStage     string     `json:"deployable_stage" gorm:"type:varchar(255)"`
	DeployableStartedAt *time.Time `json:"deployable_started_at"`
	DeployableStatus    string     `json:"deployable_status" gorm:"type:varchar(255)"`
	DeployableTag       bool       `json:"deployable_tag"`
	DeployableDuration  float64    `json:"deployable_duration"`
	QueuedDuration      float64    `json:"queued_duration"`

	DeployableUserID        int       `json:"deployable_user_id"`
	DeployableUserName      string    `json:"deployable_user_name" gorm:"type:varchar(255)"`
	DeployableUserUsername  string    `json:"deployable_user_username" gorm:"type:varchar(255)"`
	DeployableUserState     string    `json:"deployable_user_state" gorm:"type:varchar(255)"`
	DeployableUserAvatarURL string    `json:"deployable_user_avatar_url" gorm:"type:varchar(255)"`
	DeployableUserWebURL    string    `json:"deployable_user_web_url" gorm:"type:varchar(255)"`
	DeployableUserCreatedAt time.Time `json:"deployable_user_created_at"`
	//DeployableUserBio          any       `json:"deployable_user_bio"`
	//DeployableUserLocation     any       `json:"deployable_user_location"`
	DeployableUserPublicEmail  string `json:"deployable_user_public_email" gorm:"type:varchar(255)"`
	DeployableUserSkype        string `json:"deployable_user_skype" gorm:"type:varchar(255)"`
	DeployableUserLinkedin     string `json:"deployable_user_linkedin" gorm:"type:varchar(255)"`
	DeployableUserTwitter      string `json:"deployable_user_twitter" gorm:"type:varchar(255)"`
	DeployableUserWebsiteURL   string `json:"deployable_user_website_url" gorm:"type:varchar(255)"`
	DeployableUserOrganization string `json:"deployable_user_organization" gorm:"type:varchar(255)"`

	DeployablePipelineCreatedAt time.Time `json:"deployable_pipeline_created_at"`
	DeployablePipelineID        int       `json:"deployable_pipeline_id"`
	DeployablePipelineRef       string    `json:"deployable_pipeline_ref" gorm:"type:varchar(255)"`
	DeployablePipelineSha       string    `json:"deployable_pipeline_sha" gorm:"type:varchar(255)"`
	DeployablePipelineStatus    string    `json:"deployable_pipeline_status" gorm:"type:varchar(255)"`
	DeployablePipelineUpdatedAt time.Time `json:"deployable_pipeline_updated_at"`
	DeployablePipelineWebURL    string    `json:"deployable_pipeline_web_url" gorm:"type:varchar(255)"`

	UserAvatarURL string `json:"user_avatar_url" gorm:"type:varchar(255)"`
	UserID        int    `json:"user_id"`
	UserName      string `json:"user_name" gorm:"type:varchar(255)"`
	UserState     string `json:"user_state" gorm:"type:varchar(255)"`
	UserUsername  string `json:"user_username" gorm:"type:varchar(255)"`
	UserWebURL    string `json:"user_web_url" gorm:"type:varchar(255)"`
}

func (gitlabDeployment20230926) TableName() string {
	return "_tool_gitlab_deployments"
}

type addDeployment struct {
}

func (addDeployment) Up(basicRes context.BasicRes) errors.Error {
	return basicRes.GetDal().AutoMigrate(&gitlabDeployment20230926{})
}

func (addDeployment) Version() uint64 {
	return 20230926140000
}

func (addDeployment) Name() string {
	return "add deployment table in tool layer"
}
