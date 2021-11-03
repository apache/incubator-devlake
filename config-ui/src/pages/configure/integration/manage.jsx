import React, { useEffect, useState } from 'react'
import {
  useParams,
  Link,
  useHistory
} from 'react-router-dom'
import {
  Button, Card, Elevation, Colors,
  Spinner,
  Icon,
} from '@blueprintjs/core'
import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import AppCrumbs from '@/components/Breadcrumbs'
import Content from '@/components/Content'
import { ToastNotification } from '@/components/Toast'
import request from '@/utils/request'

import useConnectionManager from '@/hooks/useConnectionManager'

import { SERVER_HOST, DEVLAKE_ENDPOINT } from '@/utils/config'

import { integrationsData } from '@/pages/configure/mock-data/integrations'

import '@/styles/integration.scss'
import '@blueprintjs/popover2/lib/css/blueprint-popover2.css'

export default function ManageIntegration () {
  const history = useHistory()

  const [dbUrl, setDbUrl] = useState()
  const [port, setPort] = useState()
  const [mode, setMode] = useState()

  const { providerId } = useParams()

  const [isLoading, setIsLoading] = useState(true)
  const [errors, setErrors] = useState([])
  const [integrations, setIntegrations] = useState(integrationsData)
  const [connections, setConnections] = useState([])
  const [activeProvider, setActiveProvider] = useState(integrations[0])

  const {
    sourceLimits,
  } = useConnectionManager({
    activeProvider
  })

  const fetchConnections = async () => {
    setIsLoading(true)
    try {
      const connectionsResponse = await request.get(`${DEVLAKE_ENDPOINT}/plugins/${activeProvider.id}/sources`)
      let providerConnections = connectionsResponse.data || []
      providerConnections = providerConnections.map((conn, idx) => {
        return {
          ...conn,
          status: connectionsResponse.status === 200 || connectionsResponse.status === 201 ? 1 : 0, // conn.status
          id: conn.ID,
          name: conn.name,
          endpoint: conn.endpoint,
          errors: []
        }
      })
      setConnections(providerConnections)
      console.log('>> CONNECTIONS FETCHED', connectionsResponse)
    } catch (e) {
      console.log('>> FAILED TO FETCH CONNECTIONS', e)
      setErrors([e])
      console.log(e)
    }
    setTimeout(() => {
      setIsLoading(false)
      ToastNotification.clear()
      ToastNotification.show({ message: 'Loaded all connections.', intent: 'success', icon: 'small-tick' })
    }, 1000)
  }

  const addConnection = () => {
    history.push(`/connections/add/${activeProvider.id}`)
  }

  const editConnection = (connection, e) => {
    console.log(e.target.classList)
    if (e.target && (!e.target.classList.contains('cell-actions') || !e.target.classList.contains('actions-link'))) {
      history.push(`/connections/edit/${activeProvider.id}/${connection.id}`)
    }
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
    ToastNotification.clear()
    ToastNotification.show({ message: `Triggered Collection Process on ${connection.name}`, icon: 'info-sign' })
    console.log('>> running connection: ', id, endpoint)
  }

  const refreshConnections = () => {
    fetchConnections()
  }

  const maxConnectionsExceeded = (limit, totalConnections) => {
    return totalConnections > 0 && totalConnections >= limit
  }

  useEffect(() => {
    // Selected Provider
    // console.log(activeProvider)
    // !WARNING! DO NOT ADD fetchConnections TO DEPENDENCIES ARRAY!
    // @todo FIXME: Circular-loop issue & Migrate fetching all connections to Connection Manager Hook
    fetchConnections()
  }, [activeProvider])

  useEffect(() => {
    console.log('>> ACTIVE PROVIDER = ', providerId)
    console.log(dbUrl, port, mode)
    setIntegrations(integrations)
    setActiveProvider(integrations.find(p => p.id === providerId))
  }, [])

  useEffect(() => {
    console.log('>> CONNECTION SOURCE LIMITS', sourceLimits)
  }, [connections, sourceLimits])

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
                { href: `/integrations/${activeProvider.id}`, icon: false, text: `${activeProvider.name}`, current: true },
              ]}
            />
            <div className='headlineContainer'>
              <Link style={{ float: 'right', marginLeft: '10px', color: '#777777' }} to='/integrations'>
                <Icon icon='fast-backward' size={16} /> Go Back
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
                      disabled={maxConnectionsExceeded(sourceLimits[activeProvider.id], connections.length)}
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
                          <tr
                            key={`connection-row-${idx}`}
                            className={connection.status === 0 ? 'connection-offline' : ''}
                          >
                            <td
                              onClick={(e) => editConnection(connection, e)}
                              style={{ cursor: 'pointer' }}
                              className='cell-name'
                            >
                              {/* <Icon icon='power' color={Colors.GRAY4} size={10} style={{ float: 'right', marginLeft: '10px' }} /> */}
                              <strong>{connection.name || connection.Name}</strong>
                              <a
                                href='#'
                                data-provider={connection.id}
                                className='table-action-link'
                                onClick={(e) => editConnection(connection, e)}
                              />
                            </td>
                            <td className='cell-endpoint'><a href='#' target='_blank' rel='noreferrer'>{connection.endpoint || connection.Endpoint}</a></td>
                            <td className='cell-status'>
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
                            <td className='cell-actions'>
                              <a
                                href='#'
                                data-provider={connection.id}
                                className='table-action-link actions-link'
                                // onClick={() => configureConnection(connection)}
                                onClick={(e) => editConnection(connection, e)}
                              >
                                <Icon icon='settings' size={12} />
                                Settings
                              </a>
                              {/* <a
                                href='#'
                                data-provider={connection.id}
                                className='table-action-link actions-link'
                                onClick={() => runCollection(connection)}
                              >
                                <Icon icon='refresh' size={12} />
                                Collect
                              </a>
                              <a
                                href='#'
                                data-provider={connection.id}
                                className='table-action-link actions-link'
                                onClick={() => testConnection(connection)}
                              >
                                <Icon icon='data-connection' size={12} />
                                Test
                              </a> */}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                    {maxConnectionsExceeded(sourceLimits[activeProvider.id], connections.length) && (
                      <p style={{ margin: 0, padding: '10px', backgroundColor: '#f0f0f0', borderTop: '1px solid #cccccc' }}>
                        <Icon icon='warning-sign' size='16' color={Colors.GRAY1} style={{ marginRight: '5px' }} />
                        You have reached the maximum number of allowed connections for this provider.
                      </p>
                    )}
                  </Card>
                  <p style={{
                    textAlign: 'right',
                    margin: '12px 6px',
                    fontSize: '10px',
                    color: '#aaaaaa'
                  }}
                  >Fetched <strong>{connections.length}</strong> connections from Lake API for <strong>{activeProvider.name}</strong>
                  </p>
                </>
              )}
            </div>
          </main>
        </Content>
      </div>
    </>
  )
}
