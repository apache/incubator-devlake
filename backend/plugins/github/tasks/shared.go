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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/github/utils"
	"net/http"
)

const (
	// https://docs.github.com/en/graphql/reference/enums#checkrunstate
	// https://docs.github.com/en/graphql/reference/enums#checkstatusstate
	// https://docs.github.com/en/graphql/reference/enums#deploymentstate
	// https://docs.github.com/en/graphql/reference/enums#checkconclusionstate
	StatusCompleted      = "COMPLETED"
	StatusSuccess        = "SUCCESS"
	StatusFailure        = "FAILURE"
	StatusCancelled      = "CANCELLED"
	StatusTimedOut       = "TIMED_OUT"
	StatusStartUpFailure = "STARTUP_FAILURE"
	StatusInProgress     = "IN_PROGRESS "
	StatusQueued         = "QUEUED"
	StatusWaiting        = "WAITING "
	StatusPending        = "PENDING"
	StatusNeutral        = "NEUTRAL"
	StatusSkipped        = "SKIPPED"
	StatusStale          = "STALE"
	StatusActionRequired = "ACTION_REQUIRED"
	StatusRequested      = "REQUESTED"
	StatusAbandoned      = "ABANDONED"
	StatusActive         = "ACTIVE"
	StatusInactive       = "INACTIVE"
	StatusDestroyed      = "DESTROYED"
	StatusError          = "ERROR"
)

func GetTotalPagesFromResponse(res *http.Response, args *api.ApiCollectorArgs) (int, errors.Error) {
	link := res.Header.Get("link")
	pageInfo, err := utils.GetPagingFromLinkHeader(link)
	if err != nil {
		return 0, nil
	}
	return pageInfo.Last, nil
}

func ignoreHTTPStatus404(res *http.Response) errors.Error {
	if res.StatusCode == http.StatusUnauthorized {
		return errors.Unauthorized.New("authentication failed, please check your AccessToken")
	}
	if res.StatusCode == http.StatusNotFound {
		return api.ErrIgnoreAndContinue
	}
	return nil
}

func ignoreHTTPStatus422(res *http.Response) errors.Error {
	if res.StatusCode == http.StatusUnprocessableEntity {
		return api.ErrIgnoreAndContinue
	}
	return nil
}

func CreateRawDataSubTaskArgs(taskCtx plugin.SubTaskContext, table string) (*api.RawDataSubTaskArgs, *GithubTaskData) {
	data := taskCtx.GetData().(*GithubTaskData)
	RawDataSubTaskArgs := &api.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: models.GithubApiParams{
			Name:         data.Options.Name,
			ConnectionId: data.Options.ConnectionId,
		},
		Table: table,
	}
	return RawDataSubTaskArgs, data
}
