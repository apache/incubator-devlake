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
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
)

var ConvertPrReviewersMeta = plugin.SubTaskMeta{
	Name:             "Convert PR Reviewers",
	EntryPoint:       ConvertPrReviewers,
	EnabledByDefault: true,
	Required:         false,
	Description:      "Convert tool layer PR reviewers to domain layer",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
}

func ConvertPrReviewers(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PULL_REQUEST_TABLE)
	db := taskCtx.GetDal()
	repoId := data.Options.FullName

	cursor, err := db.Cursor(
		dal.From(&models.BitbucketPrReviewer{}),
		dal.Where("repo_id = ? and connection_id = ?", repoId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	prIdGen := didgen.NewDomainIdGenerator(&models.BitbucketPullRequest{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.BitbucketAccount{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(models.BitbucketPrReviewer{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			reviewer := inputRow.(*models.BitbucketPrReviewer)
			domainReviewer := &code.PullRequestReviewer{
				PullRequestId: prIdGen.Generate(data.Options.ConnectionId, reviewer.RepoId, reviewer.PullRequestId),
				ReviewerId:    accountIdGen.Generate(data.Options.ConnectionId, reviewer.AccountId),
				Name:          reviewer.DisplayName,
				UserName:      reviewer.DisplayName,
			}
			return []interface{}{
				domainReviewer,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
