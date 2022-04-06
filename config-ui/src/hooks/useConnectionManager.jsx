import { useState, useEffect, useCallback } from 'react'
import {
  useHistory
} from 'react-router-dom'
import { ToastNotification } from '@/components/Toast'
import { DEVLAKE_ENDPOINT } from '@/utils/config'
import request from '@/utils/request'
import { NullConnection } from '@/data/NullConnection'
import { Providers, ProviderSourceLimits } from '@/data/Providers'

import useNetworkOfflineMode from '@/hooks/useNetworkOfflineMode'

function useConnectionManager ({
  activeProvider,
  connectionId,
  // activeConnection,
  // setActiveConnection,
  // name = null,
  // endpointUrl = null,
  // token = null,
  // username = null,
  // password = null,
  // isTesting, setIsTesting,
  // isSaving, setIsSaving,
  // testStatus, setTestStatus,
  // errors, setErrors,
  // showError, setShowError
}, updateMode = false) {
  const history = useHistory()
  const { handleOfflineMode } = useNetworkOfflineMode()

  const [name, setName] = useState()
  const [endpointUrl, setEndpointUrl] = useState()
  const [proxy, setProxy] = useState()
  const [token, setToken] = useState()
  const [username, setUsername] = useState()
  const [password, setPassword] = useState()

  const [isSaving, setIsSaving] = useState(false)
  const [isFetching, setIsFetching] = useState(false)
  const [isRunning, setIsRunning] = useState(false)
  const [isTesting, setIsTesting] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)
  const [errors, setErrors] = useState([])
  const [showError, setShowError] = useState(false)
  const [testStatus, setTestStatus] = useState(0) //  0=Pending, 1=Success, 2=Failed
  const [sourceLimits, setSourceLimits] = useState(ProviderSourceLimits)

  const [activeConnection, setActiveConnection] = useState(NullConnection)
  const [allConnections, setAllConnections] = useState([])
  const [testedConnections, setTestedConnections] = useState([])
  const [connectionCount, setConnectionCount] = useState(0)
  const [connectionLimitReached, setConnectionLimitReached] = useState(false)

  const [saveComplete, setSaveComplete] = useState(false)
  const [deleteComplete, setDeleteComplete] = useState(false)

  const testConnection = useCallback((notify = true, manualPayload = {}, onSuccess = () => {}, onFail = () => {}) => {
    setIsTesting(true)
    setShowError(false)
    ToastNotification.clear()
    // TODO: run Save first
    const runTest = async () => {
      let connectionPayload
      switch (activeProvider.id) {
        case Providers.JIRA:
          connectionPayload = { endpoint: endpointUrl, auth: token, proxy: proxy }
          break
        case Providers.GITHUB:
          connectionPayload = { endpoint: endpointUrl, auth: token, proxy: proxy }
          break
        case Providers.JENKINS:
          connectionPayload = { endpoint: endpointUrl, username: username, password: password }
          break
        case Providers.GITLAB:
          connectionPayload = { endpoint: endpointUrl, auth: token, proxy: proxy }
          break
      }
      connectionPayload = { ...connectionPayload, ...manualPayload }
      const testUrl = `${DEVLAKE_ENDPOINT}/plugins/${activeProvider.id}/test`
      console.log('INFO >>> Endopoint URL & Payload for testing: ', testUrl, connectionPayload)
      const res = await request.post(testUrl, connectionPayload)
      console.log('res.data', res.data)
      if (res?.data?.Success && res.status === 200) {
        setIsTesting(false)
        setTestStatus(1)
        if (notify) {
          ToastNotification.show({ message: `Connection test OK. ${connectionPayload.endpoint}`, intent: 'success', icon: 'small-tick' })
        }
        onSuccess(res)
      } else {
        setIsTesting(false)
        setTestStatus(2)
        const errorMessage = 'Connection test FAILED. ' + res.message
        if (notify) {
          ToastNotification.show({ message: errorMessage, intent: 'danger', icon: 'error' })
        }
        onFail(res)
      }
    }
    runTest()
  }, [activeProvider.id, endpointUrl, password, proxy, token, username])

  const saveConnection = (configurationSettings = {}) => {
    setIsSaving(true)
    let connectionPayload = { ...configurationSettings }
    switch (activeProvider.id) {
      case Providers.JIRA:
        connectionPayload = { name: name, Endpoint: endpointUrl, BasicAuthEncoded: token, Proxy: proxy, ...connectionPayload }
        break
      case Providers.GITHUB:
        connectionPayload = { name: name, GITHUB_ENDPOINT: endpointUrl, GITHUB_AUTH: token, GITHUB_PROXY: proxy, ...connectionPayload }
        break
      case Providers.JENKINS:
        // eslint-disable-next-line max-len
        connectionPayload = { name: name, JENKINS_ENDPOINT: endpointUrl, JENKINS_USERNAME: username, JENKINS_PASSWORD: password, ...connectionPayload }
        break
      case Providers.GITLAB:
        connectionPayload = { name: name, GITLAB_ENDPOINT: endpointUrl, GITLAB_AUTH: token, GITLAB_PROXY: proxy, ...connectionPayload }
        break
    }

    let saveResponse = {
      success: false,
      connection: {
        ...connectionPayload
      },
      errors: []
    }

    const saveConfiguration = async (configPayload) => {
      try {
        setShowError(false)
        setErrors([])
        ToastNotification.clear()
        const s = await request.post(`${DEVLAKE_ENDPOINT}/plugins/${activeProvider.id}/sources`, configPayload)
        console.log('>> CONFIGURATION SAVED SUCCESSFULLY', configPayload, s)
        saveResponse = {
          ...saveResponse,
          success: [200, 201].includes(s.status),
          connection: { ...s.data },
          errors: s.isAxiosError ? [s.message] : []
        }
      } catch (e) {
        saveResponse.errors.push(e.message)
        setErrors(saveResponse.errors)
        console.log('>> CONFIGURATION FAILED TO SAVE', configPayload, e)
      }
    }

    const modifyConfiguration = async (configPayload) => {
      try {
        setShowError(false)
        setErrors([])
        ToastNotification.clear()
        const s = await request.put(`${DEVLAKE_ENDPOINT}/plugins/${activeProvider.id}/sources/${activeConnection.ID}`, configPayload)
        console.log('>> CONFIGURATION MODIFIED SUCCESSFULLY', configPayload, s)
        saveResponse = {
          ...saveResponse,
          success: [200, 201].includes(s.status),
          connection: { ...s.data },
          errors: s.isAxiosError ? [s.message] : []
        }
      } catch (e) {
        saveResponse.errors.push(e.message)
        setErrors(saveResponse.errors)
        console.log('>> CONFIGURATION FAILED TO UPDATE', configPayload, e)
      }
    }

    if (updateMode && activeConnection) {
      modifyConfiguration(connectionPayload)
    } else {
      saveConfiguration(connectionPayload)
    }

    setTimeout(() => {
      if (saveResponse.success && errors.length === 0) {
        ToastNotification.show({ message: 'Connection saved successfully.', intent: 'success', icon: 'small-tick' })
        setShowError(false)
        setIsSaving(false)
        setSaveComplete(saveResponse.connection)
        if (!updateMode) {
          history.push(`/integrations/${activeProvider.id}`)
        }
      } else {
        ToastNotification.show({ message: 'Connection failed to save, please try again.', intent: 'danger', icon: 'error' })
        setShowError(true)
        setIsSaving(false)
        setSaveComplete(false)
      }
    }, 2000)
  }

  const runCollection = (options = {}) => {
    setIsRunning(true)
    ToastNotification.show({ message: 'Triggered Collection Process', intent: 'info', icon: 'info' })
    console.log('>> RUNNING COLLECTION PROCESS', isRunning)
    // Run Collection Tasks...
  }

  // const fetchConnection = async () => {
  const fetchConnection = useCallback(() => {
    console.log('>> FETCHING CONNECTION....')
    try {
      setIsFetching(true)
      setErrors([])
      ToastNotification.clear()
      console.log('>> FETCHING CONNECTION SOURCE')
      const fetch = async () => {
        const f = await request.get(`${DEVLAKE_ENDPOINT}/plugins/${activeProvider.id}/sources/${connectionId}`)
        const connectionData = f.data
        console.log('>> RAW CONNECTION DATA FROM API...', connectionData)
        setActiveConnection({
          ...connectionData,
          name: connectionData.name || connectionData.Name,
          endpoint: connectionData.endpoint || connectionData.Endpoint,
          proxy: connectionData.proxy || connectionData.Proxy,
          username: connectionData.username || connectionData.Username,
          password: connectionData.password || connectionData.Password
        })
        setTimeout(() => {
          setIsFetching(false)
        }, 500)
      }
      fetch()
      // setIsFetching(false)
    } catch (e) {
      setIsFetching(false)
      setActiveConnection(NullConnection)
      setErrors([e.message])
      ToastNotification.show({ message: `${e}`, intent: 'danger', icon: 'error' })
      console.log('>> FAILED TO FETCH CONNECTION', e)
    }
  }, [activeProvider.id, connectionId])

  const fetchAllConnections = useCallback(async (notify = false) => {
    try {
      setIsFetching(true)
      setErrors([])
      ToastNotification.clear()
      console.log('>> FETCHING ALL CONNECTION SOURCES')
      const f = await request.get(`${DEVLAKE_ENDPOINT}/plugins/${activeProvider.id}/sources`)
      console.log('>> RAW ALL CONNECTIONS DATA FROM API...', f.data)
      const providerConnections = [].concat(f.data || []).map((conn, idx) => {
        return {
          ...conn,
          status: 0,
          id: conn.ID,
          name: conn.name,
          endpoint: conn.endpoint,
          errors: []
        }
      })
      // setConnections(providerConnections)
      if (notify) {
        ToastNotification.show({ message: 'Loaded all connections.', intent: 'success', icon: 'small-tick' })
      }
      setAllConnections(providerConnections)
      setConnectionCount(f.data?.length)
      setConnectionLimitReached(sourceLimits[activeProvider.id] && f.data?.length >= sourceLimits[activeProvider.id])
      setIsFetching(false)
    } catch (e) {
      console.log('>> FAILED TO FETCH ALL CONNECTIONS', e)
      ToastNotification.show({ message: `Failed to Load Connections - ${e.message}`, intent: 'danger', icon: 'error' })
      setIsFetching(false)
      setAllConnections([])
      setConnectionCount(0)
      setConnectionLimitReached(false)
      setErrors([e.message])
      handleOfflineMode(e.response.status, e.response)
    }
  }, [activeProvider.id, sourceLimits, handleOfflineMode])

  const deleteConnection = useCallback(async (connection) => {
    try {
      setIsDeleting(true)
      setErrors([])
      console.log('>> TRYING TO DELETE CONNECTION...', connection)
      const d = await request.delete(`${DEVLAKE_ENDPOINT}/plugins/${activeProvider.id}/sources/${connection.ID}`)
      console.log('>> CONNECTION DELETED...', d)
      setIsDeleting(false)
      setDeleteComplete({
        provider: activeProvider,
        connection: d.data
      })
    } catch (e) {
      setIsDeleting(false)
      setDeleteComplete(false)
      setErrors([e.message])
      console.log('>> FAILED TO DELETE CONNECTION', e)
    }
  }, [activeProvider.id])

  const getConnectionName = useCallback((connectionId, sources) => {
    const source = sources.find(s => s.id === connectionId)
    return source ? source.title : '(Instance)'
  }, [])

  const testAllConnections = useCallback((connections) => {
    console.log('>> TESTING ALL CONNECTION SOURCES...')
    connections.forEach((c, cIdx) => {
      console.log('>>> TESTING CONNECTION INSTANCE...', c)
      const notify = false
      const payload = {
        endpoint: c.Endpoint || c.endpoint,
        username: c.Username,
        password: c.Password,
        auth: c.basicAuthEncoded || c.Auth,
        proxy: c.Proxy || c.Proxy
      }
      const onSuccess = (res) => {
        setTestedConnections([...connections.filter(oC => oC.ID !== c.ID), { ...c, status: 1 }])
      }
      const onFail = (res) => {
        setTestedConnections([...connections.filter(oC => oC.ID !== c.ID), { ...c, status: 2 }])
      }
      testConnection(notify, payload, onSuccess, onFail)
    })
  }, [testConnection])

  useEffect(() => {
    if (activeConnection && activeConnection.ID !== null) {
      setName(activeConnection.name)
      setEndpointUrl(activeConnection.endpoint)
      switch (activeProvider.id) {
        case Providers.JENKINS:
          setUsername(activeConnection.username)
          setPassword(activeConnection.password)
          break
        case Providers.GITLAB:
          setToken(activeConnection.basicAuthEncoded || activeConnection.Auth)
          setProxy(activeConnection.Proxy || activeConnection.proxy)
          break
        case Providers.GITHUB:
          setToken(activeConnection.basicAuthEncoded || activeConnection.Auth)
          setProxy(activeConnection.Proxy || activeConnection.proxy)
          break
        case Providers.JIRA:
          setToken(activeConnection.basicAuthEncoded || activeConnection.Auth)
          setProxy(activeConnection.Proxy || activeConnection.proxy)
          break
      }
      ToastNotification.clear()
      ToastNotification.show({ message: `Fetched settings for ${activeConnection.name}.`, intent: 'success', icon: 'small-tick' })
      console.log('>> FETCHED CONNECTION FOR MODIFY', activeConnection)
    }
  }, [activeConnection, activeProvider.id])

  useEffect(() => {
    if (saveComplete && saveComplete.ID) {
      console.log('>>> CONNECTION MANAGER - SAVE COMPLETE EFFECT RUNNING...')
      setActiveConnection((ac) => {
        return {
          ...ac,
          ...saveComplete
        }
      })
    }
  }, [saveComplete])

  useEffect(() => {
    console.log('>> CONNECTION MANAGER - RECEIVED ACTIVE PROVIDER...', activeProvider)
    if (activeProvider && activeProvider.id) {
      // console.log(activeProvider)
    }
  }, [activeProvider])

  useEffect(() => {
    if (connectionId) {
      console.log('>>>> CONFIGURING CONNECTION ID ... ', connectionId)
      fetchConnection()
    }
  }, [connectionId, fetchConnection])

  useEffect(() => {
    console.log('>> TESTED CONNECTION RESULTS...', testedConnections)
  }, [testedConnections])

  return {
    activeConnection,
    fetchConnection,
    fetchAllConnections,
    testAllConnections,
    testConnection,
    saveConnection,
    deleteConnection,
    runCollection,
    isSaving,
    isTesting,
    isFetching,
    errors,
    showError,
    testStatus,
    name,
    endpointUrl,
    proxy,
    username,
    password,
    token,
    setName,
    setEndpointUrl,
    setProxy,
    setToken,
    setUsername,
    setPassword,
    setIsSaving,
    setIsTesting,
    setIsFetching,
    setErrors,
    setShowError,
    setTestStatus,
    setSourceLimits,
    allConnections,
    testedConnections,
    sourceLimits,
    connectionCount,
    connectionLimitReached,
    Providers,
    saveComplete,
    deleteComplete,
    getConnectionName
  }
}

export default useConnectionManager
