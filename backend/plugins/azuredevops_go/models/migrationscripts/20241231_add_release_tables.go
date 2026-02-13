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
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models/migrationscripts/archived"
)

type addReleaseTables struct {
}

func (u *addReleaseTables) Up(basicRes context.BasicRes) errors.Error {
	err := migrationhelper.AutoMigrateTables(
		basicRes,
		&archived.AzuredevopsRelease{},
		&archived.AzuredevopsReleaseDeployment{},
	)
	return err
}

func (*addReleaseTables) Version() uint64 {
	return 20241231000001
}

func (*addReleaseTables) Name() string {
	return "Add Azure DevOps Release Pipeline tables"
}
