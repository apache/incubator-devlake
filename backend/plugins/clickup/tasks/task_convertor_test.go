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
	"testing"

	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
)

func TestListTypeFor(t *testing.T) {
	bug := regexp.MustCompile(`(?i)bug`)
	incident := regexp.MustCompile(`(?i)incident`)

	tests := []struct {
		name     string
		listName string
		bugRe    *regexp.Regexp
		incRe    *regexp.Regexp
		want     string
	}{
		{"qa bugs list -> BUG", "QA Bugs", bug, incident, ticket.BUG},
		{"case-insensitive bug", "bug tracking", bug, incident, ticket.BUG},
		{"incident list -> INCIDENT", "Incident Backlog", bug, incident, ticket.INCIDENT},
		{"incident beats bug when both match", "Bug/Incident triage", bug, incident, ticket.INCIDENT},
		{"backlog matches neither -> empty", "Backlog", bug, incident, ""},
		{"sprint matches neither -> empty", "v4.3.0 Sprint 40", bug, incident, ""},
		{"nil bug regex, no incident -> empty", "QA Bugs", nil, nil, ""},
		{"only incident regex set", "Incident Backlog", nil, incident, ticket.INCIDENT},
		{"only bug regex set, name is incident -> empty", "Incident Backlog", bug, nil, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := listTypeFor(tt.listName, tt.bugRe, tt.incRe); got != tt.want {
				t.Fatalf("listTypeFor(%q) = %q, want %q", tt.listName, got, tt.want)
			}
		})
	}
}
