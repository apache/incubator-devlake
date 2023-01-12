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

package impl

import (
	"fmt"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitextractor/models"
	"github.com/apache/incubator-devlake/plugins/gitextractor/parser"
	"github.com/apache/incubator-devlake/plugins/gitextractor/store"
	"github.com/apache/incubator-devlake/plugins/gitextractor/tasks"
	"strings"
)

var _ plugin.PluginMeta = (*GitExtractor)(nil)
var _ plugin.PluginTask = (*GitExtractor)(nil)
var _ plugin.PluginModel = (*GitExtractor)(nil)

type GitExtractor struct{}

func (p GitExtractor) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{}
}

func (p GitExtractor) Description() string {
	return "extract infos from git repository"
}

// return all available subtasks, framework will run them for you in order
func (p GitExtractor) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectGitCommitMeta,
		tasks.CollectGitBranchMeta,
		tasks.CollectGitTagMeta,
		tasks.CollectGitDiffLineMeta,
	}
}

// based on task context and user input options, return data that shared among all subtasks
func (p GitExtractor) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.GitExtractorOptions
	if err := helper.Decode(options, &op, nil); err != nil {
		return nil, err
	}
	if err := op.Valid(); err != nil {
		return nil, err
	}
	storage := store.NewDatabase(taskCtx, op.RepoId)
	repo, err := NewGitRepo(taskCtx.GetLogger(), storage, op)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (p GitExtractor) Close(taskCtx plugin.TaskContext) errors.Error {
	if repo, ok := taskCtx.GetData().(*parser.GitRepo); ok {
		if err := repo.Close(); err != nil {
			return errors.Convert(err)
		}
	}
	return nil
}

func (p GitExtractor) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/gitextractor"
}

// NewGitRepo create and return a new parser git repo
func NewGitRepo(logger log.Logger, storage models.Store, op tasks.GitExtractorOptions) (*parser.GitRepo, errors.Error) {
	var err errors.Error
	var repo *parser.GitRepo
	p := parser.NewGitRepoCreator(storage, logger)
	if strings.HasPrefix(op.Url, "http") {
		repo, err = p.CloneOverHTTP(op.RepoId, op.Url, op.User, op.Password, op.Proxy)
	} else if url := strings.TrimPrefix(op.Url, "ssh://"); strings.HasPrefix(url, "git@") {
		repo, err = p.CloneOverSSH(op.RepoId, url, op.PrivateKey, op.Passphrase)
	} else if strings.HasPrefix(op.Url, "/") {
		repo, err = p.LocalRepo(op.Url, op.RepoId)
	} else {
		return nil, errors.BadInput.New(fmt.Sprintf("unsupported url [%s]", op.Url))
	}
	return repo, err
}
