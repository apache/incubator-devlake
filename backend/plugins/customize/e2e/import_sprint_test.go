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

	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/customize/impl"
	"github.com/apache/incubator-devlake/plugins/customize/service"
)

func TestImportSprintDataFlow(t *testing.T) {
	var plugin impl.Customize
	dataflowTester := e2ehelper.NewDataFlowTester(t, "customize", plugin)

	// 创建表
	dataflowTester.FlushTabler(&ticket.Sprint{})
	dataflowTester.FlushTabler(&ticket.BoardSprint{})

	// 导入必要数据
	svc := service.NewService(dataflowTester.Dal)

	// 导入全量数据
	sprintFile, err := os.Open("raw_tables/sprints.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer sprintFile.Close()
	err = svc.ImportSprint("csv-board", sprintFile, false)
	if err != nil {
		t.Fatal(err)
	}

	// 导入增量数据
	sprintIncrementalFile, err := os.Open("raw_tables/sprints_incremental.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer sprintIncrementalFile.Close()
	err = svc.ImportSprint("csv-board", sprintIncrementalFile, true)
	if err != nil {
		t.Fatal(err)
	}

	// 验证结果
	dataflowTester.VerifyTableWithRawData(
		ticket.Sprint{},
		"snapshot_tables/sprints.csv",
		[]string{
			"id",
			"url",
			"status",
			"name",
			"started_date",
			"ended_date",
			"completed_date",
			"original_board_id",
		})

	dataflowTester.VerifyTableWithRawData(
		ticket.BoardSprint{},
		"snapshot_tables/board_sprints.csv",
		[]string{
			"board_id",
			"sprint_id",
		})
}
