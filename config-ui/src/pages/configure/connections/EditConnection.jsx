import React, { useEffect, useState } from 'react'
import axios from 'axios'
import {
  useParams,
  Link,
  useHistory
} from 'react-router-dom'
import {
  Button, Card, Elevation, Colors,
  FormGroup, InputGroup, Tooltip, Label,
  Icon,
} from '@blueprintjs/core'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
import { ToastNotification } from '@/components/Toast'
import ConnectionForm from '@/pages/configure/connections/ConnectionForm'
import { integrationsData } from '@/pages/configure/mock-data/integrations'
import { connectionsData } from '@/pages/configure/mock-data/connections'
import { SERVER_HOST, DEVLAKE_ENDPOINT } from '@/utils/config'
import request from '@/utils/request'

import useConnectionManager from '@/hooks/useConnectionManager'

import '@/styles/integration.scss'
import '@/styles/connections.scss'
import '@blueprintjs/popover2/lib/css/blueprint-popover2.css'

export default function EditConnection () {
  const history = useHistory()
  const { providerId, connectionId } = useParams()

  const [name, setName] = useState()
  const [endpointUrl, setEndpointUrl] = useState()
  const [token, setToken] = useState()
  const [username, setUsername] = useState()
  const [password, setPassword] = useState()

  // const [isSaving, setIsSaving] = useState(false)
  // const [isTesting, setIsTesting] = useState(false)
  // const [errors, setErrors] = useState([])
  // const [showError, setShowError] = useState(false)
  // const [testStatus, setTestStatus] = useState(0) //  0=Pending, 1=Success, 2=Failed

  const [integrations, setIntegrations] = useState(integrationsData)
  const [activeProvider, setActiveProvider] = useState(integrations[0])

  const [activeConnection, setActiveConnection] = useState({
    id: null,
    name: null,
    endpoint: null,
    token: null,
    username: null,
    password: null,
  })

  const {
    testConnection, saveConnection,
    errors, // showErrors,
    isSaving, // setIsSaving,
    isTesting, // setIsTesting,
    showError, // setShowError,
    testStatus, // setTestStatus
  } = useConnectionManager({
    activeProvider,
    activeConnection,
    name,
    endpointUrl,
    token,
    username,
    password,
  }, true)

  const cancel = () => {
    history.push(`/integrations/${activeProvider.id}`)
  }

  useEffect(() => {
    // Selected Provider
    console.log(activeProvider)

    const fetchConnection = async () => {
      try {
        const connectionResponse = await request.get(`${DEVLAKE_ENDPOINT}/plugins/${activeProvider.id}/sources/${connectionId}`)
        const connectionData = connectionResponse.data.data
        setActiveConnection(connectionData)
        console.log('>> FETCHED CONNECTION FOR MODIFY', connectionResponse)
      } catch (e) {
        console.log('>> FAILED TO FETCH CONNECTION', e)
      }
    }
    fetchConnection()
  }, [activeProvider, providerId, connectionId])

  useEffect(() => {
    setName(activeConnection.name)
    setEndpointUrl(activeConnection.endpoint)
    switch (activeProvider.id) {
      case 'jenkins':
        setUsername(activeConnection.username)
        setPassword(activeConnection.password)
        break
      case 'gitlab':
      case 'jira':
        setToken(activeConnection.token)
        break
    }
  }, [activeConnection, activeProvider.id])

  useEffect(() => {
    fetch(`${SERVER_HOST}/api/getenv`)
      .then(response => response.json())
      .then(env => {
        // setDbUrl(env.DB_URL)
        // setPort(env.PORT)
        // setMode(env.MODE)
      })
    console.log('>>>> DETECTED PROVIDER = ', providerId)
    setActiveProvider(integrations.find(p => p.id === providerId))
  }, [])

  return (
    <>
      <div className='container'>
        <Nav />
        <Sidebar />
        <Content>
          <main className='main'>
            <AppCrumbs
              items={[
                { href: '/', icon: false, text: 'Dashboard' },
                { href: '/integrations', icon: false, text: 'Integrations' },
                { href: `/integrations/${activeProvider.id}`, icon: false, text: `${activeProvider.name}` },
                {
                  href: `/connections/edit/${activeProvider.id}/${activeConnection.id}`,
                  icon: false,
                  text: 'Edit Connection',
                  current: true
                }
              ]}
            />
            <div style={{ width: '100%' }}>
              <Link style={{ float: 'right', marginLeft: '10px', color: '#777777' }} to={`/integrations/${activeProvider.id}`}>
                <Icon icon='fast-backward' size={16} /> &nbsp; Go Back
              </Link>
              <div style={{ display: 'flex' }}>
                <div>
                  <span style={{ marginRight: '10px' }}>{activeProvider.icon}</span>
                </div>
                <div>
                  <h1 style={{ margin: 0 }}>
                    Edit {activeProvider.name} Connection
                  </h1>
                  <p className='description'>Manage the connection source for this provider.</p>
                </div>
              </div>
              <div className='editConnection' style={{ display: 'flex' }}>
                <ConnectionForm
                  name={name}
                  endpointUrl={endpointUrl}
                  token={token}
                  username={username}
                  password={password}
                  onSave={saveConnection}
                  onTest={testConnection}
                  onCancel={cancel}
                  onNameChange={setName}
                  onEndpointChange={setEndpointUrl}
                  onTokenChange={setToken}
                  onUsernameChange={setUsername}
                  onPasswordChange={setPassword}
                  isSaving={isSaving}
                  isTesting={isTesting}
                  testStatus={testStatus}
                  errors={errors}
                  showError={showError}
                  authType={activeProvider.id === 'jenkins' ? 'plain' : 'token'}
                />
              </div>
            </div>
          </main>
        </Content>
      </div>
    </>
  )
}
