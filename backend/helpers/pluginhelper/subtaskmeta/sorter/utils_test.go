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
)

func Test_topologicalSortSameElements(t *testing.T) {
	type args struct {
		dependenciesMap map[string][]string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "correct stable sort",
			args: args{
				map[string][]string{
					"Aa": {"B", "C"},
					"Ac": {"B", "C"},
					"B":  {"C"},
					"C":  {},
				},
			},
			want:    []string{"C", "B", "Aa", "Ac"},
			wantErr: false,
		},
		{
			name: "cyclic error",
			args: args{
				map[string][]string{
					"A": {"B", "C"},
					"B": {"C"},
					"C": {"A"},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := topologicalSortSameElements(tt.args.dependenciesMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("topologicalSortSameElements() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("topologicalSortSameElements() got = %v, want %v", got, tt.want)
			}
		})
	}
}
