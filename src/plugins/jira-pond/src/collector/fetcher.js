const axios = require('axios')
const ProxyAgent = require('proxy-agent')
const config = require('@config/resolveConfig').jira
const maxRetry = config.maxRetry || 3
const timeout = config.timeout || 10000

async function fetch (resourceUri) {
  let retry = 0
  let res, lastError
  while (retry < maxRetry) {
    console.log(`INFO: jira fetching data from ${resourceUri}, retry: #${retry}`)
    const abort = axios.CancelToken.source()
    const id = setTimeout(
      () => abort.cancel(`Timeout of ${timeout}ms.`),
      timeout
    )
    try {
      res = await axios.get(`${config.host}/rest/${resourceUri}`, {
        headers: {
          Accept: 'application/json',
          Authorization: `Basic ${config.basicAuth}`
        },
        agent: config.proxy && new ProxyAgent(config.proxy),
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

module.exports = { fetch, fetchPaged }
