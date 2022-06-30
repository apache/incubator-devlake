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

type Issue20220630 struct {
	StoryPoint    int64
	TmpStoryPoint uint
}

func (Issue20220630) TableName() string {
	return "issues"
}

type UpdateSchemas20220630 struct {
}

func (u *UpdateSchemas20220630) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().RenameColumn(&Issue20220630{}, "story_point", "tmp_story_point")
	if err != nil {
		return err
	}
	err = db.Migrator().AddColumn(&Issue20220630{}, "story_point")
	if err != nil {
		return err
	}
	err = db.Model(&Issue20220630{}).Where("1 = 1").UpdateColumn("story_point", gorm.Expr("tmp_story_point")).Error
	if err != nil {
		return err
	}
	return db.Migrator().DropColumn(&Issue20220630{}, "tmp_story_point")
}

func (*UpdateSchemas20220630) Version() uint64 {
	return 20220630131508
}

func (*UpdateSchemas20220630) Name() string {
	return "alter story_point from unsigned to signed"
}
