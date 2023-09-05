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
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type CircleciPipeline struct {
	ConnectionId uint64 `gorm:"primaryKey;type:BIGINT"`
	Id           string `gorm:"primaryKey;type:varchar(100)" json:"id"`
	ProjectSlug  string `gorm:"type:varchar(255)" json:"project_slug"`
	Errors       []struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `gorm:"serializer:json;type:text" json:"errors"`
	UpdatedAt         *api.Iso8601Time `json:"updated_at"`
	Number            int64            `json:"number"`
	TriggerParameters any              `gorm:"serializer:json;type:text" json:"trigger_parameters"`
	State             string           `gorm:"type:varchar(100)" json:"state"`
	CreatedAt         *api.Iso8601Time `json:"created_at"`
	Trigger           CircleciTrigger  `gorm:"serializer:json;type:text" json:"trigger"`
	Vcs               CircleciVcs      `gorm:"serializer:json;type:text" json:"vcs"`
	common.NoPKModel
}

type CircleciTrigger struct {
	Type       string `json:"type"`
	ReceivedAt string `json:"received_at"`
	Actor      struct {
		Login     string `json:"login"`
		AvatarUrl string `json:"avatar_url"`
	} `json:"actor"`
}

type CircleciVcs struct {
	ProviderName        string `json:"provider_name"`
	TargetRepositoryUrl string `json:"target_repository_url"`
	Branch              string `json:"branch"`
	ReviewId            string `json:"review_id"`
	ReviewUrl           string `json:"review_url"`
	Revision            string `json:"revision"`
	Tag                 string `json:"tag"`
	OriginRepositoryUrl string `json:"origin_repository_url"`
	Commit              struct {
		Subject string `json:"subject"`
		Body    string `json:"body"`
	} `json:"commit"`
}

func (CircleciPipeline) TableName() string {
	return "_tool_circleci_pipelines"
}
