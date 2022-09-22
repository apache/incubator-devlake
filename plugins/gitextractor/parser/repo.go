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
	"container/list"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/gitextractor/models"
	git "github.com/libgit2/git2go/v33"
	"log"
	"regexp"
	"sort"
)

type GitRepo struct {
	store   models.Store
	logger  core.Logger
	id      string
	repo    *git.Repository
	cleanup func()
}

func (r *GitRepo) CollectAll(subtaskCtx core.SubTaskContext) errors.Error {
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
	return r.CollectDiffLine(subtaskCtx)
	return nil
}

func (r *GitRepo) Close() errors.Error {
	defer func() {
		if r.cleanup != nil {
			r.cleanup()
		}
	}()
	return r.store.Close()
}

func (r *GitRepo) CountTags() (int, errors.Error) {
	tags, err := r.repo.Tags.List()
	if err != nil {
		return 0, errors.Convert(err)
	}
	return len(tags), nil
}

func (r *GitRepo) CountBranches(ctx context.Context) (int, errors.Error) {
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
			count++
		}
		return nil
	})
	return count, errors.Convert(err)
}

func (r *GitRepo) CountCommits(ctx context.Context) (int, errors.Error) {
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
		commit, _ := r.repo.LookupCommit(id)
		if commit != nil {
			count++
		}
		return nil
	})
	return count, errors.Convert(err)
}

func (r *GitRepo) CollectTags(subtaskCtx core.SubTaskContext) errors.Error {
	return errors.Convert(r.repo.Tags.Foreach(func(name string, id *git.Oid) error {
		select {
		case <-subtaskCtx.GetContext().Done():
			return subtaskCtx.GetContext().Err()
		default:
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
		r.logger.Info("tagCommit:%s", tagCommit)
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
	}))
}

func (r *GitRepo) CollectBranches(subtaskCtx core.SubTaskContext) errors.Error {
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
	}))
}

func (r *GitRepo) CollectCommits(subtaskCtx core.SubTaskContext) errors.Error {
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
		err = r.storeParentCommits(commitSha, commit)
		if err != nil {
			return err
		}
		if commit.ParentCount() > 0 {
			parent := commit.Parent(0)
			if parent != nil {
				var stats *git.DiffStats
				if stats, err = r.getDiffComparedToParent(c.Sha, commit, parent, opts, componentMap); err != nil {
					return err
				}
				c.Additions += stats.Insertions()
				c.Deletions += stats.Deletions()
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
	}))
}

