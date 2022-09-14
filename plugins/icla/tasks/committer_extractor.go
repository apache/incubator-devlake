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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/icla/models"
)

var _ core.SubTaskEntryPoint = ExtractCommitter

func ExtractCommitter(taskCtx core.SubTaskContext) errors.Error {
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:    taskCtx,
			Params: IclaApiParams{},
			Table:  RAW_COMMITTER_TABLE,
		},
		Extract: func(resData *helper.RawData) ([]interface{}, errors.Error) {
			names := &map[string]string{}
			err := errors.Convert(json.Unmarshal(resData.Data, names))
			if err != nil {
				return nil, err
			}
			extractedModels := make([]interface{}, 0)
			for userName, name := range *names {
				extractedModels = append(extractedModels, &models.IclaCommitter{
					UserName: userName,
					Name:     name,
				})
			}
			return extractedModels, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractCommitterMeta = core.SubTaskMeta{
	Name:             "ExtractCommitter",
	EntryPoint:       ExtractCommitter,
	EnabledByDefault: true,
	Description:      "Extract raw data into tool layer table {{ .plugin_name }}_{{ .extractor_data_name }}",
}
