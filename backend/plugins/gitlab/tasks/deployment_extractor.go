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

package tasks

import (
	"encoding/json"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"time"
)

var _ plugin.SubTaskEntryPoint = ExtractDeployment

func init() {
	RegisterSubtaskMeta(ExtractDeploymentMeta)
}

var ExtractDeploymentMeta = &plugin.SubTaskMeta{
	Name:             "ExtractDeployment",
	EntryPoint:       ExtractDeployment,
	EnabledByDefault: true,
	Description:      "Extract gitlab deployment from raw layer to tool layer",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	Dependencies:     []*plugin.SubTaskMeta{CollectDeploymentMeta},
}

func ExtractDeployment(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_DEPLOYMENT)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			deploymentResp := &GitlabDeploymentResp{}
			err := errors.Convert(json.Unmarshal(row.Data, deploymentResp))
			if err != nil {
				return nil, err
			}
			gitlabDeployment := deploymentResp.toGitlabDeployment(data.Options.ConnectionId)
			return []interface{}{
				gitlabDeployment,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

type GitlabDeploymentResp struct {
	CreatedAt   time.Time                   `json:"created_at"`
	UpdatedAt   time.Time                   `json:"updated_at"`
	Status      string                      `json:"status"`
	Deployable  GitlabDeploymentDeployable  `json:"deployable"`
	Environment GitlabDeploymentEnvironment `json:"environment"`
	ID          int                         `json:"id"`
	Iid         int                         `json:"iid"`
	Ref         string                      `json:"ref"`
	Sha         string                      `json:"sha"`
	User        GitlabDeploymentSimpleUser  `json:"user"`
}

func (r GitlabDeploymentResp) toGitlabDeployment(connectionId uint64) *models.GitlabDeployment {
	ret := models.GitlabDeployment{
		NoPKModel:                   common.NewNoPKModel(),
		ConnectionId:                connectionId,
		CreatedDate:                 r.CreatedAt,
		UpdatedDate:                 r.UpdatedAt,
		Status:                      r.Status,
		DeploymentId:                r.ID,
		Iid:                         r.Iid,
		Ref:                         r.Ref,
		Sha:                         r.Sha,
		Environment:                 r.Environment.Name,
		Name:                        r.Deployable.Name,
		DeployableCommitAuthorEmail: r.Deployable.Commit.AuthorEmail,
		DeployableCommitAuthorName:  r.Deployable.Commit.AuthorName,
		DeployableCommitCreatedAt:   r.Deployable.Commit.CreatedAt,
		DeployableCommitID:          r.Deployable.Commit.ID,
		DeployableCommitMessage:     r.Deployable.Commit.Message,
		DeployableCommitShortID:     r.Deployable.Commit.ShortID,
		DeployableCommitTitle:       r.Deployable.Commit.Title,
		//DeployableCoverage:          r.Deployable.Coverage,
		DeployableID:   r.Deployable.ID,
		DeployableName: r.Deployable.Name,
		DeployableRef:  r.Deployable.Ref,
		//DeployableRunner:            r.Deployable.Runner,
		DeployableStage:         r.Deployable.Stage,
		DeployableStatus:        r.Deployable.Status,
		DeployableTag:           r.Deployable.Tag,
		DeployableUserID:        r.Deployable.User.ID,
		DeployableUserName:      r.Deployable.User.Name,
		DeployableUserUsername:  r.Deployable.User.Username,
		DeployableUserState:     r.Deployable.User.State,
		DeployableUserAvatarURL: r.Deployable.User.AvatarURL,
		DeployableUserWebURL:    r.Deployable.User.WebURL,
		DeployableUserCreatedAt: r.Deployable.User.CreatedAt,
		//DeployableUserBio:           r.Deployable.User.Bio,
		//DeployableUserLocation:      r.Deployable.User.Location,
		DeployableUserPublicEmail:   r.Deployable.User.PublicEmail,
		DeployableUserSkype:         r.Deployable.User.Skype,
		DeployableUserLinkedin:      r.Deployable.User.Linkedin,
		DeployableUserTwitter:       r.Deployable.User.Twitter,
		DeployableUserWebsiteURL:    r.Deployable.User.WebsiteURL,
		DeployableUserOrganization:  r.Deployable.User.Organization,
		DeployablePipelineCreatedAt: r.Deployable.Pipeline.CreatedAt,
		DeployablePipelineID:        r.Deployable.Pipeline.ID,
		DeployablePipelineRef:       r.Deployable.Pipeline.Ref,
		DeployablePipelineSha:       r.Deployable.Pipeline.Sha,
		DeployablePipelineStatus:    r.Deployable.Pipeline.Status,
		DeployablePipelineUpdatedAt: r.Deployable.Pipeline.UpdatedAt,
		DeployablePipelineWebURL:    r.Deployable.Pipeline.WebURL,
		UserAvatarURL:               r.User.AvatarURL,
		UserID:                      r.User.ID,
		UserName:                    r.User.Name,
		UserState:                   r.User.State,
		UserUsername:                r.User.Username,
		UserWebURL:                  r.User.WebURL,
	}
	if r.Deployable.StartedAt != nil {
		ret.DeployableStartedAt = r.Deployable.StartedAt
	}
	if r.Deployable.CreatedAt != nil {
		ret.DeployableCreatedAt = r.Deployable.CreatedAt
	}
	if r.Deployable.FinishedAt != nil {
		ret.DeployableFinishedAt = r.Deployable.FinishedAt
	}
	if r.Deployable.Duration != nil {
		ret.DeployableDuration = *r.Deployable.Duration
	}
	return &ret
}

type GitlabDeploymentCommit struct {
	AuthorEmail string    `json:"author_email"`
	AuthorName  string    `json:"author_name"`
	CreatedAt   time.Time `json:"created_at"`
	ID          string    `json:"id"`
	Message     string    `json:"message"`
	ShortID     string    `json:"short_id"`
	Title       string    `json:"title"`
}

type GitlabDeploymentProject struct {
	CiJobTokenScopeEnabled bool `json:"ci_job_token_scope_enabled"`
}

type GitlabDeploymentFullUser struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Username     string    `json:"username"`
	State        string    `json:"state"`
	AvatarURL    string    `json:"avatar_url"`
	WebURL       string    `json:"web_url"`
	CreatedAt    time.Time `json:"created_at"`
	Bio          any       `json:"bio"`
	Location     any       `json:"location"`
	PublicEmail  string    `json:"public_email"`
	Skype        string    `json:"skype"`
	Linkedin     string    `json:"linkedin"`
	Twitter      string    `json:"twitter"`
	WebsiteURL   string    `json:"website_url"`
	Organization string    `json:"organization"`
}

