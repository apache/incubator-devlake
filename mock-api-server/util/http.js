const axios = require('axios')
const HttpHelper = {
  async get(url, callback){
    try {
      res = await axios.get(url)
      callback(res, null)
    } catch (error) {
      callback(null, error) 
    }
  }
}

module.exports = HttpHelper