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
	"encoding/json"
	goerror "errors"
	"fmt"
	"strings"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/logger"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/go-playground/validator/v10"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// BlueprintQuery is a query for GetBlueprints
type BlueprintQuery struct {
	Enable   *bool  `form:"enable,omitempty"`
	IsManual *bool  `form:"is_manual"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Label    string `form:"label"`
}

var (
	blueprintLog = logger.Global.Nested("blueprint")
	vld          = validator.New()
)

// CreateBlueprint accepts a Blueprint instance and insert it to database
func CreateBlueprint(blueprint *models.Blueprint) errors.Error {
	err := validateBlueprintAndMakePlan(blueprint)
	if err != nil {
		return err
	}
	dbBlueprint := parseDbBlueprint(blueprint)
	dbBlueprint, err = encryptDbBlueprint(dbBlueprint)
	if err != nil {
		return err
	}
	err = SaveDbBlueprint(dbBlueprint)
	if err != nil {
		return err
	}
	blueprint.Model = dbBlueprint.Model
	err = ReloadBlueprints(cronManager)
	if err != nil {
		return errors.Internal.Wrap(err, "error reloading blueprints")
	}
	return nil
}

// GetBlueprints returns a paginated list of Blueprints based on `query`
func GetBlueprints(query *BlueprintQuery) ([]*models.Blueprint, int64, errors.Error) {
	dbBlueprints, count, err := GetDbBlueprints(query)
	if err != nil {
		return nil, 0, errors.Convert(err)
	}
	blueprints := make([]*models.Blueprint, 0)
	for _, dbBlueprint := range dbBlueprints {
		dbBlueprint, err = decryptDbBlueprint(dbBlueprint)
		if err != nil {
			return nil, 0, err
		}
		blueprint := parseBlueprint(dbBlueprint)
		blueprints = append(blueprints, blueprint)
	}
	return blueprints, count, nil
}

// GetBlueprint returns the detail of a given Blueprint ID
func GetBlueprint(blueprintId uint64) (*models.Blueprint, errors.Error) {
	dbBlueprint, err := GetDbBlueprint(blueprintId)
	if err != nil {
		if goerror.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound.New("blueprint not found")
		}
		return nil, errors.Internal.Wrap(err, "error getting the blueprint from database")
	}
	dbBlueprint, err = decryptDbBlueprint(dbBlueprint)
	if err != nil {
		return nil, err
	}
	blueprint := parseBlueprint(dbBlueprint)
	return blueprint, nil
}

// GetBlueprintByProjectName returns the detail of a given ProjectName
func GetBlueprintByProjectName(projectName string) (*models.Blueprint, errors.Error) {
	dbBlueprint, err := GetDbBlueprintByProjectName(projectName)
	if err != nil {
		// Allow specific projectName to fail to find the corresponding blueprint
		if goerror.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Internal.Wrap(err, fmt.Sprintf("error getting the blueprint from database with project %s", projectName))
	}
	dbBlueprint, err = decryptDbBlueprint(dbBlueprint)
	if err != nil {
		return nil, err
	}
	blueprint := parseBlueprint(dbBlueprint)
	return blueprint, nil
}

func validateBlueprintAndMakePlan(blueprint *models.Blueprint) errors.Error {
	// validation
	err := vld.Struct(blueprint)
	if err != nil {
		return errors.BadInput.WrapRaw(err)
	}

	// checking if the project exist
	if blueprint.ProjectName != "" {
		_, err := GetProject(blueprint.ProjectName)
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("invalid projectName: [%s] for the blueprint [%s]", blueprint.ProjectName, blueprint.Name))
		}
	}
	if strings.ToLower(blueprint.CronConfig) == "manual" {
		blueprint.IsManual = true
	}
	if !blueprint.IsManual {
		_, err = cron.ParseStandard(blueprint.CronConfig)
		if err != nil {
			return errors.Default.Wrap(err, "invalid cronConfig")
		}
	}
	if blueprint.Mode == models.BLUEPRINT_MODE_ADVANCED {
		plan := make(core.PipelinePlan, 0)
		err = errors.Convert(json.Unmarshal(blueprint.Plan, &plan))
		if err != nil {
			return errors.Default.Wrap(err, "invalid plan")
		}
		// tasks should not be empty
		if len(plan) == 0 || len(plan[0]) == 0 {
			return errors.Default.New("empty plan")
		}
	} else if blueprint.Mode == models.BLUEPRINT_MODE_NORMAL {
		plan, err := MakePlanForBlueprint(blueprint)
		if err != nil {
			return errors.Default.Wrap(err, "invalid plan")
		}
		blueprint.Plan, err = errors.Convert01(json.Marshal(plan))
		if err != nil {
			return errors.Default.Wrap(err, "failed to markshal plan")
		}
	}

	return nil
}

func saveBlueprint(blueprint *models.Blueprint) (*models.Blueprint, errors.Error) {
	// validation
	err := validateBlueprintAndMakePlan(blueprint)
	if err != nil {
		return nil, errors.BadInput.WrapRaw(err)
	}

	// save
	dbBlueprint := parseDbBlueprint(blueprint)
	dbBlueprint, err = encryptDbBlueprint(dbBlueprint)
	if err != nil {
		return nil, err
	}
	err = SaveDbBlueprint(dbBlueprint)
	if err != nil {
		return nil, err
	}

	// reload schedule
	err = ReloadBlueprints(cronManager)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error reloading blueprints")
	}
	// done
	return blueprint, nil
}

// PatchBlueprintEnableByProjectName FIXME ...
func PatchBlueprintEnableByProjectName(projectName string, enable bool) (*models.Blueprint, errors.Error) {
	blueprint, err := GetBlueprintByProjectName(projectName)
	if err != nil {
		return nil, err
	}

	if blueprint == nil {
		return nil, errors.Default.New(fmt.Sprintf("do not surpport to set enable for projectName:[%s] ,because it has no blueprint.", projectName))
	}

	blueprint.Enable = enable

	blueprint, err = saveBlueprint(blueprint)
	if err != nil {
		return nil, err
	}

	return blueprint, nil
}

// PatchBlueprint FIXME ...
func PatchBlueprint(id uint64, body map[string]interface{}) (*models.Blueprint, errors.Error) {
	// load record from db
	blueprint, err := GetBlueprint(id)
	if err != nil {
		return nil, err
	}

	originMode := blueprint.Mode
	err = helper.DecodeMapStruct(body, blueprint)
	if err != nil {
		return nil, err
	}
	// make sure mode is not being update
	if originMode != blueprint.Mode {
		return nil, errors.Default.New("mode is not updatable")
	}

	blueprint, err = saveBlueprint(blueprint)
	if err != nil {
		return nil, err
	}

	return blueprint, nil
}

// DeleteBlueprint FIXME ...
func DeleteBlueprint(id uint64) errors.Error {
	err := DeleteDbBlueprint(id)
	if err != nil {
		return errors.Internal.Wrap(err, fmt.Sprintf("error deleting blueprint %d", id))
	}
	err = ReloadBlueprints(cronManager)
	if err != nil {
		return errors.Internal.Wrap(err, "error reloading blueprints")
	}
	return nil
}

// ReloadBlueprints FIXME ...
func ReloadBlueprints(c *cron.Cron) errors.Error {
	enable := true
	isManual := false
	dbBlueprints, _, err := GetDbBlueprints(&BlueprintQuery{Enable: &enable, IsManual: &isManual})
	if err != nil {
		return err
	}
	for _, e := range c.Entries() {
		c.Remove(e.ID)
	}
	c.Stop()
	for _, dbBlueprint := range dbBlueprints {
		dbBlueprint, err = decryptDbBlueprint(dbBlueprint)
		if err != nil {
			return err
		}
		blueprint := parseBlueprint(dbBlueprint)
		if err != nil {
			blueprintLog.Error(err, failToCreateCronJob)
			return err
		}
		if _, err := c.AddFunc(blueprint.CronConfig, func() {
			pipeline, err := createPipelineByBlueprint(blueprint)
			if err != nil {
				blueprintLog.Error(err, "run cron job failed")
			} else {
				blueprintLog.Info("Run new cron job successfully, pipeline id: %d", pipeline.ID)
			}
		}); err != nil {
			blueprintLog.Error(err, failToCreateCronJob)
			return errors.Default.Wrap(err, "created cron job failed")
		}
	}
	if len(dbBlueprints) > 0 {
		c.Start()
	}
	log.Info("total %d blueprints were scheduled", len(dbBlueprints))
	return nil
}

func createPipelineByBlueprint(blueprint *models.Blueprint) (*models.Pipeline, errors.Error) {
	var plan core.PipelinePlan
	var err errors.Error
	if blueprint.Mode == models.BLUEPRINT_MODE_NORMAL {
		plan, err = MakePlanForBlueprint(blueprint)
	} else {
		plan, err = blueprint.UnmarshalPlan()
	}
	if err != nil {
		return nil, err
	}
	newPipeline := models.NewPipeline{}
	newPipeline.Plan = plan
	newPipeline.Name = blueprint.Name
	newPipeline.BlueprintId = blueprint.ID
	newPipeline.Labels = blueprint.Labels
	pipeline, err := CreatePipeline(&newPipeline)
	// Return all created tasks to the User
	if err != nil {
		blueprintLog.Error(err, failToCreateCronJob)
		return nil, errors.Convert(err)
	}
	return pipeline, nil
}

// MakePlanForBlueprint generates pipeline plan by version
func MakePlanForBlueprint(blueprint *models.Blueprint) (core.PipelinePlan, errors.Error) {
	bpSettings := new(models.BlueprintSettings)
	err := errors.Convert(json.Unmarshal(blueprint.Settings, bpSettings))
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("settings:%s", string(blueprint.Settings)))
	}

	bpSyncPolicy := new(core.BlueprintSyncPolicy)
	bpSyncPolicy.Version = bpSettings.Version
	bpSyncPolicy.SkipOnFail = blueprint.SkipOnFail
	bpSyncPolicy.CreatedDateAfter = blueprint.CreatedDateAfter

	var plan core.PipelinePlan
	switch bpSettings.Version {
	case "1.0.0":
		// Notice: v1 not complete SkipOnFail & CreatedDateAfter
		plan, err = GeneratePlanJsonV100(bpSettings)
	case "2.0.0":
		if blueprint.ProjectName == "" {
			return nil, errors.BadInput.New("projectName is required for blueprint v2.0.0")
		}
		// load project metric plugins and convert it to a map
		metrics := make(map[string]json.RawMessage)
		projectMetrics := make([]models.ProjectMetric, 0)
		db.Find(&projectMetrics, "project_name = ? AND enable = ?", blueprint.ProjectName, true)
		for _, projectMetric := range projectMetrics {
			metrics[projectMetric.PluginName] = json.RawMessage(projectMetric.PluginOption)
		}
		plan, err = GeneratePlanJsonV200(blueprint.ProjectName, bpSyncPolicy, bpSettings, metrics)
	default:
		return nil, errors.Default.New(fmt.Sprintf("unknown version of blueprint settings: %s", bpSettings.Version))
	}
	if err != nil {
		return nil, err
	}
	return WrapPipelinePlans(bpSettings.BeforePlan, plan, bpSettings.AfterPlan)
}

// WrapPipelinePlans merges multiple pipelines and append before and after pipeline
func WrapPipelinePlans(beforePlanJson json.RawMessage, mainPlan core.PipelinePlan, afterPlanJson json.RawMessage) (core.PipelinePlan, errors.Error) {
	beforePipelinePlan := core.PipelinePlan{}
	afterPipelinePlan := core.PipelinePlan{}

	if beforePlanJson != nil {
		err := errors.Convert(json.Unmarshal(beforePlanJson, &beforePipelinePlan))
		if err != nil {
			return nil, err
		}
	}
	if afterPlanJson != nil {
		err := errors.Convert(json.Unmarshal(afterPlanJson, &afterPipelinePlan))
		if err != nil {
			return nil, err
		}
	}

	return SequencializePipelinePlans(beforePipelinePlan, mainPlan, afterPipelinePlan), nil
}

// ParallelizePipelinePlans merges multiple pipelines into one unified plan
// by assuming they can be executed in parallel
func ParallelizePipelinePlans(plans ...core.PipelinePlan) core.PipelinePlan {
	merged := make(core.PipelinePlan, 0)
	// iterate all pipelineTasks and try to merge them into `merged`
	for _, plan := range plans {
		// add all stages from plan to merged
		for index, stage := range plan {
			if index >= len(merged) {
				merged = append(merged, nil)
			}
			// add all tasks from plan to target respectively
			merged[index] = append(merged[index], stage...)
		}
	}
	return merged
}

// SequencializePipelinePlans merges multiple pipelines into one unified plan
// by assuming they must be executed in sequencial order
func SequencializePipelinePlans(plans ...core.PipelinePlan) core.PipelinePlan {
	merged := make(core.PipelinePlan, 0)
	// iterate all pipelineTasks and try to merge them into `merged`
	for _, plan := range plans {
		merged = append(merged, plan...)
	}
	return merged
}

// TriggerBlueprint triggers blueprint immediately
func TriggerBlueprint(id uint64) (*models.Pipeline, errors.Error) {
	// load record from db
	blueprint, err := GetBlueprint(id)
	if err != nil {
		return nil, err
	}
	pipeline, err := createPipelineByBlueprint(blueprint)
	// done
	return pipeline, err
}
