const COLLECTION_DEFAULTS = require('../data/collectionDefaults')
const PLUGINS = require('../data/availablePlugins')

const TriggersUtil = {
  getCollectionJson: (plugins = []) => {
    const arrayOfCollectors = []
    const arrayOfDomainConverters = []

    for (const plugin of plugins) {
      if (PLUGINS.includes(plugin)) {
        const collectorJson = module.exports.getCollectorJson(plugin)
        const domainJson = module.exports.getDomainJson(plugin)
        arrayOfCollectors.push(collectorJson)
        arrayOfDomainConverters.push(domainJson)
      }
    }
    return [arrayOfCollectors, arrayOfDomainConverters]
  },
  getCollectorJson: (name) => {
    return {
      Plugin: name,
      Options: COLLECTION_DEFAULTS[name]?.options,
    }
  },
  getDomainJson: (name) => {
    return {
      Plugin: COLLECTION_DEFAULTS[name]?.domainPlugin?.name,
      Options: COLLECTION_DEFAULTS[name]?.domainPlugin?.options,
    }
  }
}

module.exports = TriggersUtil
