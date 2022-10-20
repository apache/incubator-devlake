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
	"github.com/apache/incubator-devlake/plugins/bitbucket/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/plugins/core"
)

type addDeployment20221013 struct{}

func (*addDeployment20221013) Up(basicRes core.BasicRes) errors.Error {
	return basicRes.GetDal().AutoMigrate(&archived.BitbucketDeployment{})
}

func (*addDeployment20221013) Version() uint64 {
	return 20221013152349
}

func (*addDeployment20221013) Name() string {
	return "bitbucket add _tool_bitbucket_deployments table"
}
