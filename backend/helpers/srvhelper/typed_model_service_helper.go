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

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
)

type GenericModelInfo[M dal.Tabler] struct {
	modelName string
	tableName string
}

func NewGenericModelInfo[M dal.Tabler]() *GenericModelInfo[M] {
	var m M
	return &GenericModelInfo[M]{
		modelName: fmt.Sprintf("%T", m),
		tableName: m.TableName(),
	}
}

func (info *GenericModelInfo[M]) ModelName() string {
	return info.modelName
}

func (info *GenericModelInfo[M]) TableName() string {
	return info.tableName
}

func (info *GenericModelInfo[M]) New() interface{} {
	return new(M)
}

func (info *GenericModelInfo[M]) NewSlice() interface{} {
	return make([]*M, 0)
}

type ModelSrvHelper[M dal.Tabler] struct {
	*AnyModelSrvHelper
}

func NewModelSrvHelper[M dal.Tabler](basicRes context.BasicRes, searchColumns []string) *ModelSrvHelper[M] {
	return &ModelSrvHelper[M]{
		NewAnyModelSrvHelper(basicRes, NewGenericModelInfo[M](), searchColumns),
	}
}

// Create validates given model and insert it into database if validation passed
func (srv *ModelSrvHelper[M]) Create(model *M) errors.Error {
	return srv.CreateAny(model)
}

// Update validates given model and update it into database if validation passed
func (srv *ModelSrvHelper[M]) Update(model *M) errors.Error {
	return srv.UpdateAny(model)
}

// CreateOrUpdate validates given model and insert or update it into database if validation passed
func (srv *ModelSrvHelper[M]) CreateOrUpdate(model *M) errors.Error {
	return srv.CreateOrUpdateAny(model)
}

// DeleteModel deletes given model from database
func (srv *ModelSrvHelper[M]) DeleteModel(model *M) errors.Error {
	return srv.DeleteModelAny(model)
}

// FindByPk returns model with given primary key from database
func (srv *ModelSrvHelper[M]) FindByPk(pk ...interface{}) (*M, errors.Error) {
	model, err := srv.FindByPkAny(pk...)
	return model.(*M), err
}

// GetAll returns all models from database
func (srv *ModelSrvHelper[M]) GetAll() ([]*M, errors.Error) {
	array, err := srv.QueryAllAny()
	return array.([]*M), err
}

func (srv *ModelSrvHelper[M]) GetPage(pagination *Pagination, query ...dal.Clause) ([]*M, int64, errors.Error) {
	array, count, err := srv.QueryPageAny(pagination, query...)
	return array.([]*M), count, err
}
