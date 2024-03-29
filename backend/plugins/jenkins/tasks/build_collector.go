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
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_BUILD_TABLE = "jenkins_api_builds"

var CollectApiBuildsMeta = plugin.SubTaskMeta{
	Name:             "collectApiBuilds",
	EntryPoint:       CollectApiBuilds,
	EnabledByDefault: true,
	Description:      "Collect builds data from jenkins api, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

type SimpleJob struct {
	Name string
	Path string
}

type SimpleJenkinsApiBuild struct {
	Number    int64
	Timestamp int64 `json:"timestamp"`
}

func CollectApiBuilds(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JenkinsTaskData)
	collector, err := helper.NewStatefulApiCollectorForFinalizableEntity(helper.FinalizableApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Params: JenkinsApiParams{
				ConnectionId: data.Options.ConnectionId,
				FullName:     data.Options.JobFullName,
			},
			Ctx:   taskCtx,
			Table: RAW_BUILD_TABLE,
		},
		ApiClient: data.ApiClient,
		CollectNewRecordsByList: helper.FinalizableApiCollectorListArgs{
			PageSize:    100,
			Concurrency: 10,
			FinalizableApiCollectorCommonArgs: helper.FinalizableApiCollectorCommonArgs{
				UrlTemplate: fmt.Sprintf("%sjob/%s/api/json", data.Options.JobPath, data.Options.JobName),
				Query: func(reqData *helper.RequestData, createdAfter *time.Time) (url.Values, errors.Error) {
					query := url.Values{}
					actionFields := "lastBuiltRevision[SHA1,branch[name]],remoteUrls,mercurialRevisionNumber,causes[*],parameters[name,value]"
					treeValue := fmt.Sprintf(
						"allBuilds[timestamp,number,duration,building,estimatedDuration,fullDisplayName,result,actions[%s],changeSet[kind,revisions[revision]]]{%d,%d}",
						actionFields, reqData.Pager.Skip, reqData.Pager.Skip+reqData.Pager.Size)
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

					builds := make([]json.RawMessage, 0, len(data.Builds))
					for _, build := range data.Builds {
						var buildObj map[string]interface{}
						err := json.Unmarshal(build, &buildObj)
						if err != nil {
							return nil, errors.Convert(err)
						}
						if buildObj["result"] != nil {
							builds = append(builds, build)
						}
					}

					return builds, nil
				},
			},
			GetCreated: func(item json.RawMessage) (time.Time, errors.Error) {
				b := &SimpleJenkinsApiBuild{}
				err := json.Unmarshal(item, b)
				if err != nil {
					return time.Time{}, errors.BadInput.Wrap(err, "failed to unmarshal jenkins build")
				}
				seconds := b.Timestamp / 1000
				nanos := (b.Timestamp % 1000) * 1000000
				return time.Unix(seconds, nanos), nil
			},
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
