package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/bitbucket-server/models"
)

type addInitTables20231123 struct{}

func (script *addInitTables20231123) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&models.BitbucketServerAccount{},
		&models.BitbucketServerCommit{},
		&models.BitbucketServerConnection{},
		&models.BitbucketServerPullRequest{},
		&models.BitbucketServerPrComment{},
		&models.BitbucketServerPrCommit{},
		&models.BitbucketServerRepo{},
		&models.BitbucketServerRepoCommit{},
		&models.BitbucketServerScopeConfig{},
	)
}

func (*addInitTables20231123) Version() uint64 {
	return 20231123112623
}

func (*addInitTables20231123) Name() string {
	return "Bitbucket Server init schema 20231123"
}
