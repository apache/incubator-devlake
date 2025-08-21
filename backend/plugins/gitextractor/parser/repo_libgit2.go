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
	"regexp"
	"sort"
	"strconv"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/gitextractor/models"

	git "github.com/libgit2/git2go/v33"
)

var _ RepoCollector = (*Libgit2RepoCollector)(nil)

var TypeNotMatchError = "the requested type does not match the type in the ODB"

type Libgit2RepoCollector struct {
	id     string
	logger log.Logger

	store   models.Store
	repo    *git.Repository
	cleanup func()
}

func NewLibgit2RepoCollector(localDir string, repoId string, store models.Store, logger log.Logger) (*Libgit2RepoCollector, errors.Error) {
	repo, err := git.OpenRepository(localDir)
	if err != nil {
		return nil, errors.Convert(err)
	}
	return &Libgit2RepoCollector{
		id:     repoId,
		logger: logger,
		store:  store,
		repo:   repo,
	}, nil
}

func (r *Libgit2RepoCollector) SetCleanUp(f func()) error {
	if f != nil {
		r.cleanup = f
	}
	return nil
}

// CollectAll The main parser subtask
func (r *Libgit2RepoCollector) CollectAll(subtaskCtx plugin.SubTaskContext) error {
	subtaskCtx.SetProgress(0, -1)
	err := r.CollectTags(subtaskCtx)
	if err != nil {
		return err
	}
	err = r.CollectBranches(subtaskCtx)
	if err != nil {
		return err
	}
	err = r.CollectCommits(subtaskCtx)
	if err != nil {
		return err
	}
	opt := subtaskCtx.GetData().(*GitExtractorTaskData).Options
	if !*opt.SkipCommitStat {
		return r.CollectDiffLine(subtaskCtx)
	}
	return nil
}

// Close resources
func (r *Libgit2RepoCollector) Close(ctx context.Context) error {
	defer func() {
		if r.cleanup != nil {
			r.cleanup()
		}
	}()
	return r.store.Close()
}

// CountTags Count git tags subtask
func (r *Libgit2RepoCollector) CountTags(ctx context.Context) (int, error) {
	tags, err := r.repo.Tags.List()
	if err != nil {
		return 0, errors.Convert(err)
	}
	return len(tags), nil
}

// CountBranches count the number of branches in a git repo
func (r *Libgit2RepoCollector) CountBranches(ctx context.Context) (int, error) {
	var branchIter *git.BranchIterator
	branchIter, err := r.repo.NewBranchIterator(git.BranchAll)
	if err != nil {
		return 0, errors.Convert(err)
	}
	count := 0
	err = branchIter.ForEach(func(branch *git.Branch, branchType git.BranchType) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if branch.IsBranch() || branch.IsRemote() {
			isHead, _ := branch.IsHead()
			if !isHead {
				count++
			}
		}
		return nil
	})
	return count, errors.Convert(err)
}

// CountCommits count the number of commits in a git repo
func (r *Libgit2RepoCollector) CountCommits(ctx context.Context) (int, error) {
	odb, err := r.repo.Odb()
	if err != nil {
		return 0, errors.Convert(err)
	}
	count := 0
	err = odb.ForEach(func(id *git.Oid) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		commit, e := r.repo.LookupCommit(id)
		if e != nil && e.Error() != TypeNotMatchError {
			return errors.Convert(e)
		}
		if commit != nil {
			count++
		}
		return nil
	})
	return count, errors.Convert(err)
}

// CollectTags Collect Tags data
func (r *Libgit2RepoCollector) CollectTags(subtaskCtx plugin.SubTaskContext) error {
	return errors.Convert(r.repo.Tags.Foreach(func(name string, id *git.Oid) error {
		select {
		case <-subtaskCtx.GetContext().Done():
			return subtaskCtx.GetContext().Err()
		default:
		}
		var err1 error
		var tag *git.Tag
		var tagCommit string
		tag, err1 = r.repo.LookupTag(id)
		if err1 != nil && err1.Error() != TypeNotMatchError {
			return errors.Convert(err1)
		}
		if tag != nil {
			tagCommit = tag.TargetId().String()
		} else {
			tagCommit = id.String()
		}
		r.logger.Debug("tagCommit: %s", tagCommit)
		if tagCommit != "" {
			ref := &code.Ref{
				DomainEntityExtended: domainlayer.DomainEntityExtended{Id: fmt.Sprintf("%s:%s", r.id, name)},
				RepoId:               r.id,
				Name:                 name,
				CommitSha:            tagCommit,
				RefType:              TAG,
			}
			err1 = r.store.Refs(ref)
			if err1 != nil {
				return err1
			}
			subtaskCtx.IncProgress(1)
		}
		return nil
	}))
}

