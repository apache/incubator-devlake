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

package srvhelper

import (
	"fmt"
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/helpers/dbhelper"
	"github.com/go-playground/validator/v10"
)

type CustomValidator interface {
	CustomValidate(entity interface{}, validate *validator.Validate) errors.Error
}

type ModelInfo interface {
	New() any
	NewSlice() any
	ModelName() string
	TableName() string
}

type AnyModelSrvHelper struct {
	ModelInfo
	basicRes      context.BasicRes
	log           log.Logger
	db            dal.Dal
	validator     *validator.Validate
	pk            []dal.ColumnMeta
	pkWhere       string
	pkCount       int
	searchColumns []string
}

func NewAnyModelSrvHelper(basicRes context.BasicRes, modelInfo ModelInfo, searchColumns []string) *AnyModelSrvHelper {
	db := basicRes.GetDal()
	if db == nil {
		db = basicRes.GetDal()
	}
	pk := errors.Must1(dal.GetPrimarykeyColumns(db, modelInfo.TableName()))
	pkWhere := ""
	for _, col := range pk {
		if pkWhere != "" {
			pkWhere += " AND "
		}
		pkWhere += fmt.Sprintf("%s = ? ", col.Name())
	}

	return &AnyModelSrvHelper{
		ModelInfo:     modelInfo,
		basicRes:      basicRes,
		log:           basicRes.GetLogger().Nested(fmt.Sprintf("%s_dal", modelInfo.ModelName())),
		db:            db,
		validator:     validator.New(),
		pk:            pk,
		pkWhere:       pkWhere,
		pkCount:       len(pk),
		searchColumns: searchColumns,
	}
}

func (srv *AnyModelSrvHelper) ValidateModel(model any) errors.Error {
	// the model can validate itself
	if customValidator, ok := (interface{}(model)).(CustomValidator); ok {
		return customValidator.CustomValidate(model, srv.validator)
	}
	// basic validator
	if e := srv.validator.Struct(model); e != nil {
		return errors.BadInput.Wrap(e, "validation faild")
	}
	return nil
}

// Create validates given model and insert it into database if validation passed
func (srv *AnyModelSrvHelper) CreateAny(model any) errors.Error {
	err := srv.ValidateModel(model)
	if err != nil {
		return err
	}
	err = srv.db.Create(model)
	if err != nil {
		if srv.db.IsDuplicationError(err) {
			return errors.Default.New("The name of the current scope config is duplicated. Please modify it before saving.")
		}
		return err
	}
	return err
}

// Update validates given model and update it into database if validation passed
func (srv *AnyModelSrvHelper) UpdateAny(model any) errors.Error {
	err := srv.ValidateModel(model)
	if err != nil {
		if srv.db.IsDuplicationError(err) {
			return errors.Default.New("The name of the current scope config is duplicated. Please modify it before saving.")
		}
		return err
	}
	return srv.db.Update(model, dal.From(srv.ModelInfo.TableName()))
}

// CreateOrUpdate validates given model and insert or update it into database if validation passed
func (srv *AnyModelSrvHelper) CreateOrUpdateAny(model any) errors.Error {
	err := srv.ValidateModel(model)
	if err != nil {
		return err
	}
	return srv.db.CreateOrUpdate(model)
}

// DeleteModel deletes given model from database
func (srv *AnyModelSrvHelper) DeleteModelAny(model any) errors.Error {
	return srv.db.Delete(model)
}

// FindByPk returns model with given primary key from database
func (srv *AnyModelSrvHelper) FindByPkAny(pk ...interface{}) (any, errors.Error) {
	if len(pk) != srv.pkCount {
		return nil, errors.BadInput.New("invalid primary key")
	}
	model := srv.New()
	err := srv.db.First(model, dal.From(srv.TableName()), dal.Where(srv.pkWhere, pk...))
	if err != nil {
		if srv.db.IsErrorNotFound(err) {
			return nil, errors.NotFound.Wrap(err, fmt.Sprintf("%s not found", srv.ModelName()))
		}
		return nil, err
	}
	return model, nil
}

// GetAll returns all models from database
func (srv *AnyModelSrvHelper) QueryAllAny() (any, errors.Error) {
	array := srv.NewSlice()
	return array, srv.db.All(&array, dal.From(srv.ModelInfo.TableName()))
}

func (srv *AnyModelSrvHelper) QueryPageAny(pagination *Pagination, query ...dal.Clause) (any, int64, errors.Error) {
	query = append(query, dal.From(srv.TableName()))
	// process keyword
	searchTerm := pagination.SearchTerm
	if searchTerm != "" && len(srv.searchColumns) > 0 {
		sql := ""
		value := "%" + searchTerm + "%"
		values := make([]interface{}, len(srv.searchColumns))
		for i, field := range srv.searchColumns {
			if sql != "" {
				sql += " OR "
			}
			sql += fmt.Sprintf("%s LIKE ?", field)
			values[i] = value
		}
		sql = fmt.Sprintf("(%s)", sql)
		query = append(query,
			dal.Where(sql, values...),
		)
	}
	count, err := srv.db.Count(query...)
	if err != nil {
		return nil, 0, err
	}
	query = append(query, dal.Limit(pagination.GetLimit()), dal.Offset(pagination.GetOffset()))
	array := srv.NewSlice()
	return array, count, srv.db.All(&array, query...)
}

func (srv *AnyModelSrvHelper) NoRunningPipeline(fn func(tx dal.Transaction) errors.Error, tablesToLock ...*dal.LockTable) (err errors.Error) {
	// make sure no pipeline is running
	tablesToLock = append(tablesToLock, &dal.LockTable{Table: "_devlake_pipelines", Exclusive: true})
	txHelper := dbhelper.NewTxHelper(srv.basicRes, &err)
	defer txHelper.End()
	tx := txHelper.Begin()
	err = txHelper.LockTablesTimeout(2*time.Second, tablesToLock)
	if err != nil {
		err = errors.Conflict.Wrap(err, "lock pipelines table timedout")
		return
	}
	count := errors.Must1(tx.Count(
		dal.From("_devlake_pipelines"),
		dal.Where("status = ?", models.TASK_RUNNING),
	))
	if count > 0 {
		err = errors.Conflict.New("at least one pipeline is running")
		return
	}
	// time.Sleep(1 * time.Minute) # uncomment this line if you were to verify pipelines get blocked while deleting data
	// creating a nested transaction to avoid mysql complaining about table(s) NOT being locked
	nextedTxHelper := dbhelper.NewTxHelper(srv.basicRes, &err)
	defer nextedTxHelper.End()
	nestedTX := nextedTxHelper.Begin()
	err = fn(nestedTX)
	return
}
