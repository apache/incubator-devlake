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
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseRetryAfterSeconds(t *testing.T) {
	wait := parseRetryAfter("10", time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
	require.Equal(t, 10*time.Second, wait)
}

func TestParseRetryAfterHttpDate(t *testing.T) {
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	value := now.Add(5 * time.Second).Format(http.TimeFormat)
	wait := parseRetryAfter(value, now)
	require.True(t, wait >= 4*time.Second && wait <= 6*time.Second)
}

func TestComputeReportDateRangeDefaultLookback(t *testing.T) {
	now := time.Date(2025, 1, 10, 12, 0, 0, 0, time.UTC)
	start, until := computeReportDateRange(now, nil)
	require.Equal(t, time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC), until)
	require.Equal(t, time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC), start)
}

func TestComputeReportDateRangeUsesSince(t *testing.T) {
	now := time.Date(2025, 1, 10, 12, 0, 0, 0, time.UTC)
	since := time.Date(2025, 1, 3, 12, 0, 0, 0, time.UTC)
	start, until := computeReportDateRange(now, &since)
	require.Equal(t, time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC), until)
	require.Equal(t, time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC), start)
}

func TestComputeReportDateRangeClampsToLookback(t *testing.T) {
	now := time.Date(2025, 1, 10, 12, 0, 0, 0, time.UTC)
	since := time.Date(2024, 6, 24, 12, 0, 0, 0, time.UTC)
	start, until := computeReportDateRange(now, &since)
	require.Equal(t, time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC), until)
	require.Equal(t, time.Date(2024, 6, 24, 0, 0, 0, 0, time.UTC), start)
}

func TestComputeReportDateRangeClampsFutureSince(t *testing.T) {
	now := time.Date(2025, 1, 10, 12, 0, 0, 0, time.UTC)
	since := now.Add(24 * time.Hour)
	start, until := computeReportDateRange(now, &since)
	require.Equal(t, time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC), until)
	require.Equal(t, time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC), start)
}
