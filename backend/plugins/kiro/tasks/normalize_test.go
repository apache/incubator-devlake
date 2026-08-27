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

	"github.com/stretchr/testify/assert"
)

func TestSplitUserId(t *testing.T) {
	tests := []struct {
		name            string
		raw             string
		wantUserId      string
		wantIdentityStr string
	}{
		{
			// Every interaction log record carries this shape.
			name:            "prefixed as seen in interaction logs",
			raw:             "d-1234567890.11111111-1111-4111-8111-111111111111",
			wantUserId:      "11111111-1111-4111-8111-111111111111",
			wantIdentityStr: "d-1234567890",
		},
		{
			// The common report CSV shape.
			name:            "bare uuid passes through",
			raw:             "11111111-1111-4111-8111-111111111111",
			wantUserId:      "11111111-1111-4111-8111-111111111111",
			wantIdentityStr: "",
		},
		{
			name:            "empty",
			raw:             "",
			wantUserId:      "",
			wantIdentityStr: "",
		},
		{
			// A uuid containing dots but no d- prefix must not be split.
			name:            "no d prefix but has dot",
			raw:             "abc.def",
			wantUserId:      "abc.def",
			wantIdentityStr: "",
		},
		{
			name:            "d prefix without separator is not a prefixed id",
			raw:             "d-1234567890",
			wantUserId:      "d-1234567890",
			wantIdentityStr: "",
		},
		{
			name:            "trailing dot carries no user id",
			raw:             "d-1234567890.",
			wantUserId:      "d-1234567890.",
			wantIdentityStr: "",
		},
		{
			// Only the first separator splits; the remainder is the user id
			// verbatim even if it contains further dots.
			name:            "multiple dots split on the first",
			raw:             "d-abc.part1.part2",
			wantUserId:      "part1.part2",
			wantIdentityStr: "d-abc",
		},
		{
			name:            "surrounding whitespace trimmed",
			raw:             "  d-1234567890.11111111  ",
			wantUserId:      "11111111",
			wantIdentityStr: "d-1234567890",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserId, gotIdentityStore := SplitUserId(tt.raw)
			assert.Equal(t, tt.wantUserId, gotUserId)
			assert.Equal(t, tt.wantIdentityStr, gotIdentityStore)
		})
	}
}

func TestParseKiroBool(t *testing.T) {
	trueCases := []string{"true", "TRUE", "True", " true "}
	for _, c := range trueCases {
		t.Run("true/"+c, func(t *testing.T) {
			got, err := ParseKiroBool(c)
			assert.Nil(t, err)
			assert.True(t, got)
		})
	}

	falseCases := []string{"false", "FALSE", "False"}
	for _, c := range falseCases {
		t.Run("false/"+c, func(t *testing.T) {
			got, err := ParseKiroBool(c)
			assert.Nil(t, err)
			assert.False(t, got)
		})
	}

	// An unrecognized value must surface as an error rather than default to
	// false - a silent false would be indistinguishable from real data.
	errCases := []string{"", "yes", "1", "0", "null"}
	for _, c := range errCases {
		t.Run("error/"+c, func(t *testing.T) {
			_, err := ParseKiroBool(c)
			assert.NotNil(t, err)
		})
	}
}

func TestNormalizeTier(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		// Observed values are already canonical.
		{"POWER", "POWER"},
		{"PRO_PLUS", "PRO_PLUS"},
		// The docs spell them in CamelCase.
		{"Power", "POWER"},
		{"ProPlus", "PRO_PLUS"},
		{"ProMax", "PRO_MAX"},
		{"Pro", "PRO"},
		// Unknown tiers must pass through rather than break collection.
		{"SomeFutureTier", "SOME_FUTURE_TIER"},
		{"power", "POWER"},
		{"", ""},
		{"  POWER  ", "POWER"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			assert.Equal(t, tt.want, NormalizeTier(tt.in))
		})
	}
}

