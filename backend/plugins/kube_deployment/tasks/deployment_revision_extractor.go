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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/kube_deployment/models"
)

var _ plugin.SubTaskEntryPoint = ExtractKubeDeploymentRevisions

func ExtractKubeDeploymentRevisions(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*KubeDeploymentTaskData)
	fmt.Println("ooooppppp")
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: KubeDeploymentApiParams{
				ConnectionId:   data.Options.ConnectionId,
				DeploymentName: data.Options.DeploymentName,
				Namespace:      data.Options.Namespace,
			},
			Table: RAW_KUBE_DEPLOYMENT_REVISION_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			println("Extracting KubeDeploymentRevisions %v", string(row.Data))
			body := &models.KubeDeploymentRevision{}
			err := errors.Convert(json.Unmarshal(row.Data, body))
			if err != nil {

				print("Error unmarshalling KubeDeploymentRevisions", err)
				return nil, err
			}
			println("Extracting KubeDeploymentRevisions", body)
			body.ConnectionId = data.Options.ConnectionId
			return []interface{}{body}, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractKubeDeploymentRevisionsMeta = plugin.SubTaskMeta{
	Name:             "extractKubeDeploymentRevisions",
	EntryPoint:       ExtractKubeDeploymentRevisions,
	EnabledByDefault: true,
	Description:      "Extract raw KubeDeploymentRevisions data into tool layer table",
}
