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
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
)

type GithubRelease struct {
	common.NoPKModel `json:"-" mapstructure:"-"`
	ConnectionId     uint64     `json:"connection_id" gorm:"primaryKey"`
	GithubId         int        `json:"github_id"`
	Id               string     `json:"id" gorm:"type:varchar(255);primaryKey"`
	AuthorName       string     `json:"authorName"`
	AuthorID         string     `json:"authorId"`
	CreatedAt        time.Time  `json:"createdAt"`
	DatabaseID       int        `json:"databaseId"`
	Description      string     `json:"description"`
	DescriptionHTML  string     `json:"descriptionHTML"`
	IsDraft          bool       `json:"isDraft"`
	IsLatest         bool       `json:"isLatest"`
	IsPrerelease     bool       `json:"isPrerelease"`
	Name             string     `json:"name"`
	PublishedAt      *time.Time `json:"publishedAt"`
	ResourcePath     string     `json:"resourcePath"`
	TagName          string     `json:"tagName"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	CommitSha        string     `json:"commit_sha"`
	URL              string     `json:"url"`
}

func (GithubRelease) TableName() string {
	return "_tool_github_releases"
}
