const axios = require('axios')
const ProxyAgent = require('proxy-agent')
const config = require('@config/resolveConfig').gitlab
const { host, apiPath, token } = config

module.exports = {
  async fetch (resourceUri, maxRetry = 3) {
    let retry = 0
    let res
    while (retry < maxRetry) {
      console.log(`INFO >>> gitlab fetching data from ${resourceUri} #${retry}`)
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
    if (!res) {
      throw new Error('INFO >>> gitlab fetching data failed: retry limit exceeding', retry)
    }
    if (res.data && res.data.message) {
      throw new Error(`INFO >>> gitlab fetching data failed: status: ${res.status} messgae: ${res.data.message}`)
    }
    console.log(`INFO >>> gitlab fetched data from ${resourceUri}`)
    return res
  },

  async * fetchPaged (resourceUri, maxRetry = 3, startAt = 0, pageSize = 100) {
    resourceUri = `${resourceUri}${resourceUri.includes('?') ? '&' : '?'}`

    let page = 1

    while (true) {
      const res = await module.exports.fetch(`${resourceUri}per_page=${pageSize}&page=${page}`, maxRetry)
      // we always simply want the next page of data
      page += 1
      // If no data is returned, we must be on the last page of refults
      if (res.data.length === 0) {
        console.log('INFO: fetchPaged: No data for this page, done fetching paged')
        break
      }
      for (const item of res.data) {
        yield item
      }
    }
  }
}
