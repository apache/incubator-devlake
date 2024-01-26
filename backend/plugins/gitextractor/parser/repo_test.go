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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/impls/logruslog"
	"github.com/apache/incubator-devlake/plugins/gitextractor/models"
	"github.com/apache/incubator-devlake/plugins/gitextractor/store"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	enableRepoTest       = false
	repoId               = "test-repo-id"
	runInLocal           = true
	ctx                  = context.Background()
	subTaskCtx           = &testSubTaskContext{}
	devlakeRepoRemoteURL = "https://github.com/apache/incubator-devlake"
	output               = "./output"

	repoMericoLake        = "/Users/houlinwei/Code/go/src/github.com/merico-dev/lake"
	repoMericoLakeWebsite = "/Users/houlinwei/Code/go/src/github.com/merico-dev/website"
	simpleRepo            = "/Users/houlinwei/Code/go/src/github.com/merico-dev/test/demo1"

	logger log.Logger

	storage        models.Store
	gitRepoCreator *GitRepoCreator

	goGitStorage     models.Store
	goGitRepoCreator *GitRepoCreator
)

func TestMain(m *testing.M) {
	if !enableRepoTest {
		return
	}
	fmt.Println("test main starts")
	logger = logruslog.Global.Nested("git extractor")
	fmt.Println("logger inited")

	clearOutput()

	var err error
	storage, err = store.NewCsvStore(output + "_libgit2")
	if err != nil {
		panic(err)
	}
	defer storage.Close()
	fmt.Println("git storage inited")
	gitRepoCreator = NewGitRepoCreator(storage, logger)

	goGitStorage, err = store.NewCsvStore(output + "_gogit")
	if err != nil {
		panic(err)
	}
	defer goGitStorage.Close()
	fmt.Println("go git storage inited")
	goGitRepoCreator = NewGitRepoCreator(goGitStorage, logger)

	fmt.Printf("test main run success\n\tlogger: %+v\tstorage: %+v\tgogit storage: %+v\n", logger, storage, goGitStorage)
	m.Run()
}

func getRepos(localRepoDir string) (RepoCollector, RepoCollector) {
	var gitRepo RepoCollector
	var goGitRepo RepoCollector
	var err errors.Error

	if runInLocal {
		repoPath := localRepoDir
		gitRepo, err = gitRepoCreator.LocalRepo(repoPath, repoId)
		if err != nil {
			panic(err)
		}
		goGitRepo, err = goGitRepoCreator.LocalGoGitRepo(repoPath, repoId)
		if err != nil {
			panic(err)
		}
	} else {
		gitRepo, err = gitRepoCreator.CloneOverHTTP(subTaskCtx, repoId, devlakeRepoRemoteURL, "", "", "")
		if err != nil {
			panic(err)
		}
		goGitRepo, err = goGitRepoCreator.CloneGoGitRepoOverHTTP(subTaskCtx, repoId, devlakeRepoRemoteURL, "", "", "")
		if err != nil {
			panic(err)
		}
	}
	return goGitRepo, gitRepo
}

func TestGitRepo_CountRepoInfo(t *testing.T) {
	if !enableRepoTest {
		return
	}
	goGitRepo, gitRepo := getRepos(repoMericoLakeWebsite)

	{
		tagsCount1, err1 := gitRepo.CountTags(ctx)
		if err1 != nil {
			panic(err1)
		}
		tagsCount2, err2 := goGitRepo.CountTags(ctx)
		if err2 != nil {
			panic(err2)
		}
		t.Logf("[tagsCount] libgit2 result: %d, gogit result: %d", tagsCount1, tagsCount2)
		assert.Equalf(t, tagsCount1, tagsCount2, "unexpected")
	}

	{
		branchesCount1, err1 := gitRepo.CountBranches(ctx)
		if err1 != nil {
			panic(err1)
		}
		branchesCount2, err2 := goGitRepo.CountBranches(ctx)
		if err2 != nil {
			panic(err2)
		}
		t.Logf("[branchesCount] libgit2 result: %d, gogit result: %d", branchesCount1, branchesCount2)
		assert.Equalf(t, branchesCount1, branchesCount2, "unexpected")
	}

	{
		commitCount1, err1 := gitRepo.CountCommits(ctx)
		if err1 != nil {
			panic(err1)
		}
		commitCount2, err2 := goGitRepo.CountCommits(ctx)
		if err2 != nil {
			panic(err2)
		}
		t.Logf("[commitCount] libgit2 result: %d, gogit result: %d", commitCount1, commitCount2)
		assert.Equalf(t, commitCount1, commitCount2, "unexpected")
	}

}

