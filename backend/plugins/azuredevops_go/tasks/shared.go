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
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"net/http"
	"net/url"
)

// Build and TimeLine Record State and Result types can be found here:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/build/timeline/get?view=azure-devops-rest-7.1#taskresult
// https://learn.microsoft.com/en-us/rest/api/azure/devops/build/builds/get?view=azure-devops-rest-7.1#buildstatus
const (
	cancelling          string = "cancelling"
	completed           string = "completed"
	inProgress          string = "inProgress"
	notStarted          string = "notStarted"
	postponed           string = "postponed"
	canceled            string = "canceled"
	failed              string = "failed"
	none                string = "none"
	partiallySucceeded  string = "partiallySucceeded"
	succeeded           string = "succeeded"
	pending             string = "pending"
	abandoned           string = "abandoned"
	skipped             string = "skipped"
	succeededWithIssues string = "succeededWithIssues"
)

type CustomPageDate struct {
	ContinuationToken string
}

func CreateRawDataSubTaskArgs(taskCtx plugin.SubTaskContext, Table string) (*api.RawDataSubTaskArgs, *AzuredevopsTaskData) {
	data := taskCtx.GetData().(*AzuredevopsTaskData)
	RawDataSubTaskArgs := &api.RawDataSubTaskArgs{
		Ctx:     taskCtx,
		Options: data.Options,
		Table:   Table,
	}
	return RawDataSubTaskArgs, data
}

func ParseRawMessageFromValue(res *http.Response) ([]json.RawMessage, errors.Error) {
	var data struct {
		Value []json.RawMessage `json:"value"`
	}
	err := api.UnmarshalResponse(res, &data)
	if err != nil {
		return nil, err
	}
	return data.Value, nil
}

func ParseRawMessageFromRecords(res *http.Response) ([]json.RawMessage, errors.Error) {
	var data struct {
		Value []json.RawMessage `json:"records"`
	}
	err := api.UnmarshalResponse(res, &data)
	if err != nil {
		return nil, err
	}
	return data.Value, nil
}

func BuildPaginator(cursor bool) func(reqData *api.RequestData) (url.Values, errors.Error) {
	return func(reqData *api.RequestData) (url.Values, errors.Error) {
		query := url.Values{}
		if cursor && reqData.CustomData != nil {
			pag := reqData.CustomData.(CustomPageDate)
			query.Set("continuationToken", pag.ContinuationToken)
		} else {
			query.Set("$skip", fmt.Sprint(reqData.Pager.Skip))
		}
		query.Set("$top", fmt.Sprint(reqData.Pager.Size))
		return query, nil
	}
}

func ExtractContToken(_ *api.RequestData, prevPageResponse *http.Response) (interface{}, errors.Error) {
	t := prevPageResponse.Header.Get("X-Ms-Continuationtoken")
	if t == "" {
		return nil, api.ErrFinishCollect
	}
	return CustomPageDate{
		ContinuationToken: t,
	}, nil
}

var cicdBuildResultRule = devops.ResultRule{
	Success: []string{succeeded, partiallySucceeded},
	Failure: []string{canceled, failed, none},
	Default: devops.RESULT_DEFAULT,
}

var cicdBuildStatusRule = devops.StatusRule{
	Done:       []string{completed, cancelling},
	InProgress: []string{inProgress, notStarted, postponed},
	Default:    devops.STATUS_OTHER,
}

var cicdTaskResultRule = &devops.ResultRule{
	Success: []string{succeeded, succeededWithIssues},
	Failure: []string{abandoned, canceled, failed, skipped},
	Default: devops.RESULT_DEFAULT,
}

var cicdTaskStatusRule = &devops.StatusRule{
	Done:       []string{completed},
	InProgress: []string{pending, inProgress},
	Default:    devops.STATUS_OTHER,
}

func change203To401(res *http.Response) errors.Error {
	if res.StatusCode == http.StatusUnauthorized {
		return errors.Unauthorized.New("authentication failed, please check your AccessToken")
	}

	// When the token is invalid, Azure DevOps returns a 302 that resolves to a sign-in page with status 203
	// We want to change that to a 401 and raise an exception
	if res.StatusCode == http.StatusNonAuthoritativeInfo {
		return errors.Unauthorized.New("authentication failed, please check your AccessToken")
	}
	return nil
}

func ignoreDeletedBuilds(res *http.Response) errors.Error {
	if res.StatusCode == http.StatusNotFound {
		return api.ErrIgnoreAndContinue
	}
	return nil
}
