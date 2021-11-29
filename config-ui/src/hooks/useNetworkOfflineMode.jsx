import React, { useCallback, useState, useEffect } from 'react'
import { useHistory } from 'react-router-dom'

function useNetworkOfflineMode (offlineStatus = 504, offlineRoute = '/offline') {
  const history = useHistory()
  const [status, setStatus] = useState()
  const [response, setResponse] = useState()

  const handleOfflineMode = useCallback((statusCode, xhrResponse) => {
    setStatus(statusCode)
    setResponse(xhrResponse)
  }, [])

  useEffect(() => {
    if (status && response && status === offlineStatus) {
      history.push(offlineRoute)
    }
  }, [status, response, offlineStatus, offlineRoute, history])

  return {
    handleOfflineMode,
    response
  }
}

export default useNetworkOfflineMode
