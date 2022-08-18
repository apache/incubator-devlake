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
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type JenkinsBuildWithRepo struct {
	// collected fields
	ConnectionId      uint64    `gorm:"primaryKey"`
	Duration          float64   // build time
	DisplayName       string    `gorm:"type:varchar(255)"` // "#7"
	EstimatedDuration float64   // EstimatedDuration
	Number            int64     `gorm:"primaryKey"`
	Result            string    // Result
	Timestamp         int64     // start time
	StartTime         time.Time // convered by timestamp
	CommitSha         string    `gorm:"type:varchar(255)"`
	Type              string    `gorm:"index;type:varchar(255)"`
	Class             string    `gorm:"index;type:varchar(255)" `
	TriggeredBy       string    `gorm:"type:varchar(255)"`
	Building          bool
	Branch            string `gorm:"type:varchar(255)"`
	RepoUrl           string `gorm:"type:varchar(255)"`
	HasStages         bool
	common.NoPKModel
}

var ConvertBuildsToCICDMeta = core.SubTaskMeta{
	Name:             "convertBuildsToCICD",
	EntryPoint:       ConvertBuildsToCICD,
	EnabledByDefault: true,
	Description:      "convert builds to cicd",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func ConvertBuildsToCICD(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*JenkinsTaskData)

	clauses := []dal.Clause{
		dal.Select(`tjb.connection_id, tjb.duration, tjb.display_name, tjb.estimated_duration, tjb.number,
			tjb._raw_data_remark, tjb._raw_data_id, tjb._raw_data_table, tjb._raw_data_params,
			tjb.result, tjb.timestamp, tjb.start_time, tjbr.commit_sha, tjb.type, tjb.class, 
			tjb.triggered_by, tjb.building, tjbr.branch, tjbr.repo_url, tjb.has_stages`),
		dal.From("_tool_jenkins_builds tjb"),
		dal.Join("left join _tool_jenkins_build_repos tjbr on tjbr.build_name = tjb.display_name"),
		dal.Where("tjb.connection_id = ?", data.Options.ConnectionId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(JenkinsBuildWithRepo{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Params: JenkinsApiParams{
				ConnectionId: data.Options.ConnectionId,
			},
			Ctx:   taskCtx,
			Table: RAW_BUILD_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			jenkinsBuildWithRepo := inputRow.(*JenkinsBuildWithRepo)
			durationSec := int64(jenkinsBuildWithRepo.Duration / 1000)
			jenkinsPipelineResult := ""
			jenkinsPipelineStatus := ""
			var jenkinsPipelineFinishedDate *time.Time
			results := make([]interface{}, 0)
			if jenkinsBuildWithRepo.Result == "SUCCESS" {
				jenkinsPipelineResult = devops.SUCCESS
			} else if jenkinsBuildWithRepo.Result == "FAILURE" {
				jenkinsPipelineResult = devops.FAILURE
			} else {
				jenkinsPipelineResult = devops.ABORT
			}

			if jenkinsBuildWithRepo.Building {
				jenkinsPipelineStatus = devops.IN_PROGRESS
			} else {
				jenkinsPipelineStatus = devops.DONE
				finishTime := jenkinsBuildWithRepo.StartTime.Add(time.Duration(durationSec * int64(time.Second)))
				jenkinsPipelineFinishedDate = &finishTime
			}
			if jenkinsBuildWithRepo.TriggeredBy == "" {
				jenkinsPipeline := &devops.CICDPipeline{
					DomainEntity: domainlayer.DomainEntity{
						Id: fmt.Sprintf("%s:%s:%d:%s:%s", "jenkins", "JenkinsPipeline", jenkinsBuildWithRepo.ConnectionId,
							jenkinsBuildWithRepo.CommitSha, jenkinsBuildWithRepo.DisplayName),
					},
					Name:         jenkinsBuildWithRepo.DisplayName,
					CommitSha:    jenkinsBuildWithRepo.CommitSha,
					Branch:       jenkinsBuildWithRepo.Branch,
					Repo:         jenkinsBuildWithRepo.RepoUrl,
					Result:       jenkinsPipelineResult,
					Status:       jenkinsPipelineStatus,
					FinishedDate: jenkinsPipelineFinishedDate,
					Type:         "CI/CD",
					DurationSec:  uint64(durationSec),
					CreatedDate:  jenkinsBuildWithRepo.StartTime,
				}
				if jenkinsPipelineFinishedDate != nil {
				}
				jenkinsPipeline.RawDataOrigin = jenkinsBuildWithRepo.RawDataOrigin
				results = append(results, jenkinsPipeline)
			}

			if !jenkinsBuildWithRepo.HasStages {
				jenkinsTask := &devops.CICDTask{
					DomainEntity: domainlayer.DomainEntity{
						Id: fmt.Sprintf("%s:%s:%d:%s", "jenkins", "JenkinsTask", jenkinsBuildWithRepo.ConnectionId,
							jenkinsBuildWithRepo.DisplayName),
					},
					Name:         jenkinsBuildWithRepo.DisplayName,
					Result:       jenkinsPipelineResult,
					Status:       jenkinsPipelineStatus,
					Type:         "CI/CD",
					DurationSec:  uint64(durationSec),
					StartedDate:  jenkinsBuildWithRepo.StartTime,
					FinishedDate: jenkinsPipelineFinishedDate,
				}
				if jenkinsBuildWithRepo.TriggeredBy != "" {
					tmp := make([]*JenkinsBuildWithRepo, 0)
					clauses := []dal.Clause{
						dal.Select(`tjb.display_name, tjb.result, tjb.timestamp, tjbr.commit_sha`),
						dal.From("_tool_jenkins_builds tjb"),
						dal.Join("left join _tool_jenkins_build_repos tjbr on tjbr.build_name = tjb.display_name"),
						dal.Where("tjb.connection_id = ? and tjb.display_name = ?", data.Options.ConnectionId, jenkinsBuildWithRepo.TriggeredBy),
					}
					err = db.All(&tmp, clauses...)
					if err != nil {
						return nil, err
					}
					if len(tmp) > 0 {
						jenkinsTask.PipelineId = fmt.Sprintf("%s:%s:%d:%s:%s", "jenkins", "JenkinsPipeline", jenkinsBuildWithRepo.ConnectionId,
							tmp[0].CommitSha, tmp[0].DisplayName)
					}
				} else {
					jenkinsTask.PipelineId = fmt.Sprintf("%s:%s:%d:%s:%s", "jenkins", "JenkinsPipeline", jenkinsBuildWithRepo.ConnectionId,
						jenkinsBuildWithRepo.CommitSha, jenkinsBuildWithRepo.DisplayName)
				}
				jenkinsTask.RawDataOrigin = jenkinsBuildWithRepo.RawDataOrigin
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
