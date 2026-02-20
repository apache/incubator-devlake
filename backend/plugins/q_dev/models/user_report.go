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

type QDevUserReport struct {
	common.Model
	ConnectionId       uint64    `gorm:"primaryKey"`
	UserId             string    `gorm:"index" json:"userId"`
	Date               time.Time `gorm:"index" json:"date"`
	DisplayName        string    `gorm:"type:varchar(255)" json:"displayName"`
	ScopeId            string    `gorm:"index;type:varchar(255)" json:"scopeId"`
	ClientType         string    `gorm:"type:varchar(50)" json:"clientType"`
	SubscriptionTier   string    `gorm:"type:varchar(50)" json:"subscriptionTier"`
	ProfileId          string    `gorm:"type:varchar(512)" json:"profileId"`
	ChatConversations  int       `json:"chatConversations"`
	CreditsUsed        float64   `json:"creditsUsed"`
	OverageCap         float64   `json:"overageCap"`
	OverageCreditsUsed float64   `json:"overageCreditsUsed"`
	OverageEnabled     bool      `json:"overageEnabled"`
	TotalMessages      int       `json:"totalMessages"`
}

func (QDevUserReport) TableName() string {
	return "_tool_q_dev_user_report"
}
