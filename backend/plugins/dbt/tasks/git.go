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
	"os"
	"os/exec"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

func Git(taskCtx plugin.SubTaskContext) errors.Error {
	logger := taskCtx.GetLogger()
	data := taskCtx.GetData().(*DbtTaskData)
	if data.Options.ProjectGitURL == "" {
		return nil
	}

	// clean ProjectPath
	err := os.RemoveAll(data.Options.ProjectPath)
	if err != nil {
		logger.Error(err, "cleanup before clone dbt project failed")
		return errors.Convert(err)
	}

	// git clone from ProjectGitURL into ProjectPath
	cmd := exec.Command("git", "clone", data.Options.ProjectGitURL, data.Options.ProjectPath)
	logger.Info("start clone dbt project: %v", cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error(err, "clone dbt project failed")
		return errors.Convert(err)
	}
	logger.Info("clone dbt project success: %v", string(out))
	return nil
}

var GitMeta = plugin.SubTaskMeta{
	Name:             "Git",
	EntryPoint:       Git,
	EnabledByDefault: true,
	Description:      "Clone dbt project from git",
}
