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

func findIssueTypeMappingByInputParam(input *core.ApiResourceInput) (*models.JiraIssueTypeMapping, error) {
	// load jira connection
	jiraConnection, err := findConnectionByInputParam(input)
	if err != nil {
		return nil, err
	}
	// load jira type mapping from db
	userType := input.Params["userType"]
	if userType == "" {
		return nil, fmt.Errorf("missing userType")
	}
	jiraIssueTypeMapping := &models.JiraIssueTypeMapping{}
	err = db.First(jiraIssueTypeMapping, jiraConnection.ID, userType).Error
	if err != nil {
		return nil, err
	}

	return jiraIssueTypeMapping, nil
}

func mergeFieldsToJiraTypeMapping(
	jiraIssueTypeMapping *models.JiraIssueTypeMapping,
	connections ...map[string]interface{},
) error {
	// merge fields from connections to jiraIssueTypeMapping
	for _, connection := range connections {
		err := mapstructure.Decode(connection, jiraIssueTypeMapping)
		if err != nil {
			return err
		}
	}
	// validate
	vld := validator.New()
	err := vld.Struct(jiraIssueTypeMapping)
	if err != nil {
		return err
	}
	return nil
}

func wrapIssueTypeDuplicateErr(err error) error {
	if common.IsDuplicateError(err) {
		return fmt.Errorf("jira issue type mapping already exists")
	}
	return err
}

func saveTypeMappings(tx *gorm.DB, jiraConnectionId uint64, typeMappings interface{}) error {
	typeMappingsMap, ok := typeMappings.(map[string]interface{})
	if !ok {
		return fmt.Errorf("typeMappings is not a JSON object: %v", typeMappings)
	}
	err := tx.Where("connection_id = ?", jiraConnectionId).Delete(&models.JiraIssueTypeMapping{}).Error
	if err != nil {
		return err
	}
	for userType, typeMapping := range typeMappingsMap {
		typeMappingMap, ok := typeMapping.(map[string]interface{})
		if !ok {
			return fmt.Errorf("typeMapping is not a JSON object: %v", typeMapping)
		}
		jiraIssueTypeMapping := &models.JiraIssueTypeMapping{}
		err = mergeFieldsToJiraTypeMapping(jiraIssueTypeMapping, typeMappingMap, map[string]interface{}{
			"ConnectionID": jiraConnectionId,
			"UserType":     userType,
		})
		if err != nil {
			return err
		}
		err = wrapIssueTypeDuplicateErr(tx.Create(jiraIssueTypeMapping).Error)
		if err != nil {
			return err
		}

		statusMappings := typeMappingMap["statusMappings"]
		if statusMappings != nil {
			err = saveStatusMappings(tx, jiraConnectionId, userType, statusMappings)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func findIssueTypeMappingByConnectionId(jiraConnectionId uint64) ([]*models.JiraIssueTypeMapping, error) {
	jiraIssueTypeMappings := make([]*models.JiraIssueTypeMapping, 0)
	err := db.Where("connection_id = ?", jiraConnectionId).Find(&jiraIssueTypeMappings).Error
	if err != nil {
		return nil, err
	}
	return jiraIssueTypeMappings, nil
}

/*
POST /plugins/jira/connections/:connectionId/type-mappings
{
	"userType": "user custom type",
	"standardType": "devlake standard type"
}
*/
func PostIssueTypeMappings(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// create new
	jiraConnection, err := findConnectionByInputParam(input)
	if err != nil {
		return nil, err
	}
	jiraIssueTypeMapping := &models.JiraIssueTypeMapping{}
	err = mergeFieldsToJiraTypeMapping(jiraIssueTypeMapping, input.Body, map[string]interface{}{
		"ConnectionID": jiraConnection.ID,
	})
	if err != nil {
		return nil, err
	}
	// save
	err = wrapIssueTypeDuplicateErr(db.Create(jiraIssueTypeMapping).Error)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jiraIssueTypeMapping, Status: http.StatusCreated}, nil
}

/*
PUT /plugins/jira/connections/:connectionId/type-mappings/:userType
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
	err = mergeFieldsToJiraTypeMapping(jiraIssueTypeMapping, input.Body)
	if err != nil {
		return nil, err
	}
	// save
	err = wrapIssueTypeDuplicateErr(db.Save(jiraIssueTypeMapping).Error)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jiraIssueTypeMapping}, nil
}

/*
DELETE /plugins/jira/connections/:connectionId/type-mappings/:userType
*/
func DeleteIssueTypeMapping(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// load from db
	jiraIssueTypeMapping, err := findIssueTypeMappingByInputParam(input)
	if err != nil {
		return nil, err
	}
	err = db.Delete(jiraIssueTypeMapping).Error
	if err != nil {
		return nil, err
	}
	err = db.Where(
		"connection_id = ? AND user_type = ?",
		jiraIssueTypeMapping.ConnectionID,
		jiraIssueTypeMapping.UserType,
	).Delete(&models.JiraIssueStatusMapping{}).Error
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jiraIssueTypeMapping}, nil
}

/*
GET /plugins/jira/connections/:connectionId/type-mappings
*/
func ListIssueTypeMappings(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	jiraConnection, err := findConnectionByInputParam(input)
	if err != nil {
		return nil, err
	}
	jiraIssueTypeMappings, err := findIssueTypeMappingByConnectionId(jiraConnection.ID)
	return &core.ApiResourceOutput{Body: jiraIssueTypeMappings}, err
}
