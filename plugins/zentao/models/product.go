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
	"time"
)

type ZentaoProductRes struct {
	ID          uint64 `json:"id"`
	Program     int    `json:"program"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Bind        string `json:"bind"`
	Line        int    `json:"line"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	SubStatus   string `json:"subStatus"`
	Description string `json:"desc"`
	PO          struct {
		ID       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"PO"`
	QD struct {
		ID       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"QD"`
	RD struct {
		ID       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"RD"`
	Feedback  interface{}   `json:"feedback"`
	Acl       string        `json:"acl"`
	Whitelist []interface{} `json:"whitelist"`
	Reviewer  string        `json:"reviewer"`
	CreatedBy struct {
		ID       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"createdBy"`
	CreatedDate    *time.Time `json:"createdDate"`
	CreatedVersion string     `json:"createdVersion"`
	OrderIn        int        `json:"order"`
	Vision         string     `json:"vision"`
	Deleted        string     `json:"deleted"`
	Stories        struct {
		Active    int `json:"active"`
		Reviewing int `json:"reviewing"`
		int       `json:""`
		Draft     int `json:"draft"`
		Closed    int `json:"closed"`
		Changing  int `json:"changing"`
	} `json:"stories"`
	Plans      int  `json:"plans"`
	Releases   int  `json:"releases"`
	Builds     int  `json:"builds"`
	Cases      int  `json:"cases"`
	Projects   int  `json:"projects"`
	Executions int  `json:"executions"`
	Bugs       int  `json:"bugs"`
	Docs       int  `json:"docs"`
	Progress   int  `json:"progress"`
	CaseReview bool `json:"caseReview"`
}

type ZentaoProduct struct {
	ConnectionId   uint64 `gorm:"primaryKey"`
	Id             uint64 `json:"id" gorm:"primaryKey"`
	Program        int    `json:"program"`
	Name           string `json:"name"`
	Code           string `json:"code"`
	Bind           string `json:"bind"`
	Line           int    `json:"line"`
	Type           string `json:"type"`
	Status         string `json:"status"`
	SubStatus      string `json:"subStatus"`
	Description    string `json:"desc"`
	POId           uint64
	QDId           uint64
	RDId           uint64
	Acl            string `json:"acl"`
	Reviewer       string `json:"reviewer"`
	CreatedById    uint64
	CreatedDate    *time.Time `json:"createdDate"`
	CreatedVersion string     `json:"createdVersion"`
	OrderIn        int        `json:"order"`
	Deleted        string     `json:"deleted"`
	Plans          int        `json:"plans"`
	Releases       int        `json:"releases"`
	Builds         int        `json:"builds"`
	Cases          int        `json:"cases"`
	Projects       int        `json:"projects"`
	Executions     int        `json:"executions"`
	Bugs           int        `json:"bugs"`
	Docs           int        `json:"docs"`
	Progress       int        `json:"progress"`
	CaseReview     bool       `json:"caseReview"`
	common.NoPKModel
}

func (ZentaoProduct) TableName() string {
	return "_tool_zentao_products"
}
