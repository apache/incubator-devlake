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

func CalculateCommitsDiff(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*RefdiffTaskData)
	repoId := data.Options.RepoId
	db := taskCtx.GetDal()
	ctx := taskCtx.GetContext()
	logger := taskCtx.GetLogger()
	
	projectName := data.Options.ProjectName
	if projectName != "" {
		return nil
	}
	// mysql limit
	insertCountLimitOfCommitsDiff := int(65535 / reflect.ValueOf(code.CommitsDiff{}).NumField())

	commitPairs := data.Options.AllPairs
	commitNodeGraph := utils.NewCommitNodeGraph()

	// load commits from db
	commitParent := &code.CommitParent{}
	cursor, err := db.Cursor(
		dal.Select("cp.*"),
		dal.Join("LEFT JOIN repo_commits rc ON (rc.commit_sha = cp.commit_sha)"),
		dal.From("commit_parents cp"),
		dal.Where("rc.repo_id = ?", repoId),
	)
	if err != nil {
		panic(err)
	}
	defer cursor.Close()

	for cursor.Next() {
		select {
		case <-ctx.Done():
			return errors.Convert(ctx.Err())
		default:
		}
		err = db.Fetch(cursor, commitParent)
		if err != nil {
			return errors.Default.Wrap(err, "failed to read commit from database")
		}
		commitNodeGraph.AddParent(commitParent.CommitSha, commitParent.ParentCommitSha)
	}

	logger.Info("Create a commit node graph with node count[%d]", commitNodeGraph.Size())

	// calculate diffs for commits pairs and store them into database
	commitsDiff := &code.CommitsDiff{}
	finishedCommitDiff := &code.FinishedCommitsDiffs{}
	lenCommitPairs := len(commitPairs)
	taskCtx.SetProgress(0, lenCommitPairs)

	for _, pair := range commitPairs {
		select {
		case <-ctx.Done():
			return errors.Convert(ctx.Err())
		default:
		}
		// ref might advance, keep commit sha for debugging
		commitsDiff.NewCommitSha = pair[0]
		commitsDiff.OldCommitSha = pair[1]
		finishedCommitDiff.NewCommitSha = pair[0]
		finishedCommitDiff.OldCommitSha = pair[1]

		if commitsDiff.NewCommitSha == commitsDiff.OldCommitSha {
			// different refs might point to a same commit, it is ok
			logger.Info(
				"skipping ref pair due to they are the same %s",
				commitsDiff.NewCommitSha,
			)
			continue
		}

		lostSha, oldCount, newCount := commitNodeGraph.CalculateLostSha(pair[1], pair[0])

		commitsDiffs := []code.CommitsDiff{}
		finishedCommitDiffs := []code.FinishedCommitsDiffs{}

		commitsDiff.SortingIndex = 1
		for _, sha := range lostSha {
			commitsDiff.CommitSha = sha
			commitsDiffs = append(commitsDiffs, *commitsDiff)

			// sql limit placeholders count only 65535
			if commitsDiff.SortingIndex%insertCountLimitOfCommitsDiff == 0 {
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
			finishedCommitDiffs = append(finishedCommitDiffs, *finishedCommitDiff)
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
		taskCtx.IncProgress(1)
	}
	return nil
}

var CalculateCommitsDiffMeta = core.SubTaskMeta{
	Name:             "calculateCommitsDiff",
	EntryPoint:       CalculateCommitsDiff,
	EnabledByDefault: true,
	Description:      "Calculate diff commits between refs",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE},
}
