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
	"github.com/apache/incubator-devlake/errors"
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
		dal.From(&models.JenkinsBuildRepo{}),
		dal.Where("connection_id = ?", data.Options.ConnectionId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.JenkinsBuildRepo{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Params: JenkinsApiParams{
				ConnectionId: data.Options.ConnectionId,
			},
			Ctx:   taskCtx,
			Table: RAW_BUILD_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			jenkinsBuildRepo := inputRow.(*models.JenkinsBuildRepo)
			build := &devops.CiCDPipelineCommit{
				PipelineId: fmt.Sprintf("%s:%s:%d:%s", "jenkins", "JenkinsTask", jenkinsBuildRepo.ConnectionId,
					jenkinsBuildRepo.BuildName),
				CommitSha: jenkinsBuildRepo.CommitSha,
				Branch:    jenkinsBuildRepo.Branch,
				Repo:      jenkinsBuildRepo.RepoUrl,
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
