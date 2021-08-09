require('module-alias/register')
const config = require('@config/resolveConfig')
const dispatch = require('./collection/dispatch')
const axios = require('axios')


const loopTimeInMinutes = 60
const loopTime = 1000 * 60 * loopTimeInMinutes
const delay = 5000

module.exports = {
  loop: async () => {
    console.log('INFO: Collection Loop: Starting data enrichment loop')
    // run on startup
    setTimeout(async ()=> {
      await module.exports.collectionLoop()
    }, delay)

    // run on interval
    setInterval(async () => {
      await module.exports.collectionLoop()
    }, loopTime)
  },

  getMessageFromConfig: (config) => {
    return {
      jira: config.jira && config.jira.dataCollection,
      gitlab: config.gitlab && config.gitlab.dataCollection
    }
  },

  collectionLoop: async () => {
    const message = module.exports.getMessageFromConfig(config)
    console.log('INFO: Collection Loop: processing message', message)
    await axios.post(
      `http://localhost:${process.env.COLLECTION_PORT || 3001}`,
      message,
      { headers: { 'x-token': config.token || '' } }
    )
  }
}
module.exports.loop()