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

type AzuredevopsUser struct {
	archived.NoPKModel

	ConnectionId  uint64 `gorm:"primaryKey"`
	AzuredevopsId string `gorm:"primaryKey"`
	Origin        string `gorm:"type:varchar(255)"`
	Descriptor    string `gorm:"type:varchar(255)"`
	PrincipalName string `gorm:"type:varchar(255)"`
	MailAddress   string `gorm:"type:varchar(255)"`
	DisplayName   string `gorm:"type:varchar(255)"`
	Url           string `gorm:"type:varchar(255)"`
}

func (AzuredevopsUser) TableName() string {
	return "_tool_azuredevops_go_users"
}
