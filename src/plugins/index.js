const registryConfig = require('../../config/plugins.json')

module.exports = {
  collection: buildPluginRegistry('collector'),
  enrichment: buildPluginRegistry('enricher')
}

function buildPluginRegistry (type) {
  const pluginRegistry = pluginRegistryFactory()

  for (const pluginModule of registryConfig) {
    if (pluginModule.type !== type) {
      continue
    }

    const plugin = require(pluginModule.package)

    pluginRegistry.register(plugin[pluginModule.name])
  }

  return pluginRegistry
}

function pluginRegistryFactory () {
  return {
    plugins: {},

    register (plugin) {
      const { name, exec } = plugin
      this.plugins[name] = exec
    }
  }
}
