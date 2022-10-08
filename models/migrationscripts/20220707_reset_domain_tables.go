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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/core"

	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
)

var _ core.MigrationScript = (*resetDomainTables)(nil)

type resetDomainTables struct{}

func (*resetDomainTables) Up(basicRes core.BasicRes) errors.Error {
	db := basicRes.GetDal()
	err := db.DropTables(
		"issue_assignee_history",
		"issue_status_history",
		"issue_sprints_history",
		"users",
		&archived.Repo{},
		&archived.Commit{},
		&archived.CommitParent{},
		&archived.PullRequest{},
		&archived.PullRequestCommit{},
		&archived.PullRequestComment{},
		&archived.PullRequestLabel{},
		&archived.Note{},
		&archived.RepoCommit{},
		&archived.Ref{},
		&archived.RefsCommitsDiff{},
		&archived.CommitFile{},
		&archived.Board{},
		&archived.Issue{},
		&archived.IssueLabel{},
		&archived.IssueComment{},
		&archived.BoardIssue{},
		&archived.BoardSprint{},
		&archived.IssueChangelogs{},
		&archived.Sprint{},
		&archived.SprintIssue{},
		&archived.IssueWorklog{},
		&archived.BoardRepo{},
		&archived.PullRequestIssue{},
		&archived.IssueCommit{},
		&archived.IssueRepoCommit{},
		&archived.RefsIssuesDiffs{},
		&archived.RefsPrCherrypick{},
	)
	if err != nil {
		return err
	}

	migrationHelper := migrationhelper.NewMigrationHelper(basicRes)
	return migrationHelper.AutoMigrateTables(
		&archived.Repo{},
		&archived.Commit{},
		&archived.CommitParent{},
		&archived.PullRequest{},
		&archived.PullRequestCommit{},
		&archived.PullRequestComment{},
		&archived.PullRequestLabel{},
		&archived.Note{},
		&archived.RepoCommit{},
		&archived.Ref{},
		&archived.RefsCommitsDiff{},
		&archived.CommitFile{},
		&archived.Board{},
		&archived.Issue{},
		&archived.IssueLabel{},
		&archived.IssueComment{},
		&archived.BoardIssue{},
		&archived.BoardSprint{},
		&archived.IssueChangelogs{},
		&archived.Sprint{},
		&archived.SprintIssue{},
		&archived.IssueWorklog{},
		&archived.BoardRepo{},
		&archived.PullRequestIssue{},
		&archived.IssueCommit{},
		&archived.IssueRepoCommit{},
		&archived.RefsIssuesDiffs{},
		&archived.RefsPrCherrypick{},
		&archived.Account{},
		&archived.User{},
		&archived.Team{},
		&archived.UserAccount{},
		&archived.TeamUser{},
	)
}

func (*resetDomainTables) Version() uint64 {
	return 20220707232344
}

func (*resetDomainTables) Owner() string {
	return "Framework"
}

func (*resetDomainTables) Name() string {
	return "create domain layer init schemas"
}
