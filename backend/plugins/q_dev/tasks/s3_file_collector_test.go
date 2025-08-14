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

package tasks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeS3Prefix(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"prefix", "prefix/"},
		{"prefix/", "prefix/"},
		{"path/to/folder", "path/to/folder/"},
		{"path/to/folder/", "path/to/folder/"},
	}

	for _, test := range tests {
		result := normalizeS3Prefix(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestIsCSVFile(t *testing.T) {
	tests := []struct {
		key      string
		expected bool
	}{
		{"file.csv", true},
		{"data.CSV", false},
		{"report.csv", true},
		{"document.txt", false},
		{"path/to/file.csv", true},
		{"file.csv.backup", false},
		{"", false},
	}

	for _, test := range tests {
		result := isCSVFile(test.key)
		assert.Equal(t, test.expected, result)
	}
}
