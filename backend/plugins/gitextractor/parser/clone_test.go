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
	"errors"
	"fmt"
	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/impls/dalgorm"
	"github.com/stretchr/testify/assert"
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
	Name    string
}

func (ctx *testSubTaskContext) GetConfigReader() config.ConfigReader {
	cfg := config.GetConfig()
	return cfg
}

func (ctx *testSubTaskContext) GetConfig(name string) string {
	return config.GetConfig().GetString(name)
}

func (ctx *testSubTaskContext) GetLogger() log.Logger {
	return logger
}

func (ctx *testSubTaskContext) NestedLogger(name string) context.BasicRes {
	//TODO implement me
	panic("implement me")
}

func (ctx *testSubTaskContext) ReplaceLogger(logger log.Logger) context.BasicRes {
	//TODO implement me
	panic("implement me")
}

func (ctx *testSubTaskContext) GetDal() dal.Dal {
	//dsn := "mysql://root:admin@127.0.0.1:3306/lake?charset=utf8mb4&parseTime=True&loc=UTC"
	if runInLocal {
		dsn := "merico:merico@tcp(127.0.0.1:3306)/lake?charset=utf8mb4&parseTime=True&loc=Local"
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		return dalgorm.NewDalgorm(db)
	} else {
		panic("implement me")
	}
}

func (ctx *testSubTaskContext) GetName() string {
	return ctx.Name
}

func (ctx *testSubTaskContext) GetContext() gocontext.Context {
	return gocontext.Background()
}

func (ctx *testSubTaskContext) GetData() interface{} {
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

func (ctx *testSubTaskContext) TaskContext() plugin.TaskContext {
	//TODO implement me
	panic("implement me")
}

func Test_removePasswordFromErro1r(t *testing.T) {
	type args struct {
		err      error
		password string
		gitUrl   string
	}
	err := errors.New("some errors occurred when requesting http://devlakeuser:password@127.0.0.1/repos")
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "test-0",
			args: args{
				err:      err,
				password: "password",
				gitUrl:   "http://devlakeuser:password@127.0.0.1/repos",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err.Error() == "some errors occurred when requesting http://devlakeuser:********@127.0.0.1/repos"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = removePasswordFromError(tt.args.err, tt.args.password, tt.args.gitUrl)
			t.Logf("removePasswordFromError return err: %s", err)
			fmt.Printf("removePasswordFromError return err: %s\n", err)
			tt.wantErr(t, err)
		})
	}
}