// CollectBranches Collect branch data
func (r *Libgit2RepoCollector) CollectBranches(subtaskCtx plugin.SubTaskContext) error {
	var repoInter *git.BranchIterator
	repoInter, err := r.repo.NewBranchIterator(git.BranchAll)
	if err != nil {
		return errors.Convert(err)
	}
	return errors.Convert(repoInter.ForEach(func(branch *git.Branch, branchType git.BranchType) error {
		select {
		case <-subtaskCtx.GetContext().Done():
			return subtaskCtx.GetContext().Err()
		default:
		}
		if branch.IsBranch() || branch.IsRemote() {
			name, err1 := branch.Name()
			if err1 != nil && err1.Error() != TypeNotMatchError {
				return err1
			}
			var sha string
			if oid := branch.Target(); oid != nil {
				sha = oid.String()
			}
			ref := &code.Ref{
				DomainEntityExtended: domainlayer.DomainEntityExtended{Id: fmt.Sprintf("%s:%s", r.id, name)},
				RepoId:               r.id,
				Name:                 name,
				CommitSha:            sha,
				RefType:              BRANCH,
			}
			// commit sha may be an empty string, it's unexpected.
			ref.IsDefault, err1 = branch.IsHead()
			if err1 != nil && err1.Error() != TypeNotMatchError {
				return err1
			}
			err1 = r.store.Refs(ref)
			if err1 != nil && err1.Error() != TypeNotMatchError {
				return err1
			}
			subtaskCtx.IncProgress(1)
			return nil
		}
		return nil
	}))
}

