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
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

var _ plugin.SubTaskEntryPoint = DBGetTaskRepoCommits

var DBGetTaskRepoCommitsMeta = plugin.SubTaskMeta{
	Name:             "collectTaskRepoCommits",
	EntryPoint:       DBGetTaskRepoCommits,
	EnabledByDefault: true,
	Description:      "Get task commits data from Zentao database",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func DBGetTaskRepoCommits(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)

	// skip if no RemoteDb
	if data.RemoteDb == nil {
		return nil
	}

	divider := api.NewBatchSaveDivider(taskCtx, 500, "", "")
	defer func() {
		err1 := divider.Close()
		if err1 != nil {
			panic(err1)
		}
	}()
	handler, err := newTaskRepoCommitHandler(taskCtx, divider)
	if err != nil {
		return err
	}
	return handler.collectTaskRepoCommit(
		taskCtx.GetDal(),
		data.RemoteDb,
		data.Options.ProjectId,
		data.Options.ConnectionId,
	)
}

type taskRepoCommitHandler struct {
	rawDataParams          string
	taskRepoCommitBachSave *api.BatchSave
}

type toolZentaoTask struct {
	TaskID int64 `gorm:"type:integer"`
}

type RemoteTaskRepoCommit struct {
	Project   int64 `gorm:"type:integer"`
	IssueID   int64 `gorm:"type:integer"`
	RepoUrl   string
	CommitSha string
}

func (h taskRepoCommitHandler) collectTaskRepoCommit(db dal.Dal, rdb dal.Dal, projectId int64, connectionId uint64) errors.Error {
	taskCursor, err := db.RawCursor(`
		SELECT
			DISTINCT id AS task_id
		FROM
			_tool_zentao_tasks AS fact_task
		WHERE
			fact_task.project = ? AND
			fact_task.connection_id = ?
	`, projectId, connectionId)
	if err != nil {
		return err
	}
	defer taskCursor.Close()
	var taskIds []int64
	for taskCursor.Next() {
		var row toolZentaoTask
		err := db.Fetch(taskCursor, &row)
		if err != nil {
			return errors.Default.Wrap(err, "error fetching taskCursor")
		}
		taskIds = append(taskIds, row.TaskID)
	}

	remoteCursor, err := rdb.RawCursor(`
		SELECT
			DISTINCT dim_task_commits.task_id AS issue_id,
			dim_repo.path AS repo_url,
			dim_revision.revision AS commit_sha
		FROM (
			SELECT
				fact_action.objectID AS task_id,
				fact_action.project AS project_id,
				fact_action.product AS product_id,
				fact_action.extra AS short_commit_hexsha
			FROM
				zt_action AS fact_action
			WHERE
				fact_action.objectType IN ('task') AND
				fact_action.objectID IN ? AND
				fact_action.action IN ('gitcommited')
		) AS dim_task_commits
		INNER JOIN (
			SELECT
				fact_repo_hist.repo AS repo_id,
				fact_repo_hist.revision AS revision,
				LEFT(fact_repo_hist.revision, 10) AS short_commit_hexsha
			FROM
				zt_repohistory AS fact_repo_hist
		) AS dim_revision
		ON
			dim_task_commits.short_commit_hexsha = dim_revision.short_commit_hexsha
		INNER JOIN
			zt_repo AS dim_repo
		ON
			dim_revision.repo_id = dim_repo.id
	`, taskIds)
	if err != nil {
		return err
	}
	defer remoteCursor.Close()

	for remoteCursor.Next() {
		var remoteTaskRepoCommit RemoteTaskRepoCommit
		err = rdb.Fetch(remoteCursor, &remoteTaskRepoCommit)
		if err != nil {
			return err
		}
		taskRepoCommit := &models.ZentaoTaskRepoCommit{
			ConnectionId: connectionId,
			Project:      projectId,
			RepoUrl:      remoteTaskRepoCommit.RepoUrl,
			CommitSha:    remoteTaskRepoCommit.CommitSha,
			IssueId:      strconv.FormatInt(remoteTaskRepoCommit.IssueID, 10),
		}
		taskRepoCommit.NoPKModel.RawDataParams = h.rawDataParams
		err = h.taskRepoCommitBachSave.Add(taskRepoCommit)
		if err != nil {
			return err
		}
	}
	return h.taskRepoCommitBachSave.Flush()
}

func newTaskRepoCommitHandler(taskCtx plugin.SubTaskContext, divider *api.BatchSaveDivider) (*taskRepoCommitHandler, errors.Error) {
	data := taskCtx.GetData().(*ZentaoTaskData)

	taskRepoCommitBachSave, err := divider.ForType(reflect.TypeOf(&models.ZentaoTaskRepoCommit{}))
	if err != nil {
		return nil, err
	}
	blob, _ := json.Marshal(data.Options.GetParams())
	rawDataParams := string(blob)
	db := taskCtx.GetDal()
	err = db.Delete(&models.ZentaoTaskRepoCommit{}, dal.Where("_raw_data_params = ?", rawDataParams))
	if err != nil {
		return nil, err
	}
	return &taskRepoCommitHandler{
		rawDataParams:          rawDataParams,
		taskRepoCommitBachSave: taskRepoCommitBachSave,
	}, nil
}
