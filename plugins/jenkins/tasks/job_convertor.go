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
	. "github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
)

var ConvertJobsMeta = core.SubTaskMeta{
	Name:             "convertJobs",
	EntryPoint:       ConvertJobs,
	EnabledByDefault: true,
	Description:      "Convert tool layer table jenkins_jobs into  domain layer table jobs",
}

func ConvertJobs(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*JenkinsTaskData)

	clauses := []Clause{
		Select("*"),
		From("_tool_jenkins_jobs"),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}

	jobIdGen := didgen.NewDomainIdGenerator(&models.JenkinsJob{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.JenkinsJob{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Params: JenkinsApiParams{
				ConnectionId: data.Connection.ID,
			},
			Ctx:   taskCtx,
			Table: RAW_JOB_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			jenkinsJob := inputRow.(*models.JenkinsJob)
			job := &devops.Job{
				DomainEntity: domainlayer.DomainEntity{
					Id: jobIdGen.Generate(jenkinsJob.ConnectionId, jenkinsJob.Name),
				},
				Name: jenkinsJob.Name,
			}
			return []interface{}{
				job,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
