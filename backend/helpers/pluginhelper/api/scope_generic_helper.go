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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/domainlayer/domaininfo"
	"github.com/apache/incubator-devlake/core/plugin"
	serviceHelper "github.com/apache/incubator-devlake/helpers/pluginhelper/services"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	tablesCache       []string // these cached vars can probably be moved somewhere more centralized later
	tablesCacheLoader = new(sync.Once)
)

type NoTransformation struct{}

type (
	GenericScopeApiHelper[Conn any, Scope any, Tr any] struct {
		log              log.Logger
		db               dal.Dal
		validator        *validator.Validate
		reflectionParams *ReflectionParameters
		dbHelper         ScopeDatabaseHelper[Conn, Scope, Tr]
		connHelper       *ConnectionApiHelper
	}
	ReflectionParameters struct {
		ScopeIdFieldName  string
		ScopeIdColumnName string
		RawScopeParamName string
	}
)

type (
	requestParams struct {
		connectionId uint64
		scopeId      string
		plugin       string
	}
	deleteRequestParams struct {
		requestParams
		deleteDataOnly bool
	}

	getRequestParams struct {
		requestParams
		loadBlueprints bool
	}
)

func NewGenericScopeHelper[Conn any, Scope any, Tr any](
	basicRes context.BasicRes,
	vld *validator.Validate,
	connHelper *ConnectionApiHelper,
	dbHelper ScopeDatabaseHelper[Conn, Scope, Tr],
	params *ReflectionParameters,
) *GenericScopeApiHelper[Conn, Scope, Tr] {
	if connHelper == nil {
		return nil
	}
	if params == nil {
		panic("reflection params not provided")
	}
	tablesCacheLoader.Do(func() {
		var err errors.Error
		tablesCache, err = basicRes.GetDal().AllTables()
		if err != nil {
			panic(err)
		}
	})
	return &GenericScopeApiHelper[Conn, Scope, Tr]{
		log:              basicRes.GetLogger(),
		db:               basicRes.GetDal(),
		validator:        vld,
		reflectionParams: params,
		dbHelper:         dbHelper,
		connHelper:       connHelper,
	}
}

func (c *GenericScopeApiHelper[Conn, Scope, Tr]) PutScopes(input *plugin.ApiResourceInput, scopes []*Scope) ([]*ScopeRes[Scope], errors.Error) {
	params := c.extractFromReqParam(input)
	if params.connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}
	err := c.dbHelper.VerifyConnection(params.connectionId)
	if err != nil {
		return nil, err
	}
	// Create a map to keep track of primary key values
	keeper := make(map[string]struct{})

	// Set the CreatedDate and UpdatedDate fields to the current time for each scope
	now := time.Now()
	for _, scope := range scopes {
		// Ensure that the primary key value is unique (for validatable types)
		if c.validator != nil {
			primaryValueStr := returnPrimaryKeyValue(scope)
			if _, ok := keeper[primaryValueStr]; ok {
				return nil, errors.BadInput.New("duplicated item")
			} else {
				keeper[primaryValueStr] = struct{}{}
			}
		}
		b, _ := json.Marshal(scope)
		_ = b
		// Set the connection ID, CreatedDate, and UpdatedDate fields
		setScopeFields(scope, params.connectionId, &now, &now)

		//Verify that the primary key value is valid
		err = VerifyScope(scope, c.validator)
		if err != nil {
			return nil, err
		}
	}
	// Save the scopes to the database
	if scopes != nil && len(scopes) > 0 {
		err = c.dbHelper.SaveScope(scopes)
		if err != nil {
			return nil, err
		}
	}

	apiScopes, err := c.addTransformationName(scopes)
	if err != nil {
		return nil, err
	}

	return apiScopes, nil
}

