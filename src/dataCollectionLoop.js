require('module-alias/register')
const config = require('@config/resolveConfig').cron || {}
const axios = require('axios')

const loopTimeInMinutes = config.loopIntervalInMinutes || 60
const loopTime = 1000 * 60 * loopTimeInMinutes
const delay = 5000

module.exports = {
  loop: async () => {
    if (!config.job || (!config.job.gitlab && !config.job.jira)) {
      console.error('please define your job request body first!')
      return
    }
    try {
      console.log('INFO: Collection Loop: Starting data collection loop')
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

  collectionLoop: async () => {
    try {
      console.log('INFO: Collection Loop: processing message', config.job)
      await axios.post(
        `http://localhost:${process.env.COLLECTION_PORT || 3001}`,
        config.job,
        { headers: { 'x-token': config.token || '' } }
      )
    } catch (error) {
      console.error('*********************************************')
      console.error('ERROR: Failed to run collection loop. You may need to set up a config/docker.js file that has the right properties.')
      console.error('*********************************************')
    }
  }
}
module.exports.loop()
