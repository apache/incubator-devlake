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

package services

import (
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

// ProjectQuery used to query projects as the api project input
type ProjectQuery struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

// CreateProject accepts a project instance and insert it to database
func CreateProject(project *models.Project) errors.Error {
	enProject, err := encryptProject(project)
	if err != nil {
		return err
	}
	err = CreateDbProject(enProject)
	if err != nil {
		return err
	}
	return nil
}

// CreateProjectMetric accepts a ProjectMetric instance and insert it to database
func CreateProjectMetric(projectMetric *models.ProjectMetric) errors.Error {
	enProjectMetric, err := encryptProjectMetric(projectMetric)
	if err != nil {
		return err
	}
	err = CreateDbProjectMetric(enProjectMetric)
	if err != nil {
		return err
	}
	return nil
}

// GetProject returns a Project
func GetProject(name string) (*models.Project, errors.Error) {
	project, err := GetDbProject(name)
	if err != nil {
		return nil, errors.Convert(err)
	}

	project, err = decryptProject(project)
	if err != nil {
		return nil, errors.Convert(err)
	}

	return project, nil
}

// GetProjectMetric returns a ProjectMetric
func GetProjectMetric(projectName string, pluginName string) (*models.ProjectMetric, errors.Error) {
	projectMetric, err := GetDbProjectMetric(projectName, pluginName)
	if err != nil {
		return nil, errors.Convert(err)
	}

	projectMetric, err = decryptProjectMetric(projectMetric)
	if err != nil {
		return nil, errors.Convert(err)
	}

	return projectMetric, nil
}

// GetProjects returns a paginated list of Projects based on `query`
func GetProjects(query *ProjectQuery) ([]*models.Project, int64, errors.Error) {
	projects, count, err := GetDbProjects(query)
	if err != nil {
		return nil, 0, errors.Convert(err)
	}
	for i, project := range projects {
		projects[i], err = decryptProject(project)
		if err != nil {
			return nil, 0, err
		}
	}
	return projects, count, nil
}

// PatchProject FIXME ...
func PatchProject(name string, body map[string]interface{}) (*models.Project, errors.Error) {
	// load record from db
	project, err := GetProject(name)
	if err != nil {
		return nil, err
	}

	err = helper.DecodeMapStruct(body, project)
	if err != nil {
		return nil, err
	}

	// save
	err = SaveDbProject(project)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error saving project")
	}

	// done
	return project, nil
}

// PatchProjectMetric FIXME ...
func PatchProjectMetric(projectName string, pluginName string, body map[string]interface{}) (*models.ProjectMetric, errors.Error) {
	// load record from db
	projectMetric, err := GetDbProjectMetric(projectName, pluginName)
	if err != nil {
		return nil, err
	}

	err = helper.DecodeMapStruct(body, projectMetric)
	if err != nil {
		return nil, err
	}

	// save
	err = SaveDbProjectMetric(projectMetric)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error saving project")
	}

	// done
	return projectMetric, nil
}
