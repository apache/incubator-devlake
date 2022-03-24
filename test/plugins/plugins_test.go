package plugins

import (
	"testing"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/runner"
	"github.com/stretchr/testify/assert"
)

func TestPluginsLoading(t *testing.T) {
	cfg := config.GetConfig()
	log := logger.Global
	db, err := runner.NewGormDb(cfg, log)
	if !assert.Nil(t, err) {
		return
	}
	err = runner.LoadPlugins(cfg.GetString("PLUGIN_DIR"), cfg, log, db)
	if !assert.Nil(t, err) {
		return
	}
	assert.NotEmpty(t, core.AllPlugins())
}
