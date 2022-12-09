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

package code

import (
	"time"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/plugins/core"
)

var _ core.Scope = (*Repo)(nil)

type Repo struct {
	domainlayer.DomainEntity
	Name        string     `json:"name"`
	Url         string     `json:"url"`
	Description string     `json:"description"`
	OwnerId     string     `json:"ownerId" gorm:"type:varchar(255)"`
	Language    string     `json:"language" gorm:"type:varchar(255)"`
	ForkedFrom  string     `json:"forkedFrom"`
	CreatedDate *time.Time `json:"createdDate"`
	UpdatedDate *time.Time `json:"updatedDate"`
	Deleted     bool       `json:"deleted"`
}

func (Repo) TableName() string {
	return "repos"
}

type RepoLanguage struct {
	RepoId   string `json:"repoId" gorm:"index;type:varchar(255)"`
	Language string `json:"language" gorm:"type:varchar(255)"`
	Bytes    int
}

func (RepoLanguage) TableName() string {
	return "repo_languages"
}

func (r *Repo) ScopeId() string {
	return r.Id
}

func (r *Repo) ScopeName() string {
	return r.Name
}
