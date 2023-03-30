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
	"fmt"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

// ProjectQuery used to query projects as the api project input
type ProjectQuery struct {
	Pagination
}

// GetProjects returns a paginated list of Projects based on `query`
func GetProjects(query *ProjectQuery) ([]*models.Project, int64, errors.Error) {
	// verify input
	if err := VerifyStruct(query); err != nil {
		return nil, 0, err
	}
	clauses := []dal.Clause{
		dal.From(&models.Project{}),
	}

	count, err := db.Count(clauses...)
	if err != nil {
		return nil, 0, errors.Default.Wrap(err, "error getting DB count of project")
	}

	clauses = append(clauses,
		dal.Orderby("created_at DESC"),
		dal.Offset(query.GetSkip()),
		dal.Limit(query.GetPageSize()),
	)
	projects := make([]*models.Project, 0)
	err = db.All(&projects, clauses...)
	if err != nil {
		return nil, 0, errors.Default.Wrap(err, "error finding DB project")
	}

	return projects, count, nil
}

// CreateProject accepts a project instance and insert it to database
func CreateProject(projectInput *models.ApiInputProject) (*models.ApiOutputProject, errors.Error) {
	// verify input
	if err := VerifyStruct(projectInput); err != nil {
		return nil, err
	}

	// create transaction to updte multiple tables
	var err errors.Error
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			err = tx.Rollback()
			if err != nil {
				logger.Error(err, "PatchProject: failed to rollback")
			}
		}
	}()

	// create project first
	project := &models.Project{}
	project.BaseProject = projectInput.BaseProject
	err = db.Create(project)
	if err != nil {
		if db.IsDuplicationError(err) {
			return nil, errors.BadInput.New(fmt.Sprintf("A project with name [%s] already exists", project.Name))
		}
		return nil, errors.Default.Wrap(err, "error creating DB project")
	}

	// check if need flush the Metrics
	if projectInput.Metrics != nil {
		err = refreshProjectMetrics(tx, projectInput)
		if err != nil {
			return nil, err
		}
	}

	// all good, commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return makeProjectOutput(&projectInput.BaseProject)
}

// GetProject returns a Project
func GetProject(name string) (*models.ApiOutputProject, errors.Error) {
	// verify input
	if name == "" {
		return nil, errors.BadInput.New("project name is missing")
	}

	// load project
	project := &models.Project{}
	err := db.First(project, dal.Where("name = ?", name))
	if err != nil {
		if db.IsErrorNotFound(err) {
			return nil, errors.NotFound.Wrap(err, fmt.Sprintf("could not find project [%s] in DB", name))
		}
		return nil, errors.Default.Wrap(err, "error getting project from DB")
	}

	// convert to api output
	return makeProjectOutput(&project.BaseProject)
}

// PatchProject FIXME ...
func PatchProject(name string, body map[string]interface{}) (*models.ApiOutputProject, errors.Error) {
	projectInput := &models.ApiInputProject{}

	// load input
	err := helper.DecodeMapStruct(body, projectInput, true)
	if err != nil {
		return nil, err
	}

	// wrap all operation inside a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			err = tx.Rollback()
			if err != nil {
				logger.Error(err, "PatchProject: failed to rollback")
			}
		}
	}()

	project := &models.Project{}
	err = tx.First(project, dal.Where("name = ?", name), dal.Lock(true, false))
	if err != nil {
		return nil, err
	}

	// allowed to changed the name
	if projectInput.Name == "" {
		projectInput.Name = name
	}
	project.BaseProject = projectInput.BaseProject

	// name changed, updates the related entities as well
	if name != project.Name {
		// ProjectMetric
		err = tx.UpdateColumn(
			&models.ProjectMetricSetting{},
			"project_name", project.Name,
			dal.Where("project_name = ?", name),
		)
		if err != nil {
			return nil, err
		}

		// ProjectPrMetric
		err = tx.UpdateColumn(
			&crossdomain.ProjectPrMetric{},
			"project_name", project.Name,
			dal.Where("project_name = ?", name),
		)
		if err != nil {
			return nil, err
		}

		// ProjectIssueMetric
		err = tx.UpdateColumn(
			&crossdomain.ProjectIssueMetric{},
			"project_name", project.Name,
			dal.Where("project_name = ?", name),
		)
		if err != nil {
			return nil, err
		}

		// ProjectMapping
		err = tx.UpdateColumn(
			&crossdomain.ProjectMapping{},
			"project_name", project.Name,
			dal.Where("project_name = ?", name),
		)
		if err != nil {
			return nil, err
		}

		// Blueprint
		err = tx.UpdateColumn(
			&models.Blueprint{},
			"project_name", project.Name,
			dal.Where("project_name = ?", name),
		)
		if err != nil {
			return nil, err
		}
		// rename project
		err = tx.UpdateColumn(
			&models.Project{},
			"name", project.Name,
			dal.Where("name = ?", name),
		)
	}

	// Blueprint
	err = tx.UpdateColumn(
		&models.Blueprint{},
		"enable", projectInput.Enable,
		dal.Where("project_name = ?", name),
	)
	if err != nil {
		return nil, err
	}

	// refresh project metrics if needed
	if projectInput.Metrics != nil {
		err = refreshProjectMetrics(tx, projectInput)
		if err != nil {
			return nil, err
		}
	}

	// update project itself
	err = tx.Update(project)
	if err != nil {
		return nil, err
	}

	// commit the transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// all good, render output
	return makeProjectOutput(&projectInput.BaseProject)
}

func refreshProjectMetrics(tx dal.Transaction, projectInput *models.ApiInputProject) errors.Error {
	err := tx.Delete(&models.ProjectMetricSetting{}, dal.Where("project_name = ?", projectInput.Name))
	if err != nil {
		return err
	}

	for _, baseMetric := range *projectInput.Metrics {
		err = tx.Create(&models.ProjectMetricSetting{
			BaseProjectMetricSetting: models.BaseProjectMetricSetting{
				ProjectName: projectInput.Name,
				BaseMetric:  baseMetric,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func makeProjectOutput(baseProject *models.BaseProject) (*models.ApiOutputProject, errors.Error) {
	projectOutput := &models.ApiOutputProject{}
	projectOutput.BaseProject = *baseProject
	// load project metrics
	projectMetrics := make([]models.ProjectMetricSetting, 0)
	err := db.All(&projectMetrics, dal.Where("project_name = ?", projectOutput.Name))
	if err != nil {
		return nil, errors.Default.Wrap(err, "failed to load project metrics")
	}
	// convert metric to api output
	if len(projectMetrics) > 0 {
		baseMetric := make([]models.BaseMetric, len(projectMetrics))
		for i, projectMetric := range projectMetrics {
			baseMetric[i] = projectMetric.BaseMetric
		}
		projectOutput.Metrics = &baseMetric
	}

	// load blueprint
	projectOutput.Blueprint, err = GetBlueprintByProjectName(projectOutput.Name)
	if err != nil {
		return nil, errors.Default.Wrap(err, "Error to get blueprint by project")
	}
	return projectOutput, err
}
