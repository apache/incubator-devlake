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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_PULL_REQUEST_COMMENTS_TABLE = "bitbucket_api_pull_request_comments"

var CollectApiPrCommentsMeta = plugin.SubTaskMeta{
	Name:             "collectApiPullRequestsComments",
	EntryPoint:       CollectApiPullRequestsComments,
	EnabledByDefault: true,
	Required:         false,
	Description:      "Collect pull requests comments data from bitbucket api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
}

func CollectApiPullRequestsComments(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PULL_REQUEST_COMMENTS_TABLE)
	collectorWithState, err := helper.NewStatefulApiCollector(*rawDataSubTaskArgs, data.TimeAfter)
	if err != nil {
		return err
	}

	iterator, err := GetPullRequestsIterator(taskCtx, collectorWithState)
	if err != nil {
		return err
	}
	defer iterator.Close()

	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Incremental: collectorWithState.IsIncremental(),
		Input:       iterator,
		UrlTemplate: "repositories/{{ .Params.FullName }}/pullrequests/{{ .Input.BitbucketId }}/comments",
		Query: GetQueryFields(
			`values.id,values.type,values.created_on,values.updated_on,values.content.raw,values.pullrequest.id,values.user,` +
				`page,pagelen,size`),
		GetTotalPages:  GetTotalPagesFromResponse,
		ResponseParser: GetRawMessageFromResponse,
	})
	if err != nil {
		return err
	}

	return collectorWithState.Execute()
}
