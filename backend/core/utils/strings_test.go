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

package utils

import (
	"testing"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestRandLetterBytes(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name  string
		args  args
		want1 errors.Error
	}{
		{
			"test1",
			args{0},
			nil,
		},
		{
			"test1",
			args{-1},
			errors.Default.New("n must be greater than 0"),
		},
		{
			"test2",
			args{10},
			nil,
		},
		{
			"test3",
			args{128},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := RandLetterBytes(tt.args.n)
			t.Log(got)
			assert.Equalf(t, tt.want1, got1, "RandLetterBytes(%v)", tt.args.n)
		})
	}
}

func TestSanitizeString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test-1",
			args: args{s: ""},
			want: "",
		},
		{
			name: "test-2",
			args: args{s: "s"},
			want: "*",
		},
		{
			name: "test-3",
			args: args{s: "ss"},
			want: "**",
		},
		{
			name: "test-4",
			args: args{s: "s1s"},
			want: "s*s",
		},
		{
			name: "test-5",
			args: args{s: "s12s"},
			want: "s**s",
		},
		{
			name: "test-6",
			args: args{s: "s123s"},
			want: "s***s",
		},
		{
			name: "test-7",
			args: args{s: "s1234s"},
			want: "s1**4s",
		},
		{
			name: "test-8",
			args: args{s: "s123456789s"},
			want: "s1*******9s",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, SanitizeString(tt.args.s), "SanitizeString(%v)", tt.args.s)
		})
	}
}

func TestSubstr(t *testing.T) {
	assert.Equalf(t, "😂😂😂", Substr("😂😂😂a", 0, 3), "substr on utf8")
	assert.Equalf(t, "c", Substr("😂😂c😂", 2, 1), "substr on lattin + unicode")
	assert.Equalf(t, "abcde", Substr("abcde", 0, 100), "specified length greater than actual length")
	assert.Equalf(t, "", Substr("abcde", 100, 100), "specified start greater than actual length")
}
