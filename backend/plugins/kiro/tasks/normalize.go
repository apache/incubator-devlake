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
	"regexp"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
)

// Kiro exports the same logical value in several shapes across its two data
// sources. Every normalization helper here exists because a real sample proved
// the variation exists; see the design spec section 3 for the evidence.

// modelColumnRe matches the per-model message count columns in the user report
// CSV, e.g. "claude_opus_4.6_messages". The column set is dynamic - it grows
// and shrinks with the models a team actually uses.
var modelColumnRe = regexp.MustCompile(`(?i)^(.+)_messages$`)

// totalMessagesColumn is the aggregate column that must never be mistaken for a
// model. It matches modelColumnRe, and its value happens to equal the sum of
// all model columns, so treating it as a model yields a plausible-looking but
// wrong distribution.
const totalMessagesColumn = "total_messages"

// SplitUserId separates an identity store prefixed user id into its parts.
//
// Interaction logs always carry the prefix ("d-1234567890.11111111-..."), while
// report CSVs usually do not - but three sampled report files did. Both paths
// must normalize, otherwise the same person lands in the table twice and every
// per-person aggregate is silently wrong.
//
// Returns the bare user id, plus the identity store id when a prefix was
// present (empty string otherwise).
func SplitUserId(raw string) (userId string, identityStoreId string) {
	trimmed := strings.TrimSpace(raw)
	if !strings.HasPrefix(trimmed, "d-") {
		return trimmed, ""
	}
	// Only the first separator matters: the identity store id itself never
	// contains a dot, while the trailing UUID may not be validated here.
	idx := strings.Index(trimmed, ".")
	if idx < 0 {
		// "d-something" with no separator is not a prefixed id.
		return trimmed, ""
	}
	prefix := trimmed[:idx]
	rest := trimmed[idx+1:]
	if rest == "" {
		// A trailing dot with nothing after it carries no user id; keep the
		// original value rather than inventing an empty one.
		return trimmed, ""
	}
	return rest, prefix
}

// ParseKiroBool parses a boolean whose casing drifts between report versions:
// March samples emit "TRUE", July samples emit "true".
func ParseKiroBool(s string) (bool, errors.Error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, errors.Default.New("unrecognized boolean value: " + s)
	}
}

// NormalizeTier canonicalizes a subscription tier to UPPER_SNAKE.
//
// Observed values are already UPPER_SNAKE ("POWER", "PRO_PLUS") while the
// official docs write them in CamelCase ("Power", "ProPlus"). Unknown tiers are
// upper-cased and returned as-is rather than rejected - a new tier must not
// break collection.
func NormalizeTier(s string) string {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return ""
	}
	// Insert a separator at lower-to-upper transitions so "ProPlus" becomes
	// "PRO_PLUS", while an already-snake "PRO_PLUS" is left alone.
	var b strings.Builder
	runes := []rune(trimmed)
	for i, r := range runes {
		if i > 0 && r >= 'A' && r <= 'Z' {
			prev := runes[i-1]
			if prev >= 'a' && prev <= 'z' {
				b.WriteRune('_')
			}
		}
		b.WriteRune(r)
	}
	return strings.ToUpper(b.String())
}

// ParseModelColumn reports whether a CSV header is a per-model message count
// column, returning the model name when it is.
//
// "Total_Messages" is excluded explicitly: it matches the same pattern but is
// the aggregate, not a model.
func ParseModelColumn(header string) (string, bool) {
	trimmed := strings.TrimSpace(header)
	if strings.EqualFold(trimmed, totalMessagesColumn) {
		return "", false
	}
	m := modelColumnRe.FindStringSubmatch(trimmed)
	if m == nil {
		return "", false
	}
	name := m[1]
	if name == "" {
		return "", false
	}
	return name, true
}

// CanonicalModelName folds the two spellings of a model identifier onto one
// form so they can be compared.
//
// The report CSV uses underscores ("claude_opus_4.6") while interaction logs
// use hyphens ("claude-opus-4.6"). Values are stored as the source emits them;
// this helper is for comparison only.
func CanonicalModelName(s string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(s), "-", "_"))
}

// ParseKiroTime parses an interaction log timestamp.
//
// Log timestamps carry nanosecond precision ("2026-07-27T23:03:29.027400929Z")
// but MySQL DATETIME(6) only stores microseconds, so the value is truncated
// here rather than at the storage layer - that keeps the truncation explicit
// and identical on every database backend.
func ParseKiroTime(s string) (time.Time, errors.Error) {
	trimmed := strings.TrimSpace(s)
	t, err := time.Parse(time.RFC3339Nano, trimmed)
	if err != nil {
		return time.Time{}, errors.Default.Wrap(err, "failed to parse kiro timestamp: "+trimmed)
	}
	return t.Truncate(time.Microsecond), nil
}

// ParseReportDate parses a user report date, which is strictly ISO 8601.
//
// Deliberately strict: a lenient parser would accept the legacy
// by_user_analytic format (MM-DD-YYYY) and silently produce a wrong date.
func ParseReportDate(s string) (time.Time, errors.Error) {
	trimmed := strings.TrimSpace(s)
	t, err := time.Parse(time.DateOnly, trimmed)
	if err != nil {
		return time.Time{}, errors.Default.Wrap(err, "failed to parse kiro report date: "+trimmed)
	}
	return t, nil
}
