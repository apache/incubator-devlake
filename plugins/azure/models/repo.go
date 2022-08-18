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
	"github.com/apache/incubator-devlake/models/common"
)

type AzureRepo struct {
	ConnectionId  uint64 `gorm:"primaryKey"`
	AzureId       string `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name          string `gorm:"type:varchar(255)" json:"name"`
	Url           string `gorm:"type:varchar(255)" json:"url"`
	ProjectId     string `gorm:"type:varchar(255);index"`
	DefaultBranch string `json:"defaultBranch"`
	Size          int    `json:"size"`
	RemoteURL     string `json:"remoteUrl"`
	SshUrl        string `gorm:"type:varchar(255)" json:"sshUrl"`
	WebUrl        string `gorm:"type:varchar(255)" json:"webUrl"`
	IsDisabled    bool   `json:"isDisabled"`
	common.NoPKModel
}

func (AzureRepo) TableName() string {
	return "_tool_azure_repos"
}
