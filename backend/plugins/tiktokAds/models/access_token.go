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

type ApiAccessTokenRequest struct {
	AppId    string `json:"app_id"`
	AuthCode string `json:"auth_code"`
	Secret   string `json:"secret"`
}

type ApiAccessTokenResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      `json:data`
	RequestId string `json:request_id`
}

type Data struct {
	AccessToken   string   `json:access_token`
	Scope         []int    `json:scope`
	AdvertiserIds []string `json:advertiser_ids`
}
