package api

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"github.com/mitchellh/mapstructure"
)

func findIssueTypeMappingByInputParam(input *core.ApiResourceInput) (*models.JiraIssueTypeMapping, error) {
	// load jira source
	jiraSource, err := findSourceByInputParam(input)
	if err != nil {
		return nil, err
	}
	// load jira type mapping from db
	userType := input.Params["userType"]
	if userType == "" {
		return nil, fmt.Errorf("missing userType")
	}
	jiraIssueTypeMapping := &models.JiraIssueTypeMapping{}
	err = lakeModels.Db.First(jiraIssueTypeMapping, jiraSource.ID, userType).Error
	if err != nil {
		return nil, err
	}

	return jiraIssueTypeMapping, nil
}

func syncIssueTypeMappingFromInput(
	jiraIssueTypeMapping *models.JiraIssueTypeMapping,
	input *core.ApiResourceInput,
) error {
	// decode
	err := mapstructure.Decode(input.Body, jiraIssueTypeMapping)
	if err != nil {
		return err
	}
	// validate
	vld := validator.New()
	err = vld.Struct(jiraIssueTypeMapping)
	if err != nil {
		return err
	}
	return nil
}

func wrapIssueTypeDuplicateErr(err error) error {
	if lakeModels.IsDuplicateError(err) {
		return fmt.Errorf("jira issue type mapping already exists")
	}
	return err
}

/*
POST /plugins/jira/sources/:sourceId/type-mappings
{
	"userType": "user custom type",
	"standardType": "devlake standard type"
}
*/
func PostIssueTypeMappings(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// create new
	jiraSource, err := findSourceByInputParam(input)
	if err != nil {
		return nil, err
	}
	jiraIssueTypeMapping := &models.JiraIssueTypeMapping{JiraSourceID: jiraSource.ID}
	err = syncIssueTypeMappingFromInput(jiraIssueTypeMapping, input)
	if err != nil {
		return nil, err
	}
	// save
	err = wrapIssueTypeDuplicateErr(lakeModels.Db.Create(jiraIssueTypeMapping).Error)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jiraIssueTypeMapping}, nil
}

/*
PUT /plugins/jira/sources/:sourceId/type-mappings/:userType
{
	"standardType": "devlake standard type"
}
*/
func PutIssueTypeMapping(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// load from db
	jiraIssueTypeMapping, err := findIssueTypeMappingByInputParam(input)
	if err != nil {
		return nil, err
	}
	// update with request body
	err = syncIssueTypeMappingFromInput(jiraIssueTypeMapping, input)
	if err != nil {
		return nil, err
	}
	// save
	err = wrapIssueTypeDuplicateErr(lakeModels.Db.Save(jiraIssueTypeMapping).Error)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jiraIssueTypeMapping}, nil
}

/*
DELETE /plugins/jira/sources/:sourceId/type-mappings/:userType
*/
func DeleteIssueTypeMapping(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// load from db
	jiraIssueTypeMapping, err := findIssueTypeMappingByInputParam(input)
	if err != nil {
		return nil, err
	}
	err = lakeModels.Db.Delete(jiraIssueTypeMapping).Error
	if err != nil {
		return nil, err
	}
	err = lakeModels.Db.Where("jira_source_id = ? AND user_type = ?", jiraIssueTypeMapping.JiraSourceID, jiraIssueTypeMapping.UserType).Delete(&models.JiraIssueStatusMapping{}).Error
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jiraIssueTypeMapping}, nil
}

/*
GET /plugins/jira/sources/:sourceId/type-mappings
*/
func ListIssueTypeMappings(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	jiraSource, err := findSourceByInputParam(input)
	if err != nil {
		return nil, err
	}
	jiraIssueTypeMappings := make([]models.JiraIssueTypeMapping, 0)
	err = lakeModels.Db.Where("jira_source_id = ?", jiraSource.ID).Find(&jiraIssueTypeMappings).Error
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jiraIssueTypeMappings}, nil
}
