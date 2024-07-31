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
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

var _ RepoCloner = (*GitcliCloner)(nil)
var ErrNoData = errors.NotModified.New("No data to be collected")

type GitcliCloner struct {
	logger       log.Logger
	stateManager *api.SubtaskStateManager
}

func NewGitcliCloner(basicRes context.BasicRes) *GitcliCloner {
	return &GitcliCloner{
		logger: basicRes.GetLogger().Nested("gitcli"),
	}
}

// CloneRepoConfig is the configuration for the CloneRepo method
// the subtask should run in Full Sync mode whenever the configuration is changed
type CloneRepoConfig struct {
	UseGoGit        *bool
	SkipCommitStat  *bool
	SkipCommitFiles *bool
	NoShallowClone  bool
}

func (g *GitcliCloner) IsIncremental() bool {
	if g != nil && g.stateManager != nil {
		if g.stateManager.GetSince() != nil {
			return true
		}
		return g.stateManager.IsIncremental()
	}

	return false
}

func (g *GitcliCloner) CloneRepo(ctx plugin.SubTaskContext, localDir string) errors.Error {
	taskData := ctx.GetData().(*GitExtractorTaskData)
	var since *time.Time
	if !taskData.Options.NoShallowClone {
		stateManager, err := api.NewSubtaskStateManager(&api.SubtaskCommonArgs{
			SubTaskContext: ctx,
			Params:         taskData.Options.GitExtractorApiParams,
			SubtaskConfig: CloneRepoConfig{
				UseGoGit:        taskData.Options.UseGoGit,
				SkipCommitStat:  taskData.Options.SkipCommitStat,
				SkipCommitFiles: taskData.Options.SkipCommitFiles,
				NoShallowClone:  taskData.Options.NoShallowClone,
			},
		})
		if err != nil {
			return err
		}
		g.stateManager = stateManager
		since = stateManager.GetSince()

	}

	err := g.execGitCloneCommand(ctx, localDir, since)
	if err != nil {
		return err
	}
	// deepen the commits by 1 more step to avoid https://github.com/apache/incubator-devlake/issues/7426
	if since != nil {
		// fixes error described on https://stackoverflow.com/questions/63878612/git-fatal-error-in-object-unshallow-sha-1
		// It might be casued by the commit which being deepen has mulitple parent(e.g. a merge commit), not sure.
		if err := g.execGitCommandIn(ctx, localDir, "repack", "-d"); err != nil {
			return errors.Default.Wrap(err, "failed to repack the repo")
		}
		// deepen would fail on a EMPTY repo, ignore the error
		if err := g.execGitCommandIn(ctx, localDir, "fetch", "--deepen=1"); err != nil {
			g.logger.Error(err, "failed to deepen the cloned repo")
		}
	}

	// save state
	if g.stateManager != nil {
		return g.stateManager.Close()
	}
	return nil
}

func (g *GitcliCloner) execGitCloneCommand(ctx plugin.SubTaskContext, localDir string, since *time.Time) errors.Error {
	taskData := ctx.GetData().(*GitExtractorTaskData)
	var args []string
	if *taskData.Options.SkipCommitStat {
		args = append(args, "--filter=blob:none")
	}
	if since != nil {
		// to fetch newly added commits from ALL branches, we need to the following guide:
		//    https://stackoverflow.com/questions/23708231/git-shallow-clone-clone-depth-misses-remote-branches

		// 1. clone the repo with depth 1
		cloneArgs := append([]string{"clone", taskData.Options.Url, localDir, "--depth=1", "--bare"}, args...)
		if err := g.execGitCommand(ctx, cloneArgs...); err != nil {
			return err
		}
		// 2. configure to fetch all branches from the remote server so we can collect new commits from them
		gitConfig, err := os.OpenFile(path.Join(localDir, "config"), os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return errors.Default.Wrap(err, "failed to open git config file")
		}
		_, err = gitConfig.WriteString("\tfetch = +refs/heads/*:refs/remotes/origin/*\n")
		if err != nil {
			return errors.Default.Wrap(err, "failed to write to git config file")
		}
		// 3. fetch all branches with depth=1 so the next step would collect less commits
		// (I don't know why, but it reduced total number of commits from 18k to 7k on https://gitlab.com/gitlab-org/gitlab-foss.git with the same parameters)
		fetchBranchesArgs := append([]string{"fetch", "--depth=1", "origin"}, args...)
		if err := g.execGitCommandIn(ctx, localDir, fetchBranchesArgs...); err != nil {
			return errors.Default.Wrap(err, "failed to fetch all branches from the remote server")
		}
		// 4. fetch all new commits from all branches since the given time
		args = append([]string{"fetch", fmt.Sprintf("--shallow-since=%s", since.Format(time.RFC3339))}, args...)
		if err := g.execGitCommandIn(ctx, localDir, args...); err != nil {
			g.logger.Warn(err, "shallow fetch failed")
		}
		return nil
	} else {
		args = append([]string{"clone", taskData.Options.Url, localDir, "--bare"}, args...)
		return g.execGitCommand(ctx, args...)
	}
}

