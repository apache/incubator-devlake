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
	"context"
	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/plugins/helper"
	"gorm.io/gorm"
)

type GithubConnection20220608 struct {
	archived.Model
	Name      string `gorm:"type:varchar(100);uniqueIndex" json:"name" validate:"required"`
	Endpoint  string `mapstructure:"endpoint" env:"GITHUB_ENDPOINT" validate:"required" json:"endpoint"`
	Proxy     string `mapstructure:"proxy" env:"GITHUB_PROXY" json:"proxy"`
	RateLimit int    `comment:"api request rate limit per hour" json:"rateLimit"`
	Auth      string `mapstructure:"auth" validate:"required" env:"GITHUB_AUTH" json:"auth"`

	Config20220608 `mapstructure:",squash"`
}

type Config20220608 struct {
	PrType               string `mapstructure:"prType" env:"GITHUB_PR_TYPE" json:"prType"`
	PrComponent          string `mapstructure:"prComponent" env:"GITHUB_PR_COMPONENT" json:"prComponent"`
	IssueSeverity        string `mapstructure:"issueSeverity" env:"GITHUB_ISSUE_SEVERITY" json:"issueSeverity"`
	IssuePriority        string `mapstructure:"issuePriority" env:"GITHUB_ISSUE_PRIORITY" json:"issuePriority"`
	IssueComponent       string `mapstructure:"issueComponent" env:"GITHUB_ISSUE_COMPONENT" json:"issueComponent"`
	IssueTypeBug         string `mapstructure:"issueTypeBug" env:"GITHUB_ISSUE_TYPE_BUG" json:"issueTypeBug"`
	IssueTypeIncident    string `mapstructure:"issueTypeIncident" env:"GITHUB_ISSUE_TYPE_INCIDENT" json:"issueTypeIncident"`
	IssueTypeRequirement string `mapstructure:"issueTypeRequirement" env:"GITHUB_ISSUE_TYPE_REQUIREMENT" json:"issueTypeRequirement"`
}

func (GithubConnection20220608) TableName() string {
	return "_tool_github_connections"
}

type UpdateSchemas20220608 struct{}

func (*UpdateSchemas20220608) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().CreateTable(GithubConnection20220608{})
	if err != nil {
		return err
	}
	v := config.GetConfig()
	connection := &GithubConnection20220608{}
	err = helper.EncodeStruct(v, connection, "env")
	connection.Name = `GitHub`
	if err != nil {
		return err
	}
	// update from .env and save to db
	if connection.Endpoint != `` && connection.Auth != `` {
		db.Create(connection)
	}
	return nil
}

func (*UpdateSchemas20220608) Version() uint64 {
	return 20220608000003
}

func (*UpdateSchemas20220608) Name() string {
	return "Add connection for github"
}
