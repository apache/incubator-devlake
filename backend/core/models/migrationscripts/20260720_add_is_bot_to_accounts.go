package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.MigrationScript = (*addIsBotToAccounts)(nil)

type account20260720 struct {
	IsBot bool `gorm:"default:false"`
}

func (account20260720) TableName() string {
	return "accounts"
}

type addIsBotToAccounts struct{}

func (*addIsBotToAccounts) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if err := db.AutoMigrate(&account20260720{}); err != nil {
		return err
	}
	return nil
}

func (*addIsBotToAccounts) Version() uint64 {
	return 20260720120000
}

func (*addIsBotToAccounts) Name() string {
	return "add is_bot to accounts so bot/automation activity can be excluded from metrics, according to #8974"
}
