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
)

type AzuredevopsUser struct {
	common.NoPKModel

	ConnectionId  uint64 `gorm:"primaryKey"`
	AzuredevopsId string `gorm:"primaryKey"`
	Origin        string
	Descriptor    string
	PrincipalName string
	MailAddress   string
	DisplayName   string
	Url           string
}

func (AzuredevopsUser) TableName() string {
	return "_tool_azuredevops_go_users"
}

type AzuredevopsApiUser struct {
	SubjectKind    string `json:"subjectKind"`
	MetaType       string `json:"metaType"`
	DirectoryAlias string `json:"directoryAlias"`
	Domain         string `json:"domain"`
	PrincipalName  string `json:"principalName"`
	MailAddress    string `json:"mailAddress"`
	Origin         string `json:"origin"`
	OriginId       string `json:"originId"`
	DisplayName    string `json:"displayName"`
	Url            string `json:"url"`
	Descriptor     string `json:"descriptor"`
}

func (u AzuredevopsApiUser) ToModel() AzuredevopsUser {
	res := AzuredevopsUser{
		AzuredevopsId: u.OriginId,
		Origin:        u.Origin,
		Descriptor:    u.Descriptor,
		PrincipalName: u.PrincipalName,
		MailAddress:   u.MailAddress,
		DisplayName:   u.DisplayName,
		Url:           u.Url,
	}

	return res
}
