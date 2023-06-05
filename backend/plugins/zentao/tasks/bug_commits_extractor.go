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
	"net/url"
	"regexp"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

var _ plugin.SubTaskEntryPoint = ExtractBugCommits

var ExtractBugCommitsMeta = plugin.SubTaskMeta{
	Name:             "extractBugCommits",
	EntryPoint:       ExtractBugCommits,
	EnabledByDefault: true,
	Description:      "extract Zentao bug commits",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ExtractBugCommits(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)

	// this Extract only work for product
	if data.Options.ProductId == 0 {
		return nil
	}

	re := regexp.MustCompile(`href='(.*?)'`)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ZentaoApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProductId:    data.Options.ProductId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_BUG_COMMITS_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			res := &models.ZentaoBugCommitsRes{}
			err := json.Unmarshal(row.Data, res)
			if err != nil {
				return nil, errors.Default.WrapRaw(err)
			}
			// only linked2revision action is valid
			if res.Action != "linked2revision" {
				return nil, nil
			}

			bugCommits := &models.ZentaoBugCommit{
				ConnectionId: data.Options.ConnectionId,
				ID:           res.ID,
				ObjectType:   res.ObjectType,
				ObjectID:     res.ObjectID,
				Product:      data.Options.ProductId,
				Project:      data.Options.ProjectId,
				Execution:    res.Execution,
				Actor:        res.Actor,
				Action:       res.Action,
				Date:         res.Date,
				Comment:      res.Comment,
				ActionRead:   res.Read,
				Vision:       res.Vision,
				Efforted:     res.Efforted,
				ActionDesc:   res.Desc,
			}

			match := re.FindStringSubmatch(res.Extra)
			if len(match) > 1 {
				bugCommits.Extra = match[1]
			} else {
				return nil, nil
			}
			u, err := url.Parse(match[1])
			if err != nil {
				return nil, errors.Default.WrapRaw(err)
			}
			bugCommits.Host = u.Host
			bugCommits.RepoRevision = u.Path

			results := make([]interface{}, 0)
			results = append(results, bugCommits)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
