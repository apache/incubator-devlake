const axios = require('axios')
const ProxyAgent = require('proxy-agent')
const { merge } = require('lodash')

const configuration = {
  verified: false,
  host: null,
  apiPath: null,
  token: null,
  timeout: 10000,
  maxRetry: 3,
  skip: {
    commits: false,
    projects: false,
    mergeRequests: false,
    notes: false
  }
}

function configure (config) {
  merge(configuration, config)
  configuration.verified = false
  const { host, apiPath, token } = configuration
  if (!host) {
    throw new Error('gitlab configuration error: host is required')
  }
  if (!apiPath) {
    throw new Error('gitlab configuration error: apiPath is required')
  }
  if (!token) {
    throw new Error('gitlab configuration error: token is required')
  }
  configuration.verified = true
}

async function fetch (resourceUri) {
  if (!configuration.verified) {
    throw new Error('gitlab fetcher is not configured properly!')
  }
  const { host, apiPath, token, proxy, timeout, maxRetry } = configuration
  let retry = 0
  let res, lastError
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
        agent: proxy && new ProxyAgent(proxy),
        cancelToken: abort.token
      })
      clearTimeout(id)
      break
    } catch (error) {
      lastError = error
      if (error.response) {
        const { status } = error.response
        if (status >= 400 && status < 500) { // no point to retry on client side errors
          break
        }
      }
      await new Promise(resolve => setTimeout(resolve, 200))
      retry++
    }
  }
  if (!res) {
    if (lastError && lastError.response) {
      const { status, data } = lastError.response
      lastError = `status: ${status}, body: ${JSON.stringify(data)}`
    }
    throw new Error(`INFO >>> gitlab fetching data failed! retry: ${retry} , last error: ${lastError}`)
  }
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

module.exports = { configure, fetch, fetchPaged }
