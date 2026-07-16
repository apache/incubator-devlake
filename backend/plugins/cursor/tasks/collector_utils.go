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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/plugins/cursor/models"
)

const (
	rawUsageEventsTable = "cursor_usage_events"
	rawUserSpendTable   = "cursor_user_spend"
	rawMembersTable     = "cursor_members"
	rawDailyUsageTable  = "cursor_daily_usage"

	cursorApiPageSize         = 100
	cursorInitialBackfillDays = 90
	// cursorDailyUsageMaxDays is the maximum span allowed by POST /teams/daily-usage-data.
	cursorDailyUsageMaxDays = 30
)

type cursorRawParams struct {
	ConnectionId uint64
	ScopeId      string
	Endpoint     string
}

func (p cursorRawParams) GetParams() any {
	return p
}

func rawParamsFromTaskData(data *CursorTaskData) cursorRawParams {
	endpoint := models.DefaultEndpoint
	if data.Connection != nil {
		data.Connection.Normalize()
		endpoint = data.Connection.Endpoint
	}
	return cursorRawParams{
		ConnectionId: data.Options.ConnectionId,
		ScopeId:      data.Options.ScopeId,
		Endpoint:     endpoint,
	}
}

type cursorTimeRangeInput struct {
	StartDateMs int64 `json:"startDateMs"`
	EndDateMs   int64 `json:"endDateMs"`
}

type cursorPageState struct {
	Page int `json:"page"`
}

type cursorPagination struct {
	NumPages        int  `json:"numPages"`
	CurrentPage     int  `json:"currentPage"`
	PageSize        int  `json:"pageSize"`
	HasNextPage     bool `json:"hasNextPage"`
	HasPreviousPage bool `json:"hasPreviousPage"`
	Page            int  `json:"page"`
	TotalPages      int  `json:"totalPages"`
}

func computeUsageTimeRangeMs(since *time.Time, now time.Time) (int64, int64) {
	end := now.UTC()
	start := end.AddDate(0, 0, -cursorInitialBackfillDays)
	if since != nil && !since.IsZero() && since.After(start) {
		start = since.UTC()
	}
	return start.UnixMilli(), end.UnixMilli()
}

// splitDailyUsageTimeRangeMs splits [startMs, endMs] into chunks of at most maxDays for
// POST /teams/daily-usage-data (API limit: date range cannot exceed 30 days).
func splitDailyUsageTimeRangeMs(startMs, endMs int64, maxDays int) []cursorTimeRangeInput {
	if startMs >= endMs || maxDays <= 0 {
		return nil
	}
	maxDuration := time.Duration(maxDays) * 24 * time.Hour
	chunkStart := time.UnixMilli(startMs).UTC()
	end := time.UnixMilli(endMs).UTC()

	var chunks []cursorTimeRangeInput
	for chunkStart.Before(end) {
		chunkEnd := chunkStart.Add(maxDuration)
		if chunkEnd.After(end) {
			chunkEnd = end
		}
		chunks = append(chunks, cursorTimeRangeInput{
			StartDateMs: chunkStart.UnixMilli(),
			EndDateMs:   chunkEnd.UnixMilli(),
		})
		if !chunkEnd.Before(end) {
			break
		}
		chunkStart = chunkEnd.Add(time.Millisecond)
	}
	return chunks
}

func parseUsageEventsResponse(res *http.Response) ([]json.RawMessage, errors.Error) {
	body, err := readResponseBody(res)
	if err != nil {
		return nil, err
	}
	var response struct {
		UsageEvents []json.RawMessage `json:"usageEvents"`
	}
	if jsonErr := json.Unmarshal(body, &response); jsonErr != nil {
		return nil, errors.Default.Wrap(errors.Convert(jsonErr), "failed to decode usage events response")
	}
	return response.UsageEvents, nil
}

func parseSpendMembersResponse(res *http.Response) ([]json.RawMessage, errors.Error) {
	body, err := readResponseBody(res)
	if err != nil {
		return nil, err
	}
	var response struct {
		TeamMemberSpend []json.RawMessage `json:"teamMemberSpend"`
	}
	if jsonErr := json.Unmarshal(body, &response); jsonErr != nil {
		return nil, errors.Default.Wrap(errors.Convert(jsonErr), "failed to decode spend response")
	}
	return response.TeamMemberSpend, nil
}

func parseDailyUsageResponse(res *http.Response) ([]json.RawMessage, errors.Error) {
	body, err := readResponseBody(res)
	if err != nil {
		return nil, err
	}
	var response struct {
		Data []json.RawMessage `json:"data"`
	}
	if jsonErr := json.Unmarshal(body, &response); jsonErr != nil {
		return nil, errors.Default.Wrap(errors.Convert(jsonErr), "failed to decode daily usage response")
	}
	return response.Data, nil
}

func parseMembersResponse(res *http.Response) ([]json.RawMessage, errors.Error) {
	body, err := readResponseBody(res)
	if err != nil {
		return nil, err
	}
	var response struct {
		TeamMembers []json.RawMessage `json:"teamMembers"`
	}
	if jsonErr := json.Unmarshal(body, &response); jsonErr != nil {
		return nil, errors.Default.Wrap(errors.Convert(jsonErr), "failed to decode members response")
	}
	return response.TeamMembers, nil
}

func parseSpendMeta(res *http.Response) (int64, errors.Error) {
	body, err := readResponseBody(res)
	if err != nil {
		return 0, err
	}
	var response struct {
		SubscriptionCycleStart int64 `json:"subscriptionCycleStart"`
	}
	if jsonErr := json.Unmarshal(body, &response); jsonErr != nil {
		return 0, errors.Default.Wrap(errors.Convert(jsonErr), "failed to decode spend metadata")
	}
	return response.SubscriptionCycleStart, nil
}

func parsePaginationBody(body []byte) (cursorPagination, errors.Error) {
	var response struct {
		Pagination cursorPagination `json:"pagination"`
	}
	if jsonErr := json.Unmarshal(body, &response); jsonErr != nil {
		return cursorPagination{}, errors.Default.Wrap(errors.Convert(jsonErr), "failed to decode pagination")
	}
	return response.Pagination, nil
}

func parsePagination(res *http.Response) (cursorPagination, errors.Error) {
	body, err := readResponseBody(res)
	if err != nil {
		return cursorPagination{}, err
	}
	return parsePaginationBody(body)
}

func readResponseBody(res *http.Response) ([]byte, errors.Error) {
	if res == nil || res.Body == nil {
		return nil, errors.Default.New("response body is nil")
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Default.Wrap(err, "failed to read response body")
	}
	return body, nil
}

func parseEventTimestampMs(raw string) (time.Time, errors.Error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, errors.BadInput.New("event timestamp is empty")
	}
	ms, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return time.Time{}, errors.BadInput.Wrap(err, "invalid event timestamp")
	}
	return time.UnixMilli(ms).UTC(), nil
}

func computeEventId(timestamp, userEmail, conversationId, model string, chargedCents, requestsCosts float64) string {
	payload := fmt.Sprintf("%s|%s|%s|%s|%v|%v", timestamp, userEmail, conversationId, model, chargedCents, requestsCosts)
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}

func normalizeNullableString(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || strings.EqualFold(value, "null") {
		return ""
	}
	return value
}

func billingCycleTime(ms int64) time.Time {
	if ms <= 0 {
		return time.Time{}
	}
	return time.UnixMilli(ms).UTC()
}
