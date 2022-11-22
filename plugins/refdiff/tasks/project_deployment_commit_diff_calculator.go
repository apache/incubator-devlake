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
	"reflect"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/refdiff/utils"
)

func CalculateProjectDeploymentCommitsDiff(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*RefdiffTaskData)
	db := taskCtx.GetDal()
	ctx := taskCtx.GetContext()
	logger := taskCtx.GetLogger()

	projectName := data.Options.ProjectName
	if projectName == "" {
		return nil
	}

	cursorScope, err := db.Cursor(
		dal.Select("row_id"),
		dal.From("project_mapping"),
		dal.Where("project_name = ?", projectName),
	)
	if err != nil {
		return err
	}
	defer cursorScope.Close()

	var ExistFinishedCommitDiff []code.FinishedCommitsDiffs
	err = db.All(&ExistFinishedCommitDiff,
		dal.Select("*"),
		dal.From("finished_commits_diffs"),
	)
	if err != nil {
		return err
	}

	for cursorScope.Next() {
		var scopeId string
		err = errors.Convert(cursorScope.Scan(&scopeId))
		if err != nil {
			return err
		}

		var commitShaList []string
		err := db.All(&commitShaList,
			dal.Select("commit_sha"),
			dal.From("cicd_tasks ct"),
			dal.Join("left join cicd_pipelines cp on cp.id = ct.pipeline_id"),
			dal.Join("left join cicd_pipeline_commits cpc on cpc.pipeline_id = cp.id"),
			dal.Where("ct.type = ? and commit_sha != ? and repo_id=? ", "DEPLOYMENT", "", scopeId),
			dal.Orderby("ct.started_date"),
		)
		if err != nil {
			return err
		}

		var commitPairs []code.CommitsDiff
		var finishedCommitDiffs []code.FinishedCommitsDiffs

		for i := 0; i < len(commitShaList)-1; i++ {
			for _, item := range ExistFinishedCommitDiff {
				if commitShaList[i+1] == item.NewCommitSha && commitShaList[i] == item.OldCommitSha {
					i++
				}
			}
			commitPairs = append(commitPairs, code.CommitsDiff{NewCommitSha: commitShaList[i+1], OldCommitSha: commitShaList[i]})
			finishedCommitDiffs = append(finishedCommitDiffs, code.FinishedCommitsDiffs{NewCommitSha: commitShaList[i+1], OldCommitSha: commitShaList[i]})
		}

		insertCountLimitOfDeployCommitsDiff := int(65535 / reflect.ValueOf(code.CommitsDiff{}).NumField())
		commitNodeGraph := utils.NewCommitNodeGraph()

		var CommitParentList []code.CommitParent
		err = db.All(&CommitParentList,
			dal.Select("cp.*"),
			dal.Join("LEFT JOIN repo_commits rc ON (rc.commit_sha = cp.commit_sha)"),
			dal.From("commit_parents cp"),
			dal.Where("rc.repo_id = ?", scopeId),
		)
		if err != nil {
			return err
		}

		for i := 0; i < len(CommitParentList); i++ {
			commitNodeGraph.AddParent(CommitParentList[i].CommitSha, CommitParentList[i].ParentCommitSha)
		}
		logger.Info("Create a commit node graph with node count[%d]", commitNodeGraph.Size())

		// calculate diffs for commits pairs and store them into database
		commitsDiff := &code.CommitsDiff{}
		lenCommitPairs := len(commitPairs)
		taskCtx.SetProgress(0, lenCommitPairs)

		for _, pair := range commitPairs {
			select {
			case <-ctx.Done():
				return errors.Convert(ctx.Err())
			default:
			}

			commitsDiff.NewCommitSha = pair.NewCommitSha
			commitsDiff.OldCommitSha = pair.OldCommitSha

			if commitsDiff.NewCommitSha == commitsDiff.OldCommitSha {
				// different deploy might point to a same commit, it is ok
				logger.Info(
					"skipping ref pair due to they are the same %s",
					commitsDiff.NewCommitSha,
				)
				continue
			}

			lostSha, oldCount, newCount := commitNodeGraph.CalculateLostSha(commitsDiff.OldCommitSha, commitsDiff.NewCommitSha)

			commitsDiffs := []code.CommitsDiff{}
			commitsDiff.SortingIndex = 1
			for _, sha := range lostSha {
				commitsDiff.CommitSha = sha
				commitsDiffs = append(commitsDiffs, *commitsDiff)

				// sql limit placeholders count only 65535
				if commitsDiff.SortingIndex%insertCountLimitOfDeployCommitsDiff == 0 {
					logger.Info("commitsDiffs count in limited[%d] index[%d]--exec and clean", len(commitsDiffs), commitsDiff.SortingIndex)
					err = db.CreateIfNotExist(commitsDiffs)
					if err != nil {
						return err
					}
					commitsDiffs = []code.CommitsDiff{}
				}

				commitsDiff.SortingIndex++
			}

			if len(commitsDiffs) > 0 {
				logger.Info("insert data count [%d]", len(commitsDiffs))
				err = db.CreateIfNotExist(commitsDiffs)
				if err != nil {
					return err
				}
				err = db.CreateIfNotExist(finishedCommitDiffs)
				if err != nil {
					return err
				}
			}

			logger.Info(
				"total %d commits of difference found between [new][%s] and [old][%s(total:%d)]",
				newCount,
				commitsDiff.NewCommitSha,
				commitsDiff.OldCommitSha,
				oldCount,
			)

		}

	}
	taskCtx.IncProgress(1)
	return nil
}

var CalculateProjectDeploymentCommitsDiffMeta = core.SubTaskMeta{
	Name:             "calculateProjectDeploymentCommitsDiff",
	EntryPoint:       CalculateProjectDeploymentCommitsDiff,
	EnabledByDefault: true,
	Description:      "Calculate diff commits between project deployments",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE},
}
