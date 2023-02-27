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

package service

import (
	"regexp"
	"strings"
	"testing"
)

func TestService_checkFieldName(t *testing.T) {
	nameChecker := regexp.MustCompile(`^x_[a-zA-Z0-9_]{0,50}$`)
	tests := []struct {
		name string
		args string
		want bool
	}{
		{
			"",
			"x_abc23_e",
			true,
		},
		{
			"",
			"_abc23_e",
			false,
		},
		{
			"",
			"x__",
			true,
		},
		{
			"issue #4519",
			"x_" + strings.Repeat("a", 50),
			true,
		},
		{
			"issue #4519",
			"x_" + strings.Repeat("a", 51),
			false,
		},
		{
			"",
			"x_ space",
			false,
		},
		{
			"",
			"x_123",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := nameChecker.MatchString(tt.args)
			if got != tt.want {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}
