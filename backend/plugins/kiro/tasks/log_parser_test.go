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
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Fixtures are real log files whose text content was replaced but whose exact
// string lengths and line structure were preserved, because those are the
// values the parser derives and this test asserts.
func loadLogFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "logs", name))
	require.NoError(t, err)
	return data
}

func gzipBytes(t *testing.T, s string) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	_, err := w.Write([]byte(s))
	require.NoError(t, err)
	require.NoError(t, w.Close())
	return buf.Bytes()
}

// An empty prompt means the agent continued on its own after a tool call rather
// than the user speaking. This is the majority case - roughly 71% of sampled
// records - and it must not be confused with a user turn.
func TestParseChatLog_EmptyPrompt(t *testing.T) {
	logs, err := ParseChatLog(loadLogFixture(t, "chat_01_empty_prompt.json.gz"), testConnectionId, testScopeId)
	require.Nil(t, err)
	require.Len(t, logs, 1)

	l := logs[0]
	assert.False(t, l.HasPrompt)
	assert.Equal(t, 0, l.PromptLength)
	// Hashing the empty string would give ~71% of rows one shared hash and
	// destroy the rework signal, so it stays NULL.
	assert.Nil(t, l.PromptSha256)
	// Both heuristics read prompt text, so without a prompt they are unknown
	// rather than false.
	assert.Nil(t, l.HasSteering)
	assert.Nil(t, l.IsSpecMode)

	assert.Equal(t, "64d13ea7-dff5-4563-9285-6a9e351e87a0", l.RequestId)
	assert.Equal(t, "MANUAL", l.ChatTriggerType)
	require.NotNil(t, l.ModelId)
	assert.Equal(t, "claude-opus-5", *l.ModelId)
	assert.Equal(t, 259, l.ResponseLength)

	// Log records always carry the identity-store prefix; stripping it is what
	// lets this join against the report table.
	assert.Equal(t, "11111111-1111-4111-8111-111111111111", l.UserId)
	assert.Equal(t, "d-1234567890", l.IdentityStoreId)

	// Nanosecond input truncated for DATETIME(6).
	assert.Equal(t, 2026, l.Timestamp.Year())
	assert.Equal(t, 23, l.Timestamp.Hour())
	assert.Equal(t, 0, l.Timestamp.Nanosecond()%1000)

	// Documented but never observed, so they stay NULL - their absence is why
	// the S3 logs cannot group interactions into sessions.
	assert.Nil(t, l.ConversationId)
	assert.Nil(t, l.UtteranceId)
}

func TestParseChatLog_WithPrompt(t *testing.T) {
	logs, err := ParseChatLog(loadLogFixture(t, "chat_02_with_prompt.json.gz"), testConnectionId, testScopeId)
	require.Nil(t, err)
	require.Len(t, logs, 1)

	l := logs[0]
	assert.True(t, l.HasPrompt)
	assert.Equal(t, 1363, l.PromptLength)
	require.NotNil(t, l.PromptSha256)
	assert.Len(t, *l.PromptSha256, 64)
	// With a prompt present the heuristics carry a real verdict.
	require.NotNil(t, l.HasSteering)
	require.NotNil(t, l.IsSpecMode)
	assert.Equal(t, 64, l.ResponseLength)
}

// Kiro packs one or two records per file, so the array is always walked.
func TestParseChatLog_TwoRecords(t *testing.T) {
	logs, err := ParseChatLog(loadLogFixture(t, "chat_03_two_records.json.gz"), testConnectionId, testScopeId)
	require.Nil(t, err)
	require.Len(t, logs, 2)

	assert.NotEqual(t, logs[0].RequestId, logs[1].RequestId)
	// A zero-length assistant response is real data, not a parse failure.
	assert.Equal(t, 0, logs[0].ResponseLength)
	assert.Equal(t, 117, logs[1].ResponseLength)
}

// modelId is absent on more than half of all records. Storing an empty string
// instead of NULL would invent a phantom model accounting for most of the
// traffic.
func TestParseChatLog_MissingModelId(t *testing.T) {
	logs, err := ParseChatLog(loadLogFixture(t, "chat_04_no_model_id.json.gz"), testConnectionId, testScopeId)
	require.Nil(t, err)
	require.Len(t, logs, 1)

	l := logs[0]
	assert.Nil(t, l.ModelId, "an absent modelId must be NULL, never an empty string")
	// This record does carry followupPrompts, which a single-day sample had
	// suggested never appears.
	assert.True(t, l.HasFollowupPrompts)
	assert.True(t, l.HasPrompt)
	assert.Equal(t, 175, l.PromptLength)
}

