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

package apiv2models

import (
	"time"

	"github.com/apache/incubator-devlake/plugins/jira/models"
)

type Sprint struct {
	ID            uint64     `json:"id"`
	Self          string     `json:"self"`
	State         string     `json:"state"`
	Name          string     `json:"name"`
	StartDate     *time.Time `json:"startDate"`
	EndDate       *time.Time `json:"endDate"`
	CompleteDate  *time.Time `json:"completeDate"`
	OriginBoardID uint64     `json:"originBoardId"`
}

func (s Sprint) ToToolLayer(connectionId uint64) *models.JiraSprint {
	return &models.JiraSprint{
		ConnectionId:  connectionId,
		SprintId:      s.ID,
		Self:          s.Self,
		State:         s.State,
		Name:          s.Name,
		StartDate:     s.StartDate,
		EndDate:       s.EndDate,
		CompleteDate:  s.CompleteDate,
		OriginBoardID: s.OriginBoardID,
	}
}
