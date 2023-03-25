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

package bridge

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/utils"
	"github.com/apache/incubator-devlake/impls/logruslog"
	"github.com/spf13/viper"
)

func Bootstrap(cfg *viper.Viper, port int) errors.Error {
	scriptPath := cfg.GetString("REMOTE_PLUGINS_STARTUP_PATH")
	if scriptPath == "" {
		return errors.BadInput.New(fmt.Sprintf("missing env key: %s", "REMOTE_PLUGINS_STARTUP_PATH"))
	}
	absScriptPath := scriptPath
	if !filepath.IsAbs(scriptPath) {
		workingDir, err := errors.Convert01(os.Getwd())
		if err != nil {
			return err
		}
		absScriptPath = filepath.Join(workingDir, scriptPath)
	}
	logruslog.Global.Info("Resolved remote plugins script path: %s", absScriptPath)
	cmd := exec.Command(absScriptPath, fmt.Sprintf("http://127.0.0.1:%d", port)) //expects the plugins to live on the same host
	cmd.Dir = filepath.Dir(absScriptPath)
	result, err := utils.RunProcess(cmd, &utils.RunProcessOptions{
		OnStdout: func(b []byte) {
			logruslog.Global.Info(string(b))
		},
		OnStderr: func(b []byte) {
			logruslog.Global.Error(nil, string(b))
		},
	})
	if err != nil {
		return err
	}
	if result.GetError() != nil {
		logruslog.Global.Error(result.GetError(), "error occurred bootstrapping remote plugins")
	}
	return nil
}
