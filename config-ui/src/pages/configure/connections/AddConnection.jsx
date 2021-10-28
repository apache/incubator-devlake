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
// import { connectionsData } from '@/pages/configure/mock-data/connections'
import { SERVER_HOST, DEVLAKE_ENDPOINT } from '@/utils/config'

import useConnectionManager from '@/hooks/useConnectionManager'

import '@/styles/integration.scss'
import '@/styles/connections.scss'
import '@blueprintjs/popover2/lib/css/blueprint-popover2.css'

export default function AddConnection () {
  const history = useHistory()
  const { providerId } = useParams()

  const [name, setName] = useState()
  const [endpointUrl, setEndpointUrl] = useState()
  const [token, setToken] = useState()
  const [username, setUsername] = useState()
  const [password, setPassword] = useState()

  const [integrations, setIntegrations] = useState(integrationsData)
  const [activeProvider, setActiveProvider] = useState(integrations[0])

  const {
    testConnection, saveConnection,
    errors, // showErrors,
    isSaving, // setIsSaving,
    isTesting, // setIsTesting,
    showError, // setShowError,
    testStatus, // setTestStatus
  } = useConnectionManager({
    activeProvider,
    name,
    endpointUrl,
    token,
    username,
    password,
  })

  const cancel = () => {
    history.push(`/integrations/${activeProvider.id}`)
  }

  // const resetForm = () => {
  //   setName(null)
  //   setEndpointUrl(null)
  //   setToken(null)
  //   setUsername(null)
  //   setPassword(null)
  // }

  useEffect(() => {
    // Selected Provider
    console.log(activeProvider)
  }, [activeProvider])

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
                { href: `/connections/add/${activeProvider.id}`, icon: false, text: 'Add Connection', current: true }
              ]}
            />
            <div style={{ width: '100%' }}>
              <Link style={{ float: 'right', marginLeft: '10px', color: '#777777' }} to={`/integrations/${activeProvider.id}`}>
                <Icon icon='fast-backward' size={16} /> Go Back
              </Link>
              <div style={{ display: 'flex' }}>
                <div>
                  <span style={{ marginRight: '10px' }}>{activeProvider.icon}</span>
                </div>
                <div>
                  <h1 style={{ margin: 0 }}>
                    {activeProvider.name} Add Connection
                  </h1>
                  <p className='description'>Create a new connection source for this provider.</p>
                </div>
              </div>
              <div className='addConnection' style={{ display: 'flex' }}>
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