func (r *GitRepo) storeParentCommits(commitSha string, commit *git.Commit) errors.Error {
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

func (r *GitRepo) getDiffComparedToParent(commitSha string, commit *git.Commit, parent *git.Commit, opts *git.DiffOptions, componentMap map[string]*regexp.Regexp) (*git.DiffStats, errors.Error) {
	var err error
	var parentTree, tree *git.Tree
	parentTree, err = parent.Tree()
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
	err = r.storeCommitFilesFromDiff(commitSha, diff, componentMap)
	if err != nil {
		return nil, errors.Convert(err)
	}
	var stats *git.DiffStats
	stats, err = diff.Stats()
	if err != nil {
		return nil, errors.Convert(err)
	}
	return stats, nil
}

func (r *GitRepo) storeCommitFilesFromDiff(commitSha string, diff *git.Diff, componentMap map[string]*regexp.Regexp) errors.Error {
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
		//commitFileComponent.FilePath = file.NewFile.Path
		//commitFileComponent.CommitSha = commitSha
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

func (r *GitRepo) CollectDiffLine(subtaskCtx core.SubTaskContext) errors.Error {
	//Using this subtask,we can get every line change in every commit.
	//We maintain a snapshot structure to get which commit each deleted line belongs to
	snapshot := make(map[string] /*file path*/ *FileBlame)
	repo := r.repo
	//step 1. get the reverse commit list
	commitList := make([]git.Commit, 0)
	//get currently head commitsha, dafault is master branch
	// check branch, if not master, checkout to branch's head
	commitOid, _ := repo.Head()
	//get head commit object and add into commitList
	commit, _ := repo.LookupCommit(commitOid.Target())
	commitList = append(commitList, *commit)
	// if current head has parent, get parent commitsha
	for commit != nil && commit.ParentCount() > 0 {
		pid := commit.ParentId(0)
		commit, _ = repo.LookupCommit(pid)
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
		//process the first commit
		if curcommit.ParentCount() == 0 {
			var parentTree, tree *git.Tree
			tree, err = curcommit.Tree()
			if err != nil {
				return errors.Convert(err)
			}
			var diff *git.Diff
			opts, err := git.DefaultDiffOptions()
			if err != nil {
				return errors.Convert(err)
			}
			opts.NotifyCallback = func(diffSoFar *git.Diff, delta git.DiffDelta, matchedPathSpec string) error {
				return nil
			}

			diff, err = repo.DiffTreeToTree(parentTree, tree, &opts)
			if err != nil {
				return errors.Convert(err)
			}
			hunks := make([]git.DiffHunk, 0)
			deleted := make(difflines, 0)
			added := make([]git.DiffLine, 0)
			var lastFile string
			lastFile = ""
			err = diff.ForEach(func(file git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
				//if file doesn't exist in snapshot, create a new one
				if _, ok := snapshot[file.NewFile.Path]; !ok {
					nfb, err := NewFileBlame()
					if err != nil {
						log.Fatal(err)
					}
					snapshot[file.NewFile.Path] = (*FileBlame)(nfb)
				}
				if lastFile == "" {
					lastFile = file.NewFile.Path
				} else if lastFile != file.NewFile.Path {
					UpdateSnapshotFileBlame(curcommit, deleted, added, lastFile, snapshot)
					//reset the deleted and added,last_file now is current file
					deleted = make([]git.DiffLine, 0)
					added = make([]git.DiffLine, 0)
					hunks = make([]git.DiffHunk, 0)
					lastFile = file.NewFile.Path
				}
				hunkNum := 0
				return func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
					hunks = append(hunks, hunk)
					hunkNum++
					return func(line git.DiffLine) error {
						//update snapshot
						commitLineChange := &code.CommitLineChange{}
						commitLineChange.CommitSha = curcommit.Id().String()
						commitLineChange.ChangedType = line.Origin.String()
						commitLineChange.LineNoNew = int16(line.NewLineno)
						commitLineChange.LineNoOld = int16(line.OldLineno)
						commitLineChange.OldFilePath = file.OldFile.Path
						commitLineChange.NewFilePath = file.NewFile.Path
						commitLineChange.HunkNum = string(hunkNum)
						r.store.CommitLineChange(commitLineChange)

						if line.Origin == git.DiffLineAddition {
							added = append(added, line)
						} else if line.Origin == git.DiffLineDeletion {
							deleted = append(deleted, line)
						}
						return nil
					}, nil
				}, nil
			}, git.DiffDetailLines)
			if err != nil {
				return errors.Convert(err)
			}
			//process the last file in the diff

			UpdateSnapshotFileBlame(curcommit, deleted, added, lastFile, snapshot)
		} else if curcommit.ParentCount() > 0 {
			parent := curcommit.Parent(0)
			if parent != nil {
				var parentTree, tree *git.Tree
				var err error
				parentTree, err = parent.Tree()
				if err != nil {
					return errors.Convert(err)
				}
				tree, err = curcommit.Tree()
				if err != nil {
					return errors.Convert(err)
				}
				var diff *git.Diff
				opts := &git.DiffOptions{}
				//if we use these options, the diff will be very slow and snapshot will be wrong
				// findOpts, err := git.DefaultDiffFindOptions()
				// findOpts.Flags = git.DiffFindByConfig //git.DiffFindBreakRewrites
				// err = diff.FindSimilar(&findOpts)
				diff, err = repo.DiffTreeToTree(parentTree, tree, opts)
				if err != nil {
					return errors.Convert(err)
				}
				hunks := make([]git.DiffHunk, 0)
				deleted := make(difflines, 0)
				added := make([]git.DiffLine, 0)
				var lastFile string
				lastFile = ""

				// file callback
				err = diff.ForEach(func(file git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
					//if doesn't exist in snapshot, create a new one
					if _, ok := snapshot[file.OldFile.Path]; !ok {
						nfb, err := NewFileBlame()
						if err != nil {
							log.Fatal(err)
						}
						snapshot[file.OldFile.Path] = (*FileBlame)(nfb)
					}
					// if file renamed or copied, add the new file into snapshot
					// if we don't use find_similar opts,each file treat as NewFile
					//if file.NewFile.Path != file.OldFile.Path {
					//	snapshot[file.NewFile.Path] = snapshot[file.OldFile.Path] //*/ nfb
					//	return func(dh git.DiffHunk) (git.DiffForEachLineCallback, error) {
					//		return func(dl git.DiffLine) error {
					//			return nil
					//		}, nil
					//	}, nil
					//
					//}
					if lastFile == "" {
						lastFile = file.NewFile.Path
					} else if lastFile != file.NewFile.Path {
						UpdateSnapshotFileBlame(curcommit, deleted, added, lastFile, snapshot)
						//reset the deleted and added,last_file now is current file
						deleted = make([]git.DiffLine, 0)
						added = make([]git.DiffLine, 0)
						hunks = make([]git.DiffHunk, 0)
						lastFile = file.NewFile.Path
					}
					// hunk callback
					hunkNum := 0
					return func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
						hunks = append(hunks, hunk)
						hunkNum++
						//line call back
						return func(line git.DiffLine) error {
							//first store line message
							commitLineChange := &code.CommitLineChange{}
							commitLineChange.CommitSha = curcommit.Id().String()
							commitLineChange.ChangedType = line.Origin.String()
							commitLineChange.LineNoNew = int16(line.NewLineno)
							commitLineChange.LineNoOld = int16(line.OldLineno)
							commitLineChange.OldFilePath = file.OldFile.Path
							commitLineChange.NewFilePath = file.NewFile.Path
							commitLineChange.HunkNum = string(hunkNum)

							if line.Origin == git.DiffLineAddition {
								added = append(added, line)
							} else if line.Origin == git.DiffLineDeletion {
								t := snapshot[file.OldFile.Path]
								a := t.Find(line.OldLineno)
								if a != nil && a.Value != nil {
									temp := snapshot[file.OldFile.Path].Find(line.OldLineno)
									commitLineChange.PrevCommit = temp.Value.(string)

								} else {
									fmt.Println("err", file.OldFile.Path, line.OldLineno, curcommit.Id().String())
								}
								deleted = append(deleted, line)
							}
							r.store.CommitLineChange(commitLineChange)

							return nil
						}, nil
					}, nil
				}, git.DiffDetailLines)
				if err != nil {
					log.Fatal(err)
				}
				//finally,process the last file in  diff
				UpdateSnapshotFileBlame(curcommit, deleted, added, lastFile, snapshot)
			}
		}
	}

	//for test
	log.Print("Test Area")
	printSnapshotFileBlame("README.md", snapshot)

	return nil
}

