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

package migrationscripts

import (
	"fmt"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/domaininfo"
	"github.com/apache/incubator-devlake/core/plugin"
	"reflect"
	"strings"
)

var _ plugin.MigrationScript = (*addRawParamsTableForDomainScopes)(nil)

type addRawParamsTableForDomainScopes struct{}

type reflectedSlice any

func (script *addRawParamsTableForDomainScopes) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal().Begin()
	defer func() {
		if r := recover(); r != nil {
			err := db.Rollback()
			if err != nil {
				basicRes.GetLogger().Error(err, "error rolling back transaction")
			}
		}
	}()
	for _, domainModel := range domaininfo.GetDomainTablesInfo() {
		if scopeModel, ok := domainModel.(plugin.Scope); ok {
			scopes := toSlice(scopeModel)
			err := db.All(scopes)
			if err != nil {
				return err
			}
			sliceSize, err := forEach(scopes, func(scope plugin.Scope) errors.Error {
				reflectedScope := reflect.ValueOf(scope).Elem()
				domainId := reflectedScope.FieldByName("Id").String()
				derivedPlugin := strings.Split(domainId, ":")[0]
				if _, err = plugin.GetPlugin(derivedPlugin); err != nil {
					return errors.Default.New(fmt.Sprintf("could not infer the plugin in context from the domainId: %s in table: %s",
						derivedPlugin, scope.TableName()))
				}
				reflectedScope.FieldByName("RawDataTable").SetString(fmt.Sprintf("_raw_%s_scopes", derivedPlugin))
				return nil
			})
			if err != nil {
				return err
			}
			if sliceSize > 0 {
				err = db.Update(scopes)
				if err != nil {
					return err
				}
			}
		}
	}
	err := db.Commit()
	return err
}

func (*addRawParamsTableForDomainScopes) Version() uint64 {
	return 20230630000001
}

func (*addRawParamsTableForDomainScopes) Name() string {
	return "populated _raw_data_table column for domain scopes"
}

func toSlice(model any) reflectedSlice {
	ifc := reflect.New(reflect.SliceOf(reflect.TypeOf(model))).Interface()
	return ifc
}

func forEach[T any](r reflectedSlice, f func(x T) errors.Error) (int, errors.Error) {
	slice := reflect.ValueOf(r).Elem()
	count := slice.Len()
	for i := 0; i < count; i++ {
		err := f(slice.Index(i).Interface().(T))
		if err != nil {
			return count, err
		}
	}
	return count, nil
}
