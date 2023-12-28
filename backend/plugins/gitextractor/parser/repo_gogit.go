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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/gitextractor/models"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"regexp"
)

var b1, b3, b2, b4 []string

type GoGitRepo struct {
	goGitStore   models.Store
	goGitRepo    *gogit.Repository
	goGitCleanUp func() // fixme
}

// CountTagsWithGoGit Count git tags subtask
func (r *GitRepo) CountTagsWithGoGit() (int, error) {
	repo := r.goGitRepo
	iter, err := repo.Tags()
	if err != nil {
		return 0, err
	}
	var tagsCount int
	iter.ForEach(func(reference *plumbing.Reference) error {
		tagsCount += 1
		return nil
	})
	return tagsCount, nil
}

// CountBranchesWithGoGit count the number of branches in a git repo
func (r *GitRepo) CountBranchesWithGoGit(ctx context.Context) (int, error) {
	repo := r.goGitRepo
	refIter, err := repo.Storer.IterReferences()
	if err != nil {
		return 0, err
	}
	branchIter := storer.NewReferenceFilteredIter(
		func(r *plumbing.Reference) bool {
			return r.Name().IsBranch() || r.Name().IsRemote()
		}, refIter)
	if err != nil {
		return 0, err
	}
	var branchesCount int

	headRef, err := repo.Head()
	if err != nil {
		return 0, err
	}
	branchIter.ForEach(func(reference *plumbing.Reference) error {
		if reference.Name() != headRef.Name() {
			branchesCount += 1
		}
		return nil
	})
	return branchesCount, errors.Convert(err)
}

// CountCommitsWithGoGit count the number of commits in a git repo
func (r *GitRepo) CountCommitsWithGoGit(ctx context.Context) (int, error) {
	repo := r.goGitRepo
	iter, err := repo.CommitObjects()
	if err != nil {
		return 0, err
	}
	var count int
	iter.ForEach(func(commit *object.Commit) error {
		count += 1
		return nil
	})
	return count, nil
}

