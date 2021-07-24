const axios = require('axios')

const config = require('@config/resolveConfig').jira

module.exports = {
  async fetch (resourceUri) {
    try {
      const Authorization = config.restAuth.enabled
        ? `Basic ${new Buffer(`${config.restAuth.username}:${config.restAuth.apiToken}`).toString('base64')}`
        : `Basic ${config.basicAuth}`

      const response = await axios.get(`${config.host}/rest/api/3/${resourceUri}`, {
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json',
          Authorization
        }
      })

      return response.data
    } catch (error) {
      console.error(error)
    }
  }
}
