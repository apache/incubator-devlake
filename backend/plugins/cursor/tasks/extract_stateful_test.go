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
)

func TestIncrementalRawWhereWithSince(t *testing.T) {
	since := time.Date(2026, 7, 29, 0, 6, 3, 0, time.UTC)
	clause, args := incrementalRawWhere(&since, 18002)
	want := "(created_at >= ? OR id > ?)"
	if clause != want {
		t.Fatalf("clause mismatch:\n got: %s\nwant: %s", clause, want)
	}
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(args))
	}
	gotSince, ok := args[0].(time.Time)
	if !ok || !gotSince.Equal(since) {
		t.Fatalf("since arg mismatch: %#v", args[0])
	}
	gotMax, ok := args[1].(uint64)
	if !ok || gotMax != 18002 {
		t.Fatalf("maxPromoted arg mismatch: %#v", args[1])
	}
}

func TestIncrementalRawWhereWithoutSince(t *testing.T) {
	clause, args := incrementalRawWhere(nil, 100)
	want := "id > ?"
	if clause != want {
		t.Fatalf("clause mismatch:\n got: %s\nwant: %s", clause, want)
	}
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(args))
	}
	gotMax, ok := args[0].(uint64)
	if !ok || gotMax != 100 {
		t.Fatalf("maxPromoted arg mismatch: %#v", args[0])
	}
}

func TestIncrementalRawWhereRepairsUnpromotedGap(t *testing.T) {
	// Simulates production: tool capped at raw id 18002 while raw has 18003+.
	// A since filter alone would miss older created_at rows; the OR id > max
	// clause must include them.
	since := time.Date(2026, 7, 29, 0, 6, 3, 0, time.UTC)
	clause, args := incrementalRawWhere(&since, 18002)
	if clause != "(created_at >= ? OR id > ?)" {
		t.Fatalf("unexpected clause: %s", clause)
	}
	maxPromoted := args[1].(uint64)
	unpromotedID := uint64(18003)
	if !(unpromotedID > maxPromoted) {
		t.Fatalf("unpromoted id %d should be greater than maxPromoted %d", unpromotedID, maxPromoted)
	}
}

func TestExtractRepairGap(t *testing.T) {
	cases := []struct {
		name        string
		maxRaw      uint64
		maxPromoted uint64
		want        uint64
	}{
		{name: "equal", maxRaw: 100, maxPromoted: 100, want: 0},
		{name: "raw behind", maxRaw: 50, maxPromoted: 100, want: 0},
		{name: "small gap", maxRaw: 105, maxPromoted: 100, want: 5},
		{name: "large gap", maxRaw: 18102, maxPromoted: 18002, want: 100},
		{name: "empty tool", maxRaw: 500, maxPromoted: 0, want: 500},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := extractRepairGap(tc.maxRaw, tc.maxPromoted)
			if got != tc.want {
				t.Fatalf("extractRepairGap(%d, %d)=%d want %d", tc.maxRaw, tc.maxPromoted, got, tc.want)
			}
		})
	}
}

func TestExtractRepairGapWarnThreshold(t *testing.T) {
	cases := []struct {
		name string
		gap  uint64
		warn bool
	}{
		{name: "below", gap: 99, warn: false},
		{name: "at", gap: cursorExtractRepairGapWarnThreshold, warn: true},
		{name: "above", gap: 101, warn: true},
		{name: "zero", gap: 0, warn: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			shouldWarn := tc.gap >= cursorExtractRepairGapWarnThreshold
			if shouldWarn != tc.warn {
				t.Fatalf("gap=%d warn=%v want %v", tc.gap, shouldWarn, tc.warn)
			}
		})
	}
}
