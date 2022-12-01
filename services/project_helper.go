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
	goerror "errors"
	"fmt"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models"
	"gorm.io/gorm"
)

// CreateDbProject accepts a project instance and insert it to database
func CreateDbProject(project *models.Project) errors.Error {
	err := db.Create(project).Error
	if err != nil {
		return errors.Default.Wrap(err, "error creating DB project")
	}
	return nil
}

// CreateDbProjectMetric accepts a project metric instance and insert it to database
func CreateDbProjectMetric(projectMetric *models.ProjectMetric) errors.Error {
	err := db.Create(projectMetric).Error
	if err != nil {
		return errors.Default.Wrap(err, "error creating DB project metric")
	}
	return nil
}

// SaveDbProject save a project instance and update it to database
func SaveDbProject(project *models.Project) errors.Error {
	err := db.Save(project).Error
	if err != nil {
		return errors.Default.Wrap(err, "error saving DB project")
	}
	return nil
}

// SaveDbProjectMetric save a project instance and update it to database
func SaveDbProjectMetric(projectMetric *models.ProjectMetric) errors.Error {
	err := db.Save(projectMetric).Error
	if err != nil {
		return errors.Default.Wrap(err, "error saving DB project metric")
	}
	return nil
}

// GetDbProjects returns a paginated list of Project based on `query`
func GetDbProjects(query *ProjectQuery) ([]*models.Project, int64, errors.Error) {
	projects := make([]*models.Project, 0)
	db := db.Model(projects).Order("created_at desc")

	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, errors.Default.Wrap(err, "error getting DB count of project")
	}
	db = processDbClausesWithPager(db, query.PageSize, query.Page)

	err = db.Find(&projects).Error
	if err != nil {
		return nil, 0, errors.Default.Wrap(err, "error finding DB project")
	}

	return projects, count, nil
}

// GetDbProject returns the detail of a given project name
func GetDbProject(name string) (*models.Project, errors.Error) {
	project := &models.Project{}
	project.Name = name

	err := db.First(project).Error
	if err != nil {
		if goerror.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound.Wrap(err, fmt.Sprintf("could not find project [%s] in DB", name))
		}
		return nil, errors.Default.Wrap(err, "error getting project from DB")
	}

	return project, nil
}

// GetDbProjectMetric returns the detail of a given project name
func GetDbProjectMetric(projectName string, pluginName string) (*models.ProjectMetric, errors.Error) {
	projectMetric := &models.ProjectMetric{}
	projectMetric.ProjectName = projectName
	projectMetric.PluginName = pluginName

	err := db.First(projectMetric).Error
	if err != nil {
		if goerror.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound.Wrap(err, fmt.Sprintf("could not find project metric [%s][%s] in DB", projectName, pluginName))
		}
		return nil, errors.Default.Wrap(err, "error getting project metric from DB")
	}

	return projectMetric, nil
}

// GetDbProjectMetrics returns all of Metrics of a given project name
func GetDbProjectMetrics(projectName string) (*[]models.ProjectMetric, int64, errors.Error) {
	projectMetrics := make([]models.ProjectMetric, 0)
	db := db.Table(models.ProjectMetric{}.TableName()).Where("project_name = ?", projectName)

	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, errors.Default.Wrap(err, fmt.Sprintf("could not get project metric count for projectName [%s] in DB", projectName))
	}

	if count == 0 {
		return nil, 0, nil
	}

	err = db.Find(&projectMetrics).Error
	if err != nil {
		return nil, 0, errors.Default.Wrap(err, fmt.Sprintf("could not find project metric for projectName [%s] in DB", projectName))
	}

	return &projectMetrics, count, nil
}

func removeAllDbProjectMetricsByProjectName(projectName string) errors.Error {
	err := db.Delete(&models.ProjectMetric{}, "project_name = ?", projectName).Error
	if err != nil {
		return errors.Default.Wrap(err, "error deleting ProjectMetrics from DB")
	}
	return nil
}

// encryptProject
/*func encryptProject(project *models.Project) (*models.Project, errors.Error) {
	encKey := config.GetConfig().GetString(core.EncodeKeyEnvStr)

	describeEncrypt, err := core.Encrypt(encKey, project.Description)
	if err != nil {
		return nil, err
	}
	project.Description = describeEncrypt

	return project, nil
}

// encryptProjectMetric
func encryptProjectMetric(projectMetric *models.ProjectMetric) (*models.ProjectMetric, errors.Error) {
	encKey := config.GetConfig().GetString(core.EncodeKeyEnvStr)

	pluginOption, err := core.Encrypt(encKey, projectMetric.PluginOption)
	if err != nil {
		return nil, err
	}
	projectMetric.PluginOption = pluginOption

	return projectMetric, nil
}*/

// decryptProject
/*func decryptProject(project *models.Project) (*models.Project, errors.Error) {
	encKey := config.GetConfig().GetString(core.EncodeKeyEnvStr)

	describe, err := core.Decrypt(encKey, project.Description)
	if err != nil {
		return nil, err
	}
	project.Description = describe

	return project, nil
}

// decryptProjectMetric
func decryptProjectMetric(projectMetric *models.ProjectMetric) (*models.ProjectMetric, errors.Error) {
	encKey := config.GetConfig().GetString(core.EncodeKeyEnvStr)

	pluginOption, err := core.Decrypt(encKey, projectMetric.PluginOption)
	if err != nil {
		return nil, err
	}
	projectMetric.PluginOption = pluginOption

	return projectMetric, nil
}*/
