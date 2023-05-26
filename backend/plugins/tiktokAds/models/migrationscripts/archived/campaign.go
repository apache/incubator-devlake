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

package archived

import "github.com/apache/incubator-devlake/plugins/tiktokAds/tasks"

type Campaign struct {
	ID              uint64  `gorm:"column:campaign_id;primaryKey" json:"campaign_id"`
	SecondaryStatus string  `gorm:"column:secondary_status" json:"secondary_status"`
	Budget          float64 `gorm:"column:budget" json:"budget"`
	DeepBidType     *string `gorm:"column:deep_bid_type" json:"deep_bid_type"`
	AdvertiserID    string  `gorm:"column:advertiser_id" json:"advertiser_id"`
	RoasBid         float64 `gorm:"column:roas_bid" json:"roas_bid"`
	OperationStatus string  `gorm:"column:operation_status" json:"operation_status"`
	ObjectiveType   string  `gorm:"column:objective_type" json:"objective_type"`
	IsNewStructure  bool    `gorm:"column:is_new_structure" json:"is_new_structure"`
	CampaignName    string  `gorm:"column:campaign_name" json:"campaign_name"`
	CampaignType    string  `gorm:"column:campaign_type" json:"campaign_type"`
	BudgetMode      string  `gorm:"column:budget_mode" json:"budget_mode"`
	Objective       string  `gorm:"column:objective" json:"objective"`
	ModifyTime      string  `gorm:"column:modify_time" json:"modify_time"`
	CreateTime      string  `gorm:"column:create_time" json:"create_time"`
}

type CampaignList struct {
	List     []Campaign     `json:"list"`
	PageInfo tasks.PageInfo `json:"page_info"`
}

type CampaignResponse struct {
	Code      int          `json:"code"`
	Message   string       `json:"message"`
	RequestID string       `json:"request_id"`
	Data      CampaignList `json:"data"`
}

func (Campaign) TableName() string {
	return "_tool_tiktokAds_campaign"
}
