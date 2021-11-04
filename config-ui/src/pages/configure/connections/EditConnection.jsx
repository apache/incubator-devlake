import React, { useEffect, useState } from 'react'
import {
  useParams,
  Link,
  useHistory
} from 'react-router-dom'
import {
  Icon,
} from '@blueprintjs/core'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
// import { ToastNotification } from '@/components/Toast'
import ConnectionForm from '@/pages/configure/connections/ConnectionForm'
import { integrationsData } from '@/pages/configure/mock-data/integrations'
import { NullConnection } from '@/data/NullConnection'

import useConnectionManager from '@/hooks/useConnectionManager'

import '@/styles/integration.scss'
import '@/styles/connections.scss'
import '@blueprintjs/popover2/lib/css/blueprint-popover2.css'

export default function EditConnection () {
  const history = useHistory()
  const { providerId, connectionId } = useParams()

  // const [name, setName] = useState()
  // const [endpointUrl, setEndpointUrl] = useState()
  // const [token, setToken] = useState()
  // const [username, setUsername] = useState()
  // const [password, setPassword] = useState()

  const [integrations, setIntegrations] = useState(integrationsData)
  const [activeProvider, setActiveProvider] = useState(integrations[0])

  const [activeConnection, setActiveConnection] = useState(NullConnection)

  const {
    testConnection,
    saveConnection,
    fetchConnection,
    name,
    endpointUrl,
    username,
    password,
    token,
    errors,
    isSaving,
    isTesting,
    showError,
    testStatus,
    setName,
    setEndpointUrl,
    setUsername,
    setPassword,
    setToken
  } = useConnectionManager({
    activeProvider,
    activeConnection,
    connectionId,
    setActiveConnection,
    // name,
    // endpointUrl,
    // token,
    // username,
    // password,
  }, true)

  const cancel = () => {
    history.push(`/integrations/${activeProvider.id}`)
  }

  useEffect(() => {
    if (activeProvider && connectionId) {
      fetchConnection()
    }
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
        setToken(activeConnection.basicAuthEncoded || activeConnection.Auth)
        break
    }
  }, [activeConnection, activeProvider.id])

  useEffect(() => {
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
                  href: `/connections/edit/${activeProvider.id}/${activeConnection.ID}`,
                  icon: false,
                  text: 'Edit Connection',
                  current: true
                }
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
                    Edit <strong style={{ fontWeight: 900 }}>{activeProvider.name}</strong> Connection
                  </h1>
                  <p className='description'>Manage the connection source for this provider.</p>
                </div>
              </div>
              <div className='editConnection' style={{ display: 'flex' }}>
                <ConnectionForm
                  activeProvider={activeProvider}
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
