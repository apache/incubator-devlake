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
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
)

// this struct should be moved to `gitub_api_common.go`

var EnrichApiBuildsMeta = core.SubTaskMeta{
	Name:             "enrichApiBuilds",
	EntryPoint:       EnrichApiBuilds,
	EnabledByDefault: true,
	Description:      "Enrich  jenkins_builds",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func EnrichApiBuilds(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JenkinsTaskData)
	db := taskCtx.GetDal()
	clauses := []dal.Clause{
		dal.Select("distinct build_name"),
		dal.From(&models.JenkinsStage{}),
		dal.Where("connection_id = ?", data.Options.ConnectionId),
		dal.Groupby("build_name"),
	}
	cursor, err := db.Cursor(clauses...)
	for cursor.Next() {
		var buildName string
		err = cursor.Scan(&buildName)
		if err != nil {
			return err
		}
		if buildName == "" {
			continue
		}
		build := &models.JenkinsBuild{}
		build.HasStages = true
		err = db.Update(&models.JenkinsBuild{}, dal.Where("connection_id = ? and build_name = ?", data.Options.ConnectionId, buildName))
		if err != nil {
			return err
		}
	}
	return nil
}
