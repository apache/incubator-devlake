const axios = require('axios')

// const config = require('@config/resolveConfig').gitlab
const host = 'https://gitlab.com'
const path = '/api/v4'

module.exports = {
  async fetch (requestUri) {
    try {
      const res = await axios.get(`${host}${path}${requestUri}`)
      return res.data
    } catch (error) {
      console.error(error)
    }
  }
}
