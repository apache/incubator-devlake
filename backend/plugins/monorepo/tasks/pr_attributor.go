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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/plugin"
)

var AttributePullRequestsMeta = plugin.SubTaskMeta{
	Name:             "attributePullRequests",
	EntryPoint:       AttributePullRequests,
	EnabledByDefault: true,
	Description:      "Attribute all pull requests (open, closed and merged) to monorepo sub-projects by label",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
}

type prIdRow struct {
	Id string
}

type prLabelRow struct {
	PullRequestId string
	LabelName     string
}

// AttributePullRequests is attribution-only: it tags pull_requests.sub_project (and, from
// there, pull_request_commits.sub_project) for every pull request of the project,
// regardless of merge status. It intentionally does not compute coding/pickup/review/
// deploy/cycle time - that responsibility belongs to updateProjectPrMetricsSubProject,
// which sources those numbers from DORA's project_pr_metrics instead of recomputing them.
//
// This also removes the historical `pr.merged_date IS NOT NULL` filter: open and closed
// (but unmerged) pull requests are attributed too, so a monorepo's PR-volume dashboards
// are not silently missing everything that hasn't merged yet.
func AttributePullRequests(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()
	data := taskCtx.GetData().(*MonorepoTaskData)
	projectName := data.Options.ProjectName
	includeUnattributed := data.Options.ShouldIncludeUnattributed()

	labelsByPr, err := loadPrLabels(db, projectName)
	if err != nil {
		return err
	}

	clauses := []dal.Clause{
		dal.Select("pr.id"),
		dal.From("pull_requests pr"),
		dal.Join("JOIN project_mapping pm ON (pm.table = 'repos' AND pm.row_id = pr.base_repo_id)"),
		dal.Where("pm.project_name = ?", projectName),
	}
	var rows []prIdRow
	if err := db.All(&rows, clauses...); err != nil {
		return errors.Default.Wrap(err, "error loading pull requests to attribute")
	}

	matchedCount, unattributedCount, skippedCount := 0, 0, 0
	for _, row := range rows {
		labelMatch := data.Matcher.MatchPrLabels(labelsByPr[row.Id])
		subProject, matched, unattributed, skipped := resolveSubProject(labelMatch, includeUnattributed)
		if matched {
			matchedCount++
		}
		if unattributed {
			unattributedCount++
		}
		if skipped {
			skippedCount++
		}

		var value interface{}
		if subProject != "" {
			value = subProject
		}
		pr := &code.PullRequest{}
		pr.Id = row.Id
		if err := db.UpdateColumn(pr, "sub_project", value); err != nil {
			return errors.Default.Wrap(err, "error updating pull_requests.sub_project")
		}
	}
	logger.Info("monorepo: attributed %d pull requests (%d matched, %d unattributed, %d left unclassified)",
		len(rows), matchedCount, unattributedCount, skippedCount)

	// Propagate to pull_request_commits in one set-based statement. Written as a
	// correlated subquery / IN-subquery rather than a three-way UPDATE...JOIN so the same
	// SQL works on both MySQL and PostgreSQL without a dialect branch.
	if err := db.Exec(`
		UPDATE pull_request_commits
		SET sub_project = (
			SELECT pr.sub_project FROM pull_requests pr
			WHERE pr.id = pull_request_commits.pull_request_id
		)
		WHERE pull_request_id IN (
			SELECT pr2.id FROM pull_requests pr2
			JOIN project_mapping pm ON (pm.table = 'repos' AND pm.row_id = pr2.base_repo_id)
			WHERE pm.project_name = ?
		)
	`, projectName); err != nil {
		return errors.Default.Wrap(err, "error updating pull_request_commits.sub_project")
	}

	return nil
}

// resolveSubProject decides the sub_project value to write for a pull request, given the
// (possibly empty) result of matching its labels and whether unattributed rows are
// enabled. An empty returned subProject means "write NULL" (leave/clear unclassified).
func resolveSubProject(labelMatch string, includeUnattributed bool) (subProject string, matched, unattributed, skipped bool) {
	switch {
	case labelMatch != "":
		return labelMatch, true, false, false
	case includeUnattributed:
		return UnattributedSubProject, false, true, false
	default:
		return "", false, false, true
	}
}

func loadPrLabels(db dal.Dal, projectName string) (map[string][]string, errors.Error) {
	var rows []prLabelRow
	err := db.All(&rows,
		dal.Select("prl.pull_request_id, prl.label_name"),
		dal.From("pull_request_labels prl"),
		dal.Join("JOIN pull_requests pr ON (pr.id = prl.pull_request_id)"),
		dal.Join("JOIN project_mapping pm ON (pm.table = 'repos' AND pm.row_id = pr.base_repo_id)"),
		dal.Where("pm.project_name = ?", projectName),
	)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error loading pull request labels")
	}
	byPr := make(map[string][]string)
	for _, r := range rows {
		byPr[r.PullRequestId] = append(byPr[r.PullRequestId], r.LabelName)
	}
	return byPr, nil
}
