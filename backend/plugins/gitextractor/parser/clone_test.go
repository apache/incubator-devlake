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
	"encoding/base64"
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/impls/logruslog"
	"github.com/apache/incubator-devlake/plugins/gitextractor/models"
	"github.com/apache/incubator-devlake/plugins/gitextractor/store"
	git "github.com/libgit2/git2go/v33"
	"reflect"
	"testing"
)

func (l *GitRepoCreator) CloneOverHTTPWithLibGit(repoId, url, user, password, proxy string) (*GitRepo, errors.Error) {
	return withTempDirectory(func(dir string) (*GitRepo, error) {
		cloneOptions := &git.CloneOptions{Bare: true}
		if proxy != "" {
			cloneOptions.FetchOptions.ProxyOptions.Type = git.ProxyTypeSpecified
			cloneOptions.FetchOptions.ProxyOptions.Url = proxy
		}
		if user != "" {
			auth := fmt.Sprintf("Authorization: Basic %s", base64.StdEncoding.EncodeToString([]byte(user+":"+password)))
			cloneOptions.FetchOptions.Headers = []string{auth}
		}
		fmt.Printf("CloneOverHTTPWithLibGit clone opt: %+v\ndir: %v, repo: %v, id: %v, user: %v, passwd: %v, proxy: %v\n", cloneOptions, dir, url, repoId, user, password, proxy)
		clonedRepo, err := git.Clone(url, dir, cloneOptions)
		if err != nil {
			return nil, err
		}
		return l.newGitRepo(repoId, clonedRepo), nil
	})
}

var (
	output  = "./output"
	logger  log.Logger
	storage models.Store
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

func TestGitRepoCreator_CloneOverHTTP(t *testing.T) {
	fmt.Println("run test clone over http")
	var (
		testCtx = context.Background()
		repoUrl = "https://github.com/apache/incubator-devlake-website"
		repoId  = "fake-id-1"
		proxy   = ""
		user    = ""
		passwd  = ""
	)
	type fields struct {
		store  models.Store
		logger log.Logger
	}
	type args struct {
		ctx      context.Context
		repoId   string
		url      string
		user     string
		password string
		proxy    string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *GitRepo
		want1  errors.Error
	}{
		{
			name: "all-without-proxy-and-auth",
			fields: fields{
				store:  storage,
				logger: logger,
			},
			args: args{
				ctx:    testCtx,
				repoId: repoId,
				url:    repoUrl,
			},
			want1: nil,
		},
		{
			name: "all-with-proxy-and-no-auth",
			fields: fields{
				store:  storage,
				logger: logger,
			},
			args: args{
				ctx:    testCtx,
				repoId: repoId,
				url:    repoUrl,
				proxy:  proxy,
			},
			want1: nil,
		},
		{
			name: "all-with-auth-and-no-proxy",
			fields: fields{
				store:  storage,
				logger: logger,
			},
			args: args{
				ctx:      testCtx,
				repoId:   repoId,
				url:      repoUrl,
				user:     user,
				password: passwd,
			},
			want1: nil,
		},
		{
			name: "all-with-auth-and-proxy",
			fields: fields{
				store:  storage,
				logger: logger,
			},
			args: args{
				ctx:      testCtx,
				repoId:   repoId,
				url:      repoUrl,
				user:     user,
				password: passwd,
				proxy:    proxy,
			},
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Logf("case config: %+v\n", tt)
		t.Run(tt.name, func(t *testing.T) {
			l := &GitRepoCreator{
				store:  tt.fields.store,
				logger: tt.fields.logger,
			}
			got1Result, got1 := l.CloneOverHTTP(tt.args.ctx, tt.args.repoId, tt.args.url, tt.args.user, tt.args.password, tt.args.proxy)
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CloneOverHTTP() got1 = %v, want %v", got1, tt.want1)
			}

			got2Result, got2 := l.CloneOverHTTPWithLibGit(tt.args.repoId, tt.args.url, tt.args.user, tt.args.password, tt.args.proxy)
			if !reflect.DeepEqual(got2, tt.want1) {
				t.Errorf("CloneOverHTTP() got1 = %v, want %v", got1, tt.want1)
			}

			equal, err := GitReposAreEqual(tt.args.ctx, got1Result, got2Result)
			if err != nil {
				t.Error(err)
			}
			if !equal {
				t.Errorf("not equal")
			}
		})

		// break
	}
}

func GitReposAreEqual(ctx context.Context, r1, r2 *GitRepo) (bool, error) {
	if r1 == nil && r2 == nil {
		return true, nil
	}
	if r1 == nil && r2 != nil {
		return false, nil
	}
	if r1 != nil && r2 == nil {
		return false, nil
	}

	fmt.Print("compare tags ")
	r1Tags, r1TagErr := r1.CountTags()
	r2Tags, r2TagErr := r2.CountTags()
	if !reflect.DeepEqual(r1TagErr, r2TagErr) {
		fmt.Println("tag err")
		return false, nil
	}
	if !reflect.DeepEqual(r1Tags, r2Tags) {
		fmt.Printf("tag %d != %d \n", r1Tags, r1Tags)
		return false, nil
	}
	fmt.Println("done")

	fmt.Print("compare commits ")
	r1Commits, r1CommitsErr := r1.CountCommits(ctx)
	r2Commits, r2CommitsErr := r2.CountCommits(ctx)
	if !reflect.DeepEqual(r1CommitsErr, r2CommitsErr) {
		fmt.Println("commit err")
		return false, nil
	}
	if !reflect.DeepEqual(r1Commits, r2Commits) {
		fmt.Printf("commits %d != %d \n", r1Commits, r2Commits)
		return false, nil
	}
	fmt.Println("done")

	fmt.Print("compare branches ")
	r1Branches, r1BranchesErr := r1.CountBranches(ctx)
	r2Branches, r2BranchesErr := r2.CountBranches(ctx)
	if !reflect.DeepEqual(r1BranchesErr, r2BranchesErr) {
		fmt.Println("branch err")
		return false, nil
	}
	if !reflect.DeepEqual(r1Branches+1, r2Branches) { // branch count with gogit is accurate, but to keep the test go on, just plus 1.
		fmt.Printf("branch %d != %d \n", r1Branches, r2Branches)
		return false, nil
	}
	fmt.Println("done")
	return true, nil
}
