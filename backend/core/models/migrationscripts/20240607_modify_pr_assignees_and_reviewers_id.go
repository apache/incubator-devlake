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
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type modifyPrAssigneeAndReviewerId struct{}

func (u *modifyPrAssigneeAndReviewerId) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	err := db.DropTables(&archived.PullRequestAssignee{}, &archived.PullRequestReviewer{})
	if err != nil {
		return err
	}
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&archived.PullRequestAssignee{},
		&archived.PullRequestReviewer{},
	)
}

func (*modifyPrAssigneeAndReviewerId) Version() uint64 {
	return 20250607000041
}

func (*modifyPrAssigneeAndReviewerId) Name() string {
	return "modify pull_request_reviewers and pull_request_assignees id columns"
}
