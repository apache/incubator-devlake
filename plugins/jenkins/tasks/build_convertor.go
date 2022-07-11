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

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
)

var ConvertBuildsMeta = core.SubTaskMeta{
	Name:             "convertBuilds",
	EntryPoint:       ConvertBuilds,
	EnabledByDefault: true,
	Description:      "Convert tool layer table jenkins_builds into  domain layer table builds",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func ConvertBuilds(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*JenkinsTaskData)

	clauses := []dal.Clause{
		dal.Select("*"),
		dal.From("_tool_jenkins_builds"),
		dal.Where("connection_id = ?", data.Options.ConnectionId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	jobIdGen := didgen.NewDomainIdGenerator(&models.JenkinsJob{})
	buildIdGen := didgen.NewDomainIdGenerator(&models.JenkinsBuild{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.JenkinsBuild{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Params: JenkinsApiParams{
				ConnectionId: data.Options.ConnectionId,
			},
			Ctx:   taskCtx,
			Table: RAW_BUILD_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			jenkinsBuild := inputRow.(*models.JenkinsBuild)
			build := &devops.Build{
				DomainEntity: domainlayer.DomainEntity{
					Id: buildIdGen.Generate(jenkinsBuild.ConnectionId, jenkinsBuild.JobName, jenkinsBuild.Number),
				},
				JobId:       jobIdGen.Generate(jenkinsBuild.ConnectionId, jenkinsBuild.JobName),
				Name:        jenkinsBuild.DisplayName,
				DurationSec: uint64(jenkinsBuild.Duration / 1000),
				Status:      jenkinsBuild.Result,
				StartedDate: jenkinsBuild.StartTime,
				CommitSha:   jenkinsBuild.CommitSha,
			}
			return []interface{}{
				build,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
