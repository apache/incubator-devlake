/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
import { useCallback, useEffect, useMemo, useState } from 'react'
import { useHistory } from 'react-router-dom'
import { ToastNotification } from '@/components/Toast'
import { DEVLAKE_ENDPOINT } from '@/utils/config'
import request from '@/utils/request'
import { NullConnection } from '@/data/NullConnection'
import { ConnectionStatus, ProviderConfigMap, Providers } from '@/data/Providers'

import useNetworkOfflineMode from '@/hooks/useNetworkOfflineMode'

function useConnectionManager (
  {
    provider: activeProvider,
    connectionId,
  },
  updateMode = false
) {
  const history = useHistory()
  const { handleOfflineMode } = useNetworkOfflineMode()

  const [provider, setProvider] = useState(activeProvider)
  const [initialTokenStore, setInitialTokenStore] = useState({
    0: '',
    1: '',
    2: ''
  })

  const [isSaving, setIsSaving] = useState(false)
  const [isFetching, setIsFetching] = useState(false)
  const [isRunning, setIsRunning] = useState(false)
  const [isTesting, setIsTesting] = useState(false)
  // eslint-disable-next-line no-unused-vars
  const [isDeleting, setIsDeleting] = useState(false)
  const [errors, setErrors] = useState([])
  const [showError, setShowError] = useState(false)
  const [testStatus, setTestStatus] = useState(0) //  0=Pending, 1=Success, 2=Failed
  const [testResponse, setTestResponse] = useState()
  const [allTestResponses, setAllTestResponses] = useState({})

  const [activeConnection, setActiveConnection] = useState(NullConnection)
  const [editingConnection, setEditingConnection] = useState(NullConnection)
  const [allConnections, setAllConnections] = useState([])
  const [allProviderConnections, setAllProviderConnections] = useState([])
  const [domainRepositories, setDomainRepositories] = useState([])
  const [testedConnections, setTestedConnections] = useState([])
  const connectionCount = useMemo(() => allConnections.length, [allConnections])
  const connectionLimitReached = useMemo(() =>
    ProviderConfigMap[provider?.id]?.limit && connectionCount >= ProviderConfigMap[provider.id].limit,
  [connectionCount],
  )

  const [saveComplete, setSaveComplete] = useState(false)
  const [deleteComplete, setDeleteComplete] = useState(false)
  const connectionTestPayload = useMemo(() => ({
    endpoint: editingConnection.endpoint,
    username: editingConnection.username,
    password: editingConnection.password,
    token: editingConnection.token,
    proxy: editingConnection.proxy,
  }), [editingConnection])

  const testConnection = useCallback(
    (
      notify = true,
      manualPayload = {},
      onSuccess = () => {},
      onFail = () => {}
    ) => {
      setIsTesting(true)
      setShowError(false)
      ToastNotification.clear()
      setTestResponse(null)

      const runTest = async () => {
        const payload = Object.keys(manualPayload).length > 0 ? manualPayload : connectionTestPayload
        const testUrl = `${DEVLAKE_ENDPOINT}/plugins/${provider.id}/test`
        console.log(
          'INFO >>> Endpoint URL & Payload for testing: ',
          testUrl,
          payload
        )
        const res = await request.post(testUrl, payload)
        setTestResponse(res.data)
        if ([Providers.GITHUB].includes(provider.id)) {
          console.log('>>> SETTING TOKEN TEST RESPONSE FOR TOKEN >>>', manualPayload?.token || connectionTestPayload.token)
          setAllTestResponses(tRs => ({ ...tRs, [manualPayload?.token]: res.data }))
        }
        if (res.data?.success && res.status === 200) {
          setIsTesting(false)
          setTestStatus(1)
          if (notify) {
            ToastNotification.show({
              message: `Connection test OK. ${payload.endpoint}`,
              intent: 'success',
              icon: 'small-tick',
            })
          }
          onSuccess(res)
        } else {
          setIsTesting(false)
          setTestStatus(2)
          const errorMessage =
            'Connection test FAILED. ' + (res.data ? res.data.message : '')
          if (notify) {
            ToastNotification.show({
              message: errorMessage,
              intent: 'danger',
              icon: 'error',
            })
          }
          onFail(res)
        }
      }
      runTest()
    },
    [provider?.id, connectionTestPayload]
  )

  const setConnectionColumn = useCallback((columnName, value) => {
    console.log('>> SET EDITING CONNECTION COLUMN', columnName, value)
    if (editingConnection[columnName] !== value) {
      setEditingConnection({
        ...editingConnection,
        [columnName]: value,
      })
    }
  }, [editingConnection])

  const saveConnection = (configurationSettings = {}) => {
    setIsSaving(true)
    let connectionPayload = {
      name: editingConnection.name,
      endpoint: editingConnection.endpoint,
      rateLimitPerHour: editingConnection.rateLimitPerHour,
      proxy: editingConnection.proxy,
      ...configurationSettings,
    }
    switch (provider.id) {
      case Providers.JIRA:
        connectionPayload = {
          username: editingConnection.username,
          password: editingConnection.password,
          ...connectionPayload,
        }
        break
      case Providers.TAPD:
        connectionPayload = {
          username: editingConnection.username,
          password: editingConnection.password,
          ...connectionPayload,
        }
        break
      case Providers.GITHUB:
        connectionPayload = {
          token: editingConnection.token,
          ...connectionPayload,
        }
        break
      case Providers.JENKINS:
      // eslint-disable-next-line max-len
        connectionPayload = {
          username: editingConnection.username,
          password: editingConnection.password,
          ...connectionPayload,
        }
        break
      case Providers.GITLAB:
        connectionPayload = {
          token: editingConnection.token,
          ...connectionPayload,
        }
        break
    }

    let saveResponse = {
      success: false,
      connection: {
        ...connectionPayload,
      },
      errors: [],
    }

    const saveConfiguration = async (configPayload) => {
      try {
        setShowError(false)
        setErrors([])
        ToastNotification.clear()
        const s = await request.post(
          `${DEVLAKE_ENDPOINT}/plugins/${provider.id}/connections`,
          configPayload
        )
        console.log('>> CONFIGURATION SAVED SUCCESSFULLY', configPayload, s)
        saveResponse = {
          ...saveResponse,
          success: [200, 201].includes(s.status),
          connection: { ...s.data },
          errors: s.isAxiosError ? [s.message] : [],
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
        // eslint-disable-next-line max-len
        const s = await request.patch(
          `${DEVLAKE_ENDPOINT}/plugins/${provider.id}/connections/${
            activeConnection.id || activeConnection.ID
          }`,
          configPayload
        )
        const silentRefetch = true
        console.log('>> CONFIGURATION MODIFIED SUCCESSFULLY', configPayload, s)
        saveResponse = {
          ...saveResponse,
          success: [200, 201].includes(s.status),
          connection: { ...s.data },
          errors: s.isAxiosError ? [s.message] : [],
        }
        fetchConnection(silentRefetch)
      } catch (e) {
        saveResponse.errors.push(e.message)
        setErrors(saveResponse.errors)
        console.log('>> CONFIGURATION FAILED TO UPDATE', configPayload, e)
      }
    }

    let savePromise
    if (updateMode && activeConnection?.id !== null) {
      savePromise = modifyConfiguration(connectionPayload)
    } else {
      savePromise = saveConfiguration(connectionPayload)
    }
    savePromise.then(() => {
      if (saveResponse.success && errors.length === 0) {
        ToastNotification.show({
          message: 'Connection saved successfully.',
          intent: 'success',
          icon: 'small-tick',
        })
        setShowError(false)
        setIsSaving(false)
        setSaveComplete(saveResponse.connection)
        if (
          [Providers.GITHUB, Providers.JIRA, Providers.GITLAB].includes(provider.id) &&
          editingConnection.token !== '' &&
          editingConnection.token?.toString().split(',').length > 1
        ) {
          testConnection()
        }
        if (!updateMode) {
          history.replace(`/integrations/${provider.id}`)
        }
      }
    })

    setTimeout(() => {
      if (!saveResponse.success || errors.length !== 0) {
        ToastNotification.show({
          message: 'Connection failed to save, please try again.',
          intent: 'danger',
          icon: 'error',
        })
        setShowError(true)
        setIsSaving(false)
        setSaveComplete(false)
      }
    }, 2000)
  }

  const runCollection = (options = {}) => {
    setIsRunning(true)
    ToastNotification.show({
      message: 'Triggered Collection Process',
      intent: 'info',
      icon: 'info',
    })
    console.log('>> RUNNING COLLECTION PROCESS', isRunning)
    // Run Collection Tasks...
  }

  const fetchConnection = useCallback(
    (silent = false, notify = false, cId = null) => {
      console.log(`>> FETCHING CONNECTION [PROVIDER = ${provider.id}]....`)
      try {
        setIsFetching(!silent)
        setErrors([])
        console.log('>> FETCHING CONNECTION SOURCE')
        const fetch = async () => {
          const f = await request.get(
            `${DEVLAKE_ENDPOINT}/plugins/${provider.id}/connections/${cId || connectionId}`
          )
          const connectionData = f.data
          console.log('>> RAW CONNECTION DATA FROM API...', connectionData)
          // @todo: cleanup legacy api pascal-case parameters
          setActiveConnection({
            ...connectionData,
            ID: connectionData.ID || connectionData.id,
            name: connectionData.name || connectionData.Name,
            endpoint: connectionData.endpoint || connectionData.Endpoint,
            proxy: connectionData.proxy || connectionData.Proxy,
            rateLimit: connectionData.rateLimitPerHour,
            username: connectionData.username || connectionData.Username,
            password: connectionData.password || connectionData.Password,
            token: connectionData.token || connectionData.auth
          })
          setTimeout(() => {
            setIsFetching(false)
          }, 500)
        }
        fetch()
      } catch (e) {
        setIsFetching(false)
        setActiveConnection(NullConnection)
        setErrors([e.message])
        ToastNotification.clear()
        ToastNotification.show({
          message: `${e}`,
          intent: 'danger',
          icon: 'error',
        })
        console.log('>> FAILED TO FETCH CONNECTION', e)
      }
    },
    [provider?.id, connectionId]
  )

  const fetchAllConnections = useCallback(
    async (notify = false, allSources = false) => {
      try {
        setIsFetching(true)
        setErrors([])
        ToastNotification.clear()
        console.log('>> FETCHING ALL CONNECTION SOURCES')
        if (allSources) {
          // @todo: build promises dynamically from $integrationsData
          const aC = await Promise.all([
            request.get(
              `${DEVLAKE_ENDPOINT}/plugins/${Providers.JIRA}/connections`
            ),
            request.get(
              `${DEVLAKE_ENDPOINT}/plugins/${Providers.GITLAB}/connections`
            ),
            request.get(
              `${DEVLAKE_ENDPOINT}/plugins/${Providers.JENKINS}/connections`
            ),
            request.get(
              `${DEVLAKE_ENDPOINT}/plugins/${Providers.GITHUB}/connections`
            ),
            request.get(
              `${DEVLAKE_ENDPOINT}/plugins/${Providers.TAPD}/connections`
            ),
          ])
          const builtConnections = aC
            .map((providerResponse) => [].concat(providerResponse.data || []).map(c => ({
              ...c,
              connectionId: c.id,
              provider: providerResponse.config?.url?.split('/')[3],
              // @todo: inject realtime connection status...
              status: ConnectionStatus.ONLINE
            })))
          setAllProviderConnections(builtConnections.flat())
          console.log(
            '>> ALL SOURCE CONNECTIONS: FETCHING ALL CONNECTION FROM ALL DATA SOURCES'
          )
          console.log('>> ALL SOURCE CONNECTIONS: ', aC)
        } else {
          const c = await request.get(
            `${DEVLAKE_ENDPOINT}/plugins/${provider.id}/connections`
          )

          console.log('>> RAW ALL CONNECTIONS DATA FROM API...', c?.data)
          const providerConnections = []
            .concat(Array.isArray(c?.data) ? c?.data : [])
            .map((conn, idx) => {
              return {
                ...conn,
                status: ConnectionStatus.OFFLINE,
                ID: conn.ID || conn.id,
                id: conn.id,
                name: conn.name,
                endpoint: conn.endpoint,
                errors: [],
              }
            })
          setAllConnections(providerConnections)
        }
        if (notify) {
          ToastNotification.show({
            message: 'Loaded all connections.',
            intent: 'success',
            icon: 'small-tick',
          })
        }
        setIsFetching(false)
      } catch (e) {
        console.log('>> FAILED TO FETCH ALL CONNECTIONS', e)
        ToastNotification.show({
          message: `Failed to Load Connections - ${e.message}`,
          intent: 'danger',
          icon: 'error',
        })
        setIsFetching(false)
        setAllConnections([])
        setErrors([e.message])
        handleOfflineMode(e.response?.status, e.response)
      }
    },
    [provider?.id, handleOfflineMode]
  )

  const deleteConnection = useCallback(
    async (connection) => {
      try {
        setIsDeleting(true)
        setErrors([])
        console.log('>> TRYING TO DELETE CONNECTION...', connection)
        const d = await request.delete(
          `${DEVLAKE_ENDPOINT}/plugins/${provider.id}/connections/${
            connection.ID || connection.id
          }`
        )
        console.log('>> CONNECTION DELETED...', d)
        setIsDeleting(false)
        setDeleteComplete({
          provider,
          connection: d.data,
        })
      } catch (e) {
        setIsDeleting(false)
        setDeleteComplete(false)
        setErrors([e.message])
        console.log('>> FAILED TO DELETE CONNECTION', e)
      }
    },
    [provider?.id]
  )

  const getConnectionName = useCallback((connectionId, connections) => {
    const source = connections.find((s) => s.id === connectionId)
    return source ? source.title : '(Instance)'
  }, [])

  const testAllConnections = useCallback(
    (connections) => {
      console.log('>> TESTING ALL CONNECTION SOURCES...')
      connections.forEach((c, cIdx) => {
        console.log('>>> TESTING CONNECTION INSTANCE...', c)
        const notify = false
        const payload = {
          endpoint: c.endpoint,
          username: c.username,
          password: c.password,
          token: c.token,
          proxy: c.proxy,
        }
        const onSuccess = (res) => {
          setTestedConnections((testedConnections) => [
            ...new Set([
              ...testedConnections.filter((oC) => oC.id !== c.id),
              { ...c, status: ConnectionStatus.ONLINE },
            ]),
          ])
        }
        const onFail = (res) => {
          setTestedConnections((testedConnections) => [
            ...new Set([
              ...testedConnections.filter((oC) => oC.ID !== c.ID),
              { ...c, status: ConnectionStatus.DISCONNECTED },
            ]),
          ])
        }
        testConnection(notify, payload, onSuccess, onFail)
      })
    },
    [testConnection]
  )

  const fetchDomainLayerRepositories = useCallback(() => {
    console.log('>> FETCHING DOMAIN LAYER REPOS....')
    try {
      setIsFetching(true)
      setErrors([])
      ToastNotification.clear()
      const fetch = async () => {
        const r = await request.get(`${DEVLAKE_ENDPOINT}/domainlayer/repos`)
        console.log('>> RAW REPOSITORY DATA FROM API...', r.data?.repos)
        setDomainRepositories(r.data?.repos || [])
        setTimeout(() => {
          setIsFetching(false)
        }, 500)
      }
      fetch()
    } catch (e) {
      setIsFetching(false)
      setDomainRepositories([])
      setErrors([e.message])
      ToastNotification.show({
        message: `${e}`,
        intent: 'danger',
        icon: 'error',
      })
      console.log('>> FAILED TO FETCH DOMAIN LAYER REPOS', e)
    }
  }, [])

  const clearConnection = useCallback(() => {
    setEditingConnection(NullConnection)
    setInitialTokenStore({
      0: '',
      1: '',
      2: ''
    })
  }, [])

  useEffect(() => {
    if (activeConnection && activeConnection.id !== null) {
      const connectionToken = activeConnection.auth || activeConnection.token || activeConnection.basicAuthEncoded
      setEditingConnection(activeConnection)
      switch (provider.id) {
        // case Providers.JENKINS:
        //   setUsername(activeConnection.username)
        //   setPassword(activeConnection.password)
        //   setRateLimit(activeConnection.rateLimitPerHour)
        //   break
        // case Providers.GITLAB:
        //   setToken(activeConnection.basicAuthEncoded || activeConnection.token || activeConnection.auth)
        //   setProxy(activeConnection.Proxy || activeConnection.proxy)
        //   setRateLimit(activeConnection.rateLimitPerHour)
        //   break
        case Providers.GITHUB:
          // setToken(connectionToken)
          setInitialTokenStore(connectionToken?.split(',')?.reduce((tS, cT, id) => ({ ...tS, [id]: cT }), {}))
          // setProxy(activeConnection.Proxy || activeConnection.proxy)
          // setRateLimit(activeConnection.rateLimitPerHour)
          break
        // case Providers.JIRA:
        // // setToken(activeConnection.basicAuthEncoded || activeConnection.token)
        //   setUsername(activeConnection.username)
        //   setPassword(activeConnection.password)
        //   setProxy(activeConnection.Proxy || activeConnection.proxy)
        //   setRateLimit(activeConnection.rateLimitPerHour)
        //   break
        // case Providers.TAPD:
        //   setUsername(activeConnection.username)
        //   setPassword(activeConnection.password)
        //   setProxy(activeConnection.Proxy || activeConnection.proxy)
        //   setRateLimit(activeConnection.rateLimitPerHour)
        //   break
      }
      // ToastNotification.clear()
      // ToastNotification.show({ message: `Fetched settings for ${activeConnection.name}.`, intent: 'success', icon: 'small-tick' })
      console.log('>> FETCHED CONNECTION FOR MODIFY', activeConnection)
    }
  }, [activeConnection, provider?.id])

  useEffect(() => {
    if (saveComplete && saveComplete.ID) {
      console.log('>>> CONNECTION MANAGER - SAVE COMPLETE EFFECT RUNNING...')
      setActiveConnection((ac) => {
        return {
          ...ac,
          ...saveComplete,
        }
      })
    }
  }, [saveComplete])

  useEffect(() => {
    console.log(
      '>> CONNECTION MANAGER - SELECTING ACTIVE PROVIDER...',
      provider
    )
    if (provider && provider?.id) {
      // console.log(activeProvider)
    }
  }, [provider])

  useEffect(() => {
    if (connectionId !== null && connectionId !== undefined) {
      console.log('>>>> CONFIGURING CONNECTION ID ... ', connectionId)
      fetchConnection()
    }
  }, [connectionId, fetchConnection])

  useEffect(() => {
    console.log('>> TESTED CONNECTION RESULTS...', testedConnections)
  }, [testedConnections])

  return {
    // activeConnection,
    // fetchConnection,
    fetchAllConnections,
    // fetchDomainLayerRepositories,
    testAllConnections,
    testConnection,
    saveConnection,
    deleteConnection,
    // runCollection,
    isSaving,
    isTesting,
    isFetching,
    isDeleting,
    errors,
    showError,
    testStatus,
    // provider,
    // setActiveConnection,
    setProvider,

    editingConnection,
    setConnectionColumn,
    initialTokenStore,
    setInitialTokenStore,

    // setIsSaving,
    // setIsTesting,
    // setIsFetching,
    // setErrors,
    // setShowError,
    // setTestStatus,
    // setTestResponse,
    // setAllTestResponses,
    setSaveComplete,
    allConnections,
    allProviderConnections,
    // domainRepositories,
    testedConnections,
    // connectionCount,
    connectionLimitReached,
    // connectionTestPayload,
    Providers,
    // saveComplete,
    deleteComplete,
    // getConnectionName,
    // clearConnection,
    testResponse, // ???
    allTestResponses // ???
  }
}

export default useConnectionManager
