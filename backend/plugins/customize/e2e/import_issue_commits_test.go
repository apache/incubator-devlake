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
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/customize/impl"
	"github.com/apache/incubator-devlake/plugins/customize/service"
)

func TestImportIssueCommitDataFlow(t *testing.T) {
	var plugin impl.Customize
	dataflowTester := e2ehelper.NewDataFlowTester(t, "customize", plugin)

	// create table `issue_commits`
	dataflowTester.FlushTabler(&crossdomain.IssueCommit{})
	svc := service.NewService(dataflowTester.Dal)

	f, err1 := os.Open("raw_tables/issues_commits.csv")
	if err1 != nil {
		t.Fatal(err1)
	}
	defer f.Close()
	// import data
	err := svc.ImportIssueCommit(`{"ConnectionId":1,"BoardId":8}`, f)
	if err != nil {
		t.Fatal(err)
	}
	dataflowTester.VerifyTableWithRawData(
		crossdomain.IssueCommit{},
		"snapshot_tables/issue_commits.csv",
		[]string{
			"issue_id",
			"commit_sha",
		})
}
