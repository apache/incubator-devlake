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
	"github.com/apache/incubator-devlake/errors"
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
)

const RAW_JOB_TABLE = "jenkins_api_jobs"

var CollectApiJobsMeta = core.SubTaskMeta{
	Name:             "collectApiJobs",
	EntryPoint:       CollectApiJobs,
	EnabledByDefault: true,
	Description:      "Collect jobs data from jenkins api",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func CollectApiJobs(taskCtx core.SubTaskContext) errors.Error {
	it := helper.NewQueueIterator()
	it.Push(models.NewFolderInput(""))
	data := taskCtx.GetData().(*JenkinsTaskData)
	incremental := false
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Params: JenkinsApiParams{
				ConnectionId: data.Options.ConnectionId,
			},
			Ctx:   taskCtx,
			Table: RAW_JOB_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Incremental: incremental,
		// jenkins api is special, 1. If the concurrency is larger than 1, then it will report 500.
		Concurrency: 1,

		UrlTemplate: "{{ .Input.Path }}api/json",
		Input:       it,
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			treeValue := fmt.Sprintf(
				"jobs[name,class,url,color,base,jobs,upstreamProjects[name]]{%d,%d}",
				reqData.Pager.Skip, reqData.Pager.Skip+reqData.Pager.Size)
			query.Set("tree", treeValue)
			return query, nil
		},
		Header: func(reqData *helper.RequestData) (http.Header, errors.Error) {
			input, ok := reqData.Input.(*models.FolderInput)
			if ok {
				return http.Header{
					"Path": {
						input.Path,
					},
				}, nil
			} else {
				return nil, errors.Default.New("empty FolderInput")
			}
		},

		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data struct {
				Jobs []json.RawMessage `json:"jobs"`
			}
			err := helper.UnmarshalResponse(res, &data)
			if err != nil {
				return nil, err
			}
			BasePath := res.Request.Header.Get("Path")
			for _, rawJobs := range data.Jobs {
				job := &models.Job{}
				err := errors.Convert(json.Unmarshal(rawJobs, job))
				if err != nil {
					return nil, err
				}

				if job.Jobs != nil {
					it.Push(models.NewFolderInput(BasePath + "job/" + job.Name + "/"))
				}
			}
			return data.Jobs, nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
