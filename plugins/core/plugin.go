package core

import (
	"context"
)

// Deprecated: this all in one interface is bloated, implement PluginMeta instead
// and then, based on what a plugin can offer,  implement PluginInit/PluginTask/PluginApi on demand
type Plugin interface {
	PluginMeta
	PluginApi
	Init()
	Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error
}
