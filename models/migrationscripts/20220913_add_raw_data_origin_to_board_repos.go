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
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/plugins/core"
)

var _ core.MigrationScript = (*addRawDataOriginToBoardRepos)(nil)

// addRawDataOriginToBoardRepos add raw data fields to board_repos
type addRawDataOriginToBoardRepos struct{}

type boardRepo20220913 struct {
	archived.NoPKModel
}

func (boardRepo20220913) TableName() string {
	return "board_repos"
}

func (*addRawDataOriginToBoardRepos) Up(basicRes core.BasicRes) errors.Error {
	return basicRes.GetDal().AutoMigrate(&boardRepo20220913{})
}

func (*addRawDataOriginToBoardRepos) Version() uint64 {
	return 20220913232735
}

func (*addRawDataOriginToBoardRepos) Name() string {
	return "modify board repos"
}
