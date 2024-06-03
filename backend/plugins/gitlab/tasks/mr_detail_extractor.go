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
	"regexp"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiMergeRequestDetailsMeta)
}

var ExtractApiMergeRequestDetailsMeta = plugin.SubTaskMeta{
	Name:             "Extract MR Details",
	EntryPoint:       ExtractApiMergeRequestDetails,
	EnabledByDefault: true,
	Description:      "Extract raw merge request Details data into tool layer table GitlabMergeRequest and GitlabReviewer",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
	Dependencies:     []*plugin.SubTaskMeta{&ExtractApiCommitsMeta},
}

func ExtractApiMergeRequestDetails(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_MERGE_REQUEST_DETAIL_TABLE)
	config := data.Options.ScopeConfig
	var labelTypeRegex *regexp.Regexp
	var labelComponentRegex *regexp.Regexp
	var prType = config.PrType
	var err error
	if len(prType) > 0 {
		labelTypeRegex, err = regexp.Compile(prType)
		if err != nil {
			return errors.Default.Wrap(err, "regexp Compile prType failed")
		}
	}
	var prComponent = config.PrComponent
	if len(prComponent) > 0 {
		labelComponentRegex, err = regexp.Compile(prComponent)
		if err != nil {
			return errors.Default.Wrap(err, "regexp Compile prComponent failed")
		}
	}
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			mr := &MergeRequestRes{}
			err := errors.Convert(json.Unmarshal(row.Data, mr))
			if err != nil {
				return nil, err
			}

			gitlabMergeRequest, err := convertMergeRequest(mr)
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0, len(mr.Reviewers)+1)
			gitlabMergeRequest.ConnectionId = data.Options.ConnectionId
			gitlabMergeRequest.IsDetailRequired = true
			results = append(results, gitlabMergeRequest)
			for _, label := range mr.Labels {
				results = append(results, &models.GitlabMrLabel{
					MrId:         gitlabMergeRequest.GitlabId,
					LabelName:    label,
					ConnectionId: data.Options.ConnectionId,
				})
				// if pr.Type has not been set and prType is set in .env, process the below
				if labelTypeRegex != nil && labelTypeRegex.MatchString(label) {
					gitlabMergeRequest.Type = label
				}
				// if pr.Component has not been set and prComponent is set in .env, process
				if labelComponentRegex != nil && labelComponentRegex.MatchString(label) {
					gitlabMergeRequest.Component = label
				}
			}
			for _, reviewer := range mr.Reviewers {
				gitlabReviewer := &models.GitlabReviewer{
					ConnectionId:   data.Options.ConnectionId,
					ReviewerId:     reviewer.ReviewerIdId,
					MergeRequestId: mr.GitlabId,
					ProjectId:      data.Options.ProjectId,
					Username:       reviewer.Username,
					Name:           reviewer.Name,
					State:          reviewer.State,
					AvatarUrl:      reviewer.AvatarUrl,
					WebUrl:         reviewer.WebUrl,
				}
				results = append(results, gitlabReviewer)
			}

			for _, assignee := range mr.Assignees {
				gitlabAssignee := &models.GitlabAssignee{
					ConnectionId:   data.Options.ConnectionId,
					AssigneeId:     assignee.AssigneeId,
					MergeRequestId: mr.GitlabId,
					ProjectId:      data.Options.ProjectId,
					Username:       assignee.Username,
					Name:           assignee.Name,
					State:          assignee.State,
					AvatarUrl:      assignee.AvatarUrl,
					WebUrl:         assignee.WebUrl,
				}
				results = append(results, gitlabAssignee)
			}

			return results, nil
		},
	})

	if err != nil {
		return errors.Convert(err)
	}

	return extractor.Execute()
}
