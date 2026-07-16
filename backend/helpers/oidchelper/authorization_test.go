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

package oidchelper

import "testing"

func TestIsUserAllowed(t *testing.T) {
	cases := []struct {
		name  string
		cfg   Config
		email string
		want  bool
	}{
		{
			name:  "no restrictions",
			cfg:   Config{},
			email: "user@example.com",
			want:  true,
		},
		{
			name: "allowed email",
			cfg: Config{
				AllowEmails: map[string]struct{}{
					"user@example.com": {},
				},
			},
			email: "user@example.com",
			want:  true,
		},
		{
			name: "blocked email",
			cfg: Config{
				AllowEmails: map[string]struct{}{
					"user@example.com": {},
				},
			},
			email: "other@example.com",
			want:  false,
		},
		{
			name: "allowed domain",
			cfg: Config{
				AllowDomains: map[string]struct{}{
					"example.com": {},
				},
			},
			email: "user@example.com",
			want:  true,
		},
		{
			name: "blocked domain",
			cfg: Config{
				AllowDomains: map[string]struct{}{
					"example.com": {},
				},
			},
			email: "user@other.com",
			want:  false,
		},
		{
			name: "email case insensitive",
			cfg: Config{
				AllowEmails: map[string]struct{}{
					"user@example.com": {},
				},
			},
			email: "USER@example.com",
			want:  true,
		},
		{
			name: "domain case insensitive",
			cfg: Config{
				AllowDomains: map[string]struct{}{
					"example.com": {},
				},
			},
			email: "user@EXAMPLE.COM",
			want:  true,
		},
		{
			name: "invalid email",
			cfg: Config{
				AllowDomains: map[string]struct{}{
					"example.com": {},
				},
			},
			email: "not-an-email",
			want:  false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.cfg.IsUserAllowed(tc.email); got != tc.want {
				t.Errorf(
					"IsUserAllowed(%q) = %v, want %v",
					tc.email,
					got,
					tc.want,
				)
			}
		})
	}
}
