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
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
	"reflect"
)

// this struct should be moved to `gitub_api_common.go`

var EnrichApiBuildOriginTriggerMeta = core.SubTaskMeta{
	Name:             "enrichApiBuildOriginTrigger",
	EntryPoint:       EnrichApiBuildOriginTrigger,
	EnabledByDefault: true,
	Description:      "Enrich jenkins build with origin trigger",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func EnrichApiBuildOriginTrigger(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JenkinsTaskData)
	db := taskCtx.GetDal()
	clauses := []dal.Clause{
		dal.From(&models.JenkinsBuild{}),
		dal.Where("connection_id = ?", data.Options.ConnectionId),
	}
	cursor, err := db.Cursor(clauses...)
	taskCtx.SetProgress(0, -1)
	iterator, err := helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(models.JenkinsBuild{}))
	if err != nil {
		return err
	}
	defer iterator.Close()
	for iterator.HasNext() {
		i, err := iterator.Fetch()
		build := i.(*models.JenkinsBuild)

		if build.TriggeredBy == "" {
			continue
		}
		buildTmp := &models.JenkinsBuild{}
		buildTmp.DisplayName = build.TriggeredBy
		for {
			err = db.First(buildTmp)
			if err != nil {
				return err
			}
			if buildTmp.TriggeredBy == "" {
				break
			}
			buildTmp.DisplayName = buildTmp.TriggeredBy
		}
		build.TriggeredBy = buildTmp.DisplayName
		err = db.Update(build)
		if err != nil {
			return err
		}
		taskCtx.IncProgress(1)
	}
	return nil
}
