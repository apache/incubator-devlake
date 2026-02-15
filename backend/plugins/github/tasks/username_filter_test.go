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
	"testing"
)

func TestShouldSkipByUsername_EmptyList(t *testing.T) {
	resetExcludedUsernamesForTest()
	os.Setenv("GITHUB_PR_EXCLUDELIST", "")
	defer os.Unsetenv("GITHUB_PR_EXCLUDELIST")

	if shouldSkipByUsername("renovate[bot]") {
		t.Error("Expected false for any username when excludelist is empty")
	}
}

func TestShouldSkipByUsername_SingleMatch(t *testing.T) {
	resetExcludedUsernamesForTest()
	os.Setenv("GITHUB_PR_EXCLUDELIST", "renovate[bot]")
	defer os.Unsetenv("GITHUB_PR_EXCLUDELIST")

	if !shouldSkipByUsername("renovate[bot]") {
		t.Error("Expected true for renovate[bot]")
	}
	if shouldSkipByUsername("human-user") {
		t.Error("Expected false for human-user")
	}
}

func TestShouldSkipByUsername_MultipleUsernames(t *testing.T) {
	resetExcludedUsernamesForTest()
	os.Setenv("GITHUB_PR_EXCLUDELIST", "renovate[bot],dependabot[bot],github-actions[bot]")
	defer os.Unsetenv("GITHUB_PR_EXCLUDELIST")

	if !shouldSkipByUsername("renovate[bot]") {
		t.Error("Expected true for renovate[bot]")
	}
	if !shouldSkipByUsername("dependabot[bot]") {
		t.Error("Expected true for dependabot[bot]")
	}
	if !shouldSkipByUsername("github-actions[bot]") {
		t.Error("Expected true for github-actions[bot]")
	}
	if shouldSkipByUsername("human-user") {
		t.Error("Expected false for human-user")
	}
}

func TestShouldSkipByUsername_CaseInsensitive(t *testing.T) {
	resetExcludedUsernamesForTest()
	os.Setenv("GITHUB_PR_EXCLUDELIST", "renovate[bot]")
	defer os.Unsetenv("GITHUB_PR_EXCLUDELIST")

	if !shouldSkipByUsername("Renovate[bot]") {
		t.Error("Expected true for Renovate[bot] (case insensitive)")
	}
	if !shouldSkipByUsername("RENOVATE[BOT]") {
		t.Error("Expected true for RENOVATE[BOT] (case insensitive)")
	}
}

func TestShouldSkipByUsername_WhitespaceTrimming(t *testing.T) {
	resetExcludedUsernamesForTest()
	os.Setenv("GITHUB_PR_EXCLUDELIST", " renovate[bot] , dependabot[bot] ")
	defer os.Unsetenv("GITHUB_PR_EXCLUDELIST")

	if !shouldSkipByUsername("renovate[bot]") {
		t.Error("Expected true for renovate[bot] with whitespace in config")
	}
	if !shouldSkipByUsername("dependabot[bot]") {
		t.Error("Expected true for dependabot[bot] with whitespace in config")
	}
}

func TestShouldSkipByUsername_EmptyUsername(t *testing.T) {
	resetExcludedUsernamesForTest()
	os.Setenv("GITHUB_PR_EXCLUDELIST", "renovate[bot]")
	defer os.Unsetenv("GITHUB_PR_EXCLUDELIST")

	if shouldSkipByUsername("") {
		t.Error("Expected false for empty username")
	}
}
