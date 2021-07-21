const axios = require('axios')

module.exports = {
  async fetchOne (host, path, modelName, id) {
    try {
      const res = await axios.get(`${host}/${path}/${modelName}/${id}`)
      return res.data
    } catch (error) {
      console.error(error)
    }
  },
  async fetchAll (host, path, modelName){
    try {
      const res = await axios.get(`${host}/${path}/${modelName}`)
      return res.data
    } catch (error) {
      console.error(error)
    }
  } 
}
