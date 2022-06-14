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
	"net/http"

	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"gorm.io/gorm"
)

func GetTotalPagesFromResponse(res *http.Response, args *helper.ApiCollectorArgs) (int, error) {
	body := &JiraPagination{}
	err := helper.UnmarshalResponse(res, body)
	if err != nil {
		return 0, err
	}
	pages := body.Total / args.PageSize
	if body.Total%args.PageSize > 0 {
		pages++
	}
	return pages, nil
}

func GetStdStatus(statusKey string) string {
	if statusKey == "done" {
		return ticket.DONE
	} else if statusKey == "new" {
		return ticket.TODO
	} else {
		return ticket.IN_PROGRESS
	}
}

func GetStatusInfo(db *gorm.DB) ([]models.JiraStatus, error) {
	return nil, nil
}
