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

package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*addIssueImpacts)(nil)

type issueImpacts20241010 struct {
	ConnectionId    uint64 `gorm:"primaryKey"`
	IssueKey        string `gorm:"primaryKey;type:varchar(100)"`
	SoftwareQuality string `gorm:"primaryKey;type:varchar(255)"`
	Severity        string `gorm:"type:varchar(100)"`
	archived.NoPKModel
}

func (issueImpacts20241010) TableName() string {
	return "_tool_sonarqube_issue_impacts"
}

type addIssueImpacts struct {
}

func (script *addIssueImpacts) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes, &issueImpacts20241010{})
}

func (*addIssueImpacts) Version() uint64 {
	return 20241010162943
}

func (*addIssueImpacts) Name() string {
	return "add issue_impacts table for sonarcloud"
}
