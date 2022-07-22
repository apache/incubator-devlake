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

	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/logger"
	"github.com/apache/incubator-devlake/plugins/gitextractor/models"
	"github.com/apache/incubator-devlake/plugins/gitextractor/store"
	"github.com/apache/incubator-devlake/plugins/gitextractor/tasks"
	"github.com/apache/incubator-devlake/plugins/helper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	url := flag.String("url", "", "-url")
	proxy := flag.String("proxy", "", "-proxy")
	id := flag.String("id", "", "-id")
	user := flag.String("user", "", "-user")
	password := flag.String("password", "", "-password")
	output := flag.String("output", "", "-output")
	db := flag.String("db", "", "-db")
	flag.Parse()
	log := logger.Global.Nested("git extractor")
	var storage models.Store
	var err error
	if *url == "" {
		panic("url is missing")
	}
	if *id == "" {
		panic("id is missing")
	}
	if *output != "" {
		storage, err = store.NewCsvStore(*output)
		if err != nil {
			panic(err)
		}
	} else if *db != "" {
		database, err := gorm.Open(mysql.Open(*db))
		if err != nil {
			panic(err)
		}
		basicRes := helper.NewDefaultBasicRes(nil, log, database)
		storage = store.NewDatabase(basicRes, *url)
	} else {
		panic("either specify `-output` or `-db` argument as destination")
	}
	defer storage.Close()
	ctx := context.Background()
	subTaskCtx := helper.NewStandaloneSubTaskContext(
		ctx,
		config.GetConfig(),
		log,
		nil,
		"git extractor",
		nil,
	)
	repo, err := newGitRepo(log, storage, tasks.GitExtractorOptions{
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
