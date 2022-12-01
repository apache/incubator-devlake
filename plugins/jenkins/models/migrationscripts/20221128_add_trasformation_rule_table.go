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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/jenkins/models/migrationscripts/archived"
)

type jenkinsJob20221128 struct {
	TransformationRuleId uint64
}

func (jenkinsJob20221128) TableName() string {
	return "_tool_jenkins_jobs"
}

type addTransformationRule20221128 struct{}

func (*addTransformationRule20221128) Up(basicRes core.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes, &jenkinsJob20221128{}, &archived.JenkinsTransformationRule{})
}

func (*addTransformationRule20221128) Version() uint64 {
	return 20221128113500
}

func (*addTransformationRule20221128) Name() string {
	return "add table _tool_jenkins_transformation_rules, add transformation_rule_id to _tool_jenkins_jobs"
}
