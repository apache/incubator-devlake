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
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
)

type SlackChannel struct {
	archived.NoPKModel `json:"-"`
	ConnectionId       uint64 `gorm:"primaryKey"`
	Id                 string `json:"id" gorm:"primaryKey"`
	Name               string `json:"name"`
	IsChannel          bool   `json:"is_channel"`
	IsGroup            bool   `json:"is_group"`
	IsIm               bool   `json:"is_im"`
	IsMpim             bool   `json:"is_mpim"`
	IsPrivate          bool   `json:"is_private"`
	Created            int    `json:"created"`
	IsArchived         bool   `json:"is_archived"`
	IsGeneral          bool   `json:"is_general"`
	Unlinked           int    `json:"unlinked"`
	NameNormalized     string `json:"name_normalized"`
	IsShared           bool   `json:"is_shared"`
	IsOrgShared        bool   `json:"is_org_shared"`
	IsPendingExtShared bool   `json:"is_pending_ext_shared"`
	ContextTeamId      string `json:"context_team_id"`
	Updated            int64  `json:"updated"`
	Creator            string `json:"creator"`
	IsExtShared        bool   `json:"is_ext_shared"`
	IsMember           bool   `json:"is_member"`
	NumMembers         int    `json:"num_members"`
}

func (SlackChannel) TableName() string {
	return "_tool_slack_channels"
}
