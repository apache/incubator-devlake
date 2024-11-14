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

var _ plugin.SubTaskEntryPoint = DBGetStoryRepoCommits

var DBGetStoryRepoCommitsMeta = plugin.SubTaskMeta{
	Name:             "collectStoryRepoCommits",
	EntryPoint:       DBGetStoryRepoCommits,
	EnabledByDefault: true,
	Description:      "Get story commits data from Zentao database",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func DBGetStoryRepoCommits(taskCtx plugin.SubTaskContext) errors.Error {
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
	handler, err := newStoryRepoCommitHandler(taskCtx, divider)
	if err != nil {
		return err
	}
	return handler.collectStoryRepoCommit(
		taskCtx.GetDal(),
		data.RemoteDb,
		data.Options.ProjectId,
		data.Options.ConnectionId,
	)
}

type storyRepoCommitHandler struct {
	rawDataParams           string
	storyRepoCommitBachSave *api.BatchSave
}

type toolZentaoStory struct {
	StoryID int64 `gorm:"type:integer"`
}

type RemoteStoryRepoCommit struct {
	Project   int64 `gorm:"type:integer"`
	IssueID   int64 `gorm:"type:integer"`
	RepoUrl   string
	CommitSha string
}

func (h storyRepoCommitHandler) collectStoryRepoCommit(db dal.Dal, rdb dal.Dal, projectId int64, connectionId uint64) errors.Error {
	storyCursor, err := db.RawCursor(`
		SELECT
			DISTINCT fact_project_story.story_id AS story_id
		FROM
			_tool_zentao_project_stories AS fact_project_story
		LEFT JOIN
			_tool_zentao_stories AS fact_story
		ON
			fact_project_story.story_id = fact_story.id AND
			fact_project_story.connection_id = fact_story.connection_id
		WHERE
			fact_project_story.project_id = ? AND
			fact_project_story.connection_id = ?
	`, projectId, connectionId)
	if err != nil {
		return err
	}
	defer storyCursor.Close()
	var storyIds []int64
	for storyCursor.Next() {
		var row toolZentaoStory
		err := db.Fetch(storyCursor, &row)
		if err != nil {
			return errors.Default.Wrap(err, "error fetching storyCursor")
		}
		storyIds = append(storyIds, row.StoryID)
	}

	remoteCursor, err := rdb.RawCursor(`
		SELECT
			DISTINCT dim_story_commits.story_id AS issue_id,
			dim_repo.path AS repo_url,
			dim_revision.revision AS commit_sha
		FROM (
			SELECT
				fact_action.objectID AS story_id,
				fact_action.project AS project_id,
				fact_action.product AS product_id,
				fact_action.extra AS short_commit_hexsha
			FROM
				zt_action AS fact_action
			WHERE
				fact_action.objectType IN ('story', 'requirement') AND
				fact_action.objectID IN ? AND
				fact_action.action IN ('gitcommited')
		) AS dim_story_commits
		INNER JOIN (
			SELECT
				fact_repo_hist.repo AS repo_id,
				fact_repo_hist.revision AS revision,
				LEFT(fact_repo_hist.revision, 10) AS short_commit_hexsha
			FROM
				zt_repohistory AS fact_repo_hist
		) AS dim_revision
		ON
			dim_story_commits.short_commit_hexsha = dim_revision.short_commit_hexsha
		INNER JOIN
			zt_repo AS dim_repo
		ON
			dim_revision.repo_id = dim_repo.id
	`, storyIds)
	if err != nil {
		return err
	}
	defer remoteCursor.Close()

	for remoteCursor.Next() {
		var remoteStoryRepoCommit RemoteStoryRepoCommit
		err = rdb.Fetch(remoteCursor, &remoteStoryRepoCommit)
		if err != nil {
			return err
		}
		storyRepoCommit := &models.ZentaoStoryRepoCommit{
			ConnectionId: connectionId,
			Project:      projectId,
			RepoUrl:      remoteStoryRepoCommit.RepoUrl,
			CommitSha:    remoteStoryRepoCommit.CommitSha,
			IssueId:      strconv.FormatInt(remoteStoryRepoCommit.IssueID, 10),
		}
		storyRepoCommit.NoPKModel.RawDataParams = h.rawDataParams
		err = h.storyRepoCommitBachSave.Add(storyRepoCommit)
		if err != nil {
			return err
		}
	}
	return h.storyRepoCommitBachSave.Flush()
}

func newStoryRepoCommitHandler(taskCtx plugin.SubTaskContext, divider *api.BatchSaveDivider) (*storyRepoCommitHandler, errors.Error) {
	data := taskCtx.GetData().(*ZentaoTaskData)

	storyRepoCommitBachSave, err := divider.ForType(reflect.TypeOf(&models.ZentaoStoryRepoCommit{}))
	if err != nil {
		return nil, err
	}
	blob, _ := json.Marshal(data.Options.GetParams())
	rawDataParams := string(blob)
	db := taskCtx.GetDal()
	err = db.Delete(&models.ZentaoStoryRepoCommit{}, dal.Where("_raw_data_params = ?", rawDataParams))
	if err != nil {
		return nil, err
	}
	return &storyRepoCommitHandler{
		rawDataParams:           rawDataParams,
		storyRepoCommitBachSave: storyRepoCommitBachSave,
	}, nil
}
