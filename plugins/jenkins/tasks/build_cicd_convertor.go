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
	"regexp"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
)

var ConvertBuildsToCICDMeta = core.SubTaskMeta{
	Name:             "convertBuildsToCICD",
	EntryPoint:       ConvertBuildsToCICD,
	EnabledByDefault: true,
	Description:      "convert builds to cicd",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func ConvertBuildsToCICD(taskCtx core.SubTaskContext) (err errors.Error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*JenkinsTaskData)
	var deployTagRegexp *regexp.Regexp
	deploymentPattern := data.Options.DeploymentPattern
	if len(deploymentPattern) > 0 {
		deployTagRegexp, err = errors.Convert01(regexp.Compile(deploymentPattern))
		if err != nil {
			return errors.Default.Wrap(err, "regexp compile deploymentPattern failed")
		}
	}

	clauses := []dal.Clause{
		dal.From("_tool_jenkins_builds"),
		dal.Where("_tool_jenkins_builds.connection_id = ?", data.Options.ConnectionId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
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
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			jenkinsBuild := inputRow.(*models.JenkinsBuild)
			durationSec := int64(jenkinsBuild.Duration / 1000)
			jenkinsPipelineResult := ""
			jenkinsPipelineStatus := ""
			var jenkinsPipelineFinishedDate *time.Time
			results := make([]interface{}, 0)
			if jenkinsBuild.Result == "SUCCESS" {
				jenkinsPipelineResult = devops.SUCCESS
			} else if jenkinsBuild.Result == "FAILURE" {
				jenkinsPipelineResult = devops.FAILURE
			} else {
				jenkinsPipelineResult = devops.ABORT
			}

			if jenkinsBuild.Building {
				jenkinsPipelineStatus = devops.IN_PROGRESS
			} else {
				jenkinsPipelineStatus = devops.DONE
				finishTime := jenkinsBuild.StartTime.Add(time.Duration(durationSec * int64(time.Second)))
				jenkinsPipelineFinishedDate = &finishTime
			}
			jenkinsPipeline := &devops.CICDPipeline{
				DomainEntity: domainlayer.DomainEntity{
					Id: buildIdGen.Generate(jenkinsBuild.ConnectionId,
						jenkinsBuild.FullDisplayName),
				},
				Name:         jenkinsBuild.JobName,
				Result:       jenkinsPipelineResult,
				Status:       jenkinsPipelineStatus,
				FinishedDate: jenkinsPipelineFinishedDate,
				DurationSec:  uint64(durationSec),
				CreatedDate:  jenkinsBuild.StartTime,
			}

			if jenkinsBuild.TriggeredBy != "" {
				domainPipelineRelation := &devops.CICDPipelineRelationship{
					ParentPipelineId: buildIdGen.Generate(jenkinsBuild.ConnectionId,
						jenkinsBuild.TriggeredBy),
					ChildPipelineId: buildIdGen.Generate(jenkinsBuild.ConnectionId,
						jenkinsBuild.FullDisplayName),
				}
				results = append(results, domainPipelineRelation)
			}
			jenkinsPipeline.RawDataOrigin = jenkinsBuild.RawDataOrigin
			results = append(results, jenkinsPipeline)

			if !jenkinsBuild.HasStages {
				jenkinsTask := &devops.CICDTask{
					DomainEntity: domainlayer.DomainEntity{
						Id: buildIdGen.Generate(jenkinsBuild.ConnectionId,
							jenkinsBuild.FullDisplayName),
					},
					Name:         jenkinsBuild.JobName,
					Result:       jenkinsPipelineResult,
					Status:       jenkinsPipelineStatus,
					DurationSec:  uint64(durationSec),
					StartedDate:  jenkinsBuild.StartTime,
					FinishedDate: jenkinsPipelineFinishedDate,
				}
				if deployTagRegexp != nil {
					if deployFlag := deployTagRegexp.FindString(jenkinsBuild.JobName); deployFlag != "" {
						jenkinsTask.Type = devops.DEPLOYMENT
					}
				}

				jenkinsTask.PipelineId = buildIdGen.Generate(jenkinsBuild.ConnectionId, jenkinsBuild.FullDisplayName)
				jenkinsTask.RawDataOrigin = jenkinsBuild.RawDataOrigin
				results = append(results, jenkinsTask)

			}

			return results, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
