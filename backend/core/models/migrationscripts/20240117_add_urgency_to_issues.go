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
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.MigrationScript = (*addUrgencyToIssues)(nil)

type issue20240117 struct {
	Urgency string `gorm:"type:varchar(255)"`
}

func (issue20240117) TableName() string {
	return "issues"
}

type addUrgencyToIssues struct{}

func (*addUrgencyToIssues) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if err := db.AutoMigrate(&issue20240117{}); err != nil {
		return err
	}
	return nil
}

func (*addUrgencyToIssues) Version() uint64 {
	return 20240117142100
}

func (*addUrgencyToIssues) Name() string {
	return "add urgency to issues table"
}
