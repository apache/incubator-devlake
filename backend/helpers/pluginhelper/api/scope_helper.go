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

package api

import (
	"fmt"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"

	"reflect"
)

// ScopeApiHelper is used to write the CURD of connection
type ScopeApiHelper struct {
	log        log.Logger
	db         dal.Dal
	validator  *validator.Validate
	connHelper *ConnectionApiHelper
}

// NewScopeHelper creates a ScopeHelper for connection management
func NewScopeHelper(
	basicRes context.BasicRes,
	vld *validator.Validate,
	connHelper *ConnectionApiHelper,
) *ScopeApiHelper {
	if vld == nil {
		vld = validator.New()
	}
	if connHelper == nil {
		return nil
	}
	return &ScopeApiHelper{
		log:        basicRes.GetLogger(),
		db:         basicRes.GetDal(),
		validator:  vld,
		connHelper: connHelper,
	}
}

// Put saves the given scopes to the database. It expects a slice of struct pointers
// as the scopes argument. It also expects a fieldName argument, which is used to extract
// the connection ID from the input.Params map.
func (c *ScopeApiHelper) Put(input *plugin.ApiResourceInput, apiScope interface{}, connection interface{}) errors.Error {
	err := errors.Convert(mapstructure.Decode(input.Body, apiScope))
	if err != nil {
		return errors.BadInput.Wrap(err, "decoding Github repo error")
	}
	// Ensure that the scopes argument is a slice
	v := reflect.ValueOf(apiScope)
	scopesValue := v.Elem().FieldByName("Data")
	if scopesValue.Kind() != reflect.Slice {
		panic("expected a slice")
	}
	// Extract the connection ID from the input.Params map
	connectionId, _ := ExtractParam(input.Params)
	if connectionId == 0 {
		return errors.BadInput.New("invalid connectionId or scopeId")
	}
	err = c.VerifyConnection(connection, connectionId)
	if err != nil {
		return err
	}
	// Create a map to keep track of primary key values
	keeper := make(map[string]struct{})

	// Set the CreatedDate and UpdatedDate fields to the current time for each scope
	now := time.Now()
	for i := 0; i < scopesValue.Len(); i++ {
		// Get the reflect.Value of the i-th struct pointer in the slice
		structValue := scopesValue.Index(i)

		// Ensure that the structValue is a pointer to a struct
		if structValue.Kind() != reflect.Ptr || structValue.Elem().Kind() != reflect.Struct {
			panic("expected a pointer to a struct")
		}

		// Ensure that the primary key value is unique
		primaryValueStr := ReturnPrimaryKeyValue(structValue.Elem().Interface())
		if _, ok := keeper[primaryValueStr]; ok {
			return errors.BadInput.New("duplicated item")
		} else {
			keeper[primaryValueStr] = struct{}{}
		}

		// Set the connection ID, CreatedDate, and UpdatedDate fields
		SetScopeFields(structValue.Interface(), connectionId, &now, &now)

		// Verify that the primary key value is valid
		err = VerifyPrimaryKeyValue(structValue.Elem().Interface())
		if err != nil {
			return err
		}
	}

	// Save the scopes to the database
	return c.save(scopesValue.Interface(), c.db.Create)
}

func (c *ScopeApiHelper) Update(input *plugin.ApiResourceInput, fieldName string, connection interface{}, scope interface{}) errors.Error {
	connectionId, scopeId := ExtractParam(input.Params)

	if connectionId == 0 || len(scopeId) == 0 || scopeId == "0" {
		return errors.BadInput.New("invalid connectionId")
	}
	err := c.VerifyConnection(connection, connectionId)
	if err != nil {
		return err
	}

	err = c.db.First(scope, dal.Where(fmt.Sprintf("connection_id = ? AND %s = ?", fieldName), connectionId, scopeId))
	if err != nil {
		return errors.Default.New("getting Scope error")
	}
	err = DecodeMapStruct(input.Body, scope)
	if err != nil {
		return errors.Default.Wrap(err, "patch scope error")
	}
	err = VerifyPrimaryKeyValue(scope)
	if err != nil {
		return err
	}
	err = c.db.Update(scope)
	if err != nil {
		return errors.Default.Wrap(err, "error on saving Scope")
	}
	return nil
}

