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

package tasks

import (
	"bufio"
	"encoding/json"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/apache/incubator-devlake/errors"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/spf13/viper"
)

func DbtConverter(taskCtx core.SubTaskContext) errors.Error {
	log := taskCtx.GetLogger()
	taskCtx.SetProgress(0, -1)
	data := taskCtx.GetData().(*DbtTaskData)
	models := data.Options.SelectedModels
	projectPath := data.Options.ProjectPath
	projectName := data.Options.ProjectName
	projectTarget := data.Options.ProjectTarget
	projectVars := data.Options.ProjectVars
	args := data.Options.Args
	err := errors.Convert(os.Chdir(projectPath))
	if err != nil {
		return err
	}
	_, err = errors.Convert01(os.Stat("profiles.yml"))
	// if profiles.yml not exist, create it manually
	if err != nil {
		dbUrl := taskCtx.GetConfig("DB_URL")
		u, err := errors.Convert01(url.Parse(dbUrl))
		if err != nil {
			return err
		}
		dbType := u.Scheme
		dbUsername := u.User.Username()
		dbPassword, _ := u.User.Password()
		dbServer, dbPort, _ := net.SplitHostPort(u.Host)
		dbDataBase := u.Path[1:]
		var dbSchema string
		flag := strings.Compare(dbType, "mysql")
		if flag == 0 {
			// mysql database
			dbSchema = dbDataBase
		} else {
			// other database
			mapQuery, err := errors.Convert01(url.ParseQuery(u.RawQuery))
			if err != nil {
				return err
			}
			if value, ok := mapQuery["search_path"]; ok {
				if len(value) < 1 {
					return errors.Default.New("DB_URL search_path parses error")
				}
				dbSchema = value[0]
			} else {
				dbSchema = "public"
			}
		}
		config := viper.New()
		config.Set(projectName+".target", projectTarget)
		config.Set(projectName+".outputs."+projectTarget+".type", dbType)
		dbPortInt, _ := strconv.Atoi(dbPort)
		config.Set(projectName+".outputs."+projectTarget+".port", dbPortInt)
		config.Set(projectName+".outputs."+projectTarget+".password", dbPassword)
		config.Set(projectName+".outputs."+projectTarget+".schema", dbSchema)
		if flag == 0 {
			config.Set(projectName+".outputs."+projectTarget+".server", dbServer)
			config.Set(projectName+".outputs."+projectTarget+".username", dbUsername)
			config.Set(projectName+".outputs."+projectTarget+".database", dbDataBase)
		} else {
			config.Set(projectName+".outputs."+projectTarget+".host", dbServer)
			config.Set(projectName+".outputs."+projectTarget+".user", dbUsername)
			config.Set(projectName+".outputs."+projectTarget+".dbname", dbDataBase)
		}
		err = errors.Convert(config.WriteConfigAs("profiles.yml"))
		if err != nil {
			return err
		}
	}
	// if package.yml exist, install dbt dependencies
	_, err = errors.Convert01(os.Stat("packages.yml"))
	if err == nil {
		cmd := exec.Command("dbt", "deps")
		err = errors.Convert(cmd.Start())
		if err != nil {
			return err
		}
	}
	dbtExecParams := []string{"dbt", "run", "--project-dir", projectPath}
	if projectVars != nil {
		jsonProjectVars, err := json.Marshal(projectVars)
		if err != nil {
			return errors.Default.New("parameters vars json marshal error")
		}
		dbtExecParams = append(dbtExecParams, "--vars")
		dbtExecParams = append(dbtExecParams, string(jsonProjectVars))
	}
	if models != nil {
		dbtExecParams = append(dbtExecParams, "--select")
		dbtExecParams = append(dbtExecParams, models...)
	}
	if args != nil {
		dbtExecParams = append(dbtExecParams, args...)
	}
	cmd := exec.Command(dbtExecParams[0], dbtExecParams[1:]...)
	log.Info("dbt run script: ", cmd)
	stdout, _ := cmd.StdoutPipe()
	err = errors.Convert(cmd.Start())
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(stdout)
	var errStr string
	for scanner.Scan() {
		line := scanner.Text()
		log.Info(line)
		if strings.Contains(line, "ERROR") || errStr != "" {
			errStr += line + "\n"
		}
		if strings.Contains(line, "of") && strings.Contains(line, "OK") {
			taskCtx.IncProgress(1)
		}
	}
	if err := errors.Convert(scanner.Err()); err != nil {
		return err
	}
	err = errors.Convert(cmd.Wait())
	if err != nil {
		return errors.Internal.New(errStr)
	}
	if !cmd.ProcessState.Success() {
		log.Error(nil, "dbt run task error, please check!!!")
	}

	return nil
}

var DbtConverterMeta = core.SubTaskMeta{
	Name:             "DbtConverter",
	EntryPoint:       DbtConverter,
	EnabledByDefault: true,
	Description:      "Convert data by dbt",
}
