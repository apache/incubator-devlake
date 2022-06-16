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

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractProjectMeta = core.SubTaskMeta{
	Name:             "extractApiProject",
	EntryPoint:       ExtractApiProject,
	EnabledByDefault: true,
	Description:      "Extract raw project data into tool layer table GitlabProject",
}

func ExtractApiProject(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECT_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			// create gitlab commit
			gitlabApiProject := &GitlabApiProject{}
			err := json.Unmarshal(row.Data, gitlabApiProject)
			if err != nil {
				return nil, err
			}
			gitlabProject := convertProject(gitlabApiProject)
			gitlabProject.ConnectionId = data.Options.ConnectionId
			results := make([]interface{}, 0, 1)
			results = append(results, gitlabProject)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
