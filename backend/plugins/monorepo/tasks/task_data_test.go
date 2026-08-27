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

// twoServices is the canonical monorepo configuration used across these tests:
// serviceA is declared first, so it wins any tie.
func twoServices() []SubProjectConfig {
	return []SubProjectConfig{
		{
			Name:             "serviceA",
			PrLabels:         []string{"serviceA"},
			DeployJobPattern: "^deploy-serviceA$",
		},
		{
			Name:             "serviceB",
			PrLabels:         []string{"serviceB", "svc-b"},
			DeployJobPattern: "^deploy-serviceB$",
		},
	}
}

func TestMatchDeployJob(t *testing.T) {
	matcher, err := NewSubProjectMatcher(twoServices())
	assert.Nil(t, err)

	cases := []struct {
		name     string
		jobName  string
		expected []string
	}{
		{"matches serviceA", "deploy-serviceA", []string{"serviceA"}},
		{"matches serviceB", "deploy-serviceB", []string{"serviceB"}},
		{"build job is not a deployment", "build-serviceA", nil},
		{"unrelated job matches nothing", "run-tests", nil},
		{"anchored pattern rejects a superstring", "deploy-serviceAB", nil},
		{"empty job name matches nothing", "", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, matcher.MatchDeployJob(tc.jobName))
		})
	}
}

// A pipeline that runs both services' deploy jobs produces a row for each. The two jobs
// arrive as separate rows, so each resolves to exactly one sub-project.
func TestMatchDeployJob_PipelineDeployingBothServices(t *testing.T) {
	matcher, err := NewSubProjectMatcher(twoServices())
	assert.Nil(t, err)

	assert.Equal(t, []string{"serviceA"}, matcher.MatchDeployJob("deploy-serviceA"))
	assert.Equal(t, []string{"serviceB"}, matcher.MatchDeployJob("deploy-serviceB"))
}

// Overlapping patterns are reported faithfully rather than silently resolved, so a
// misconfiguration is visible in the data instead of hidden.
func TestMatchDeployJob_OverlappingPatterns(t *testing.T) {
	matcher, err := NewSubProjectMatcher([]SubProjectConfig{
		{Name: "serviceA", DeployJobPattern: "deploy"},
		{Name: "serviceB", DeployJobPattern: "^deploy-serviceB$"},
	})
	assert.Nil(t, err)

	assert.Equal(t, []string{"serviceA", "serviceB"}, matcher.MatchDeployJob("deploy-serviceB"))
}

func TestMatchDeployJob_NoPatternNeverMatches(t *testing.T) {
	matcher, err := NewSubProjectMatcher([]SubProjectConfig{
		{Name: "labelsOnly", PrLabels: []string{"labelsOnly"}},
	})
	assert.Nil(t, err)

	assert.Nil(t, matcher.MatchDeployJob("deploy-labelsOnly"))
}

func TestMatchPrLabels(t *testing.T) {
	matcher, err := NewSubProjectMatcher(twoServices())
	assert.Nil(t, err)

	cases := []struct {
		name     string
		labels   []string
		expected string
	}{
		{"single matching label", []string{"serviceA"}, "serviceA"},
		{"alias label resolves to its sub-project", []string{"svc-b"}, "serviceB"},
		{"matching label among unrelated ones", []string{"bug", "serviceB", "urgent"}, "serviceB"},
		{"no matching label", []string{"bug", "urgent"}, ""},
		{"no labels at all", nil, ""},
		{"empty label slice", []string{}, ""},
		{"matching is case sensitive", []string{"servicea"}, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, matcher.MatchPrLabels(tc.labels))
		})
	}
}

// A PR labelled for several sub-projects is assigned to exactly one: the earliest in
// configuration order. Labels carry no size signal, so declaration order is the tie-break.
func TestMatchPrLabels_TieBreakIsConfigOrder(t *testing.T) {
	both := []string{"serviceB", "serviceA"}

	matcher, err := NewSubProjectMatcher(twoServices())
	assert.Nil(t, err)
	assert.Equal(t, "serviceA", matcher.MatchPrLabels(both))

	// Reversing the configuration reverses the winner, proving order drives the result
	// rather than the order of labels on the PR.
	reversed := []SubProjectConfig{twoServices()[1], twoServices()[0]}
	reversedMatcher, err := NewSubProjectMatcher(reversed)
	assert.Nil(t, err)
	assert.Equal(t, "serviceB", reversedMatcher.MatchPrLabels(both))
}

