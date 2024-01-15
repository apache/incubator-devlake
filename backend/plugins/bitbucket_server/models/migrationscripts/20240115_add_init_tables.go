package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/bitbucket_server/models"
)

type addInitTables20240115 struct{}

func (script *addInitTables20240115) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&models.BitbucketServerUser{},
		&models.BitbucketServerCommit{},
		&models.BitbucketServerBranch{},
		&models.BitbucketServerConnection{},
		&models.BitbucketServerPullRequest{},
		&models.BitbucketServerPrComment{},
		&models.BitbucketServerPrCommit{},
		&models.BitbucketServerRepo{},
		&models.BitbucketServerRepoCommit{},
		&models.BitbucketServerScopeConfig{},
	)
}

func (*addInitTables20240115) Version() uint64 {
	return 20240115
}

func (*addInitTables20240115) Name() string {
	return "Bitbucket Server init schema 20240115"
}
