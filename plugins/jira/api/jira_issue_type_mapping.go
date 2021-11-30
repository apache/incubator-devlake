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

func mergeFieldsToJiraTypeMapping(
	jiraIssueTypeMapping *models.JiraIssueTypeMapping,
	sources ...map[string]interface{},
) error {
	// merge fields from sources to jiraIssueTypeMapping
	for _, source := range sources {
		err := mapstructure.Decode(source, jiraIssueTypeMapping)
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

func saveTypeMappings(tx *gorm.DB, jiraSourceId uint64, typeMappings interface{}) error {
	typeMappingsMap, ok := typeMappings.(map[string]interface{})
	if !ok {
		return fmt.Errorf("typeMappings is not a JSON object: %v", typeMappings)
	}
	err := tx.Where("source_id = ?", jiraSourceId).Delete(&models.JiraIssueTypeMapping{}).Error
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
			"SourceID": jiraSourceId,
			"UserType": userType,
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
			err = saveStatusMappings(tx, jiraSourceId, userType, statusMappings)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func findIssueTypeMappingBySourceId(jiraSourceId uint64) ([]*models.JiraIssueTypeMapping, error) {
	jiraIssueTypeMappings := make([]*models.JiraIssueTypeMapping, 0)
	err := lakeModels.Db.Where("source_id = ?", jiraSourceId).Find(&jiraIssueTypeMappings).Error
	if err != nil {
		return nil, err
	}
	return jiraIssueTypeMappings, nil
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
	jiraIssueTypeMapping := &models.JiraIssueTypeMapping{}
	err = mergeFieldsToJiraTypeMapping(jiraIssueTypeMapping, input.Body, map[string]interface{}{
		"SourceID": jiraSource.ID,
	})
	if err != nil {
		return nil, err
	}
	// save
	err = wrapIssueTypeDuplicateErr(lakeModels.Db.Create(jiraIssueTypeMapping).Error)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: jiraIssueTypeMapping, Status: http.StatusCreated}, nil
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
	err = mergeFieldsToJiraTypeMapping(jiraIssueTypeMapping, input.Body)
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
	err = lakeModels.Db.Where(
		"source_id = ? AND user_type = ?",
		jiraIssueTypeMapping.SourceID,
		jiraIssueTypeMapping.UserType,
	).Delete(&models.JiraIssueStatusMapping{}).Error
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
	jiraIssueTypeMappings, err := findIssueTypeMappingBySourceId(jiraSource.ID)
	return &core.ApiResourceOutput{Body: jiraIssueTypeMappings}, err
}
