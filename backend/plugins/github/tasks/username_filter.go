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
	"os"
	"strings"
	"sync"
)

var (
	excludedUsernames     []string
	excludedUsernamesOnce sync.Once
	excludedUsernamesMu   sync.RWMutex
)

// initExcludedUsernames reads and parses the GITHUB_PR_EXCLUDELIST environment variable
func initExcludedUsernames() {
	excludedUsernamesOnce.Do(func() {
		loadExcludedUsernames()
	})
}

// loadExcludedUsernames parses the environment variable (called by initExcludedUsernames or tests)
func loadExcludedUsernames() {
	excludedUsernamesMu.Lock()
	defer excludedUsernamesMu.Unlock()

	envValue := os.Getenv("GITHUB_PR_EXCLUDELIST")
	if envValue == "" {
		excludedUsernames = []string{}
		return
	}

	usernames := strings.Split(envValue, ",")
	excludedUsernames = make([]string, 0, len(usernames))
	for _, username := range usernames {
		trimmed := strings.TrimSpace(username)
		if trimmed != "" {
			excludedUsernames = append(excludedUsernames, strings.ToLower(trimmed))
		}
	}
}

// resetExcludedUsernamesForTest resets the cache for testing purposes
func resetExcludedUsernamesForTest() {
	excludedUsernamesMu.Lock()
	defer excludedUsernamesMu.Unlock()
	excludedUsernames = nil
	excludedUsernamesOnce = sync.Once{}
}

// shouldSkipByUsername checks if the given username should be filtered out
// Returns true if the username matches any entry in the GITHUB_PR_EXCLUDELIST
func shouldSkipByUsername(username string) bool {
	initExcludedUsernames()

	if username == "" {
		return false
	}

	excludedUsernamesMu.RLock()
	defer excludedUsernamesMu.RUnlock()

	lowerUsername := strings.ToLower(username)
	for _, excluded := range excludedUsernames {
		if lowerUsername == excluded {
			return true
		}
	}
	return false
}
