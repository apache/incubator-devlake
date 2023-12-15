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
	"context"
	"fmt"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/impls/logruslog"
	"github.com/apache/incubator-devlake/plugins/gitextractor/models"
	"github.com/apache/incubator-devlake/plugins/gitextractor/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	output  = "./output"
	logger  log.Logger
	storage models.Store
	ctx     = context.Background()
)

func TestMain(m *testing.M) {
	fmt.Println("test main starts")
	logger = logruslog.Global.Nested("git extractor")
	fmt.Println("logger inited")
	var err error
	storage, err = store.NewCsvStore(output)
	if err != nil {
		panic(err)
	}
	fmt.Println("storage inited")
	defer storage.Close()
	fmt.Printf("test main run success, logger: %+v, storage: %+v\n", logger, storage)
	m.Run()
}

func TestGitRepo_CountRepoInfo(t *testing.T) {
	repoPath := "/Users/houlinwei/Code/go/src/github.com/merico-dev/lake"
	repoId := "test-repo-id"
	gitRepo, err := NewGitRepoCreator(storage, logger).LocalRepo(repoPath, repoId)
	if err != nil {
		panic(err)
	}
	tagsCount1, err1 := gitRepo.CountTags()
	if err1 != nil {
		panic(err1)
	}
	tagsCount2, err2 := gitRepo.CountTagsWithGoGit()
	if err2 != nil {
		panic(err2)
	}
	t.Logf("[tagsCount] libgit2 result: %d, gogit result: %d", tagsCount1, tagsCount2)
	assert.Equal(t, tagsCount1, tagsCount2)

	branchesCount1, err1 := gitRepo.CountBranches(ctx)
	if err1 != nil {
		panic(err1)
	}
	branchesCount2, err2 := gitRepo.CountBranchesWithGoGit(ctx)
	if err2 != nil {
		panic(err2)
	}
	t.Logf("[branchesCount] libgit2 result: %d, gogit result: %d", branchesCount1, branchesCount2)
	assert.Equal(t, branchesCount1, branchesCount2)

	commitCount1, err1 := gitRepo.CountCommits(ctx)
	if err1 != nil {
		panic(err1)
	}
	commitCount2, err2 := gitRepo.CountCommitsWithGoGit(ctx)
	if err2 != nil {
		panic(err2)
	}
	t.Logf("[commitCount] libgit2 result: %d, gogit result: %d", commitCount1, commitCount2)
	assert.Equal(t, commitCount1, commitCount2)

}

func compare(b1, b2 []string) {
	fmt.Println(len(b1), len(b2))
	for _, b := range b2 {
		var found bool
		for _, bb := range b1 {
			if bb == b {
				found = true
			}
		}
		if !found {
			fmt.Printf("%s from b2, not found in b1\n", b)
		}
	}

	for _, b := range b1 {
		var found bool
		for _, bb := range b2 {
			if bb == b {
				found = true
			}
		}
		if !found {
			fmt.Printf("%s from b1, not found in b2\n", b)
		}
	}
}
