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
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*addNamespaceRevisionToDeploymentRevision)(nil)

type addNamespaceRevisionToDeploymentRevision struct{}

type deployment_revision20230606 struct {
	Namespace      string `gorm:"type:varchar(255)"`
	RevisionNumber int32
}

func (deployment_revision20230606) TableName() string {
	return "_tool_kube_deployment_revisions"
}

func (u *addNamespaceRevisionToDeploymentRevision) Up(baseRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(baseRes, &deployment_revision20230606{})
}

func (*addNamespaceRevisionToDeploymentRevision) Version() uint64 {
	return 20230606160400
}

func (*addNamespaceRevisionToDeploymentRevision) Name() string {
	return "add add namespace and revision to deployment revisions"
}
