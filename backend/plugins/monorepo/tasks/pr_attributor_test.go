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

	"github.com/stretchr/testify/assert"
)

// TestResolveSubProject exercises the three buckets a pull request can land in during
// attribution: matched by label, unattributed (when enabled), or left unclassified (when
// unattributed rows are disabled). This is the pure decision logic behind
// AttributePullRequests; the actual DB reads/writes around it are covered by the e2e
// dataflow test, since AttributePullRequests no longer contains any other pure logic to
// unit test - the coding/pickup/review/deploy/cycle-time computation that used to live
// here (and was unit tested via firstDeploymentAfter/computeTimeSpan) has been retired in
// favor of updateProjectPrMetricsSubProject sourcing those numbers from DORA.
func TestResolveSubProject(t *testing.T) {
	cases := []struct {
		name                string
		labelMatch          string
		includeUnattributed bool
		wantSubProject      string
		wantMatched         bool
		wantUnattributed    bool
		wantSkipped         bool
	}{
		{
			name:                "label match wins regardless of includeUnattributed",
			labelMatch:          "serviceA",
			includeUnattributed: false,
			wantSubProject:      "serviceA",
			wantMatched:         true,
		},
		{
			name:                "no match, unattributed enabled",
			labelMatch:          "",
			includeUnattributed: true,
			wantSubProject:      UnattributedSubProject,
			wantUnattributed:    true,
		},
		{
			name:                "no match, unattributed disabled leaves it unclassified",
			labelMatch:          "",
			includeUnattributed: false,
			wantSubProject:      "",
			wantSkipped:         true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			subProject, matched, unattributed, skipped := resolveSubProject(tc.labelMatch, tc.includeUnattributed)
			assert.Equal(t, tc.wantSubProject, subProject)
			assert.Equal(t, tc.wantMatched, matched)
			assert.Equal(t, tc.wantUnattributed, unattributed)
			assert.Equal(t, tc.wantSkipped, skipped)
		})
	}
}
