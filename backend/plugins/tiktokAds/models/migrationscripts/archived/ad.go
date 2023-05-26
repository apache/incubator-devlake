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

type AdResponse struct {
	Code int `json:"code"`
	Data struct {
		AdIDs     []string `json:"ad_ids"`
		Creatives []Ad     `json:"creatives"`
		NeedAudit bool     `json:"need_audit"`
	} `json:"data"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

type Ad struct {
	Id                        string `gorm:"column:ad_id;primaryKey" json:"ad_id"`
	CreativeId                uint64 `gorm:"column:creative_id;primaryKey" json:"creative_id"`
	AdFormat                  string `gorm:"column:ad_format" json:"ad_format"`
	AdName                    string `gorm:"column:ad_name" json:"ad_name"`
	AdText                    string `gorm:"column:ad_text" json:"ad_text"`
	AdTexts                   []byte `gorm:"column:ad_texts" json:"ad_texts"`
	AdGroupID                 string `gorm:"column:adgroup_id" json:"adgroup_id"`
	AdGroupName               string `gorm:"column:adgroup_name" json:"adgroup_name"`
	AdvertiserID              string `gorm:"column:advertiser_id" json:"advertiser_id"`
	AppName                   string `gorm:"column:app_name" json:"app_name"`
	AvatarIconWebURI          string `gorm:"column:avatar_icon_web_uri" json:"avatar_icon_web_uri"`
	CallToAction              string `gorm:"column:call_to_action" json:"call_to_action"`
	CallToActionID            []byte `gorm:"column:call_to_action_id" json:"call_to_action_id"`
	CampaignID                string `gorm:"column:campaign_id" json:"campaign_id"`
	CampaignName              string `gorm:"column:campaign_name" json:"campaign_name"`
	CardID                    string `gorm:"column:card_id" json:"card_id"`
	CatalogID                 string `gorm:"column:catalog_id" json:"catalog_id"`
	ClickTrackingURL          string `gorm:"column:click_tracking_url" json:"click_tracking_url"`
	CreateTime                string `gorm:"column:create_time" json:"create_time"`
	CreativeType              []byte `gorm:"column:creative_type" json:"creative_type"`
	DisplayName               string `gorm:"column:display_name" json:"display_name"`
	DynamicDestination        string `gorm:"column:dynamic_destination" json:"dynamic_destination"`
	DynamicFormat             string `gorm:"column:dynamic_format" json:"dynamic_format"`
	ExternalAction            string `gorm:"column:external_action" json:"external_action"`
	FallbackType              string `gorm:"column:fallback_type" json:"fallback_type"`
	IdentityID                string `gorm:"column:identity_id" json:"identity_id"`
	IdentityType              string `gorm:"column:identity_type" json:"identity_type"`
	ImageIDs                  []byte `gorm:"column:image_ids" json:"image_ids"`
	ImpressionTrackingURL     string `gorm:"column:impression_tracking_url" json:"impression_tracking_url"`
	IsACO                     bool   `gorm:"column:is_aco" json:"is_aco"`
	IsCreativeAuthorized      bool   `gorm:"column:is_creative_authorized" json:"is_creative_authorized"`
	IsNewStructure            bool   `gorm:"column:is_new_structure" json:"is_new_structure"`
	LandingPageURL            string `gorm:"column:landing_page_url" json:"landing_page_url"`
	LandingPageURLs           []byte `gorm:"column:landing_page_urls" json:"landing_page_urls"`
	ModifyTime                string `gorm:"column:modify_time" json:"modify_time"`
	Deeplink                  string `gorm:"column:deeplink" json:"deeplink"`
	DeeplinkType              string `gorm:"column:deeplink_type" json:"deeplink_type"`
	OperationStatus           string `gorm:"column:operation_status" json:"operation_status"`
	PageID                    []byte `gorm:"column:page_id" json:"page_id"`
	PlayableURL               string `gorm:"column:playable_url" json:"playable_url"`
	ProductSetID              string `gorm:"column:product_set_id" json:"product_set_id"`
	ProductSpecificType       string `gorm:"column:product_specific_type" json:"product_specific_type"`
	ProfileImage              string `gorm:"column:profile_image" json:"profile_image"`
	ShoppingAdsFallbackType   string `gorm:"column:shopping_ads_fallback_type" json:"shopping_ads_fallback_type"`
	ShoppingAdsDeeplinkType   string `gorm:"column:shopping_ads_deeplink_type" json:"shopping_ads_deeplink_type"`
	ShoppingAdsVideoPackageID string `gorm:"column:shopping_ads_video_package_id" json:"shopping_ads_video_package_id"`
	VastMoatEnabled           bool   `gorm:"column:vast_moat_enabled" json:"vast_moat_enabled"`
	VerticalVideoStrategy     string `gorm:"column:vertical_video_strategy" json:"vertical_video_strategy"`
	VideoID                   []byte `gorm:"column:video_id" json:"video_id"`
}

func (Ad) TableName() string {
	return "_tool_tiktokAds_ad"
}
