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
	plugin "github.com/apache/incubator-devlake/core/plugin"
)

// All return all the migration scripts of framework
func All() []plugin.MigrationScript {
	return []plugin.MigrationScript{
		new(addFrameworkTables),
		new(renamePipelineStepToStage),
		new(addSubtaskToTaskTable),
		new(addBlueprintMode),
		new(renameTasksToPlan),
		new(resetDomainTables),
		new(addCommitFileComponent),
		new(removeNotes),
		new(addProjectMapping),
		new(renameColumnsOfPullRequestIssue),
		new(addNoPKModelToCommitParent),
		new(addSubtasksTable),
		new(addCICDTables),
		new(renameColumnsOfPrCommentIssueComment),
		new(modifyTablesForDora),
		new(addTypeToBoard),
		new(encryptBlueprint),
		new(encryptPipeline),
		new(modifyCicdPipeline),
		new(modifyCICDTasks),
		new(addOriginChangeValueForPr),
		new(fixCommitFileIdTooLong),
		new(addRawDataOriginToBoardRepos),
		new(renamePipelineCommits),
		new(commitLineChange),
		new(changeLeadTimeMinutesToInt64),
		new(addRepoSnapshot),
		new(createCollectorState),
		new(removeCicdPipelineRelation),
		new(addCicdScopeDropBuildsJobs),
		new(addSkipOnFail),
		new(modifyCommitsDiffs),
		new(addProjectPrMetric),
		new(addProjectTables),
		new(addProjectToBluePrint),
		new(addProjectIssueMetric),
		new(addLabels),
		new(renameFiledsInProjectPrMetric),
		new(addEnableToProjectMetric),
		new(addCollectorMeta20221125),
		new(addOriginalProject),
		new(addErrorName),
		new(encryptTask221221),
		new(renameProjectMetrics),
		new(addOriginalTypeToIssue221230),
		new(addTimeAfterToCollectorMeta20230213),
		new(addCodeQuality),
		new(modifyIssueStorypointToFloat64),
		new(addCommitShaIndex),
		new(removeCreatedDateAfterFromCollectorMeta20230223),
		new(addHostNamespaceRepoName),
		new(renameCollectorTapStateTable),
		new(renameCicdPipelineRepoToRepoUrl),
		new(addCicdDeploymentCommits),
		new(renameDeploymentIdForPrProjectMetric),
		new(addCommitAuthoredDate),
		new(addOriginalStatusToPullRequest20230508),
		new(addIssueAssignee20230402),
		new(addCalendarMonths),
		new(modifyPrLabelsAndComments),
		new(renameFinishedCommitsDiffs),
		new(addUpdatedDateToIssueComments),
		new(addApiKeyTables),
		new(addIssueRelationship),
		new(tasksUsesJSON),
		new(modifyCicdPipelinesToText),
		new(dropTapStateTable),
		new(addCICDDeploymentsTable),
		new(normalizeBpSettings),
		new(addSyncPolicy),
		new(addIssueCustomArrayField),
		new(removePositionFromPullRequestComments),
		new(changeDurationSecToFloat64),
		new(addSomeDateFieldsToDevopsTables),
		new(addOriginalStatusAndResultToDevOpsTables),
		new(addQueuedDurationSecFieldToDevopsTables),
		new(addCommitMsgtoDeploymentCommit),
		new(modifyIssueOriginalTypeLength),
		new(addCommitMsgtoPipelineCommit),
		new(modfiyFieldsSort),
		new(modifyIssueLeadTimeMinutesToUint),
		new(addUrgencyToIssues),
		new(modifyRefsIdLength),
		new(addOriginalEnvironmentToCicdDeploymentsAndCicdDeploymentCommits),
		new(addSubtabknameToDeployment),
		new(addStore),
		new(addSubtaskField),
		new(addDisplayTitleAndUrl),
		new(addSubtaskStates),
		new(addCicdRelease),
		new(addCommitShaToCicdRelease),
		new(updateIssueKeyType),
		new(updatePluginOptionInProjectMetricSetting),
		new(modifyCicdDeploymentCommitsRepoUrlLength),
		new(modifyCicdPipelineCommitsRepoUrlLength),
		new(addPrAssigneeAndReviewer),
		new(modifyPrAssigneeAndReviewerId),
		new(addPullRequestIdIndexToPullRequestCommits),
		new(addPullRequestIdIndexToPullRequestComments),
	}
}
