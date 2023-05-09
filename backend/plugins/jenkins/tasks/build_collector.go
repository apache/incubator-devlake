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
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
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
	db := taskCtx.GetDal()
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
		TimeAfter: data.TimeAfter,
		CollectNewRecordsByList: helper.FinalizableApiCollectorListArgs{
			PageSize:    100,
			Concurrency: 10,
			FinalizableApiCollectorCommonArgs: helper.FinalizableApiCollectorCommonArgs{
				UrlTemplate: fmt.Sprintf("%sjob/%s/api/json", data.Options.JobPath, data.Options.JobName),
				Query: func(reqData *helper.RequestData, createdAfter *time.Time) (url.Values, errors.Error) {
					query := url.Values{}
					treeValue := fmt.Sprintf(
						"allBuilds[timestamp,number,duration,building,estimatedDuration,fullDisplayName,result,actions[lastBuiltRevision[SHA1,branch[name]],remoteUrls,mercurialRevisionNumber,causes[*]],changeSet[kind,revisions[revision]]]{%d,%d}",
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
		CollectUnfinishedDetails: helper.FinalizableApiCollectorDetailArgs{
			BuildInputIterator: func() (helper.Iterator, errors.Error) {
				cursor, err := db.Cursor(
					dal.Select("number"),
					dal.From(&models.JenkinsBuild{}),
					dal.Where(
						"full_name = ? AND connection_id = ? AND result != 'SUCCESS' AND result != 'FAILURE'",
						data.Options.JobFullName, data.Options.ConnectionId,
					),
				)
				if err != nil {
					return nil, err
				}
				return helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(SimpleJenkinsApiBuild{}))
			},
			FinalizableApiCollectorCommonArgs: helper.FinalizableApiCollectorCommonArgs{
				UrlTemplate: fmt.Sprintf("%sjob/%s/{{ .Input.Number }}/api/json?tree=number,url,result,timestamp,id,duration,estimatedDuration,building",
					data.Options.JobPath, data.Options.JobName),
				ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
					body, err := io.ReadAll(res.Body)
					if err != nil {
						return nil, errors.Convert(err)
					}
					res.Body.Close()
					return []json.RawMessage{body}, nil
				},
			},
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
