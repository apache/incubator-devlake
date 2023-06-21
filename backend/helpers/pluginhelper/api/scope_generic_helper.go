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
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer/domaininfo"
	"github.com/apache/incubator-devlake/core/plugin"
	serviceHelper "github.com/apache/incubator-devlake/helpers/pluginhelper/services"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

type NoScopeConfig struct{}

type (
	GenericScopeApiHelper[Conn any, Scope any, ScopeConfig any] struct {
		log              log.Logger
		db               dal.Dal
		validator        *validator.Validate
		reflectionParams *ReflectionParameters
		dbHelper         ScopeDatabaseHelper[Conn, Scope, ScopeConfig]
		bpManager        *serviceHelper.BlueprintManager
		connHelper       *ConnectionApiHelper
		opts             *ScopeHelperOptions
		plugin           string
	}
	// as of golang v1.20, embedding generic fields is not supported
	// let's divide the struct into two parts for swagger doc to work
	// https://stackoverflow.com/questions/66118867/go-generics-is-it-possible-to-embed-generic-structs
	ScopeResDoc[ScopeConfig any] struct {
		ScopeConfig *ScopeConfig        `mapstructure:"scopeConfig,omitempty" json:"scopeConfig"`
		Blueprints  []*models.Blueprint `mapstructure:"blueprints,omitempty" json:"blueprints"`
	}
	// Alias, for swagger purposes
	ScopeRefDoc                          = serviceHelper.BlueprintProjectPairs
	ScopeRes[Scope any, ScopeConfig any] struct {
		Scope                    *Scope                   `mapstructure:",squash"` // ideally we need this field to be embedded in the struct
		ScopeResDoc[ScopeConfig] `mapstructure:",squash"` // however, only this type of embeding is supported as of golang 1.20
	}
	ReflectionParameters struct {
		// This corresponds to the struct field of the scope struct's ID field
		ScopeIdFieldName string `validate:"required"`
		// This corresponds to the database column name of the scope struct's ID (typically primary key) field
		ScopeIdColumnName string `validate:"required"`
		// This corresponds to the scope field on the ApiParams struct of a plugin.
		RawScopeParamName string `validate:"required"`
	}
	ScopeHelperOptions struct {
		// Define this if the raw params doesn't store the ScopeId but a different attribute of the Scope (e.g. Name)
		GetScopeParamValue func(db dal.Dal, scopeId string) (string, errors.Error)
		IsRemote           bool
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

func NewGenericScopeHelper[Conn any, Scope any, ScopeConfig any](
	basicRes context.BasicRes,
	vld *validator.Validate,
	connHelper *ConnectionApiHelper,
	dbHelper ScopeDatabaseHelper[Conn, Scope, ScopeConfig],
	params *ReflectionParameters,
	opts *ScopeHelperOptions,
) *GenericScopeApiHelper[Conn, Scope, ScopeConfig] {
	if connHelper == nil {
		panic("nil connHelper")
	}
	if params == nil {
		panic("reflection params not provided")
	}
	err := vld.Struct(params)
	if err != nil {
		panic(err)
	}
	if opts == nil {
		opts = &ScopeHelperOptions{}
	}
	return &GenericScopeApiHelper[Conn, Scope, ScopeConfig]{
		log:              basicRes.GetLogger(),
		db:               basicRes.GetDal(),
		validator:        vld,
		reflectionParams: params,
		dbHelper:         dbHelper,
		bpManager:        serviceHelper.NewBlueprintManager(basicRes.GetDal()),
		connHelper:       connHelper,
		opts:             opts,
		plugin:           connHelper.pluginName,
	}
}

func (gs *GenericScopeApiHelper[Conn, Scope, ScopeConfig]) DbHelper() ScopeDatabaseHelper[Conn, Scope, ScopeConfig] {
	return gs.dbHelper
}

// hacky, temporary solution
func (gs *GenericScopeApiHelper[Conn, Scope, ScopeConfig]) syncRawDataOrigin(scopes ...*Scope) {
	for _, scope := range scopes {
		toolScope, ok := interface{}(scope).(plugin.ToolLayerScope)
		if !ok {
			panic(fmt.Errorf("%v does not implement the plugin.ToolLayerScope interface", scope))
		}
		setRawDataOrigin(toolScope, common.RawDataOrigin{
			RawDataTable:  fmt.Sprintf("_raw_%s_scopes", gs.plugin),
			RawDataParams: plugin.MarshalScopeParams(toolScope.ScopeParams()),
		})
	}
}

func (gs *GenericScopeApiHelper[Conn, Scope, ScopeConfig]) PutScopes(input *plugin.ApiResourceInput, scopes []*Scope) ([]*ScopeRes[Scope, ScopeConfig], errors.Error) {
	params, err := gs.extractFromReqParam(input, false)
	if err != nil {
		return nil, err
	}
	err = gs.dbHelper.VerifyConnection(params.connectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("error verifying connection for connection ID %d", params.connectionId))
	}
	err = gs.validatePrimaryKeys(scopes)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	for _, scope := range scopes {
		// Set the connection ID, CreatedDate, and UpdatedDate fields
		gs.setScopeFields(scope, params.connectionId, &now, &now)
		err = gs.verifyScope(scope, gs.validator)
		if err != nil {
			return nil, errors.Default.Wrap(err, "error verifying scope")
		}
	}
	// Save the scopes to the database
	if len(scopes) > 0 {
		gs.syncRawDataOrigin(scopes...)
		if err != nil {
			return nil, errors.Default.Wrap(err, "error saving scope")
		}
		err = gs.dbHelper.SaveScope(scopes)
		if err != nil {
			return nil, errors.Default.Wrap(err, "error saving scope")
		}
	}
	apiScopes, err := gs.addScopeConfig(scopes...)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error associating scope config to scope")
	}
	return apiScopes, nil
}

