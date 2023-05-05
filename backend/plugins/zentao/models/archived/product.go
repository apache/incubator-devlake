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
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type ZentaoProduct struct {
	archived.NoPKModel `json:"-"`
	ConnectionId       uint64 `json:"connectionid" mapstructure:"connectionid" gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Id                 int64  `json:"id" mapstructure:"id" gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Program            int    `json:"program" mapstructure:"program"`
	Name               string `json:"name" mapstructure:"name"`
	Code               string `json:"code" mapstructure:"code"`
	Bind               string `json:"bind" mapstructure:"bind"`
	Line               int    `json:"line" mapstructure:"line"`
	Type               string `json:"type" mapstructure:"type"`
	ProductType        string `json:"productType" mapstructure:"productType"`
	Status             string `json:"status" mapstructure:"status"`
	SubStatus          string `json:"subStatus" mapstructure:"subStatus"`
	Description        string `json:"desc" mapstructure:"desc"`
	POId               int64
	QDId               int64
	RDId               int64
	Acl                string `json:"acl" mapstructure:"acl"`
	Reviewer           string `json:"reviewer" mapstructure:"reviewer"`
	CreatedById        int64
	CreatedDate        *helper.Iso8601Time `json:"createdDate" mapstructure:"createdDate"`
	CreatedVersion     string              `json:"createdVersion" mapstructure:"createdVersion"`
	OrderIn            int                 `json:"order" mapstructure:"order"`
	Deleted            string              `json:"deleted" mapstructure:"deleted"`
	Plans              int                 `json:"plans" mapstructure:"plans"`
	Releases           int                 `json:"releases" mapstructure:"releases"`
	Builds             int                 `json:"builds" mapstructure:"builds"`
	Cases              int                 `json:"cases" mapstructure:"cases"`
	Projects           int                 `json:"projects" mapstructure:"projects"`
	Executions         int                 `json:"executions" mapstructure:"executions"`
	Bugs               int                 `json:"bugs" mapstructure:"bugs"`
	Docs               int                 `json:"docs" mapstructure:"docs"`
	Progress           float64             `json:"progress" mapstructure:"progress"`
	CaseReview         bool                `json:"caseReview" mapstructure:"caseReview"`
}

func (ZentaoProduct) TableName() string {
	return "_tool_zentao_products"
}
