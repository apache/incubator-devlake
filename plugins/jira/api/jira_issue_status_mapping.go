package api

import (
	"fmt"
	"github.com/merico-dev/lake/models/common"
	"net/http"

	"github.com/go-playground/validator/v10"
	lakeModels "github.com/merico-dev/lake/models"
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
	err = lakeModels.Db.First(
		jiraIssueStatusMapping,
		jiraIssueTypeMapping.SourceID,
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
	sources ...map[string]interface{},
) error {
	// merge fields from sources to jiraIssueStatusMapping
	for _, source := range sources {
		err := mapstructure.Decode(source, jiraIssueStatusMapping)
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

func saveStatusMappings(tx *gorm.DB, jiraSourceId uint64, userType string, statusMappings interface{}) error {
	statusMappingsMap, ok := statusMappings.(map[string]interface{})
	if !ok {
		return fmt.Errorf("statusMappings is not a JSON object: %v", statusMappings)
	}
	err := tx.Where(
		"source_id = ? AND user_type = ?",
		jiraSourceId,
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
			"SourceID":   jiraSourceId,
			"UserType":   userType,
			"UserStatus": userStatus,
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

func findIssueStatusMappingBySourceIdAndUserType(
	jiraSourceId uint64,
	userType string,
) ([]*models.JiraIssueStatusMapping, error) {
	jiraIssueStatusMappings := make([]*models.JiraIssueStatusMapping, 0)
	err := lakeModels.Db.Where(
		"source_id = ? AND user_type = ?",
		jiraSourceId,
		userType,
	).Find(&jiraIssueStatusMappings).Error
	return jiraIssueStatusMappings, err
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
	jiraIssueStatusMapping := &models.JiraIssueStatusMapping{}
	err = mergeFieldsToJiraStatusMapping(jiraIssueStatusMapping, input.Body, map[string]interface{}{
		"SourceID": jiraIssueTypeMapping.SourceID,
		"UserType": jiraIssueTypeMapping.UserType,
	})
	if err != nil {
		return nil, err
	}
	// save
	err = wrapIssueStatusDuplicateErr(lakeModels.Db.Create(jiraIssueStatusMapping).Error)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jiraIssueStatusMapping, Status: http.StatusCreated}, nil
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
	err = mergeFieldsToJiraStatusMapping(jiraIssueStatusMapping, input.Body)
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
	jiraIssueStatusMappings, err := findIssueStatusMappingBySourceIdAndUserType(
		jiraIssueTypeMapping.SourceID,
		jiraIssueTypeMapping.UserType,
	)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jiraIssueStatusMappings}, nil
}
