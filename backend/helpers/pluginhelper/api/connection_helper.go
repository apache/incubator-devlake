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
	"strconv"

	"github.com/apache/incubator-devlake/helpers/pluginhelper/services"
	"github.com/apache/incubator-devlake/server/api/shared"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/go-playground/validator/v10"
)

// ConnectionApiHelper is used to write the CURD of connection
type ConnectionApiHelper struct {
	encryptionSecret string
	log              log.Logger
	db               dal.Dal
	validator        *validator.Validate
	bpManager        *services.BlueprintManager
	pluginName       string
}

// NewConnectionHelper creates a ConnectionHelper for connection management
func NewConnectionHelper(
	basicRes context.BasicRes,
	vld *validator.Validate,
	pluginName string,
) *ConnectionApiHelper {
	if vld == nil {
		vld = validator.New()
	}
	return &ConnectionApiHelper{
		encryptionSecret: basicRes.GetConfig(plugin.EncodeKeyEnvStr),
		log:              basicRes.GetLogger(),
		db:               basicRes.GetDal(),
		validator:        vld,
		bpManager:        services.NewBlueprintManager(basicRes.GetDal()),
		pluginName:       pluginName,
	}
}

// Create a connection record based on request body
func (c *ConnectionApiHelper) Create(connection interface{}, input *plugin.ApiResourceInput) errors.Error {
	return c.CreateWithTx(nil, connection, input)
}

// Create a connection record based on request body
func (c *ConnectionApiHelper) CreateWithTx(tx dal.Transaction, connection interface{}, input *plugin.ApiResourceInput) errors.Error {
	// update fields from request body
	db := c.db
	if tx != nil {
		db = tx
	}
	err := c.merge(connection, input.Body)
	if err != nil {
		return err
	}
	if err := c.save(connection, db.Create); err != nil {
		c.log.Error(err, "create connection")
		return err
	}
	return nil
}

// Patch (Modify) a connection record based on request body
func (c *ConnectionApiHelper) Patch(connection interface{}, input *plugin.ApiResourceInput) errors.Error {
	err := c.First(connection, input.Params)
	if err != nil {
		return err
	}

	err = c.merge(connection, input.Body)
	if err != nil {
		return err
	}
	return c.save(connection, c.db.CreateOrUpdate)
}

// First finds connection from db  by parsing request input and decrypt it
func (c *ConnectionApiHelper) First(connection interface{}, params map[string]string) errors.Error {
	connectionId := params["connectionId"]
	if connectionId == "" {
		return errors.BadInput.New("missing connectionId")
	}
	id, err := strconv.ParseUint(connectionId, 10, 64)
	if err != nil || id < 1 {
		return errors.BadInput.New("invalid connectionId")
	}
	return c.FirstById(connection, id)
}

// FirstById finds connection from db by id and decrypt it
func (c *ConnectionApiHelper) FirstById(connection interface{}, id uint64) errors.Error {
	return CallDB(c.db.First, connection, dal.Where("id = ?", id))
}

// List returns all connections with password/token decrypted
func (c *ConnectionApiHelper) List(connections interface{}) errors.Error {
	return CallDB(c.db.All, connections)
}

// Delete connection.
func (c *ConnectionApiHelper) deleteConnection(connection interface{}) (*services.BlueprintProjectPairs, errors.Error) {
	connectionId := reflectField(connection, "ID").Uint()
	referencingBps := c.bpManager.GetBlueprintsByConnection(c.pluginName, connectionId)
	if len(referencingBps) > 0 {
		return services.NewBlueprintProjectPairs(referencingBps), errors.Conflict.New("Found one or more blueprint/project references to this connection")
	}
	src, err := c.getPluginSource()
	if err != nil {
		return nil, err
	}
	if scopeModel := src.Scope(); scopeModel != nil {
		// ensure the connection has no scopes using it
		count := errors.Must1(c.db.Count(dal.From(scopeModel.TableName()), dal.Where("connection_id = ?", connectionId)))
		if count > 0 {
			return nil, errors.Conflict.New("Please delete all data scope(s) before you delete this Data Connection.")
		}
	}
	if scopeConfigModel := src.ScopeConfig(); scopeConfigModel != nil {
		// remove scope-configs that use this connection
		err = CallDB(c.db.Delete, scopeConfigModel, dal.Where("connection_id = ?", connectionId))
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("error deleting scope-configs for plugin %s using connection %d", c.pluginName, connectionId))
		}
	}
	return nil, CallDB(c.db.Delete, connection)
}

// TODO: combine connection/scope/scopeConfig helper
func (c *ConnectionApiHelper) Delete(connection interface{}, input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	err := c.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	var refs *services.BlueprintProjectPairs
	refs, err = c.deleteConnection(connection)
	if err != nil {
		return &plugin.ApiResourceOutput{Body: &shared.ApiBody{
			Success: false,
			Message: err.Error(),
			Data:    refs,
		}, Status: err.GetType().GetHttpCode()}, nil
	}
	return &plugin.ApiResourceOutput{Body: connection}, err
}

func (c *ConnectionApiHelper) merge(connection interface{}, body map[string]interface{}) errors.Error {
	connection = models.UnwrapObject(connection)
	if connectionValidator, ok := connection.(plugin.ConnectionValidator); ok {
		err := Decode(body, connection, nil)
		if err != nil {
			return err
		}
		return connectionValidator.ValidateConnection(connection, c.validator)
	}
	return Decode(body, connection, c.validator)
}

func (c *ConnectionApiHelper) save(connection interface{}, method func(entity interface{}, clauses ...dal.Clause) errors.Error) errors.Error {
	err := CallDB(method, connection)
	if err != nil {
		if c.db.IsDuplicationError(err) {
			return errors.BadInput.New("the connection name already exists")
		}
		return err
	}
	return nil
}

func (c *ConnectionApiHelper) getPluginSource() (plugin.PluginSource, errors.Error) {
	pluginMeta, _ := plugin.GetPlugin(c.pluginName)
	pluginSrc, ok := pluginMeta.(plugin.PluginSource)
	if !ok {
		return nil, errors.Default.New(fmt.Sprintf("plugin %s doesn't implement PluginSource", c.pluginName))
	}
	return pluginSrc, nil
}