func (c *GenericScopeApiHelper[Conn, Scope, Tr]) GetScopes(input *plugin.ApiResourceInput) ([]*ScopeRes[Scope], errors.Error) {
	params := c.extractFromGetReqParam(input)
	if params.connectionId == 0 {
		return nil, errors.BadInput.New("invalid path params: \"connectionId\" not set")
	}
	err := c.dbHelper.VerifyConnection(params.connectionId)
	if err != nil {
		return nil, err
	}
	scopes, err := c.dbHelper.ListScopes(input, params.connectionId)
	if err != nil {
		return nil, err
	}
	apiScopes, err := c.addTransformationName(scopes)
	if err != nil {
		return nil, err
	}
	if params.loadBlueprints {
		scopesById := c.mapByScopeId(apiScopes)
		var scopeIds []string
		for id := range scopesById {
			scopeIds = append(scopeIds, id)
		}
		blueprintMap, err := serviceHelper.NewBlueprintManager(c.db).GetBlueprintsByScopes(scopeIds...)
		if err != nil {
			return nil, errors.Default.Wrap(err, "error getting blueprints for scope")
		}
		apiScopes = nil
		for scopeId, scope := range scopesById {
			if bps, ok := blueprintMap[scopeId]; ok {
				scope.Blueprints = bps
				delete(blueprintMap, scopeId)
			}
			apiScopes = append(apiScopes, scope)
		}
		if len(blueprintMap) > 0 {
			var danglingIds []string
			for bpId := range blueprintMap {
				danglingIds = append(danglingIds, bpId)
			}
			c.log.Warn(nil, "The following dangling scopes were found: %v", danglingIds)
		}
	}
	return apiScopes, nil
}

func (c *GenericScopeApiHelper[Conn, Scope, Tr]) GetScope(input *plugin.ApiResourceInput) (*ScopeRes[Scope], errors.Error) {
	params := c.extractFromGetReqParam(input)
	if params == nil || params.connectionId == 0 {
		return nil, errors.BadInput.New("invalid path params: \"connectionId\" not set")
	}
	if len(params.scopeId) == 0 || params.scopeId == "0" {
		return nil, errors.BadInput.New("invalid path params: \"scopeId\" not set/invalid")
	}
	err := c.dbHelper.VerifyConnection(params.connectionId)
	if err != nil {
		return nil, err
	}
	db := c.db
	scope, err := c.dbHelper.GetScope(params.connectionId, params.scopeId)
	if err != nil {
		return nil, err
	}
	valueRepoRuleId := reflect.ValueOf(scope).FieldByName("TransformationRuleId")
	transformationRuleName := ""
	if valueRepoRuleId.IsValid() {
		repoRuleId := valueRepoRuleId.Uint()
		var rule any
		if repoRuleId > 0 {
			rule, err = c.dbHelper.GetTransformationRule(repoRuleId)
			if err != nil {
				return nil, err
			}
		}
		transformationRuleName = reflect.ValueOf(rule).FieldByName("Name").String()
	}
	var blueprints []*models.Blueprint
	if params.loadBlueprints {
		blueprintMap, err := serviceHelper.NewBlueprintManager(db).GetBlueprintsByScopes(params.scopeId)
		if err != nil {
			return nil, errors.Default.Wrap(err, "error getting blueprints for scope")
		}
		if len(blueprintMap) == 1 {
			blueprints = blueprintMap[params.scopeId]
		}
	}
	scopeRes := &ScopeRes[Scope]{
		Scope:                  scope,
		TransformationRuleName: transformationRuleName,
		Blueprints:             blueprints,
	}
	return scopeRes, nil
}

