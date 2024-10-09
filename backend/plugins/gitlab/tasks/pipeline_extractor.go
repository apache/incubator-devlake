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
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiPipelinesMeta)
}

type ApiDetailedStatus struct {
	Icon        string
	Text        string
	Label       string
	Group       string
	Tooltip     string
	HasDetails  bool   `json:"has_details"`
	DetailsPath string `json:"details_path"`
	Favicon     string
}

type ApiPipeline struct {
	Id             int `json:"id"`
	Ref            string
	Sha            string
	Status         string
	Tag            bool
	Duration       int
	QueuedDuration *float64 `json:"queued_duration"`
	WebUrl         string   `json:"web_url"`
	Source         string   `json:"source"`

	CreatedAt  *common.Iso8601Time `json:"created_at"`
	UpdatedAt  *common.Iso8601Time `json:"updated_at"`
	StartedAt  *common.Iso8601Time `json:"started_at"`
	FinishedAt *common.Iso8601Time `json:"finished_at"`

	ApiDetailedStatus
}

var ExtractApiPipelinesMeta = plugin.SubTaskMeta{
	Name:             "Extract Pipelines",
	EntryPoint:       ExtractApiPipelines,
	EnabledByDefault: true,
	Description:      "Extract raw pipelines data into tool layer table GitlabPipeline",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	Dependencies:     []*plugin.SubTaskMeta{&CollectApiPipelinesMeta},
}

func ExtractApiPipelines(subtaskCtx plugin.SubTaskContext) errors.Error {
	subtaskCommonArgs, data := CreateSubtaskCommonArgs(subtaskCtx, RAW_PIPELINE_TABLE)

	extractor, err := api.NewStatefulApiExtractor(&api.StatefulApiExtractorArgs[ApiPipeline]{
		SubtaskCommonArgs: subtaskCommonArgs,
		Extract: func(gitlabApiPipeline *ApiPipeline, row *api.RawData) ([]interface{}, errors.Error) {
			pipelineProject := convertApiPipelineToGitlabPipelineProject(gitlabApiPipeline, data.Options.ConnectionId, data.Options.ProjectId)
			return []interface{}{pipelineProject}, nil
		},
	})

	if err != nil {
		return err
	}

	err = extractor.Execute()
	if err != nil {
		return err
	}

	return nil
}
