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
	"encoding/json"
	"testing"
	"time"

	helper "github.com/apache/devlake/helpers/pluginhelper/api"
	"github.com/apache/devlake/plugins/cursor/models"
)

func TestRawParamsFromTaskDataIncludesEndpoint(t *testing.T) {
	taskData := &CursorTaskData{
		Options: &CursorOptions{
			ConnectionId: 1,
			ScopeId:      "team",
		},
		Connection: &models.CursorConnection{
			CursorConn: models.CursorConn{
				RestConnection: helper.RestConnection{Endpoint: models.DefaultEndpoint},
			},
		},
	}

	params := rawParamsFromTaskData(taskData)
	encoded, err := json.Marshal(params.GetParams())
	if err != nil {
		t.Fatalf("marshal params: %v", err)
	}

	expected := `{"ConnectionId":1,"ScopeId":"team","Endpoint":"https://api.cursor.com"}`
	if string(encoded) != expected {
		t.Fatalf("raw params mismatch:\n got: %s\nwant: %s", encoded, expected)
	}
}

func TestSplitDailyUsageTimeRangeMsChunksLongRanges(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
	chunks := splitDailyUsageTimeRangeMs(start.UnixMilli(), end.UnixMilli(), cursorDailyUsageMaxDays)

	if len(chunks) != 4 {
		t.Fatalf("expected 4 chunks for ~90-day span, got %d", len(chunks))
	}
	if chunks[0].StartDateMs != start.UnixMilli() {
		t.Fatalf("first chunk start mismatch: got %d want %d", chunks[0].StartDateMs, start.UnixMilli())
	}
	if chunks[len(chunks)-1].EndDateMs != end.UnixMilli() {
		t.Fatalf("last chunk end mismatch: got %d want %d", chunks[len(chunks)-1].EndDateMs, end.UnixMilli())
	}
	for i := 1; i < len(chunks); i++ {
		if chunks[i].StartDateMs <= chunks[i-1].StartDateMs {
			t.Fatalf("chunk %d does not advance start time", i)
		}
		maxSpan := time.Duration(cursorDailyUsageMaxDays) * 24 * time.Hour
		span := time.UnixMilli(chunks[i].StartDateMs).Sub(time.UnixMilli(chunks[i-1].StartDateMs))
		if span > maxSpan+time.Millisecond {
			t.Fatalf("gap between chunk %d and %d exceeds %d days", i-1, i, cursorDailyUsageMaxDays)
		}
	}
}

func TestSplitDailyUsageTimeRangeMsEmptyRange(t *testing.T) {
	if chunks := splitDailyUsageTimeRangeMs(100, 100, cursorDailyUsageMaxDays); len(chunks) != 0 {
		t.Fatalf("expected no chunks for empty range, got %d", len(chunks))
	}
}

func TestNewUsageEventsDateRangeIteratorChunksLongRanges(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
	chunks := splitDailyUsageTimeRangeMs(start.UnixMilli(), end.UnixMilli(), cursorDailyUsageMaxDays)

	if len(chunks) != 4 {
		t.Fatalf("expected 4 chunks for ~90-day usage-events span, got %d", len(chunks))
	}
}

func TestDailyUsagePostBodyUsesEpochMilliseconds(t *testing.T) {
	body := dailyUsagePostBody(&helper.RequestData{
		Input: cursorTimeRangeInput{StartDateMs: 1710720000000, EndDateMs: 1710892800000},
		Pager: &helper.Pager{Page: 1, Size: 100},
	})
	if body["startDate"] != int64(1710720000000) {
		t.Fatalf("startDate should be epoch ms int64, got %#v", body["startDate"])
	}
	if body["endDate"] != int64(1710892800000) {
		t.Fatalf("endDate should be epoch ms int64, got %#v", body["endDate"])
	}
}

func TestComputeUsageTimeRangeMsFullSyncUsesBackfill(t *testing.T) {
	now := time.Date(2026, 7, 30, 12, 0, 0, 0, time.UTC)
	startMs, endMs := computeUsageTimeRangeMs(nil, now, false)
	start := time.UnixMilli(startMs).UTC()
	end := time.UnixMilli(endMs).UTC()
	if !end.Equal(now) {
		t.Fatalf("end mismatch: got %v want %v", end, now)
	}
	wantStart := now.AddDate(0, 0, -cursorInitialBackfillDays)
	if !start.Equal(wantStart) {
		t.Fatalf("full sync start mismatch: got %v want %v", start, wantStart)
	}
}

