const axios = require('axios')

const config = require('@config/resolveConfig').jira

module.exports = {
  async fetch(resourceUri) {
    try {
      const response = await axios.get(`${config.host}/rest/api/3/${resourceUri}`, {
        headers: {
          Accept: 'application/json',
          Authorization: `Basic ${config.basicAuth}`
        }
      })

      return response.data
    } catch (error) {
      console.error(error)
    }
  }
}