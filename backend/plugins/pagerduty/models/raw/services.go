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

type Service struct {
	Id                      string    `json:"id"`
	Summary                 string    `json:"summary"`
	Type                    string    `json:"type"`
	Self                    string    `json:"self"`
	HtmlUrl                 string    `json:"html_url"`
	Name                    string    `json:"name"`
	AutoResolveTimeout      int       `json:"auto_resolve_timeout"`
	AcknowledgementTimeout  int       `json:"acknowledgement_timeout"`
	CreatedAt               time.Time `json:"created_at"`
	Status                  string    `json:"status"`
	AlertCreation           string    `json:"alert_creation"`
	AlertGroupingParameters struct {
		Type string `json:"type"`
	} `json:"alert_grouping_parameters"`
	Integrations []struct {
		Id      string `json:"id"`
		Type    string `json:"type"`
		Summary string `json:"summary"`
		Self    string `json:"self"`
		HtmlUrl string `json:"html_url"`
	} `json:"integrations"`
	EscalationPolicy struct {
		Id      string `json:"id"`
		Type    string `json:"type"`
		Summary string `json:"summary"`
		Self    string `json:"self"`
		HtmlUrl string `json:"html_url"`
	} `json:"escalation_policy"`
	Teams []struct {
		Id      string `json:"id"`
		Type    string `json:"type"`
		Summary string `json:"summary"`
		Self    string `json:"self"`
		HtmlUrl string `json:"html_url"`
	} `json:"teams"`
	IncidentUrgencyRule struct {
		Type               string `json:"type"`
		DuringSupportHours struct {
			Type    string `json:"type"`
			Urgency string `json:"urgency"`
		} `json:"during_support_hours"`
		OutsideSupportHours struct {
			Type    string `json:"type"`
			Urgency string `json:"urgency"`
		} `json:"outside_support_hours"`
	} `json:"incident_urgency_rule"`
	SupportHours struct {
		Type       string `json:"type"`
		TimeZone   string `json:"time_zone"`
		StartTime  string `json:"start_time"`
		EndTime    string `json:"end_time"`
		DaysOfWeek []int  `json:"days_of_week"`
	} `json:"support_hours"`
	ScheduledActions []struct {
		Type string `json:"type"`
		At   struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"at"`
		ToUrgency string `json:"to_urgency"`
	} `json:"scheduled_actions"`
	AutoPauseNotificationsParameters struct {
		Enabled bool `json:"enabled"`
		Timeout int  `json:"timeout"`
	} `json:"auto_pause_notifications_parameters"`
}