func (gs *GenericScopeApiHelper[Conn, Scope, ScopeConfig]) UpdateScope(input *plugin.ApiResourceInput) (*ScopeRes[Scope, ScopeConfig], errors.Error) {
	params, err := gs.extractFromReqParam(input, true)
	if err != nil {
		return nil, err
	}
	err = gs.dbHelper.VerifyConnection(params.connectionId)
	if err != nil {
		return nil, err
	}
	scope, err := gs.dbHelper.GetScope(params.connectionId, params.scopeId)
	if err != nil {
		return nil, err
	}
	err = DecodeMapStruct(input.Body, scope, false)
	if err != nil {
		return nil, errors.Default.Wrap(err, "patch scope error")
	}
	err = gs.verifyScope(scope, gs.validator)
	if err != nil {
		return nil, errors.Default.Wrap(err, "Invalid scope")
	}
	gs.syncRawDataOrigin(scope)
	err = gs.dbHelper.UpdateScope(params.connectionId, params.scopeId, scope)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving Scope")
	}
	scopeRes, err := gs.addScopeConfig(scope)
	if err != nil {
		return nil, err
	}
	return scopeRes[0], nil
}

func (gs *GenericScopeApiHelper[Conn, Scope, ScopeConfig]) GetScopes(input *plugin.ApiResourceInput) ([]*ScopeRes[Scope, ScopeConfig], errors.Error) {
	params, err := gs.extractFromGetReqParam(input, false)
	if err != nil {
		return nil, errors.BadInput.New("invalid path params: \"connectionId\" not set")
	}
	err = gs.dbHelper.VerifyConnection(params.connectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("error verifying connection for connection ID %d", params.connectionId))
	}
	scopes, err := gs.dbHelper.ListScopes(input, params.connectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("error verifying connection for connection ID %d", params.connectionId))
	}
	apiScopes, err := gs.addScopeConfig(scopes...)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error associating scope configs with scopes")
	}
	// return empty array rather than nil in case of no scopes
	if len(apiScopes) > 0 && params.loadBlueprints {
		scopesById := gs.mapByScopeId(apiScopes)
		var scopeIds []string
		for id := range scopesById {
			scopeIds = append(scopeIds, id)
		}
		blueprintMap, err := gs.bpManager.GetBlueprintsByScopes(params.connectionId, params.plugin, scopeIds...)
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("error getting blueprints for scopes from connection %d", params.connectionId))
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
			gs.log.Warn(nil, "The following dangling scopes were found: %v", danglingIds)
		}
	}
	return apiScopes, nil
}

