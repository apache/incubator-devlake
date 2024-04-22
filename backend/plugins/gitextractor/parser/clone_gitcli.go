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
	"bufio"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

var _ RepoCloner = (*GitcliCloner)(nil)
var ErrShallowInfoProcessing = errors.BadInput.New("No data found for the selected time range. Please revise the 'Time Range' on your Project/Blueprint/Configuration page or in the API parameter.")
var ErrNoDataOnIncrementalMode = errors.NotModified.New("No data found since the previous run.")

type GitcliCloner struct {
	logger       log.Logger
	stateManager *api.CollectorStateManager
}

func NewGitcliCloner(basicRes context.BasicRes) *GitcliCloner {
	return &GitcliCloner{
		logger: basicRes.GetLogger().Nested("gitcli"),
	}
}

func (g *GitcliCloner) CloneRepo(ctx plugin.SubTaskContext, localDir string) errors.Error {
	taskData := ctx.GetData().(*GitExtractorTaskData)
	var since *time.Time
	if !taskData.Options.NoShallowClone {
		// load state
		stateManager, err := api.NewCollectorStateManager(
			ctx,
			ctx.TaskContext().SyncPolicy(),
			"gitextractor",
			fmt.Sprintf(
				`{"RepoId: "%s","SkipCommitStat": %v, "SkipCommitFiles": %v}`,
				taskData.Options.RepoId,
				*taskData.Options.SkipCommitStat,
				*taskData.Options.SkipCommitFiles,
			),
		)
		if err != nil {
			return err
		}
		g.stateManager = stateManager
		since = stateManager.GetSince()
	}

	cmd, err := g.buildCloneCommand(ctx, localDir, since)
	if err != nil {
		return err
	}
	err = g.execCloneCommand(cmd)
	if err != nil {
		// it is likely that nothing to collect on incrmental mode
		if errors.Is(err, ErrShallowInfoProcessing) && g.stateManager != nil && g.stateManager.IsIncremental() {
			return ErrNoDataOnIncrementalMode
		}
		return err
	}

	// save state
	if g.stateManager != nil {
		return g.stateManager.Close()
	}
	return nil
}

func (g *GitcliCloner) buildCloneCommand(ctx plugin.SubTaskContext, localDir string, since *time.Time) (*exec.Cmd, errors.Error) {
	taskData := ctx.GetData().(*GitExtractorTaskData)
	args := []string{"clone", taskData.Options.Url, localDir, "--bare", "--progress"}
	env := []string{}
	// support proxy
	if taskData.ParsedURL.Scheme == "http" || taskData.ParsedURL.Scheme == "https" {
		if taskData.Options.Proxy != "" {
			env = append(env, fmt.Sprintf("HTTPS_PROXY=%s", taskData.Options.Proxy))
		}
	} else if taskData.ParsedURL.Scheme == "ssh" {
		var sshCmdArgs []string
		if taskData.Options.Proxy != "" {
			parsedProxyURL, e := url.Parse(taskData.Options.Proxy)
			if e != nil {
				return nil, errors.BadInput.Wrap(e, "failed to parse the proxy URL")
			}
			proxyCommand := "corkscrew"
			sshCmdArgs = append(sshCmdArgs, "-o", fmt.Sprintf(`ProxyCommand="%s %s %s %%h %%p"`, proxyCommand, parsedProxyURL.Hostname(), parsedProxyURL.Port()))
		}
		// support private key
		if taskData.Options.PrivateKey != "" {
			pkFile, err := os.CreateTemp("", "gitext-pk")
			if err != nil {
				g.logger.Error(err, "create temp private key file error")
				return nil, errors.Default.New("failed to handle the private key")
			}
			if _, e := pkFile.WriteString(taskData.Options.PrivateKey + "\n"); e != nil {
				g.logger.Error(err, "write private key file error")
				return nil, errors.Default.New("failed to write the  private key")
			}
			pkFile.Close()
			if e := os.Chmod(pkFile.Name(), 0600); e != nil {
				g.logger.Error(err, "chmod private key file error")
				return nil, errors.Default.New("failed to modify the private key")
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
					return nil, errors.Default.New("failed to decrypt the private key")
				}
			}
			defer os.Remove(pkFile.Name())
			sshCmdArgs = append(sshCmdArgs, fmt.Sprintf("-i %s -o StrictHostKeyChecking=no", pkFile.Name()))
		}
		if len(sshCmdArgs) > 0 {
			env = append(env, fmt.Sprintf("GIT_SSH_COMMAND=ssh %s", strings.Join(sshCmdArgs, " ")))
		}
	}
	// support time after and diff sync
	if since != nil {
		args = append(args, fmt.Sprintf("--shallow-since=%s", since.Format(time.RFC3339)))
	}
	// support skipping blobs collection
	if *taskData.Options.SkipCommitStat {
		args = append(args, "--filter=blob:none")
	}
	// fmt.Printf("args: %v\n", args)
	g.logger.Debug("git %v", args)
	cmd := exec.CommandContext(ctx.GetContext(), "git", args...)
	cmd.Env = env
	return cmd, nil
}

func (g *GitcliCloner) execCloneCommand(cmd *exec.Cmd) errors.Error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		g.logger.Error(err, "stdout pipe error")
		return errors.Default.New("stdout pipe error")
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		g.logger.Error(err, "stderr pipe error")
		return errors.Default.New("stderr pipe error")
	}
	combinedOutput := new(strings.Builder)
	stdoutScanner := bufio.NewScanner(stdout)
	stdoutScanner.Split(bufio.ScanLines)
	stderrScanner := bufio.NewScanner(stderr)
	stderrScanner.Split(bufio.ScanLines)
	done := make(chan bool)
	go func() {
		for stdoutScanner.Scan() {
			// TODO: extract progress?
			combinedOutput.WriteString(fmt.Sprintf("stdout: %s\n", stdoutScanner.Text()))
		}
		done <- true
	}()
	go func() {
		// TODO: extract progress?
		for stderrScanner.Scan() {
			combinedOutput.WriteString(fmt.Sprintf("stderr: %s\n", stderrScanner.Text()))
		}
		done <- true
	}()
	if e := cmd.Start(); e != nil {
		g.logger.Error(e, "failed to start\n%s", combinedOutput.String())
		return errors.Default.New("failed to start")
	}
	<-done
	<-done
	err = cmd.Wait()
	if err != nil {
		g.logger.Error(err, "git exited with error\n%s", combinedOutput.String())
		if strings.Contains(combinedOutput.String(), "stderr: fatal: error processing shallow info: 4") ||
			strings.Contains(combinedOutput.String(), "stderr: fatal: the remote end hung up unexpectedly") {
			return ErrShallowInfoProcessing
		}
		return errors.Default.New("git exit error")
	}
	return nil
}
