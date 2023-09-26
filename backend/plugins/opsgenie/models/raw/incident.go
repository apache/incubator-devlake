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

type Incident struct {
	Id               *string    `json:"id"`
	Description      *string    `json:"description"`
	ImpactedServices *[]string  `json:"impactedServices"`
	TinyId           *string    `json:"tinyId"`
	Message          *string    `json:"message"`
	Status           *string    `json:"status"`
	Tags             []any      `json:"tags"`
	CreatedAt        *time.Time `json:"createdAt"`
	UpdatedAt        *time.Time `json:"updatedAt"`
	Priority         *string    `json:"priority"`
	OwnerTeam        *string    `json:"ownerTeam"`
	Responders       *[]struct {
		Type *string `json:"type"`
		Id   *string `json:"id"`
	} `json:"responders"`
	ExtraProperties struct {
	} `json:"extraProperties"`
	Links struct {
		Web *string `json:"web"`
		Api *string `json:"api"`
	} `json:"links"`
	Actions *[]any `json:"actions"`
}
