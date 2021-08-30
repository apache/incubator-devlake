package main

import (
	"fmt"
	"testing"
)

func TestGitlabPlugin(t *testing.T) {
	PluginEntry.Init()

	projectId := 20103385

	options := make(map[string]interface{})
	options["projectId"] = projectId

	c := make(chan float32)

	go func() {
		PluginEntry.Execute(options, c)
	}()

	for p := range c {
		fmt.Printf("running plugin Jira, progress: %v\n", p*100)
	}

}
