package tasks

import (
	"context"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/code"
	aeModels "github.com/merico-dev/lake/plugins/ae/models"
	"github.com/merico-dev/lake/plugins/core"
)

// NOTE: This only works on Commits in the Domain layer. You need to run Github or Gitlab collection and Domain layer enrichemnt first.
func SetDevEqOnCommits(ctx context.Context) error {
	commit := &code.Commit{}
	aeCommit := &aeModels.AECommit{}

	// Get all the commits from the domain layer
	cursor, err := lakeModels.Db.Model(aeCommit).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	// Loop over them
	for cursor.Next() {
		// we do a non-blocking checking, this should be applied to all converters from all plugins
		select {
		case <-ctx.Done():
			return core.TaskCanceled
		default:
		}
		// uncomment following line if you want to test out canceling feature for this task
		//time.Sleep(1 * time.Second)

		err = lakeModels.Db.ScanRows(cursor, aeCommit)
		if err != nil {
			return err
		}

		err := lakeModels.Db.Model(commit).Where("sha = ?", aeCommit.HexSha).Update("dev_eq", aeCommit.DevEq).Error
		if err != nil {
			return err
		}
	}

	return nil
}
