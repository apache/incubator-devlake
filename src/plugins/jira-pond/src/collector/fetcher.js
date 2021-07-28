const axios = require('axios')
const ProxyAgent = require('proxy-agent')

const config = require('@config/resolveConfig').jira

module.exports = {
  async fetch (resourceUri) {
    console.log('[jira] Fetching data from ', resourceUri)
    try {
      const response = await axios.get(`${config.host}/rest/api/3/${resourceUri}`, {
        headers: {
          Accept: 'application/json',
          Authorization: `Basic ${config.basicAuth}`
        },
        agent: config.proxy && new ProxyAgent(config.proxy),
        timeout: config.timeout
      })

      return response.data
    } catch (error) {
      console.error(error)
    }
  }
}