type GitlabDeploymentPipeline struct {
	CreatedAt time.Time `json:"created_at"`
	ID        int       `json:"id"`
	Ref       string    `json:"ref"`
	Sha       string    `json:"sha"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	WebURL    string    `json:"web_url"`
}

type GitlabDeploymentDeployable struct {
	Commit     GitlabDeploymentCommit   `json:"commit"`
	Coverage   any                      `json:"coverage"`
	CreatedAt  *time.Time               `json:"created_at"`
	FinishedAt *time.Time               `json:"finished_at"`
	ID         int                      `json:"id"`
	Name       string                   `json:"name"`
	Ref        string                   `json:"ref"`
	Runner     any                      `json:"runner"`
	Stage      string                   `json:"stage"`
	StartedAt  *time.Time               `json:"started_at"`
	Status     string                   `json:"status"`
	Tag        bool                     `json:"tag"`
	Project    GitlabDeploymentProject  `json:"project"`
	User       GitlabDeploymentFullUser `json:"user"`
	Pipeline   GitlabDeploymentPipeline `json:"pipeline"`

	AllowFailure bool       `json:"allow_failure"`
	ErasedAt     *time.Time `json:"erased_at"`

	Duration          *float64                    `json:"duration"`
	QueuedDuration    *float64                    `json:"queued_duration"`
	ArtifactsExpireAt *time.Time                  `json:"artifacts_expire_at "`
	TagList           []string                    `json:"tag_list"`
	Artifacts         []GitlabDeploymentArtifacts `json:"artifacts"`
}

type GitlabDeploymentEnvironment struct {
	ExternalURL string `json:"external_url"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
}

type GitlabDeploymentArtifacts struct {
	FileType   string `json:"file_type"`
	Size       int    `json:"size"`
	Filename   string `json:"filename"`
	FileFormat any    `json:"file_format"`
}

type GitlabDeploymentSimpleUser struct {
	AvatarURL string `json:"avatar_url"`
	ID        int    `json:"id"`
	Name      string `json:"name"`
	State     string `json:"state"`
	Username  string `json:"username"`
	WebURL    string `json:"web_url"`
}
