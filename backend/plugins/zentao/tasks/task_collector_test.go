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

package tasks

import (
	"github.com/apache/incubator-devlake/plugins/zentao/models"
	"reflect"
	"testing"
)

func Test_extractChildren(t *testing.T) {
	type args struct {
		allTaskRecords map[int64]models.ZentaoTaskRes
	}

	id1 := models.ZentaoTaskRes{Id: 1}
	id2 := models.ZentaoTaskRes{Id: 2}
	id3 := models.ZentaoTaskRes{Id: 3}
	id4 := models.ZentaoTaskRes{Id: 4}
	id5 := models.ZentaoTaskRes{Id: 5}
	id6 := models.ZentaoTaskRes{Id: 6}
	id1WithChildren3And4 := models.ZentaoTaskRes{
		Id:       1,
		Children: []*models.ZentaoTaskRes{&id3, &id4},
	}

	id2WithChildren3And4And5 := models.ZentaoTaskRes{
		Id:       2,
		Children: []*models.ZentaoTaskRes{&id3, &id4, &id5},
	}

	id3WithChildren3And4And6AndId1WithChildren3And4 := models.ZentaoTaskRes{
		Id:       3,
		Children: []*models.ZentaoTaskRes{&id3, &id4, &id6, &id1WithChildren3And4},
	}

	//id6 := models.ZentaoTaskRes{Id: 6}
	tests := []struct {
		name    string
		args    args
		want    map[int64]models.ZentaoTaskRes
		wantErr bool
	}{
		{
			name:    "s-0",
			args:    args{allTaskRecords: map[int64]models.ZentaoTaskRes{}},
			want:    map[int64]models.ZentaoTaskRes{},
			wantErr: false,
		},
		{
			name: "s-0-1",
			args: args{allTaskRecords: map[int64]models.ZentaoTaskRes{
				1: id1,
			}},
			want: map[int64]models.ZentaoTaskRes{
				1: id1,
			},
			wantErr: false,
		},
		{
			name: "s-1",
			args: args{allTaskRecords: map[int64]models.ZentaoTaskRes{
				1: id1,
				2: id2,
			}},
			want: map[int64]models.ZentaoTaskRes{
				1: id1,
				2: id2,
			},
			wantErr: false,
		},
		{
			name: "s-2",
			args: args{allTaskRecords: map[int64]models.ZentaoTaskRes{
				1: id1WithChildren3And4,
				2: id2,
			}},
			want: map[int64]models.ZentaoTaskRes{
				1: id1WithChildren3And4,
				2: id2,
				3: id3,
				4: id4,
			},
			wantErr: false,
		},
		{
			name: "s-3",
			args: args{allTaskRecords: map[int64]models.ZentaoTaskRes{
				1: id1WithChildren3And4,
				2: id2WithChildren3And4And5,
			}},
			want: map[int64]models.ZentaoTaskRes{
				1: id1WithChildren3And4,
				2: id2WithChildren3And4And5,
				3: id3,
				4: id4,
				5: id5,
			},
			wantErr: false,
		},
		{
			name: "s-4",
			args: args{allTaskRecords: map[int64]models.ZentaoTaskRes{
				1: id1WithChildren3And4,
				2: id2WithChildren3And4And5,
				3: id3WithChildren3And4And6AndId1WithChildren3And4,
			}},
			want: map[int64]models.ZentaoTaskRes{
				1: id1WithChildren3And4,
				2: id2WithChildren3And4And5,
				3: id3WithChildren3And4And6AndId1WithChildren3And4,
				4: id4,
				5: id5,
				6: id6,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractChildren(tt.args.allTaskRecords)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractChildren() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractChildren() got = %v, want %v", got, tt.want)
				t.Log("==========got vs want==========")
				for k, v := range got {
					var find bool
					for wantk, wantv := range tt.want {
						if k == wantk {
							find = true
							if !reflect.DeepEqual(v, wantv) {
								t.Log("\n", "got", k, v, "\n", "want", wantk, wantv)
							}
						}
					}
					if !find {
						t.Log("not found", "got", k, v)
					}
				}
			}
		})
	}
}

func Test_extractChildrenWithDFS(t *testing.T) {
	type args struct {
		task models.ZentaoTaskRes
	}

	id1 := models.ZentaoTaskRes{Id: 1}
	id2 := models.ZentaoTaskRes{Id: 2}
	id3 := models.ZentaoTaskRes{Id: 3}
	id4 := models.ZentaoTaskRes{Id: 4}
	id5 := models.ZentaoTaskRes{Id: 5}
	id6 := models.ZentaoTaskRes{Id: 6}
	id7 := models.ZentaoTaskRes{Id: 7}
	id1WithChildren3And4 := models.ZentaoTaskRes{
		Id:       1,
		Children: []*models.ZentaoTaskRes{&id3, &id4},
	}
	id1WithChildren7 := models.ZentaoTaskRes{
		Id:       1,
		Children: []*models.ZentaoTaskRes{&id7},
	}

	id2WithChildren3And4And5 := models.ZentaoTaskRes{
		Id:       2,
		Children: []*models.ZentaoTaskRes{&id3, &id4, &id5},
	}

	id3WithChildren4And6AndId1WithChildren3And4 := models.ZentaoTaskRes{
		Id:       3,
		Children: []*models.ZentaoTaskRes{&id4, &id6, &id1WithChildren7},
	}

	tests := []struct {
		name    string
		args    args
		want    []models.ZentaoTaskRes
		wantErr bool
	}{
		{
			name:    "0",
			args:    args{task: id2},
			want:    []models.ZentaoTaskRes{id2},
			wantErr: false,
		},
		{
			name:    "1",
			args:    args{task: id1},
			want:    []models.ZentaoTaskRes{id1},
			wantErr: false,
		},
		{
			name:    "2",
			args:    args{task: id1WithChildren3And4},
			want:    []models.ZentaoTaskRes{id1WithChildren3And4, id3, id4},
			wantErr: false,
		},
		{
			name:    "3",
			args:    args{task: id2WithChildren3And4And5},
			want:    []models.ZentaoTaskRes{id2WithChildren3And4And5, id3, id4, id5},
			wantErr: false,
		},
		{
			name:    "4",
			args:    args{task: id3WithChildren4And6AndId1WithChildren3And4},
			want:    []models.ZentaoTaskRes{id3WithChildren4And6AndId1WithChildren3And4, id7, id4, id6, id1WithChildren7},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractChildrenWithDFS(tt.args.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractChildrenWithDFS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, v := range got {
				var find bool
				for _, wantv := range tt.want {
					if v.Id == wantv.Id {
						find = true
						if !reflect.DeepEqual(v, wantv) {
							t.Errorf("got: %v, want: %v", v, wantv)
						}
					}
				}
				if !find {
					t.Errorf("not found: %v", v)
				}
			}
		})
	}
}