func TestNewSubProjectMatcher_Validation(t *testing.T) {
	cases := []struct {
		name        string
		subProjects []SubProjectConfig
		expectErr   bool
	}{
		{
			name:        "valid configuration",
			subProjects: twoServices(),
		},
		{
			name:        "empty configuration is allowed here, rejected by option decoding",
			subProjects: nil,
		},
		{
			name:        "missing name",
			subProjects: []SubProjectConfig{{PrLabels: []string{"x"}}},
			expectErr:   true,
		},
		{
			name: "duplicate names",
			subProjects: []SubProjectConfig{
				{Name: "serviceA", DeployJobPattern: "^a$"},
				{Name: "serviceA", DeployJobPattern: "^b$"},
			},
			expectErr: true,
		},
		{
			name:        "invalid deploy job regex",
			subProjects: []SubProjectConfig{{Name: "serviceA", DeployJobPattern: "^deploy-(unclosed"}},
			expectErr:   true,
		},
		{
			name:        "name 'unattributed' collides with the sentinel and is rejected",
			subProjects: []SubProjectConfig{{Name: "unattributed", DeployJobPattern: "^deploy-x$"}},
			expectErr:   true,
		},
		{
			name:        "name 'All' collides with the dashboard-side sentinel and is rejected",
			subProjects: []SubProjectConfig{{Name: "All", DeployJobPattern: "^deploy-x$"}},
			expectErr:   true,
		},
		{
			name: "empty prLabels entry is rejected",
			subProjects: []SubProjectConfig{
				{Name: "serviceA", PrLabels: []string{"serviceA", ""}},
			},
			expectErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			matcher, err := NewSubProjectMatcher(tc.subProjects)
			if tc.expectErr {
				assert.NotNil(t, err)
				assert.Nil(t, matcher)
				return
			}
			assert.Nil(t, err)
			assert.NotNil(t, matcher)
		})
	}
}

func TestDecodeAndValidateTaskOptions(t *testing.T) {
	t.Run("valid options", func(t *testing.T) {
		op, err := DecodeAndValidateTaskOptions(map[string]interface{}{
			"projectName": "monorepo",
			"subProjects": []interface{}{
				map[string]interface{}{
					"name":             "serviceA",
					"prLabels":         []interface{}{"serviceA"},
					"deployJobPattern": "^deploy-serviceA$",
				},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, "monorepo", op.ProjectName)
		assert.Len(t, op.SubProjects, 1)
		assert.Equal(t, "serviceA", op.SubProjects[0].Name)
		assert.Equal(t, []string{"serviceA"}, op.SubProjects[0].PrLabels)
		assert.Equal(t, "^deploy-serviceA$", op.SubProjects[0].DeployJobPattern)
		assert.True(t, op.ShouldIncludeUnattributed(),
			"includeUnattributed must default to true when the caller doesn't set it")
	})

	t.Run("includeUnattributed explicit false is preserved", func(t *testing.T) {
		op, err := DecodeAndValidateTaskOptions(map[string]interface{}{
			"projectName": "monorepo",
			"subProjects": []interface{}{
				map[string]interface{}{"name": "serviceA"},
			},
			"includeUnattributed": false,
		})
		assert.Nil(t, err)
		assert.False(t, op.ShouldIncludeUnattributed())
	})

	t.Run("missing projectName is rejected", func(t *testing.T) {
		_, err := DecodeAndValidateTaskOptions(map[string]interface{}{
			"subProjects": []interface{}{
				map[string]interface{}{"name": "serviceA"},
			},
		})
		assert.NotNil(t, err)
	})

	t.Run("missing subProjects is rejected", func(t *testing.T) {
		_, err := DecodeAndValidateTaskOptions(map[string]interface{}{
			"projectName": "monorepo",
		})
		assert.NotNil(t, err)
	})
}
