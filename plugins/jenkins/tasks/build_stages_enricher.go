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
	goerror "errors"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
	"gorm.io/gorm"
)

var EnrichApiBuildWithStagesMeta = core.SubTaskMeta{
	Name:             "enrichApiBuildWithStages",
	EntryPoint:       EnrichApiBuildWithStages,
	EnabledByDefault: true,
	Description:      "Enrich  jenkins build with stages",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func EnrichApiBuildWithStages(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JenkinsTaskData)
	db := taskCtx.GetDal()
	clauses := []dal.Clause{
		dal.Select("distinct build_name"),
		dal.From(&models.JenkinsStage{}),
		dal.Where("connection_id = ?", data.Options.ConnectionId),
		dal.Groupby("build_name"),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	taskCtx.SetProgress(0, -1)

	for cursor.Next() {
		var buildName string
		err = errors.Convert(cursor.Scan(&buildName))
		if err != nil {
			return err
		}
		if buildName == "" {
			continue
		}
		build := &models.JenkinsBuild{}
		build.ConnectionId = data.Options.ConnectionId
		build.FullDisplayName = buildName
		err = db.First(build)
		if goerror.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		if err != nil {
			return errors.Convert(err)
		}
		build.HasStages = true

		err = db.Update(build)
		if err != nil {
			return errors.Convert(err)
		}
		taskCtx.IncProgress(1)
	}
	return nil
}
