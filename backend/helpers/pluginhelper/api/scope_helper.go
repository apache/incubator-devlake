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
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"

	"reflect"
)

// ScopeApiHelper is used to write the CURD of connection
type ScopeApiHelper[Conn any, Scope any, Tr any] struct {
	log        log.Logger
	db         dal.Dal
	validator  *validator.Validate
	connHelper *ConnectionApiHelper
}

// NewScopeHelper creates a ScopeHelper for connection management
func NewScopeHelper[Conn any, Scope any, Tr any](
	basicRes context.BasicRes,
	vld *validator.Validate,
	connHelper *ConnectionApiHelper,
) *ScopeApiHelper[Conn, Scope, Tr] {
	if vld == nil {
		vld = validator.New()
	}
	if connHelper == nil {
		return nil
	}
	return &ScopeApiHelper[Conn, Scope, Tr]{
		log:        basicRes.GetLogger(),
		db:         basicRes.GetDal(),
		validator:  vld,
		connHelper: connHelper,
	}
}

type ScopeRes[T any] struct {
	Scope                  T      `mapstructure:",squash"`
	TransformationRuleName string `mapstructure:"transformationRuleName,omitempty"`
}

type ScopeReq[T any] struct {
	Data []*T `json:"data"`
}

// Put saves the given scopes to the database. It expects a slice of struct pointers
// as the scopes argument. It also expects a fieldName argument, which is used to extract
// the connection ID from the input.Params map.
func (c *ScopeApiHelper[Conn, Scope, Tr]) Put(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var req struct {
		Data []*Scope `json:"data"`
	}
	err := errors.Convert(DecodeMapStruct(input.Body, &req))
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "decoding Github repo error")
	}
	// Extract the connection ID from the input.Params map
	connectionId, _ := extractFromReqParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId or scopeId")
	}
	err = c.VerifyConnection(connectionId)
	if err != nil {
		return nil, err
	}
	// Create a map to keep track of primary key values
	keeper := make(map[string]struct{})

	// Set the CreatedDate and UpdatedDate fields to the current time for each scope
	now := time.Now()
	for _, v := range req.Data {
		// Ensure that the primary key value is unique
		primaryValueStr := returnPrimaryKeyValue(*v)
		if _, ok := keeper[primaryValueStr]; ok {
			return nil, errors.BadInput.New("duplicated item")
		} else {
			keeper[primaryValueStr] = struct{}{}
		}

		// Set the connection ID, CreatedDate, and UpdatedDate fields
		setScopeFields(v, connectionId, &now, &now)

		// Verify that the primary key value is valid
		err = VerifyScope(v, c.validator)
		if err != nil {
			return nil, err
		}
	}
	err = c.save(&req.Data)
	if err != nil {
		return nil, err
	}

	// Save the scopes to the database
	return &plugin.ApiResourceOutput{Body: req.Data, Status: http.StatusOK}, nil
}

func (c *ScopeApiHelper[Conn, Scope, Tr]) Update(input *plugin.ApiResourceInput, fieldName string) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, scopeId := extractFromReqParam(input.Params)

	if connectionId == 0 || len(scopeId) == 0 || scopeId == "0" {
		return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusInternalServerError}, errors.BadInput.New("invalid connectionId")
	}
	err := c.VerifyConnection(connectionId)
	if err != nil {
		return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusInternalServerError}, err
	}
	var scope Scope
	err = c.db.First(&scope, dal.Where(fmt.Sprintf("connection_id = ? AND %s = ?", fieldName), connectionId, scopeId))
	if err != nil {
		return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusInternalServerError}, errors.Default.New("getting Scope error")
	}
	err = DecodeMapStruct(input.Body, &scope)
	if err != nil {
		return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusInternalServerError}, errors.Default.Wrap(err, "patch scope error")
	}
	err = VerifyScope(&scope, c.validator)
	if err != nil {
		return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusInternalServerError}, errors.Default.Wrap(err, "Invalid scope")
	}

	err = c.db.Update(scope)
	if err != nil {
		return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusInternalServerError}, errors.Default.Wrap(err, "error on saving Scope")
	}
	return &plugin.ApiResourceOutput{Body: &scope, Status: http.StatusOK}, nil
}

func (c *ScopeApiHelper[Conn, Scope, Tr]) GetScopeList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, _ := extractFromReqParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid path params")
	}
	err := c.VerifyConnection(connectionId)
	if err != nil {
		return nil, err
	}
	limit, offset := GetLimitOffset(input.Query, "pageSize", "page")
	var scopes []*Scope
	err = c.db.All(&scopes, dal.Where("connection_id = ?", connectionId), dal.Limit(limit), dal.Offset(offset))
	if err != nil {
		return nil, err
	}

	var ruleIds []uint64
	for _, scope := range scopes {
		valueRepoRuleId := reflect.ValueOf(scope).Elem().FieldByName("TransformationRuleId")
		if !valueRepoRuleId.IsValid() {
			return &plugin.ApiResourceOutput{Body: scopes, Status: http.StatusOK}, nil
		}
		ruleId := reflect.ValueOf(scope).Elem().FieldByName("TransformationRuleId").Uint()
		if ruleId > 0 {
			ruleIds = append(ruleIds, ruleId)
		}
	}
	var rules []*Tr
	if len(ruleIds) > 0 {
		err = c.db.All(&rules, dal.Where("id IN (?)", ruleIds))
		if err != nil {
			return nil, err
		}
	}
	names := make(map[uint64]string)
	for _, rule := range rules {
		// Get the reflect.Value of the i-th struct pointer in the slice
		names[reflect.ValueOf(rule).Elem().FieldByName("ID").Uint()] = reflect.ValueOf(rule).Elem().FieldByName("Name").String()
	}
	apiScopes := make([]ScopeRes[Scope], 0)
	for _, scope := range scopes {
		apiScopes = append(apiScopes, ScopeRes[Scope]{*scope, names[reflect.ValueOf(scope).Elem().FieldByName("TransformationRuleId").Uint()]})
	}
	return &plugin.ApiResourceOutput{Body: apiScopes, Status: http.StatusOK}, nil
}