func NewFileBlame() (*models.FileBlame, error) {
	fb := &models.FileBlame{Idx: 0, It: &list.Element{}, Lines: list.New()}
	fb.It = fb.Lines.Front()
	fb.Idx = 0
	return fb, nil
}

type difflines []git.DiffLine

// some essential functions for custom sort
func (difflines difflines) Len() int {
	return len(difflines)
}

func (difflines difflines) Less(i, j int) bool {
	return difflines[i].OldLineno > difflines[j].OldLineno
}

func (difflines difflines) Swap(i, j int) {
	temp := difflines[i]
	difflines[i] = difflines[j]
	difflines[j] = temp
}

func UpdateSnapshotFileBlame(currentCommit *git.Commit, deleted difflines, added difflines, lastFile string, snapshot map[string]*FileBlame) {
	sort.Sort(deleted)

	for _, line := range deleted {
		snapshot[lastFile].RemoveLine(line.OldLineno)
	}
	for _, line := range added {
		snapshot[lastFile].AddLine(line.NewLineno, currentCommit.Id().String())
	}
}

type FileBlame struct {
	Idx   int
	It    *list.Element
	Lines *list.List
}

func (fb *FileBlame) Walk(num int) {
	for fb.Idx < num && fb.It != fb.Lines.Back() {
		fb.Idx++
		fb.It = fb.It.Next()
	}
	for fb.Idx > num && fb.It != fb.Lines.Front() {
		fb.Idx--
		fb.It = fb.It.Prev()
	}
}

func (fb *FileBlame) Find(num int) *list.Element {
	fb.Walk(num)
	if fb.Idx == num && fb.It != nil {
		return fb.It
	}
	return nil
}

func (fb *FileBlame) AddLine(num int, commit string) {
	fb.Walk(num)
	flag := false
	for fb.It == fb.Lines.Back() && fb.Idx < num {
		flag = true
		fb.It = fb.Lines.PushBack(nil)
		fb.Idx++

	}

	if fb.It == nil {
		fb.It = fb.Lines.PushBack(commit)
	} else if flag {
		fb.It.Value = commit
	} else {
		fb.It = fb.Lines.InsertBefore(commit, fb.It)
	}
}

func (fb *FileBlame) RemoveLine(num int) {

	fb.Walk(num)
	a := fb.It
	if fb.Idx < 0 || num < 1 {
		log.Fatal("RemoveLine Error")
		return
	}
	if fb.Idx == num && fb.It != nil {
		if fb.Lines.Len() == 1 {
			fb.Idx = 0
			fb.Lines.Init()
			fb.It = fb.Lines.Front()
			return
		}
		if fb.Idx == 1 {
			fb.It = fb.It.Next()
			fb.Lines.Remove(fb.It.Prev())
			return
		}
		if fb.It != fb.Lines.Back() {
			fb.It = fb.It.Next()
		} else {
			fb.It = fb.It.Prev()
			fb.Idx--
		}

		fb.Lines.Remove(a)

	}
}
func printSnapshotFileBlame(filepath string, snapshot map[string]*FileBlame) {
	if _, ok := snapshot[filepath]; !ok {
		log.Println("file not found")
		return
	}
	temp := snapshot[filepath]
	cnt := 0
	for e := temp.Lines.Front(); e != nil; e = e.Next() {
		cnt++
		log.Println(cnt, e.Value)
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
