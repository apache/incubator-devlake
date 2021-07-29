const axios = require('axios')
const ProxyAgent = require('proxy-agent')

const config = require('@config/resolveConfig').jira

module.exports = {
  async fetch (resourceUri, maxRetry = 3) {
    let retry = 0
    let res
    while (retry < maxRetry) {
      console.log(`[jira] Fetching data from ${resourceUri} #${retry}`)
      try {
        res = await axios.get(`${config.host}/rest/${resourceUri}`, {
          headers: {
            Accept: 'application/json',
            Authorization: `Basic ${config.basicAuth}`
          },
          agent: config.proxy && new ProxyAgent(config.proxy),
          timeout: config.timeout
        })
        break
      } catch (error) {
        retry++
      }
    }
    if (res.data && res.data.message) {
      throw new Error(`status: ${res.status} messgae: ${res.data.message}`)
    }
    return res
  },

  async * fetchPaged (resourceUri, prop = 'values', maxRetry = 3, startAt = 0, pageSize = 100) {
    resourceUri = `${resourceUri}${resourceUri.includes('?') ? '&' : '?'}`

    let total = Number.MAX_VALUE

    while (startAt < total) {
      const res = await module.exports.fetch(`${resourceUri}maxResults=${pageSize}&startAt=${startAt}`, maxRetry)
      total = res.data.total || 0
      const list = res.data[prop]
      startAt += list.length
      for (const item of list) {
        yield item
      }
    }
  }
}
