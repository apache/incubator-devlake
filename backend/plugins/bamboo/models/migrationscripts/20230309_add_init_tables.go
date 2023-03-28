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
	"github.com/apache/incubator-devlake/plugins/bamboo/models/migrationscripts/archived"
)

type addInitTables struct{}

func (u *addInitTables) Up(baseRes context.BasicRes) errors.Error {
	// will be deleted after finish bamboo
	_ = baseRes.GetDal().DropTables(
		&archived.BambooPlan{},
		&archived.BambooJob{},
		&archived.BambooPlanBuild{},
		&archived.BambooPlanBuildVcsRevision{},
		&archived.BambooJobBuild{},
		&archived.BambooTransformationRule{},
		&archived.BambooDeployEnvironment{},
		&archived.BambooDeployBuild{},
	)
	return migrationhelper.AutoMigrateTables(
		baseRes,
		&archived.BambooConnection{},
		&archived.BambooProject{},
		&archived.BambooPlan{},
		&archived.BambooJob{},
		&archived.BambooPlanBuild{},
		&archived.BambooPlanBuildVcsRevision{},
		&archived.BambooJobBuild{},
		&archived.BambooTransformationRule{},
		&archived.BambooDeployEnvironment{},
		&archived.BambooDeployBuild{},
	)
}

func (*addInitTables) Version() uint64 {
	return 20230315205035
}

func (*addInitTables) Name() string {
	return "bamboo init schemas"
}