// all testes pass
func TestGitRepo_CollectRepoInfo(t *testing.T) {
	if !enableRepoTest {
		return
	}
	goGitRepo, gitRepo := getRepos(simpleRepo)

	{
		// finished
		subTaskCtxCollectTags := &testSubTaskContext{}
		if err1 := gitRepo.CollectTags(subTaskCtxCollectTags); err1 != nil {
			panic(err1)
		}
		subTaskCtxCollectTagsWithGoGit := &testSubTaskContext{}
		if err2 := goGitRepo.CollectTags(subTaskCtxCollectTagsWithGoGit); err2 != nil {
			panic(err2)
		}
		t.Logf("[CollectTags] libgit2 result: %+v, gogit result: %+v", subTaskCtxCollectTags, subTaskCtxCollectTagsWithGoGit)
		assert.Equalf(t, subTaskCtxCollectTags.total, subTaskCtxCollectTagsWithGoGit.total, "unexpected")
	}

	{
		// finished
		subTaskCtxCollectBranches := &testSubTaskContext{}
		if err1 := gitRepo.CollectBranches(subTaskCtxCollectBranches); err1 != nil {
			panic(err1)
		}
		subTaskCtxCollectBranchesWithGoGit := &testSubTaskContext{}
		if err2 := goGitRepo.CollectBranches(subTaskCtxCollectBranchesWithGoGit); err2 != nil {
			panic(err2)
		}
		t.Logf("[CollectBranches] libgit2 result: %+v, gogit result: %+v", subTaskCtxCollectBranches, subTaskCtxCollectBranchesWithGoGit)
		assert.Equalf(t, subTaskCtxCollectBranches.total, subTaskCtxCollectBranchesWithGoGit.total, "unexpected")
	}

	{
		subTaskCtxCollectCommits := &testSubTaskContext{}
		if err1 := gitRepo.CollectCommits(subTaskCtxCollectCommits); err1 != nil {
			panic(err1)
		}
		subTaskCtxCCollectCommitsWithGoGit := &testSubTaskContext{}
		if err2 := goGitRepo.CollectCommits(subTaskCtxCCollectCommitsWithGoGit); err2 != nil {
			panic(err2)
		}

		t.Logf("[CollectCommits] libgit2 result: %+v, gogit result: %+v", subTaskCtxCollectCommits, subTaskCtxCCollectCommitsWithGoGit)
		fmt.Println(subTaskCtxCollectCommits.total, subTaskCtxCCollectCommitsWithGoGit.total)
		assert.Equalf(t, subTaskCtxCollectCommits.total, subTaskCtxCCollectCommitsWithGoGit.total, "unexpected")
	}

	{
		subTaskCtxCollectDiffLine := &testSubTaskContext{}
		if err1 := gitRepo.CollectDiffLine(subTaskCtxCollectDiffLine); err1 != nil {
			panic(err1)
		}
		subTaskCtxCollectDiffLineWithGoGit := &testSubTaskContext{}
		if err2 := goGitRepo.CollectDiffLine(subTaskCtxCollectDiffLineWithGoGit); err2 != nil {
			panic(err2)
		}

		t.Logf("[CollectDiffLine] libgit2 result: %+v, gogit result: %+v", subTaskCtxCollectDiffLine, subTaskCtxCollectDiffLineWithGoGit)
		fmt.Println(subTaskCtxCollectDiffLine.total, subTaskCtxCollectDiffLineWithGoGit.total)
		assert.Equalf(t, subTaskCtxCollectDiffLine.total, subTaskCtxCollectDiffLineWithGoGit.total, "unexpected")
	}
}

func clearOutput() {
	os.RemoveAll(fmt.Sprintf("./output_libgit2"))
	os.RemoveAll(fmt.Sprintf("./output_gogit"))
}

func TestGitRepo_CollectCommits(t *testing.T) {
	if !enableRepoTest {
		return
	}
	repoPath := simpleRepo
	gitRepo, err := gitRepoCreator.LocalRepo(repoPath, repoId)
	if err != nil {
		panic(err)
	}
	goGitRepo, err := goGitRepoCreator.LocalGoGitRepo(repoPath, repoId)
	if err != nil {
		panic(err)
	}

	{
		subTaskCtxCollectCommits := &testSubTaskContext{}
		if err1 := gitRepo.CollectCommits(subTaskCtxCollectCommits); err1 != nil {
			panic(err1)
		}

		subTaskCtxCCollectCommitsWithGoGit := &testSubTaskContext{}
		if err2 := goGitRepo.CollectCommits(subTaskCtxCCollectCommitsWithGoGit); err2 != nil {
			panic(err2)
		}

		t.Logf("[CollectCommits] libgit2 result: %+v, gogit result: %+v", subTaskCtxCollectCommits, subTaskCtxCCollectCommitsWithGoGit)
		fmt.Println(subTaskCtxCollectCommits.total, subTaskCtxCCollectCommitsWithGoGit.total)
		assert.Equalf(t, subTaskCtxCollectCommits.total, subTaskCtxCCollectCommitsWithGoGit.total, "unexpected")
	}
}

func TestGitRepo_CollectDiffLine(t *testing.T) {
	if !enableRepoTest {
		return
	}
	repoPath := simpleRepo
	gitRepo, err := gitRepoCreator.LocalRepo(repoPath, repoId)
	if err != nil {
		panic(err)
	}
	goGitRepo, err := goGitRepoCreator.LocalGoGitRepo(repoPath, repoId)
	if err != nil {
		panic(err)
	}

	{
		subTaskCtxCollectDiffLine := &testSubTaskContext{}
		if err1 := gitRepo.CollectDiffLine(subTaskCtxCollectDiffLine); err1 != nil {
			panic(err1)
		}
		//t.Logf("[CollectDiffLine] libgit2 result: %+v", subTaskCtxCollectDiffLine)

		subTaskCtxCollectDiffLineWithGoGit := &testSubTaskContext{}
		if err2 := goGitRepo.CollectDiffLine(subTaskCtxCollectDiffLineWithGoGit); err2 != nil {
			panic(err2)
		}

		t.Logf("[CollectCommits] libgit2 result: %+v, gogit result: %+v", subTaskCtxCollectDiffLine, subTaskCtxCollectDiffLineWithGoGit)
		fmt.Println(subTaskCtxCollectDiffLine.total, subTaskCtxCollectDiffLineWithGoGit.total)
		assert.Equalf(t, subTaskCtxCollectDiffLine.total, subTaskCtxCollectDiffLineWithGoGit.total, "unexpected")
	}
}