func TestParseChatLog_PromptHashing(t *testing.T) {
	// The same prompt text must hash identically, which is what makes repeated
	// submission detectable as a rework signal.
	body := `{"records":[{"generateAssistantResponseEventRequest":` +
		`{"prompt":"fix the retry logic","chatTriggerType":"MANUAL",` +
		`"userId":"d-abc.user-1","timeStamp":"2026-07-27T23:03:29.027400929Z"},` +
		`"generateAssistantResponseEventResponse":{"assistantResponse":"ok","requestId":"%s"}}]}`

	first, err := ParseChatLog(gzipBytes(t, fmt.Sprintf(body, "req-1")), testConnectionId, testScopeId)
	require.Nil(t, err)
	second, err := ParseChatLog(gzipBytes(t, fmt.Sprintf(body, "req-2")), testConnectionId, testScopeId)
	require.Nil(t, err)

	require.Len(t, first, 1)
	require.Len(t, second, 1)
	require.NotNil(t, first[0].PromptSha256)
	require.NotNil(t, second[0].PromptSha256)
	assert.Equal(t, *first[0].PromptSha256, *second[0].PromptSha256)
	assert.NotEqual(t, first[0].RequestId, second[0].RequestId)
}

func TestParseChatLog_Heuristics(t *testing.T) {
	tests := []struct {
		name            string
		prompt          string
		wantHasSteering bool
		wantIsSpecMode  bool
	}{
		{"steering reference", "please follow .kiro/steering/style.md", true, false},
		{"spec reference", "implement .kiro/specs/auth/tasks.md", false, true},
		{"both", "read .kiro/steering and .kiro/specs", true, true},
		{"neither", "just fix this bug", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := `{"records":[{"generateAssistantResponseEventRequest":` +
				`{"prompt":"` + tt.prompt + `","chatTriggerType":"MANUAL",` +
				`"userId":"d-abc.user-1","timeStamp":"2026-07-27T23:03:29Z"},` +
				`"generateAssistantResponseEventResponse":{"assistantResponse":"ok","requestId":"r1"}}]}`
			logs, err := ParseChatLog(gzipBytes(t, body), testConnectionId, testScopeId)
			require.Nil(t, err)
			require.Len(t, logs, 1)
			require.NotNil(t, logs[0].HasSteering)
			require.NotNil(t, logs[0].IsSpecMode)
			assert.Equal(t, tt.wantHasSteering, *logs[0].HasSteering)
			assert.Equal(t, tt.wantIsSpecMode, *logs[0].IsSpecMode)
		})
	}
}

func TestParseCompletionLog_NonEmpty(t *testing.T) {
	logs, err := ParseCompletionLog(loadLogFixture(t, "completion_01_non_empty.json.gz"), testConnectionId, testScopeId)
	require.Nil(t, err)
	require.Len(t, logs, 1)

	l := logs[0]
	assert.Equal(t, "e0d1e760-a907-47d4-baa5-93df0f38b274", l.RequestId)
	// A bare file name with no directory path, so it cannot be resolved to a
	// unique repository file.
	assert.Equal(t, "mcp.json", l.FileName)
	assert.Equal(t, "json", l.FileExtension)
	assert.Equal(t, 1, l.CompletionsCount)
	assert.Equal(t, 347, l.ReturnedCharCount)
	assert.Equal(t, 15, l.ReturnedLineCount)
	assert.Equal(t, 5324, l.LeftContextLength)
	assert.Equal(t, 12, l.RightContextLength)
	assert.False(t, l.HasCustomization)
	assert.Equal(t, "11111111-1111-4111-8111-111111111111", l.UserId)
}

// A completion record is written when the suggestion is requested, not when it
// is accepted, so an empty array is normal. The record is the denominator and
// must still be stored.
func TestParseCompletionLog_EmptyCompletions(t *testing.T) {
	logs, err := ParseCompletionLog(loadLogFixture(t, "completion_02_empty.json.gz"), testConnectionId, testScopeId)
	require.Nil(t, err)
	require.Len(t, logs, 1, "a record with no completions is still stored")

	l := logs[0]
	assert.Equal(t, 0, l.CompletionsCount)
	assert.Equal(t, 0, l.ReturnedCharCount)
	assert.Equal(t, 0, l.ReturnedLineCount)
	// Context was still sent, which is what distinguishes this from a truncated
	// record.
	assert.Equal(t, 5323, l.LeftContextLength)
}

