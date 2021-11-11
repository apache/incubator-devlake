import React, { useEffect, useState, useCallback } from 'react'
import {
  useParams,
  Link,
  useHistory
} from 'react-router-dom'
import {
  Button,
  Icon,
  Intent,
  Card,
  Elevation,
} from '@blueprintjs/core'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
import useConnectionManager from '@/hooks/useConnectionManager'
import useSettingsManager from '@/hooks/useSettingsManager'
import ConnectionForm from '@/pages/configure/connections/ConnectionForm'

import { integrationsData } from '@/data/integrations'
// import { NullConnection } from '@/data/NullConnection'
import { NullSettings } from '@/data/NullSettings'
import {
  Providers,
  ProviderSourceLimits,
  ProviderFormLabels,
  ProviderFormPlaceholders
} from '@/data/Providers'

import '@/styles/integration.scss'
import '@/styles/connections.scss'
import '@/styles/configure.scss'

import '@blueprintjs/popover2/lib/css/blueprint-popover2.css'

export default function ConfigureConnection () {
  const history = useHistory()
  const { providerId, connectionId } = useParams()

  const [integrations, setIntegrations] = useState(integrationsData)
  const [activeProvider, setActiveProvider] = useState(integrations.find(p => p.id === providerId))
  // const [activeConnection, setActiveConnection] = useState(NullConnection)
  const [connections, setConnections] = useState([])
  const [showConnectionSettings, setShowConnectionSettings] = useState(true)

  const [settings, setSettings] = useState(NullSettings)

  const {
    testConnection,
    saveConnection,
    fetchConnection,
    activeConnection,
    name,
    endpointUrl,
    username,
    password,
    token,
    errors,
    testStatus,
    isSaving: isSavingConnection,
    isTesting: isTestingConnection,
    setName,
    setEndpointUrl,
    setUsername,
    setPassword,
    setToken,
    saveComplete: saveConnectionComplete,
    showError: showConnectionError
  } = useConnectionManager({
    activeProvider,
    // activeConnection,
    connectionId,
    // setActiveConnection,
  }, true)

  const {
    saveSettings,
    // errors: settingsErrors,
    isSaving,
    isTesting,
    showError,
  } = useSettingsManager({
    activeProvider,
    activeConnection,
    settings
  })

  const cancel = () => {
    history.push(`/integrations/${activeProvider.id}`)
  }

  const renderProviderSettings = (providerId, activeProvider) => {
    console.log('>>> RENDERING PROVIDER SETTINGS...')
    let settingsComponent = null
    if (activeProvider && activeProvider.settings) {
      settingsComponent = activeProvider.settings({
        activeProvider,
        activeConnection,
        isSaving,
        setSettings
      })
    } else {
      // @todo create & display "fallback/empty settings" view
      console.log('>> WARNING: NO PROVIDER SETTINGS RENDERED, PROVIDER = ', activeProvider)
    }
    return settingsComponent
  }

  useEffect(() => {
    console.log('>>>> DETECTED PROVIDER ID = ', providerId)
    console.log('>>>> DETECTED CONNECTION ID = ', connectionId)
    if (connectionId && providerId) {
      setActiveProvider(integrations.find(p => p.id === providerId))
    } else {
      console.log('NO PARAMS!')
    }
  }, [connectionId, providerId, integrations])

  useEffect(() => {

  }, [settings])

  useEffect(() => {

  }, [connections])

  // useEffect(() => {
  //   // CONNECTION SAVED!
  // }, [saveConnectionComplete])

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
                  href: `/connections/configure/${activeProvider.id}/${activeConnection && activeConnection.ID}`,
                  icon: false,
                  text: `${activeConnection ? activeConnection.name : 'Configure'} Settings`,
                  current: true
                }
              ]}
            />
            <div className='configureConnection' style={{ width: '100%' }}>
              <Link style={{ float: 'right', marginLeft: '10px', color: '#777777' }} to={`/integrations/${activeProvider.id}`}>
                <Icon icon='fast-backward' size={16} /> Connection List
              </Link>
              <div style={{ display: 'flex', justifyContent: 'flex-start' }}>
                <div>
                  <span style={{ marginRight: '10px' }}>{activeProvider.icon}</span>
                </div>
                <div style={{ justifyContent: 'flex-start' }}>
                  <h1 style={{ margin: 0 }}>Manage <strong style={{ fontWeight: 900 }}>{activeProvider.name}</strong> Settings </h1>
                  {activeConnection && (
                    <>
                      <h2 style={{ margin: 0 }}>{activeConnection.name}</h2>
                      <p className='description'>Manage settings and options for this connection.</p>
                    </>
                  )}
                </div>
              </div>
              {activeProvider && activeConnection && (
                <>
                  <Card
                    interactive={false}
                    elevation={Elevation.ZERO}
                    style={{ backgroundColor: '#f8f8f8', width: '100%', marginBottom: '20px' }}
                  >
                    <Button
                      type='button'
                      icon={showConnectionSettings ? 'eye-on' : 'eye-off'}
                      intent={showConnectionSettings ? Intent.PRIMARY : Intent.DISABLED}
                      style={{ margin: '2px', float: 'right' }}
                      onClick={() => setShowConnectionSettings(!showConnectionSettings)}
                      minimal
                      small
                    />
                    {showConnectionSettings
                      ? (
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
                            isSaving={isSavingConnection}
                            isTesting={isTestingConnection}
                            testStatus={testStatus}
                            errors={errors}
                            showError={showConnectionError}
                            authType={activeProvider.id === Providers.JENKINS ? 'plain' : 'token'}
                            showLimitWarning={false}
                            sourceLimits={ProviderSourceLimits}
                            labels={ProviderFormLabels[activeProvider.id]}
                            placeholders={ProviderFormPlaceholders[activeProvider.id]}
                          />
                        </div>
                        )
                      : (
                        <>
                          <h2 style={{ margin: 0 }}>Configure Connection</h2>
                          <p className='description' style={{ margin: 0 }}>
                            ( Click the <strong>Visibility</strong> icon to your right to edit connection )
                          </p>
                        </>
                        )}
                  </Card>
                  <div style={{ marginTop: '30px' }}>
                    {renderProviderSettings(providerId, activeProvider)}
                  </div>
                  <div className='form-actions-block' style={{ display: 'flex', marginTop: '60px', justifyContent: 'space-between' }}>
                    <div>
                      {/* <Button
                        icon={getConnectionStatusIcon()}
                        text='Test Connection'
                        onClick={testConnection}
                        loading={isTesting}
                        disabled={isTesting || isSaving}
                      /> */}
                    </div>
                    <div>
                      <Button icon='remove' text='Cancel' onClick={cancel} disabled={isSaving} />
                      <Button
                        icon='cloud-upload'
                        intent={Intent.PRIMARY}
                        text='Save Settings'
                        loading={isSaving}
                        disabled={isSaving || providerId === 'jenkins'}
                        onClick={saveSettings}
                        style={{ marginLeft: '10px' }}
                      />
                    </div>
                  </div>

                </>
              )}
            </div>
          </main>
        </Content>
      </div>
    </>
  )
}
