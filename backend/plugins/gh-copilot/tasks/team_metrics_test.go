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
	"time"

	"github.com/stretchr/testify/require"
)

func TestTeamMetricsWindow(t *testing.T) {
	now := time.Date(2026, 3, 5, 15, 30, 0, 0, time.UTC)
	since, until := teamMetricsWindow(now)

	require.Equal(t, "2026-02-05", since)
	require.Equal(t, "2026-03-04", until)
}

func TestParseTeamMetricsDate(t *testing.T) {
	parsed, err := parseTeamMetricsDate("2024-06-24")
	require.NoError(t, err)
	require.Equal(t, time.Date(2024, 6, 24, 0, 0, 0, 0, time.UTC), parsed)
}

func TestParseTeamMetricsDateInvalid(t *testing.T) {
	_, err := parseTeamMetricsDate("2024/06/24")
	require.Error(t, err)
}

func TestParseTeamCustomModelTrainingDateNilOrBlank(t *testing.T) {
	parsed, err := parseTeamCustomModelTrainingDate(nil)
	require.NoError(t, err)
	require.Nil(t, parsed)

	blank := "   "
	parsed, err = parseTeamCustomModelTrainingDate(&blank)
	require.NoError(t, err)
	require.Nil(t, parsed)
}

func TestParseTeamCustomModelTrainingDateDateOnly(t *testing.T) {
	dateOnly := "2024-02-01"
	parsed, err := parseTeamCustomModelTrainingDate(&dateOnly)
	require.NoError(t, err)
	require.NotNil(t, parsed)
	require.Equal(t, time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), *parsed)
}

func TestParseTeamCustomModelTrainingDateRFC3339(t *testing.T) {
	rfc3339 := "2024-02-01T10:11:12Z"
	parsed, err := parseTeamCustomModelTrainingDate(&rfc3339)
	require.NoError(t, err)
	require.NotNil(t, parsed)
	require.Equal(t, "2024-02-01", parsed.UTC().Format("2006-01-02"))
}

func TestParseTeamCustomModelTrainingDateInvalid(t *testing.T) {
	invalid := "not-a-date"
	_, err := parseTeamCustomModelTrainingDate(&invalid)
	require.Error(t, err)
}
