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

package helper

import (
	"bytes"
	"context"
	"encoding/json"
	goerror "errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/apache/incubator-devlake/core/config"
	corectx "github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/migration"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/runner"
	contextimpl "github.com/apache/incubator-devlake/impls/context"
	"github.com/apache/incubator-devlake/impls/dalgorm"
	"github.com/apache/incubator-devlake/impls/logruslog"
	"github.com/apache/incubator-devlake/server/api"
	remotePlugin "github.com/apache/incubator-devlake/server/services/remote/plugin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	throwawayDir           string
	initService            = new(sync.Once)
	dbTruncationExclusions = []string{
		"_devlake_migration_history",
		"_devlake_locking_stub",
		"_devlake_locking_history",
	}
)

func init() {
	tempDir, err := errors.Convert01(os.MkdirTemp("", "devlake_test"+"_*"))
	if err != nil {
		panic(err)
	}
	throwawayDir = tempDir
}

// DevlakeClient FIXME
type (
	DevlakeClient struct {
		Endpoint string
		db       *gorm.DB
		log      log.Logger
		cfg      *viper.Viper
		testCtx  *testing.T
		basicRes corectx.BasicRes
		timeout  time.Duration
	}
	LocalClientConfig struct {
		ServerPort           uint
		DbURL                string
		CreateServer         bool
		DropDb               bool
		TruncateDb           bool
		Plugins              map[string]plugin.PluginMeta
		AdditionalMigrations func() []plugin.MigrationScript
		Timeout              time.Duration
	}
	RemoteClientConfig struct {
		Endpoint string
	}
)

// ConnectRemoteServer returns a client to an existing server based on the config
func ConnectRemoteServer(t *testing.T, sbConfig *RemoteClientConfig) *DevlakeClient {
	return &DevlakeClient{
		Endpoint: sbConfig.Endpoint,
		db:       nil,
		log:      nil,
		testCtx:  t,
	}
}

// ConnectLocalServer spins up a local server from the config and returns a client connected to it
func ConnectLocalServer(t *testing.T, sbConfig *LocalClientConfig) *DevlakeClient {
	t.Helper()
	fmt.Printf("Using test temp directory: %s\n", throwawayDir)
	logger := logruslog.Global.Nested("test")
	cfg := config.GetConfig()
	cfg.Set("DB_URL", sbConfig.DbURL)
	db, err := runner.NewGormDb(cfg, logger)
	require.NoError(t, err)
	t.Cleanup(func() {
		d, err := db.DB()
		require.NoError(t, err)
		require.NoError(t, d.Close())
	})
	addr := fmt.Sprintf("http://localhost:%d", sbConfig.ServerPort)
	d := &DevlakeClient{
		Endpoint: addr,
		db:       db,
		log:      logger,
		cfg:      cfg,
		basicRes: contextimpl.NewDefaultBasicRes(cfg, logger, dalgorm.NewDalgorm(db)),
		testCtx:  t,
		timeout:  sbConfig.Timeout,
	}
	d.configureEncryption()
	d.initPlugins(sbConfig)
	if sbConfig.DropDb || sbConfig.TruncateDb {
		d.prepareDB(sbConfig)
	}
	if sbConfig.CreateServer {
		cfg.Set("PORT", sbConfig.ServerPort)
		cfg.Set("PLUGIN_DIR", throwawayDir)
		cfg.Set("LOGGING_DIR", throwawayDir)
		go func() {
			initService.Do(api.CreateApiService)
		}()
	}
	req, err2 := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/proceed-db-migration", addr), nil)
	require.NoError(t, err2)
	d.forceSendHttpRequest(20, req, func(err errors.Error) bool {
		e := err.Unwrap()
		return goerror.Is(e, syscall.ECONNREFUSED)
	})
	d.runMigrations(sbConfig)
	return d
}

// SetTimeout override the timeout of api requests
func (d *DevlakeClient) SetTimeout(timeout time.Duration) {
	d.timeout = timeout
}

// RunPlugin manually execute a plugin directly (local server only)
func (d *DevlakeClient) RunPlugin(ctx context.Context, pluginName string, pluginTask plugin.PluginTask, options map[string]interface{}, subtaskNames ...string) errors.Error {
	if len(subtaskNames) == 0 {
		subtaskNames = GetSubtaskNames(pluginTask.SubTaskMetas()...)
	}
	optionsJson, err := json.Marshal(options)
	if err != nil {
		return errors.Convert(err)
	}
	subtasksJson, err := json.Marshal(subtaskNames)
	if err != nil {
		return errors.Convert(err)
	}
	task := &models.Task{
		Plugin:   pluginName,
		Options:  string(optionsJson),
		Subtasks: subtasksJson,
	}
	return runner.RunPluginSubTasks(
		ctx,
		d.basicRes,
		task,
		pluginTask,
		nil,
	)
}

