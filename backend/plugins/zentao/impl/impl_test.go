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

import "testing"

func Test_getZentaoWebURL(t *testing.T) {
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
			name:    "without-zentao",
			args:    args{endpoint: "http://54.158.1.10:30001/api.php/v1/"},
			want:    "http://54.158.1.10:30001",
			wantErr: false,
		},
		{
			name:    "with-zentao",
			args:    args{endpoint: "http://54.158.1.10:30001/zentao/api.php/v1/"},
			want:    "http://54.158.1.10:30001/zentao",
			wantErr: false,
		},
		{
			name:    "with-others",
			args:    args{endpoint: "http://54.158.1.10:30001/abc/api.php/v1/"},
			want:    "http://54.158.1.10:30001/abc",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getZentaoHomePage(tt.args.endpoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("getZentaoHomePage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getZentaoHomePage() got = %v, want %v", got, tt.want)
			}
		})
	}
}