func (gs *GenericScopeApiHelper[Conn, Scope, ScopeConfig]) GetScope(input *plugin.ApiResourceInput) (*ScopeRes[Scope, ScopeConfig], errors.Error) {
	params, err := gs.extractFromGetReqParam(input, true)
	if err != nil {
		return nil, err
	}
	err = gs.dbHelper.VerifyConnection(params.connectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("error verifying connection for connection ID %d", params.connectionId))
	}
	scope, err := gs.dbHelper.GetScope(params.connectionId, params.scopeId)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("error retrieving scope with scope ID %s", params.scopeId))
	}
	apiScopes, err := gs.addScopeConfig(scope)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("error associating scope config with scope %s", params.scopeId))
	}
	scopeRes := apiScopes[0]
	if params.loadBlueprints {
		blueprintMap, err := gs.bpManager.GetBlueprintsByScopes(params.connectionId, params.plugin, params.scopeId)
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("error getting blueprints for scope with scope ID %s", params.scopeId))
		}
		if len(blueprintMap) == 1 {
			scopeRes.Blueprints = blueprintMap[params.scopeId]
		}
	}
	return scopeRes, nil
}

func (gs *GenericScopeApiHelper[Conn, Scope, ScopeConfig]) DeleteScope(input *plugin.ApiResourceInput) (*serviceHelper.BlueprintProjectPairs, errors.Error) {
	params, err := gs.extractFromDeleteReqParam(input)
	if err != nil {
		return nil, err
	}
	err = gs.dbHelper.VerifyConnection(params.connectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("error verifying connection for connection ID %d", params.connectionId))
	}
	// TODO: the Scope type in this helper must implement the ToolLayerScope interface and we should enforce that.
	scope, err := gs.dbHelper.GetScope(params.connectionId, params.scopeId)
	if err != nil {
		return nil, err
	}
	toolScope, ok := interface{}(scope).(plugin.ToolLayerScope)
	if !ok {
		panic(fmt.Errorf("%v does not implement the plugin.ToolLayerScope interface", scope))
	}
	// now we can as scope to state its `Params` for data bloodline identificaltion
	if refs, err := gs.getScopeReferences(input.GetPlugin(), params.connectionId, params.scopeId); err != nil || refs != nil {
		if err != nil {
			return nil, err
		}
		return refs, errors.Conflict.New("Found one or more references to this scope")
	}
	if err = gs.deleteScopeData(toolScope); err != nil {
		return nil, err
	}
	if !params.deleteDataOnly {
		// Delete the scope itself
		// err = gs.dbHelper.DeleteScope(params.connectionId, params.scopeId)
		err = gs.db.Delete(scope)
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("error deleting scope %s", params.scopeId))
		}
		err = gs.updateBlueprints(params.connectionId, params.plugin, params.scopeId)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (gs *GenericScopeApiHelper[Conn, Scope, ScopeConfig]) addScopeConfig(scopes ...*Scope) ([]*ScopeRes[Scope, ScopeConfig], errors.Error) {
	apiScopes := make([]*ScopeRes[Scope, ScopeConfig], len(scopes))
	for i, scope := range scopes {
		apiScopes[i] = &ScopeRes[Scope, ScopeConfig]{
			Scope: scope,
		}
		scIdField := reflectField(scope, "ScopeConfigId")
		if scIdField.IsValid() && scIdField.Uint() > 0 {
			scopeConfig, err := gs.dbHelper.GetScopeConfig(scIdField.Uint())
			if err != nil {
				return nil, err
			}
			apiScopes[i].ScopeConfig = scopeConfig
		}
	}
	return apiScopes, nil
}

