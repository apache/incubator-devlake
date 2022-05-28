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

type GithubIssue20220524 struct {
	AuthorId   int
	AuthorName string `gorm:"type:varchar(255)"`
}

func (GithubIssue20220524) TableName() string {
	return "_tool_github_issues"
}

type UpdateSchemas20220524 struct{}

func (*UpdateSchemas20220524) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().AddColumn(GithubIssue20220524{}, "author_id")
	if err != nil {
		return err
	}
	err = db.Migrator().AddColumn(GithubIssue20220524{}, "author_name")
	if err != nil {
		return err
	}
	return nil
}

func (*UpdateSchemas20220524) Version() uint64 {
	return 20220524000002
}

func (*UpdateSchemas20220524) Name() string {
	return "Add column `author_id`/`author_name` in `GithubIssue`"
}
