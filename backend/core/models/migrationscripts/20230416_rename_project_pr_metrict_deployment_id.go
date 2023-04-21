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

var _ plugin.MigrationScript = (*renameDeploymentIdForPrProjectMetric)(nil)

type renameDeploymentIdForPrProjectMetric struct{}

func (*renameDeploymentIdForPrProjectMetric) Up(basicRes context.BasicRes) errors.Error {
	return basicRes.GetDal().RenameColumn("project_pr_metrics", "deployment_id", "deployment_commit_id")
}

func (*renameDeploymentIdForPrProjectMetric) Version() uint64 {
	return 20230416080701
}

func (*renameDeploymentIdForPrProjectMetric) Name() string {
	return "Rename project_pr_metrics.deployment_id to deployment_commit_id"
}
