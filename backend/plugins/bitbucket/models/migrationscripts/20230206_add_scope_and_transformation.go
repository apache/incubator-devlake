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
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models/migrationscripts/archived"
)

type BitbucketRepo20230206 struct {
	TransformationRuleId uint64 `json:"transformationRuleId,omitempty" mapstructure:"transformationRuleId,omitempty"`
	CloneUrl             string `json:"cloneUrl" gorm:"type:varchar(255)" mapstructure:"cloneUrl,omitempty"`
	Owner                string `json:"owner" mapstructure:"owner,omitempty"`
}

func (BitbucketRepo20230206) TableName() string {
	return "_tool_bitbucket_repos"
}

type BitbucketIssue20230206 struct {
	StdState string `gorm:"type:varchar(255)"`
}

func (BitbucketIssue20230206) TableName() string {
	return "_tool_bitbucket_issues"
}

type addScope20230206 struct{}

func (*addScope20230206) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	err := db.RenameColumn("_tool_bitbucket_repos", "owner_id", "owner")
	if err != nil {
		return err
	}

	return migrationhelper.AutoMigrateTables(
		basicRes,
		&BitbucketRepo20230206{},
		&BitbucketIssue20230206{},
		&archived.BitbucketTransformationRule{},
	)
}

func (*addScope20230206) Version() uint64 {
	return 20230206000008
}

func (*addScope20230206) Name() string {
	return "add scope and table _tool_bitbucket_transformation_rules"
}
