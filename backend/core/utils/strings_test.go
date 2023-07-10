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
			assert.Equalf(t, unwrap(tt.want1), unwrap(got1), "RandLetterBytes(%v)", tt.args.n)
		})
	}
}

func unwrap(err errors.Error) error {
	if err == nil {
		return nil
	}
	return err.Unwrap()
}
