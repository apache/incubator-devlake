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

var _ plugin.SubTaskEntryPoint = DBGetBugRepoCommits

var DBGetBugRepoCommitsMeta = plugin.SubTaskMeta{
	Name:             "collectBugRepoCommits",
	EntryPoint:       DBGetBugRepoCommits,
	EnabledByDefault: true,
	Description:      "Get bug commits data from Zentao database",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func DBGetBugRepoCommits(taskCtx plugin.SubTaskContext) errors.Error {
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
	handler, err := newBugRepoCommitHandler(taskCtx, divider)
	if err != nil {
		return err
	}
	return handler.collectBugRepoCommit(
		taskCtx.GetDal(),
		data.RemoteDb,
		data.Options.ProjectId,
		data.Options.ConnectionId,
	)
}

type bugRepoCommitHandler struct {
	rawDataParams         string
	bugRepoCommitBachSave *api.BatchSave
}

type toolZentaoBug struct {
	BugID int64 `gorm:"type:integer"`
}

type RemoteBugRepoCommit struct {
	Project   int64 `gorm:"type:integer"`
	IssueID   int64 `gorm:"type:integer"`
	RepoUrl   string
	CommitSha string
}

func (h bugRepoCommitHandler) collectBugRepoCommit(db dal.Dal, rdb dal.Dal, projectId int64, connectionId uint64) errors.Error {
	bugCursor, err := db.RawCursor(`
		SELECT
			DISTINCT id AS bug_id
		FROM
			_tool_zentao_bugs AS fact_bug
		WHERE
			fact_bug.project = ? AND
			fact_bug.connection_id = ?
	`, projectId, connectionId)
	if err != nil {
		return err
	}
	defer bugCursor.Close()
	var bugIds []int64
	for bugCursor.Next() {
		var row toolZentaoBug
		err := db.Fetch(bugCursor, &row)
		if err != nil {
			return errors.Default.Wrap(err, "error fetching bugCursor")
		}
		bugIds = append(bugIds, row.BugID)
	}

	remoteCursor, err := rdb.RawCursor(`
		SELECT
			DISTINCT dim_bug_commits.bug_id AS issue_id,
			dim_repo.path AS repo_url,
			dim_revision.revision AS commit_sha
		FROM (
			SELECT
				fact_action.objectID AS bug_id,
				fact_action.project AS project_id,
				fact_action.product AS product_id,
				fact_action.extra AS short_commit_hexsha
			FROM
				zt_action AS fact_action
			WHERE
				fact_action.objectType IN ('bug') AND
				fact_action.objectID IN ? AND
				fact_action.action IN ('gitcommited')
		) AS dim_bug_commits
		INNER JOIN (
			SELECT
				fact_repo_hist.repo AS repo_id,
				fact_repo_hist.revision AS revision,
				LEFT(fact_repo_hist.revision, 10) AS short_commit_hexsha
			FROM
				zt_repohistory AS fact_repo_hist
		) AS dim_revision
		ON
			dim_bug_commits.short_commit_hexsha = dim_revision.short_commit_hexsha
		INNER JOIN
			zt_repo AS dim_repo
		ON
			dim_revision.repo_id = dim_repo.id
	`, bugIds)
	if err != nil {
		return err
	}
	defer remoteCursor.Close()

	for remoteCursor.Next() {
		var remoteBugRepoCommit RemoteBugRepoCommit
		err = rdb.Fetch(remoteCursor, &remoteBugRepoCommit)
		if err != nil {
			return err
		}
		bugRepoCommit := &models.ZentaoBugRepoCommit{
			ConnectionId: connectionId,
			Project:      projectId,
			RepoUrl:      remoteBugRepoCommit.RepoUrl,
			CommitSha:    remoteBugRepoCommit.CommitSha,
			IssueId:      strconv.FormatInt(remoteBugRepoCommit.IssueID, 10),
		}
		bugRepoCommit.NoPKModel.RawDataParams = h.rawDataParams
		err = h.bugRepoCommitBachSave.Add(bugRepoCommit)
		if err != nil {
			return err
		}
	}
	return h.bugRepoCommitBachSave.Flush()
}

func newBugRepoCommitHandler(taskCtx plugin.SubTaskContext, divider *api.BatchSaveDivider) (*bugRepoCommitHandler, errors.Error) {
	data := taskCtx.GetData().(*ZentaoTaskData)

	bugRepoCommitBachSave, err := divider.ForType(reflect.TypeOf(&models.ZentaoBugRepoCommit{}))
	if err != nil {
		return nil, err
	}
	blob, _ := json.Marshal(data.Options.GetParams())
	rawDataParams := string(blob)
	db := taskCtx.GetDal()
	err = db.Delete(&models.ZentaoBugRepoCommit{}, dal.Where("_raw_data_params = ?", rawDataParams))
	if err != nil {
		return nil, err
	}
	return &bugRepoCommitHandler{
		rawDataParams:         rawDataParams,
		bugRepoCommitBachSave: bugRepoCommitBachSave,
	}, nil
}