func (c *GenericScopeApiHelper[Conn, Scope, Tr]) DeleteScope(input *plugin.ApiResourceInput) errors.Error {
	params := c.extractFromDeleteReqParam(input)
	if params == nil || params.connectionId == 0 {
		return errors.BadInput.New("invalid path params: \"connectionId\" not set")
	}
	if len(params.scopeId) == 0 || params.scopeId == "0" {
		return errors.BadInput.New("invalid path params: \"scopeId\" not set/invalid")
	}
	err := c.dbHelper.VerifyConnection(params.connectionId)
	if err != nil {
		return err
	}
	db := c.db
	bpManager := serviceHelper.NewBlueprintManager(db)
	blueprintsMap, err := bpManager.GetBlueprintsByScopes(params.scopeId)
	if err != nil {
		return err
	}
	blueprints := blueprintsMap[params.scopeId]
	// find all tables for this plugin
	tables, err := getPluginTables(params.plugin)
	if err != nil {
		return err
	}
	// delete all the plugin records referencing this scope
	if c.reflectionParams.RawScopeParamName != "" {
		for _, table := range tables {
			err = db.Exec(createDeleteQuery(table, c.reflectionParams.RawScopeParamName, params.scopeId))
			if err != nil {
				return err
			}
		}
	}
	if !params.deleteDataOnly {
		// DeleteScope the scope itself
		err = c.dbHelper.DeleteScope(params.connectionId, params.scopeId)
		if err != nil {
			return err
		}
		// update the blueprints (remove scope reference from them)
		for _, blueprint := range blueprints {
			settings, _ := blueprint.UnmarshalSettings()
			err = settings.UpdateConnections(func(c *plugin.BlueprintConnectionV200) errors.Error {
				var filteredScopes []*plugin.BlueprintScopeV200
				for _, bpScope := range c.Scopes {
					if bpScope.Id != params.scopeId { // keep the ones NOT equal to this scope
						filteredScopes = append(filteredScopes, bpScope)
					}
				}
				c.Scopes = filteredScopes
				return nil
			})
			if err != nil {
				return errors.Default.Wrap(err, fmt.Sprintf("error removing scope %s from blueprint %d", params.scopeId, blueprint.ID))
			}
		}
	}
	return nil
}

func (c *GenericScopeApiHelper[Conn, Scope, Tr]) addTransformationName(scopes []*Scope) ([]*ScopeRes[Scope], errors.Error) {
	var ruleIds []uint64
	for _, scope := range scopes {
		valueRepoRuleId := reflectField(scope, "TransformationRuleId")
		if !valueRepoRuleId.IsValid() {
			break
		}
		ruleId := reflectField(scope, "TransformationRuleId").Uint()
		if ruleId > 0 {
			ruleIds = append(ruleIds, ruleId)
		}
	}
	var rules []*Tr
	var err errors.Error
	if len(ruleIds) > 0 {
		rules, err = c.dbHelper.ListTransformationRules(ruleIds)
		if err != nil {
			return nil, err
		}
	}
	names := make(map[uint64]string)
	for _, rule := range rules {
		// Get the reflect.Value of the i-th struct pointer in the slice
		names[reflectField(rule, "ID").Uint()] = reflectField(rule, "Name").String()
	}
	apiScopes := make([]*ScopeRes[Scope], 0)
	for _, scope := range scopes {
		txRuleField := reflectField(scope, "TransformationRuleId")
		txRuleName := ""
		if txRuleField.IsValid() {
			txRuleName = names[txRuleField.Uint()]
		}
		apiScopes = append(apiScopes, &ScopeRes[Scope]{
			Scope:                  *scope,
			TransformationRuleName: txRuleName,
		})
	}
	return apiScopes, nil
}

func (c *GenericScopeApiHelper[Conn, Scope, Tr]) mapByScopeId(scopes []*ScopeRes[Scope]) map[string]*ScopeRes[Scope] {
	scopeMap := map[string]*ScopeRes[Scope]{}
	for _, scope := range scopes {
		scopeId := fmt.Sprintf("%v", reflectField(scope.Scope, c.reflectionParams.ScopeIdFieldName).Interface())
		scopeMap[scopeId] = scope
	}
	return scopeMap
}

func (c *GenericScopeApiHelper[Conn, Scope, Tr]) extractFromReqParam(input *plugin.ApiResourceInput) *requestParams {
	connectionId, err := strconv.ParseUint(input.Params["connectionId"], 10, 64)
	if err != nil || connectionId == 0 {
		connectionId = 0
	}
	scopeId, _ := input.Params["scopeId"]
	pluginName := input.Params["plugin"]
	return &requestParams{
		connectionId: connectionId,
		scopeId:      scopeId,
		plugin:       pluginName,
	}
}

