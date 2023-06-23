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

package subtaskmeta_sorter

import (
	"reflect"
	"testing"

	"github.com/apache/incubator-devlake/core/plugin"
)

func Test_dependencyTableTopologicalSort(t *testing.T) {
	pluginA := plugin.SubTaskMeta{
		Name:             string(prefixCollect) + "A",
		DependencyTables: []string{"Table1"},
	}
	pluginB := plugin.SubTaskMeta{
		Name:             string(prefixCollect) + "B",
		DependencyTables: []string{"table2"},
	}
	pluginC := plugin.SubTaskMeta{
		Name:             string(prefixCollect) + "C",
		DependencyTables: []string{"table1", "table2"},
	}
	pluginD := plugin.SubTaskMeta{
		Name:             string(prefixCollect) + "D",
		DependencyTables: []string{"table1", "table2"},
	}
	type args struct {
		metas []*plugin.SubTaskMeta
	}
	tests := []struct {
		name    string
		args    args
		want    []plugin.SubTaskMeta
		wantErr bool
	}{
		{
			name: "correct stable sort",
			args: args{
				metas: []*plugin.SubTaskMeta{
					&pluginA, &pluginB, &pluginC,
				},
			},
			want:    []plugin.SubTaskMeta{pluginA, pluginB, pluginC},
			wantErr: false,
		},
		{
			name: "cycle error",
			args: args{
				metas: []*plugin.SubTaskMeta{
					&pluginC, &pluginD,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dependencyTableTopologicalSort(tt.args.metas)
			if (err != nil) != tt.wantErr {
				t.Errorf("dependencyTableTopologicalSort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dependencyTableTopologicalSort() got = %v, want %v", got, tt.want)
			}
		})
	}
}
