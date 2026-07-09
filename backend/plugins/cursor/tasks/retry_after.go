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
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
)

func handleCursorRetryAfter(res *http.Response, logger log.Logger, now nowFunc, sleep sleepFunc) errors.Error {
	if res == nil || res.StatusCode != http.StatusTooManyRequests {
		return nil
	}
	if now == nil {
		now = time.Now
	}
	if sleep == nil {
		sleep = time.Sleep
	}
	wait := parseRetryAfter(res.Header.Get("Retry-After"), now().UTC())
	if wait > 0 && logger != nil {
		logger.Warn(nil, "Cursor returned 429; sleeping %s per Retry-After", wait.String())
	}
	if wait > 0 {
		sleep(wait)
	}
	return errors.HttpStatus(http.StatusTooManyRequests).New("Cursor rate limited the request")
}

func parseRetryAfter(value string, now time.Time) time.Duration {
	value = trimSpace(value)
	if value == "" {
		return 0
	}
	if seconds, err := strconv.Atoi(value); err == nil {
		return time.Duration(seconds) * time.Second
	}
	if retryAt, err := http.ParseTime(value); err == nil {
		wait := retryAt.Sub(now)
		if wait > 0 {
			return wait
		}
	}
	return 0
}

func trimSpace(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

type nowFunc func() time.Time
type sleepFunc func(time.Duration)
