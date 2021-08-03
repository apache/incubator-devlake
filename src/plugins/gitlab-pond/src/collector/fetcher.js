const axios = require('axios')
const ProxyAgent = require('proxy-agent')
const config = require('@config/resolveConfig').gitlab
const { host, apiPath, token } = config
const maxRetry = config.maxRetry || 3
const timeout = config.timeout || 10000

async function fetch (resourceUri) {
  let retry = 0
  let res
  while (retry < maxRetry) {
    console.log(`INFO >>> gitlab fetching data from ${resourceUri} #${retry}`)
    const abort = axios.CancelToken.source()
    const id = setTimeout(
      () => abort.cancel(`Timeout of ${timeout}ms.`),
      timeout
    )
    try {
      res = await axios.get(`${host}/${apiPath}/${resourceUri}`, {
        headers: { 'PRIVATE-TOKEN': token },
        agent: config.proxy && new ProxyAgent(config.proxy),
        cancelToken: abort.token
      })
      clearTimeout(id)
      break
    } catch (error) {
      console.log('ERROR: ', error)
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
