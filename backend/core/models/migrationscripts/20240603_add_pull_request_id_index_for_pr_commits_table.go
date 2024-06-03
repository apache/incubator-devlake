package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
)

type addPullRequestIdIndexToPullRequestCommits struct{}

func (u *addPullRequestIdIndexToPullRequestCommits) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	err := db.Exec("ALTER TABLE pull_request_commits DROP PRIMARY KEY;")
	if err != nil {
		return err
	}
	err = db.Exec("ALTER TABLE pull_request_commits ADD PRIMARY KEY (pull_request_id, commit_sha);")
	if err != nil {
		return err
	}

	return nil
}

func (*addPullRequestIdIndexToPullRequestCommits) Version() uint64 {
	return 20240602103400
}

func (*addPullRequestIdIndexToPullRequestCommits) Name() string {
	return "changing pull_request_commits primary key columns order"
}
