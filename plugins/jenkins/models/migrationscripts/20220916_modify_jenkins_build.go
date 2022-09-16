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
	"context"
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jenkins/api"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"time"
)

type modifyJenkinsBuild struct{}

type JenkinsBuildOld struct {
	archived.NoPKModel
	// collected fields
	ConnectionId      uint64    `gorm:"primaryKey"`
	JobName           string    `gorm:"primaryKey;type:varchar(255)"`
	Duration          float64   // build time
	DisplayName       string    `gorm:"type:varchar(255)"` // "#7"
	EstimatedDuration float64   // EstimatedDuration
	Number            int64     `gorm:"primaryKey"`
	Result            string    // Result
	Timestamp         int64     // start time
	StartTime         time.Time // convered by timestamp
	CommitSha         string    `gorm:"type:varchar(255)"`
	Type              string    `gorm:"index;type:varchar(255)"`
	Class             string    `gorm:"index;type:varchar(255)" `
	TriggeredBy       string    `gorm:"type:varchar(255)"`
	Building          bool
	HasStages         bool
}

func (JenkinsBuildOld) TableName() string {
	return "_tool_jenkins_builds"
}

type JenkinsBuil0916 struct {
	archived.NoPKModel
	// collected fields
	ConnectionId      uint64    `gorm:"primaryKey"`
	JobName           string    `gorm:"index;type:varchar(255)"`
	Duration          float64   // build time
	FullDisplayName   string    `gorm:"primaryKey;type:varchar(255)"` // "#7"
	EstimatedDuration float64   // EstimatedDuration
	Number            int64     `gorm:"index"`
	Result            string    // Result
	Timestamp         int64     // start time
	StartTime         time.Time // convered by timestamp
	Type              string    `gorm:"index;type:varchar(255)"`
	Class             string    `gorm:"index;type:varchar(255)" `
	TriggeredBy       string    `gorm:"type:varchar(255)"`
	Building          bool
	HasStages         bool
}

func (JenkinsBuil0916) TableName() string {
	return "_tool_jenkins_builds"
}

func (*modifyJenkinsBuild) Up(ctx context.Context, db *gorm.DB) errors.Error {
	cursor, err := db.Model(&JenkinsBuildOld{}).Rows()
	if err != nil {
		return errors.Convert(err)
	}
	err = db.Migrator().RenameTable(&JenkinsBuildOld{}, "_tool_jenkins_builds_old")
	if err != nil {
		return errors.Default.Wrap(err, "fail to rename _tool_jenkins_builds")
	}
	err = db.Migrator().AutoMigrate(&JenkinsBuil0916{})
	if err != nil {
		return errors.Default.Wrap(err, "fail to create _tool_jenkins_builds")
	}
	batch, err := helper.NewBatchSave(api.BasicRes, reflect.TypeOf(&JenkinsBuil0916{}), 300)
	if err != nil {
		return errors.Default.Wrap(err, "error getting batch from table")
	}
	defer batch.Close()
	for cursor.Next() {
		build := JenkinsBuildOld{}
		err = db.ScanRows(cursor, &build)
		if err != nil {
			return errors.Convert(err)
		}
		newBuild := &JenkinsBuil0916{
			NoPKModel:         build.NoPKModel,
			ConnectionId:      build.ConnectionId,
			JobName:           build.JobName,
			Duration:          build.Duration,
			FullDisplayName:   build.DisplayName,
			EstimatedDuration: build.EstimatedDuration,
			Number:            build.Number,
			Result:            build.Result,
			Timestamp:         build.Timestamp,
			StartTime:         build.StartTime,
			Type:              build.Type,
			Class:             build.Class,
			TriggeredBy:       build.TriggeredBy,
			Building:          build.Building,
			HasStages:         build.HasStages,
		}
		if strings.Contains(build.DisplayName, build.JobName) {
			newBuild.FullDisplayName = build.DisplayName
		} else {
			newBuild.FullDisplayName = fmt.Sprintf("%s %s", build.JobName, build.DisplayName)
		}
		err = batch.Add(&newBuild)
		if err != nil {
			return errors.Convert(err)
		}
	}
	if err != nil {
		return errors.Convert(err)
	}

	return nil
}

func (*modifyJenkinsBuild) Version() uint64 {
	return 20220916231237
}

func (*modifyJenkinsBuild) Name() string {
	return "Jenkins modify build primary key"
}