func (c *ScopeApiHelper) GetScopeList(input *plugin.ApiResourceInput, connection interface{}, scopes interface{}, rules interface{}) (map[uint64]string, errors.Error) {
	connectionId, _ := ExtractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid path params")
	}
	err := c.VerifyConnection(connection, connectionId)
	if err != nil {
		return nil, err
	}
	limit, offset := GetLimitOffset(input.Query, "pageSize", "page")
	err = c.db.All(scopes, dal.Where("connection_id = ?", connectionId), dal.Limit(limit), dal.Offset(offset))
	if err != nil {
		return nil, err
	}

	scopesValue := reflect.ValueOf(reflect.ValueOf(scopes).Elem().Interface())
	if scopesValue.Kind() != reflect.Slice {
		panic("expected a slice")
	}
	var ruleIds []uint64
	for i := 0; i < scopesValue.Len(); i++ {
		// Get the reflect.Value of the i-th struct pointer in the slice
		structValue := scopesValue.Index(i)

		// Ensure that the structValue is a pointer to a struct
		if structValue.Kind() != reflect.Ptr || structValue.Elem().Kind() != reflect.Struct {
			panic("expected a pointer to a struct")
		}
		ruleId := structValue.Elem().FieldByName("TransformationRuleId").Uint()
		if ruleId > 0 {
			ruleIds = append(ruleIds, ruleId)
		}
	}

	if len(ruleIds) > 0 {
		err = c.db.All(rules, dal.Where("id IN (?)", ruleIds))
		if err != nil {
			return nil, err
		}
	}
	rulesValue := reflect.ValueOf(reflect.ValueOf(rules).Elem().Interface())
	if scopesValue.Kind() != reflect.Slice {
		panic("expected a slice")
	}
	names := make(map[uint64]string)
	for i := 0; i < rulesValue.Len(); i++ {
		// Get the reflect.Value of the i-th struct pointer in the slice
		structValue := rulesValue.Index(i)
		names[structValue.FieldByName("ID").Uint()] = structValue.FieldByName("Name").String()
	}
	return names, nil
}

func (c *ScopeApiHelper) GetScope(input *plugin.ApiResourceInput, fieldName string, connection interface{}, scope interface{}, rule interface{}) errors.Error {
	connectionId, scopeId := ExtractParam(input.Params)
	if connectionId == 0 || len(scopeId) == 0 || scopeId == "0" {
		return errors.BadInput.New("invalid path params")
	}
	err := c.VerifyConnection(connection, connectionId)
	if err != nil {
		return err
	}
	db := c.db
	query := dal.Where(fmt.Sprintf("connection_id = ? AND %s = ?", fieldName), connectionId, scopeId)
	err = db.First(scope, query)
	if db.IsErrorNotFound(err) {
		return errors.NotFound.New("Scope not found")
	}
	if err != nil {
		return err
	}
	repoRuleId := reflect.ValueOf(scope).Elem().FieldByName("TransformationRuleId").Uint()
	if repoRuleId > 0 {
		err = db.First(rule, dal.Where("id = ?", repoRuleId))
		if err != nil {
			return errors.NotFound.New("transformationRule not found")
		}
	}
	return nil
}

func (c *ScopeApiHelper) VerifyConnection(connection interface{}, connId uint64) errors.Error {
	err := c.connHelper.FirstById(&connection, connId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.BadInput.New("Invalid Connection Id")
		}
		return err
	}
	return nil
}

