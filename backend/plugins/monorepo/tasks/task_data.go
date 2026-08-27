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
	"fmt"
	"regexp"

	"github.com/apache/incubator-devlake/core/errors"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type MonorepoApiParams struct {
	ProjectName string
}

const (
	// UnattributedSubProject is the sentinel sub_project value written for PRs/deployments
	// that belong to a monorepo project (i.e. one with SubProjects configured) but matched
	// none of the configured sub-projects. Whether it is written at all is controlled by
	// MonorepoOptions.IncludeUnattributed.
	UnattributedSubProject = "unattributed"
	// AllSubProjectsLabel is the dashboard-side label shown for rows with no sub_project at
	// all (single-repo projects, or rows the monorepo plugin has not processed). It is never
	// written to the database — dashboards derive it via COALESCE(sub_project, 'All') — but
	// it is reserved here too so it cannot be configured as a real sub-project name and
	// collide with that convention.
	AllSubProjectsLabel = "All"
)

// SubProjectConfig declares one logical project living inside a monorepo.
type SubProjectConfig struct {
	// Name identifies the sub-project in the output tables and dashboards.
	Name string `json:"name" mapstructure:"name"`
	// PrLabels are the pull request labels that mark a PR as belonging to this
	// sub-project. Matching is exact and case-sensitive.
	PrLabels []string `json:"prLabels" mapstructure:"prLabels"`
	// DeployJobPattern is a regular expression matched against cicd_tasks.name to
	// recognise this sub-project's deployment jobs, e.g. "^deploy-serviceA$".
	DeployJobPattern string `json:"deployJobPattern" mapstructure:"deployJobPattern"`
}

type MonorepoOptions struct {
	ProjectName string `json:"projectName" mapstructure:"projectName"`
	// SubProjects is ordered: when a pull request carries the labels of more than one
	// sub-project, the earliest entry in this list wins.
	SubProjects []SubProjectConfig `json:"subProjects" mapstructure:"subProjects"`
	// IncludeUnattributed controls whether PRs/deployments that belong to this monorepo
	// project but matched none of the configured sub-projects get sub_project =
	// UnattributedSubProject (true, the default) or are left unclassified / skipped
	// entirely (false, the pre-existing behaviour). A pointer so decoding can distinguish
	// "the caller didn't set this" (nil, defaults to true) from an explicit false.
	IncludeUnattributed *bool `json:"includeUnattributed" mapstructure:"includeUnattributed"`
}

// ShouldIncludeUnattributed returns the effective value of IncludeUnattributed, defaulting
// to true when the option was not set.
func (op *MonorepoOptions) ShouldIncludeUnattributed() bool {
	return op.IncludeUnattributed == nil || *op.IncludeUnattributed
}

type MonorepoTaskData struct {
	Options *MonorepoOptions
	Matcher *SubProjectMatcher
}

// SubProjectMatcher resolves deployments and pull requests to sub-projects. It holds
// the compiled form of the configuration so the regexes are built once per task rather
// than once per row.
type SubProjectMatcher struct {
	names        []string
	prLabels     []map[string]struct{}
	deployJobRes []*regexp.Regexp
}

// NewSubProjectMatcher compiles the sub-project configuration, validating it along the way.
func NewSubProjectMatcher(subProjects []SubProjectConfig) (*SubProjectMatcher, errors.Error) {
	m := &SubProjectMatcher{
		names:        make([]string, 0, len(subProjects)),
		prLabels:     make([]map[string]struct{}, 0, len(subProjects)),
		deployJobRes: make([]*regexp.Regexp, 0, len(subProjects)),
	}
	seen := make(map[string]struct{}, len(subProjects))
	for i, sp := range subProjects {
		if sp.Name == "" {
			return nil, errors.BadInput.New(fmt.Sprintf("subProjects[%d]: name is required", i))
		}
		if sp.Name == UnattributedSubProject || sp.Name == AllSubProjectsLabel {
			return nil, errors.BadInput.New(fmt.Sprintf(
				"subProjects[%d]: name %q is reserved and cannot be used as a sub-project name", i, sp.Name))
		}
		if _, dup := seen[sp.Name]; dup {
			return nil, errors.BadInput.New(fmt.Sprintf("subProjects[%d]: duplicate name %q", i, sp.Name))
		}
		seen[sp.Name] = struct{}{}

		var jobRe *regexp.Regexp
		if sp.DeployJobPattern != "" {
			compiled, err := regexp.Compile(sp.DeployJobPattern)
			if err != nil {
				return nil, errors.BadInput.Wrap(err, fmt.Sprintf(
					"subProjects[%d] (%s): invalid deployJobPattern", i, sp.Name))
			}
			jobRe = compiled
		}

		labels := make(map[string]struct{}, len(sp.PrLabels))
		for j, l := range sp.PrLabels {
			if l == "" {
				return nil, errors.BadInput.New(fmt.Sprintf(
					"subProjects[%d] (%s): prLabels[%d] must not be empty", i, sp.Name, j))
			}
			labels[l] = struct{}{}
		}

		m.names = append(m.names, sp.Name)
		m.prLabels = append(m.prLabels, labels)
		m.deployJobRes = append(m.deployJobRes, jobRe)
	}
	return m, nil
}

// MatchDeployJob returns every sub-project whose DeployJobPattern matches jobName.
//
// More than one match is possible and is reported faithfully: a single pipeline running
// both deploy-serviceA and deploy-serviceB genuinely deploys two sub-projects. If a
// single job name matches two patterns, that indicates overlapping configuration.
func (m *SubProjectMatcher) MatchDeployJob(jobName string) []string {
	var matched []string
	for i, re := range m.deployJobRes {
		if re != nil && re.MatchString(jobName) {
			matched = append(matched, m.names[i])
		}
	}
	return matched
}

// MatchPrLabels returns the single sub-project a pull request belongs to, or "" when no
// sub-project claims it. When several sub-projects match, the earliest one in the
// configured order wins — labels carry no size signal that could rank them otherwise.
func (m *SubProjectMatcher) MatchPrLabels(labels []string) string {
	if len(labels) == 0 {
		return ""
	}
	present := make(map[string]struct{}, len(labels))
	for _, l := range labels {
		present[l] = struct{}{}
	}
	for i, wanted := range m.prLabels {
		for l := range wanted {
			if _, ok := present[l]; ok {
				return m.names[i]
			}
		}
	}
	return ""
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*MonorepoOptions, errors.Error) {
	var op MonorepoOptions
	if err := helper.Decode(options, &op, nil); err != nil {
		return nil, errors.Default.Wrap(err, "error decoding monorepo task options")
	}
	if op.ProjectName == "" {
		return nil, errors.BadInput.New("projectName is required for the monorepo plugin")
	}
	if len(op.SubProjects) == 0 {
		return nil, errors.BadInput.New("at least one entry in subProjects is required for the monorepo plugin")
	}
	if op.IncludeUnattributed == nil {
		defaultTrue := true
		op.IncludeUnattributed = &defaultTrue
	}
	return &op, nil
}
