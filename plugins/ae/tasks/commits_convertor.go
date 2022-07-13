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
	aeModels "github.com/apache/incubator-devlake/plugins/ae/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
)

// NOTE: This only works on Commits in the Domain layer. You need to run Github or Gitlab collection and Domain layer enrichemnt first.
func ConvertCommits(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*AeTaskData)

	aeCommit := &aeModels.AECommit{}

	// Get all the commits from the domain layer
	cursor, err := db.Cursor(
		dal.From(aeCommit),
		dal.Where("ae_project_id = ?", data.Options.ProjectId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	taskCtx.SetProgress(0, -1)
	ctx := taskCtx.GetContext()
	// Loop over them
	for cursor.Next() {
		// we do a non-blocking checking, this should be applied to all converters from all plugins
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		// uncomment following line if you want to test out canceling feature for this task
		//time.Sleep(1 * time.Second)

		err = db.Fetch(cursor, aeCommit)
		if err != nil {
			return err
		}

		err = db.Exec("UPDATE commits SET dev_eq = ? WHERE sha = ? ", aeCommit.DevEq, aeCommit.HexSha)
		if err != nil {
			return err
		}

		taskCtx.IncProgress(1)
	}

	return nil
}

var ConvertCommitsMeta = core.SubTaskMeta{
	Name:             "convertCommits",
	EntryPoint:       ConvertCommits,
	EnabledByDefault: true,
	Description:      "Update domain layer commits dev_eq field according to ae_commits",
}
