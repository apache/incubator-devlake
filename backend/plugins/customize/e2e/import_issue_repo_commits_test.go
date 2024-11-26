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

package e2e

import (
	"os"
	"testing"

	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/customize/impl"
	"github.com/apache/incubator-devlake/plugins/customize/service"
)

func TestImportIssueRepoCommitDataFlow(t *testing.T) {
	var plugin impl.Customize
	dataflowTester := e2ehelper.NewDataFlowTester(t, "customize", plugin)

	// create tables `issue_repo_commits` and `issue_commits`
	dataflowTester.FlushTabler(&crossdomain.IssueRepoCommit{})
	dataflowTester.FlushTabler(&crossdomain.IssueCommit{})
	dataflowTester.FlushTabler(&ticket.BoardIssue{})

	dataflowTester.ImportCsvIntoTabler("raw_tables/issue_repo_commits_original.csv", &crossdomain.IssueRepoCommit{})
	dataflowTester.ImportCsvIntoTabler("raw_tables/board_issues.csv", &ticket.BoardIssue{})

	svc := service.NewService(dataflowTester.Dal)

	issueRepoCommitsFile, err1 := os.Open("raw_tables/issue_repo_commits.csv")
	if err1 != nil {
		t.Fatal(err1)
	}
	defer issueRepoCommitsFile.Close()
	// import data
	err := svc.ImportIssueRepoCommit("csv-board", issueRepoCommitsFile, false)
	if err != nil {
		t.Fatal(err)
	}

	// import data incrementally
	issueRepoCommitsIncrementalFile, err2 := os.Open("raw_tables/issue_repo_commits_incremental.csv")
	if err2 != nil {
		t.Fatal(err2)
	}
	defer issueRepoCommitsIncrementalFile.Close()
	err = svc.ImportIssueRepoCommit("csv-board", issueRepoCommitsIncrementalFile, true)
	if err != nil {
		t.Fatal(err)
	}

	dataflowTester.VerifyTableWithRawData(
		crossdomain.IssueRepoCommit{},
		"snapshot_tables/issue_repo_commits.csv",
		[]string{
			"issue_id",
			"repo_url",
			"commit_sha",
			"host",
			"namespace",
			"repo_name",
		})
	dataflowTester.VerifyTableWithRawData(
		crossdomain.IssueCommit{},
		"snapshot_tables/issue_commits_from_import_issue_repo_commit.csv",
		[]string{
			"issue_id",
			"commit_sha",
		})
}
