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

package devops

import (
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"time"
)

type CicdRelease struct {
	domainlayer.DomainEntity
	PublishedAt time.Time `json:"publishedAt"`

	CicdScopeId string `gorm:"index;type:varchar(255)"`

	Name         string `gorm:"type:varchar(255)"`
	DisplayTitle string `gorm:"type:varchar(255)"`
	Description  string `json:"description"`
	URL          string `json:"url"`

	IsDraft      bool `json:"isDraft"`
	IsLatest     bool `json:"isLatest"`
	IsPrerelease bool `json:"isPrerelease"`

	AuthorID string `json:"id" gorm:"type:varchar(255)"`

	RepoId string `gorm:"type:varchar(255)"`
	//RepoUrl string `gorm:"index;not null"`

	TagName   string `json:"tagName"`
	CommitSha string `gorm:"uniqueIndex;type:varchar(255)"`
	//CommitMsg string
	//RefName   string `gorm:"type:varchar(255)"`
}

func (CicdRelease) TableName() string {
	return "cicd_releases"
}
