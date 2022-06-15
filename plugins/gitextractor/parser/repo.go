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
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitextractor/models"
	git "github.com/libgit2/git2go/v33"
)

type GitRepo struct {
	store   models.Store
	ctx     context.Context
	logger  core.Logger
	id      string
	repo    *git.Repository
	cleanup func()
}

func (r *GitRepo) CollectAll(subtaskCtx core.SubTaskContext) error {
	subtaskCtx.SetProgress(0, -1)
	err := r.CollectTags(subtaskCtx)
	if err != nil {
		return err
	}
	err = r.CollectBranches(subtaskCtx)
	if err != nil {
		return err
	}
	return r.CollectCommits(subtaskCtx)
}

func (r *GitRepo) Close() error {
	defer func() {
		if r.cleanup != nil {
			r.cleanup()
		}
	}()
	return r.store.Close()
}

func (r *GitRepo) CountTags() (int, error) {
	tags, err := r.repo.Tags.List()
	if err != nil {
		return 0, err
	}
	return len(tags), nil
}

func (r *GitRepo) CountBranches() (int, error) {
	var branchIter *git.BranchIterator
	branchIter, err := r.repo.NewBranchIterator(git.BranchAll)
	if err != nil {
		return 0, err
	}
	count := 0
	err = branchIter.ForEach(func(branch *git.Branch, branchType git.BranchType) error {
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		default:
			break
		}
		if branch.IsBranch() || branch.IsRemote() {
			count++
		}
		return nil
	})
	return count, err
}

func (r *GitRepo) CountCommits() (int, error) {
	odb, err := r.repo.Odb()
	if err != nil {
		return 0, err
	}
	count := 0
	err = odb.ForEach(func(id *git.Oid) error {
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		default:
			break
		}
		commit, _ := r.repo.LookupCommit(id)
		if commit != nil {
			count++
		}
		return nil
	})
	return count, err
}

func (r *GitRepo) CollectTags(subtaskCtx core.SubTaskContext) error {
	return r.repo.Tags.Foreach(func(name string, id *git.Oid) error {
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		default:
			break
		}
		var err1 error
		var tag *git.Tag
		var tagCommit string
		tag, _ = r.repo.LookupTag(id)
		if tag != nil {
			tagCommit = tag.TargetId().String()
		} else {
			tagCommit = id.String()
		}
		r.logger.Info("tagCommit", tagCommit)
		if tagCommit != "" {
			ref := &code.Ref{
				DomainEntity: domainlayer.DomainEntity{Id: fmt.Sprintf("%s:%s", r.id, name)},
				RepoId:       r.id,
				Name:         name,
				CommitSha:    tagCommit,
				RefType:      TAG,
			}
			err1 = r.store.Refs(ref)
			if err1 != nil {
				return err1
			}
			subtaskCtx.IncProgress(1)
		}
		return nil
	})
}

func (r *GitRepo) CollectBranches(subtaskCtx core.SubTaskContext) error {
	var repoInter *git.BranchIterator
	repoInter, err := r.repo.NewBranchIterator(git.BranchAll)
	if err != nil {
		return err
	}
	return repoInter.ForEach(func(branch *git.Branch, branchType git.BranchType) error {
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		default:
			break
		}
		if branch.IsBranch() || branch.IsRemote() {
			name, err1 := branch.Name()
			if err1 != nil {
				return err1
			}
			var sha string
			if oid := branch.Target(); oid != nil {
				sha = oid.String()
			}
			ref := &code.Ref{
				DomainEntity: domainlayer.DomainEntity{Id: fmt.Sprintf("%s:%s", r.id, name)},
				RepoId:       r.id,
				Name:         name,
				CommitSha:    sha,
				RefType:      BRANCH,
			}
			ref.IsDefault, _ = branch.IsHead()
			err1 = r.store.Refs(ref)
			if err1 != nil {
				return err1
			}
			subtaskCtx.IncProgress(1)
			return nil
		}
		return nil
	})
}

