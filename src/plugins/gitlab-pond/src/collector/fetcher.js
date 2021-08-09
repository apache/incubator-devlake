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
    console.log(`INFO >>> GitLab fetching data from: ${resourceUri}, retry: #${retry}`)
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
    throw new Error('INFO >>> GitLab fetching data failed. Retry limit exceeding. retry: #', retry)
  }
  if (res.data && res.data.message) {
    throw new Error(`INFO >>> GitLab fetching data failed. Status: ${res.status} Message: ${res.data.message}`)
  }
  console.log(`INFO >>> GitLab fetched data from ${resourceUri}`)
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
