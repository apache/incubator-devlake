package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.MigrationScript = (*addAdditionsDeletionsToMr)(nil)

type mrAdditionsDeletions struct {
	Additions int
	Deletions int
}

func (mrAdditionsDeletions) TableName() string {
	return "_tool_gitlab_merge_requests"
}

type addAdditionsDeletionsToMr struct{}

func (*addAdditionsDeletionsToMr) Up(basicRes context.BasicRes) errors.Error {
	return errors.Convert(basicRes.GetDal().AutoMigrate(&mrAdditionsDeletions{}))
}

func (*addAdditionsDeletionsToMr) Version() uint64 {
	return 20260615000001
}

func (*addAdditionsDeletionsToMr) Name() string {
	return "gitlab: add additions and deletions to merge requests"
}