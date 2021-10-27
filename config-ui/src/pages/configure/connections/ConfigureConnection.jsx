import React, { useEffect, useState } from 'react'
import {
  useParams,
  Link,
  useHistory
} from 'react-router-dom'
import {
  Button,
  Icon,
} from '@blueprintjs/core'
// import { SERVER_HOST, DEVLAKE_ENDPOINT } from '@/utils/config'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'

import { ToastNotification } from '@/components/Toast'

import JiraSettings from '@/pages/configure/settings/jira'
import GitlabSettings from '@/pages/configure/settings/gitlab'
import JenkinsSettings from '@/pages/configure/settings/jenkins'

import { integrationsData } from '@/pages/configure/mock-data/integrations'
import { connectionsData } from '@/pages/configure/mock-data/connections'

import '@/styles/integration.scss'
import '@/styles/connections.scss'
import '@/styles/configure.scss'

import '@blueprintjs/popover2/lib/css/blueprint-popover2.css'

export default function ConfigureConnection () {
  const history = useHistory()
  const { providerId, connectionId } = useParams()

  const [isSaving, setIsSaving] = useState(false)
  const [errors, setErrors] = useState([])
  const [showError, setShowError] = useState(false)

  const [integrations, setIntegrations] = useState(integrationsData)
  const [activeProvider, setActiveProvider] = useState(integrations[0])
  const [activeConnection, setActiveConnection] = useState()
  const [connections, setConnections] = useState(connectionsData)

  const cancel = () => {
    history.push(`/integrations/${activeProvider.id}`)
  }

  const saveConfiguration = () => {
    setIsSaving(true)
    const saveResponse = {
      success: true,
      connection: {
        ...activeConnection
      },
      errors: []
    }
    setErrors(saveResponse.errors)
    setTimeout(() => {
      if (saveResponse.success && errors.length === 0) {
        ToastNotification.show({ message: 'Configuration saved successfully.', intent: 'success', icon: 'small-tick' })
        setShowError(false)
        setIsSaving(false)
        // resetForm()
        // history.push(`/integrations/${activeProvider.id}`)
      } else {
        ToastNotification.show({ message: 'Unable to save configuration.', intent: 'danger', icon: 'error' })
        setShowError(true)
        setIsSaving(false)
      }
    }, 2000)
  }

  const renderProviderSettings = (providerId) => {
    let settingsComponent = null
    switch (providerId) {
      case 'jira' : settingsComponent = <JiraSettings provider={activeProvider} connection={activeConnection} isSaving={isSaving} />
        break
      case 'gitlab' : settingsComponent = <GitlabSettings provider={activeProvider} connection={activeConnection} isSaving={isSaving} />
        break
      case 'jenkins' : settingsComponent = <JenkinsSettings provider={activeProvider} connection={activeConnection} isSaving={isSaving} />
        break
    }
    return settingsComponent
  }

  useEffect(() => {
    // fetch(`${SERVER_HOST}/api/getenv`)
    //   .then(response => response.json())
    //   .then(env => {
    //     // setDbUrl(env.DB_URL)
    //     // setPort(env.PORT)
    //     // setMode(env.MODE)
    //   })
    console.log('>>>> DETECTED PROVIDER ID = ', providerId)
    console.log('>>>> DETECTED CONNECTION ID = ', connectionId)
    if (connectionId && providerId) {
      setActiveProvider(integrations.find(p => p.id === providerId))
      setActiveConnection(connections.find(c => c.id === 3))
    } else {
      console.log('NO PARAMS!')
    }
  }, [connectionId, providerId, integrations, connections])

  useEffect(() => {
    // Selected Provider
    // console.log('>> active connection', activeConnection)
    console.log('>> active connection', activeConnection)
  }, [activeConnection])

  useEffect(() => {

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
                  href: `/connections/configure/${activeProvider.id}`,
                  icon: false,
                  text: `${activeConnection ? activeConnection.name : 'Configure'} Settings`,
                  current: true
                }
              ]}
            />
            <div className='configureConnection' style={{ width: '100%' }}>
              <Link style={{ float: 'right', marginLeft: '10px', color: '#777777' }} to={`/integrations/${activeProvider.id}`}>
                <Icon icon='fast-backward' size={16} /> &nbsp; Connection List
              </Link>
              <div style={{ display: 'flex', justifyContent: 'flex-start' }}>
                <div>
                  <span style={{ marginRight: '10px' }}>{activeProvider.icon}</span>
                </div>
                <div style={{ justifyContent: 'flex-start' }}>
                  <h1 style={{ margin: 0 }}>Manage {activeProvider.name} Settings </h1>
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
                  {/* <Card interactive={false} elevation={Elevation.TWO} style={{ width: '50%', marginBottom: '20px' }}>
                    <h5>Edit Connection</h5>
                  </Card> */}

                  <div style={{ marginTop: '30px' }}>
                    {renderProviderSettings(providerId)}
                  </div>
                  <div style={{ display: 'flex', marginTop: '60px', justifyContent: 'space-between', maxWidth: '50%' }}>
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
                        icon='cloud-upload' intent='primary' text='Save Settings'
                        loading={isSaving}
                        disabled={isSaving || providerId === 'jenkins'}
                        onClick={saveConfiguration}
                        style={{ backgroundColor: '#E8471C', marginLeft: '10px' }}
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
