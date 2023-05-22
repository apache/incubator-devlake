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

package api

import (
	"testing"
)

func Test_stripZeroByte(t *testing.T) {
	type foo struct {
		Field      string
		FromString string
		ToString   string
		From       *string
	}
	from := "Earth\u0000"
	type args struct {
		ifc interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_stripZeroByte",
			args: args{
				ifc: &foo{
					Field:      "home\u0000",
					FromString: "Earth",
					ToString:   "Mars\u0000",
					From:       &from,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stripZeroByte(tt.args.ifc)
			if tt.args.ifc.(*foo).Field != "home" {
				t.Errorf("stripZeroByte() = %v, want %v", tt.args.ifc.(*foo).Field, "home")
			}
			if tt.args.ifc.(*foo).FromString != "Earth" {
				t.Errorf("stripZeroByte() = %v, want %v", tt.args.ifc.(*foo).FromString, "Earth")
			}
			if tt.args.ifc.(*foo).ToString != "Mars" {
				t.Errorf("stripZeroByte() = %v, want %v", tt.args.ifc.(*foo).ToString, "Mars")
			}
			if *tt.args.ifc.(*foo).From != from {
				t.Errorf("stripZeroByte() = %v, want %v", *tt.args.ifc.(*foo).From, "Earth")
			}
		})
	}
}
