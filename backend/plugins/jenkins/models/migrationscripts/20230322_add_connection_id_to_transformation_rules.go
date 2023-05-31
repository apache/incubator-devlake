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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type addConnectionIdToTransformationRule struct{}

type job20230322 struct {
	ConnectionId         uint64 `gorm:"primaryKey" mapstructure:"connectionId,omitempty" validate:"required" json:"connectionId"`
	FullName             string `gorm:"primaryKey;type:varchar(255)" mapstructure:"jobFullName" validate:"required" json:"jobFullName"` // "path1/path2/job name"
	TransformationRuleId uint64 `mapstructure:"transformationRuleId,omitempty" json:"transformationRuleId,omitempty"`
}

func (job20230322) TableName() string {
	return "_tool_jenkins_jobs"
}

type transformationRule20230322 struct {
	ConnectionId uint64
}

func (transformationRule20230322) TableName() string {
	return "_tool_jenkins_transformation_rules"
}

func (u *addConnectionIdToTransformationRule) Up(baseRes context.BasicRes) errors.Error {
	err := migrationhelper.AutoMigrateTables(baseRes, &transformationRule20230322{})
	if err != nil {
		return err
	}
	var scopes []job20230322
	err = baseRes.GetDal().All(&scopes)
	if err != nil {
		return err
	}
	// get all rules that are not referenced.
	idMap := make(map[uint64]uint64)
	for _, scope := range scopes {
		if scope.TransformationRuleId > 0 && idMap[scope.TransformationRuleId] == 0 {
			idMap[scope.TransformationRuleId] = scope.ConnectionId
		}
	}
	// set connection_id for rules
	for trId, cId := range idMap {
		err = baseRes.GetDal().UpdateColumn(
			&job20230322{}, "connection_id", cId,
			dal.Where("id = ?", trId))
		if err != nil {
			return err
		}
	}
	// delete all rules that are not referenced.
	return baseRes.GetDal().Delete(&job20230322{}, dal.Where("connection_id IS NULL OR connection_id = 0"))
}

func (*addConnectionIdToTransformationRule) Version() uint64 {
	return 20230322150357
}

func (*addConnectionIdToTransformationRule) Name() string {
	return "add connection_id to _tool_jenkins_transformation_rules"
}
