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

package runner

import (
	"net/url"
	"testing"
)

func Test_addLocal(t *testing.T) {
	values, _ := url.ParseQuery("charset=utf8mb4&parseTime=True")
	type args struct {
		query url.Values
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test add local",
			args: args{
				query: values,
			},
			want: "charset=utf8mb4&loc=Local&parseTime=True",
		},
		{
			name: "test add local",
			args: args{
				query: url.Values{
					"local": []string{"abc"},
				},
			},
			want: "loc=Local&local=abc",
		},
		{
			name: "test add local",
			args: args{
				query: url.Values{
					"loc": []string{"abc"},
				},
			},
			want: "loc=abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addLocal(tt.args.query); got != tt.want {
				t.Errorf("addLocal() = %v, want %v", got, tt.want)
			}
		})
	}
}
