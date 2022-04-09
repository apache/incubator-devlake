package api

import (
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/models/common"

	"github.com/go-playground/validator/v10"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

func findIssueStatusMappingFromInput(input *core.ApiResourceInput) (*models.TapdIssueStatusMapping, error) {
	// load type mapping
	tapdIssueTypeMapping, err := findIssueTypeMappingByInputParam(input)
	if err != nil {
		return nil, err
	}
	// load status mapping from db
	userStatus := input.Params["userStatus"]
	if userStatus == "" {
		return nil, fmt.Errorf("missing userStatus")
	}
	tapdIssueStatusMapping := &models.TapdIssueStatusMapping{}
	err = db.First(
		tapdIssueStatusMapping,
		tapdIssueTypeMapping.SourceID,
		tapdIssueTypeMapping.UserType,
		userStatus,
	).Error
	if err != nil {
		return nil, err
	}

	return tapdIssueStatusMapping, nil
}

func mergeFieldsToTapdStatusMapping(
	tapdIssueStatusMapping *models.TapdIssueStatusMapping,
	sources ...map[string]interface{},
) error {
	// merge fields from sources to tapdIssueStatusMapping
	for _, source := range sources {
		err := mapstructure.Decode(source, tapdIssueStatusMapping)
		if err != nil {
			return err
		}
	}
	// validate
	vld := validator.New()
	err := vld.Struct(tapdIssueStatusMapping)
	if err != nil {
		return err
	}
	return nil
}

func wrapIssueStatusDuplicateErr(err error) error {
	if common.IsDuplicateError(err) {
		return fmt.Errorf("tapd issue status mapping already exists")
	}
	return err
}

func saveStatusMappings(tx *gorm.DB, tapdSourceId uint64, userType string, statusMappings interface{}) error {
	statusMappingsMap, ok := statusMappings.(map[string]interface{})
	if !ok {
		return fmt.Errorf("statusMappings is not a JSON object: %v", statusMappings)
	}
	err := tx.Where(
		"source_id = ? AND user_type = ?",
		tapdSourceId,
		userType).Delete(&models.TapdIssueStatusMapping{}).Error
	if err != nil {
		return err
	}
	for userStatus, statusMapping := range statusMappingsMap {
		statusMappingMap, ok := statusMapping.(map[string]interface{})
		if !ok {
			return fmt.Errorf("statusMapping is not a JSON object: %v", statusMappings)
		}
		tapdIssueStatusMapping := &models.TapdIssueStatusMapping{}
		err = mergeFieldsToTapdStatusMapping(tapdIssueStatusMapping, statusMappingMap, map[string]interface{}{
			"SourceID":   tapdSourceId,
			"UserType":   userType,
			"UserStatus": userStatus,
		})
		if err != nil {
			return err
		}
		err = tx.Create(tapdIssueStatusMapping).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func findIssueStatusMappingBySourceIdAndUserType(
	tapdSourceId uint64,
	userType string,
) ([]*models.TapdIssueStatusMapping, error) {
	tapdIssueStatusMappings := make([]*models.TapdIssueStatusMapping, 0)
	err := db.Where(
		"source_id = ? AND user_type = ?",
		tapdSourceId,
		userType,
	).Find(&tapdIssueStatusMappings).Error
	return tapdIssueStatusMappings, err
}

/*
POST /plugins/tapd/sources/:sourceId/type-mappings/:userType/status-mappings
{
	"userStatus": "user custom status",
	"standardStatus": "devlake standard status"
}
*/
func PostIssueStatusMappings(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	tapdIssueTypeMapping, err := findIssueTypeMappingByInputParam(input)
	if err != nil {
		return nil, err
	}
	tapdIssueStatusMapping := &models.TapdIssueStatusMapping{}
	err = mergeFieldsToTapdStatusMapping(tapdIssueStatusMapping, input.Body, map[string]interface{}{
		"SourceID": tapdIssueTypeMapping.SourceID,
		"UserType": tapdIssueTypeMapping.UserType,
	})
	if err != nil {
		return nil, err
	}
	// save
	err = wrapIssueStatusDuplicateErr(db.Create(tapdIssueStatusMapping).Error)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: tapdIssueStatusMapping, Status: http.StatusCreated}, nil
}

/*
PUT /plugins/tapd/sources/:sourceId/type-mappings/:userType/status-mappings/:userStatus
{
	"standardStatus": "devlake standard status"
}
*/
func PutIssueStatusMapping(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// load from db
	tapdIssueStatusMapping, err := findIssueStatusMappingFromInput(input)
	if err != nil {
		return nil, err
	}
	// update with request body
	err = mergeFieldsToTapdStatusMapping(tapdIssueStatusMapping, input.Body)
	if err != nil {
		return nil, err
	}
	// save
	err = wrapIssueStatusDuplicateErr(db.Save(tapdIssueStatusMapping).Error)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: tapdIssueStatusMapping}, nil
}

/*
DELETE /plugins/tapd/sources/:sourceId/type-mappings/:userType/status-mappings/:userStatus
*/
func DeleteIssueStatusMapping(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	tapdIssueStatusMapping, err := findIssueStatusMappingFromInput(input)
	if err != nil {
		return nil, err
	}
	err = db.Delete(tapdIssueStatusMapping).Error
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: tapdIssueStatusMapping}, nil
}

/*
GET /plugins/tapd/sources/:sourceId/type-mappings/:userType/status-mappings
*/
func ListIssueStatusMappings(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	tapdIssueTypeMapping, err := findIssueTypeMappingByInputParam(input)
	if err != nil {
		return nil, err
	}
	tapdIssueStatusMappings, err := findIssueStatusMappingBySourceIdAndUserType(
		tapdIssueTypeMapping.SourceID,
		tapdIssueTypeMapping.UserType,
	)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: tapdIssueStatusMappings}, nil
}
