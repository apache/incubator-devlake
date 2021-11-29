const COLLECTION_DEFAULTS = require('../data/collectionDefaults')
const PLUGINS = require('../data/availablePlugins')

const TriggersUtil = {
  getCollectionJson: (plugins = []) => {
    const arrayOfCollectors = []

    for (const plugin of plugins) {
      if (PLUGINS.includes(plugin)) {
        const collectorJson = module.exports.getCollectorJson(plugin)
        arrayOfCollectors.push(collectorJson)
      }
    }
    return [arrayOfCollectors]
  },
  getCollectorJson: (name) => {
    return {
      Plugin: name,
      Options: COLLECTION_DEFAULTS[name]?.options,
    }
  }
}

module.exports = TriggersUtil
