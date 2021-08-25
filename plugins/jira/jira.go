package main // must be main for plugin entry point

import (
	"fmt"
	"time"
)

// A pseudo type for Plugin Interface implementation
type Jira string

func (jira Jira) Description() string {
	return "To collect and enrich data from JIRA"
}

func (jira Jira) Execute(options map[string]interface{}, progress chan<- float32) {
	fmt.Println("start jira plugin execution")
	time.Sleep(1 * time.Second)
	progress <- 0.1
	time.Sleep(1 * time.Second)
	progress <- 0.5
	time.Sleep(1 * time.Second)
	progress <- 1
	fmt.Println("end jira plugin execution")
	close(progress)
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry Jira
