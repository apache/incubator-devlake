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

package models

import (
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
)

// CopilotOrgMetrics captures daily organization-level Copilot adoption metrics.
type CopilotOrgMetrics struct {
	ConnectionId uint64    `gorm:"primaryKey" json:"connectionId"`
	ScopeId      string    `gorm:"primaryKey;type:varchar(255)" json:"scopeId"`
	Date         time.Time `gorm:"primaryKey;type:date" json:"date"`

	TotalActiveUsers         int `json:"totalActiveUsers"`
	TotalEngagedUsers        int `json:"totalEngagedUsers"`
	CompletionSuggestions    int `json:"completionSuggestions"`
	CompletionAcceptances    int `json:"completionAcceptances"`
	CompletionLinesSuggested int `json:"completionLinesSuggested"`
	CompletionLinesAccepted  int `json:"completionLinesAccepted"`
	IdeChats                 int `json:"ideChats"`
	IdeChatCopyEvents        int `json:"ideChatCopyEvents"`
	IdeChatInsertionEvents   int `json:"ideChatInsertionEvents"`
	IdeChatEngagedUsers      int `json:"ideChatEngagedUsers"`
	DotcomChats              int `json:"dotcomChats"`
	DotcomChatEngagedUsers   int `json:"dotcomChatEngagedUsers"`
	SeatActiveCount          int `json:"seatActiveCount"`
	SeatTotal                int `json:"seatTotal"`

	common.NoPKModel
}

func (CopilotOrgMetrics) TableName() string {
	return "_tool_copilot_org_metrics"
}
