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

package models

import (
	"github.com/apache/devlake/core/models/common"
)

// Sprint Report bucket values, mirroring the four buckets returned by
// GET rest/greenhopper/1.0/rapid/charts/sprintreport. Jira computes these
// once, at sprint close, so persisting them (rather than reconstructing
// from resolution_date) is what makes committed/completed velocity exact.
const (
	SprintReportBucketCompleted              = "completed"
	SprintReportBucketNotCompleted           = "notCompleted"
	SprintReportBucketPunted                 = "punted"
	SprintReportBucketCompletedInOtherSprint = "completedInOtherSprint"
)

// JiraSprintReport is a frozen, per-(board, sprint, issue) snapshot taken
// from Jira's Sprint Report at sprint close. Unlike JiraSprintIssue (which
// is derived from each issue's live resolution_date and therefore
// mis-attributes carryover issues), this table stores Jira's own
// point-in-time bucketing, so it doesn't drift.
type JiraSprintReport struct {
	common.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey"`
	BoardId      uint64 `gorm:"primaryKey"`
	SprintId     uint64 `gorm:"primaryKey"`
	IssueId      uint64 `gorm:"primaryKey"`

	IssueKey string `gorm:"type:varchar(255)"`
	// Bucket is one of the SprintReportBucket* constants above.
	Bucket string `gorm:"type:varchar(32);index"`
	Done   bool

	// StoryPointsAtSprintStart is estimateStatistic.statFieldValue.value in
	// Jira's response ("BOS points") — the estimate as it stood when the
	// sprint began, i.e. what should be summed for *committed* velocity.
	StoryPointsAtSprintStart *float64
	// StoryPointsAtSprintEnd is currentEstimateStatistic.statFieldValue.value
	// ("EOS points") — the estimate as of sprint close, i.e. what should be
	// summed (for Bucket == completed) for *completed* velocity.
	StoryPointsAtSprintEnd *float64
}

func (JiraSprintReport) TableName() string {
	return "_tool_jira_sprint_reports"
}
