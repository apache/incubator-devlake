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
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/capability"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"net"
	"net/http"
	neturl "net/url"
	"os"
	"regexp"
	"strings"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/spf13/cast"
	ssh2 "golang.org/x/crypto/ssh"
)

// We have done comparison experiments for git2go and go-git, and the results show that git2go has better performance.
// We kept go-git because it supports cloning via key-based SSH.

const DefaultUser = "git"

func cloneOverSSH(ctx plugin.SubTaskContext, url, dir, passphrase string, pk []byte) errors.Error {
	key, err := ssh.NewPublicKeys(DefaultUser, pk, passphrase)
	if err != nil {
		return errors.Convert(err)
	}
	key.HostKeyCallbackHelper = ssh.HostKeyCallbackHelper{
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh2.PublicKey) error {
			return nil
		},
	}
	var data []byte
	buf := bytes.NewBuffer(data)
	done := make(chan struct{}, 1)
	go refreshCloneProgress(ctx, done, buf)
	_, err = gogit.PlainCloneContext(ctx.GetContext(), dir, true, &gogit.CloneOptions{
		URL:      url,
		Auth:     key,
		Progress: buf,
	})
	done <- struct{}{}
	if err != nil {
		return errors.Convert(err)
	}
	return nil
}

func (l *GitRepoCreator) cloneOverHTTP(ctx plugin.SubTaskContext, withGoGit bool, repoId, url, user, password, proxy string) (RepoCollector, errors.Error) {
	return withTempDirectory(func(dir string) (RepoCollector, error) {
		var data []byte
		buf := bytes.NewBuffer(data)
		done := make(chan struct{}, 1)
		go refreshCloneProgress(ctx, done, buf)
		cloneOptions := &gogit.CloneOptions{
			URL:      url,
			Progress: buf,
		}
		if proxy != "" {
			proxyUrl, err := neturl.Parse(proxy)
			if err != nil {
				l.logger.Error(err, "parse proxy")
				return nil, fmt.Errorf("parse %s err: %w", proxyUrl, err)
			}
			customClient := &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxyUrl),
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},

				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}
			client.InstallProtocol("https", githttp.NewClient(customClient))
		}
		if user != "" {
			cloneOptions.Auth = &githttp.BasicAuth{
				Username: user,
				Password: password,
			}
		}
		// fmt.Printf("CloneOverHTTP clone opt: %+v\ndir: %v, repo: %v, id: %v, user: %v, passwd: %v, proxy: %v\n", cloneOptions, dir, url, repoId, user, password, proxy)
		if isAzureRepo(ctx.GetContext(), url) {
			// https://github.com/go-git/go-git/issues/64
			// https://github.com/go-git/go-git/blob/master/_examples/azure_devops/main.go#L34
			transport.UnsupportedCapabilities = []capability.Capability{
				capability.ThinPack,
			}
		}
		_, err := gogit.PlainCloneContext(ctx.GetContext(), dir, true, cloneOptions)
		done <- struct{}{}
		if err != nil {
			l.logger.Error(err, "PlainCloneContext")
			// Some sensitive information such as password will be released in this err.
			if err.Error() == "repository not found" {
				// do nothing, it's a safe error message.
			} else {
				err = fmt.Errorf("plain clone git error")
			}
			return nil, err
		}
		if withGoGit {
			return l.LocalGoGitRepo(dir, repoId)
		}
		return l.LocalRepo(dir, repoId)
	})
}

func (l *GitRepoCreator) cloneOverSSH(ctx plugin.SubTaskContext, withGoGit bool, repoId, url, privateKey, passphrase string) (RepoCollector, errors.Error) {
	return withTempDirectory(func(dir string) (RepoCollector, error) {
		pk, err := base64.StdEncoding.DecodeString(privateKey)
		if err != nil {
			return nil, err
		}
		err = cloneOverSSH(ctx, url, dir, passphrase, pk)
		if err != nil {
			return nil, err
		}
		if withGoGit {
			return l.LocalGoGitRepo(dir, repoId)
		}
		return l.LocalRepo(dir, repoId)
	})
}

func setCloneProgress(subTaskCtx plugin.SubTaskContext, cloneProgressInfo string) {
	if cloneProgressInfo == "" {
		return
	}
	re, err := regexp.Compile(`\d+/\d+`) // find strings like 12/123.
	if err != nil {
		panic(err)
	}
	progress := re.FindAllString(cloneProgressInfo, -1)
	lenProgress := len(progress)
	if lenProgress == 0 {
		return
	}
	latestProgress := progress[lenProgress-1]
	latestProgressInfo := strings.Split(latestProgress, "/")
	if len(latestProgressInfo) == 2 {
		step := latestProgressInfo[0]
		total := latestProgressInfo[1]
		subTaskCtx.SetProgress(cast.ToInt(step), cast.ToInt(total))
	}
}

func refreshCloneProgress(subTaskCtx plugin.SubTaskContext, done chan struct{}, buf *bytes.Buffer) {
	ticker := time.NewTicker(time.Second * 1)
	func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if buf != nil {
					cloneProgressInfo := buf.String()
					setCloneProgress(subTaskCtx, cloneProgressInfo)
				}
			}
		}
	}()
}

func isAzureRepo(ctx context.Context, repoUrl string) bool {
	return strings.Contains(repoUrl, "dev.azure.com")
}

func withTempDirectory(f func(tempDir string) (RepoCollector, error)) (RepoCollector, errors.Error) {
	dir, err := os.MkdirTemp("", "gitextractor")
	if err != nil {
		return nil, errors.Convert(err)
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
	defer func() {
		if err != nil {
			cleanup()
		}
	}()
	repo, err := f(dir)
	if err != nil {
		return nil, errors.Convert(err)
	}
	if err := repo.SetCleanUp(cleanup); err != nil {
		return nil, errors.Convert(err)
	}
	return repo, nil
}
