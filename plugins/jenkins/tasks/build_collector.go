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
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"net/http"
	"net/url"
)

const RAW_BUILD_TABLE = "jenkins_api_builds"

var CollectApiBuildsMeta = core.SubTaskMeta{
	Name:             "collectApiBuilds",
	EntryPoint:       CollectApiBuilds,
	EnabledByDefault: true,
	Description:      "Collect builds data from jenkins api",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

type SimpleJob struct {
	Name string
	Path string
}

func CollectApiBuilds(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JenkinsTaskData)

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Params: JenkinsApiParams{
				ConnectionId: data.Options.ConnectionId,
				FullName:     data.Options.JobFullName,
			},
			Ctx:   taskCtx,
			Table: RAW_BUILD_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    100,
		UrlTemplate: fmt.Sprintf("%sjob/%s/api/json", data.Options.JobPath, data.Options.JobName),
		/*
			(Optional) Return query string for request, or you can plug them into UrlTemplate directly
		*/
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			treeValue := fmt.Sprintf(
				"allBuilds[number,timestamp,duration,building,estimatedDuration,fullName,fullDisplayName,result,actions[lastBuiltRevision[SHA1,branch[name]],remoteUrls,mercurialRevisionNumber,causes[*]],changeSet[kind,revisions[revision]]]{%d,%d}",
				reqData.Pager.Skip, reqData.Pager.Skip+reqData.Pager.Size)
			query.Set("tree", treeValue)
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data struct {
				Builds []json.RawMessage `json:"allBuilds"`
			}
			err := helper.UnmarshalResponse(res, &data)
			if err != nil {
				return nil, err
			}
			return data.Builds, nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
