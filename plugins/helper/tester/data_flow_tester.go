package tester

import (
	"testing"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/runner"
	"gorm.io/gorm"
)

// DataFlowTester helps you write subtasks e2e tests with ease
type DataFlowTester struct {
	Db     *gorm.DB
	T      *testing.T
	Plugin core.PluginMeta
}

// NewDataFlowTester create a *DataFlowTester to help developer test their subtasks data flow
func NewDataFlowTester(t *testing.T, plugin core.PluginMeta) *DataFlowTester {
	cfg := config.GetConfig()
	db, err := runner.NewGormDb(cfg, logger.Global)
	if err != nil {
		panic(err)
	}
	return &DataFlowTester{
		Db:     db,
		T:      t,
		Plugin: plugin,
	}
}

func (t *DataFlowTester) ImportCsv(csvRelPath string, tableName string) {
}

func (t *DataFlowTester) Subtask(subtaskMeta core.SubTaskMeta) {
}

func (T *DataFlowTester) VerifyTable(tableName string, csvRelPaht string) {
}
