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
)

const RAW_FINDINGS_TABLE = "checkmarxone_api_findings"

var CollectFindingsMeta = plugin.SubTaskMeta{
	Name:             "collectFindings",
	EntryPoint:       CollectFindings,
	EnabledByDefault: true,
	Description:      "Collect findings from CheckmarxOne API",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_SECURITY},
}

func CollectFindings(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*CheckmarxoneTaskData)
	logger := taskCtx.GetLogger()

	findings, err := data.ApiClient.GetFindings(data.Options.ProjectId)
	if err != nil {
		logger.Error(err, "failed to fetch findings")
		return err
	}

	for _, finding := range findings {
		select {
		case <-taskCtx.GetContext().Done():
			return taskCtx.GetContext().Err()
		default:
		}

		err := taskCtx.SaveRawData(RAW_FINDINGS_TABLE, finding)
		if err != nil {
			logger.Error(err, "failed to save raw data")
			return err
		}
	}

	return nil
}