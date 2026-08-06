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

package apiv2models

import "time"

// SprintReportInput drives the per-(board, sprint) Sprint Report collector.
// It's what gets iterated over via the DAL cursor, and re-attached to each
// raw row so the extractor knows which board/sprint a response belongs to.
type SprintReportInput struct {
	BoardId  uint64 `json:"board_id"`
	SprintId uint64 `json:"sprint_id"`
	// UpdateTime is the sprint's CompleteDate; used as the incremental-sync
	// watermark since a sprint report only exists/changes once a sprint closes.
	UpdateTime *time.Time `json:"update_time"`
}

// SprintReportStatFieldValue mirrors Jira's
// {statFieldValue: {value, text}} shape used for per-issue point estimates.
type SprintReportStatFieldValue struct {
	Value *float64 `json:"value"`
	Text  string   `json:"text"`
}

type SprintReportStatistic struct {
	StatFieldValue SprintReportStatFieldValue `json:"statFieldValue"`
}

// SprintReportIssue is one entry inside any of the four bucket lists
// (completedIssues, issuesNotCompletedInCurrentSprint, puntedIssues,
// issuesCompletedInAnotherSprint) in the Sprint Report response.
type SprintReportIssue struct {
	Id       uint64 `json:"id"`
	Key      string `json:"key"`
	TypeName string `json:"typeName"`
	Done     bool   `json:"done"`
	// EstimateStatistic is the issue's estimate as of sprint *start*
	// ("BOS points" / committed).
	EstimateStatistic SprintReportStatistic `json:"estimateStatistic"`
	// CurrentEstimateStatistic is the issue's estimate as of sprint *close*
	// ("EOS points" / completed).
	CurrentEstimateStatistic SprintReportStatistic `json:"currentEstimateStatistic"`
}

// SprintReportContents is the "contents" object of the Sprint Report
// response — the frozen snapshot Jira takes at sprint close.
type SprintReportContents struct {
	CompletedIssues                   []SprintReportIssue `json:"completedIssues"`
	IssuesNotCompletedInCurrentSprint []SprintReportIssue `json:"issuesNotCompletedInCurrentSprint"`
	PuntedIssues                      []SprintReportIssue `json:"puntedIssues"`
	IssuesCompletedInAnotherSprint    []SprintReportIssue `json:"issuesCompletedInAnotherSprint"`
}

// SprintReport is the top-level response of
// GET rest/greenhopper/1.0/rapid/charts/sprintreport?rapidViewId=&sprintId=
type SprintReport struct {
	Contents SprintReportContents `json:"contents"`
}