func (gs *GenericScopeApiHelper[Conn, Scope, ScopeConfig]) getScopeReferences(pluginName string, connectionId uint64, scopeId string) (*serviceHelper.BlueprintProjectPairs, errors.Error) {
	blueprintMap, err := gs.bpManager.GetBlueprintsByScopes(connectionId, pluginName, scopeId)
	if err != nil {
		return nil, err
	}
	blueprints := blueprintMap[scopeId]
	if len(blueprints) == 0 {
		return nil, nil
	}
	return serviceHelper.NewBlueprintProjectPairs(blueprints), nil
}

func (gs *GenericScopeApiHelper[Conn, Scope, ScopeConfig]) mapByScopeId(scopes []*ScopeRes[Scope, ScopeConfig]) map[string]*ScopeRes[Scope, ScopeConfig] {
	scopeMap := map[string]*ScopeRes[Scope, ScopeConfig]{}
	for _, scope := range scopes {
		scopeId := fmt.Sprintf("%v", reflectField(scope.Scope, gs.reflectionParams.ScopeIdFieldName).Interface())
		scopeMap[scopeId] = scope
	}
	return scopeMap
}

func (gs *GenericScopeApiHelper[Conn, Scope, ScopeConfig]) extractFromReqParam(input *plugin.ApiResourceInput, withScopeId bool) (*requestParams, errors.Error) {
	connectionId, err := strconv.ParseUint(input.Params["connectionId"], 10, 64)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "Invalid \"connectionId\"")
	}
	if connectionId == 0 {
		return nil, errors.BadInput.New("\"connectionId\" cannot be 0")
	}
	var scopeId string
	if withScopeId {
		scopeId = input.Params["scopeId"]
		// Path params that use `/*param` handlers instead of `/:param` start with a /, so remove it
		if scopeId[0] == '/' {
			scopeId = scopeId[1:]
		}
	}
	pluginName := input.GetPlugin()
	return &requestParams{
		connectionId: connectionId,
		plugin:       pluginName,
		scopeId:      scopeId,
	}, nil
}

func (gs *GenericScopeApiHelper[Conn, Scope, ScopeConfig]) extractFromDeleteReqParam(input *plugin.ApiResourceInput) (*deleteRequestParams, errors.Error) {
	params, err := gs.extractFromReqParam(input, true)
	if err != nil {
		return nil, err
	}
	var deleteDataOnly bool
	{
		ddo, ok := input.Query["delete_data_only"]
		if ok {
			deleteDataOnly, err = errors.Convert01(strconv.ParseBool(ddo[0]))
			if err != nil {
				deleteDataOnly = false
			}
		}
	}
	return &deleteRequestParams{
		requestParams:  *params,
		deleteDataOnly: deleteDataOnly,
	}, nil
}

func (gs *GenericScopeApiHelper[Conn, Scope, ScopeConfig]) extractFromGetReqParam(input *plugin.ApiResourceInput, withScopeId bool) (*getRequestParams, errors.Error) {
	params, err := gs.extractFromReqParam(input, withScopeId)
	if err != nil {
		return nil, err
	}
	var loadBlueprints bool
	{
		lbps, ok := input.Query["blueprints"]
		if ok {
			loadBlueprints, err = errors.Convert01(strconv.ParseBool(lbps[0]))
			if err != nil {
				loadBlueprints = false
			}
		}
	}
	return &getRequestParams{
		requestParams:  *params,
		loadBlueprints: loadBlueprints,
	}, nil
}

