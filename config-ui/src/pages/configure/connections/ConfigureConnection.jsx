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
import useConnectionManager from '@/hooks/useConnectionManager'
import useSettingsManager from '@/hooks/useSettingsManager'

// import { ToastNotification } from '@/components/Toast'

import JiraSettings from '@/pages/configure/settings/jira'
import GitlabSettings from '@/pages/configure/settings/gitlab'
import JenkinsSettings from '@/pages/configure/settings/jenkins'

import { integrationsData } from '@/pages/configure/mock-data/integrations'
// import { connectionsData } from '@/pages/configure/mock-data/connections'

import '@/styles/integration.scss'
import '@/styles/connections.scss'
import '@/styles/configure.scss'

import '@blueprintjs/popover2/lib/css/blueprint-popover2.css'

export default function ConfigureConnection () {
  const history = useHistory()
  const { providerId, connectionId } = useParams()

  // const [isSaving, setIsSaving] = useState(false)
  // const [errors, setErrors] = useState([])
  // const [showError, setShowError] = useState(false)

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
  const [connections, setConnections] = useState([])

  const [settings, setSettings] = useState({
    JIRA_BASIC_AUTH_ENCODED: null,
    JIRA_ISSUE_EPIC_KEY_FIELD: null,
    JIRA_ISSUE_TYPE_MAPPING: null,
    JIRA_ISSUE_STORYPOINT_COEFFICIENT: null,
    JIRA_ISSUE_STORYPOINT_FIELD: null,
    JIRA_BOARD_GITLAB_PROJECTS: null,
  })

  const {
    fetchConnection,
  } = useConnectionManager({
    activeProvider,
    activeConnection,
    connectionId,
    setActiveConnection,
  })

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

  const renderProviderSettings = (providerId) => {
    let settingsComponent = null
    switch (providerId) {
      case 'jira' :
        settingsComponent = (
          <JiraSettings
            provider={activeProvider}
            connection={activeConnection}
            isSaving={isSaving}
            onSettingsChange={setSettings}
          />
        )
        break
      case 'gitlab' :
        settingsComponent = (
          <GitlabSettings
            provider={activeProvider}
            connection={activeConnection}
            isSaving={isSaving}
            onSettingsChange={setSettings}
          />
        )
        break
      case 'jenkins' :
        settingsComponent = (
          <JenkinsSettings
            provider={activeProvider}
            connection={activeConnection}
            isSaving={isSaving}
            onSettingsChange={setSettings}
          />
        )
        break
    }
    return settingsComponent
  }

  useEffect(() => {
    console.log('>>>> DETECTED PROVIDER ID = ', providerId)
    console.log('>>>> DETECTED CONNECTION ID = ', connectionId)
    if (connectionId && providerId) {
      setActiveProvider(integrations.find(p => p.id === providerId))
      // !WARNING! DO NOT ADD fetchConnection TO DEPENDENCIES ARRAY!
      // @todo FIXME: Fix Hook Circular-loop Behavior inside effect when added to dependencies
      fetchConnection()
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

  }, [settings])

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
                  href: `/connections/configure/${activeProvider.id}/${activeConnection && activeConnection.id}`,
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
                        onClick={saveSettings}
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
