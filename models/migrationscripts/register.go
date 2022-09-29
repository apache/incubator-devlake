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
	"github.com/apache/incubator-devlake/plugins/core"
)

// All return all the migration scripts of framework
func All() []core.MigrationScript {
	return []core.MigrationScript{
		new(addFrameTables),
		new(renameStepToStage),
		new(addSubtasksField),
		// new(updateBlueprintMode),
		// new(renameTasksToPlan),
		// new(addDomainTables),
		// new(commitfileComponent),
		// new(removeNotes),
		// new(addProjectMapping),
		// new(renameColumnsOfPullRequestIssue),
		// new(addNoPKModelToCommitParent),
		// new(addSubtasksTable),
		// new(addCICD),
		// new(renameColumnsOfPrCommentIssueComment),
		// new(modifyTablesForDora),
		// new(addTypeFieldInBoard),
		// new(modifyPipeline),
		// new(encryptBLueprint),
		// new(encryptPipeline),
		// new(modifyCICDTasks),
		// new(modifyBoardRepos),
		// new(addOriginChangeValueForPr),
		// new(addCommitFilePathLength),
		// new(renamePipelineCommits),
	}
}
