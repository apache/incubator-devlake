package plugins

import (
	"path"
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins"
	"github.com/merico-dev/lake/plugins/core"
)

func TestPluginsLoading(t *testing.T) {
	err := plugins.LoadPlugins(plugins.PluginDir())
	if err != nil {
		t.Errorf("Failed to LoadPlugins %v", err)
	}
	if len(core.AllPlugins()) == 0 {
		t.Errorf("No plugin found")
	}

	// name := "jira"
	// options := map[string]interface{}{
	// 	"boardId": 20,
	// }
	// progress := make(chan float32)
	// fmt.Printf("start runing plugin %v\n", name)
	// go func() {
	// 	_ = plugins.RunPlugin(name, options, progress)
	// }()
	// for p := range progress {
	// 	fmt.Printf("running plugin %v, progress: %v\n", name, p*100)
	// }
	// fmt.Printf("end running plugin %v\n", name)
	// if err != nil {
	// 	t.Error(err)
	// }

}

func TestPluginDir(t *testing.T) {
	pluginDir := plugins.PluginDir()
	assert.Equal(t, path.IsAbs(pluginDir), true)
	assert.Equal(t, strings.HasSuffix(pluginDir, config.V.GetString("PLUGIN_DIR")), true)
}
