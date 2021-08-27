package main // must be main for plugin entry point

import (
	"fmt"
	"time"

	"github.com/merico-dev/lake/plugins/jira/tasks"
)

// A pseudo type for Plugin Interface implementation
type Jira string

func (plugin Jira) Description() string {
	return "To collect and enrich data from JIRA"
}

func (plugin Jira) Execute(options map[string]interface{}, progress chan<- float32) {
	boardId, ok := options["boardId"]
	if !ok {
		fmt.Println("boardId is required for jira execution")
		return
	}
	boardIdInt := int(boardId.(float64))
	if boardIdInt < 0 {
		fmt.Println("boardId is invalid")
		return
	}
	fmt.Println("start jira plugin execution")
	err := tasks.CollectBoard(boardIdInt)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
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
var PluginEntry Jira //nolint
