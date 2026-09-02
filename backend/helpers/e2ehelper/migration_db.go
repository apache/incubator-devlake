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

package e2ehelper

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/apache/devlake/core/config"
	"github.com/apache/devlake/core/plugin"
	"github.com/apache/devlake/core/runner"
	"github.com/apache/devlake/impls/dalgorm"
	"github.com/apache/devlake/impls/logruslog"
	"gorm.io/gorm"
)

// NewIsolatedMigrationDb creates a dedicated, empty database next to the one
// referenced by E2E_DB_URL (`<e2e-db>_<suffix>`) and returns a connection to it.
//
// Tests that execute the REAL migration scripts must not share the regular e2e
// database: the other plugin e2e tests seed/AutoMigrate tables (domain layer
// included) without recording anything in `_devlake_migration_history`, so a
// subsequent migration run fails with errors such as
// "Table 'cicd_pipeline_commits' already exists" when it tries to create or
// rename a table that is already there.
//
// The database is dropped again when the test finishes. If E2E_DB_URL is not
// set the test is skipped.
func NewIsolatedMigrationDb(t *testing.T, suffix string) *gorm.DB {
	cfg := config.GetConfig()
	e2eDbUrl := cfg.GetString("E2E_DB_URL")
	if e2eDbUrl == "" {
		t.Skip("E2E_DB_URL is not set; skipping migration schema check")
	}
	u, err := url.Parse(e2eDbUrl)
	if err != nil {
		t.Fatalf("unable to parse E2E_DB_URL: %v", err)
	}
	isolatedName := fmt.Sprintf("%s_%s", strings.TrimPrefix(u.Path, "/"), suffix)
	quotedName := quoteDbName(u.Scheme, isolatedName)

	gormConf := &gorm.Config{SkipDefaultTransaction: true}
	adminDb, err := runner.MakeDbConnection(e2eDbUrl, gormConf)
	if err != nil {
		t.Fatalf("unable to connect to E2E_DB_URL: %v", err)
	}
	if err = adminDb.Exec("DROP DATABASE IF EXISTS " + quotedName).Error; err != nil {
		t.Fatalf("unable to drop leftover database %s: %v", isolatedName, err)
	}
	if err = adminDb.Exec("CREATE DATABASE " + quotedName).Error; err != nil {
		t.Fatalf("unable to create database %s: %v", isolatedName, err)
	}
	closeDb(adminDb)

	isolatedUrl := *u
	isolatedUrl.Path = "/" + isolatedName
	db, err := runner.MakeDbConnection(isolatedUrl.String(), gormConf)
	if err != nil {
		t.Fatalf("unable to connect to %s: %v", isolatedName, err)
	}

	// migration scripts and models read DB_URL from the global config, keep it
	// consistent with the connection we hand out and restore it afterwards.
	previousDbUrl := cfg.GetString("DB_URL")
	cfg.Set("DB_URL", isolatedUrl.String())

	// Some migration scripts refuse to run without an encryption secret
	// (e.g. jira 20220716: "jira v0.11 invalid encKey"). CI does not
	// necessarily provide one, so fall back to a deterministic test value.
	// dalgorm.Init registers the `encdec` GORM serializer used by connection
	// models - without it migrations fail with "invalid serializer type encdec"
	// (runner.CreateBasicRes does not register it, only CreateAppBasicRes does).
	if cfg.GetString(plugin.EncodeKeyEnvStr) == "" {
		cfg.Set(plugin.EncodeKeyEnvStr, "devlake-e2e-test-encryption-secret")
	}
	dalgorm.Init(cfg.GetString(plugin.EncodeKeyEnvStr))

	t.Cleanup(func() {
		cfg.Set("DB_URL", previousDbUrl)
		closeDb(db)
		cleanupDb, cleanupErr := runner.MakeDbConnection(e2eDbUrl, gormConf)
		if cleanupErr != nil {
			t.Logf("unable to connect for dropping %s: %v", isolatedName, cleanupErr)
			return
		}
		defer closeDb(cleanupDb)
		if dropErr := cleanupDb.Exec("DROP DATABASE IF EXISTS " + quotedName).Error; dropErr != nil {
			t.Logf("unable to drop database %s: %v", isolatedName, dropErr)
		}
	})

	logruslog.Global.Info("running migrations against isolated database %s", isolatedName)
	return db
}

func quoteDbName(scheme string, name string) string {
	// database names are derived from E2E_DB_URL + a constant suffix, but quote
	// them anyway to stay safe with reserved words.
	if strings.EqualFold(scheme, "mysql") {
		return "`" + strings.ReplaceAll(name, "`", "") + "`"
	}
	return `"` + strings.ReplaceAll(name, `"`, "") + `"`
}

func closeDb(db *gorm.DB) {
	if sqlDb, err := db.DB(); err == nil {
		_ = sqlDb.Close()
	}
}
