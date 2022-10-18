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
	"github.com/apache/incubator-devlake/plugins/core"
)

var _ core.MigrationScript = (*addOriginChangeValueForPr)(nil)

type addOriginChangeValueForPr struct{}

type PullRequest0913 struct {
	OrigCodingTimespan int64
	OrigReviewLag      int64
	OrigReviewTimespan int64
	OrigDeployTimespan int64
}

func (PullRequest0913) TableName() string {
	return "pull_requests"
}

func (*addOriginChangeValueForPr) Up(basicRes core.BasicRes) errors.Error {
	return basicRes.GetDal().AutoMigrate(&PullRequest0913{})
}

func (*addOriginChangeValueForPr) Version() uint64 {
	return 20220913235535
}

func (*addOriginChangeValueForPr) Name() string {
	return "add origin change lead time for pr"
}
