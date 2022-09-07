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
	"fmt"
	"github.com/apache/incubator-devlake/plugins/gitea/models/migrationscripts/archived"

	"gorm.io/gorm"
)

type addInitTables struct{}

func (*addInitTables) Up(ctx context.Context, db *gorm.DB) error {
	rawTableList := []string{
		"_raw_gitea_api_commit",
		"_raw_gitea_api_issues",
		"_raw_gitea_api_repo",
		"_raw_gitea_api_comments",
		"_raw_gitea_api_commits",
		"_raw_gitea_issue_comments",
	}
	for _, v := range rawTableList {
		err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", v)).Error
		if err != nil {
			return err
		}
	}

	err := db.Migrator().DropTable(
		&archived.GiteaRepo{},
		&archived.GiteaCommit{},
		&archived.GiteaRepoCommit{},
		&archived.GiteaIssue{},
		&archived.GiteaIssueComment{},
		&archived.GiteaCommitStat{},
		&archived.GiteaIssueLabel{},
		&archived.GiteaReviewer{},
		&archived.GiteaConnection{},
	)

	if err != nil {
		return err
	}

	err = db.Migrator().AutoMigrate(
		&archived.GiteaRepo{},
		&archived.GiteaCommit{},
		&archived.GiteaRepoCommit{},
		&archived.GiteaAccount{},
		&archived.GiteaIssue{},
		&archived.GiteaIssueComment{},
		&archived.GiteaCommitStat{},
		&archived.GiteaIssueLabel{},
		&archived.GiteaReviewer{},
		&archived.GiteaConnection{},
	)

	if err != nil {
		return err
	}

	return nil
}

func (*addInitTables) Version() uint64 {
	return 20220830163407
}

func (*addInitTables) Name() string {
	return "Gitea init schemas"
}
