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

func TestComputeMetricsDateRangeDefaultLookback(t *testing.T) {
	now := time.Date(2025, 1, 10, 12, 0, 0, 0, time.UTC)
	start, until := computeMetricsDateRange(now, nil)
	require.Equal(t, now, until)
	require.Equal(t, now.AddDate(0, 0, -100), start)
}

func TestComputeMetricsDateRangeUsesSince(t *testing.T) {
	now := time.Date(2025, 1, 10, 12, 0, 0, 0, time.UTC)
	since := now.AddDate(0, 0, -7)
	start, until := computeMetricsDateRange(now, &since)
	require.Equal(t, now, until)
	require.Equal(t, since, start)
}

func TestComputeMetricsDateRangeClampsToLookback(t *testing.T) {
	now := time.Date(2025, 1, 10, 12, 0, 0, 0, time.UTC)
	since := now.AddDate(0, 0, -200)
	start, until := computeMetricsDateRange(now, &since)
	require.Equal(t, now, until)
	require.Equal(t, now.AddDate(0, 0, -100), start)
}

func TestComputeMetricsDateRangeClampsFutureSince(t *testing.T) {
	now := time.Date(2025, 1, 10, 12, 0, 0, 0, time.UTC)
	since := now.Add(24 * time.Hour)
	start, until := computeMetricsDateRange(now, &since)
	require.Equal(t, now, until)
	require.Equal(t, now, start)
}
