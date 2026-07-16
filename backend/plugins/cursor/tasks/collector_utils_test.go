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

	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/cursor/models"
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
