import React, { useEffect, useState } from 'react'
import {
  useParams,
  Link,
  useHistory
} from 'react-router-dom'
import {
  Colors,
  Icon,
} from '@blueprintjs/core'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
import ConnectionForm from '@/pages/configure/connections/ConnectionForm'
import { integrationsData } from '@/data/integrations'
import {
  Providers,
  ProviderLabels,
  ProviderSourceLimits,
  ProviderFormLabels,
  ProviderFormPlaceholders
} from '@/data/Providers'

import useConnectionManager from '@/hooks/useConnectionManager'
import useConnectionValidation from '@/hooks/useConnectionValidation'

import '@/styles/integration.scss'
import '@/styles/connections.scss'

export default function AddConnection () {
  const history = useHistory()
  const { providerId } = useParams()

  const [integrations, setIntegrations] = useState(integrationsData)
  const [activeProvider, setActiveProvider] = useState(integrations.find(p => p.id === providerId))

  const {
    testConnection, saveConnection,
    errors,
    isSaving,
    isTesting,
    showError,
    testStatus,
    name,
    endpointUrl,
    token,
    username,
    password,
    setName,
    setEndpointUrl,
    setUsername,
    setPassword,
    setToken,
    fetchAllConnections,
    connectionLimitReached,
    // Providers
  } = useConnectionManager({
    activeProvider,
  })

  const {
    validate,
    errors: validationErrors,
    isValid: isValidForm
  } = useConnectionValidation({
    activeProvider,
    name,
    endpointUrl,
    token,
    username,
    password
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
    if (activeProvider && activeProvider.id) {
      fetchAllConnections()
      switch (activeProvider.id) {
        case Providers.GITHUB:
          setName(ProviderLabels.GITHUB)
          break
        case Providers.GITLAB:
          setName(ProviderLabels.GITLAB)
          break
        case Providers.JENKINS:
          setName(ProviderLabels.JENKINS)
          break
        case Providers.JIRA:
        default:
          setName('')
          break
      }
    }
  }, [activeProvider.id])

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
                  isLocked={connectionLimitReached}
                  isValid={isValidForm}
                  activeProvider={activeProvider}
                  name={name}
                  endpointUrl={endpointUrl}
                  token={token}
                  username={username}
                  password={password}
                  onSave={saveConnection}
                  onTest={testConnection}
                  onCancel={cancel}
                  onValidate={validate}
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
                  authType={activeProvider.id === Providers.JENKINS ? 'plain' : 'token'}
                  sourceLimits={ProviderSourceLimits}
                  labels={ProviderFormLabels[activeProvider.id]}
                  placeholders={ProviderFormPlaceholders[activeProvider.id]}
                />
              </div>
              {validationErrors.length > 0 && (
                <div className='validation-errors'>
                  <p style={{ margin: '5px 0 5px 0', textAlign: 'right' }}>
                    <Icon icon='warning-sign' size={13} color={Colors.ORANGE4} style={{ marginRight: '6px', marginBottom: '2px' }} />
                    {validationErrors[0]}
                  </p>
                </div>
              )}
            </div>
          </main>
        </Content>
      </div>
    </>
  )
}
