package main

import (
	"testing"
)

func TestExecut(t *testing.T) {
	opts := make(map[string]interface{})
	progress := make(chan float32)

	PluginEntry.Init()
	go PluginEntry.Execute(opts, progress)
	<-progress
}
