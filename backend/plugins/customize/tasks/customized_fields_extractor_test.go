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

	"github.com/apache/incubator-devlake/core/dal"
	mockdal "github.com/apache/incubator-devlake/mocks/core/dal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestExtractCustomizedFieldsSkipsOrphanedRows is a regression test for the case where a
// domain-layer row's `_raw_data_id` has no matching raw record (LEFT JOIN yields NULL data).
// Such a row must be skipped and extraction must continue with the remaining rows, instead of
// aborting the whole scan on the first orphaned row.
func TestExtractCustomizedFieldsSkipsOrphanedRows(t *testing.T) {
	mockRows := new(mockdal.Rows)
	mockRows.On("Next").Return(true).Times(2)
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Close").Return(nil).Once()

	mockDal := new(mockdal.Dal)
	mockDal.On("GetColumns", mock.Anything, mock.Anything).Return([]dal.ColumnMeta{}, nil)
	mockDal.On("Cursor", mock.Anything).Return(mockRows, nil).Once()

	// Row 1 is orphaned: no matching raw record, so `data` is NULL (nil). It must be skipped.
	mockDal.On("Fetch", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		dst := args.Get(1).(*map[string]interface{})
		*dst = map[string]interface{}{
			"id":           "row-orphaned",
			"_raw_data_id": int64(1),
			"data":         nil,
		}
	}).Return(nil).Once()
	// Row 2 has valid raw data and matches the filter; it must still be processed.
	mockDal.On("Fetch", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		dst := args.Get(1).(*map[string]interface{})
		*dst = map[string]interface{}{
			"id":           "row-valid",
			"_raw_data_id": int64(2),
			"data":         `{"name":"hello"}`,
		}
	}).Return(nil).Once()

	execCalls := 0
	mockDal.On("Exec", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		execCalls++
	}).Return(nil)

	err := extractCustomizedFields(
		context.Background(),
		mockDal,
		"boards", // non-issues table -> simple column update path
		"_raw_jira_api_issues",
		`{"ConnectionId":1`,
		map[string]string{"x_test": "name"},
	)

	assert.Nil(t, err)
	// The valid row after the orphaned one must have been updated exactly once.
	// Before the fix, the orphaned row hit `default: return nil` and aborted the scan,
	// so Exec was never called.
	assert.Equal(t, 1, execCalls)
	mockDal.AssertExpectations(t)
}
