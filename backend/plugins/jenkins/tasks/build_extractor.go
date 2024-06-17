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
	"fmt"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
)

var ExtractApiBuildsMeta = plugin.SubTaskMeta{
	Name:             "extractApiBuilds",
	EntryPoint:       ExtractApiBuilds,
	EnabledByDefault: true,
	Description:      "Extract raw builds data into tool layer table jenkins_builds",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func ExtractApiBuilds(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JenkinsTaskData)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Params: JenkinsApiParams{
				ConnectionId: data.Options.ConnectionId,
				FullName:     data.Options.JobFullName,
			},
			Ctx:   taskCtx,
			Table: RAW_BUILD_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			body := &models.ApiBuildResponse{}
			err := errors.Convert(json.Unmarshal(row.Data, body))
			if err != nil {
				return nil, err
			}

			jobName := data.Options.JobName
			jobPath := data.Options.JobPath
			fullName := fmt.Sprintf(`%s#%d`, data.Options.JobFullName, body.Number)

			input := &SimpleJob{}
			err1 := json.Unmarshal(row.Input, input)
			if err1 == nil && input.Class == WORKFLOW_JOB {
				// For jobs from multi-branch workflow, the job name and path must come from the input,
				// otherwise it will be set to the multi-branch workflow name and path
				jobName = input.Name
				jobPath = input.Path
				fullName = fmt.Sprintf(`%s#%d`, input.FullName, body.Number)
			}

			results := make([]interface{}, 0)
			strList := strings.Split(body.Class, ".")
			class := strList[len(strList)-1]
			build := &models.JenkinsBuild{
				ConnectionId:      data.Options.ConnectionId,
				JobName:           jobName,
				JobPath:           jobPath,
				Duration:          body.Duration,
				FullName:          fullName,
				EstimatedDuration: body.EstimatedDuration,
				Number:            body.Number,
				Result:            body.Result,
				Timestamp:         body.Timestamp,
				Class:             class,
				Building:          body.Building,
				StartTime:         time.Unix(body.Timestamp/1000, 0),
			}
			// we also need to collect the commit info from the build which does not have changeSet
			// changeSet describes the changes that were made in the build
			// process hudson.plugins.git.util.BuildData
			for _, a := range body.Actions {
				if a.LastBuiltRevision == nil {
					continue
				}
				sha := ""
				branch := ""
				if a.LastBuiltRevision.SHA1 != "" {
					sha = a.LastBuiltRevision.SHA1
				}
				if a.MercurialRevisionNumber != "" {
					sha = a.MercurialRevisionNumber
				}

				if len(a.LastBuiltRevision.Branches) > 0 {
					branch = a.LastBuiltRevision.Branches[0].Name
				}
				for _, url := range a.RemoteUrls {
					if url != "" {
						buildCommitRemoteUrl := models.JenkinsBuildCommit{
							ConnectionId: data.Options.ConnectionId,
							BuildName:    build.FullName,
							CommitSha:    sha,
							RepoUrl:      url,
							Branch:       branch,
							Number:       body.Number,
						}
						results = append(results, &buildCommitRemoteUrl)
					}
				}
			}
			// process hudson.model.CauseAction
			for _, a := range body.Actions {
				if len(a.Causes) > 0 {
					for _, cause := range a.Causes {
						if cause.UpstreamProject != "" {
							triggeredByBuild := fmt.Sprintf("%s #%d", cause.UpstreamProject, cause.UpstreamBuild)
							build.TriggeredBy = triggeredByBuild
						}
					}
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
