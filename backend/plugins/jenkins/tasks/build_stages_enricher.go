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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
)

var EnrichApiBuildWithStagesMeta = plugin.SubTaskMeta{
	Name:             "enrichApiBuildWithStages",
	EntryPoint:       EnrichApiBuildWithStages,
	EnabledByDefault: true,
	Description:      "Enrich  jenkins build with stages",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func EnrichApiBuildWithStages(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JenkinsTaskData)
	db := taskCtx.GetDal()

	clauses := []dal.Clause{
		dal.Select("tjb.*"),
		dal.From(`_tool_jenkins_builds tjb`),
		dal.Join(`inner join _tool_jenkins_stages tjs 
				on tjs.build_name = tjb.full_name 
				and tjs.connection_id = tjb.connection_id`),
	}

	if data.Options.Class == WORKFLOW_MULTI_BRANCH_PROJECT {
		clauses = append(clauses,
			dal.Where(`tjb.connection_id = ? and tjb.full_name like ?`,
				data.Options.ConnectionId, fmt.Sprintf("%s%%", data.Options.JobFullName)))
	} else {
		clauses = append(clauses,
			dal.Where(`tjb.connection_id = ? and tjb.job_path = ? and tjb.job_name = ?`,
				data.Options.ConnectionId, data.Options.JobPath, data.Options.JobName))
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	taskCtx.SetProgress(0, -1)

	for cursor.Next() {
		build := &models.JenkinsBuild{}
		err = db.Fetch(cursor, build)
		if err != nil {
			return err
		}
		build.HasStages = true
		err = db.CreateOrUpdate(build)
		if err != nil {
			return err
		}
	}

	return nil
}
