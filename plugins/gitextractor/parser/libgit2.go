package parser

import (
	"context"
	"fmt"
	git "github.com/libgit2/git2go/v33"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/plugins/gitextractor/models"
)

const (
	BRANCH = "BRANCH"
	TAG    = "TAG"
)

type LibGit2 struct {
	store models.Store
}

func NewLibGit2(store models.Store) *LibGit2 {
	return &LibGit2{store: store}
}

func (l *LibGit2) LocalRepo(ctx context.Context, repoPath, repoId string) error {
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return err
	}
	return l.run(ctx, repo, repoId)
}

func (l *LibGit2) run(ctx context.Context, repo *git.Repository, repoId string) error {
	defer l.store.Close()
	repoInter, err := repo.NewBranchIterator(git.BranchAll)
	if err != nil {
		return err
	}
	err = repo.Tags.Foreach(func(name string, id *git.Oid) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			break
		}
		ref := &code.Ref{
			DomainEntity: domainlayer.DomainEntity{Id: fmt.Sprintf("%s:%s", repoId, name)},
			RepoId:       repoId,
			Ref:          name,
			CommitSha:    id.String(),
			RefType:      TAG,
		}
		return l.store.Refs(ref)
	})
	if err != nil {
		return err
	}
	err = repoInter.ForEach(func(branch *git.Branch, branchType git.BranchType) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			break
		}
		if branch.IsBranch() {
			name, err := branch.Name()
			if err != nil {
				return err
			}
			var sha string
			if oid := branch.Target(); oid != nil {
				sha = oid.String()
			}
			ref := &code.Ref{
				DomainEntity: domainlayer.DomainEntity{Id: fmt.Sprintf("%s:%s", repoId, name)},
				RepoId:       repoId,
				Ref:          name,
				CommitSha:    sha,
				RefType:      BRANCH,
			}
			ref.IsDefault, _ = branch.IsHead()
			return l.store.Refs(ref)
		}
		return nil
	})
	if err != nil {
		return err
	}
	odb, err := repo.Odb()
	if err != nil {
		return err
	}
	opts, err := git.DefaultDiffOptions()
	if err != nil {
		return err
	}
	opts.NotifyCallback = func(diffSoFar *git.Diff, delta git.DiffDelta, matchedPathSpec string) error {
		return nil
	}
	err = odb.ForEach(func(id *git.Oid) error {
		logger.Info("process commit:", id.String())
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			break
		}
		if id == nil {
			return nil
		}
		commit, err := repo.LookupCommit(id)
		if err != nil {
			return nil
		}
		c := &code.Commit{
			Sha:     id.String(),
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
		var commitParents []*code.CommitParent
		for i := uint(0); i < commit.ParentCount(); i++ {
			parent := commit.Parent(i)
			if parent == nil {
				continue
			}
			if parentId := parent.Id(); parentId != nil {
				commitParents = append(commitParents, &code.CommitParent{
					CommitSha:       id.String(),
					ParentCommitSha: parentId.String(),
				})
			}
			parentTree, err := parent.Tree()
			if err != nil {
				continue
			}
			tree, err := commit.Tree()
			if err != nil {
				continue
			}

			diff, err := repo.DiffTreeToTree(parentTree, tree, &opts)
			if err != nil {
				continue
			}
			//commitFile := new(code.CommitFile)
			//err = diff.ForEach(func(file git.DiffDelta, progress float64) (
			//	git.DiffForEachHunkCallback, error) {
			//	if commitFile.CommitSha != "" {
			//		err = l.store.CommitFiles(commitFile)
			//		if err != nil {
			//			logger.Error("CommitFiles error:", err)
			//		}
			//	}
			//	commitFile.CommitSha = id.String()
			//	commitFile.FilePath = file.NewFile.Path
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
			//if err != nil {
			//	return err
			//}
			//if commitFile.CommitSha != "" {
			//	err = l.store.CommitFiles(commitFile)
			//	if err != nil {
			//		logger.Error("CommitFiles error:", err)
			//	}
			//}
			stats, err := diff.Stats()
			if err != nil {
				continue
			}
			c.Additions += stats.Insertions()
			c.Deletions += stats.Deletions()
		}
		err = l.store.Commits(c)
		if err != nil {
			return err
		}
		repoCommit := &code.RepoCommit{
			RepoId:    repoId,
			CommitSha: c.Sha,
		}
		err = l.store.RepoCommits(repoCommit)
		if err != nil {
			return err
		}
		err = l.store.CommitParents(commitParents)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
