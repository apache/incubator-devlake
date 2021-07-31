const axios = require('axios')
const ProxyAgent = require('proxy-agent')
const config = require('@config/resolveConfig').gitlab
const { host, apiPath, token } = config
const maxRetry = config.maxRetry || 3

async function fetch (resourceUri) {
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
}

async function * fetchPaged (resourceUri, page = 1, pageSize = 100) {
  resourceUri = `${resourceUri}${resourceUri.includes('?') ? '&' : '?'}`

  while (page) {
    const res = await fetch(`${resourceUri}per_page=${pageSize}&page=${page}`)
    page = res.headers['x-next-page']
    for (const item of res.data) {
      yield item
    }
  }
}

module.exports = { fetch, fetchPaged }
