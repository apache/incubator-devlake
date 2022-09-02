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
import { useState, useEffect, useCallback, useMemo } from 'react'
import { useHistory } from 'react-router-dom'
import { ToastNotification } from '@/components/Toast'
import { DEVLAKE_ENDPOINT } from '@/utils/config'
import request from '@/utils/request'
import Connection from '@/models/Connection'
import ProviderListConnection from '@/models/ProviderListConnection'
import {
  Providers,
  ProviderConnectionLimits,
  ConnectionStatus,
  ConnectionStatusLabels
} from '@/data/Providers'

import useNetworkOfflineMode from '@/hooks/useNetworkOfflineMode'

function useConnectionManager (
  {
    activeProvider,
    connectionId,
  },
  updateMode = false
) {
  const history = useHistory()
  const { handleOfflineMode } = useNetworkOfflineMode()

  const [provider, setProvider] = useState(activeProvider)
  const [name, setName] = useState()
  // @todo: refactor to endpoint and setEndpoint
  const [endpointUrl, setEndpointUrl] = useState()
  const [proxy, setProxy] = useState()
  const [rateLimitPerHour, setRateLimitPerHour] = useState(0)
  const [token, setToken] = useState()
  const defaultTokenStore = useMemo(() => ({ 0: '', 1: '', 2: '' }), [])
  const [initialTokenStore, setInitialTokenStore] = useState(defaultTokenStore)
  const [username, setUsername] = useState()
  const [password, setPassword] = useState()

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
  const [sourceLimits, setConnectionLimits] = useState(ProviderConnectionLimits)

  const [activeConnection, setActiveConnection] = useState(new Connection())
  const [allConnections, setAllConnections] = useState([])
  const [allProviderConnections, setAllProviderConnections] = useState([])
  const [domainRepositories, setDomainRepositories] = useState([])
  const [testedConnections, setTestedConnections] = useState([])
  const connectionCount = useMemo(() => allConnections.length, [allConnections])
  const connectionLimitReached = useMemo(() =>
    sourceLimits[provider?.id] && connectionCount >= sourceLimits[provider?.id],
  [provider?.id, sourceLimits, connectionCount],
  )
  const [connectionsList, setConnectionsList] = useState([])

  const [saveComplete, setSaveComplete] = useState(false)
  const [deleteComplete, setDeleteComplete] = useState(false)
  const connectionTestPayload = useMemo(() => ({
    endpoint: endpointUrl,
    username,
    password,
    token,
    proxy
  }), [endpointUrl, password, proxy, token, username])
  const connectionSavePayload = useMemo(() => ({
    name,
    endpoint: endpointUrl,
    username,
    password,
    token,
    proxy,
    rateLimitPerHour,
  }), [name, endpointUrl, username, password, token, proxy, rateLimitPerHour])

  const testConnection = useCallback(
    (
      notify = true,
      manualPayload = {},
      onSuccess = () => {},
      onFail = () => {}
    ) => {
      setIsTesting(true)
      setShowError(false)
      // ToastNotification.clear()
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

  const notifyConnectionSaveSuccess = useCallback((message = 'Connection saved successfully.') => {
    ToastNotification.show({
      message: message,
      intent: 'success',
      icon: 'small-tick',
    })
  }, [])

  const notifyConnectionSaveFailure = useCallback((message = 'Connection failed to save, please try again.') => {
    ToastNotification.show({
      message: message,
      intent: 'danger',
      icon: 'error',
    })
  }, [])

  const saveConnection = useCallback((configurationSettings = {}) => {
    setIsSaving(true)

    let saveResponse = {
      success: false,
      connection: {
        ...connectionSavePayload,
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
        setIsSaving(false)
        setSaveComplete(saveResponse.success ? s.data : false)
        if (!saveResponse.success) { notifyConnectionSaveFailure(s.data || s.message) }
      } catch (e) {
        saveResponse.errors.push(e.message)
        setErrors(saveResponse.errors)
        setIsSaving(false)
        setSaveComplete(false)
        notifyConnectionSaveFailure(e.message)
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
            activeConnection.id
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
        setIsSaving(false)
        setSaveComplete(saveResponse.success ? s.data : false)
        if (!saveResponse.success) { notifyConnectionSaveFailure(s.data || s.message) }
      } catch (e) {
        saveResponse.errors.push(e.message)
        setErrors(saveResponse.errors)
        setIsSaving(false)
        setSaveComplete(false)
        notifyConnectionSaveFailure(e.message)
        console.log('>> CONFIGURATION FAILED TO UPDATE', configPayload, e)
      }
    }

    if (updateMode && activeConnection?.id !== null) {
      modifyConfiguration(connectionSavePayload)
    } else {
      saveConfiguration(connectionSavePayload)
    }
  }, [
    activeConnection?.id,
    connectionSavePayload,
    fetchConnection,
    notifyConnectionSaveFailure,
    provider?.id,
    updateMode
  ])

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
          setActiveConnection(new Connection({
            ...connectionData
          }))
          setTimeout(() => {
            setIsFetching(false)
          }, 500)
        }
        fetch()
      } catch (e) {
        setIsFetching(false)
        setActiveConnection(new Connection())
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
        // ToastNotification.clear()
        console.log('>> FETCHING ALL CONNECTION SOURCES')
        let c = null
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
            .map((providerResponse) => [].concat(providerResponse.data || []).map(c => new Connection({
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
          c = await request.get(
            `${DEVLAKE_ENDPOINT}/plugins/${provider.id}/connections`
          )
          console.log('>> RAW ALL CONNECTIONS DATA FROM API...', c?.data)
          const providerConnections = []
            .concat(Array.isArray(c?.data) ? c?.data : [])
            .map((conn, idx) =>
              new Connection({
                ...conn,
                status: ConnectionStatus.OFFLINE,
                id: conn.id,
                name: conn.name,
                endpoint: conn.endpoint,
                errors: [],
              })
            )
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
            connection.id
          }`
        )
        console.log('>> CONNECTION DELETED...', d)
        setIsDeleting(false)
        setDeleteComplete({
          provider: activeProvider,
          connection: d.data,
        })
      } catch (e) {
        setIsDeleting(false)
        setDeleteComplete(false)
        setErrors([e.message])
        console.log('>> FAILED TO DELETE CONNECTION', e)
      }
    },
    [provider?.id, activeProvider]
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
              ...testedConnections.filter((oC) => oC.id !== c.id),
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
    setName('')
    setEndpointUrl('')
    setUsername('')
    setPassword('')
    setToken('')
    setInitialTokenStore(defaultTokenStore)
    setProxy('')
    setRateLimitPerHour(0)
  }, [defaultTokenStore])

  useEffect(() => {
    if (activeConnection && activeConnection.id !== null) {
      setName(activeConnection?.name)
      setEndpointUrl(activeConnection?.endpoint)
      setRateLimitPerHour(activeConnection?.rateLimitPerHour)
      setProxy(activeConnection?.proxy)
      setUsername(activeConnection?.username)
      setPassword(activeConnection?.password)
      setToken(activeConnection?.token)
      setInitialTokenStore(activeConnection?.token
        ? activeConnection?.token?.split(',')?.reduce((tS, cT, id) => ({ ...tS, [id]: cT }), {})
        : defaultTokenStore
      )
      console.log('>> FETCHED CONNECTION FOR MODIFY', activeConnection)
    }
  }, [activeConnection, defaultTokenStore])

  useEffect(() => {
    if (saveComplete?.id) {
      console.log('>>> CONNECTION MANAGER - SAVE COMPLETE EFFECT RUNNING...', saveComplete)
      setActiveConnection((ac) =>
        new Connection({
          ...ac,
          ...saveComplete
        })
      )
      if (!updateMode) {
        history.replace(`/integrations/${provider.id}`)
        notifyConnectionSaveSuccess()
      } else {
        notifyConnectionSaveSuccess()
      }
    }
  }, [saveComplete, updateMode, history, provider?.id, notifyConnectionSaveSuccess])

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

  useEffect(() => {
    console.log('>> CONNECTION MANAGER, ACTIVE PROVIDER CHANGED ====>', activeProvider)
    setProvider(activeProvider)
  }, [activeProvider])

  useEffect(() => {
    console.log('>>> ALL DATA PROVIDER CONNECTIONS...', allProviderConnections)
    setConnectionsList(
      allProviderConnections?.map((c, cIdx) => new ProviderListConnection({
        ...c,
        id: cIdx,
        key: cIdx,
        connectionId: c.id,
        name: c.name,
        title: c.name,
        value: c.id,
        status:
          ConnectionStatusLabels[c.status] ||
          ConnectionStatusLabels[ConnectionStatus.OFFLINE],
        statusResponse: null,
        provider: c.provider,
        providerId: c.provider,
        plugin: c.provider,
      }))
    )
  }, [allProviderConnections])

  return {
    activeConnection,
    fetchConnection,
    fetchAllConnections,
    fetchDomainLayerRepositories,
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
    rateLimitPerHour,
    username,
    password,
    token,
    initialTokenStore,
    provider,
    setActiveConnection,
    setProvider,
    setName,
    setEndpointUrl,
    setProxy,
    setRateLimitPerHour,
    setToken,
    setInitialTokenStore,
    setUsername,
    setPassword,
    setIsSaving,
    setIsTesting,
    setIsFetching,
    setErrors,
    setShowError,
    setTestStatus,
    setTestResponse,
    setAllTestResponses,
    setConnectionLimits,
    setConnectionsList,
    setSaveComplete,
    allConnections,
    allProviderConnections,
    connectionsList,
    domainRepositories,
    testedConnections,
    sourceLimits,
    connectionCount,
    connectionLimitReached,
    connectionTestPayload,
    Providers,
    saveComplete,
    deleteComplete,
    getConnectionName,
    clearConnection,
    testResponse,
    allTestResponses
  }
}

export default useConnectionManager
