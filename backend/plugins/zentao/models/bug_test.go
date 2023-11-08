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

package models

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestApiAccount_UnmarshalJSON(t *testing.T) {
	type bug struct {
		NotifyEmail string `json:"notifyEmail"`
		OpenedBy    *ApiAccount
	}

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    bug
		wantErr bool
	}{
		{
			"string",
			args{
				data: []byte(`{
								  "notifyEmail": "caitlyn.langosh@yahoo.com",
								  "openedBy": "admin"
								}`),
			},
			bug{
				NotifyEmail: "caitlyn.langosh@yahoo.com",
				OpenedBy:    &ApiAccount{Account: "admin"},
			},
			false,
		},
		{
			"empty string",
			args{
				data: []byte(`{
								  "notifyEmail": "caitlyn.langosh@yahoo.com",
								  "openedBy": ""
								}`),
			},
			bug{
				NotifyEmail: "caitlyn.langosh@yahoo.com",
				OpenedBy:    &ApiAccount{},
			},
			false,
		},
		{
			"struct",
			args{
				data: []byte(`{
								  "notifyEmail": "caitlyn.langosh@yahoo.com",
								  "openedBy": {
									"id": 1,
									"account": "admin",
									"avatar": "https://example.com/avatar.png",
									"realname": "root"
								  }
								}`),
			},
			bug{
				NotifyEmail: "caitlyn.langosh@yahoo.com",
				OpenedBy: &ApiAccount{
					ID:       1,
					Account:  "admin",
					Avatar:   "https://example.com/avatar.png",
					Realname: "root",
				},
			},
			false,
		},
		{
			"null",
			args{
				data: []byte(`{
								  "notifyEmail": "caitlyn.langosh@yahoo.com",
								  "openedBy": null
								}`),
			},
			bug{
				NotifyEmail: "caitlyn.langosh@yahoo.com",
				OpenedBy:    nil,
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dst bug
			if err := json.Unmarshal(tt.args.data, &dst); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(dst, tt.want) {
				t.Errorf("UnmarshalJSON() got = %v, want %v", dst, tt.want)
			}
			t.Logf("%+v\n", dst.OpenedBy)
		})
	}
}
