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

package models

import (
	"encoding/json"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type Operation interface {
	Execute(basicRes context.BasicRes) errors.Error
}

type ExecuteOperation struct {
	Sql         string  `json:"sql"`
	Dialect     *string `json:"dialect"`
	IgnoreError bool    `json:"ignore_error"`
}

func (o ExecuteOperation) Execute(basicRes context.BasicRes) errors.Error {
	var err errors.Error
	db := basicRes.GetDal()
	if o.Dialect != nil {
		if db.Dialect() == *o.Dialect {
			err = db.Exec(o.Sql)
		}
	} else {
		err = db.Exec(o.Sql)
	}
	if o.IgnoreError {
		return nil
	}
	return err
}

var _ Operation = (*ExecuteOperation)(nil)

type AddColumnOperation struct {
	Table      string         `json:"table"`
	Column     string         `json:"column"`
	ColumnType dal.ColumnType `json:"column_type"`
}

func (o AddColumnOperation) Execute(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if db.HasColumn(o.Table, o.Column) {
		return nil
	}
	return db.AddColumn(o.Table, o.Column, o.ColumnType)
}

type DropColumnOperation struct {
	Table  string `json:"table"`
	Column string `json:"column"`
}

func (o DropColumnOperation) Execute(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if db.HasColumn(o.Table, o.Column) {
		return db.DropColumns(o.Table, o.Column)
	}
	return nil
}

var _ Operation = (*DropColumnOperation)(nil)

type DropTableOperation struct {
	Table  string `json:"table"`
	Column string `json:"column"`
}

func (o DropTableOperation) Execute(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if db.HasTable(o.Table) {
		return db.DropTables(o.Table)
	}
	return nil
}

var _ Operation = (*DropTableOperation)(nil)

type RenameTableOperation struct {
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

func (o RenameTableOperation) Execute(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if !db.HasTable(o.OldName) {
		return nil
	}
	return db.RenameTable(o.OldName, o.NewName)
}

type CreateTableOperation struct {
	ModelInfo *DynamicModelInfo `json:"model_info"`
}

func (o CreateTableOperation) Execute(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if db.HasTable(o.ModelInfo.TableName) {
		basicRes.GetLogger().Warn(nil, "table %s already exists. It won't be created.", o.ModelInfo.TableName)
		return nil
	}
	model, err := o.ModelInfo.LoadDynamicTabler(common.NoPKModel{})
	if err != nil {
		return err
	}
	// uncomment to debug "modelDump" as needed
	modelDump := models.DumpInfo(model.New())
	_ = modelDump
	err = api.CallDB(db.AutoMigrate, model.New())
	if err != nil {
		return err
	}
	return nil
}

var _ Operation = (*RenameTableOperation)(nil)

type RemoteMigrationScript struct {
	operations []Operation
	version    uint64
	name       string
}

type rawRemoteMigrationScript struct {
	Operations []json.RawMessage `json:"operations"`
	Version    uint64            `json:"version"`
	Name       string            `json:"name"`
}

func (s *RemoteMigrationScript) UnmarshalJSON(data []byte) error {
	var rawScript rawRemoteMigrationScript
	err := json.Unmarshal(data, &rawScript)
	if err != nil {
		return err
	}
	s.version = rawScript.Version
	s.name = rawScript.Name
	s.operations = make([]Operation, len(rawScript.Operations))
	for i, operationRaw := range rawScript.Operations {
		operationMap := make(map[string]interface{})
		err = json.Unmarshal(operationRaw, &operationMap)
		if err != nil {
			return err
		}
		operationType := operationMap["type"].(string)
		var operation Operation
		switch operationType {
		case "execute":
			operation = &ExecuteOperation{}
		case "add_column":
			operation = &AddColumnOperation{}
		case "drop_column":
			operation = &DropColumnOperation{}
		case "drop_table":
			operation = &DropTableOperation{}
		case "rename_table":
			operation = &RenameTableOperation{}
		case "create_table":
			operation = &CreateTableOperation{}
		default:
			return errors.BadInput.New("unsupported operation type")
		}
		err = json.Unmarshal(operationRaw, operation)
		if err != nil {
			return err
		}
		s.operations[i] = operation
	}
	return nil
}

func (s *RemoteMigrationScript) Up(basicRes context.BasicRes) errors.Error {
	for _, operation := range s.operations {
		err := operation.Execute(basicRes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *RemoteMigrationScript) Version() uint64 {
	return s.version
}

func (s *RemoteMigrationScript) Name() string {
	return s.name
}

var _ plugin.MigrationScript = (*RemoteMigrationScript)(nil)
