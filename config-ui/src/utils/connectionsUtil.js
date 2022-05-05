import { DEVLAKE_ENDPOINT } from './config.js'
import request from './request'
import PLUGINS from '../data/availablePlugins'

const ConnectionsUtil = {
  getPluginConnections: async () => {
    const pluginsToSet = []
    const errors = []
    for (const plugin of PLUGINS) {
      try {
        const res = await request.get(`${DEVLAKE_ENDPOINT}/plugins/${plugin}/connections`)
        if (res?.data?.length > 0) {
          pluginsToSet.push(plugin)
        }
      } catch (error) {
        errors.push(error)
      }
    }
    if (errors.length > 0) {
      console.log('errors', errors)
    }
    return pluginsToSet
  }
}

export default ConnectionsUtil
