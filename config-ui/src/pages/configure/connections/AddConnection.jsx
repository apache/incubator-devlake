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
import { ToastNotification } from '@/components/Toast'

import { SERVER_HOST, DEVLAKE_ENDPOINT } from '@/utils/config'

import { ReactComponent as GitlabProvider } from '@/images/integrations/gitlab.svg'
import { ReactComponent as JenkinsProvider } from '@/images/integrations/jenkins.svg'
import { ReactComponent as JiraProvider } from '@/images/integrations/jira.svg'

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
      success: false,
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

  const getConnectionStatusIcon = () => {
    let statusIcon = <Icon icon='full-circle' size='10' color={Colors.RED3} />
    switch (testStatus) {
      case 1:
        statusIcon = <Icon icon='full-circle' size='10' color={Colors.GREEN3} />
        break
      case 2:
        statusIcon = <Icon icon='full-circle' size='10' color={Colors.RED3} />
        break
      case 0:
      default:
        statusIcon = <Icon icon='full-circle' size='10' color={Colors.GRAY3} />
        break
    }
    return statusIcon
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
              <h1><span style={{ marginRight: '10px' }}>{activeProvider.icon}</span> {activeProvider.name} Add Connection </h1>
              <p className='description'>Create a new connection for this provider.</p>
              <div className='addConnection' style={{ display: 'flex' }}>
                <form className='form'>

                  <div className='headlineContainer'>
                    <h2 className='headline'>Configure Instance</h2>
                    <p className='description'>Account & Authentication settings</p>
                  </div>

                  {showError && (
                    <div className='bp3-callout bp3-intent-danger' style={{ margin: '20px 0', maxWidth: '50%' }}>
                      <h4 className='bp3-heading'>Operation Failed</h4>
                      Your connection could not be saved.
                      {errors.length > 0 && (
                        <ul>
                          {errors.map((errorMessage, idx) => (
                            <li key={`save-error-message-${idx}`}>{errorMessage}</li>
                          ))}
                        </ul>)}
                    </div>)}

                  <div className='formContainer'>
                    <FormGroup
                      disabled={isTesting || isSaving}
                      label=''
                      inline={true}
                      labelFor='connection-name'
                      helperText='JIRA_INSTANCE_NAME'
                      className='formGroup'
                      contentClassName='formGroupContent'
                    >
                      <Label style={{ display: 'inline' }}>
                        Connection&nbsp;Name <span className='requiredStar'>*</span>
                      </Label>
                      <InputGroup
                        id='connection-name'
                        disabled={isTesting || isSaving}
                        placeholder='Enter Instance Name eg. JIRA-AWS-US-EAST'
                        defaultValue={name}
                        onChange={(e) => setName(e.target.value)}
                        className='input'
                        fill
                      />
                    </FormGroup>
                  </div>

                  <div className='formContainer'>
                    <FormGroup
                      disabled={isTesting || isSaving}
                      label=''
                      inline={true}
                      labelFor='connection-endpoint'
                      helperText='JIRA_ENDPOINT'
                      className='formGroup'
                      contentClassName='formGroupContent'
                    >
                      <Label>
                        Endpoint&nbsp;URL <span className='requiredStar'>*</span>
                      </Label>
                      <InputGroup
                        id='connection-endpoint'
                        disabled={isTesting || isSaving}
                        placeholder='Enter Endpoint URL eg. https://merico.atlassian.net/rest'
                        defaultValue={endpointUrl}
                        onChange={(e) => setEndpointUrl(e.target.value)}
                        className='input'
                        fill
                      />
                      <a href='#' style={{ margin: '5px 0 5px 5px' }}><Icon icon='info-sign' size='16' /></a>
                    </FormGroup>
                  </div>

                  <div className='formContainer'>
                    <FormGroup
                      disabled={isTesting || isSaving}
                      label=''
                      inline={true}
                      labelFor='connection-token'
                      helperText='JIRA_BASIC_AUTH_ENCODED'
                      className='formGroup'
                      contentClassName='formGroupContent'
                    >
                      <Label>
                        Basic&nbsp;Auth&nbsp;Token <span className='requiredStar'>*</span>
                      </Label>
                      <InputGroup
                        id='connection-token'
                        disabled={isTesting || isSaving}
                        placeholder='Enter Auth Token eg. EJrLG8DNeXADQcGOaaaX4B47'
                        defaultValue={token}
                        onChange={(e) => setToken(e.target.value)}
                        className='input'
                        fill
                      />
                      <a href='#' style={{ margin: '5px 0 5px 5px' }}><Icon icon='info-sign' size='16' /></a>
                    </FormGroup>
                  </div>

                  <div style={{ marginTop: '20px', marginBottom: '20px' }}>
                    <h3 style={{ margin: 0 }}>Username & Password</h3>
                    <span className='description' style={{ margin: 0, color: Colors.GRAY2 }}>
                      If this connection uses login credentials to generate a token or uses PLAN Auth, specify it here.
                    </span>
                  </div>
                  <div className='formContainer'>
                    <FormGroup
                      label=''
                      disabled={isTesting || isSaving}
                      inline={true}
                      labelFor='connection-username'
                      helperText='ACCOUNT_NAME'
                      className='formGroup'
                      contentClassName='formGroupContent'
                    >
                      <Label style={{ display: 'inline' }}>
                        Username
                      </Label>
                      <InputGroup
                        id='connection-username'
                        disabled={isTesting || isSaving}
                        placeholder='Enter Username'
                        defaultValue={username}
                        onChange={(e) => setUsername(e.target.value)}
                        className='input'
                      />
                    </FormGroup>
                  </div>
                  <div className='formContainer'>
                    <FormGroup
                      disabled={isTesting || isSaving}
                      label=''
                      inline={true}
                      labelFor='connection-password'
                      helperText='ACCOUNT_PASSWORD'
                      className='formGroup'
                      contentClassName='formGroupContent'
                    >
                      <Label style={{ display: 'inline' }}>
                        Password
                      </Label>
                      <InputGroup
                        id='connection-password'
                        type='password'
                        disabled={isTesting || isSaving}
                        placeholder='Enter Password'
                        defaultValue={password}
                        onChange={(e) => setPassword(e.target.value)}
                        className='input'
                      />
                    </FormGroup>
                  </div>
                  <div style={{ display: 'flex', marginTop: '30px', justifyContent: 'space-between', maxWidth: '50%' }}>
                    <div>
                      <Button
                        icon={getConnectionStatusIcon()}
                        text='Test Connection'
                        onClick={testConnection}
                        loading={isTesting}
                        disabled={isTesting || isSaving}
                      />
                    </div>
                    <div>
                      <Button icon='remove' text='Cancel' onClick={cancel} disabled={isSaving || isTesting} />
                      <Button
                        icon='cloud-upload' intent='primary' text='Save'
                        loading={isSaving}
                        disabled={isSaving || isTesting}
                        onClick={saveConnection}
                        style={{ backgroundColor: '#E8471C', marginLeft: '10px' }}
                      />
                    </div>
                  </div>
                </form>
              </div>
            </div>
          </main>
        </Content>
      </div>
    </>
  )
}