func (g *GitcliCloner) execGitCommand(ctx plugin.SubTaskContext, args ...string) errors.Error {
	return g.execGitCommandIn(ctx, "", args...)
}

func (g *GitcliCloner) execGitCommandIn(ctx plugin.SubTaskContext, workingDir string, args ...string) errors.Error {
	taskData := ctx.GetData().(*GitExtractorTaskData)
	env := []string{}
	if args[0] == "clone" || args[0] == "fetch" {
		// support proxy
		if taskData.ParsedURL.Scheme == "http" || taskData.ParsedURL.Scheme == "https" {
			if taskData.Options.Proxy != "" {
				env = append(env, fmt.Sprintf("HTTPS_PROXY=%s", taskData.Options.Proxy))
			}
			if taskData.ParsedURL.Scheme == "https" && ctx.GetConfigReader().GetBool("IN_SECURE_SKIP_VERIFY") {
				args = append(args, "-c http.sslVerify=false")
			}
		} else if taskData.ParsedURL.Scheme == "ssh" {
			var sshCmdArgs []string
			if taskData.Options.Proxy != "" {
				parsedProxyURL, e := url.Parse(taskData.Options.Proxy)
				if e != nil {
					return errors.BadInput.Wrap(e, "failed to parse the proxy URL")
				}
				proxyCommand := "corkscrew"
				sshCmdArgs = append(sshCmdArgs, "-o", fmt.Sprintf(`ProxyCommand="%s %s %s %%h %%p"`, proxyCommand, parsedProxyURL.Hostname(), parsedProxyURL.Port()))
			}
			// support private key
			if taskData.Options.PrivateKey != "" {
				pkFile, err := os.CreateTemp("", "gitext-pk")
				if err != nil {
					g.logger.Error(err, "create temp private key file error")
					return errors.Default.New("failed to handle the private key")
				}
				if _, e := pkFile.WriteString(taskData.Options.PrivateKey + "\n"); e != nil {
					g.logger.Error(err, "write private key file error")
					return errors.Default.New("failed to write the  private key")
				}
				pkFile.Close()
				if e := os.Chmod(pkFile.Name(), 0600); e != nil {
					g.logger.Error(err, "chmod private key file error")
					return errors.Default.New("failed to modify the private key")
				}

				if taskData.Options.Passphrase != "" {
					pp := exec.CommandContext(
						ctx.GetContext(),
						"ssh-keygen", "-p",
						"-P", taskData.Options.Passphrase,
						"-N", "",
						"-f", pkFile.Name(),
					)
					if ppout, pperr := pp.CombinedOutput(); pperr != nil {
						g.logger.Error(pperr, "change private key passphrase error")
						g.logger.Info(string(ppout))
						return errors.Default.New("failed to decrypt the private key")
					}
				}
				defer os.Remove(pkFile.Name())
				sshCmdArgs = append(sshCmdArgs, fmt.Sprintf("-i %s -o StrictHostKeyChecking=no", pkFile.Name()))
			}
			if len(sshCmdArgs) > 0 {
				env = append(env, fmt.Sprintf("GIT_SSH_COMMAND=ssh %s", strings.Join(sshCmdArgs, " ")))
			}
		}
	}
	g.logger.Debug("git %v", args)
	cmd := exec.CommandContext(ctx.GetContext(), "git", args...)
	cmd.Env = env
	cmd.Dir = workingDir
	return g.execCommand(cmd)
}

func (g *GitcliCloner) execCommand(cmd *exec.Cmd) errors.Error {
	output, err := cmd.CombinedOutput()
	if err != nil {
		outputString := string(output)
		if strings.Contains(outputString, "fatal: error processing shallow info: 4") ||
			strings.Contains(outputString, "fatal: the remote end hung up unexpectedly") {
			return ErrNoData
		}
		return errors.Default.New(fmt.Sprintf("git cmd %v in %s failed: %s", sanitizeArgs(cmd.Args), cmd.Dir, outputString))
	}
	return nil
}

func sanitizeArgs(args []string) []string {
	var ret []string
	for _, arg := range args {
		u, err := url.Parse(arg)
		if err == nil && u != nil && u.User != nil {
			password, ok := u.User.Password()
			if ok {
				arg = strings.Replace(arg, password, strings.Repeat("*", len(password)), -1)
			}
		}
		ret = append(ret, arg)
	}
	return ret
}
