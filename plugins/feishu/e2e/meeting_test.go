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
	"github.com/apache/incubator-devlake/plugins/feishu/models"
	"testing"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/feishu/impl"
	"github.com/apache/incubator-devlake/plugins/feishu/tasks"
)

func TestEventDataFlow(t *testing.T) {
	var plugin impl.Feishu
	dataflowTester := e2ehelper.NewDataFlowTester(t, "gitlab", plugin)

	taskData := &tasks.FeishuTaskData{
		Options: &tasks.FeishuOptions{
			ConnectionId: 1,
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_feishu_meeting_top_user_item.csv", "_raw_feishu_meeting_top_user_item")

	// verify extraction
	dataflowTester.FlushTabler(&models.FeishuMeetingTopUserItem{})
	dataflowTester.Subtask(tasks.ExtractMeetingTopUserItemMeta, taskData)
	dataflowTester.VerifyTable(
		models.FeishuMeetingTopUserItem{},
		"./snapshot_tables/_tool_feishu_meeting_top_user_items.csv",
		[]string{"connection_id", "start_time", "name"},
		[]string{
			"meeting_count",
			"meeting_duration",
			"user_type",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
}
