import React, { useEffect, useState } from 'react'
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
  Position,
  Alignment,
  Spinner,
  Icon,
} from '@blueprintjs/core'
import { Popover2, Tooltip2 } from '@blueprintjs/popover2'
// import { FormGroup, InputGroup, Button, Tooltip, Position, Label } from '@blueprintjs/core'
// import { Column, Table } from '@blueprintjs/table'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import Content from '@/components/Content'
import { ToastNotification } from '@/components/Toast'

import { SERVER_HOST } from '@/utils/config'

import { ReactComponent as GitlabProvider } from '@/images/integrations/gitlab.svg'
import { ReactComponent as JenkinsProvider } from '@/images/integrations/jenkins.svg'
import { ReactComponent as JiraProvider } from '@/images/integrations/jira.svg'

import { integrationsData } from '@/pages/configure/mock-data/integrations'
import { connectionsData } from '@/pages/configure/mock-data/connections'

import '@/styles/integration.scss'
import '@blueprintjs/popover2/lib/css/blueprint-popover2.css'

export default function ManageIntegration () {
  const history = useHistory()

  const [dbUrl, setDbUrl] = useState()
  const [port, setPort] = useState()
  const [mode, setMode] = useState()

  const [isLoading, setIsLoading] = useState(true)

  const { providerId } = useParams()

  const [errors, setErrors] = useState([])

  const [integrations, setIntegrations] = useState(integrationsData)

  const [connections, setConnections] = useState(connectionsData)

  const [activeProvider, setActiveProvider] = useState(integrations[0])

  const fetchConnections = () => {
    // Fetch Connection List from BE
    setIsLoading(true)
    setTimeout(() => {
      setIsLoading(false)
      ToastNotification.clear()
      ToastNotification.show({ message: 'Loaded all connections.', intent: 'success', icon: 'small-tick' })
    }, 1000)
  }

  const addConnection = () => {
    history.push(`/connections/add/${activeProvider.id}`)
  }

  const configureConnection = (connection) => {
    const { id, endpoint } = connection
    history.push(`/connections/configure/${activeProvider.id}/${id}`)
    console.log('>> editing/modifying connection: ', id, endpoint)
  }

  const testConnection = (connection) => {
    const { id, endpoint } = connection
    console.log('>> testing connection: ', id, endpoint)
  }

  const runCollection = (connection) => {
    const { id, endpoint } = connection
    console.log('>> running connection: ', id, endpoint)
  }

  const refreshConnections = () => {
    fetchConnections()
  }

  useEffect(() => {
    // Selected Provider
    // console.log(activeProvider)
    fetchConnections()
  }, [activeProvider])

  useEffect(() => {
    fetch(`${SERVER_HOST}/api/getenv`)
      .then(response => response.json())
      .then(env => {
        setDbUrl(env.DB_URL)
        setPort(env.PORT)
        setMode(env.MODE)
      })
    console.log('>> ACTIVE PROVIDER = ', providerId)
    console.log(dbUrl, port, mode)
    setActiveProvider(integrations.find(p => p.id === providerId))

    // Fetch Connections for Active Provider
    const providerConnections = []

    // Process Raw Connectin Data => Connection Objects
    providerConnections.map((conn, idx) => {
      return {
        id: idx,
        name: conn.name,
        endpoint: conn.endpoint,
        status: conn.status, // 0=Offline, 1=Online, 2=Error, 3=Collecting
        errors: []
      }
    })

    // Set Live Connection Objects List
    // setConnections(providerConnections)
  }, [])

  return (
    <>
      <div className='container'>
        <Nav />
        <Sidebar />
        <Content>
          <main className='main'>
            <div className='headlineContainer'>
              <Link style={{ float: 'right', marginLeft: '10px', color: '#777777' }} to='/integrations'>
                <Icon icon='fast-backward' size={16} /> &nbsp; Go Back
              </Link>
              <div style={{ display: 'flex' }}>
                <div>
                  <span style={{ marginRight: '10px' }}>{activeProvider.icon}</span> 
                </div>
                <div>
                  <h1 style={{ margin: 0 }}>
                    {activeProvider.name} Integration
                  </h1>
                  <p className='description'>Manage integration and connections.</p>
                </div>
              </div>
            </div>
            <div className='manageProvider'>
              {errors && errors.length > 0 && (
                <Card interactive={false} elevation={Elevation.TWO} style={{ width: '100%', marginBottom: '20px' }}>
                  <div style={{}}>
                    <h4 className='bp3-heading'>
                      <Icon icon='warning-sign' size={18} color={Colors.RED5} style={{ marginRight: '10px' }} />
                      Warning &mdash; This integration has issues
                    </h4>
                    <p className='bp3-ui-text bp3-text-large' style={{ margin: 0 }}>
                      Please see below for all messages that will need to be resolved.
                    </p>
                  </div>
                </Card>
              )}
              {isLoading && (
                <Card interactive={false} elevation={Elevation.TWO} style={{ width: '100%', marginBottom: '20px' }}>
                  <div style={{}}>
                    <div style={{ display: 'flex' }}>
                      <Spinner intent='primary' size={24} />
                      <h4 className='bp3-heading' style={{ marginLeft: '10px' }}>
                        Loading Connections ...
                      </h4>
                    </div>

                    <p className='bp3-ui-text bp3-text-large' style={{ margin: 0 }}>
                      Please wait while the connections are loaded.
                    </p>

                  </div>
                </Card>
              )}
              {!isLoading && connections && connections.length === 0 && (
                <Card interactive={false} elevation={Elevation.TWO} style={{ width: '100%', marginBottom: '20px' }}>
                  <div style={{}}>
                    <h4 className='bp3-heading'>
                      <Icon icon='offline' size={18} color={Colors.GRAY3} style={{ marginRight: '10px' }} /> No Connection Entries
                    </h4>
                    <p className='bp3-ui-text bp3-text-large' style={{ margin: 0 }}>
                      Please check your connection settings and try again.
                      Also verify the authentication token (if applicable) for accuracy.
                    </p>
                    <p className='bp3-monospace-text' style={{ margin: '0 0 20px 0', fontSize: '10px', color: Colors.GRAY4 }}>
                      If the problem persists, please contact our team on <strong>GitHub</strong>
                    </p>
                    <p>
                      <Button
                        onClick={addConnection}
                        rightIcon='add'
                        intent='primary'
                        text='Add Connection'
                        style={{ marginRight: '10px' }}
                      />
                      <Button rightIcon='refresh' text='Refresh Connections' minimal onClick={refreshConnections} />
                    </p>

                  </div>
                </Card>
              )}
              {!isLoading && connections && connections.length > 0 && (
                <>
                  <p>
                    <Button
                      onClick={addConnection}
                      rightIcon='add'
                      intent='primary'
                      text='Add Connection'
                      style={{ marginRight: '10px' }}
                    />
                    <Button rightIcon='refresh' text='Refresh Connections' minimal onClick={refreshConnections} />
                  </p>
                  <Card interactive={false} elevation={Elevation.TWO} style={{ width: '100%', padding: '2px' }}>
                    <table className='bp3-html-table bp3-html-table-bordered connections-table' style={{ width: '100%' }}>
                      <thead>
                        <tr>
                          <th>Connection Name</th>
                          <th>Endpoint</th>
                          <th>Status</th>
                          <th />
                        </tr>
                      </thead>
                      <tbody>
                        {connections.map((connection, idx) => (
                          <tr key={`connection-row-${idx}`} className={connection.status === 0 ? 'connection-offline' : ''}>
                            <td>
                              <strong>{connection.name}</strong>
                            </td>
                            <td><a href='#' target='_blank' rel='noreferrer'>{connection.endpoint}</a></td>
                            <td>
                              {connection.status === 0 && (
                                <strong style={{ color: Colors.GRAY4 }}>Offline</strong>
                              )}
                              {connection.status === 1 && (
                                <strong style={{ color: Colors.GREEN3 }}>Online</strong>
                              )}
                              {connection.status === 2 && (
                                <strong style={{ color: Colors.RED3 }}>Error</strong>
                              )}
                              {connection.status === 3 && (
                                <strong style={{ color: Colors.BLUE3 }}>
                                  <Icon icon='array' size={14} color={Colors.GRAY2} /> Collecting...
                                </strong>
                              )}
                            </td>
                            <td>
                              <a
                                href='#'
                                data-provider={connection.id}
                                className='table-action-link'
                                onClick={() => configureConnection(connection)}
                              >
                                <Icon icon='settings' size={12} />
                                Settings
                              </a>
                              <a
                                href='#'
                                data-provider={connection.id}
                                className='table-action-link'
                                onClick={() => runCollection(connection)}
                              >
                                <Icon icon='refresh' size={12} />
                                Collect
                              </a>
                              <a
                                href='#'
                                data-provider={connection.id}
                                className='table-action-link'
                                onClick={() => testConnection(connection)}
                              >
                                <Icon icon='data-connection' size={12} />
                                Test
                              </a>
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </Card>
                </>
              )}
            </div>
          </main>
        </Content>
      </div>
    </>
  )
}
