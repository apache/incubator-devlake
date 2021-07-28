const axios = require('axios')
// const { gitlab: { host, apiPath, token } } = require('../../../../../../config/resolveConfig')
const config = require('@config/resolveConfig').gitlab
const { host, apiPath, token } = config
const privateTokenHeader = { 'PRIVATE-TOKEN': token }

module.exports = {
  async fetch (url) {
    try {
      const res = await axios.get(`${host}/${apiPath}/${url}`, { headers: privateTokenHeader })
      return res.data
    } catch (error) {
      console.error(error)
    }
  }
  // async fetchOne (host, path, modelName, id) {
  //   try {
  //     const res = await axios.get(`${host}/${path}/${modelName}/${id}`)
  //     return res.data
  //   } catch (error) {
  //     console.error(error)
  //   }
  // },
  // async fetchAll (host, path, modelName) {
  //   try {
  //     const res = await axios.get(`${host}/${path}/${modelName}`)
  //     return res.data
  //   } catch (error) {
  //     console.error(error)
  //   }
  // }
}
