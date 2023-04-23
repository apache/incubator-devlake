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
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
)

type BambooJobBuild struct {
	ConnectionId             uint64     `gorm:"primaryKey"`
	JobBuildKey              string     `gorm:"primaryKey"`
	PlanBuildKey             string     `gorm:"index"`
	Expand                   string     `json:"expand"`
	Number                   int        `json:"number"`
	BuildNumber              int        `json:"buildNumber"`
	JobName                  string     `gorm:"index"`
	JobKey                   string     `gorm:"index"`
	PlanName                 string     `gorm:"index"`
	PlanKey                  string     `gorm:"index"`
	ProjectName              string     `gorm:"index"`
	ProjectKey               string     `gorm:"index"`
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
	Type                     string     `gorm:"type:varchar(255)"`
	Environment              string     `gorm:"type:varchar(255)"`
	JobResultKey             string
	common.NoPKModel
}

func (BambooJobBuild) TableName() string {
	return "_tool_bamboo_job_builds"
}

type ApiBambooJobBuild ApiBambooPlanBuild

func (apiRes *ApiBambooJobBuild) Convert() *BambooJobBuild {
	return &BambooJobBuild{
		ConnectionId:             0,
		JobBuildKey:              apiRes.Key,
		Expand:                   apiRes.Expand,
		Number:                   apiRes.Number,
		BuildNumber:              apiRes.BuildNumber,
		JobName:                  apiRes.PlanName,
		JobKey:                   "",
		PlanName:                 "",
		PlanKey:                  "",
		ProjectName:              apiRes.ProjectName,
		ProjectKey:               "",
		BuildResultKey:           apiRes.BuildResultKey,
		LifeCycleState:           apiRes.LifeCycleState,
		BuildStartedTime:         apiRes.BuildStartedTime,
		PrettyBuildStartedTime:   apiRes.PrettyBuildStartedTime,
		BuildCompletedTime:       apiRes.BuildCompletedTime,
		BuildCompletedDate:       apiRes.BuildCompletedDate,
		PrettyBuildCompletedTime: apiRes.PrettyBuildCompletedTime,
		BuildDurationInSeconds:   apiRes.BuildDurationInSeconds,
		BuildDuration:            apiRes.BuildDuration,
		BuildDurationDescription: apiRes.BuildDurationDescription,
		BuildRelativeTime:        apiRes.BuildRelativeTime,
		VcsRevisionKey:           apiRes.VcsRevisionKey,
		BuildTestSummary:         apiRes.BuildTestSummary,
		SuccessfulTestCount:      apiRes.SuccessfulTestCount,
		FailedTestCount:          apiRes.FailedTestCount,
		QuarantinedTestCount:     apiRes.QuarantinedTestCount,
		SkippedTestCount:         apiRes.SkippedTestCount,
		Continuable:              apiRes.Continuable,
		OnceOff:                  apiRes.OnceOff,
		Restartable:              apiRes.Restartable,
		NotRunYet:                apiRes.NotRunYet,
		Finished:                 apiRes.Finished,
		Successful:               apiRes.Successful,
		BuildReason:              apiRes.BuildReason,
		ReasonSummary:            apiRes.ReasonSummary,
		State:                    apiRes.State,
		BuildState:               apiRes.BuildState,
		JobResultKey:             apiRes.PlanResultKey.Key,
	}
}
