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

package parser

import (
	gocontext "context"
	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/impls/dalgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func Test_setCloneProgress(t *testing.T) {
	type args struct {
		subTaskCtx        plugin.SubTaskContext
		cloneProgressInfo string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test-0",
			args: args{
				subTaskCtx: &testSubTaskContext{},
				cloneProgressInfo: `
					Enumerating objects: 103, done.
					Counting objects: 100% (103/103), done.
					Compressing objects: 100% (81/81), done.
				`,
			},
		},
		{
			name: "test-1",
			args: args{
				subTaskCtx: &testSubTaskContext{},
				cloneProgressInfo: `
					Enumerating objects: 103, done.
					Counting objects: 100% (103/103), done.
				`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setCloneProgress(tt.args.subTaskCtx, tt.args.cloneProgressInfo)
		})
	}
}

type testSubTaskContext struct {
	current int
	total   int
}

func (testSubTaskContext) GetConfigReader() config.ConfigReader {
	//TODO implement me
	cfg := config.GetConfig()
	return cfg
}

func (testSubTaskContext) GetConfig(name string) string {
	//TODO implement me
	panic("implement me")
}

func (testSubTaskContext) GetLogger() log.Logger {
	//TODO implement me
	panic("implement me")
}

func (testSubTaskContext) NestedLogger(name string) context.BasicRes {
	//TODO implement me
	panic("implement me")
}

func (testSubTaskContext) ReplaceLogger(logger log.Logger) context.BasicRes {
	//TODO implement me
	panic("implement me")
}

func (testSubTaskContext) GetDal() dal.Dal {
	//dsn := "mysql://root:admin@127.0.0.1:3306/lake?charset=utf8mb4&parseTime=True&loc=UTC"
	dsn := "root:admin@tcp(127.0.0.1:3306)/lake?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return dalgorm.NewDalgorm(db)
}

func (testSubTaskContext) GetName() string {
	//TODO implement me
	panic("implement me")
}

func (testSubTaskContext) GetContext() gocontext.Context {
	//TODO implement me
	return gocontext.Background()
}

func (testSubTaskContext) GetData() interface{} {
	//TODO implement me
	panic("implement me")
}

func (ctx *testSubTaskContext) SetProgress(current int, total int) {
	ctx.current = current
	ctx.total = total
}

func (ctx *testSubTaskContext) IncProgress(quantity int) {
	ctx.current += quantity
	ctx.total += quantity
}

func (testSubTaskContext) TaskContext() plugin.TaskContext {
	//TODO implement me
	panic("implement me")
}
