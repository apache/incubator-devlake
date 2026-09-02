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
	"net/http"

	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/plugin"
	"github.com/apache/devlake/helpers/pluginhelper/api"
	"github.com/apache/devlake/plugins/kiro/models"
)

// PostConnections creates a new connection.
// @Summary create kiro connection
// @Description Create kiro connection
// @Tags plugins/kiro
// @Param body body models.KiroConnection true "json body"
// @Success 200  {object} models.KiroConnection
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/kiro/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.KiroConnection{}
	// Wrapped as BadInput so a struct-tag validation failure reports 400 rather
	// than 500 - the difference between "fix your input" and "the server broke".
	if err := api.Decode(input.Body, connection, vld); err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid connection payload")
	}
	if err := validateConnection(&connection.KiroConn); err != nil {
		return nil, errors.BadInput.Wrap(err, "connection validation failed")
	}
	if err := connectionHelper.Create(connection, input); err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection.Sanitize(), Status: http.StatusOK}, nil
}

// PatchConnection updates an existing connection.
// @Summary patch kiro connection
// @Description Patch kiro connection
// @Tags plugins/kiro
// @Param id path int true "connection ID"
// @Param body body models.KiroConnection true "json body"
// @Success 200  {object} models.KiroConnection
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/kiro/connections/{id} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.KiroConnection{}
	if err := connectionHelper.First(connection, input.Params); err != nil {
		return nil, err
	}
	if err := (&models.KiroConnection{}).MergeFromRequest(connection, input.Body); err != nil {
		return nil, errors.Convert(err)
	}
	if err := validateConnection(&connection.KiroConn); err != nil {
		return nil, errors.BadInput.Wrap(err, "connection validation failed")
	}
	if err := connectionHelper.SaveWithCreateOrUpdate(connection); err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection.Sanitize(), Status: http.StatusOK}, nil
}

// DeleteConnection removes a connection.
// @Summary delete a kiro connection
// @Description Delete a kiro connection
// @Tags plugins/kiro
// @Param id path int true "connection ID"
// @Success 200  {object} models.KiroConnection
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 409  {object} srvhelper.DsRefs "References exist to this connection"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/kiro/connections/{id} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	conn := &models.KiroConnection{}
	output, err := connectionHelper.Delete(conn, input)
	if err != nil {
		return output, err
	}
	output.Body = conn.Sanitize()
	return output, nil
}

// ListConnections lists all connections.
// @Summary get all kiro connections
// @Description Get all kiro connections
// @Tags plugins/kiro
// @Success 200  {object} []models.KiroConnection
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/kiro/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connections []models.KiroConnection
	if err := connectionHelper.List(&connections); err != nil {
		return nil, err
	}
	for i := range connections {
		connections[i] = connections[i].Sanitize()
	}
	return &plugin.ApiResourceOutput{Body: connections}, nil
}

// GetConnection returns one connection.
// @Summary get kiro connection detail
// @Description Get kiro connection detail
// @Tags plugins/kiro
// @Param id path int true "connection ID"
// @Success 200  {object} models.KiroConnection
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/kiro/connections/{id} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.KiroConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection.Sanitize()}, nil
}

// validateConnection checks the fields collection cannot proceed without.
//
// Identity Store fields are deliberately not required: they only resolve display
// names, and identity for joining to git history comes from the report's
// User_Email column. Requiring them would block a working setup.
func validateConnection(conn *models.KiroConn) error {
	if conn.AccessKeyId == "" {
		return errors.BadInput.New("AccessKeyId is required")
	}
	if conn.SecretAccessKey == "" {
		return errors.BadInput.New("SecretAccessKey is required")
	}
	if conn.Region == "" {
		return errors.BadInput.New("Region is required")
	}
	if conn.Bucket == "" {
		return errors.BadInput.New("Bucket is required")
	}
	// A partial Identity Store configuration is a mistake worth reporting: it
	// silently yields no display names, which looks like a data problem rather
	// than a configuration one.
	if (conn.IdentityStoreId == "") != (conn.IdentityStoreRegion == "") {
		return errors.BadInput.New("IdentityStoreId and IdentityStoreRegion must be set together")
	}
	return nil
}
