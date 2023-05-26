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

import (
	"github.com/apache/incubator-devlake/plugins/tiktokAds/tasks"
	"time"
)

type Metrics struct {
	AdName      string  `json:"ad_name" gorm:"column:ad_name"`
	Clicks      float64 `json:"clicks" gorm:"column:clicks"`
	Spend       float64 `json:"spend" gorm:"column:spend"`
	Impressions float64 `json:"impressions" gorm:"column:impressions"`
	CTR         float64 `json:"ctr" gorm:"column:ctr"`
}

type Dimensions struct {
	AdvertiserID       string `json:"advertiser_id,omitempty"`
	CampaignID         string `json:"campaign_id,omitempty"`
	AdGroupID          string `json:"adgroup_id,omitempty"`
	AdID               string `json:"ad_id,omitempty"`
	StatTimeDay        string `json:"stat_time_day,omitempty"`
	StatTimeHour       string `json:"stat_time_hour,omitempty"`
	AC                 string `json:"ac,omitempty"`
	Age                string `json:"age,omitempty"`
	CountryCode        string `json:"country_code,omitempty"`
	InterestCategory   string `json:"interest_category,omitempty"`
	InterestCategoryV2 string `json:"interest_category_v2,omitempty"`
	Gender             string `json:"gender,omitempty"`
	Language           string `json:"language,omitempty"`
	Placement          string `json:"placement,omitempty"`
	Platform           string `json:"platform,omitempty"`
	ContextualTag      string `json:"contextual_tag,omitempty"`
}

type ListItem struct {
	Metrics    Metrics    `json:"metrics"`
	Dimensions Dimensions `json:"dimensions"`
}

type Data struct {
	PageInfo tasks.PageInfo `json:"page_info"`
	List     []ListItem     `json:"list"`
}

