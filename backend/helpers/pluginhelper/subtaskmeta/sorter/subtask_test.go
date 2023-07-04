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
	"reflect"
	"testing"

	"github.com/apache/incubator-devlake/core/plugin"
)

func Test_topologicalSort(t *testing.T) {
	pluginA := plugin.SubTaskMeta{
		Name:         "A",
		Dependencies: []*plugin.SubTaskMeta{},
	}
	pluginB := plugin.SubTaskMeta{
		Name:         "B",
		Dependencies: []*plugin.SubTaskMeta{&pluginA},
	}
	pluginC := plugin.SubTaskMeta{
		Name:         "C",
		Dependencies: []*plugin.SubTaskMeta{&pluginA},
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
			name: "correct stable order",
			args: args{
				metas: []*plugin.SubTaskMeta{&pluginA, &pluginC, &pluginB},
			},
			want: []plugin.SubTaskMeta{
				pluginA, pluginB, pluginC,
			},
			wantErr: false,
		},
		{
			name: "duplicate error",
			args: args{
				metas: []*plugin.SubTaskMeta{&pluginA, &pluginA, &pluginB},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "cycle error",
			args: args{[]*plugin.SubTaskMeta{
				{
					Name: "D",
					Dependencies: []*plugin.SubTaskMeta{{
						Name: "E",
					}},
				},
				{
					Name: "E",
					Dependencies: []*plugin.SubTaskMeta{{
						Name: "D",
					}},
				},
			}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dependenciesTopologicalSort(tt.args.metas)
			if (err != nil) != tt.wantErr {
				t.Errorf("dependenciesTopologicalSort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dependenciesTopologicalSort() got = %v, want %v", got, tt.want)
			}
		})
	}
}
