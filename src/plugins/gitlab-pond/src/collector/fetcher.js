const axios = require('axios')
const ProxyAgent = require('proxy-agent')
const config = require('@config/resolveConfig').gitlab
const { host, apiPath, token } = config

module.exports = {
  async fetch (resourceUri, maxRetry = 3) {
    let retry = 0
    let res
    while (retry < maxRetry) {
      console.log(`[gitlab] Fetching data from ${resourceUri} #${retry}`)
      try {
        res = await axios.get(`${host}/${apiPath}/${resourceUri}`, {
          headers: { 'PRIVATE-TOKEN': token },
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
    console.log(`[gitlab] Fetched data from ${resourceUri}`)
    return res
  },

  async * fetchPaged (resourceUri, maxRetry = 3, startAt = 0, pageSize = 100) {
    resourceUri = `${resourceUri}${resourceUri.includes('?') ? '&' : '?'}`

    let page = 1

    while (page) {
      const res = await module.exports.fetch(`${resourceUri}per_page=${pageSize}&page=${page}`, maxRetry)
      page = res.headers['x-next-page']
      for (const item of res.data) {
        yield item
      }
    }
  }
}
