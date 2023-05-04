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

package bridge

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
)

type CmdInvoker struct {
	resolveCmd  func(methodName string, args ...string) (string, []string)
	cancelled   bool
	workingPath string
}

func NewCmdInvoker(execPath string) *CmdInvoker {
	// Split the path into dir and file
	dir, file := path.Split(execPath)
	resolveCmd := func(methodName string, args ...string) (string, []string) {
		allArgs := []string{methodName}
		allArgs = append(allArgs, args...)
		return fmt.Sprintf("./%s", file), allArgs
	}

	return &CmdInvoker{
		resolveCmd:  resolveCmd,
		workingPath: dir,
	}
}

func (c *CmdInvoker) Call(methodName string, ctx plugin.ExecContext, args ...any) *CallResult {
	serializedArgs, err := serialize(args...)
	if err != nil {
		return &CallResult{
			Err: err,
		}
	}
	executable, inputArgs := c.resolveCmd(methodName, serializedArgs...)
	cmdCtx := DefaultContext.GetContext()
	cmd := exec.CommandContext(cmdCtx, executable, inputArgs...)
	if c.workingPath != "" {
		cmd.Dir = c.workingPath
	}
	response, err := utils.RunProcess(cmd, &utils.RunProcessOptions{
		OnStdout: func(b []byte) {
			msg := string(b)
			c.logRemoteMessage(ctx.GetLogger(), msg)
		},
		OnStderr: func(b []byte) {
			msg := string(b)
			c.logRemoteError(ctx.GetLogger(), msg)
		},
		UseFdOut: true,
	})
	if err != nil {
		return NewCallResult(nil, err)
	}
	err = response.GetError()
	if err != nil {
		return &CallResult{
			Err: errors.Default.Wrap(err, fmt.Sprintf("get error when invoking remote function %s", methodName)),
		}
	}
	return NewCallResult(response.GetFdOut(), nil)
}

func (c *CmdInvoker) Stream(methodName string, ctx plugin.ExecContext, args ...any) *MethodStream {
	recvChannel := make(chan *StreamResult)
	stream := &MethodStream{
		outbound: nil,
		inbound:  recvChannel,
	}
	serializedArgs, err := serialize(args...)
	if err != nil {
		recvChannel <- NewStreamResult(nil, err)
		return stream
	}
	executable, inputArgs := c.resolveCmd(methodName, serializedArgs...)
	cmdCtx := DefaultContext.GetContext() // grabbing context off of ctx kills the cmd after a couple of seconds... why?
	cmd := exec.CommandContext(cmdCtx, executable, inputArgs...)
	if c.workingPath != "" {
		cmd.Dir = c.workingPath
	}
	processHandle, err := utils.StreamProcess(cmd, &utils.StreamProcessOptions{
		OnStdout: func(b []byte) {
			msg := string(b)
			c.logRemoteMessage(ctx.GetLogger(), msg)
		},
		OnStderr: func(b []byte) {
			msg := string(b)
			c.logRemoteError(ctx.GetLogger(), msg)
		},
		UseFdOut: true,
	})
	if err != nil {
		recvChannel <- NewStreamResult(nil, err)
		return stream
	}
	go func() {
		defer close(recvChannel)
		for msg := range processHandle.Receive() {
			if err = msg.GetError(); err != nil {
				recvChannel <- NewStreamResult(nil, err)
			}
			if !c.cancelled {
				select {
				case <-ctx.GetContext().Done():
					err = processHandle.Cancel()
					if err != nil {
						recvChannel <- NewStreamResult(nil, errors.Default.Wrap(err, "error cancelling python target"))
						return
					}
					c.cancelled = true
					// continue until the stream gets closed by the child
				default:
				}
			}
			response := msg.GetFdOut()
			if response != nil {
				recvChannel <- NewStreamResult(response, nil)
			}
		}
	}()
	return stream
}

func serialize(args ...any) ([]string, errors.Error) {
	var serializedArgs []string
	for _, arg := range args {
		serializedArg, err := json.Marshal(arg)
		if err != nil {
			return nil, errors.Convert(err)
		}
		serializedArgs = append(serializedArgs, string(serializedArg))
	}
	return serializedArgs, nil
}

func (c *CmdInvoker) logRemoteMessage(logger log.Logger, msg string) {
	logger.Info(msg)
}

func (c *CmdInvoker) logRemoteError(logger log.Logger, msg string) {
	logger.Error(nil, msg)
}

var _ Invoker = (*CmdInvoker)(nil)
