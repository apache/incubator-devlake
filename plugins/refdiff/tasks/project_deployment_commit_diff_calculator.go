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
	"time"

	"github.com/apache/incubator-devlake/errors"

	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/refdiff/utils"
)

type DeploymentCommitPairRev struct {
	CommitSha   string `gorm:"type:varchar(40)"`
	PipelineId  string `gorm:"type:varchar(255)"`
	TaskId      string `gorm:"type:varchar(255)"`
	TaskName    string `gorm:"type:varchar(255)"`
	StartedDate *time.Time
}

func CalculateProjectDeploymentCommitsDiff(taskCtx core.SubTaskContext) errors.Error {

	// 1. select scopeId from project_mapping where project_name = $projectName
	// 2. for repoId in scopeList:
	/*
	   select
	   commit_sha, cp.id as pipeline_id, ct.id as task_id, ct.started_date
	   FROM
	   cicd_tasks ct
	   left join cicd_pipelines cp on cp.id = ct.pipeline_id
	   left join cicd_pipeline_commits cpc on cpc.pipeline_id = cp.id
	   where
	   ct.type = "DEPLOYMENT" and commit_sha != "" and repo_id=$repoId  -- github:GithubRepo:1:484251804
	   order by
	   ct.started_date
	*/

	// 3, 根据新旧commit_sha 之间，计算的commit_sha
	// 4, 写表到deploy_commits_diff

	data := taskCtx.GetData().(*RefdiffTaskData)
	db := taskCtx.GetDal()
	ctx := taskCtx.GetContext()
	logger := taskCtx.GetLogger()

	projectName := data.Options.ProjectName
	cursorScope, err := db.Cursor(
		dal.Select("row_id"),
		dal.From("project_mapping"),
		dal.Where("project_name = ?", projectName),
	)
	if err != nil {
		panic(err)
	}
	defer cursorScope.Close()
	for cursorScope.Next() {
		var scopeId string
		err = errors.Convert(cursorScope.Scan(&scopeId))
		if err != nil {
			return err
		}

		cursorDeployment, err := db.Cursor(
			dal.Select("commit_sha, cp.id as pipeline_id, ct.id as task_id, ct.name as task_name, ct.started_date"),
			dal.From("cicd_tasks ct"),
			dal.Join("cicd_pipelines cp on cp.id = ct.pipeline_id"),
			dal.Join("cicd_pipeline_commits cpc on cpc.pipeline_id = cp.id"),
			dal.Where("ct.type = ? and ct.environment = ? and commit_sha != ? and repo_id=? ", "DEPLOYMENT", "PRODUCTION", "", scopeId),
			dal.Orderby("ct.started_date"),
		)
		if err != nil {
			panic(err)
		}
		defer cursorDeployment.Close()

		commitPairs := make(DeploymentCommitPairs, 0)
		for cursorDeployment.Next() {
			var dcpr1 DeploymentCommitPairRev
			err := errors.Convert(cursorDeployment.Scan(&dcpr1))
			if err != nil {
				return err
			}
			for cursorDeployment.Next() {
				var dcpr2 DeploymentCommitPairRev
				err := errors.Convert(cursorDeployment.Scan(&dcpr2))
				if err != nil {
					return err
				}
				commitPairs = append(commitPairs, code.DeployCommitsDiff{ProjectName: projectName, ScopeId: scopeId, NewPipelineId: dcpr2.PipelineId, OldPipelineId: dcpr1.PipelineId, NewPipelineCommitSha: dcpr2.CommitSha, OldPipelineCommitSha: dcpr1.CommitSha, TaskId: dcpr1.TaskId, TaskName: dcpr1.TaskName})
			}

		}

		insertCountLimitOfDeployCommitsDiff := int(65535 / reflect.ValueOf(code.DeployCommitsDiff{}).NumField())
		commitNodeGraph := utils.NewCommitNodeGraph()

		// load commits from db
		commitParent := &code.CommitParent{}
		cursor, err := db.Cursor(
			dal.Select("cp.*"),
			dal.Join("LEFT JOIN repo_commits rc ON (rc.commit_sha = cp.commit_sha)"),
			dal.From("commit_parents cp"),
			dal.Where("rc.repo_id = ?", scopeId),
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
		commitsDiff := &code.DeployCommitsDiff{}
		lenCommitPairs := len(commitPairs)
		taskCtx.SetProgress(0, lenCommitPairs)

		for _, pair := range commitPairs {
			select {
			case <-ctx.Done():
				return errors.Convert(ctx.Err())
			default:
			}

			commitsDiff.ProjectName = pair.ProjectName
			commitsDiff.ScopeId = pair.ScopeId
			commitsDiff.NewPipelineId = pair.NewPipelineId
			commitsDiff.OldPipelineId = pair.OldPipelineId
			commitsDiff.NewPipelineCommitSha = pair.NewPipelineCommitSha
			commitsDiff.OldPipelineCommitSha = pair.OldPipelineCommitSha
			commitsDiff.TaskId = pair.TaskId
			commitsDiff.TaskName = pair.TaskName

			// delete records before creation
			err := db.Delete(&code.DeployCommitsDiff{ProjectName: commitsDiff.ProjectName, TaskId: commitsDiff.TaskId})
			if err != nil {
				return err
			}

			if commitsDiff.NewPipelineCommitSha == commitsDiff.OldPipelineCommitSha {
				// different deploy might point to a same commit, it is ok
				logger.Info(
					"skipping ref pair due to they are the same %s %s => %s",
					commitsDiff.NewPipelineId,
					commitsDiff.OldPipelineId,
					commitsDiff.NewPipelineCommitSha,
				)
				continue
			}

			lostSha, oldCount, newCount := commitNodeGraph.CalculateLostSha(commitsDiff.OldPipelineCommitSha, commitsDiff.NewPipelineCommitSha)

			commitsDiffs := []code.DeployCommitsDiff{}

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
					commitsDiffs = []code.DeployCommitsDiff{}
				}

				commitsDiff.SortingIndex++
			}

			if len(commitsDiffs) > 0 {
				logger.Info("insert data count [%d]", len(commitsDiffs))
				err = db.CreateIfNotExist(commitsDiffs)
				if err != nil {
					return err
				}
			}

			logger.Info(
				"total %d commits of difference found between [new][%s] and [old][%s(total:%d)]",
				newCount,
				commitsDiff.NewPipelineId,
				commitsDiff.OldPipelineId,
				oldCount,
			)
			taskCtx.IncProgress(1)
		}

	}

	return nil
}

var CalculateProjectDeploymentCommitsDiffMeta = core.SubTaskMeta{
	Name:             "calculateProjectDeploymentCommitsDiff",
	EntryPoint:       CalculateProjectDeploymentCommitsDiff,
	EnabledByDefault: true,
	Description:      "Calculate diff commits between project deployments",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE},
}
