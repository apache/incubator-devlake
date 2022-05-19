package tasks

import (
	"fmt"
	"reflect"

	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/refdiff/utils"
	"gorm.io/gorm/clause"
)

func CalculateCommitsDiff(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*RefdiffTaskData)
	repoId := data.Options.RepoId
	pairs := data.Options.Pairs
	db := taskCtx.GetDb()
	ctx := taskCtx.GetContext()
	logger := taskCtx.GetLogger()
	insertCountLimitOfRefsCommitsDiff := int(65535 / reflect.ValueOf(code.RefsCommitsDiff{}).NumField())
	// convert ref pairs into commit pairs
	ref2sha := func(refName string) (string, error) {
		ref := &code.Ref{}
		if refName == "" {
			return "", fmt.Errorf("ref name is empty")
		}
		ref.Id = fmt.Sprintf("%s:%s", repoId, refName)
		err := db.First(ref).Error
		if err != nil {
			return "", fmt.Errorf("faild to load Ref info for repoId:%s, refName:%s", repoId, refName)
		}
		return ref.CommitSha, nil
	}
	commitPairs := make([][4]string, 0, len(pairs))
	for i, refPair := range pairs {
		// get new ref's commit sha
		newCommit, err := ref2sha(refPair.NewRef)
		if err != nil {
			return fmt.Errorf("failed to load commit sha for NewRef on pair #%d: %w", i, err)
		}
		// get old ref's commit sha
		oldCommit, err := ref2sha(refPair.OldRef)
		if err != nil {
			return fmt.Errorf("failed to load commit sha for OleRef on pair #%d: %w", i, err)
		}
		commitPairs = append(commitPairs, [4]string{newCommit, oldCommit, refPair.NewRef, refPair.OldRef})
	}

	commitNodeGraph := utils.NewCommitNodeGraph()

	// load commits from db
	commitParent := &code.CommitParent{}
	cursor, err := db.Table("commit_parents cp").
		Select("cp.*").
		Joins("LEFT JOIN repo_commits rc ON (rc.commit_sha = cp.commit_sha)").
		Where("rc.repo_id = ?", repoId).
		Rows()
	if err != nil {
		panic(err)
	}
	defer cursor.Close()

	for cursor.Next() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		err = db.ScanRows(cursor, commitParent)
		if err != nil {
			return fmt.Errorf("failed to read commit from database: %v", err)
		}
		commitNodeGraph.AddParent(commitParent.CommitSha, commitParent.ParentCommitSha)
	}

	logger.Info("refdiff", fmt.Sprintf("Create a commit node graph with node count[%d]", commitNodeGraph.Size()))

	// calculate diffs for commits pairs and store them into database
	commitsDiff := &code.RefsCommitsDiff{}
	lenCommitPairs := len(commitPairs)
	taskCtx.SetProgress(0, lenCommitPairs)

	for _, pair := range commitPairs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		// ref might advance, keep commit sha for debugging
		commitsDiff.NewRefCommitSha = pair[0]
		commitsDiff.OldRefCommitSha = pair[1]
		commitsDiff.NewRefId = fmt.Sprintf("%s:%s", repoId, pair[2])
		commitsDiff.OldRefId = fmt.Sprintf("%s:%s", repoId, pair[3])

		// delete records before creation
		err = db.Delete(&code.RefsCommitsDiff{NewRefId: commitsDiff.NewRefId, OldRefId: commitsDiff.OldRefId}).Error
		if err != nil {
			return err
		}

		if commitsDiff.NewRefCommitSha == commitsDiff.OldRefCommitSha {
			// different refs might point to a same commit, it is ok
			logger.Info(
				"refdiff",
				fmt.Sprintf(
					"skipping ref pair due to they are the same %s %s => %s",
					commitsDiff.NewRefId,
					commitsDiff.OldRefId,
					commitsDiff.NewRefCommitSha,
				),
			)
			continue
		}

		lostSha, oldCount, newCount := commitNodeGraph.CalculateLostSha(pair[1], pair[0])

		commitsDiffs := []code.RefsCommitsDiff{}

		commitsDiff.SortingIndex = 1
		for _, sha := range lostSha {
			commitsDiff.CommitSha = sha
			commitsDiffs = append(commitsDiffs, *commitsDiff)

			// sql limit placeholders count only 65535
			if commitsDiff.SortingIndex%insertCountLimitOfRefsCommitsDiff == 0 {
				logger.Info("refdiff", fmt.Sprintf("commitsDiffs count in limited[%d] index[%d]--exec and clean", len(commitsDiffs), commitsDiff.SortingIndex))
				err = db.Clauses(clause.OnConflict{DoNothing: true}).Create(commitsDiffs).Error
				if err != nil {
					return err
				}
				commitsDiffs = []code.RefsCommitsDiff{}
			}

			commitsDiff.SortingIndex++
		}

		if len(commitsDiffs) > 0 {
			logger.Info("refdiff", fmt.Sprintf("insert data count [%d]", len(commitsDiffs)))
			err = db.Clauses(clause.OnConflict{DoNothing: true}).Create(commitsDiffs).Error
			if err != nil {
				return err
			}
		}

		logger.Info("refdiff", fmt.Sprintf(
			"total %d commits of difference found between %s and %s(total:%d)",
			newCount,
			commitsDiff.NewRefCommitSha,
			commitsDiff.OldRefCommitSha,
			oldCount,
		))
		taskCtx.IncProgress(1)
	}
	return nil
}

var CalculateCommitsDiffMeta = core.SubTaskMeta{
	Name:             "calculateCommitsDiff",
	EntryPoint:       CalculateCommitsDiff,
	EnabledByDefault: true,
	Description:      "Calculate diff commits between refs",
}
