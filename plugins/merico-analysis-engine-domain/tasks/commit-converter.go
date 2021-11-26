package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	domainlayerBase "github.com/merico-dev/lake/plugins/domainlayer/models/base"
	"github.com/merico-dev/lake/plugins/domainlayer/models/code"
	"github.com/merico-dev/lake/plugins/domainlayer/okgen"
	"github.com/merico-dev/lake/plugins/merico-analysis-engine/models"
	"gorm.io/gorm/clause"
)

func ConvertCommits() error {
	aeCommit := &models.AECommit{}

	cursor, err := lakeModels.Db.Model(aeCommit).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	// iterate all rows
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, aeCommit)
		if err != nil {
			return err
		}

		commitOriginkeyGenerator := okgen.NewOriginKeyGenerator(aeCommit)

		// TODO: find if there is a commit existing that we can enhance with dev eq
		// TODO: update with dev eq and save

		commit := &code.Commit{
			DomainEntity: domainlayerBase.DomainEntity{
				OriginKey: commitOriginkeyGenerator.Generate(aeCommit.HexSha),
			},
		}

		err = lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(commit).Error
		if err != nil {
			return err
		}
	}
	return nil
}
