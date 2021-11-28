package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/domainlayer/models/code"
	aeModels "github.com/merico-dev/lake/plugins/merico-analysis-engine/models"
	"gorm.io/gorm/clause"
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
		err := lakeModels.Db.First(&aeCommit).Where("hex_sha = ", commit.Sha).Error
		if err != nil {
			return err
		}

		// Update the matches
		commitToUpdate := &code.Commit{
			Sha: aeCommit.HexSha,
		}

		err = lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(commitToUpdate).Error
		if err != nil {
			return err
		}

	}

	return nil
}
