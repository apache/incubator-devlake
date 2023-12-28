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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/plugins/gitextractor/models"
	gogit "github.com/go-git/go-git/v5"
	git "github.com/libgit2/git2go/v33"
)

const (
	BRANCH      = "BRANCH"
	TAG         = "TAG"
	EnableGoGit = true
)

type GitRepoCreator struct {
	store      models.Store
	goGitStore models.Store
	logger     log.Logger
}

func NewGitRepoCreator(store, goGitStore models.Store, logger log.Logger) *GitRepoCreator {
	return &GitRepoCreator{
		store:      store,
		logger:     logger,
		goGitStore: goGitStore,
	}
}

// LocalRepo open a local repository
func (l *GitRepoCreator) LocalRepo(repoPath, repoId string) (*GitRepo, errors.Error) {
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		l.logger.Error(err, "OpenRepository")
		return nil, errors.Convert(err)
	}
	var goGitRepo *gogit.Repository
	if EnableGoGit {
		var err error
		goGitRepo, err = gogit.PlainOpen(repoPath)
		if err != nil {
			return nil, errors.Convert(err)
		}
	}
	return l.newGitRepo(repoId, repo, goGitRepo), nil
}

func (l *GitRepoCreator) newGitRepo(repoId string, repo *git.Repository, goGitRespo *gogit.Repository) *GitRepo {
	return &GitRepo{
		store:  l.store,
		logger: l.logger,
		id:     repoId,
		repo:   repo,

		goGitRepo:  goGitRespo,
		goGitStore: l.goGitStore,
	}
}
