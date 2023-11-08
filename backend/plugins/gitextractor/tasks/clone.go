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
	"context"
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/gitextractor/models"
	"github.com/apache/incubator-devlake/plugins/gitextractor/parser"
	"github.com/apache/incubator-devlake/plugins/gitextractor/store"
	"strings"
)

var CloneGitRepoMeta = plugin.SubTaskMeta{
	Name:             "cloneGitRepo",
	EntryPoint:       CloneGitRepo,
	EnabledByDefault: true,
	Required:         true,
	Description:      "clone a git repo, make it available to later tasks",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
}

func CloneGitRepo(subTaskCtx plugin.SubTaskContext) errors.Error {
	taskData, ok := subTaskCtx.GetData().(*GitExtractorTaskData)
	if !ok {
		panic("git repo reference not found on context")
	}
	op := taskData.Options
	storage := store.NewDatabase(subTaskCtx, op.RepoId)
	repo, err := NewGitRepo(subTaskCtx.GetContext(), subTaskCtx.GetLogger(), storage, op)
	if err != nil {
		return err
	}
	taskData.GitRepo = repo
	subTaskCtx.TaskContext().SetData(taskData)
	return nil
}

// NewGitRepo create and return a new parser git repo
func NewGitRepo(ctx context.Context, logger log.Logger, storage models.Store, op *GitExtractorOptions) (*parser.GitRepo, errors.Error) {
	var err errors.Error
	var repo *parser.GitRepo
	p := parser.NewGitRepoCreator(storage, logger)
	if strings.HasPrefix(op.Url, "http") {
		repo, err = p.CloneOverHTTP(ctx, op.RepoId, op.Url, op.User, op.Password, op.Proxy)
	} else if url := strings.TrimPrefix(op.Url, "ssh://"); strings.HasPrefix(url, "git@") {
		repo, err = p.CloneOverSSH(ctx, op.RepoId, url, op.PrivateKey, op.Passphrase)
	} else if strings.HasPrefix(op.Url, "/") {
		repo, err = p.LocalRepo(op.Url, op.RepoId)
	} else {
		return nil, errors.BadInput.New(fmt.Sprintf("unsupported url [%s]", op.Url))
	}
	return repo, err
}
