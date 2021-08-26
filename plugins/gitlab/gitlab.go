package main // must be main for plugin entry point

// A pseudo type for Plugin Interface implementation
type Gitlab string

func (plugin Gitlab) Description() string {
	return "To collect and enrich data from Gitlab"
}

func (plugin Gitlab) Execute(options map[string]interface{}, progress chan<- float32) {
	progress <- 1
	close(progress)
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry Gitlab //nolint
