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

package impl

import (
	"testing"

	"github.com/apache/incubator-devlake/helpers/pluginhelper/subtaskmeta/sorter"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
)

func Test_genSubtaskList(t *testing.T) {
	testSortedSubtaskList, err := sorter.NewTableSorter(tasks.SubTaskMetaList).Sort()
	if err != nil {
		t.Error(err)
		return
	}
	subtaskNameList := make([]string, len(testSortedSubtaskList))
	for index, item := range testSortedSubtaskList {
		subtaskNameList[index] = item.Name
	}
	t.Logf("got subtask list %s", subtaskNameList)
}
