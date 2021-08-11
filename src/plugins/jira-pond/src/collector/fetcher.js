const axios = require('axios')
const ProxyAgent = require('proxy-agent')
const { merge } = require('lodash')

const configuration = {
  verified: false,
  host: null,
  basicAuth: null,
  proxy: null,
  timeout: 10000,
  maxRetry: 3
}

async function configure (config) {
  merge(configuration, config)
  configuration.verified = false
  const { host, basicAuth } = configuration
  if (!host) {
    throw new Error('jira configuration error: host is required')
  }
  if (!basicAuth) {
    throw new Error('jira configuration error: basicAuth is required')
  }
  configuration.verified = true
}

async function fetch (resourceUri) {
  if (!configuration.verified) {
    throw new Error('not configured!')
  }
  let retry = 0
  let res, lastError
  const { host, basicAuth, proxy, timeout, maxRetry } = configuration
  while (retry < maxRetry) {
    console.log(`INFO: jira fetching data from ${resourceUri}, retry: #${retry}`)
    const abort = axios.CancelToken.source()
    const id = setTimeout(
      () => abort.cancel(`Timeout of ${timeout}ms.`),
      timeout
    )
    try {
      res = await axios.get(`${host}/rest/${resourceUri}`, {
        headers: {
          Accept: 'application/json',
          Authorization: `Basic ${basicAuth}`
        },
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
    throw new Error(`INFO >>> jira fetching data failed! retry: ${retry} , last error: ${lastError}`)
  }
  return res
}

async function * fetchPaged (resourceUri, prop = 'values', startAt = 0, pageSize = 100) {
  resourceUri = `${resourceUri}${resourceUri.includes('?') ? '&' : '?'}`

  let total = Number.MAX_VALUE

  while (startAt < total) {
    const res = await fetch(`${resourceUri}maxResults=${pageSize}&startAt=${startAt}`)
    total = res.data.total || 0
    const list = res.data[prop]
    startAt += list.length
    for (const item of list) {
      yield item
    }
  }
}

module.exports = { configure, fetch, fetchPaged }
