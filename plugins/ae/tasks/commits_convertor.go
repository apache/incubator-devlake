package tasks

import (
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	aeModels "github.com/apache/incubator-devlake/plugins/ae/models"
	"github.com/apache/incubator-devlake/plugins/core"
)

// NOTE: This only works on Commits in the Domain layer. You need to run Github or Gitlab collection and Domain layer enrichemnt first.
func ConvertCommits(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*AeTaskData)

	commit := &code.Commit{}
	aeCommit := &aeModels.AECommit{}

	// Get all the commits from the domain layer
	cursor, err := db.Model(aeCommit).Where("ae_project_id = ?", data.Options.ProjectId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	taskCtx.SetProgress(0, -1)
	ctx := taskCtx.GetContext()
	// Loop over them
	for cursor.Next() {
		// we do a non-blocking checking, this should be applied to all converters from all plugins
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		// uncomment following line if you want to test out canceling feature for this task
		//time.Sleep(1 * time.Second)

		err = db.ScanRows(cursor, aeCommit)
		if err != nil {
			return err
		}

		err := db.Model(commit).Where("sha = ?", aeCommit.HexSha).Update("dev_eq", aeCommit.DevEq).Error
		if err != nil {
			return err
		}
		taskCtx.IncProgress(1)
	}

	return nil
}

var ConvertCommitsMeta = core.SubTaskMeta{
	Name:             "convertCommits",
	EntryPoint:       ConvertCommits,
	EnabledByDefault: true,
	Description:      "Update domain layer commits dev_eq field according to ae_commits",
}
