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

package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeAndDecode(t *testing.T) {
	TestStr := "The string for testing"
	var err error

	var TestEncode string
	var TestDecode string

	encKey := RandomEncKey()
	// encryption test
	TestEncode, err = Encrypt(encKey, TestStr)
	assert.Empty(t, err)

	// decrypt test
	TestDecode, err = Decrypt(encKey, TestEncode)
	assert.Empty(t, err)

	// Verify decryption result
	assert.Equal(t, string(TestDecode), TestStr)
}

func TestRandomCapsStr(t *testing.T) {
	type args struct {
		len int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"",
			args{128},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(RandomCapsStr(tt.args.len))
		})
	}
}

func TestEncode(t *testing.T) {
	encKey := RandomEncKey()
	type args struct {
		Input string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"",
			args{"bGlhbmcuemhhbmdAbWVyaWNvLmRldjprYUU2eWpNY1VYV2FCNUhIS3BGRkQ1RTg="},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encrypt(encKey, tt.args.Input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}
