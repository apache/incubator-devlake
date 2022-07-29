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
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
)

// this struct should be moved to `gitub_api_common.go`

var ExtractApiBuildsMeta = core.SubTaskMeta{
	Name:             "extractApiBuilds",
	EntryPoint:       ExtractApiBuilds,
	EnabledByDefault: true,
	Description:      "Extract raw builds data into tool layer table jenkins_builds",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func ExtractApiBuilds(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JenkinsTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Params: JenkinsApiParams{
				ConnectionId: data.Options.ConnectionId,
			},
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			/*
				Table store raw data
			*/
			Table: RAW_BUILD_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			body := &models.ApiBuildResponse{}
			err := json.Unmarshal(row.Data, body)
			if err != nil {
				return nil, err
			}
			input := &SimpleJob{}
			err = json.Unmarshal(row.Input, input)
			if err != nil {
				return nil, err
			}

			results := make([]interface{}, 0)
			strList := strings.Split(body.Class, ".")
			class := strList[len(strList)-1]
			build := &models.JenkinsBuild{
				ConnectionId:      data.Options.ConnectionId,
				Type:              class,
				JobName:           input.Name,
				Duration:          body.Duration,
				DisplayName:       body.DisplayName,
				EstimatedDuration: body.EstimatedDuration,
				Number:            body.Number,
				Result:            body.Result,
				Timestamp:         body.Timestamp,
				StartTime:         time.Unix(body.Timestamp/1000, 0),
			}
			vcs := body.ChangeSet.Kind
			if vcs == "git" || vcs == "hg" {
				for _, a := range body.Actions {
					sha := ""
					if a.LastBuiltRevision.SHA1 != "" {
						sha = a.LastBuiltRevision.SHA1
					}
					if a.MercurialRevisionNumber != "" {
						sha = a.MercurialRevisionNumber
					}
					build.CommitSha = sha
					for _, url := range a.RemoteUrls {
						if url != "" {
							buildCommitRemoteUrl := models.JenkinsBuildCommitRepoUrl{
								ConnectionId: data.Options.ConnectionId,
								BuildName:    build.DisplayName,
								CommitSha:    sha,
								RemoteUrl:    url,
							}
							results = append(results, &buildCommitRemoteUrl)
						}
					}
					if a.TriggeredBuilds != nil && len(a.TriggeredBuilds) > 0 {
						for _, b := range a.TriggeredBuilds {
							if b.DisplayName == "" {
								continue
							}
							buildTrigger := models.JenkinsBuildTriggeredBuilds{
								ConnectionId:       data.Options.ConnectionId,
								BuildName:          build.DisplayName,
								TriggeredBuildName: b.DisplayName,
							}
							results = append(results, &buildTrigger)
						}
					}
				}
			} else if vcs == "svn" {
				if len(body.ChangeSet.Revisions) > 0 {
					build.CommitSha = strconv.Itoa(body.ChangeSet.Revisions[0].Revision)
				}
			}

			results = append(results, build)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
