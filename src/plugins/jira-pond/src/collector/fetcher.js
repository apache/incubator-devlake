const axios = require('axios')

const config = require('@config/resolveConfig').jira

module.exports = {
  async fetch (resourceUri) {
    try {
      const Authorization = config.restAuth.enabled
        ? `Basic ${Buffer.from(`${config.restAuth.username}:${config.restAuth.apiToken}`).toString('base64')}`
        : `Basic ${config.basicAuth}`

      const response = await axios.get(`${config.host}/rest/api/3/${resourceUri}`, {
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json',
          Authorization
        }
      })

      const hasAuthFailures = (
        header,
        value,
        loginReasonXHeader = 'x-seraph-loginreason',
        deniedReason = 'AUTHENTICATION_DENIED'
      ) => {
        return header === loginReasonXHeader && value === deniedReason
      }

      if (response &&
          response.headers.find((value, header) => hasAuthFailures(header.toLowerCase(), value))
      ) {
        throw new Error('CAPTCHA Triggered! Too many failed authentication attempts w/ REST API.')
      }

      return response.data
    } catch (error) {
      console.error(error)
    }
  }
}
