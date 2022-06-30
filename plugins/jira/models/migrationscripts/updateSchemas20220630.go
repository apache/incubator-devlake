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

type JiraIssue20220630 struct {
	StdStoryPoint    int64
	TmpStdStoryPoint uint
}

func (JiraIssue20220630) TableName() string {
	return "_tool_jira_issues"
}

type UpdateSchemas20220630 struct {
}

func (u *UpdateSchemas20220630) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().RenameColumn(&JiraIssue20220630{}, "std_story_point", "tmp_std_story_point")
	if err != nil {
		return err
	}
	err = db.Migrator().AddColumn(&JiraIssue20220630{}, "std_story_point")
	if err != nil {
		return err
	}
	err = db.Model(&JiraIssue20220630{}).Where("1 = 1").UpdateColumn("std_story_point", gorm.Expr("tmp_std_story_point")).Error
	if err != nil {
		return err
	}
	return db.Migrator().DropColumn(&JiraIssue20220630{}, "tmp_std_story_point")
}

func (*UpdateSchemas20220630) Version() uint64 {
	return 20220630130656
}

func (*UpdateSchemas20220630) Name() string {
	return "alter std_story_point from unsigned to signed"
}
