package plugins

import (
	"fmt"
	"path"
	"runtime"
	"testing"
)

func TestPluginsLoading(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	pluginsDir := path.Dir(filename)
	err := LoadPlugins(pluginsDir)
	if err != nil {
		t.Errorf("Failed to LoadPlugins %v", err)
	}
	if len(Plugins) == 0 {
		t.Errorf("No plugin found")
		return
	}

	name := "jira"
	options := map[string]interface{}{
		"boardId": 20,
	}
	progress := make(chan float32)
	fmt.Printf("start runing plugin %v\n", name)
	go RunPlugin(name, options, progress)
	for p := range progress {
		fmt.Printf("running plugin %v, progress: %v\n", name, p*100)
	}
	fmt.Printf("end running plugin %v\n", name)
	if err != nil {
		t.Error(err)
	}

}