type Response struct {
	Message   string    `json:"message"`
	Code      int       `json:"code"`
	Data      Data      `json:"data"`
	RequestID string    `json:"request_id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Report struct {
	// dimension
	AdvertiserID       string `json:"advertiser_id,omitempty" gorm:"column:advertiser_id;primaryKey"`
	CampaignID         string `json:"campaign_id,omitempty" gorm:"column:campaign_id;primaryKey"`
	AdGroupID          string `json:"adgroup_id,omitempty" gorm:"column:adgroup_id;primaryKey"`
	AdID               string `json:"ad_id,omitempty" gorm:"column:ad_id;primaryKey"`
	StatTimeDay        string `json:"stat_time_day,omitempty" gorm:"column:stat_time_day"`
	StatTimeHour       string `json:"stat_time_hour,omitempty" gorm:"column:stat_time_hour"`
	AC                 string `json:"ac,omitempty" gorm:"column:ac"`
	Age                string `json:"age,omitempty" gorm:"column:age"`
	CountryCode        string `json:"country_code,omitempty" gorm:"column:country_code"`
	InterestCategory   string `json:"interest_category,omitempty" gorm:"column:interest_category"`
	InterestCategoryV2 string `json:"interest_category_v2,omitempty" gorm:"column:interest_category_v2"`
	Gender             string `json:"gender,omitempty" gorm:"column:gender"`
	Language           string `json:"language,omitempty" gorm:"column:language"`
	Placement          string `json:"placement,omitempty" gorm:"column:placement"`
	Platform           string `json:"platform,omitempty" gorm:"column:platform"`
	ContextualTag      string `json:"contextual_tag,omitempty" gorm:"column:contextual_tag"`
	// metrics
	Spend                      string `gorm:"column:spend"`                          //Total Cost
	CashSpend                  string `gorm:"column:cash_spend"`                     //Cost Charged by Cash
	VoucherSpend               string `gorm:"column:voucher_spend"`                  //Cost Charged by Voucher
	CPC                        string `gorm:"column:cpc"`                            //CPC
	CPM                        string `gorm:"column:cpm"`                            //CPM
	Impressions                string `gorm:"column:impressions"`                    //Impressions
	GrossImpressions           string `gorm:"column:gross_impressions"`              //Gross Impressions (Includes Invalid Impressions)
	Clicks                     string `gorm:"column:clicks"`                         //Clicks
	CTR                        string `gorm:"column:ctr"`                            //CTR (%)
	Reach                      string `gorm:"column:reach"`                          //Reach
	CostPer1000Reached         string `gorm:"column:cost_per_1000_reached"`          //Cost per 1,000 people reached
	Conversion                 string `gorm:"column:conversion"`                     //Conversions
	CostPerConversion          string `gorm:"column:cost_per_conversion"`            //CPA
	ConversionRate             string `gorm:"column:conversion_rate"`                //CVR (%)
	RealTimeConversion         string `gorm:"column:real_time_conversion"`           //Real-time Conversions
	RealTimeCostPerConversion  string `gorm:"column:real_time_cost_per_conversion"`  //Real-time CPA
	RealTimeConversionRate     string `gorm:"column:real_time_conversion_rate"`      //Real-time CVR (%)
	Result                     string `gorm:"column:result"`                         //Result
	CostPerResult              string `gorm:"column:cost_per_result"`                //Cost Per Result
	ResultRate                 string `gorm:"column:result_rate"`                    //Result Rate (%)
	RealTimeResult             string `gorm:"column:real_time_result"`               //Real-time Result
	RealTimeCostPerResult      string `gorm:"column:real_time_cost_per_result"`      //Real-time Cost Per Result
	RealTimeResultRate         string `gorm:"column:real_time_result_rate"`          //Real-time Result Rate (%)
	SecondaryGoalResult        string `gorm:"column:secondary_goal_result"`          //Secondary Goal Result
	CostPerSecondaryGoalResult string `gorm:"column:cost_per_secondary_goal_result"` //Cost per Secondary Goal Result
	SecondaryGoalResultRate    string `gorm:"column:secondary_goal_result_rate"`     //Secondary Goal Result Rate (%)
	Frequency                  string `gorm:"column:frequency"`                      //Frequency
	Currency                   string `gorm:"column:currency"`                       //currency
	CampaignName               string `gorm:"column:campaign_name"`                  //Campaign name, Supported at Campaign, Ad Group and Ad level.
	//CampaignID                 string `gorm:"column:campaign_id"` //Campaign ID, Supported at Ad Group and Ad level.
	ObjectiveType        string `gorm:"column:objective_type"`         //Advertising objective, Supported at Campaign, Ad Group and Ad level.
	SplitTest            string `gorm:"column:split_test"`             //Split test status, Supported at Campaign, Ad Group and Ad level.
	CampaignBudget       string `gorm:"column:campaign_budget"`        //Campaign budget, Supported at Campaign, Ad Group and Ad level.
	CampaignDedicateType string `gorm:"column:campaign_dedicate_type"` //Campaign type, iOS14 Dedicated Campaign or regular campaign. Supported at Campaign, Ad Group and Ad level.
	AppPromotionType     string `gorm:"column:app_promotion_type"`     //App promotion type, Supported at Campaign, Ad Group and Ad level. Enum values: APP_INSTALL, APP_RETARGETING. APP_INSTALL and APP_RETARGETING will be returned when objective_type is APP_PROMOTION. Otherwise, UNSET will be returned.
	AdGroupName          string `gorm:"column:adgroup_name"`           //Ad group name, Supported at Ad Group and Ad level.
	//AdGroupID                  string `gorm:"column:adgroup_id"` //Ad group ID, Supported at Ad level.
	PlacementType              string `gorm:"column:placement_type"`                 //Placement type, Supported at Ad Group and Ad level.
	PromotionType              string `gorm:"column:promotion_type"`                 //Promotion type, It can be app, website, or others. Supported at Ad Group and Ad levels in both synchronous and asynchronous reports.
	OptStatus                  string `gorm:"column:opt_status"`                     //Automated creative optimization, Supported at Ad Group and Ad level.
	DpaTargetAudienceType      string `gorm:"column:dpa_target_audience_type"`       //Target audience type for DPA, Supported at Ad Group or Ad levels in both synchronous and asynchronous reports.
	Budget                     string `gorm:"column:budget"`                         //Ad group budget, Supported at Ad Group and Ad level.
	SmartTarget                string `gorm:"column:smart_target"`                   //Optimization goal, Supported at Ad Group and Ad level.
	PricingCategory            string `gorm:"column:pricing_category"`               //Billing Event, Supported at Ad Group and Ad level.
	BidStrategy                string `gorm:"column:bid_strategy"`                   //Bid strategy, Supported at Ad Group and Ad level.
	Bid                        string `gorm:"column:bid"`                            //Bid, Supported at Ad Group and Ad level.
	AeoType                    string `gorm:"column:aeo_type"`                       //App Event Optimization Type, Supported at Ad Group and Ad level. (Already supported at Ad Group level, and will be supported at Ad level)
	AdName                     string `gorm:"column:ad_name"`                        //Ad name, Supported at Ad level.
	AdText                     string `gorm:"column:ad_text"`                        //Ad title, Supported at Ad level.
	CallToAction               string `gorm:"column:call_to_action"`                 //Call to action, Supported at Ad level.
	TikTokAppID                string `gorm:"column:tt_app_id"`                      //TikTok App ID, Supported at Ad Group and Ad level. Returned if the promotion type of one Ad Group is App.
	TikTokAppName              string `gorm:"column:tt_app_name"`                    //TikTok App Name, Supported at Ad Group and Ad level. Returned if the promotion type of one Ad Group is App.
	MobileAppID                string `gorm:"column:mobile_app_id"`                  //Mobile App ID. Examples are, App Store: https://apps.apple.com/us/app/angry-birds/id343200656; Google Play：https://play.google.com/store/apps/details?id=com.rovio.angrybirds. Supported at Ad Group and Ad level. Returned if the promotion type of one Ad Group is App.
	ImageMode                  string `gorm:"column:image_mode"`                     //Format, Supported at Campaign, Ad Group and Ad level.
	OnsiteDownloadStart        string `gorm:"column:onsite_download_start"`          //Total App Store Click (Onsite), The number of app store click events attributed to your TikTok ads.
	CostPerOnsiteDownloadStart string `gorm:"column:cost_per_onsite_download_start"` //Cost per App Store Click (Onsite), The average cost of each app store click event attributed to your TikTok ads.
	OnsiteDownloadStartRate    string `gorm:"column:onsite_download_start_rate"`     //App Store Click Rate (Onsite) (%), The percentage of app store click events out of all click events attributed to your TikTok ads.
}

func (Report) TableName() string {
	return "_tool_tiktokAds_report"
}
