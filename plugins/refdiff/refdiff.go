package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/quad"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm/clause"
)

type RefPair struct {
	NewRef string
	OldRef string
}

type RefDiffOptions struct {
	RepoId string
	Pairs  []RefPair
}

// make sure interface is implemented
var _ core.Plugin = (*RefDiff)(nil)

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry RefDiff //nolint

type RefDiff string

func (rd RefDiff) Description() string {
	return "Calculate commits diff for specified ref pairs based on `commits` and `commit_parents` tables"
}

func (rd RefDiff) Init() {
}

func (rd RefDiff) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	var op RefDiffOptions
	var err error

	// decode options
	err = mapstructure.Decode(options, &op)
	if err != nil {
		return fmt.Errorf("failed to parse option: %v", err)
	}

	// validation
	if op.RepoId == "" {
		return fmt.Errorf("repoId is required")
	}

	// convert ref pairs into commit pairs
	ref2sha := func(refName string) (string, error) {
		ref := &code.Ref{}
		if refName == "" {
			return "", fmt.Errorf("ref name is empty")
		}
		ref.Id = fmt.Sprintf("%s:%s", op.RepoId, refName)
		err = models.Db.First(ref).Error
		if err != nil {
			return "", fmt.Errorf("faild to load Ref info for %s", ref.Id)
		}
		return ref.CommitSha, nil
	}
	commitPairs := make([][4]string, 0, len(op.Pairs))
	for i, refPair := range op.Pairs {
		newCommit, err := ref2sha(refPair.NewRef)
		if err != nil {
			return fmt.Errorf("failed to load commit sha for NewRef on pair #%d: %w", i, err)
		}
		oldCommit, err := ref2sha(refPair.OldRef)
		if err != nil {
			return fmt.Errorf("failed to load commit sha for OleRef on pair #%d: %w", i, err)
		}
		commitPairs = append(commitPairs, [4]string{newCommit, oldCommit, refPair.NewRef, refPair.OldRef})
	}

	// create a in memory graph database
	store, err := cayley.NewMemoryGraph()
	if err != nil {
		return fmt.Errorf("failed to init graph store: %v", err)
	}

	// load commits from db
	commitParent := &code.CommitParent{}
	cursor, err := models.Db.Table("commit_parents cp").
		Select("cp.*").
		Joins("LEFT JOIN repo_commits rc ON (rc.commit_sha = cp.commit_sha)").
		Where("rc.repo_id = ?", op.RepoId).
		Rows()
	if err != nil {
		panic(err)
	}
	defer cursor.Close()

	for cursor.Next() {
		err = models.Db.ScanRows(cursor, commitParent)
		if err != nil {
			return fmt.Errorf("failed to read commit from database: %v", err)
		}
		err = store.AddQuad(quad.Make(commitParent.CommitSha, "childOf", commitParent.ParentCommitSha, nil))
		if err != nil {
			return fmt.Errorf("failed to add commit to graph store: %v", err)
		}
	}

	// calculate diffs for commits pairs and store them into database
	commitsDiff := &code.RefsCommitsDiff{}
	ancestors := cayley.StartMorphism().Out(quad.String("childOf"))
	for _, pair := range commitPairs {
		// ref might advance, keep commit sha for debugging
		commitsDiff.NewCommitSha = pair[0]
		commitsDiff.OldCommitSha = pair[1]
		commitsDiff.NewRefName = pair[2]
		commitsDiff.OldRefName = pair[3]

		newCommit := cayley.StartPath(store, quad.String(commitsDiff.NewCommitSha)).FollowRecursive(ancestors, -1, []string{})
		oldCommit := cayley.StartPath(store, quad.String(commitsDiff.OldCommitSha)).FollowRecursive(ancestors, -1, []string{})

		p := newCommit.Except(oldCommit)

		// delete records before creation
		err = models.Db.Exec(
			"DELETE FROM refs_commits_diffs WHERE new_ref_name = ? AND old_ref_name = ?",
			commitsDiff.NewRefName,
			commitsDiff.OldRefName,
		).Error
		if err != nil {
			return err
		}

		if commitsDiff.NewCommitSha == commitsDiff.OldCommitSha {
			// different refs might point to a same commit, it is ok
			logger.Info(
				"refdiff",
				fmt.Sprintf(
					"skipping ref pair due to they are the same %s %s => %s",
					commitsDiff.NewRefName,
					commitsDiff.OldRefName,
					commitsDiff.NewCommitSha,
				),
			)
			continue
		}

		// cayley produces a result that contains old commit sha but not new one
		// that is the opposite of what `git log oldcommit..newcommit would produces`
		// don't know  why exactly cayley does it this way, but we have to handle it anyway
		// 1. adding new commit sha
		commitsDiff.CommitSha = commitsDiff.NewCommitSha
		err = models.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(commitsDiff).Error
		if err != nil {
			return err
		}
		index := 1
		err = p.Iterate(nil).EachValue(nil, func(value quad.Value) {
			commitsDiff.CommitSha = fmt.Sprintf("%s", quad.NativeOf(value))
			// 2. ignoring old commit sha
			if commitsDiff.CommitSha == commitsDiff.OldCommitSha {
				return
			}
			commitsDiff.SortingIndex = index
			err = models.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(commitsDiff).Error
			if err != nil {
				panic(err)
			}
			index++
		})
		if err != nil {
			return err
		}
		logger.Info("refdiff", fmt.Sprintf("total %d commits of difference found between %s and %s", index, commitsDiff.NewCommitSha, commitsDiff.OldCommitSha))
	}

	return nil
}

// PkgPath information lost when compiled as plugin(.so)
func (rd RefDiff) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/jira"
}

func (rd RefDiff) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return nil
}

// standalone mode for debugging
func main() {
	var err error

	args := os.Args[1:]
	if len(args) < 2 {
		panic(fmt.Errorf("Usage: refdiff <repo_id> <new_ref_name> <old_ref_name>"))
	}
	repoId, newRefName, oldRefName := args[0], args[1], args[2]

	err = core.RegisterPlugin("refdiff", PluginEntry)
	if err != nil {
		panic(err)
	}
	PluginEntry.Init()
	progress := make(chan float32)
	go func() {
		err := PluginEntry.Execute(
			map[string]interface{}{
				"repoId": repoId,
				"pairs": []map[string]string{
					{
						"NewRef": newRefName,
						"OldRef": oldRefName,
					},
				},
			},
			progress,
			context.Background(),
		)
		if err != nil {
			panic(err)
		}
		close(progress)
	}()
	for p := range progress {
		fmt.Println(p)
	}
}