func (gs *GenericScopeApiHelper[Conn, Scope, ScopeConfig]) getRawParams(connectionId uint64, scopeId any) string {
	paramsMap := map[string]any{
		"ConnectionId":                        connectionId,
		gs.reflectionParams.RawScopeParamName: scopeId,
	}
	b, err := json.Marshal(paramsMap)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (gs *GenericScopeApiHelper[Conn, Scope, ScopeConfig]) setScopeFields(p interface{}, connectionId uint64, createdDate *time.Time, updatedDate *time.Time) {
	pType := reflect.TypeOf(p)
	if pType.Kind() != reflect.Ptr {
		panic("expected a pointer to a struct")
	}
	pValue := reflectValue(p)
	// set connectionId
	connIdField := pValue.FieldByName("ConnectionId")
	connIdField.SetUint(connectionId)

	// set raw params
	rawParams := pValue.FieldByName("RawDataParams")
	if !rawParams.IsValid() {
		panic("scope is missing the field \"RawDataParams\"")
	}
	scopeIdField := pValue.FieldByName(gs.reflectionParams.ScopeIdFieldName)
	rawParams.Set(reflect.ValueOf(gs.getRawParams(connectionId, scopeIdField.Interface())))

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

func (gs *GenericScopeApiHelper[Conn, Scope, ScopeConfig]) verifyScope(scope interface{}, vld *validator.Validate) errors.Error {
	if gs.opts.IsRemote {
		return nil
	}
	pType := reflect.TypeOf(scope)
	if pType.Kind() != reflect.Ptr {
		panic("expected a pointer to a struct")
	}
	if err := vld.Struct(scope); err != nil {
		return errors.Default.Wrap(err, "error validating target")
	}
	return nil
}

func (gs *GenericScopeApiHelper[Conn, Scope, ScopeConfig]) validatePrimaryKeys(scopes []*Scope) errors.Error {
	if gs.opts.IsRemote {
		return nil
	}
	keeper := make(map[string]struct{})
	for _, scope := range scopes {
		// Ensure that the primary key value is unique
		primaryValueStr := returnPrimaryKeyValue(scope)
		if _, ok := keeper[primaryValueStr]; ok {
			return errors.BadInput.New("duplicate scope was requested")
		} else {
			keeper[primaryValueStr] = struct{}{}
		}
	}
	return nil
}

func (gs *GenericScopeApiHelper[Conn, Scope, ScopeConfig]) updateBlueprints(connectionId uint64, pluginName string, scopeId string) errors.Error {
	blueprintsMap, err := gs.bpManager.GetBlueprintsByScopes(connectionId, pluginName, scopeId)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error retrieving scope with scope ID %s", scopeId))
	}
	blueprints := blueprintsMap[scopeId]
	// update the blueprints (remove scope reference from them)
	for _, blueprint := range blueprints {
		settings, _ := blueprint.UnmarshalSettings()
		var changed bool
		err = settings.UpdateConnections(func(c *plugin.BlueprintConnectionV200) errors.Error {
			var retainedScopes []*plugin.BlueprintScopeV200
			for _, bpScope := range c.Scopes {
				if bpScope.Id == scopeId { // we'll be removing this one
					changed = true
				} else {
					retainedScopes = append(retainedScopes, bpScope)
				}
			}
			c.Scopes = retainedScopes
			return nil
		})
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("error removing scope %s from blueprint %d", scopeId, blueprint.ID))
		}
		if changed {
			err = blueprint.UpdateSettings(&settings)
			if err != nil {
				return errors.Default.Wrap(err, fmt.Sprintf("error writing new settings into blueprint %s", blueprint.Name))
			}
			err = gs.bpManager.SaveDbBlueprint(blueprint)
			if err != nil {
				return errors.Default.Wrap(err, fmt.Sprintf("error saving the updated blueprint %s", blueprint.Name))
			}
		}
	}
	return nil
}

