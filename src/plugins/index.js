require('module-alias/register')
const registryConfig = require('../../config/plugins')
const dbConnector = require('@mongo/connection')
const enrichedDb = require('@db/postgres')

module.exports = {
  collection: buildPluginRegistry('collector'),
  enrichment: buildPluginRegistry('enricher')
}

async function buildPluginRegistry (type) {
  const pluginRegistry = pluginRegistryFactory()

  for (const pluginModule of registryConfig) {
    if (pluginModule.type !== type) {
      continue
    }

    const plugin = require(pluginModule.package)

    pluginRegistry.register(plugin[pluginModule.name])
  }

  await pluginRegistry.initialize()

  return pluginRegistry
}

function pluginRegistryFactory () {
  return {
    plugins: {},

    register (plugin) {
      const { name, exec } = plugin
      this.plugins[name] = exec
    },

    // initialize all plugins due to some of them need to do some preparation
    // all database migration should be moved to plugin initialization in the
    // future
    async initialize () {
      const {
        db: rawDb, client
      } = await dbConnector.connect()

      try {
        for (const key of Object.keys(this.plugins)) {
          const plugin = this.plugins[key]
          if (typeof plugin.initialize === 'function') {
            plugin.initialize(rawDb, enrichedDb, this.plugins)
          }
        }
      } finally {
        dbConnector.disconnect(client)
      }
    }
  }
}

module.exports = { buildPluginRegistry }
