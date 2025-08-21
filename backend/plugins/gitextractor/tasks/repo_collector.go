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

package tasks

import (
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/gitextractor/parser"
)

func CollectGitCommits(subTaskCtx plugin.SubTaskContext) errors.Error {
	if subTaskCtx.TaskContext().GetData().(*parser.GitExtractorTaskData).SkipAllSubtasks {
		return nil
	}
	repo := getGitRepo(subTaskCtx)
	if count, err := repo.CountCommits(subTaskCtx.GetContext()); err != nil {
		subTaskCtx.GetLogger().Error(err, "unable to get commit count")
		subTaskCtx.SetProgress(0, -1)
		return errors.Convert(err)
	} else {
		subTaskCtx.SetProgress(0, count)
	}
	return errors.Convert(repo.CollectCommits(subTaskCtx))
}

func CollectGitBranches(subTaskCtx plugin.SubTaskContext) errors.Error {
	if subTaskCtx.TaskContext().GetData().(*parser.GitExtractorTaskData).SkipAllSubtasks {
		return nil
	}
	repo := getGitRepo(subTaskCtx)
	if count, err := repo.CountBranches(subTaskCtx.GetContext()); err != nil {
		subTaskCtx.GetLogger().Error(err, "unable to get branch count")
		subTaskCtx.SetProgress(0, -1)
		return errors.Convert(err)
	} else {
		subTaskCtx.SetProgress(0, count)
	}
	return errors.Convert(repo.CollectBranches(subTaskCtx))
}

func CollectGitTags(subTaskCtx plugin.SubTaskContext) errors.Error {
	if subTaskCtx.TaskContext().GetData().(*parser.GitExtractorTaskData).SkipAllSubtasks {
		return nil
	}
	repo := getGitRepo(subTaskCtx)
	if count, err := repo.CountTags(subTaskCtx.GetContext()); err != nil {
		subTaskCtx.GetLogger().Error(err, "unable to get tag count")
		subTaskCtx.SetProgress(0, -1)
		return errors.Convert(err)
	} else {
		subTaskCtx.SetProgress(0, count)
	}
	return errors.Convert(repo.CollectTags(subTaskCtx))
}

func CollectGitDiffLines(subTaskCtx plugin.SubTaskContext) errors.Error {
	if subTaskCtx.TaskContext().GetData().(*parser.GitExtractorTaskData).SkipAllSubtasks {
		return nil
	}
	repo := getGitRepo(subTaskCtx)
	opt := subTaskCtx.GetData().(*parser.GitExtractorTaskData).Options
	if !*opt.SkipCommitStat {
		subTaskCtx.SetProgress(0, -1)
		return errors.Convert(repo.CollectDiffLine(subTaskCtx))
	}
	return nil
}

func getGitRepo(subTaskCtx plugin.SubTaskContext) parser.RepoCollector {
	taskData, ok := subTaskCtx.GetData().(*parser.GitExtractorTaskData)
	if !ok {
		subTaskCtx.GetLogger().Error(nil, "git repo reference not found on context")
		return nil
	}
	if taskData.GitRepo == nil {
		subTaskCtx.GetLogger().Error(nil, "git repo is empty, skipping Collect Commits subtask")
		return nil
	}
	return taskData.GitRepo
}

var CollectGitCommitMeta = plugin.SubTaskMeta{
	Name:             "Collect Commits",
	EntryPoint:       CollectGitCommits,
	EnabledByDefault: true,
	Description:      "collect git commits into Domain Layer Tables",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE, plugin.DOMAIN_TYPE_CROSS},
	Dependencies:     []*plugin.SubTaskMeta{&CloneGitRepoMeta},
}

var CollectGitBranchMeta = plugin.SubTaskMeta{
	Name:             "Collect Branches",
	EntryPoint:       CollectGitBranches,
	EnabledByDefault: true,
	Description:      "collect git branch into Domain Layer Tables",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
	Dependencies:     []*plugin.SubTaskMeta{&CloneGitRepoMeta},
}

var CollectGitTagMeta = plugin.SubTaskMeta{
	Name:             "Collect Tags",
	EntryPoint:       CollectGitTags,
	EnabledByDefault: true,
	Description:      "collect git tag into Domain Layer Tables",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
	Dependencies:     []*plugin.SubTaskMeta{&CloneGitRepoMeta},
}

var CollectGitDiffLineMeta = plugin.SubTaskMeta{
	Name:             "Collect DiffLine",
	EntryPoint:       CollectGitDiffLines,
	EnabledByDefault: false,
	Description:      "collect git commit diff line into Domain Layer Tables",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
	Dependencies:     []*plugin.SubTaskMeta{&CloneGitRepoMeta},
}
