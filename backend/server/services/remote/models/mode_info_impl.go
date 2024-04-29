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
	"reflect"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
)

var _ srvhelper.ModelInfo = (*RemoteModelInfo)(nil)

type RemoteModelInfo struct {
	modelName string
	modelType reflect.Type
	tableName string
}

func NewRemoteModelInfo(modelName string, modelType reflect.Type, tableName string) *RemoteModelInfo {
	return &RemoteModelInfo{
		modelName: modelName,
		modelType: modelType,
		tableName: tableName,
	}
}

func GenerateRemoteModelInfo[ParentType any](di *DynamicModelInfo) (*RemoteModelInfo, errors.Error) {
	modelName := di.JsonSchema["title"].(string)
	parentType := reflect.TypeOf((*ParentType)(nil)).Elem()
	modelType, err := GenerateStructType(di.JsonSchema, parentType)
	if err != nil {
		return nil, err
	}
	tableName := di.TableName
	return NewRemoteModelInfo(modelName, modelType, tableName), nil
}

// ModelName implements srvhelper.ModelInfo.
func (r *RemoteModelInfo) ModelName() string {
	return r.modelName
}

// New implements srvhelper.ModelInfo.
func (r *RemoteModelInfo) New() any {
	return reflect.New(r.modelType).Interface()
}

// NewSlice implements srvhelper.ModelInfo.
func (r *RemoteModelInfo) NewSlice() any {
	return reflect.New(reflect.SliceOf(reflect.PointerTo(r.modelType)))
}

// TableName implements srvhelper.ModelInfo.
func (r *RemoteModelInfo) TableName() string {
	return r.tableName
}
