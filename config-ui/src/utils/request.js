import axios from 'axios'

const headers = {}

export default {
  post: async (url, body) => {
    return await axios.post(
      url,
      body,
      {
        headers
      }
    )
  },
  get: async (url) => {
    return await axios.get(
      url,
      {
        headers
      }
    )
  },
  put: async (url, body) => {
    return await axios.put(
      url,
      body,
      {
        headers
      }
    )
  }
}
