import axios from 'axios'

const headers = {}

const handleErrorResponse = (e) => {
  let errorResponse = { success: false, message: e.message, data: null, status: 504 }
  if (e.response) {
    errorResponse = { ...errorResponse, data: e.response ? e.response.data : null, status: e.response ? e.response.status : 504 }
  }
  return errorResponse
}

export default {
  post: async (url, body) => {
    return await axios.post(
      url,
      body,
      {
        headers
      }
    ).catch(e => handleErrorResponse(e))
  },
  get: async (url) => {
    return await axios.get(
      url,
      {
        headers
      }
    ).catch(e => handleErrorResponse(e))
  },
  put: async (url, body) => {
    return await axios.put(
      url,
      body,
      {
        headers
      }
    ).catch(e => handleErrorResponse(e))
  },
  patch: async (url, body) => {
    return await axios.patch(
      url,
      body,
      {
        headers
      }
    ).catch(e => handleErrorResponse(e))
  },
  delete: async (url, body) => {
    return await axios.delete(
      url,
      body,
      {
        headers
      }
    ).catch(e => handleErrorResponse(e))
  }
}
