import React, { useEffect, useState } from 'react'
import axios from 'axios'
import {
  // BrowserRouter as Router,
  // Switch,
  // Route,
  useParams,
  Link,
  useHistory
} from 'react-router-dom'
import {
  Button, mergeRefs, Card, Elevation, Colors,
  FormGroup, InputGroup, Tooltip, Label,
  Position,
  Alignment,
  Icon,
  Toaster,
  ToasterPosition,
  IToasterProps,
  IToastProps
} from '@blueprintjs/core'
// import { Column, Table } from '@blueprintjs/table'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import Content from '@/components/Content'
import JiraSettings from '@/pages/configure/settings/jira'
import GitlabSettings from '@/pages/configure/settings/gitlab'
import JenkinsSettings from '@/pages/configure/settings/jenkins'

import { ToastNotification } from '@/components/Toast'

import { SERVER_HOST, DEVLAKE_ENDPOINT } from '@/utils/config'

import { ReactComponent as GitlabProvider } from '@/images/integrations/gitlab.svg'
import { ReactComponent as JenkinsProvider } from '@/images/integrations/jenkins.svg'
import { ReactComponent as JiraProvider } from '@/images/integrations/jira.svg'

import '@/styles/integration.scss'
import '@/styles/connections.scss'

import '@blueprintjs/popover2/lib/css/blueprint-popover2.css'

export default function ConfigureConnection () {
  const history = useHistory()
  const { providerId, connectionId } = useParams()

  const [isSaving, setIsSaving] = useState(false)
  const [errors, setErrors] = useState([])
  const [showError, setShowError] = useState(false)

  const [integrations, setIntegrations] = useState([
    {
      id: 'gitlab',
      name: 'Gitlab',
      icon: <GitlabProvider className='providerIconSvg' width='30' height='30' style={{ float: 'left', marginTop: '5px' }} />
    },
    {
      id: 'jenkins',
      name: 'Jenkins',
      icon: <JenkinsProvider className='providerIconSvg' width='30' height='30' style={{ float: 'left', marginTop: '5px' }} />
    },
    {
      id: 'jira',
      name: 'JIRA',
      icon: <JiraProvider className='providerIconSvg' width='30' height='30' style={{ float: 'left', marginTop: '5px' }} />
    },
  ])

  const [activeProvider, setActiveProvider] = useState(integrations[0])
  const [activeConnection, setActiveConnection] = useState()

  const [connections, setConnections] = useState([
    {
      id: 0,
      name: 'JIRA Development Server',
      endpoint: 'https://jira-test-a345vf.merico.dev',
      status: 1,
      errors: []
    },
    {
      id: 1,
      name: 'JIRA Staging Server',
      endpoint: 'https://jira-staging-93xt5a.merico.dev',
      status: 2,
      errors: []
    },
    {
      id: 2,
      name: 'JIRA Production Server',
      endpoint: 'https://jira-prod-z51gox.merico.dev',
      status: 0,
      errors: []
    },
    {
      id: 3,
      name: 'JIRA Demo Instance 591',
      endpoint: 'https://jira-demo-591.merico.dev',
      status: 0,
      errors: []
    },
    {
      id: 4,
      name: 'JIRA Demo Instance 142',
      endpoint: 'https://jira-demo-142.merico.dev',
      status: 0,
      errors: []
    },
    {
      id: 5,
      name: 'JIRA Demo Instance 111',
      endpoint: 'https://jira-demo-111.merico.dev',
      status: 0,
      errors: []
    },
    {
      id: 6,
      name: 'JIRA Demo Instance 784',
      endpoint: 'https://jira-demo-784.merico.dev',
      status: 3,
      errors: []
    },
  ])

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
            <div style={{ width: '100%' }}>
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
