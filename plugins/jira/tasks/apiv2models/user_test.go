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

package apiv2models

import "testing"

func TestUser_getAccountId(t *testing.T) {
	type fields struct {
		Self         string
		Key          string
		Name         string
		EmailAddress string
		AccountId    string
		AccountType  string
		AvatarUrls   struct {
			Four8X48  string `json:"48x48"`
			Two4X24   string `json:"24x24"`
			One6X16   string `json:"16x16"`
			Three2X32 string `json:"32x32"`
		}
		DisplayName string
		Active      bool
		Deleted     bool
		TimeZone    string
		Locale      string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"",
			fields{EmailAddress: "abc"},
			"abc",
		},
		{"",
			fields{EmailAddress: "abc", AccountId: "abcd"},
			"abcd",
		},
		{"",
			fields{AccountId: "abcd"},
			"abcd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Account{
				Self:         tt.fields.Self,
				Key:          tt.fields.Key,
				Name:         tt.fields.Name,
				EmailAddress: tt.fields.EmailAddress,
				AccountId:    tt.fields.AccountId,
				AccountType:  tt.fields.AccountType,
				AvatarUrls:   tt.fields.AvatarUrls,
				DisplayName:  tt.fields.DisplayName,
				Active:       tt.fields.Active,
				Deleted:      tt.fields.Deleted,
				TimeZone:     tt.fields.TimeZone,
				Locale:       tt.fields.Locale,
			}
			if got := u.getAccountId(); got != tt.want {
				t.Errorf("getAccountId() = %v, want %v", got, tt.want)
			}
		})
	}
}