func (r *GitRepo) CollectCommits(subtaskCtx core.SubTaskContext) error {
	opts, err := getDiffOpts()
	if err != nil {
		return err
	}
	odb, err := r.repo.Odb()
	if err != nil {
		return err
	}
	return odb.ForEach(func(id *git.Oid) error {
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		default:
			break
		}
		commit, _ := r.repo.LookupCommit(id)
		if commit == nil {
			return nil
		}
		commitSha := commit.Id().String()
		r.logger.Debug("process commit: %s", commitSha)
		c := &code.Commit{
			Sha:     commitSha,
			Message: commit.Message(),
		}
		author := commit.Author()
		if author != nil {
			c.AuthorName = author.Name
			c.AuthorEmail = author.Email
			c.AuthorId = author.Email
			c.AuthoredDate = author.When
		}
		committer := commit.Committer()
		if committer != nil {
			c.CommitterName = committer.Name
			c.CommitterEmail = committer.Email
			c.CommitterId = committer.Email
			c.CommittedDate = committer.When
		}
		if err != r.storeParentCommits(commitSha, commit) {
			return err
		}
		if commit.ParentCount() > 0 {
			parent := commit.Parent(0)
			if parent != nil {
				var stats *git.DiffStats
				if stats, err = r.getDiffComparedToParent(c.Sha, commit, parent, opts); err != nil {
					return err
				} else {
					c.Additions += stats.Insertions()
					c.Deletions += stats.Deletions()
				}
			}
		}
		err = r.store.Commits(c)
		if err != nil {
			return err
		}
		repoCommit := &code.RepoCommit{
			RepoId:    r.id,
			CommitSha: c.Sha,
		}
		err = r.store.RepoCommits(repoCommit)
		if err != nil {
			return err
		}
		subtaskCtx.IncProgress(1)
		return nil
	})
}

func (r *GitRepo) storeParentCommits(commitSha string, commit *git.Commit) error {
	var commitParents []*code.CommitParent
	for i := uint(0); i < commit.ParentCount(); i++ {
		parent := commit.Parent(i)
		if parent != nil {
			if parentId := parent.Id(); parentId != nil {
				commitParents = append(commitParents, &code.CommitParent{
					CommitSha:       commitSha,
					ParentCommitSha: parentId.String(),
				})
			}
		}
	}
	return r.store.CommitParents(commitParents)
}

func (r *GitRepo) getDiffComparedToParent(commitSha string, commit *git.Commit, parent *git.Commit, opts *git.DiffOptions) (*git.DiffStats, error) {
	var err error
	var parentTree, tree *git.Tree
	parentTree, err = parent.Tree()
	if err != nil {
		return nil, err
	}
	tree, err = commit.Tree()
	if err != nil {
		return nil, err
	}
	var diff *git.Diff
	diff, err = r.repo.DiffTreeToTree(parentTree, tree, opts)
	if err != nil {
		return nil, err
	}
	err = r.storeCommitFilesFromDiff(commitSha, diff)
	if err != nil {
		return nil, err
	}
	var stats *git.DiffStats
	stats, err = diff.Stats()
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *GitRepo) storeCommitFilesFromDiff(commitSha string, diff *git.Diff) error {
	var commitFile *code.CommitFile
	var err error
	err = diff.ForEach(func(file git.DiffDelta, progress float64) (
		git.DiffForEachHunkCallback, error) {
		if commitFile != nil {
			err = r.store.CommitFiles(commitFile)
			if err != nil {
				r.logger.Error("CommitFiles error:", err)
				return nil, err
			}
		}
		commitFile = new(code.CommitFile)
		commitFile.CommitSha = commitSha
		commitFile.FilePath = file.NewFile.Path
		return func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
			return func(line git.DiffLine) error {
				if line.Origin == git.DiffLineAddition {
					commitFile.Additions += line.NumLines
				}
				if line.Origin == git.DiffLineDeletion {
					commitFile.Deletions += line.NumLines
				}
				return nil
			}, nil
		}, nil
	}, git.DiffDetailLines)
	if commitFile != nil {
		err = r.store.CommitFiles(commitFile)
		if err != nil {
			r.logger.Error("CommitFiles error:", err)
		}
	}
	return err
}

func getDiffOpts() (*git.DiffOptions, error) {
	opts, err := git.DefaultDiffOptions()
	if err != nil {
		return nil, err
	}
	opts.NotifyCallback = func(diffSoFar *git.Diff, delta git.DiffDelta, matchedPathSpec string) error {
		return nil
	}
	return &opts, nil
}
