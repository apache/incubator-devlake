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

func findIssueTypeMappingByInputParam(input *core.ApiResourceInput) (*models.TapdIssueTypeMapping, error) {
	// load tapd source
	tapdSource, err := findSourceByInputParam(input)
	if err != nil {
		return nil, err
	}
	// load tapd type mapping from db
	userType := input.Params["userType"]
	if userType == "" {
		return nil, fmt.Errorf("missing userType")
	}
	tapdIssueTypeMapping := &models.TapdIssueTypeMapping{}
	err = db.First(tapdIssueTypeMapping, tapdSource.ID, userType).Error
	if err != nil {
		return nil, err
	}

	return tapdIssueTypeMapping, nil
}

func mergeFieldsToTapdTypeMapping(
	tapdIssueTypeMapping *models.TapdIssueTypeMapping,
	sources ...map[string]interface{},
) error {
	// merge fields from sources to tapdIssueTypeMapping
	for _, source := range sources {
		err := mapstructure.Decode(source, tapdIssueTypeMapping)
		if err != nil {
			return err
		}
	}
	// validate
	vld := validator.New()
	err := vld.Struct(tapdIssueTypeMapping)
	if err != nil {
		return err
	}
	return nil
}

func wrapIssueTypeDuplicateErr(err error) error {
	if common.IsDuplicateError(err) {
		return fmt.Errorf("tapd issue type mapping already exists")
	}
	return err
}

func saveTypeMappings(tx *gorm.DB, tapdSourceId uint64, typeMappings interface{}) error {
	typeMappingsMap, ok := typeMappings.(map[string]interface{})
	if !ok {
		return fmt.Errorf("typeMappings is not a JSON object: %v", typeMappings)
	}
	err := tx.Where("source_id = ?", tapdSourceId).Delete(&models.TapdIssueTypeMapping{}).Error
	if err != nil {
		return err
	}
	for userType, typeMapping := range typeMappingsMap {
		typeMappingMap, ok := typeMapping.(map[string]interface{})
		if !ok {
			return fmt.Errorf("typeMapping is not a JSON object: %v", typeMapping)
		}
		tapdIssueTypeMapping := &models.TapdIssueTypeMapping{}
		err = mergeFieldsToTapdTypeMapping(tapdIssueTypeMapping, typeMappingMap, map[string]interface{}{
			"SourceID": tapdSourceId,
			"UserType": userType,
		})
		if err != nil {
			return err
		}
		err = wrapIssueTypeDuplicateErr(tx.Create(tapdIssueTypeMapping).Error)
		if err != nil {
			return err
		}

		statusMappings := typeMappingMap["statusMappings"]
		if statusMappings != nil {
			err = saveStatusMappings(tx, tapdSourceId, userType, statusMappings)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func findIssueTypeMappingBySourceId(tapdSourceId uint64) ([]*models.TapdIssueTypeMapping, error) {
	tapdIssueTypeMappings := make([]*models.TapdIssueTypeMapping, 0)
	err := db.Where("source_id = ?", tapdSourceId).Find(&tapdIssueTypeMappings).Error
	if err != nil {
		return nil, err
	}
	return tapdIssueTypeMappings, nil
}

/*
POST /plugins/tapd/sources/:sourceId/type-mappings
{
	"userType": "user custom type",
	"standardType": "devlake standard type"
}
*/
func PostIssueTypeMappings(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// create new
	tapdSource, err := findSourceByInputParam(input)
	if err != nil {
		return nil, err
	}
	tapdIssueTypeMapping := &models.TapdIssueTypeMapping{}
	err = mergeFieldsToTapdTypeMapping(tapdIssueTypeMapping, input.Body, map[string]interface{}{
		"SourceID": tapdSource.ID,
	})
	if err != nil {
		return nil, err
	}
	// save
	err = wrapIssueTypeDuplicateErr(db.Create(tapdIssueTypeMapping).Error)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: tapdIssueTypeMapping, Status: http.StatusCreated}, nil
}

/*
PUT /plugins/tapd/sources/:sourceId/type-mappings/:userType
{
	"standardType": "devlake standard type"
}
*/
func PutIssueTypeMapping(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// load from db
	tapdIssueTypeMapping, err := findIssueTypeMappingByInputParam(input)
	if err != nil {
		return nil, err
	}
	// update with request body
	err = mergeFieldsToTapdTypeMapping(tapdIssueTypeMapping, input.Body)
	if err != nil {
		return nil, err
	}
	// save
	err = wrapIssueTypeDuplicateErr(db.Save(tapdIssueTypeMapping).Error)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: tapdIssueTypeMapping}, nil
}

/*
DELETE /plugins/tapd/sources/:sourceId/type-mappings/:userType
*/
func DeleteIssueTypeMapping(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// load from db
	tapdIssueTypeMapping, err := findIssueTypeMappingByInputParam(input)
	if err != nil {
		return nil, err
	}
	err = db.Delete(tapdIssueTypeMapping).Error
	if err != nil {
		return nil, err
	}
	err = db.Where(
		"source_id = ? AND user_type = ?",
		tapdIssueTypeMapping.SourceID,
		tapdIssueTypeMapping.UserType,
	).Delete(&models.TapdIssueStatusMapping{}).Error
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: tapdIssueTypeMapping}, nil
}

/*
GET /plugins/tapd/sources/:sourceId/type-mappings
*/
func ListIssueTypeMappings(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	tapdSource, err := findSourceByInputParam(input)
	if err != nil {
		return nil, err
	}
	tapdIssueTypeMappings, err := findIssueTypeMappingBySourceId(tapdSource.ID)
	return &core.ApiResourceOutput{Body: tapdIssueTypeMappings}, err
}