func TestParseCompletionLog_LineCounting(t *testing.T) {
	tests := []struct {
		name        string
		completions string
		wantLines   int
		wantChars   int
	}{
		{"single line", `["abc"]`, 1, 3},
		{"two lines", `["a\nb"]`, 2, 3},
		{"trailing newline counts the empty final line", `["a\n"]`, 2, 2},
		{"two completions summed", `["a\nb","c"]`, 3, 4},
		{"empty array", `[]`, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := `{"records":[{"generateCompletionsEventRequest":` +
				`{"fileName":"f.ts","leftContext":"","rightContext":"",` +
				`"userId":"d-abc.user-1","timeStamp":"2026-03-19T13:49:58Z"},` +
				`"generateCompletionsEventResponse":{"completions":` + tt.completions + `,"requestId":"r1"}}]}`
			logs, err := ParseCompletionLog(gzipBytes(t, body), testConnectionId, testScopeId)
			require.Nil(t, err)
			require.Len(t, logs, 1)
			assert.Equal(t, tt.wantLines, logs[0].ReturnedLineCount)
			assert.Equal(t, tt.wantChars, logs[0].ReturnedCharCount)
		})
	}
}

func TestParseCompletionLog_Customization(t *testing.T) {
	// customizationArn appears on completion records but never on chat records.
	body := `{"records":[{"generateCompletionsEventRequest":` +
		`{"fileName":"f.ts","userId":"d-abc.u1","timeStamp":"2026-03-19T13:49:58Z",` +
		`"customizationArn":"arn:aws:codewhisperer:us-east-1:1:customization/abc"},` +
		`"generateCompletionsEventResponse":{"completions":[],"requestId":"r1"}}]}`
	logs, err := ParseCompletionLog(gzipBytes(t, body), testConnectionId, testScopeId)
	require.Nil(t, err)
	require.Len(t, logs, 1)
	assert.True(t, logs[0].HasCustomization)
}

func TestParseLog_EdgeCases(t *testing.T) {
	t.Run("empty records array", func(t *testing.T) {
		logs, err := ParseChatLog(gzipBytes(t, `{"records":[]}`), testConnectionId, testScopeId)
		assert.Nil(t, err)
		assert.Empty(t, logs)
	})

	t.Run("record without requestId is skipped", func(t *testing.T) {
		// requestId is the primary key, so such a record cannot be stored or
		// deduplicated.
		body := `{"records":[{"generateAssistantResponseEventRequest":` +
			`{"prompt":"","userId":"d-abc.u1","timeStamp":"2026-07-27T23:03:29Z"},` +
			`"generateAssistantResponseEventResponse":{"assistantResponse":"ok"}}]}`
		logs, err := ParseChatLog(gzipBytes(t, body), testConnectionId, testScopeId)
		assert.Nil(t, err)
		assert.Empty(t, logs)
	})

	t.Run("record missing the response half is skipped", func(t *testing.T) {
		body := `{"records":[{"generateAssistantResponseEventRequest":` +
			`{"prompt":"x","userId":"d-abc.u1","timeStamp":"2026-07-27T23:03:29Z"}}]}`
		logs, err := ParseChatLog(gzipBytes(t, body), testConnectionId, testScopeId)
		assert.Nil(t, err)
		assert.Empty(t, logs)
	})

	t.Run("bad timestamp is an error", func(t *testing.T) {
		body := `{"records":[{"generateAssistantResponseEventRequest":` +
			`{"prompt":"","userId":"d-abc.u1","timeStamp":"not-a-time"},` +
			`"generateAssistantResponseEventResponse":{"assistantResponse":"ok","requestId":"r1"}}]}`
		_, err := ParseChatLog(gzipBytes(t, body), testConnectionId, testScopeId)
		assert.NotNil(t, err)
	})

	t.Run("non-gzip input is an error", func(t *testing.T) {
		_, err := ParseChatLog([]byte("not gzipped"), testConnectionId, testScopeId)
		assert.NotNil(t, err)
	})

	t.Run("malformed json is an error", func(t *testing.T) {
		_, err := ParseChatLog(gzipBytes(t, `{"records":`), testConnectionId, testScopeId)
		assert.NotNil(t, err)
	})

	t.Run("connection and scope are propagated", func(t *testing.T) {
		body := `{"records":[{"generateAssistantResponseEventRequest":` +
			`{"prompt":"","userId":"d-abc.u1","timeStamp":"2026-07-27T23:03:29Z"},` +
			`"generateAssistantResponseEventResponse":{"assistantResponse":"ok","requestId":"r1"}}]}`
		logs, err := ParseChatLog(gzipBytes(t, body), 42, "scope-x")
		require.Nil(t, err)
		require.Len(t, logs, 1)
		assert.Equal(t, uint64(42), logs[0].ConnectionId)
		assert.Equal(t, "scope-x", logs[0].ScopeId)
	})
}
