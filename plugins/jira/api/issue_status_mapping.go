package api

import (
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/models/common"

	"github.com/go-playground/validator/v10"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
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
	err = db.First(
		jiraIssueStatusMapping,
		jiraIssueTypeMapping.ConnectionID,
		jiraIssueTypeMapping.UserType,
		userStatus,
	).Error
	if err != nil {
		return nil, err
	}

	return jiraIssueStatusMapping, nil
}

func mergeFieldsToJiraStatusMapping(
	jiraIssueStatusMapping *models.JiraIssueStatusMapping,
	connections ...map[string]interface{},
) error {
	// merge fields from connections to jiraIssueStatusMapping
	for _, connection := range connections {
		err := mapstructure.Decode(connection, jiraIssueStatusMapping)
		if err != nil {
			return err
		}
	}
	// validate
	vld := validator.New()
	err := vld.Struct(jiraIssueStatusMapping)
	if err != nil {
		return err
	}
	return nil
}

func wrapIssueStatusDuplicateErr(err error) error {
	if common.IsDuplicateError(err) {
		return fmt.Errorf("jira issue status mapping already exists")
	}
	return err
}

func saveStatusMappings(tx *gorm.DB, jiraConnectionId uint64, userType string, statusMappings interface{}) error {
	statusMappingsMap, ok := statusMappings.(map[string]interface{})
	if !ok {
		return fmt.Errorf("statusMappings is not a JSON object: %v", statusMappings)
	}
	err := tx.Where(
		"connection_id = ? AND user_type = ?",
		jiraConnectionId,
		userType).Delete(&models.JiraIssueStatusMapping{}).Error
	if err != nil {
		return err
	}
	for userStatus, statusMapping := range statusMappingsMap {
		statusMappingMap, ok := statusMapping.(map[string]interface{})
		if !ok {
			return fmt.Errorf("statusMapping is not a JSON object: %v", statusMappings)
		}
		jiraIssueStatusMapping := &models.JiraIssueStatusMapping{}
		err = mergeFieldsToJiraStatusMapping(jiraIssueStatusMapping, statusMappingMap, map[string]interface{}{
			"ConnectionID": jiraConnectionId,
			"UserType":     userType,
			"UserStatus":   userStatus,
		})
		if err != nil {
			return err
		}
		err = tx.Create(jiraIssueStatusMapping).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func findIssueStatusMappingByConnectionIdAndUserType(
	jiraConnectionId uint64,
	userType string,
) ([]*models.JiraIssueStatusMapping, error) {
	jiraIssueStatusMappings := make([]*models.JiraIssueStatusMapping, 0)
	err := db.Where(
		"connection_id = ? AND user_type = ?",
		jiraConnectionId,
		userType,
	).Find(&jiraIssueStatusMappings).Error
	return jiraIssueStatusMappings, err
}

/*
POST /plugins/jira/connections/:connectionId/type-mappings/:userType/status-mappings
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
	jiraIssueStatusMapping := &models.JiraIssueStatusMapping{}
	err = mergeFieldsToJiraStatusMapping(jiraIssueStatusMapping, input.Body, map[string]interface{}{
		"ConnectionID": jiraIssueTypeMapping.ConnectionID,
		"UserType":     jiraIssueTypeMapping.UserType,
	})
	if err != nil {
		return nil, err
	}
	// save
	err = wrapIssueStatusDuplicateErr(db.Create(jiraIssueStatusMapping).Error)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jiraIssueStatusMapping, Status: http.StatusCreated}, nil
}

/*
PUT /plugins/jira/connections/:connectionId/type-mappings/:userType/status-mappings/:userStatus
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
	err = mergeFieldsToJiraStatusMapping(jiraIssueStatusMapping, input.Body)
	if err != nil {
		return nil, err
	}
	// save
	err = wrapIssueStatusDuplicateErr(db.Save(jiraIssueStatusMapping).Error)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jiraIssueStatusMapping}, nil
}

/*
DELETE /plugins/jira/connections/:connectionId/type-mappings/:userType/status-mappings/:userStatus
*/
func DeleteIssueStatusMapping(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	jiraIssueStatusMapping, err := findIssueStatusMappingFromInput(input)
	if err != nil {
		return nil, err
	}
	err = db.Delete(jiraIssueStatusMapping).Error
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jiraIssueStatusMapping}, nil
}

/*
GET /plugins/jira/connections/:connectionId/type-mappings/:userType/status-mappings
*/
func ListIssueStatusMappings(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	jiraIssueTypeMapping, err := findIssueTypeMappingByInputParam(input)
	if err != nil {
		return nil, err
	}
	jiraIssueStatusMappings, err := findIssueStatusMappingByConnectionIdAndUserType(
		jiraIssueTypeMapping.ConnectionID,
		jiraIssueTypeMapping.UserType,
	)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jiraIssueStatusMappings}, nil
}
