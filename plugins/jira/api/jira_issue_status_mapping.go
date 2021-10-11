package api

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"github.com/mitchellh/mapstructure"
)

func findIssueStatusMappingFromInput(input *core.ApiResourceInput) (*models.JiraIssueStatusMapping, error) {
	// load type mapping
	jiraIssueTypeMapping, err := findIssueTypeMappingByInputParam(input)
	if err != nil {
		return nil, err
	}
	// load status mapping from db
	userStatus := input.Params["userStatus"]
	if userStatus == "" {
		return nil, fmt.Errorf("missing userStatus")
	}
	jiraIssueStatusMapping := &models.JiraIssueStatusMapping{}
	err = lakeModels.Db.First(
		jiraIssueStatusMapping,
		jiraIssueTypeMapping.JiraSourceID,
		jiraIssueTypeMapping.UserType,
		userStatus,
	).Error
	if err != nil {
		return nil, err
	}

	return jiraIssueStatusMapping, nil
}

func syncIssueStatusMappingFromInput(jiraIssueStatusMapping *models.JiraIssueStatusMapping, input *core.ApiResourceInput) error {
	// decode
	err := mapstructure.Decode(input.Body, jiraIssueStatusMapping)
	if err != nil {
		return err
	}
	// validate
	vld := validator.New()
	err = vld.Struct(jiraIssueStatusMapping)
	if err != nil {
		return err
	}
	return nil
}

func wrapIssueStatusDuplicateErr(err error) error {
	if lakeModels.IsDuplicateError(err) {
		return fmt.Errorf("jira issue status mapping already exists")
	}
	return err
}

/*
POST /plugins/jira/sources/:sourceId/type-mappings/:userType/status-mappings
{
	"userStatus": "user custom status",
	"standardStatus": "devlake standard status"
}
*/
func PostIssueStatusMappings(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	jiraIssueTypeMapping, err := findIssueTypeMappingByInputParam(input)
	if err != nil {
		return nil, err
	}
	jiraIssueStatusMapping := &models.JiraIssueStatusMapping{
		JiraSourceID: jiraIssueTypeMapping.JiraSourceID,
		UserType:     jiraIssueTypeMapping.UserType,
	}
	err = syncIssueStatusMappingFromInput(jiraIssueStatusMapping, input)
	if err != nil {
		return nil, err
	}
	// save
	err = wrapIssueStatusDuplicateErr(lakeModels.Db.Create(jiraIssueStatusMapping).Error)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jiraIssueStatusMapping}, nil
}

/*
PUT /plugins/jira/sources/:sourceId/type-mappings/:userType/status-mappings/:userStatus
{
	"standardStatus": "devlake standard status"
}
*/
func PutIssueStatusMapping(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// load from db
	jiraIssueStatusMapping, err := findIssueStatusMappingFromInput(input)
	if err != nil {
		return nil, err
	}
	// update with request body
	err = syncIssueStatusMappingFromInput(jiraIssueStatusMapping, input)
	if err != nil {
		return nil, err
	}
	// save
	err = wrapIssueStatusDuplicateErr(lakeModels.Db.Save(jiraIssueStatusMapping).Error)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jiraIssueStatusMapping}, nil
}

/*
DELETE /plugins/jira/sources/:sourceId/type-mappings/:userType/status-mappings/:userStatus
*/
func DeleteIssueStatusMapping(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	jiraIssueStatusMapping, err := findIssueStatusMappingFromInput(input)
	if err != nil {
		return nil, err
	}
	err = lakeModels.Db.Delete(jiraIssueStatusMapping).Error
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jiraIssueStatusMapping}, nil
}

/*
GET /plugins/jira/sources/:sourceId/type-mappings/:userType/status-mappings
*/
func ListIssueStatusMappings(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	jiraIssueTypeMapping, err := findIssueTypeMappingByInputParam(input)
	if err != nil {
		return nil, err
	}
	jiraIssueStatusMappings := make([]models.JiraIssueStatusMapping, 0)
	err = lakeModels.Db.Where(
		"jira_source_id = ? AND user_type = ?",
		jiraIssueTypeMapping.JiraSourceID,
		jiraIssueTypeMapping.UserType,
	).Find(&jiraIssueStatusMappings).Error
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jiraIssueStatusMappings}, nil
}
