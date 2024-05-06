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
	"fmt"
	"reflect"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
)

var _ srvhelper.ScopeModelInfo = (*RemoteScopeModelInfo)(nil)

type RemoteScopeModelInfo struct {
	*RemoteModelInfo
}

func NewRemoteScopeModelInfo[ParentType any](di *DynamicModelInfo) (*RemoteScopeModelInfo, errors.Error) {
	mi, err := GenerateRemoteModelInfo[ParentType](di)
	if err != nil {
		return nil, err
	}
	return &RemoteScopeModelInfo{RemoteModelInfo: mi}, nil
}

// GetConnectionId implements srvhelper.ScopeModelInfo.
func (r *RemoteScopeModelInfo) GetConnectionId(scope any) uint64 {
	return reflect.ValueOf(scope).Elem().FieldByName("ConnectionId").Uint()
}

// GetScopeConfigId implements srvhelper.ScopeModelInfo.
func (r *RemoteScopeModelInfo) GetScopeConfigId(scope any) uint64 {
	return reflect.ValueOf(scope).Elem().FieldByName("ScopeConfigId").Uint()
}

// GetScopeId implements srvhelper.ScopeModelInfo.
func (r *RemoteScopeModelInfo) GetScopeId(scope any) string {
	return fmt.Sprintf("%v", reflect.ValueOf(scope).Elem().FieldByName("ID").Interface())
}

// GetScopeParams implements srvhelper.ScopeModelInfo.
func (r *RemoteScopeModelInfo) GetScopeParams(scope any) interface{} {
	return &ApiParams{
		ConnectionId: r.GetConnectionId(scope),
		ScopeId:      r.GetScopeId(scope),
	}
}