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
	"net/url"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
)

type modfiyFieldsSort struct{}

func (*modfiyFieldsSort) Up(baseRes context.BasicRes) errors.Error {
	dbUrl := baseRes.GetConfig("DB_URL")
	if dbUrl == "" {
		return errors.BadInput.New("DB_URL is required")
	}
	u, errParse := url.Parse(dbUrl)
	if errParse != nil {
		return errors.Convert(errParse)
	}
	db := baseRes.GetDal()
	if u.Scheme == "mysql" {
		// issues
		err := db.Exec("alter table issues modify original_type varchar(500) after type;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table issues modify story_point DOUBLE after priority;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table issues modify lead_time_minutes bigint after updated_date;")
		if err != nil {
			return err
		}
		// pull_requests
		err = db.Exec("alter table pull_requests modify base_ref varchar(255) after base_repo_id;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table pull_requests modify base_commit_sha varchar(40) after base_ref;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table pull_requests modify head_ref varchar(255) after head_repo_id;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table pull_requests modify head_commit_sha varchar(40) after head_ref;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table pull_requests modify merge_commit_sha varchar(40) after head_commit_sha;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table pull_requests modify original_status varchar(100) after status;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table pull_requests modify type varchar(100) after original_status;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table pull_requests modify component varchar(100) after type;")
		if err != nil {
			return err
		}
		// cicd deployment commits
		err = db.Exec("alter table cicd_deployment_commits modify original_status varchar(100) after status;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table cicd_deployment_commits modify original_result varchar(100) after result;")
		if err != nil {
			return err
		}

		err = db.Exec("alter table cicd_deployment_commits modify queued_date DATETIME(3) after created_date;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table cicd_deployment_commits modify queued_duration_sec DOUBLE after queued_date;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table cicd_deployment_commits modify duration_sec DOUBLE after finished_date;")
		if err != nil {
			return err
		}

		// cicd deployments
		err = db.Exec("alter table cicd_deployments modify original_status varchar(100) after status;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table cicd_deployments modify original_result varchar(100) after result;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table cicd_deployments modify queued_date DATETIME(3) after created_date;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table cicd_deployments modify queued_duration_sec DOUBLE after queued_date;")
		if err != nil {
			return err
		}

		// cicd pipelines
		err = db.Exec("alter table cicd_pipelines modify original_status varchar(100) after status;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table cicd_pipelines modify original_result varchar(100) after result;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table cicd_pipelines modify queued_date DATETIME(3) after created_date;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table cicd_pipelines modify queued_duration_sec DOUBLE after queued_date;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table cicd_pipelines modify started_date DATETIME(3) after queued_duration_sec;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table cicd_pipelines modify duration_sec DOUBLE after finished_date;")
		if err != nil {
			return err
		}
		// cicd tasks
		err = db.Exec("alter table cicd_tasks modify original_status varchar(100) after status;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table cicd_tasks modify original_result varchar(100) after result;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table cicd_tasks modify created_date DATETIME(3) after type;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table cicd_tasks modify queued_date DATETIME(3) after created_date;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table cicd_tasks modify queued_duration_sec DOUBLE after queued_date;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table cicd_tasks modify started_date DATETIME(3) after queued_duration_sec;")
		if err != nil {
			return err
		}
		err = db.Exec("alter table cicd_tasks modify duration_sec DOUBLE after finished_date;")
		if err != nil {
			return err
		}
	}

	return nil
}

func (*modfiyFieldsSort) Version() uint64 {
	return 20240116000011
}

func (*modfiyFieldsSort) Name() string {
	return "fix some tables fields sort"
}
