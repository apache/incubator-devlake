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
	"github.com/apache/incubator-devlake/plugins/tapd/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type InitSchemas struct{}

func (*InitSchemas) Up(ctx context.Context, db *gorm.DB) error {
	return db.Migrator().AutoMigrate(
		&archived.TapdWorkspace{},
		&archived.TapdWorklog{},
		&archived.TapdWorkspaceIteration{},
		&archived.TapdUser{},
		&archived.TapdBugChangelog{},
		&archived.TapdBugChangelogItem{},
		&archived.TapdStoryChangelog{},
		&archived.TapdStoryChangelogItem{},
		&archived.TapdTaskChangelog{},
		&archived.TapdTaskChangelogItem{},
		&archived.TapdIssue{},
		&archived.TapdIteration{},
		&archived.TapdConnection{},
		&archived.TapdBug{},
		&archived.TapdStory{},
		&archived.TapdTask{},
		&archived.TapdTaskLabel{},
		&archived.TapdBugLabel{},
		&archived.TapdStoryLabel{},
		&archived.TapdBugStatus{},
		&archived.TapdStoryStatus{},
		&archived.TapdBugCommit{},
		&archived.TapdStoryCommit{},
		&archived.TapdTaskCommit{},
		&archived.TapdWorkSpaceBug{},
		&archived.TapdWorkSpaceStory{},
		&archived.TapdWorkSpaceTask{},
		&archived.TapdIterationBug{},
		&archived.TapdIterationStory{},
		&archived.TapdIterationTask{},
		&archived.TapdStoryCustomFields{},
		&archived.TapdBugCustomFields{},
		&archived.TapdTaskCustomFields{},
	)
}

func (*InitSchemas) Version() uint64 {
	return 20220420231138
}

func (*InitSchemas) Name() string {
	return "Tapd init schemas"
}
