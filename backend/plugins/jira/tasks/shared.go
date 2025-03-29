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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"net/http"
)

func GetTotalPagesFromResponse(res *http.Response, args *api.ApiCollectorArgs) (int, errors.Error) {
	body := &JiraPagination{}
	err := api.UnmarshalResponse(res, body)
	if err != nil {
		return 0, err
	}
	pages := body.Total / args.PageSize
	if body.Total%args.PageSize > 0 {
		pages++
	}
	return pages, nil
}

func getStdStatus(statusKey string) string {
	if statusKey == "done" {
		return ticket.DONE
	} else if statusKey == "new" {
		return ticket.TODO
	} else {
		return ticket.IN_PROGRESS
	}
}

func isServer(jiraServerInfo *models.JiraServerInfo, apiclient *api.ApiAsyncClient, db dal.Dal, connectionID uint64) (bool, errors.Error) {
	if jiraServerInfo != nil {
		return jiraServerInfo.IsDeploymentServer(), nil
	}
	// try to fetch jiraServerInfo from remote api
	if apiclient != nil {
		info, code, err := GetJiraServerInfo(apiclient)
		if err != nil || code != http.StatusOK || info == nil {
			return false, errors.HttpStatus(code).Wrap(err, "fail to get Jira server info")
		}
		return info.IsDeploymentServer(), nil
	}
	// fetch from db
	info, err := getJiraServerInfoFromDB(db, connectionID)
	if err != nil {
		return false, err
	}
	if info == nil {
		return false, nil
	}
	return info.IsDeploymentServer(), nil
}

func getJiraServerInfoFromDB(db dal.Dal, connectionID uint64) (*models.JiraServerInfo, errors.Error) {
	var info models.JiraServerInfo
	if err := db.First(&info, dal.Where("connection_id = ?", connectionID)); err != nil {
		if db.IsErrorNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &info, nil
}