func (d *DevlakeClient) configureEncryption() {
	v := config.GetConfig()
	encKey := v.GetString(plugin.EncodeKeyEnvStr)
	if encKey == "" {
		// Randomly generate a bunch of encryption keys and set them to config
		encKey = plugin.RandomEncKey()
		v.Set(plugin.EncodeKeyEnvStr, encKey)
		err := config.WriteConfig(v)
		if err != nil {
			panic(err)
		}
	}
}

func (d *DevlakeClient) forceSendHttpRequest(retries uint, req *http.Request, onError func(err errors.Error) bool) {
	d.testCtx.Helper()
	for {
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			if !onError(errors.Default.WrapRaw(err)) {
				require.NoError(d.testCtx, err)
			}
		} else {
			if res.StatusCode != http.StatusOK {
				panic(fmt.Sprintf("received HTTP status %d", res.StatusCode))
			}
			return
		}
		retries--
		if retries == 0 {
			panic("retry limit exceeded")
		}
		fmt.Printf("retrying http call to %s\n", req.URL.String())
		time.Sleep(1 * time.Second)
	}
}

func (d *DevlakeClient) initPlugins(sbConfig *LocalClientConfig) {
	remotePlugin.Init(d.basicRes)
	d.testCtx.Helper()
	if sbConfig.Plugins != nil {
		for name, p := range sbConfig.Plugins {
			require.NoError(d.testCtx, plugin.RegisterPlugin(name, p))
		}
	}
	for _, p := range plugin.AllPlugins() {
		if pi, ok := p.(plugin.PluginInit); ok {
			err := pi.Init(d.basicRes)
			require.NoError(d.testCtx, err)
		}
	}
}

func (d *DevlakeClient) runMigrations(sbConfig *LocalClientConfig) {
	d.testCtx.Helper()
	basicRes := contextimpl.NewDefaultBasicRes(d.cfg, d.log, dalgorm.NewDalgorm(d.db))
	getMigrator := func() plugin.Migrator {
		migrator, err := migration.NewMigrator(basicRes)
		require.NoError(d.testCtx, err)
		return migrator
	}
	{
		migrator := getMigrator()
		for pluginName, pluginInst := range sbConfig.Plugins {
			if migratable, ok := pluginInst.(plugin.PluginMigration); ok {
				migrator.Register(migratable.MigrationScripts(), pluginName)
			}
		}
		require.NoError(d.testCtx, migrator.Execute())
	}
	{
		migrator := getMigrator()
		if sbConfig.AdditionalMigrations != nil {
			scripts := sbConfig.AdditionalMigrations()
			migrator.Register(scripts, "extra migrations")
		}
		require.NoError(d.testCtx, migrator.Execute())
	}
}

func (d *DevlakeClient) prepareDB(cfg *LocalClientConfig) {
	d.testCtx.Helper()
	migrator := d.db.Migrator()
	tables, err := migrator.GetTables()
	require.NoError(d.testCtx, err)
	if cfg.DropDb {
		d.log.Info("Dropping %d tables", len(tables))
		var tablesRaw []any
		for _, table := range tables {
			tablesRaw = append(tablesRaw, table)
		}
		err = migrator.DropTable(tablesRaw...)
		require.NoError(d.testCtx, err)
	} else if cfg.TruncateDb {
		d.log.Info("Truncating %d tables", len(tables)-len(dbTruncationExclusions))
		for _, table := range tables {
			excluded := false
			for _, exclusion := range dbTruncationExclusions {
				if exclusion == table {
					excluded = true
					break
				}
			}
			if !excluded {
				err = d.db.Exec("DELETE FROM " + table).Error
				require.NoError(d.testCtx, err)
			}
		}
	}
}

func sendHttpRequest[Res any](t *testing.T, timeout time.Duration, debug debugInfo, httpMethod string, endpoint string, body any) Res {
	t.Helper()
	b := ToJson(body)
	if debug.print {
		coloredPrintf("calling:\n\t%s %s\nwith:\n%s\n", httpMethod, endpoint, string(ToCleanJson(debug.inlineJson, body)))
	}
	timer := time.After(timeout)
	request, err := http.NewRequest(httpMethod, endpoint, bytes.NewReader(b))
	require.NoError(t, err)
	request.Close = true
	request.Header.Add("Content-Type", "application/json")
	for {
		response, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		if timeout > 0 {
			select {
			case <-timer:
			default:
				if response.StatusCode >= 300 {
					require.NoError(t, response.Body.Close())
					response.Close = true
					time.Sleep(1 * time.Second)
					continue
				}
			}
		}
		require.True(t, response.StatusCode < 300, "unexpected http status code: %d", response.StatusCode)
		var result Res
		b, _ = io.ReadAll(response.Body)
		require.NoError(t, json.Unmarshal(b, &result))
		if debug.print {
			coloredPrintf("result: %s\n", ToCleanJson(debug.inlineJson, b))
		}
		require.NoError(t, response.Body.Close())
		response.Close = true
		return result
	}
}

func coloredPrintf(msg string, args ...any) {
	msg = fmt.Sprintf(msg, args...)
	colorifier := "\033[1;33m%+v\033[0m" //yellow
	fmt.Printf(colorifier, msg)
}

type debugInfo struct {
	print      bool
	inlineJson bool
}
