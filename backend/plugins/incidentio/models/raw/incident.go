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

import (
	"time"
)

type Incident struct {
	Id                      string                   `json:"id"`
	Reference               string                   `json:"reference"`
	Name                    string                   `json:"name"`
	Summary                 *string                  `json:"summary"`
	Permalink               *string                  `json:"permalink"`
	Mode                    string                   `json:"mode"`
	CreatedAt               time.Time                `json:"created_at"`
	UpdatedAt               time.Time                `json:"updated_at"`
	IncidentStatus          *IncidentStatus          `json:"incident_status"`
	Severity                *Severity                `json:"severity"`
	IncidentType            *IncidentTypeRef         `json:"incident_type"`
	Creator                 *Creator                 `json:"creator"`
	IncidentTimestampValues []IncidentTimestampValue `json:"incident_timestamp_values"`
}

type IncidentStatus struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

type Severity struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Rank int64  `json:"rank"`
}

type IncidentTypeRef struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Creator struct {
	User *CreatorUser `json:"user"`
}

type CreatorUser struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type IncidentTimestampValue struct {
	IncidentTimestamp IncidentTimestamp `json:"incident_timestamp"`
	Value             *TimestampValue   `json:"value"`
}

type IncidentTimestamp struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type TimestampValue struct {
	Value time.Time `json:"value"`
}
