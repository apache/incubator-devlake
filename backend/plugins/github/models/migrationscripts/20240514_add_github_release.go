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
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"time"
)

type addReleaseTable struct{}

type release20240514 struct {
	archived.NoPKModel `json:"-" mapstructure:"-"`
	ConnectionId       uint64    `json:"connection_id" gorm:"primaryKey"`
	GithubId           int       `json:"github_id"`
	Id                 string    `json:"id" gorm:"type:varchar(255);primaryKey"`
	AuthorName         string    `json:"authorName"`
	AuthorID           string    `json:"authorId"`
	CreatedAt          time.Time `json:"createdAt"`
	DatabaseID         int       `json:"databaseId"`
	Description        string    `json:"description"`
	DescriptionHTML    string    `json:"descriptionHTML"`
	IsDraft            bool      `json:"isDraft"`
	IsLatest           bool      `json:"isLatest"`
	IsPrerelease       bool      `json:"isPrerelease"`
	Name               string    `json:"name"`
	PublishedAt        time.Time `json:"publishedAt"`
	ResourcePath       string    `json:"resourcePath"`
	TagName            string    `json:"tagName"`
	UpdatedAt          time.Time `json:"updatedAt"`
	URL                string    `json:"url"`
}

func (release20240514) TableName() string {
	return "_tool_github_releases"
}

func (u *addReleaseTable) Up(baseRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(baseRes, &release20240514{})
}

func (*addReleaseTable) Version() uint64 {
	return 20240514182300
}

func (*addReleaseTable) Name() string {
	return "add table _tool_github_releases"
}
