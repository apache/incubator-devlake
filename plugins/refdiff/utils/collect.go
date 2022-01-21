package main

// script to collect local git repo data into db for development/testing

import (
	"fmt"

	git2go "github.com/libgit2/git2go/v33"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"gorm.io/gorm/clause"
)

func collect(repoId string, repoPath string) error {

	repo, err := git2go.OpenRepository(repoPath)
	if err != nil {
		return err
	}

	ref := &code.Ref{}

	ref.RefType = "BRANCH"
	branchIter, err := repo.NewBranchIterator(git2go.BranchAll)
	if err != nil {
		return err
	}
	err = branchIter.ForEach(func(b *git2go.Branch, bt git2go.BranchType) error {
		name, err := b.Name()
		if err != nil {
			return err
		}
		ref.Id = fmt.Sprintf("%s:%s", repoId, name)
		if b.Target() != nil {
			ref.CommitSha = b.Target().String()
		}
		ref.IsDefault, _ = b.IsHead()
		fmt.Printf("branch %s : %s\n", name, ref.CommitSha)
		return models.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(ref).Error
	})
	branchIter.Free()
	if err != nil {
		return err
	}

	ref.RefType = "TAG"
	repo.Tags.Foreach(func(name string, id *git2go.Oid) error {
		ref.Id = fmt.Sprintf("%s:%s", repoId, name)
		ref.CommitSha = id.String()
		if err != nil {
			return err
		}
		fmt.Printf("tag: %s : %s\n", name, ref.CommitSha)
		return models.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(ref).Error
	})

	odb, err := repo.Odb()
	if err != nil {
		return err
	}

	commit := &code.Commit{}
	repoCommit := &code.RepoCommit{
		RepoId: repoId,
	}
	commitParent := &code.CommitParent{}
	total := 0

	err = odb.ForEach(func(id *git2go.Oid) error {
		if id == nil {
			return nil
		}
		c, err := repo.LookupCommit(id)
		if err != nil {
			// not a commit
			return nil
		}
		insertions := 0
		deletions := 0
		parents := make([]string, c.ParentCount())
		for i := uint(0); i < c.ParentCount(); i++ {
			parent := c.Parent(i)
			parentTree, err := parent.Tree()
			if err != nil {
				return err
			}
			tree, err := c.Tree()
			if err != nil {
				return err
			}
			diff, err := repo.DiffTreeToTree(parentTree, tree, &git2go.DiffOptions{
				Flags: git2go.DiffIgnoreSubmodules | git2go.DiffIgnoreFilemode,
			})
			if err != nil {
				return err
			}
			stats, err := diff.Stats()
			if err != nil {
				return err
			}
			insertions += stats.Insertions()
			deletions += stats.Deletions()

			parents[i] = fmt.Sprintf("%v", parent.Id())
		}
		commit.Sha = c.Id().String()
		commit.Additions = insertions
		commit.Deletions = deletions
		commit.Message = c.Message()
		commit.AuthorName = c.Author().Name
		commit.AuthorEmail = c.Author().Email
		commit.AuthoredDate = c.Author().When
		commit.CommitterName = c.Committer().Name
		commit.CommitterEmail = c.Committer().Email
		commit.CommittedDate = c.Committer().When
		err = models.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(commit).Error
		if err != nil {
			return err
		}
		repoCommit.CommitSha = commit.Sha
		err = models.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(repoCommit).Error
		if err != nil {
			return err
		}

		err = models.Db.Where("commit_sha = ?", commit.Sha).Delete(commitParent).Error
		if err != nil {
			return err
		}
		commitParent.CommitSha = commit.Sha
		for _, parentSha := range parents {
			commitParent.ParentCommitSha = parentSha
			err = models.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(commitParent).Error
			if err != nil {
				return err
			}
		}
		total++
		fmt.Printf("\r%v  %d", commit.Sha, total)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// tidb
	//err := collect("github:GithubRepository:41986369", "/home/klesh/Projects/merico/tidb")
	// devlake
	err := collect("github:GithubRepository:384111310", ".")
	if err != nil {
		panic(err)
	}
	fmt.Println("\ndone")
}
