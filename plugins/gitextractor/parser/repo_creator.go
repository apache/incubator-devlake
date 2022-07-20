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

package parser

import (
	git "github.com/libgit2/git2go/v33"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitextractor/models"
)

const (
	BRANCH = "BRANCH"
	TAG    = "TAG"
)

type GitRepoCreator struct {
	store  models.Store
	logger core.Logger
}

func NewGitRepoCreator(store models.Store, logger core.Logger) *GitRepoCreator {
	return &GitRepoCreator{
		store:  store,
		logger: logger,
	}
}

func (l *GitRepoCreator) LocalRepo(repoPath, repoId string) (*GitRepo, error) {
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return nil, err
	}
	return l.newGitRepo(repoId, repo), nil
}

func (l *GitRepoCreator) newGitRepo(repoId string, repo *git.Repository) *GitRepo {
	return &GitRepo{
		store:  l.store,
		logger: l.logger,
		id:     repoId,
		repo:   repo,
	}
}
