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
	"context"
	"testing"

	"github.com/apache/devlake/core/dal"
	mockdal "github.com/apache/devlake/mocks/core/dal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// fetchedRow is one row the mocked cursor hands back, in the shape d.Fetch produces.
type fetchedRow map[string]interface{}

// newExtractorMocks wires a Dal whose cursor yields the given rows in order. It returns the Dal
// mock plus a pointer to a slice collecting every UPDATE statement executed, so a test can tell
// which rows actually made it through.
func newExtractorMocks(t *testing.T, rows []fetchedRow) (*mockdal.Dal, *[]string) {
	t.Helper()

	idColumn := new(mockdal.ColumnMeta)
	idColumn.On("Name").Return("id").Maybe()
	idColumn.On("PrimaryKey").Return(true, true).Maybe()

	cursor := new(mockdal.Rows)
	// Next reports true once per row, then false to end the scan.
	call := 0
	cursor.On("Next").Return(func() bool {
		ok := call < len(rows)
		call++
		return ok
	}).Maybe()
	cursor.On("Close").Return(nil).Maybe()
	cursor.On("Err").Return(nil).Maybe()

	d := new(mockdal.Dal)
	d.On("GetColumns", mock.Anything, mock.Anything).
		Return([]dal.ColumnMeta{idColumn}, nil).Maybe()
	d.On("Cursor", mock.Anything).Return(cursor, nil).Maybe()

	// Fetch copies the row for the current cursor position into the destination map.
	fetched := 0
	d.On("Fetch", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		dst, ok := args.Get(1).(*map[string]interface{})
		if !ok {
			t.Fatalf("Fetch called with %T, want *map[string]interface{}", args.Get(1))
		}
		for k, v := range rows[fetched] {
			(*dst)[k] = v
		}
		fetched++
	}).Return(nil).Maybe()

	executed := make([]string, 0)
	d.On("Exec", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		executed = append(executed, args.Get(0).(string))
	}).Return(nil).Maybe()

	return d, &executed
}

func Test_extractCustomizedFields_skipsOrphanedRowsAndKeepsScanning(t *testing.T) {
	// The cursor LEFT JOINs the raw table, so a domain row whose raw record has been deleted comes
	// back with a NULL data column. Before the fix that row ended the whole scan, leaving every
	// later row unpopulated while the subtask still reported success.
	rows := []fetchedRow{
		{"_raw_data_id": int64(1), "id": "ISSUE-1", "data": nil},
		{"_raw_data_id": int64(2), "id": "ISSUE-2", "data": `{"foo":"bar"}`},
		{"_raw_data_id": int64(3), "id": "ISSUE-3", "data": nil},
		{"_raw_data_id": int64(4), "id": "ISSUE-4", "data": `{"foo":"baz"}`},
	}
	d, executed := newExtractorMocks(t, rows)

	orphaned, err := extractCustomizedFields(
		context.Background(), d, "issues", "_raw_jira_api_issues", "params",
		map[string]string{"x_custom": "foo"})

	assert.NoError(t, err)
	assert.Equal(t, 2, orphaned, "both rows without a raw record should be counted")
	assert.Len(t, *executed, 2, "the two rows that do have raw data should still be updated")
}

func Test_extractCustomizedFields_allRowsOrphanedIsNotAnError(t *testing.T) {
	// A rule whose raw records have all been cleaned up is not a failure; it just has nothing to
	// extract. The caller reports the count so it does not pass unnoticed.
	rows := []fetchedRow{
		{"_raw_data_id": int64(1), "id": "ISSUE-1", "data": nil},
		{"_raw_data_id": int64(2), "id": "ISSUE-2", "data": nil},
	}
	d, executed := newExtractorMocks(t, rows)

	orphaned, err := extractCustomizedFields(
		context.Background(), d, "issues", "_raw_jira_api_issues", "params",
		map[string]string{"x_custom": "foo"})

	assert.NoError(t, err)
	assert.Equal(t, 2, orphaned)
	assert.Empty(t, *executed)
}

func Test_extractCustomizedFields_noOrphansReportsZero(t *testing.T) {
	rows := []fetchedRow{
		{"_raw_data_id": int64(1), "id": "ISSUE-1", "data": `{"foo":"bar"}`},
		{"_raw_data_id": int64(2), "id": "ISSUE-2", "data": `{"foo":"baz"}`},
	}
	d, executed := newExtractorMocks(t, rows)

	orphaned, err := extractCustomizedFields(
		context.Background(), d, "issues", "_raw_jira_api_issues", "params",
		map[string]string{"x_custom": "foo"})

	assert.NoError(t, err)
	assert.Zero(t, orphaned)
	assert.Len(t, *executed, 2)
}

func Test_extractCustomizedFields_honoursCancelledContext(t *testing.T) {
	rows := []fetchedRow{
		{"_raw_data_id": int64(1), "id": "ISSUE-1", "data": `{"foo":"bar"}`},
	}
	d, executed := newExtractorMocks(t, rows)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := extractCustomizedFields(
		ctx, d, "issues", "_raw_jira_api_issues", "params",
		map[string]string{"x_custom": "foo"})

	assert.ErrorIs(t, err, context.Canceled)
	assert.Empty(t, *executed)
}
