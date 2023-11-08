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

package sorter

import (
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"reflect"
	"testing"
)

func Test_tableTopologicalSort(t *testing.T) {
	pluginA := plugin.SubTaskMeta{
		Name:             "A",
		DependencyTables: []string{},
		ProductTables:    []string{"_TOOL_TEST_TABLE", "_TOOL_TEST_TABLE2"},
	}
	pluginB := plugin.SubTaskMeta{
		Name:             "B",
		DependencyTables: []string{"_TOOL_TEST_TABLE"},
		ProductTables:    []string{"_TOOL_TEST_TABLE2", "_TOOL_TEST_TABLE3"},
	}
	pluginC := plugin.SubTaskMeta{
		Name:             "C",
		DependencyTables: []string{"_TOOL_TEST_TABLE"},
		ProductTables:    []string{"_TOOL_TEST_TABLE3"},
	}
	pluginD := plugin.SubTaskMeta{
		Name:             "D",
		DependencyTables: []string{"_TOOL_TEST_TABLE2", "_TOOL_TEST_TABLE3"},
		ProductTables:    []string{"_TOOL_TEST_TABLE4"},
	}

	type args struct {
		metas []*plugin.SubTaskMeta
	}
	tests := []struct {
		name  string
		args  args
		want  []plugin.SubTaskMeta
		want1 errors.Error
	}{
		{
			name:  "test sorter",
			args:  args{[]*plugin.SubTaskMeta{&pluginA, &pluginB, &pluginC, &pluginD}},
			want:  []plugin.SubTaskMeta{pluginA, pluginB, pluginC, pluginD},
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tableTopologicalSort(tt.args.metas)
			if len(got) != len(tt.want) {
				t.Errorf("tableTopologicalSort() got = %v, want %v", got, tt.want)
			}
			for index, item := range got {
				if item.Name != tt.want[index].Name {
					t.Errorf("tableTopologicalSort() got = %v, want %v, not equal with index = %d", got, tt.want, index)
				}
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("tableTopologicalSort() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
