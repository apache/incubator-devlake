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

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

// TestArrayColumnScan documents the behaviour of the scanner used for PostgreSQL
// array columns in export(). The PostgreSQL driver hands array columns to
// database/sql as the raw array literal (a string), so the destination has to be
// an sql.Scanner able to parse that literal.
func TestArrayColumnScan(t *testing.T) {
	cases := []struct {
		name     string
		literal  string
		expected []string
	}{
		{"simple", `{a,b}`, []string{"a", "b"}},
		{"empty", `{}`, []string{}},
		{"single empty element", `{""}`, []string{""}},
		{"quoted separators", `{"a,x","b\"y"}`, []string{"a,x", `b"y`}},
		{"utf8 and spaces", `{ä ö,"x y"}`, []string{"ä ö", "x y"}},
	}

	typeMap := pgtype.NewMap()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var arr []string
			scanner := typeMap.SQLScanner(&arr)

			assert.NoError(t, scanner.Scan(tc.literal))
			assert.Equal(t, tc.expected, arr)

			// database/sql may deliver the very same literal as []byte.
			var fromBytes []string
			assert.NoError(t, typeMap.SQLScanner(&fromBytes).Scan([]byte(tc.literal)))
			assert.Equal(t, tc.expected, fromBytes)
		})
	}
}

// TestArrayColumnScanNullElement pins down that NULL elements are still rejected:
// the previous implementation could not handle them either, so this is not a
// behaviour change.
func TestArrayColumnScanNullElement(t *testing.T) {
	var arr []string
	err := pgtype.NewMap().SQLScanner(&arr).Scan(`{NULL}`)
	assert.Error(t, err)
}
