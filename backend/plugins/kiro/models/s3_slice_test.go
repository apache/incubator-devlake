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

package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func intPtr(i int) *int { return &i }

func TestS3SliceNormalize(t *testing.T) {
	t.Run("month scope derives id and name", func(t *testing.T) {
		s := &KiroS3Slice{AccountId: "123456789012", Year: 2026, Month: intPtr(7)}
		assert.NoError(t, s.normalize(true))
		assert.Equal(t, "123456789012_2026_07", s.Id)
		assert.Equal(t, "123456789012 2026-07", s.Name)
		assert.Equal(t, "2026/07", s.TimePath())
	})

	t.Run("year scope widens to the whole year", func(t *testing.T) {
		s := &KiroS3Slice{AccountId: "123456789012", Year: 2026}
		assert.NoError(t, s.normalize(true))
		assert.Equal(t, "123456789012_2026", s.Id)
		assert.Equal(t, "2026", s.TimePath())
	})

	t.Run("single digit month is zero padded", func(t *testing.T) {
		s := &KiroS3Slice{AccountId: "acct", Year: 2026, Month: intPtr(3)}
		assert.NoError(t, s.normalize(true))
		// Padding matters: Kiro's own S3 partitions are two-digit, so an
		// unpadded value would build a prefix that matches nothing.
		assert.Equal(t, "acct_2026_03", s.Id)
		assert.Equal(t, "2026/03", s.TimePath())
	})

	t.Run("explicit id is preserved", func(t *testing.T) {
		s := &KiroS3Slice{Id: "custom-id", AccountId: "acct", Year: 2026}
		assert.NoError(t, s.normalize(true))
		assert.Equal(t, "custom-id", s.Id)
	})

	t.Run("base path is trimmed", func(t *testing.T) {
		s := &KiroS3Slice{BasePath: " /nested/path/ ", AccountId: "acct", Year: 2026}
		assert.NoError(t, s.normalize(true))
		assert.Equal(t, "nested/path", s.BasePath)
	})
}

func TestS3SliceValidation(t *testing.T) {
	t.Run("account id required in strict mode", func(t *testing.T) {
		s := &KiroS3Slice{Year: 2026}
		assert.Error(t, s.normalize(true))
	})

	t.Run("year required in strict mode", func(t *testing.T) {
		s := &KiroS3Slice{AccountId: "acct"}
		assert.Error(t, s.normalize(true))
	})

	// Non-strict mode runs on read, where rejecting an existing row would make
	// it unreadable rather than merely incomplete.
	t.Run("missing fields tolerated in non-strict mode", func(t *testing.T) {
		s := &KiroS3Slice{}
		assert.NoError(t, s.normalize(false))
	})

	for _, month := range []int{0, 13, -1} {
		t.Run("invalid month rejected in both modes", func(t *testing.T) {
			s := &KiroS3Slice{AccountId: "acct", Year: 2026, Month: intPtr(month)}
			assert.Error(t, s.normalize(true))
			// An out-of-range month cannot produce a valid prefix, so it is
			// rejected on read too rather than silently building a bad path.
			s2 := &KiroS3Slice{AccountId: "acct", Year: 2026, Month: intPtr(month)}
			assert.Error(t, s2.normalize(false))
		})
	}
}

func TestS3SliceScopeInterface(t *testing.T) {
	s := KiroS3Slice{AccountId: "123456789012", Year: 2026, Month: intPtr(7)}
	assert.NoError(t, s.normalize(true))

	assert.Equal(t, "123456789012_2026_07", s.ScopeId())
	assert.Equal(t, "123456789012 2026-07", s.ScopeName())
	assert.Equal(t, "123456789012 2026-07", s.ScopeFullName())
	assert.Equal(t, "_tool_kiro_s3_slices", s.TableName())

	params, ok := s.ScopeParams().(*KiroS3SliceParams)
	assert.True(t, ok)
	assert.Equal(t, "123456789012", params.AccountId)
	assert.Equal(t, 2026, params.Year)
	assert.Equal(t, 7, *params.Month)
}
