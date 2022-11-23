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
	"github.com/apache/incubator-devlake/errors"
	"net/url"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const RAW_BUG_STATUS_LAST_STEP_TABLE = "tapd_api_bug_status_last_steps"

var _ core.SubTaskEntryPoint = CollectBugStatusLastStep

func CollectBugStatusLastStep(taskCtx core.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_BUG_STATUS_LAST_STEP_TABLE, false)
	logger := taskCtx.GetLogger()
	logger.Info("collect bugStatus")
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,

		UrlTemplate: "workflows/last_steps",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
			query.Set("system", "bug")
			return query, nil
		},
		ResponseParser: GetRawMessageDirectFromResponse,
	})
	if err != nil {
		logger.Error(err, "collect bug workflow last steps")
		return err
	}
	return collector.Execute()
}

var CollectBugStatusLastStepMeta = core.SubTaskMeta{
	Name:             "collectBugStatusLastStep",
	EntryPoint:       CollectBugStatusLastStep,
	EnabledByDefault: true,
	Description:      "collect Tapd bugStatus",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}
