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
	"github.com/apache/incubator-devlake/core/plugin"
)

type Operation interface {
	Execute(dal.Dal) errors.Error
}

type ExecuteOperation struct {
	Sql         string  `json:"sql"`
	Dialect     *string `json:"dialect"`
	IgnoreError bool    `json:"ignore_error"`
}

func (o ExecuteOperation) Execute(dal dal.Dal) errors.Error {
	var err errors.Error
	if o.Dialect != nil {
		if dal.Dialect() == *o.Dialect {
			err = dal.Exec(o.Sql)
		}
	} else {
		err = dal.Exec(o.Sql)
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

func (o AddColumnOperation) Execute(dal dal.Dal) errors.Error {
	if dal.HasColumn(o.Table, o.Column) {
		err := dal.DropColumns(o.Table, o.Column)
		if err != nil {
			return err
		}
	}
	return dal.AddColumn(o.Table, o.Column, o.ColumnType)
}

type DropColumnOperation struct {
	Table  string `json:"table"`
	Column string `json:"column"`
}

func (o DropColumnOperation) Execute(dal dal.Dal) errors.Error {
	if dal.HasColumn(o.Table, o.Column) {
		return dal.DropColumns(o.Table, o.Column)
	}
	return nil
}

var _ Operation = (*DropColumnOperation)(nil)

type DropTableOperation struct {
	Table  string `json:"table"`
	Column string `json:"column"`
}

func (o DropTableOperation) Execute(dal dal.Dal) errors.Error {
	if dal.HasTable(o.Table) {
		return dal.DropTables(o.Table)
	}
	return nil
}

var _ Operation = (*DropTableOperation)(nil)

type RenameColumnOperation struct {
	Table   string `json:"table"`
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

func (o RenameColumnOperation) Execute(dal dal.Dal) errors.Error {
	if !dal.HasColumn(o.Table, o.OldName) {
		return nil
	}
	if dal.HasColumn(o.Table, o.NewName) {
		err := dal.DropColumns(o.Table, o.NewName)
		if err != nil {
			return err
		}
	}
	return dal.RenameColumn(o.Table, o.OldName, o.NewName)
}

var _ Operation = (*RenameTableOperation)(nil)

type RenameTableOperation struct {
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

func (o RenameTableOperation) Execute(dal dal.Dal) errors.Error {
	if !dal.HasTable(o.OldName) {
		return nil
	}
	if dal.HasTable(o.NewName) {
		err := dal.DropTables(o.NewName)
		if err != nil {
			return err
		}
	}
	return dal.RenameTable(o.OldName, o.NewName)
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
		err := json.Unmarshal(operationRaw, &operationMap)
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
		case "rename_column":
			operation = &RenameColumnOperation{}
		case "rename_table":
			operation = &RenameTableOperation{}
		default:
			return errors.BadInput.New("unsupported operation type: " + operationType)
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
	db := basicRes.GetDal()
	for _, operation := range s.operations {
		err := operation.Execute(db)
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
