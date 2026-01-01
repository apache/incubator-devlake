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
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

var CollectApiJobsMeta = plugin.SubTaskMeta{
	Name:             "collectApiJobs",
	EntryPoint:       CollectApiJobs,
	EnabledByDefault: true,
	Description:      "Collect jobs data from multibranch projects using jenkins api, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func CollectApiJobs(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JenkinsTaskData)
	logger := taskCtx.GetLogger()

	if data.Options.Class != WORKFLOW_MULTI_BRANCH_PROJECT {
		logger.Debug("class must be %s, got %s", WORKFLOW_MULTI_BRANCH_PROJECT, data.Options.Class)
		return nil
	}

	collector, err := helper.NewStatefulApiCollectorForFinalizableEntity(helper.FinalizableApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Params: JenkinsApiParams{
				ConnectionId: data.Options.ConnectionId,
				FullName:     data.Options.JobFullName,
			},
			Ctx:   taskCtx,
			Table: RAW_JOB_TABLE,
		},
		ApiClient: data.ApiClient,
		CollectNewRecordsByList: helper.FinalizableApiCollectorListArgs{
			FinalizableApiCollectorCommonArgs: helper.FinalizableApiCollectorCommonArgs{
				UrlTemplate: fmt.Sprintf("%sjob/%s/api/json", data.Options.JobPath, data.Options.JobName),
				Query: func(reqData *helper.RequestData, createdAfter *time.Time) (url.Values, errors.Error) {
					query := url.Values{}
					treeValue := "jobs[fullName,name,class,url,color,description]"
					query.Set("tree", treeValue)
					return query, nil
				},
				ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
					var resData struct {
						Jobs []json.RawMessage `json:"jobs"`
					}
					err := helper.UnmarshalResponse(res, &resData)
					if err != nil {
						return nil, err
					}

					// Compile branch filter pattern once for this batch
					var branchPattern *regexp.Regexp
					if data.Options.ScopeConfig != nil && data.Options.ScopeConfig.BranchFilterPattern != "" {
						var compileErr error
						branchPattern, compileErr = regexp.Compile(data.Options.ScopeConfig.BranchFilterPattern)
						if compileErr != nil {
							logger.Warn(nil, "Invalid branch filter pattern: %s, will include all jobs", data.Options.ScopeConfig.BranchFilterPattern)
						}
					}

					jobs := make([]json.RawMessage, 0, len(resData.Jobs))
					for _, job := range resData.Jobs {
						var jobObj map[string]interface{}
						err := json.Unmarshal(job, &jobObj)
						if err != nil {
							return nil, errors.Convert(err)
						}

						logger.Debug("%v", jobObj)
						if jobObj["color"] != "notbuilt" && jobObj["color"] != "nobuilt_anime" {
							// Apply branch filter pattern if configured
							if shouldIncludeJob(jobObj, branchPattern, logger) {
								jobs = append(jobs, job)
							}
						}
					}

					return jobs, nil
				},
			},
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}

// shouldIncludeJob determines whether a job should be included based on the branch filter pattern
func shouldIncludeJob(jobObj map[string]interface{}, branchPattern *regexp.Regexp, logger log.Logger) bool {
	// If no branch filter pattern is configured, include all jobs
	if branchPattern == nil {
		return true
	}

	// Get the job name for pattern matching
	jobName, ok := jobObj["name"].(string)
	if !ok {
		// If we can't get the job name, include it by default
		logger.Warn(nil, "Could not extract job name for filtering, including job by default")
		return true
	}

	// Match the job name against the pattern
	matched := branchPattern.MatchString(jobName)
	logger.Debug("Job '%s' %s branch filter pattern", jobName,
		map[bool]string{true: "matches", false: "does not match"}[matched])

	return matched
}
