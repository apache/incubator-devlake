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
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type ZentaoProductRes struct {
	ID             int64               `json:"id" mapstructure:"id"`
	Program        int                 `json:"program" mapstructure:"program"`
	Name           string              `json:"name" mapstructure:"name"`
	Code           string              `json:"code" mapstructure:"code"`
	Bind           string              `json:"bind" mapstructure:"bind"`
	Line           int                 `json:"line" mapstructure:"line"`
	Type           string              `json:"type" mapstructure:"type"`
	Status         string              `json:"status" mapstructure:"status"`
	SubStatus      string              `json:"subStatus" mapstructure:"subStatus"`
	Description    string              `json:"desc" mapstructure:"desc"`
	PO             *ZentaoAccount      `json:"po" mapstructure:"po"`
	QD             *ZentaoAccount      `json:"qd" mapstructure:"qd"`
	RD             *ZentaoAccount      `json:"rd" mapstructure:"rd"`
	Feedback       interface{}         `json:"feedback" mapstructure:"feedback"`
	Acl            string              `json:"acl" mapstructure:"acl"`
	Whitelist      []interface{}       `json:"whitelist" mapstructure:"whitelist"`
	Reviewer       string              `json:"reviewer" mapstructure:"reviewer"`
	CreatedBy      *ZentaoAccount      `json:"createdBy" mapstructure:"createdBy"`
	CreatedDate    *helper.Iso8601Time `json:"createdDate" mapstructure:"createdDate"`
	CreatedVersion string              `json:"createdVersion" mapstructure:"createdVersion"`
	OrderIn        int                 `json:"order" mapstructure:"order"`
	Vision         string              `json:"vision" mapstructure:"vision"`
	Deleted        string              `json:"deleted" mapstructure:"deleted"`
	Stories        struct {
		Active    int `json:"active" mapstructure:"active"`
		Reviewing int `json:"reviewing" mapstructure:"reviewing"`
		int       `json:""`
		Draft     int `json:"draft" mapstructure:"draft"`
		Closed    int `json:"closed" mapstructure:"closed"`
		Changing  int `json:"changing" mapstructure:"changing"`
	} `json:"stories"`
	Plans      int     `json:"plans" mapstructure:"plans"`
	Releases   int     `json:"releases" mapstructure:"releases"`
	Builds     int     `json:"builds" mapstructure:"builds"`
	Cases      int     `json:"cases" mapstructure:"cases"`
	Projects   int     `json:"projects" mapstructure:"projects"`
	Executions int     `json:"executions" mapstructure:"executions"`
	Bugs       int     `json:"bugs" mapstructure:"bugs"`
	Docs       int     `json:"docs" mapstructure:"docs"`
	Progress   float64 `json:"progress" mapstructure:"progress"`
	CaseReview bool    `json:"caseReview" mapstructure:"caseReview"`
}

func getAccountId(account *ZentaoAccount) int64 {
	if account != nil {
		return account.ID
	}
	return 0
}

type ZentaoProduct struct {
	common.NoPKModel `json:"-"`
	ConnectionId     uint64 `json:"connectionId" mapstructure:"connectionId" gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Id               int64  `json:"id" mapstructure:"id" gorm:"primaryKey;type:BIGINT  NOT NULL;autoIncrement:false"`
	Program          int    `json:"program" mapstructure:"program"`
	Name             string `json:"name" mapstructure:"name"`
	Code             string `json:"code" mapstructure:"code"`
	Bind             string `json:"bind" mapstructure:"bind"`
	Line             int    `json:"line" mapstructure:"line"`
	Type             string `json:"type" mapstructure:"type"`
	ProductType      string `json:"productType" mapstructure:"productType"`
	Status           string `json:"status" mapstructure:"status"`
	SubStatus        string `json:"subStatus" mapstructure:"subStatus"`
	Description      string `json:"desc" mapstructure:"desc"`
	POId             int64
	QDId             int64
	RDId             int64
	Acl              string `json:"acl" mapstructure:"acl"`
	Reviewer         string `json:"reviewer" mapstructure:"reviewer"`
	CreatedById      int64
	CreatedDate      *helper.Iso8601Time `json:"createdDate" mapstructure:"createdDate"`
	CreatedVersion   string              `json:"createdVersion" mapstructure:"createdVersion"`
	OrderIn          int                 `json:"order" mapstructure:"order"`
	Deleted          string              `json:"deleted" mapstructure:"deleted"`
	Plans            int                 `json:"plans" mapstructure:"plans"`
	Releases         int                 `json:"releases" mapstructure:"releases"`
	Builds           int                 `json:"builds" mapstructure:"builds"`
	Cases            int                 `json:"cases" mapstructure:"cases"`
	Projects         int                 `json:"projects" mapstructure:"projects"`
	Executions       int                 `json:"executions" mapstructure:"executions"`
	Bugs             int                 `json:"bugs" mapstructure:"bugs"`
	Docs             int                 `json:"docs" mapstructure:"docs"`
	Progress         float64             `json:"progress" mapstructure:"progress"`
	CaseReview       bool                `json:"caseReview" mapstructure:"caseReview"`
	ScopeConfigId    uint64              `json:"scopeConfigId,omitempty" mapstructure:"scopeConfigId"`
}

func (ZentaoProduct) TableName() string {
	return "_tool_zentao_products"
}