func (c *ScopeApiHelper[Conn, Scope, Tr]) GetScope(input *plugin.ApiResourceInput, fieldName string) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, scopeId := extractFromReqParam(input.Params)
	if connectionId == 0 || len(scopeId) == 0 || scopeId == "0" {
		return nil, errors.BadInput.New("invalid path params")
	}
	err := c.VerifyConnection(connectionId)
	if err != nil {
		return nil, err
	}
	db := c.db

	query := dal.Where(fmt.Sprintf("connection_id = ? AND %s = ?", fieldName), connectionId, scopeId)
	var scope Scope
	err = db.First(&scope, query)
	if db.IsErrorNotFound(err) {
		return nil, errors.NotFound.New("Scope not found")
	}
	if err != nil {
		return nil, err
	}
	valueRepoRuleId := reflect.ValueOf(scope).FieldByName("TransformationRuleId")
	if !valueRepoRuleId.IsValid() {
		return &plugin.ApiResourceOutput{Body: scope, Status: http.StatusOK}, nil
	}
	repoRuleId := reflect.ValueOf(scope).FieldByName("TransformationRuleId").Uint()
	var rule Tr
	if repoRuleId > 0 {
		err = db.First(&rule, dal.Where("id = ?", repoRuleId))
		if err != nil {
			return nil, errors.NotFound.New("transformationRule not found")
		}
	}
	scopeRes := &ScopeRes[Scope]{scope, reflect.ValueOf(rule).FieldByName("Name").String()}
	return &plugin.ApiResourceOutput{Body: scopeRes, Status: http.StatusOK}, nil
}

func (c *ScopeApiHelper[Conn, Scope, Tr]) VerifyConnection(connId uint64) errors.Error {
	var conn Conn
	err := c.connHelper.FirstById(&conn, connId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.BadInput.New("Invalid Connection Id")
		}
		return err
	}
	return nil
}

func (c *ScopeApiHelper[Conn, Scope, Tr]) save(scope interface{}) errors.Error {
	err := c.db.CreateOrUpdate(scope)
	if err != nil {
		if c.db.IsDuplicationError(err) {
			return errors.BadInput.New("the scope already exists")
		}
		return err
	}
	return nil
}

func extractFromReqParam(params map[string]string) (uint64, string) {
	connectionId, err := strconv.ParseUint(params["connectionId"], 10, 64)
	if err != nil {
		return 0, ""
	}
	scopeId := params["scopeId"]
	return connectionId, scopeId
}

func setScopeFields(p interface{}, connectionId uint64, createdDate *time.Time, updatedDate *time.Time) {
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
	if createdDateField.IsValid() {
		createdDateField.Set(reflect.ValueOf(createdDate))
	}

	// set UpdatedDate
	updatedDateField := pValue.FieldByName("UpdatedDate")
	if !updatedDateField.IsValid() {
		return
	}
	if updatedDate == nil {
		// if updatedDate is nil, set UpdatedDate to be nil
		updatedDateField.Set(reflect.Zero(updatedDateField.Type()))
	} else {
		// if updatedDate is not nil, set UpdatedDate to be the value
		updatedDateFieldValue := reflect.ValueOf(updatedDate)
		updatedDateField.Set(updatedDateFieldValue)
	}
}

// returnPrimaryKeyValue returns a string containing the primary key value(s) of a struct, concatenated with "-" between them.
// This function receives an interface{} type argument p, which can be a pointer to any struct.
// The function uses reflection to iterate through the fields of the struct, and checks if each field is tagged as "primaryKey".
func returnPrimaryKeyValue(p interface{}) string {
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

func VerifyScope(scope interface{}, vld *validator.Validate) errors.Error {
	if vld != nil {
		pType := reflect.TypeOf(scope)
		if pType.Kind() != reflect.Ptr {
			panic("expected a pointer to a struct")
		}
		if err := vld.Struct(scope); err != nil {
			return errors.Default.Wrap(err, "error validating target")
		}
	}
	return nil
}

// Implement MarshalJSON method to flatten all fields
func (sr *ScopeRes[T]) MarshalJSON() ([]byte, error) {
	var flatMap map[string]interface{}
	err := mapstructure.Decode(sr, &flatMap)
	if err != nil {
		return nil, err
	}
	// Encode the flattened map to JSON
	result, err := json.Marshal(flatMap)
	if err != nil {
		return nil, err
	}

	return result, nil
}
