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
	"github.com/apache/incubator-devlake/plugins/helper"
)

type ZentaoProductRes struct {
	ID             int64               `json:"id"`
	Program        int                 `json:"program"`
	Name           string              `json:"name"`
	Code           string              `json:"code"`
	Bind           string              `json:"bind"`
	Line           int                 `json:"line"`
	Type           string              `json:"type"`
	Status         string              `json:"status"`
	SubStatus      string              `json:"subStatus"`
	Description    string              `json:"desc"`
	PO             *ZentaoAccount      `json:"PO"`
	QD             *ZentaoAccount      `json:"QD"`
	RD             *ZentaoAccount      `json:"RD"`
	Feedback       interface{}         `json:"feedback"`
	Acl            string              `json:"acl"`
	Whitelist      []interface{}       `json:"whitelist"`
	Reviewer       string              `json:"reviewer"`
	CreatedBy      *ZentaoAccount      `json:"createdBy"`
	CreatedDate    *helper.Iso8601Time `json:"createdDate"`
	CreatedVersion string              `json:"createdVersion"`
	OrderIn        int                 `json:"order"`
	Vision         string              `json:"vision"`
	Deleted        string              `json:"deleted"`
	Stories        struct {
		Active    int `json:"active"`
		Reviewing int `json:"reviewing"`
		int       `json:""`
		Draft     int `json:"draft"`
		Closed    int `json:"closed"`
		Changing  int `json:"changing"`
	} `json:"stories"`
	Plans      int     `json:"plans"`
	Releases   int     `json:"releases"`
	Builds     int     `json:"builds"`
	Cases      int     `json:"cases"`
	Projects   int     `json:"projects"`
	Executions int     `json:"executions"`
	Bugs       int     `json:"bugs"`
	Docs       int     `json:"docs"`
	Progress   float64 `json:"progress"`
	CaseReview bool    `json:"caseReview"`
}

type ZentaoProduct struct {
	ConnectionId   uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Id             int64  `json:"id" gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Program        int    `json:"program"`
	Name           string `json:"name"`
	Code           string `json:"code"`
	Bind           string `json:"bind"`
	Line           int    `json:"line"`
	Type           string `json:"type"`
	Status         string `json:"status"`
	SubStatus      string `json:"subStatus"`
	Description    string `json:"desc"`
	POId           int64
	QDId           int64
	RDId           int64
	Acl            string `json:"acl"`
	Reviewer       string `json:"reviewer"`
	CreatedById    int64
	CreatedDate    *helper.Iso8601Time `json:"createdDate"`
	CreatedVersion string              `json:"createdVersion"`
	OrderIn        int                 `json:"order"`
	Deleted        string              `json:"deleted"`
	Plans          int                 `json:"plans"`
	Releases       int                 `json:"releases"`
	Builds         int                 `json:"builds"`
	Cases          int                 `json:"cases"`
	Projects       int                 `json:"projects"`
	Executions     int                 `json:"executions"`
	Bugs           int                 `json:"bugs"`
	Docs           int                 `json:"docs"`
	Progress       float64             `json:"progress"`
	CaseReview     bool                `json:"caseReview"`
	common.NoPKModel
}

func (ZentaoProduct) TableName() string {
	return "_tool_zentao_products"
}
