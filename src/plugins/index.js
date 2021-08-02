require('module-alias/register')
const _merge = require('lodash/merge')
const registryConfig = require('../../config/plugins')
const dbConnector = require('@mongo/connection')
const enrichedDb = require('@db/postgres')

const plugins = {}
const collection = {}
const enrichment = {}

registryConfig.forEach((item, index) => {
  const { package: packageName, name: pluginName, configuration: pluginConfiguration } = item
  if (!packageName) {
    throw new Error(`package fields is missing for plugin #${index}, check your config/plugin.js file`)
  }
  if (!pluginName) {
    throw new Error(`name fields is missing for plugin #${index}, check your config/plugin.js file`)
  }

  // load module
  const plugin = require(packageName)
  plugin.packageName = packageName
  plugin.pluginName = pluginName

  // setup configuration for plugin
  //   plugin can have a default confgiration, which is subjected to be overwritten
  //   TODO: check if configuration keys are perfectly matched
  _merge(plugin.configuration, pluginConfiguration)

  // register plugin
  plugins[pluginName] = plugin

  // distribute collector
  if (plugin.collector) {
    const { name: collectorName, exec: collector } = plugin.collector
    if (!collectorName) {
      throw new Error(`Error: name is missing for collector of ${pluginName}`)
    }
    if (collection[collectorName]) {
      throw new Error(`Error: conflicted plugin name ${pluginName}`)
    }
    if (!collector) {
      throw new Error(`Error: exec is missing for collector of ${pluginName}`)
    }
    collector.collectorName = collectorName
    collection[pluginName] = collector
  }

  // distribute enricher
  if (plugin.enricher) {
    const { name: enricherName, exec: enricher } = plugin.enricher
    if (!enricherName) {
      throw new Error(`Error: name is missing for enricher of ${pluginName}`)
    }
    if (!enricher) {
      throw new Error(`Error: exec is missing for enricher of ${pluginName}`)
    }
    enricher.enricherName = enricherName
    enrichment[pluginName] = enricher
  }
})

// run all initializations of all plugins
async function initialize () {
  console.log('INFO: initializing plugins')
  const {
    db: rawDb, client
  } = await dbConnector.connect()

  try {
    // assuming no dependencies among plugins
    await Promise.all(
      Object.values(plugins)
        .filter(plugin => plugin.initialize)
        .map(plugin => plugin.initialize(rawDb, enrichedDb, plugins))
    )
  } finally {
    dbConnector.disconnect(client)
  }
  console.log('INFO: initializing plugins done!')
}

module.exports = { collection, enrichment, initialize }

// initialize all plugin with `node src/plugins/index.js`
if (require.main === module) {
  initialize()
}
