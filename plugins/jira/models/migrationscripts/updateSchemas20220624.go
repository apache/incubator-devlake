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

	"gorm.io/gorm"
)

type JiraConnection20220624 struct {
	EpicKeyField               string `gorm:"type:varchar(50);" json:"epicKeyField"`
	StoryPointField            string `gorm:"type:varchar(50);" json:"storyPointField"`
	RemotelinkCommitShaPattern string `gorm:"type:varchar(255);comment='golang regexp, the first group will be recognized as commit sha, ref https://github.com/google/re2/wiki/Syntax'" json:"remotelinkCommitShaPattern"`
}

func (JiraConnection20220624) TableName() string {
	return "_tool_jira_connections"
}

type UpdateSchemas20220624 struct {
}

func (u *UpdateSchemas20220624) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().DropColumn(&JiraConnection20220624{}, "epic_key_field")
	if err != nil {
		return err
	}
	err = db.Migrator().DropColumn(&JiraConnection20220624{}, "story_point_field")
	if err != nil {
		return err
	}
	return db.Migrator().DropColumn(&JiraConnection20220624{}, "remotelink_commit_sha_pattern")
}

func (*UpdateSchemas20220624) Version() uint64 {
	return 20220624102636
}

func (*UpdateSchemas20220624) Name() string {
	return "remove epic_key_field, story_point_field, remotelink_commit_sha_pattern"
}