// CollectCommits Collect data from each commit, we can also get the diff line
func (r *Libgit2RepoCollector) CollectCommits(subtaskCtx plugin.SubTaskContext) error {
	taskOpts := subtaskCtx.GetData().(*GitExtractorTaskData).Options
	opts, err := getDiffOpts()
	if err != nil {
		return err
	}
	db := subtaskCtx.GetDal()
	components := make([]code.Component, 0)
	err = db.All(&components, dal.From(components), dal.Where("repo_id= ?", r.id))
	if err != nil {
		return err
	}
	componentMap := make(map[string]*regexp.Regexp)
	for _, component := range components {
		componentMap[component.Name] = regexp.MustCompile(component.PathRegex)
	}
	odb, err := errors.Convert01(r.repo.Odb())
	if err != nil {
		return err
	}
	return errors.Convert(odb.ForEach(func(id *git.Oid) error {
		select {
		case <-subtaskCtx.GetContext().Done():
			return subtaskCtx.GetContext().Err()
		default:
		}
		commit, err1 := r.repo.LookupCommit(id)
		if err1 != nil && err1.Error() != TypeNotMatchError {
			return errors.Convert(err1)
		}
		if commit == nil {
			return nil
		}
		var parent *git.Commit
		if commit.ParentCount() > 0 {
			parent = commit.Parent(0)
			// Skip calculating commit statistics when there are parent commits, but the first one cannot be fetched from the ODB.
			// This usually happens during a shallow clone for incremental collection. Otherwise, we might end up overwriting
			// the correct addition/deletion data in the database with an absurdly large addition number.
			if parent == nil {
				r.logger.Info("skip commit %s because it has no parent commit", commit.Id().String())
				return nil
			}
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
		err = r.storeParentCommits(commitSha, commit)
		if err != nil {
			return err
		}

		if !*taskOpts.SkipCommitStat {
			var stats *git.DiffStats
			if stats, err = r.getDiffComparedToParent(taskOpts, c.Sha, commit, parent, opts, componentMap); err != nil {
				return err
			}
			r.logger.Debug("state: %#+v\n", stats.Deletions())
			c.Additions += stats.Insertions()
			c.Deletions += stats.Deletions()
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
	}))
}

func (r *Libgit2RepoCollector) storeParentCommits(commitSha string, commit *git.Commit) errors.Error {
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

func (r *Libgit2RepoCollector) getDiffComparedToParent(taskOpts *GitExtractorOptions, commitSha string, commit *git.Commit, parent *git.Commit, opts *git.DiffOptions, componentMap map[string]*regexp.Regexp) (*git.DiffStats, errors.Error) {
	var err error
	var parentTree, tree *git.Tree
	if parent != nil {
		parentTree, err = parent.Tree()
	}
	if err != nil {
		return nil, errors.Convert(err)
	}
	tree, err = commit.Tree()
	if err != nil {
		return nil, errors.Convert(err)
	}
	var diff *git.Diff
	diff, err = r.repo.DiffTreeToTree(parentTree, tree, opts)
	if err != nil {
		return nil, errors.Convert(err)
	}
	if !*taskOpts.SkipCommitFiles {
		err = r.storeCommitFilesFromDiff(commitSha, diff, componentMap)
		if err != nil {
			return nil, errors.Convert(err)
		}
	}
	var stats *git.DiffStats
	stats, err = diff.Stats()
	if err != nil {
		return nil, errors.Convert(err)
	}
	return stats, nil
}

func (r *Libgit2RepoCollector) storeCommitFilesFromDiff(commitSha string, diff *git.Diff, componentMap map[string]*regexp.Regexp) errors.Error {
	var commitFile *code.CommitFile
	var commitFileComponent *code.CommitFileComponent
	var err error
	err = diff.ForEach(func(file git.DiffDelta, progress float64) (
		git.DiffForEachHunkCallback, error) {
		if commitFile != nil {
			err = r.store.CommitFiles(commitFile)
			if err != nil {
				r.logger.Error(err, "CommitFiles error")
				return nil, err
			}
		}

		commitFile = new(code.CommitFile)
		commitFile.CommitSha = commitSha
		commitFile.FilePath = file.NewFile.Path

		// With some long path,the varchar(255) was not enough both ID and file_path
		// So we use the hash to compress the path in ID and add length of file_path.
		// Use commitSha and the sha256 of FilePath to create id
		shaFilePath := sha256.New()
		shaFilePath.Write([]byte(file.NewFile.Path))
		commitFile.Id = commitSha + ":" + hex.EncodeToString(shaFilePath.Sum(nil))

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
	if commitFileComponent != nil {
		err = r.store.CommitFileComponents(commitFileComponent)
		if err != nil {
			r.logger.Error(err, "CommitFileComponents error")
		}
	}
	if commitFile != nil {
		err = r.store.CommitFiles(commitFile)
		if err != nil {
			r.logger.Error(err, "CommitFiles error")
		}
	}
	return errors.Convert(err)
}

// CollectDiffLine get line diff data from a specific branch
func (r *Libgit2RepoCollector) CollectDiffLine(subtaskCtx plugin.SubTaskContext) error {
	//Using this subtask,we can get every line change in every commit.
	//We maintain a snapshot structure to get which commit each deleted line belongs to
	snapshot := make(map[string] /*file path*/ *models.FileBlame)
	repo := r.repo
	//step 1. get the reverse commit list
	commitList := make([]git.Commit, 0)
	// get current head commit sha, default is master branch
	// check branch, if not master, checkout to branch's head
	commitOid, err1 := repo.Head()
	if err1 != nil && err1.Error() != TypeNotMatchError {
		return errors.Convert(err1)
	}
	//get head commit object and add into commitList
	commit, err1 := repo.LookupCommit(commitOid.Target())
	if err1 != nil && err1.Error() != TypeNotMatchError {
		return errors.Convert(err1)
	}
	commitList = append(commitList, *commit)
	// if current head has parents, get parent commitsha
	for commit != nil && commit.ParentCount() > 0 {
		pid := commit.ParentId(0)
		commit, err1 = repo.LookupCommit(pid)
		if err1 != nil && err1.Error() != TypeNotMatchError {
			return errors.Convert(err1)
		}
		commitList = append(commitList, *commit)
	}
	// reverse commitList
	for i, j := 0, len(commitList)-1; i < j; i, j = i+1, j-1 {
		commitList[i], commitList[j] = commitList[j], commitList[i]
	}
	//step 2. get the diff of each commit
	// for each commit, get the diff
	for _, commitsha := range commitList {
		curcommit, err := repo.LookupCommit(commitsha.Id())
		if err != nil {
			return errors.Convert(err)
		}
		if curcommit.ParentCount() == 0 || curcommit.ParentCount() > 0 {
			var parentTree, tree *git.Tree
			tree, err = curcommit.Tree()
			if err != nil {
				return errors.Convert(err)
			}
			var diff *git.Diff
			//FIXME error type convert
			opts, err := git.DefaultDiffOptions()
			opts.NotifyCallback = func(diffSoFar *git.Diff, delta git.DiffDelta, matchedPathSpec string) error {
				return nil
			}
			if err != nil {
				return errors.Convert(err)
			}
			if curcommit.ParentCount() > 0 {
				parent := curcommit.Parent(0)
				parentTree, err = parent.Tree()
			}
			diff, err = repo.DiffTreeToTree(parentTree, tree, &opts)
			if err != nil {
				return errors.Convert(err)
			}
			deleted := make(models.DiffLines, 0)
			added := make(models.DiffLines, 0)
			var lastFile string
			lastFile = ""
			err = diff.ForEach(func(file git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
				// if it doesn't exist in snapshot, create a new one
				if _, ok := snapshot[file.OldFile.Path]; !ok {
					fileBlame, err := models.NewFileBlame()
					if err != nil {
						r.logger.Info("Create FileBlame Error")
						return nil, err
					}
					snapshot[file.OldFile.Path] = (*models.FileBlame)(fileBlame)
				}
				if lastFile == "" {
					lastFile = file.NewFile.Path
				} else if lastFile != file.NewFile.Path {
					updateSnapshotFileBlame(curcommit, deleted, added, lastFile, snapshot)
					// reset the deleted and added,last_file now is current file
					deleted = make([]git.DiffLine, 0)
					added = make([]git.DiffLine, 0)
					lastFile = file.NewFile.Path
				}
				hunkNum := 0
				return func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
					hunkNum++
					return func(line git.DiffLine) error {
						commitLineChange := &code.CommitLineChange{}
						commitLineChange.CommitSha = curcommit.Id().String()
						commitLineChange.ChangedType = line.Origin.String()
						commitLineChange.LineNoNew = line.NewLineno
						commitLineChange.LineNoOld = line.OldLineno
						commitLineChange.OldFilePath = file.OldFile.Path
						commitLineChange.NewFilePath = file.NewFile.Path
						commitLineChange.HunkNum = hunkNum
						commitLineChange.Id = curcommit.Id().String() + ":" + file.NewFile.Path + ":" + strconv.Itoa(line.OldLineno) + ":" + strconv.Itoa(line.NewLineno)
						if line.Origin == git.DiffLineAddition {
							added = append(added, line)
						} else if line.Origin == git.DiffLineDeletion {
							fb := snapshot[file.OldFile.Path]
							l := fb.Find(line.OldLineno)
							if l != nil && l.Value != nil {
								temp := snapshot[file.OldFile.Path].Find(line.OldLineno)
								commitLineChange.PrevCommit = temp.Value.(string)
							} else {
								r.logger.Info("err", file.OldFile.Path, line.OldLineno, curcommit.Id().String())
							}
							deleted = append(deleted, line)
						}
						err = r.store.CommitLineChange(commitLineChange)
						if err != nil {
							return errors.Convert(err)
						}
						return nil
					}, nil
				}, nil
			}, git.DiffDetailLines)
			if err != nil {
				return errors.Convert(err)
			}
			//finally,process the last file in diff
			updateSnapshotFileBlame(curcommit, deleted, added, lastFile, snapshot)
		}
	}
	r.logger.Info("line change collect success")
	db := subtaskCtx.GetDal()
	err := db.Delete(&code.RepoSnapshot{}, dal.Where("repo_id= ?", r.id))
	if err != nil {
		return errors.Convert(err)
	}
	for fp := range snapshot {
		temp := snapshot[fp]
		count := 0
		for e := temp.Lines.Front(); e != nil; e = e.Next() {
			count++
			snapshotLine := &code.RepoSnapshot{}
			snapshotLine.RepoId = r.id
			snapshotLine.LineNo = count
			snapshotLine.CommitSha = e.Value.(string)
			snapshotLine.FilePath = fp
			err := r.store.RepoSnapshot(snapshotLine)
			if err != nil {
				r.logger.Info("error")
				return err
			}
		}

	}

	r.logger.Info("collect snapshot finished")
	return nil
}

func updateSnapshotFileBlame(currentCommit *git.Commit, deleted models.DiffLines, added models.DiffLines, lastFile string, snapshot map[string]*models.FileBlame) {
	sort.Sort(deleted)
	for _, line := range deleted {
		snapshot[lastFile].RemoveLine(line.OldLineno)
	}
	for _, line := range added {
		snapshot[lastFile].AddLine(line.NewLineno, currentCommit.Id().String())
	}
}

func getDiffOpts() (*git.DiffOptions, errors.Error) {
	opts, err := git.DefaultDiffOptions()
	if err != nil {
		return nil, errors.Convert(err)
	}
	opts.NotifyCallback = func(diffSoFar *git.Diff, delta git.DiffDelta, matchedPathSpec string) error {
		return nil
	}
	return &opts, nil
}