func (gs *GenericScopeApiHelper[Conn, Scope, ScopeConfig]) deleteScopeData(scope plugin.ToolLayerScope) errors.Error {
	// find all tables for this plugin
	tables, err := gs.getAffectedTables(gs.plugin)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error getting database tables managed by plugin %s", gs.plugin))
	}
	scopeParams := plugin.MarshalScopeParams(scope.ScopeParams())
	err = gs.transactionalDelete(tables, scopeParams)
	if err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("error deleting data bound to scope %s for plugin %s", scopeParams, gs.plugin))
	}
	return nil
}

func (gs *GenericScopeApiHelper[Conn, Scope, ScopeConfig]) transactionalDelete(tables []string, rawDataParams string) errors.Error {
	tx := gs.db.Begin()
	for _, table := range tables {
		var where string
		var params []interface{}
		if strings.HasPrefix(table, "_raw_") {
			// raw table: shhould check connection and scope
			where = "params = ?"
			params = []interface{}{rawDataParams}
		} else if strings.HasPrefix(table, "_tool_") {
			// tool layer table: should check connection and scope
			where = "_raw_data_params = ?"
			params = []interface{}{rawDataParams}
		} else {
			// framework tables: should check plugin, connection and scope
			if table == (models.CollectorLatestState{}.TableName()) {
				// diff sync state
				where = "raw_data_table LIKE ? AND raw_data_params = ?"
			} else {
				// domain layer table
				where = "_raw_data_table LIKE ? AND _raw_data_params = ?"
			}
			rawDataTablePrefix := fmt.Sprintf("_raw_%s%%", gs.plugin)
			params = []interface{}{rawDataTablePrefix, rawDataParams}
		}
		sql := fmt.Sprintf("DELETE FROM %s WHERE %s", table, where)
		err := tx.Exec(sql, params...)
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				gs.log.Warn(err2, fmt.Sprintf("error rolling back table data deletion transaction. sql: %s params: %v", sql, params))
			}
			return err
		}
	}
	err := tx.Commit()
	if err != nil {
		return errors.Default.Wrap(err, "error committing delete transaction for plugin tables")
	}
	return nil
}

// Implement MarshalJSON method to flatten all fields
func (sr *ScopeRes[T, Y]) MarshalJSON() ([]byte, error) {
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

func (gs *GenericScopeApiHelper[Conn, Scope, ScopeConfig]) getAffectedTables(pluginName string) ([]string, errors.Error) {
	var tables []string
	meta, err := plugin.GetPlugin(pluginName)
	if err != nil {
		return nil, err
	}
	if pluginModel, ok := meta.(plugin.PluginModel); !ok {
		return nil, errors.Default.New(fmt.Sprintf("plugin \"%s\" does not implement listing its tables", pluginName))
	} else {
		// Unfortunately, can't cache the tables because Python creates some tables on a per-demand basis, so such a cache would possibly get outdated.
		// It's a rare scenario in practice, but might as well play it safe and sacrifice some performance here
		var allTables []string
		if allTables, err = gs.db.AllTables(); err != nil {
			return nil, err
		}
		// collect raw tables
		for _, table := range allTables {
			if strings.HasPrefix(table, "_raw_"+pluginName) {
				tables = append(tables, table)
			}
		}
		// collect tool tables
		toolModels := pluginModel.GetTablesInfo()
		for _, toolModel := range toolModels {
			if !isScopeModel(toolModel) && hasField(toolModel, "RawDataParams") {
				tables = append(tables, toolModel.TableName())
			}
		}
		// collect domain tables
		for _, domainModel := range domaininfo.GetDomainTablesInfo() {
			// we only care about tables with RawOrigin
			ok = hasField(domainModel, "RawDataParams")
			if ok {
				tables = append(tables, domainModel.TableName())
			}
		}
		// additional tables
		tables = append(tables, models.CollectorLatestState{}.TableName())
	}
	gs.log.Debug("Discovered %d tables used by plugin \"%s\": %v", len(tables), pluginName, tables)
	return tables, nil
}

func isScopeModel(obj dal.Tabler) bool {
	_, ok := obj.(plugin.ToolLayerScope)
	return ok
}
