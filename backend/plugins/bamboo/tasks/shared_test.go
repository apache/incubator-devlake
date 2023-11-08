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

import "testing"

func Test_getBambooWebURL(t *testing.T) {
	type args struct {
		endpoint string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "without-bamboo",
			args:    args{endpoint: "http://54.158.1.10:30001/rest/api/latest"},
			want:    "http://54.158.1.10:30001",
			wantErr: false,
		},
		{
			name:    "without-bamboo-1",
			args:    args{endpoint: "http://54.158.1.10:30001/rest/api/latest/"},
			want:    "http://54.158.1.10:30001",
			wantErr: false,
		},
		{
			name:    "with-bamboo",
			args:    args{endpoint: "http://54.158.1.10:30001/bamboo/rest/api/latest"},
			want:    "http://54.158.1.10:30001/bamboo",
			wantErr: false,
		},
		{
			name:    "with-bamboo-1",
			args:    args{endpoint: "http://54.158.1.10:30001/bamboo/rest/api/latest/"},
			want:    "http://54.158.1.10:30001/bamboo",
			wantErr: false,
		},
		{
			name:    "with-others",
			args:    args{endpoint: "http://54.158.1.10:30001/abc/rest/api/latest"},
			want:    "http://54.158.1.10:30001/abc",
			wantErr: false,
		},
		{
			name:    "with-others-1",
			args:    args{endpoint: "http://54.158.1.10:30001/abc/rest/api/latest/repos"},
			want:    "http://54.158.1.10:30001/abc",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBambooHomePage(tt.args.endpoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("getbambooHomePage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getbambooHomePage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateFakeRepoURL(t *testing.T) {
	type args struct {
		endpoint string
		repoID   int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "t-1",
			args: args{
				endpoint: "http://127.0.0.1:8080/abc",
				repoID:   123,
			},
			want:    "fake://127.0.0.1:8080/repos/123",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateFakeRepoUrl(tt.args.endpoint, tt.args.repoID)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateFakeRepoUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("generateFakeRepoUrl() got = %v, want %v", got, tt.want)
			}
		})
	}
}
