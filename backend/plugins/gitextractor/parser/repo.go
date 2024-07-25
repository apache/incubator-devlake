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

	"github.com/apache/incubator-devlake/core/plugin"
)

const (
	BRANCH = "BRANCH"
	TAG    = "TAG"
)

type RepoCollector interface {
	SetCleanUp(func()) error
	Close(ctx context.Context) error

	CollectAll(subtaskCtx plugin.SubTaskContext) error

	CountTags(ctx context.Context) (int, error)
	CountBranches(ctx context.Context) (int, error)
	CountCommits(ctx context.Context) (int, error)

	CollectTags(subtaskCtx plugin.SubTaskContext) error
	CollectBranches(subtaskCtx plugin.SubTaskContext) error
	CollectCommits(subtaskCtx plugin.SubTaskContext) error
	CollectDiffLine(subtaskCtx plugin.SubTaskContext) error
}
