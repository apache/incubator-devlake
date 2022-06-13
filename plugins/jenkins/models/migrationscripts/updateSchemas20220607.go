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

	"github.com/apache/incubator-devlake/plugins/helper"
	"gorm.io/gorm"
)

type JenkinsConnection20220607 struct {
	helper.RestConnection
	helper.BasicAuth
}

func (JenkinsConnection20220607) TableName() string {
	return "_tool_jenkins_connections"
}

type UpdateSchemas20220607 struct{}

func (*UpdateSchemas20220607) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().CreateTable(&JenkinsConnection20220607{})
	if err != nil {
		return err
	}
	return nil
}

func (*UpdateSchemas20220607) Version() uint64 {
	return 20220607154646
}

func (*UpdateSchemas20220607) Name() string {
	return "add table _tool_jenkins_connections"
}