func TestParseModelColumn(t *testing.T) {
	// All eight model/mode columns observed across the full report history.
	modelColumns := map[string]string{
		"claude_opus_4.6_messages":   "claude_opus_4.6",
		"claude_opus_4.7_messages":   "claude_opus_4.7",
		"claude_opus_4.8_messages":   "claude_opus_4.8",
		"claude_opus_5_messages":     "claude_opus_5",
		"claude_sonnet_4.6_messages": "claude_sonnet_4.6",
		"gpt_5.6_sol_messages":       "gpt_5.6_sol",
		"auto_messages":              "auto",
		"simple_task_messages":       "simple_task",
	}
	for header, want := range modelColumns {
		t.Run("model/"+header, func(t *testing.T) {
			got, ok := ParseModelColumn(header)
			assert.True(t, ok)
			assert.Equal(t, want, got)
		})
	}

	// Total_Messages matches the pattern but is the aggregate. If it were
	// treated as a model named "Total", its value would equal the sum of all
	// real models - a wrong result that looks entirely plausible downstream.
	// The exclusion is case-insensitive because the pattern is.
	notModels := []string{
		"Total_Messages",
		"total_messages",
		"TOTAL_MESSAGES",
		" Total_Messages ",
		"Chat_Conversations",
		"Credits_Used",
		"User_Email",
		"_messages",
		"messages",
		"",
	}
	for _, header := range notModels {
		t.Run("not_model/"+header, func(t *testing.T) {
			got, ok := ParseModelColumn(header)
			assert.False(t, ok, "header %q must not be treated as a model", header)
			assert.Equal(t, "", got)
		})
	}
}

func TestCanonicalModelName(t *testing.T) {
	// The two sources spell the same model differently; canonical form lets
	// them be compared without rewriting either source's stored value.
	assert.Equal(t,
		CanonicalModelName("claude-opus-4.6"),
		CanonicalModelName("claude_opus_4.6"),
	)
	assert.Equal(t,
		CanonicalModelName("simple-task"),
		CanonicalModelName("simple_task"),
	)
	// Distinct models must stay distinct.
	assert.NotEqual(t,
		CanonicalModelName("claude-opus-4.6"),
		CanonicalModelName("claude-opus-4.8"),
	)
	assert.Equal(t, "claude_opus_4.6", CanonicalModelName("Claude-Opus-4.6"))
}

func TestParseKiroTime(t *testing.T) {
	// The real shape emitted by interaction logs.
	t.Run("nanosecond precision truncates to microsecond", func(t *testing.T) {
		got, err := ParseKiroTime("2026-07-27T23:03:29.027400929Z")
		assert.Nil(t, err)
		want := time.Date(2026, 7, 27, 23, 3, 29, 27400000, time.UTC)
		assert.True(t, got.Equal(want), "got %v want %v", got, want)
		// Nanosecond remainder beyond microsecond precision must be gone, so
		// the value round-trips through DATETIME(6) unchanged.
		assert.Equal(t, 0, got.Nanosecond()%1000)
	})

	t.Run("microsecond precision preserved", func(t *testing.T) {
		got, err := ParseKiroTime("2026-07-27T23:03:29.027400Z")
		assert.Nil(t, err)
		assert.Equal(t, 27400000, got.Nanosecond())
	})

	t.Run("no fractional seconds", func(t *testing.T) {
		got, err := ParseKiroTime("2026-07-27T23:03:29Z")
		assert.Nil(t, err)
		assert.Equal(t, 0, got.Nanosecond())
		assert.Equal(t, 29, got.Second())
	})

	for _, bad := range []string{"", "not-a-time", "2026-07-27", "07-27-2026"} {
		t.Run("error/"+bad, func(t *testing.T) {
			_, err := ParseKiroTime(bad)
			assert.NotNil(t, err)
		})
	}
}

func TestParseReportDate(t *testing.T) {
	t.Run("iso date", func(t *testing.T) {
		got, err := ParseReportDate("2026-07-27")
		assert.Nil(t, err)
		assert.Equal(t, 2026, got.Year())
		assert.Equal(t, time.July, got.Month())
		assert.Equal(t, 27, got.Day())
	})

	// The legacy by_user_analytic report used MM-DD-YYYY. That report is out of
	// scope, but a lenient parser would accept its format here and produce a
	// wrong date without any error - so rejection is asserted explicitly.
	t.Run("legacy MM-DD-YYYY rejected", func(t *testing.T) {
		_, err := ParseReportDate("07-27-2026")
		assert.NotNil(t, err)
	})

	for _, bad := range []string{"", "2026/07/27", "2026-13-45", "not-a-date"} {
		t.Run("error/"+bad, func(t *testing.T) {
			_, err := ParseReportDate(bad)
			assert.NotNil(t, err)
		})
	}
}
