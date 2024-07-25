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

package models

import (
	"github.com/apache/incubator-devlake/core/models/common"
	"time"
)

type AzuredevopsPullRequest struct {
	common.NoPKModel

	ConnectionId    uint64 `gorm:"primaryKey"`
	AzuredevopsId   int    `gorm:"primaryKey"`
	RepositoryId    string
	CreationDate    *time.Time
	MergeCommitSha  string
	TargetRefName   string
	Description     string
	Status          string
	SourceRefName   string
	SourceCommitSha string
	TargetCommitSha string
	Type            string
	CreatedById     string
	CreatedByName   string
	ClosedDate      *time.Time
	Title           string
	ForkRepoId      string
	Url             string
}

func (AzuredevopsPullRequest) TableName() string {
	return "_tool_azuredevops_go_pull_requests"
}

type AzuredevopsApiPullRequest struct {
	Labels []struct {
		Id     string `json:"id"`
		Name   string `json:"name"`
		Active bool   `json:"active"`
	} `json:"labels"`
	Repository struct {
		Id      string `json:"id"`
		Name    string `json:"name"`
		Url     string `json:"url"`
		Project struct {
			Id             string `json:"id"`
			Name           string `json:"name"`
			State          string `json:"state"`
			Visibility     string `json:"visibility"`
			LastUpdateTime string `json:"lastUpdateTime"`
		} `json:"project"`
	} `json:"repository"`
	PullRequestId int    `json:"pullRequestId"`
	CodeReviewId  int    `json:"codeReviewId"`
	Status        string `json:"status"`
	CreatedBy     struct {
		DisplayName string `json:"displayName"`
		Url         string `json:"url"`
		Id          string `json:"id"`
		UniqueName  string `json:"uniqueName"`
		ImageUrl    string `json:"imageUrl"`
		Descriptor  string `json:"descriptor"`
	} `json:"createdBy"`
	AzuredevopsCreationDate *common.Iso8601Time `json:"creationDate"`
	ClosedDate              *common.Iso8601Time `json:"closedDate"`
	Title                   string              `json:"title"`
	Description             string              `json:"description"`
	SourceRefName           string              `json:"sourceRefName"`
	TargetRefName           string              `json:"targetRefName"`
	MergeStatus             string              `json:"mergeStatus"`
	IsDraft                 bool                `json:"isDraft"`
	MergeId                 string              `json:"mergeId"`
	LastMergeSourceCommit   struct {
		CommitId string `json:"commitId"`
		Url      string `json:"url"`
	} `json:"lastMergeSourceCommit"`
	LastMergeTargetCommit struct {
		CommitId string `json:"commitId"`
		Url      string `json:"url"`
	} `json:"lastMergeTargetCommit"`
	LastMergeCommit struct {
		CommitId string `json:"commitId"`
		Url      string `json:"url"`
	} `json:"lastMergeCommit"`
	Url                 string     `json:"url"`
	SupportsIterations  bool       `json:"supportsIterations"`
	CompletionQueueTime *time.Time `json:"completionQueueTime"`
}
