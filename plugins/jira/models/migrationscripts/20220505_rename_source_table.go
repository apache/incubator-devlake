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
	"github.com/apache/incubator-devlake/plugins/jira/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type JiraConnectionV010 struct {
	common.Model
	Name                       string `gorm:"type:varchar(100);uniqueIndex" json:"name" validate:"required"`
	Endpoint                   string `json:"endpoint" validate:"required"`
	BasicAuthEncoded           string `json:"basicAuthEncoded" validate:"required"`
	EpicKeyField               string `gorm:"type:varchar(50);" json:"epicKeyField"`
	StoryPointField            string `gorm:"type:varchar(50);" json:"storyPointField"`
	RemotelinkCommitShaPattern string `gorm:"type:varchar(255);comment='golang regexp, the first group will be recognized as commit sha, ref https://github.com/google/re2/wiki/Syntax'" json:"remotelinkCommitShaPattern"`
	Proxy                      string `json:"proxy"`
	RateLimit                  int    `comment:"api request rate limt per hour" json:"rateLimit"`
}

func (JiraConnectionV010) TableName() string {
	return "_tool_jira_connections"
}

type renameSourceTable struct{}

func (*renameSourceTable) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().RenameTable(archived.JiraSource{}, JiraConnectionV010{})
	if err != nil {
		return err
	}
	return nil
}

func (*renameSourceTable) Version() uint64 {
	return 20220505212344
}

func (*renameSourceTable) Owner() string {
	return "Jira"
}

func (*renameSourceTable) Name() string {
	return "Rename source to connection "
}
