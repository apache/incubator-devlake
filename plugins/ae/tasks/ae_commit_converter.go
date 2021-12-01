package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/code"
	aeModels "github.com/merico-dev/lake/plugins/ae/models"
)

// NOTE: This only works on Commits in the Domain layer. You need to run Github or Gitlab collection and Domain layer enrichemnt first.
func SetDevEqOnCommits() error {
	var commits []code.Commit

	// Get all the commits from the domain layer
	err := lakeModels.Db.Find(&commits).Error
	if err != nil {
		return err
	}

	// Loop over them
	for _, commit := range commits {

		// see if there is a match between Commit and AECommit
		var aeCommit aeModels.AECommit
		results := lakeModels.Db.Where("hex_sha = ?", commit.Sha).First(&aeCommit)

		// Check to see if a record was found
		if results.RowsAffected > 0 {
			commit.DevEq = aeCommit.DevEq
			err := lakeModels.Db.Save(&commit).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}
