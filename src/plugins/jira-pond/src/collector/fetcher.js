const axios = require('axios')
const ProxyAgent = require('proxy-agent')
const config = require('@config/resolveConfig').jira
const maxRetry = config.maxRetry || 3
const timeout = config.timeout || 10000

async function fetch (resourceUri) {
  let retry = 0
  let res
  while (retry < maxRetry) {
    console.log(`INFO >>> jira fetching data from ${resourceUri} #${retry}`)
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
      retry++
    }
  }
  if (!res) {
    throw new Error('INFO >>> jira fetching data failed: retry limit exceeding', retry)
  }
  if (res.data && res.data.message) {
    throw new Error(`INFO >>> jira fetching data failed: status: ${res.status} messgae: ${res.data.message}`)
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
