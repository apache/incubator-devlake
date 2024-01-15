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
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket_server/models"
)

var ExtractApiBranchesMeta = plugin.SubTaskMeta{
	Name:             "extractApiBranches",
	EntryPoint:       ExtractApiBranches,
	EnabledByDefault: true,
	Required:         false,
	Description:      "Extract raw branch data into tool layer table bitbucket_branches",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
}

type BranchResponse struct {
	BitbucketId     string `json:"id"`
	DisplayId       string `json:"displayId"`
	Type            string `json:"type"`
	LatestCommit    string `json:"latestCommit"`
	LatestChangeset string `json:"latestChangeset"`
	IsDefault       bool   `json:"isDefault"`
}

func ExtractApiBranches(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_BRANCHES_TABLE)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			branch := &BranchResponse{}
			err := errors.Convert(json.Unmarshal(row.Data, branch))
			if err != nil {
				return nil, err
			} else if strings.ToLower(branch.Type) != "branch" {
				return []interface{}{}, nil
			}

			results := make([]interface{}, 0, 4)

			bitbucketBranch := &models.BitbucketServerBranch{
				Branch: branch.BitbucketId,

				ConnectionId: data.Options.ConnectionId,
				RepoId:       data.Options.FullName,

				LatestCommit: branch.LatestCommit,
				IsDefault:    branch.IsDefault,
			}

			results = append(results, bitbucketBranch)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