func (c *GenericScopeApiHelper[Conn, Scope, Tr]) extractFromDeleteReqParam(input *plugin.ApiResourceInput) *deleteRequestParams {
	params := c.extractFromReqParam(input)
	var err errors.Error
	var deleteDataOnly bool
	{
		ddo, ok := input.Query["delete_data_only"]
		if ok {
			deleteDataOnly, err = errors.Convert01(strconv.ParseBool(ddo[0]))
		}
		if err != nil {
			deleteDataOnly = false
		}
	}
	return &deleteRequestParams{
		requestParams:  *params,
		deleteDataOnly: deleteDataOnly,
	}
}

func (c *GenericScopeApiHelper[Conn, Scope, Tr]) extractFromGetReqParam(input *plugin.ApiResourceInput) *getRequestParams {
	params := c.extractFromReqParam(input)
	var err errors.Error
	var loadBlueprints bool
	{
		lbps, ok := input.Query["blueprints"]
		if ok {
			loadBlueprints, err = errors.Convert01(strconv.ParseBool(lbps[0]))
		}
		if err != nil {
			loadBlueprints = false
		}
	}
	return &getRequestParams{
		requestParams:  *params,
		loadBlueprints: loadBlueprints,
	}
}

func setScopeFields(p interface{}, connectionId uint64, createdDate *time.Time, updatedDate *time.Time) {
	pType := reflect.TypeOf(p)
	if pType.Kind() != reflect.Ptr {
		panic("expected a pointer to a struct")
	}
	pValue := reflectValue(p)
	// set connectionId
	connIdField := pValue.FieldByName("ConnectionId")
	connIdField.SetUint(connectionId)

	// TODO might need to change these to CreatedAt and UpdatedAt

	// set CreatedDate
	createdDateField := pValue.FieldByName("CreatedDate")
	if createdDateField.IsValid() && createdDateField.Type().AssignableTo(reflect.TypeOf(createdDate)) {
		createdDateField.Set(reflect.ValueOf(createdDate))
	}

	// set UpdatedDate
	updatedDateField := pValue.FieldByName("UpdatedDate")
	if !updatedDateField.IsValid() || (updatedDate != nil && !updatedDateField.Type().AssignableTo(reflect.TypeOf(updatedDate))) {
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
	t := reflectType(p)
	v := reflectValue(p)
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

func createDeleteQuery(tableName string, scopeIdKey string, scopeId string) string {
	column := "_raw_data_params"
	if strings.HasPrefix(tableName, "_raw_") {
		column = "params"
	}
	query := `DELETE FROM ` + tableName + ` WHERE ` + column + ` LIKE '%"` + scopeIdKey + `":"` + scopeId + `"%'`
	return query
}

func getPluginTables(pluginName string) ([]string, errors.Error) {
	var tables []string
	meta, err := plugin.GetPlugin(pluginName)
	if err != nil {
		return nil, err
	}
	if pluginModel, ok := meta.(plugin.PluginModel); !ok {
		return nil, errors.Default.New(fmt.Sprintf("plugin \"%s\" does not implement listing its tables", pluginName))
	} else {
		// collect raw tables
		for _, table := range tablesCache {
			if strings.HasPrefix(table, "_raw_"+pluginName) {
				tables = append(tables, table)
			}
		}
		// collect tool tables
		tablesInfo := pluginModel.GetTablesInfo()
		for _, table := range tablesInfo {
			tables = append(tables, table.TableName())
		}
		// collect domain tables
		for _, domainTable := range domaininfo.GetDomainTablesInfo() {
			// we only care about tables with RawOrigin
			_, ok = reflect.TypeOf(domainTable).Elem().FieldByName("RawDataParams")
			if ok {
				tables = append(tables, domainTable.TableName())
			}
		}
	}
	return tables, nil
}

func reflectField(obj any, fieldName string) reflect.Value {
	return reflectValue(obj).FieldByName(fieldName)
}

func reflectValue(obj any) reflect.Value {
	val := reflect.ValueOf(obj)
	kind := val.Kind()
	for kind == reflect.Ptr || kind == reflect.Interface {
		val = val.Elem()
		kind = val.Kind()
	}
	return val
}

func reflectType(obj any) reflect.Type {
	typ := reflect.TypeOf(obj)
	kind := typ.Kind()
	for kind == reflect.Ptr {
		typ = typ.Elem()
		kind = typ.Kind()
	}
	return typ
}
