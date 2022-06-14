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

	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/jenkins/models/migrationscripts/archived"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InitSchemas struct{}

func (*InitSchemas) Up(ctx context.Context, db *gorm.DB) error {

	rawTableList := []string{
		"_raw_jenkins_api_jobs",
		"_raw_jenkins_api_builds",
	}
	for _, v := range rawTableList {
		err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", v)).Error
		if err != nil {
			return err
		}
	}

	err := db.Migrator().DropTable(
		&archived.JenkinsJob{},
		&archived.JenkinsBuild{},
	)
	if err != nil {
		return err
	}

	err = db.Migrator().AutoMigrate(
		&archived.JenkinsJob{},
		&archived.JenkinsBuild{},
		&archived.JenkinsConnection{},
	)
	if err != nil {
		return err
	}

	conn := &archived.JenkinsConnection{}
	v := config.GetConfig()

	conn.Name = "init jenkins connection"
	conn.ID = 1
	conn.Endpoint = v.GetString("JENKINS_ENDPOINT")
	conn.Proxy = v.GetString("JENKINS_PROXY")
	conn.RateLimit = v.GetInt("JENKINS_API_REQUESTS_PER_HOUR")
	conn.Username = v.GetString("JENKINS_USERNAME")
	encKey := v.GetString(core.EncodeKeyEnvStr)
	conn.Password, err = core.Encrypt(encKey, v.GetString("JENKINS_PASSWORD"))
	if err != nil {
		return err
	}
	err = db.Clauses(clause.OnConflict{DoNothing: true}).Create(conn).Error

	if err != nil {
		return err
	}

	return nil

}

func (*InitSchemas) Version() uint64 {
	return 20220614201236
}

func (*InitSchemas) Name() string {
	return "Jenkins init schemas"
}
