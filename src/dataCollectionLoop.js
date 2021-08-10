require('module-alias/register')
const config = require('@config/resolveConfig')
const axios = require('axios')

const loopTimeInMinutes = config.lake.loopIntervalInMinutes || 60
const loopTime = 1000 * 60 * loopTimeInMinutes
const delay = 5000

module.exports = {
  loop: async () => {
    try {
      console.log('INFO: Collection Loop: Starting data enrichment loop')
      // run on startup
      setTimeout(async () => {
        await module.exports.collectionLoop()
      }, delay)

      // run on interval
      setInterval(async () => {
        await module.exports.collectionLoop()
      }, loopTime)
    } catch (error) {
      console.error('ERROR: Failed to run loop', error)
    }
  },

  getMessageFromConfig: (config) => {
    return {
      jira: config.jira && config.jira.dataCollection,
      gitlab: config.gitlab && config.gitlab.dataCollection
    }
  },

  collectionLoop: async () => {
    try {
      const message = module.exports.getMessageFromConfig(config)
      console.log('INFO: Collection Loop: processing message', message)
      await axios.post(
        `http://localhost:${process.env.COLLECTION_PORT || 3001}`,
        message,
        { headers: { 'x-token': config.token || '' } }
      )
    } catch (error) {
      console.error('*********************************************')
      console.error('ERROR: Failed to run collection loop. You may need to set up a config/docker.js file that has the right properties. The main ones to look at would be the config.jira.dataCollection and config.gitlab.dataCollection. You can also look at config/docker.sample.js')
      console.error('*********************************************')
    }
  }
}
module.exports.loop()
