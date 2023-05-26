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

type AdGroup struct {
	ID                    uint64   `gorm:"column:adgroup_id;primaryKey" json:"adgroup_id"`
	Status                string   `gorm:"column:status" json:"status"`
	ExternalAction        string   `gorm:"column:external_action" json:"external_action"`
	PixelID               *string  `gorm:"column:pixel_id" json:"pixel_id"`
	OptimizeGoal          string   `gorm:"column:optimize_goal" json:"optimize_goal"`
	AppID                 uint64   `gorm:"column:app_id" json:"app_id"`
	CreateTime            string   `gorm:"column:create_time" json:"create_time"`
	ConversionBid         float64  `gorm:"column:conversion_bid" json:"conversion_bid"`
	CampaignName          string   `gorm:"column:campaign_name" json:"campaign_name"`
	Keywords              []string `gorm:"column:keywords" json:"keywords"`
	OperationSystem       []string `gorm:"column:operation_system" json:"operation_system"`
	CreativeMaterialMode  string   `gorm:"column:creative_material_mode" json:"creative_material_mode"`
	Placement             []string `gorm:"column:placement" json:"placement"`
	ConnectionType        []string `gorm:"column:connection_type" json:"connection_type"`
	DeepCPABid            float64  `gorm:"column:deep_cpabid" json:"deep_cpabid"`
	ExternalType          string   `gorm:"column:external_type" json:"external_type"`
	AdGroupName           string   `gorm:"column:adgroup_name" json:"adgroup_name"`
	Languages             []string `gorm:"column:languages" json:"languages"`
	Location              []int    `gorm:"column:location" json:"location"`
	AppDownloadURL        string   `gorm:"column:app_download_url" json:"app_download_url"`
	BudgetMode            string   `gorm:"column:budget_mode" json:"budget_mode"`
	ExcludedAudience      []string `gorm:"column:excluded_audience" json:"excluded_audience"`
	EnableInventoryFilter bool     `gorm:"column:enable_inventory_filter" json:"enable_inventory_filter"`
	BillingEvent          string   `gorm:"column:billing_event" json:"billing_event"`
	Bid                   float64  `gorm:"column:bid" json:"bid"`
	CampaignID            uint64   `gorm:"column:campaign_id" json:"campaign_id"`
	AdvertiserID          uint64   `gorm:"column:advertiser_id" json:"advertiser_id"`
	DeepBidType           string   `gorm:"column:deep_bid_type" json:"deep_bid_type"`
	ScheduleStartTime     string   `gorm:"column:schedule_start_time" json:"schedule_start_time"`
	ScheduleEndTime       string   `gorm:"column:schedule_end_time" json:"schedule_end_time"`
	DevicePrice           *float64 `gorm:"column:device_price" json:"device_price"`
	IsCommentDisable      int      `gorm:"column:is_comment_disable" json:"is_comment_disable"`
	BidType               string   `gorm:"column:bid_type" json:"bid_type"`
	SkipLearningPhase     int      `gorm:"column:skip_learning_phase" json:"skip_learning_phase"`
	Package               string   `gorm:"column:package" json:"package"`
	Gender                *string  `gorm:"column:gender" json:"gender"`
	Age                   *string  `gorm:"column:age" json:"age"`
	PlacementType         string   `gorm:"column:placement_type" json:"placement_type"`
	Budget                float64  `gorm:"column:budget" json:"budget"`
	Pacing                string   `gorm:"column:pacing" json:"pacing"`
	ScheduleType          string   `gorm:"column:schedule_type" json:"schedule_type"`
	InterestCategoryV2    []string `gorm:"column:interest_category_v2" json:"interest_category_v2"`
	Audience              []string `gorm:"column:audience" json:"audience"`
	ModifyTime            string   `gorm:"column:modify_time" json:"modify_time"`
	DeepExternalAction    *string  `gorm:"column:deep_external_action" json:"deep_external_action"`
	Dayparting            *string  `gorm:"column:dayparting" json:"dayparting"`
}

type AdGroupList struct {
	List     []AdGroup      `json:"list"`
	PageInfo tasks.PageInfo `json:"page_info"`
}

type AdGroupResponse struct {
	Message   string      `json:"message"`
	Code      int         `json:"code"`
	Data      AdGroupList `json:"data"`
	RequestID string      `json:"request_id"`
}

func (AdGroup) TableName() string {
	return "_tool_tiktokAds_ad_group"
}