// CollectTagsWithGoGit Collect Tags data
func (r *GitRepo) CollectTagsWithGoGit(subtaskCtx plugin.SubTaskContext) error {
	repo := r.goGitRepo
	store := r.goGitStore
	tagIter, err := repo.Tags()
	if err != nil {
		return err
	}
	if err := tagIter.ForEach(func(ref *plumbing.Reference) error {
		tagCommit := ref.Hash().String()
		_, err := repo.CommitObject(ref.Hash())
		if err != nil && errors.Is(err, plumbing.ErrObjectNotFound) {
			h, err := repo.ResolveRevision(plumbing.Revision(ref.Name()))
			if err != nil {
				return err
			}
			tagCommit = h.String()
		}
		name := ref.Name().String()
		if tagCommit != "" {
			codeRef := &code.Ref{
				DomainEntity: domainlayer.DomainEntity{Id: fmt.Sprintf("%s:%s", r.id, name)},
				RepoId:       r.id,
				Name:         name,
				CommitSha:    tagCommit,
				RefType:      TAG,
			}
			err = store.Refs(codeRef)
			if err != nil {
				return err
			}
			subtaskCtx.IncProgress(1)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// CollectBranchesWithGoGit Collect branch data
func (r *GitRepo) CollectBranchesWithGoGit(subtaskCtx plugin.SubTaskContext) error {
	repo := r.goGitRepo
	store := r.goGitStore
	refIter, err := repo.Storer.IterReferences()
	if err != nil {
		return err
	}
	branchIter := storer.NewReferenceFilteredIter(
		func(r *plumbing.Reference) bool {
			return r.Name().IsBranch() || r.Name().IsRemote()
		}, refIter)
	if err != nil {
		return err
	}
	headRef, err := repo.Head()
	if err != nil {
		return err
	}
	if err := branchIter.ForEach(func(ref *plumbing.Reference) error {
		name := ref.Name().Short()
		sha := ref.Hash().String()
		_, err := repo.CommitObject(ref.Hash())
		if err != nil && errors.Is(err, plumbing.ErrObjectNotFound) {
			// handle commit sha like "0000000000000000000000000000000000000000"
			h, err := repo.ResolveRevision(plumbing.Revision(ref.Name()))
			if err != nil {
				return err
			}
			sha = h.String()
		}
		codeRef := &code.Ref{
			DomainEntity: domainlayer.DomainEntity{Id: fmt.Sprintf("%s:%s", r.id, name)},
			RepoId:       r.id,
			Name:         name,
			CommitSha:    sha,
			RefType:      BRANCH,
			IsDefault:    ref.Name() == headRef.Name(),
		}
		if err := store.Refs(codeRef); err != nil {
			return err
		}
		subtaskCtx.IncProgress(1)
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// CollectCommitsWithGoGit Collect data from each commit, we can also get the diff line
func (r *GitRepo) CollectCommitsWithGoGit(subtaskCtx plugin.SubTaskContext) (err error) {
	componentMap, err := r.getComponentMap(subtaskCtx)
	if err != nil {
		return err
	}
	skipCommitFiles := subtaskCtx.GetConfigReader().GetBool(SkipCommitFiles)
	repo := r.goGitRepo
	store := r.goGitStore
	opts := object.DefaultDiffTreeOptions

	commitsObjectsIter, err := repo.CommitObjects()
	if err != nil {
		return err
	}
	if err := commitsObjectsIter.ForEach(func(commit *object.Commit) error {
		commitSha := commit.Hash.String()

		r.logger.Debug("process commit: %s", commitSha)
		b2 = append(b2, commitSha)

		codeCommit := &code.Commit{
			Sha:            commitSha,
			Message:        commit.Message,
			AuthorName:     commit.Author.Name,
			AuthorEmail:    commit.Author.Email,
			AuthorId:       commit.Author.Email,
			AuthoredDate:   commit.Author.When,
			CommitterName:  commit.Committer.Name,
			CommitterEmail: commit.Committer.Email,
			CommitterId:    commit.Committer.Email,
			CommittedDate:  commit.Committer.When,
		}
		if err = r.storeParentCommitsWithGoGit(commitSha, commit); err != nil {
			return err
		}

		var parent *object.Commit
		if commit.NumParents() > 0 {
			parent, err = commit.Parent(0)
			if err != nil {
				return err
			}
		}
		stats, err := commit.Stats()
		if err != nil {
			return err
		} else {
			for _, stat := range stats {
				fmt.Printf("stat: %+v\n", stat)
				codeCommit.Additions += stat.Addition
				codeCommit.Deletions += stat.Deletion
			}
		}
		err = store.Commits(codeCommit)
		if err != nil {
			return err
		}

		if err := r.getDiffComparedToParentWithGoGit(subtaskCtx.GetContext(), skipCommitFiles, codeCommit.Sha, commit, parent, opts, componentMap); err != nil {
			return err
		}

		codeRepoCommit := &code.RepoCommit{
			RepoId:    r.id,
			CommitSha: commitSha,
		}
		err = store.RepoCommits(codeRepoCommit)
		if err != nil {
			return err
		}

		subtaskCtx.IncProgress(1)
		return nil
	}); err != nil {
		return err
	}
	return
}

func (r *GitRepo) storeParentCommitsWithGoGit(commitSha string, commit *object.Commit) error {
	if commit == nil {
		return nil
	}
	var commitParents []*code.CommitParent
	for i := 0; i < commit.NumParents(); i++ {
		parent, err := commit.Parent(i)
		if err != nil {
			return err
		}
		if parent != nil {
			if parentCommitSha := parent.Hash.String(); parentCommitSha != "" {
				commitParents = append(commitParents, &code.CommitParent{
					CommitSha:       commitSha,
					ParentCommitSha: parentCommitSha,
				})
			}
		}
	}
	return r.store.CommitParents(commitParents)
}

func (r *GitRepo) getDiffComparedToParentWithGoGit(ctx context.Context, skipCommitFiles bool, commitSha string, commit *object.Commit, parent *object.Commit, opts *object.DiffTreeOptions, componentMap map[string]*regexp.Regexp) (err error) {
	if skipCommitFiles {
		return nil
	}

	if _, err := commit.Stats(); err != nil {
		return err
	}

	commitTree, err := commit.Tree()
	if err != nil {
		return err
	}
	parentTree, err := parent.Tree()
	if err != nil {
		return err
	}
	changes, err := object.DiffTreeWithOptions(ctx, parentTree, commitTree, opts)
	if err != nil {
		return err
	}

	if err = r.storeCommitFilesFromDiffWitGoGit(commitSha, changes, componentMap); err != nil {
		return err
	}
	return nil
}

func (r *GitRepo) storeCommitFilesFromDiffWitGoGit(commitSha string, changes object.Changes, componentMap map[string]*regexp.Regexp) (err error) {

	store := r.goGitStore
	var commitFile *code.CommitFile
	var commitFileComponent *code.CommitFileComponent

	for _, change := range changes {
		if commitFile != nil {
			if err := store.CommitFiles(commitFile); err != nil {
				r.logger.Error(err, "CommitFiles error")
				return err
			}
		}
		commitFile = new(code.CommitFile)
		commitFile.CommitSha = commitSha
		_, toFile, err := change.Files()
		if err != nil {
			return err
		}
		if toFile != nil {
			filePath := toFile.Name
			commitFile.FilePath = filePath
			// With some long path,the varchar(255) was not enough both ID and file_path
			// So we use the hash to compress the path in ID and add length of file_path.
			// Use commitSha and the sha256 of FilePath to create id
			// fixme: maybe we can use file's hash directly
			shaFilePath := sha256.New()
			shaFilePath.Write([]byte(filePath))
			commitFile.Id = commitSha + ":" + hex.EncodeToString(shaFilePath.Sum(nil))
		}

		// handle component
		commitFileComponent = new(code.CommitFileComponent)
		for component, reg := range componentMap {
			if reg.MatchString(commitFile.FilePath) {
				commitFileComponent.ComponentName = component
				break
			}
		}
		commitFileComponent.CommitFileId = commitFile.Id
		if commitFileComponent.ComponentName == "" {
			commitFileComponent.ComponentName = "Default"
		}

	}
	//err = diff.ForEach(func(file git.DiffDelta, progress float64) (
	//	git.DiffForEachHunkCallback, error) {
	//	if commitFile != nil {
	//		err = r.store.CommitFiles(commitFile)
	//		if err != nil {
	//			r.logger.Error(err, "CommitFiles error")
	//			return nil, err
	//		}
	//	}
	//
	//	commitFile = new(code.CommitFile)
	//	commitFile.CommitSha = commitSha
	//	commitFile.FilePath = file.NewFile.Path
	//
	//	// With some long path,the varchar(255) was not enough both ID and file_path
	//	// So we use the hash to compress the path in ID and add length of file_path.
	//	// Use commitSha and the sha256 of FilePath to create id
	//	shaFilePath := sha256.New()
	//	shaFilePath.Write([]byte(file.NewFile.Path))
	//	commitFile.Id = commitSha + ":" + hex.EncodeToString(shaFilePath.Sum(nil))
	//
	//	commitFileComponent = new(code.CommitFileComponent)
	//	for component, reg := range componentMap {
	//		if reg.MatchString(commitFile.FilePath) {
	//			commitFileComponent.ComponentName = component
	//			break
	//		}
	//	}
	//	commitFileComponent.CommitFileId = commitFile.Id
	//	if commitFileComponent.ComponentName == "" {
	//		commitFileComponent.ComponentName = "Default"
	//	}
	//	return func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
	//		return func(line git.DiffLine) error {
	//			if line.Origin == git.DiffLineAddition {
	//				commitFile.Additions += line.NumLines
	//			}
	//			if line.Origin == git.DiffLineDeletion {
	//				commitFile.Deletions += line.NumLines
	//			}
	//			return nil
	//		}, nil
	//	}, nil
	//}, git.DiffDetailLines)

	if commitFileComponent != nil {
		err = store.CommitFileComponents(commitFileComponent)
		if err != nil {
			r.logger.Error(err, "CommitFileComponents error")
		}
	}
	if commitFile != nil {
		err = store.CommitFiles(commitFile)
		if err != nil {
			r.logger.Error(err, "CommitFiles error")
		}
	}
	return errors.Convert(err)
}
