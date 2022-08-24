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
	"encoding/json"
	"net/url"
	"strings"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitextractor/models"
	"github.com/apache/incubator-devlake/plugins/gitextractor/parser"
	"github.com/apache/incubator-devlake/plugins/gitextractor/store"
	"github.com/apache/incubator-devlake/plugins/gitextractor/tasks"
	"github.com/mitchellh/mapstructure"
)

const (
	credentialMaskString = "*****"
)

var _ core.PluginMeta = (*GitExtractor)(nil)
var _ core.PluginTask = (*GitExtractor)(nil)
var _ core.CredentialMasker = (*GitExtractor)(nil)

type GitExtractor struct{}

func (plugin GitExtractor) GetTablesInfo() []core.Tabler {
	return []core.Tabler{}
}

func (plugin GitExtractor) Description() string {
	return "extract infos from git repository"
}

// return all available subtasks, framework will run them for you in order
func (plugin GitExtractor) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CollectGitCommitMeta,
		tasks.CollectGitBranchMeta,
		tasks.CollectGitTagMeta,
	}
}

// based on task context and user input options, return data that shared among all subtasks
func (plugin GitExtractor) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	var op tasks.GitExtractorOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}
	err = op.Valid()
	if err != nil {
		return nil, err
	}
	storage := store.NewDatabase(taskCtx, op.Url)
	repo, err := newGitRepo(taskCtx.GetLogger(), storage, op)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (plugin GitExtractor) Close(taskCtx core.TaskContext) error {
	if repo, ok := taskCtx.GetData().(*parser.GitRepo); ok {
		if err := repo.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (plugin GitExtractor) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/gitextractor"
}

func (plugin GitExtractor) Mask(input []byte) []byte {
	var op tasks.GitExtractorOptions
	err := json.Unmarshal(input, &op)
	if err != nil {
		return input
	}
	if op.Password != "" {
		op.Password = credentialMaskString
	}
	if op.Passphrase != "" {
		op.Passphrase = credentialMaskString
	}
	if op.PrivateKey != "" {
		op.PrivateKey = credentialMaskString
	}
	if u, err := url.Parse(op.Url); err == nil {
		if u.User != nil {
			_, hasPassword := u.User.Password()
			if hasPassword {
				u.User = url.UserPassword(u.User.Username(), credentialMaskString)
				op.Url = u.String()
			}
		}
	}
	blob, err := json.Marshal(op)
	if err != nil {
		return input
	}
	return blob
}

func newGitRepo(logger core.Logger, storage models.Store, op tasks.GitExtractorOptions) (*parser.GitRepo, error) {
	var err error
	var repo *parser.GitRepo
	p := parser.NewGitRepoCreator(storage, logger)
	if strings.HasPrefix(op.Url, "http") {
		repo, err = p.CloneOverHTTP(op.RepoId, op.Url, op.User, op.Password, op.Proxy)
	} else if url := strings.TrimPrefix(op.Url, "ssh://"); strings.HasPrefix(url, "git@") {
		repo, err = p.CloneOverSSH(op.RepoId, url, op.PrivateKey, op.Passphrase)
	} else if strings.HasPrefix(op.Url, "/") {
		repo, err = p.LocalRepo(op.Url, op.RepoId)
	}
	return repo, err
}

// PluginEntry is a variable exported for Framework to search and load
var PluginEntry GitExtractor //nolint
