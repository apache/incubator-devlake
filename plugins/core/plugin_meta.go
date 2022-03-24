package core

// Minimal features a plugin should comply, should be implemented by all plugins
type PluginMeta interface {
	Description() string
	// PkgPath information lost when compiled as plugin(.so)
	RootPkgPath() string
}
