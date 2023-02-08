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

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/plugins/jenkins/models"

	"github.com/apache/incubator-devlake/core/errors"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	aha "github.com/apache/incubator-devlake/helpers/pluginhelper/api/apihelperabstract"
)

func GetJobs(
	apiClient aha.ApiClientAbstract,
	path string,
	pageSize int,
	callback func(job *models.Job) errors.Error,
) errors.Error {
	for i := 0; ; i += pageSize {
		var data struct {
			Jobs []json.RawMessage `json:"jobs"`
		}

		// set query
		query := url.Values{}
		treeValue := fmt.Sprintf("jobs[name,class,url,color,base,jobs,upstreamProjects[name]]{%d,%d}", i, i+pageSize)
		query.Set("tree", treeValue)

		res, err := apiClient.Get(path+"/api/json", query, nil)
		if err != nil {
			return err
		}

		err = helper.UnmarshalResponse(res, &data)
		if err != nil {
			// In some directories, after testing
			// it is found that if the second page is empty, a 500 error will be returned.
			// So we don't think it's an error to return 500 in this case
			if i > 0 && res.StatusCode == http.StatusInternalServerError {
				return nil
			}
			return err
		}

		for _, rawJobs := range data.Jobs {
			job := &models.Job{}
			err1 := json.Unmarshal(rawJobs, job)
			if err1 != nil {
				return errors.Convert(err1)
			}

			err = callback(job)
			if err != nil {
				return err
			}
		}

		if len(data.Jobs) < pageSize {
			return nil
		}
	}
}

func GetJob(
	apiClient aha.ApiClientAbstract,
	path string,
	name string,
	fullName string,
	pageSize int,
	callback func(job *models.Job, isPath bool) errors.Error,
) errors.Error {
	var err errors.Error

	return GetJobs(apiClient, path, pageSize, func(job *models.Job) errors.Error {
		if job.Name != name {
			return nil
		}
		job.Path = path
		job.FullName = fullName

		if job.Jobs != nil {
			err = callback(job, true)
			if err != nil {
				return err
			}
		} else {
			err = callback(job, false)
		}

		return err
	})
}

// request all jobs
func GetAllJobs(
	apiClient aha.ApiClientAbstract,
	path string,
	beforename string,
	pageSize int,
	callback func(job *models.Job, isPath bool) errors.Error,
) errors.Error {
	var err errors.Error
	return GetJobs(apiClient, path, pageSize, func(job *models.Job) errors.Error {
		job.Path = path
		job.FullName = beforename + job.Name

		if job.Jobs != nil {
			err = callback(job, true)
			if err != nil {
				return err
			}
			err = GetAllJobs(apiClient, path+"job/"+job.Name+"/", beforename+job.Name+"/", pageSize, callback)
		} else {
			err = callback(job, false)
		}

		return err
	})
}