func TestComputeUsageTimeRangeMsIncrementalLookbackBuffer(t *testing.T) {
	// since is yesterday: without lookback we'd only request ~1 day.
	// With the buffer we look back cursorLookbackDays.
	now := time.Date(2026, 7, 30, 0, 0, 0, 0, time.UTC)
	since := time.Date(2026, 7, 29, 0, 0, 0, 0, time.UTC)
	startMs, endMs := computeUsageTimeRangeMs(&since, now, true)
	start := time.UnixMilli(startMs).UTC()
	end := time.UnixMilli(endMs).UTC()
	if !end.Equal(now) {
		t.Fatalf("end mismatch: got %v want %v", end, now)
	}
	wantStart := now.AddDate(0, 0, -cursorLookbackDays)
	if !start.Equal(wantStart) {
		t.Fatalf("incremental lookback start mismatch: got %v want %v", start, wantStart)
	}
}

func TestComputeUsageTimeRangeMsIncrementalDoesNotExpandPastBackfill(t *testing.T) {
	now := time.Date(2026, 7, 30, 12, 0, 0, 0, time.UTC)
	// since far in the past → start clamped to 90-day min; lookback does not apply.
	since := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	startMs, _ := computeUsageTimeRangeMs(&since, now, true)
	start := time.UnixMilli(startMs).UTC()
	wantStart := now.AddDate(0, 0, -cursorInitialBackfillDays)
	if !start.Equal(wantStart) {
		t.Fatalf("start should stay at 90-day min: got %v want %v", start, wantStart)
	}
}

func TestComputeUsageTimeRangeMsFullSyncUsesSinceWithoutLookback(t *testing.T) {
	now := time.Date(2026, 7, 30, 0, 0, 0, 0, time.UTC)
	since := time.Date(2026, 7, 29, 0, 0, 0, 0, time.UTC)
	startMs, _ := computeUsageTimeRangeMs(&since, now, false)
	start := time.UnixMilli(startMs).UTC()
	if !start.Equal(since) {
		t.Fatalf("full sync should use since without lookback: got %v want %v", start, since)
	}
}

func TestUsageCollectionWindowForIncrementalLookback(t *testing.T) {
	now := time.Date(2026, 7, 30, 0, 0, 0, 0, time.UTC)
	since := time.Date(2026, 7, 29, 0, 0, 0, 0, time.UTC)
	window := usageCollectionWindowFor(&since, now, true)
	wantStart := now.AddDate(0, 0, -cursorLookbackDays).UnixMilli()
	if window.StartMs != wantStart {
		t.Fatalf("start mismatch: got %d want %d", window.StartMs, wantStart)
	}
	if window.EndMs != now.UnixMilli() {
		t.Fatalf("end mismatch: got %d want %d", window.EndMs, now.UnixMilli())
	}
	if window.ChunkCount != 1 {
		t.Fatalf("expected 1 chunk for 7-day lookback, got %d", window.ChunkCount)
	}
}

func TestUsageCollectionWindowForFullSyncBackfill(t *testing.T) {
	now := time.Date(2026, 7, 30, 12, 0, 0, 0, time.UTC)
	window := usageCollectionWindowFor(nil, now, false)
	wantStart := now.AddDate(0, 0, -cursorInitialBackfillDays).UnixMilli()
	if window.StartMs != wantStart {
		t.Fatalf("start mismatch: got %d want %d", window.StartMs, wantStart)
	}
	// Exact 90-day span fits in 3×30-day chunks.
	if window.ChunkCount != 3 {
		t.Fatalf("expected 3 chunks for 90-day backfill, got %d", window.ChunkCount)
	}
}

func TestUsageCollectionWindowForFullSyncUsesSinceWithoutLookback(t *testing.T) {
	now := time.Date(2026, 7, 30, 0, 0, 0, 0, time.UTC)
	since := time.Date(2026, 7, 29, 0, 0, 0, 0, time.UTC)
	window := usageCollectionWindowFor(&since, now, false)
	if window.StartMs != since.UnixMilli() {
		t.Fatalf("full sync should use since without lookback: got %d want %d", window.StartMs, since.UnixMilli())
	}
	if window.ChunkCount != 1 {
		t.Fatalf("expected 1 chunk, got %d", window.ChunkCount)
	}
}
