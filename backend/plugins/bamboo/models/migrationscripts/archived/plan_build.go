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

package archived

import (
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"time"
)

type BambooPlanBuild struct {
	ConnectionId             uint64 `gorm:"primaryKey"`
	PlanBuildKey             string `gorm:"primaryKey"`
	Expand                   string `json:"expand"`
	Number                   int    `json:"number"`
	BuildNumber              int    `json:"buildNumber"`
	PlanName                 string `json:"planName"`
	PlanKey                  string
	ProjectName              string `json:"projectName"`
	ProjectKey               string
	BuildResultKey           string     `json:"buildResultKey"`
	LifeCycleState           string     `json:"lifeCycleState"`
	BuildStartedTime         *time.Time `json:"buildStartedTime"`
	PrettyBuildStartedTime   string     `json:"prettyBuildStartedTime"`
	BuildCompletedTime       *time.Time `json:"buildCompletedTime"`
	BuildCompletedDate       *time.Time `json:"buildCompletedDate"`
	PrettyBuildCompletedTime string     `json:"prettyBuildCompletedTime"`
	BuildDurationInSeconds   int        `json:"buildDurationInSeconds"`
	BuildDuration            int        `json:"buildDuration"`
	BuildDurationDescription string     `json:"buildDurationDescription"`
	BuildRelativeTime        string     `json:"buildRelativeTime"`
	VcsRevisionKey           string     `json:"vcsRevisionKey"`
	BuildTestSummary         string     `json:"buildTestSummary"`
	SuccessfulTestCount      int        `json:"successfulTestCount"`
	FailedTestCount          int        `json:"failedTestCount"`
	QuarantinedTestCount     int        `json:"quarantinedTestCount"`
	SkippedTestCount         int        `json:"skippedTestCount"`
	Continuable              bool       `json:"continuable"`
	OnceOff                  bool       `json:"onceOff"`
	Restartable              bool       `json:"restartable"`
	NotRunYet                bool       `json:"notRunYet"`
	Finished                 bool       `json:"finished"`
	Successful               bool       `json:"successful"`
	BuildReason              string     `json:"buildReason"`
	ReasonSummary            string     `json:"reasonSummary"`
	State                    string     `json:"state"`
	BuildState               string     `json:"buildState"`
	PlanResultKey            string
	archived.NoPKModel
}

func (BambooPlanBuild) TableName() string {
	return "_tool_bamboo_plan_builds"
}