func (c *ScopeApiHelper) save(scope interface{}, method func(entity interface{}, clauses ...dal.Clause) errors.Error) errors.Error {
	err := c.db.CreateOrUpdate(scope)
	if err != nil {
		if c.db.IsDuplicationError(err) {
			return errors.BadInput.New("the scope already exists")
		}
		return err
	}
	return nil
}

func ExtractParam(params map[string]string) (uint64, string) {
	connectionId, err := strconv.ParseUint(params["connectionId"], 10, 64)
	if err != nil {
		return 0, ""
	}
	scopeId := params["scopeId"]
	return connectionId, scopeId
}

// VerifyPrimaryKeyValue function verifies that the primary key value of a given struct instance is not zero or empty.
func VerifyPrimaryKeyValue(i interface{}) errors.Error {
	var value reflect.Value
	pType := reflect.TypeOf(i)
	if pType.Kind() == reflect.Ptr {
		value = reflect.ValueOf(reflect.ValueOf(i).Elem().Interface())
	} else {
		value = reflect.ValueOf(i)
	}
	// Loop through the fields of the input struct using reflection
	for j := 0; j < value.NumField(); j++ {
		field := value.Field(j)
		tag := value.Type().Field(j).Tag.Get("gorm")

		// Check if the field is tagged as a primary key using the GORM tag "primaryKey"
		if strings.Contains(tag, "primaryKey") {
			// If the field value is zero or nil, return an error indicating that the primary key value is invalid
			if field.Interface() == reflect.Zero(field.Type()).Interface() || field.Interface() == nil {
				return errors.Default.New("primary key value is zero or empty")
			}
		}
	}
	// If all primary key values are valid, return nil (no error)
	return nil
}

func SetScopeFields(p interface{}, connectionId uint64, createdDate *time.Time, updatedDate *time.Time) {
	pType := reflect.TypeOf(p)
	if pType.Kind() != reflect.Ptr {
		panic("expected a pointer to a struct")
	}
	pValue := reflect.ValueOf(p).Elem()

	// set connectionId
	connIdField := pValue.FieldByName("ConnectionId")
	connIdField.SetUint(connectionId)

	// set CreatedDate
	createdDateField := pValue.FieldByName("CreatedDate")
	createdDateField.Set(reflect.ValueOf(createdDate))

	// set UpdatedDate
	updatedDateField := pValue.FieldByName("UpdatedDate")
	if updatedDate == nil {
		// if updatedDate is nil, set UpdatedDate to be nil
		updatedDateField.Set(reflect.Zero(updatedDateField.Type()))
	} else {
		// if updatedDate is not nil, set UpdatedDate to be the value
		updatedDateFieldValue := reflect.ValueOf(updatedDate)
		updatedDateField.Set(updatedDateFieldValue)
	}
}

// ReturnPrimaryKeyValue returns a string containing the primary key value(s) of a struct, concatenated with "-" between them.
// This function receives an interface{} type argument p, which can be a pointer to any struct.
// The function uses reflection to iterate through the fields of the struct, and checks if each field is tagged as "primaryKey".
func ReturnPrimaryKeyValue(p interface{}) string {
	result := ""
	// get the type and value of the input interface using reflection
	t := reflect.TypeOf(p)
	v := reflect.ValueOf(p)
	// iterate over each field in the struct type
	for i := 0; i < t.NumField(); i++ {
		// get the i-th field
		field := t.Field(i)

		// check if the field is marked as "primaryKey" in the struct tag
		if strings.Contains(string(field.Tag), "primaryKey") {
			// if this is the first primaryKey field encountered, set the result to be its value
			if result == "" {
				result = fmt.Sprintf("%v", v.Field(i).Interface())
			} else {
				// if this is not the first primaryKey field, append its value to the result with a "-" separator
				result = fmt.Sprintf("%s-%v", result, v.Field(i).Interface())
			}
		}
	}

	// return the final primary key value as a string
	return result
}
