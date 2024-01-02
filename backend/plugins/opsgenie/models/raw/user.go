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

package raw

import "time"

type User struct {
	Blocked  *bool   `json:"blocked"`
	Verified *bool   `json:"verified"`
	Id       *string `json:"id"`
	Username *string `json:"username"`
	FullName *string `json:"fullName"`
	Role     *struct {
		Id   *string `json:"id"`
		Name *string `json:"name"`
	} `json:"role"`
	TimeZone    *string `json:"timeZone"`
	Locale      *string `json:"locale"`
	UserAddress *struct {
		Country *string `json:"country"`
		State   *string `json:"state"`
		City    *string `json:"city"`
		Line    *string `json:"line"`
		ZipCode *string `json:"zipCode"`
	} `json:"userAddress"`
	CreatedAt *time.Time `json:"createdAt"`
}
