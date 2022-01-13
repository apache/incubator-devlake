package main

import (
	"fmt"

	git2go "github.com/libgit2/git2go/v28"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"gorm.io/gorm/clause"
)

type CommitIterator func(commit *git2go.Commit, insertions, deletions int, parents []string) error

func IterateOdbCommits(repoPath string, iterator CommitIterator) error {
	fmt.Println("GetCommitsByGit2Go")
	repo, err := git2go.OpenRepository(repoPath)
	if err != nil {
		return err
	}
	odb, err := repo.Odb()
	if err != nil {
		return err
	}
	err = odb.ForEach(func(id *git2go.Oid) error {
		if id == nil {
			return nil
		}
		commit, err := repo.LookupCommit(id)
		if err != nil {
			// not a commit
			return nil
		}
		insertions := 0
		deletions := 0
		parents := make([]string, commit.ParentCount())
		for i := uint(0); i < commit.ParentCount(); i++ {
			parent := commit.Parent(i)
			parentTree, err := parent.Tree()
			if err != nil {
				return err
			}
			tree, err := commit.Tree()
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
		return iterator(commit, insertions, deletions, parents)
	})
	if err != nil {
		return err
	}
	return nil
}

func CollectCommits(repoId string, repoPath string) error {
	commit := &code.Commit{}
	repoCommit := &code.RepoCommit{
		RepoId: repoId,
	}
	commitParent := &code.CommitParent{}
	total := 0
	err := IterateOdbCommits(repoPath, func(c *git2go.Commit, insertions, deletions int, parents []string) error {
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
		err := models.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(commit).Error
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
			err = models.Db.Create(commitParent).Error
			if err != nil {
				return err
			}
		}
		total++
		fmt.Printf("\r%v  %d", commit.Sha, total)
		return nil
	})
	return err
}

func GetCommitsByRevWalk(repoPath string) {
	fmt.Println("GetCommitsByGit2Go")
	repo, err := git2go.OpenRepository(repoPath)
	if err != nil {
		panic(err)
	}
	walk, err := repo.Walk()
	if err != nil {
		panic(err)
	}

	err = walk.Iterate(func(commit *git2go.Commit) bool {
		fmt.Printf("hello")
		insertions := 0
		deletions := 0
		parents := ""
		for i := uint(0); i < commit.ParentCount(); i++ {
			parent := commit.Parent(i)
			parentTree, err := parent.Tree()
			if err != nil {
				panic(err)
			}
			tree, err := commit.Tree()
			if err != nil {
				panic(err)
			}
			diff, err := repo.DiffTreeToTree(parentTree, tree, &git2go.DiffOptions{
				Flags: git2go.DiffIgnoreSubmodules | git2go.DiffIgnoreFilemode,
			})
			if err != nil {
				panic(err)
			}
			stats, err := diff.Stats()
			if err != nil {
				panic(err)
			}
			insertions += stats.Insertions()
			deletions += stats.Deletions()

			if parents != "" {
				parents += ", "
			}
			parents += fmt.Sprintf("%v", parent.Id())
		}
		fmt.Printf("%v %v   added: %v  deleted: %v parents: %v\n", commit.Id(), commit.Author().Email, insertions, deletions, parents)
		return true
	})
	walk.Free()
	if err != nil {
		panic(err)
	}
	repo.Free()
}

func GetCommitsByOdbIteration(repoPath string) {
	IterateOdbCommits(repoPath, func(commit *git2go.Commit, insertions, deletions int, parents []string) error {
		fmt.Printf("%v %v   added: %v  deleted: %v parents: %v\n", commit.Id(), commit.Author().Email, insertions, deletions, parents)
		return nil
	})
}

func main() {
	var test *code.Commit
	println(test.AuthorId)
	err := CollectCommits("github:GithubRepository:384111310", "../")
	if err != nil {
		panic(err)
	}
	fmt.Println("done")
}
