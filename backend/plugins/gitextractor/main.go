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

package main

import (
	"context"
	"flag"
	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/runner"
	contextimpl "github.com/apache/incubator-devlake/impls/context"
	"github.com/apache/incubator-devlake/impls/dalgorm"
	"github.com/apache/incubator-devlake/impls/logruslog"
	"github.com/apache/incubator-devlake/plugins/gitextractor/impl"
	"github.com/apache/incubator-devlake/plugins/gitextractor/models"
	"github.com/apache/incubator-devlake/plugins/gitextractor/store"
	"github.com/apache/incubator-devlake/plugins/gitextractor/tasks"
)

// PluginEntry is a variable exported for Framework to search and load
var PluginEntry impl.GitExtractor //nolint

func main() {
	url := flag.String("url", "", "-url")
	proxy := flag.String("proxy", "", "-proxy")
	id := flag.String("id", "", "-id")
	user := flag.String("user", "", "-user")
	password := flag.String("password", "", "-password")
	output := flag.String("output", "", "-output")
	dbUrl := flag.String("db", "", "-db")
	flag.Parse()
	cfg := config.GetConfig()
	logger := logruslog.Global.Nested("git extractor")
	var storage models.Store
	var err error
	if *url == "" {
		panic("url is missing")
	}
	if *id == "" {
		panic("id is missing")
	}
	db, err := runner.NewGormDb(cfg, logger)
	if err != nil {
		panic(err)
	}
	basicRes := contextimpl.NewDefaultBasicRes(cfg, logger, dalgorm.NewDalgorm(db))
	if *output != "" {
		storage, err = store.NewCsvStore(*output)
		if err != nil {
			panic(err)
		}
	} else if *dbUrl != "" {
		cfg.Set("DB_URL", *dbUrl)
	}
	// If we didn't specify output or dburl, we will use db by default
	if storage == nil {
		storage = store.NewDatabase(basicRes, *id)
	}
	defer storage.Close()
	ctx := context.Background()
	subTaskCtx := contextimpl.NewStandaloneSubTaskContext(
		ctx,
		basicRes,
		"git extractor",
		nil,
	)
	repo, err := impl.NewGitRepo(logger, storage, tasks.GitExtractorOptions{
		RepoId:   *id,
		Url:      *url,
		User:     *user,
		Password: *password,
		Proxy:    *proxy,
	})
	if err != nil {
		panic(err)
	}
	defer repo.Close()
	if err = repo.CollectAll(subTaskCtx); err != nil {
		panic(err)
	}
}
