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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"reflect"

	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
)

var ConvertBuildReposMeta = core.SubTaskMeta{
	Name:             "convertBuildRepos",
	EntryPoint:       ConvertBuildRepos,
	EnabledByDefault: true,
	Description:      "Convert tool layer table jenkins_builds into  domain layer table builds",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func ConvertBuildRepos(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*JenkinsTaskData)

	clauses := []dal.Clause{
		dal.Select("*"),
		dal.From(&models.JenkinsBuildCommit{}),
		dal.Join(`left join _tool_jenkins_builds tjb 
						on _tool_jenkins_build_commits.build_name = tjb.full_name 
						and _tool_jenkins_build_commits.connection_id = tjb.connection_id`),
		dal.Where(`_tool_jenkins_build_commits.connection_id = ?
							and tjb.job_path = ? and tjb.job_name = ?`,
			data.Options.ConnectionId, data.Options.JobPath, data.Options.JobName),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	buildIdGen := didgen.NewDomainIdGenerator(&models.JenkinsBuild{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.JenkinsBuildCommit{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Params: JenkinsApiParams{
				ConnectionId: data.Options.ConnectionId,
				FullName:     data.Options.JobFullName,
			},
			Ctx:   taskCtx,
			Table: RAW_BUILD_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			jenkinsBuildCommit := inputRow.(*models.JenkinsBuildCommit)
			build := &devops.CiCDPipelineCommit{
				PipelineId: buildIdGen.Generate(jenkinsBuildCommit.ConnectionId, jenkinsBuildCommit.BuildName),
				CommitSha:  jenkinsBuildCommit.CommitSha,
				Branch:     jenkinsBuildCommit.Branch,
				RepoUrl:    jenkinsBuildCommit.RepoUrl,
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
