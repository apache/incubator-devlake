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

	"github.com/stretchr/testify/assert"
)

// Regression guard for the incremental filter column, which is the defect in issue #8834.
//
// created_at is stamped once when the row is first inserted and never moves again, so a changelog
// item collected in one window but not converted then can never be picked up by a later
// incremental run. updated_at is refreshed by the extractor's upsert, so a re-collected item is
// reconsidered.
//
// Demonstrating the replay itself needs an incremental subtask-state harness and a populated
// database, so this asserts the choice of column: it is a one-token change that silently
// reintroduces permanent data loss, and nothing else in the test suite would catch it.
func TestIncrementalFilterUsesUpdatedAtNotCreatedAt(t *testing.T) {
	body, err := os.ReadFile("issue_changelog_convertor.go")
	assert.NoError(t, err)
	source := string(body)

	assert.Contains(t, source, "_tool_jira_issue_changelog_items.updated_at >= ?",
		"the incremental filter must use updated_at")
	assert.NotContains(t, source, "_tool_jira_issue_changelog_items.created_at >= ?",
		"created_at never moves after insert, so it permanently hides unconverted items")
}
