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
import Content from '@/components/Content'
import { ToastNotification } from '@/components/Toast'
import ConnectionForm from '@/pages/configure/connections/ConnectionForm'
import { integrationsData } from '@/pages/configure/mock-data/integrations'
import { connectionsData } from '@/pages/configure/mock-data/connections'
import { SERVER_HOST, DEVLAKE_ENDPOINT } from '@/utils/config'

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

  const [isSaving, setIsSaving] = useState(false)
  const [isTesting, setIsTesting] = useState(false)
  const [errors, setErrors] = useState([])
  const [showError, setShowError] = useState(false)

  const [testStatus, setTestStatus] = useState(0) //  0=Pending, 1=Success, 2=Failed

  const [integrations, setIntegrations] = useState(integrationsData)

  const [activeProvider, setActiveProvider] = useState(integrations[0])

  const testConnection = () => {
    setIsTesting(true)
    // Get Testing Endpoint from BE
    // Issue GET and/or POST to test endpoint
    // POST payload to BE
    const connectionTestPayload = {
      name,
      endpointUrl,
      token,
      username,
      password
    }
    // const testResponse = await axios.post(`${DEVLAKE_ENDPOINT}/connection/test`, connectionTestPayload)
    const testResponse = {
      success: true,
      connection: {
        ...connectionTestPayload
      },
      errors: []
    }
    console.log(testResponse)
    setTimeout(() => {
      if (testResponse.success) {
        setIsTesting(false)
        setTestStatus(1)
        ToastNotification.show({ message: 'Connection test OK.', intent: 'success', icon: 'small-tick' })
      } else {
        setIsTesting(false)
        setTestStatus(2)
        ToastNotification.show({ message: 'Connection test FAILED.', intent: 'danger', icon: 'error' })
      }
    }, 2000)
  }

  const saveConnection = async () => {
    setIsSaving(true)
    // POST payload to BE
    const connectionPayload = {
      name,
      endpointUrl,
      token,
      username,
      password
    }
    // const saveResponse = await axios.post(`${DEVLAKE_ENDPOINT}/connection`, connectionPayload)
    // console.log(saveResponse)
    const saveResponse = {
      success: true,
      connection: {
        ...connectionPayload
      },
      errors: []
    }
    setErrors(saveResponse.errors)
    setTimeout(() => {
      if (saveResponse.success && errors.length === 0) {
        ToastNotification.show({ message: 'Connection added successfully.', intent: 'success', icon: 'small-tick' })
        setShowError(false)
        setIsSaving(false)
        resetForm()
        // REDIRECT back to Active Provider Settings
        history.push(`/integrations/${activeProvider.id}`)
      } else {
        ToastNotification.show({ message: 'Connection failed to add.', intent: 'danger', icon: 'error' })
        setShowError(true)
        setIsSaving(false)
      }
    }, 2000)
  }

  const cancel = () => {
    history.push(`/integrations/${activeProvider.id}`)
  }

  const resetForm = () => {
    setName(null)
    setEndpointUrl(null)
    setToken(null)
    setUsername(null)
    setPassword(null)
  }

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
                    {activeProvider.name} Add Connection
                  </h1>
                  <p className='description'>Create a new connection for this provider.</p>
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
