import React, { useEffect, useState } from 'react'
import {
  Button, Colors,
  FormGroup, InputGroup, Label,
  Icon,
} from '@blueprintjs/core'

import '@/styles/integration.scss'
import '@/styles/connections.scss'
import '@blueprintjs/popover2/lib/css/blueprint-popover2.css'

export default function ConnectionForm (props) {
  const {
    name,
    endpointUrl,
    token,
    username,
    password,
    isSaving,
    isTesting,
    showError,
    errors,
    testStatus,
    onSave = () => {},
    onCancel = () => {},
    onTest = () => {},
    onNameChange = () => {},
    onEndpointChange = () => {},
    onTokenChange = () => {},
    onUsernameChange = () => {},
    onPasswordChange = () => {},
    authType = 'token'
  } = props

  const [allowedAuthTypes, setAllowedAuthTypes] = useState(['token', 'plain'])

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

  useEffect(() => {
    if (!allowedAuthTypes.includes(authType)) {
      console.log('INVALID AUTH TYPE!')
    }
  }, [authType, allowedAuthTypes])

  useEffect(() => {
    setAllowedAuthTypes(['token', 'plain'])
  }, [])

  return (
    <>
      <form className='form form-add-connection'>
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
            helperText='NAME'
            className='formGroup'
            contentClassName='formGroupContent'
          >
            <Label style={{ display: 'inline' }}>
              Connection&nbsp;Name <span className='requiredStar'>*</span>
            </Label>
            <InputGroup
              id='connection-name'
              disabled={isTesting || isSaving}
              placeholder='Enter Instance Name eg. ISSUES-AWS-US-EAST'
              defaultValue={name}
              onChange={(e) => onNameChange(e.target.value)}
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
            helperText='ENDPOINT_URL'
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
              onChange={(e) => onEndpointChange(e.target.value)}
              className='input'
              fill
            />
            {/* <a href='#' style={{ margin: '5px 0 5px 5px' }}><Icon icon='info-sign' size='16' /></a> */}
          </FormGroup>
        </div>

        {authType === 'token' && (
          <div className='formContainer'>
            <FormGroup
              disabled={isTesting || isSaving}
              label=''
              inline={true}
              labelFor='connection-token'
              helperText='TOKEN'
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
                onChange={(e) => onTokenChange(e.target.value)}
                className='input'
                fill
                required
              />
              {/* <a href='#' style={{ margin: '5px 0 5px 5px' }}><Icon icon='info-sign' size='16' /></a> */}
            </FormGroup>
          </div>
        )}
        {authType === 'plain' && (
          <>
            <div style={{ marginTop: '20px', marginBottom: '20px' }}>
              <h3 style={{ margin: 0 }}>Username & Password</h3>
              <span className='description' style={{ margin: 0, color: Colors.GRAY2 }}>
                If this connection uses login credentials to generate a token or uses PLAIN Auth, specify it here.
              </span>
            </div>
            <div className='formContainer'>
              <FormGroup
                label=''
                disabled={isTesting || isSaving}
                inline={true}
                labelFor='connection-username'
                helperText='USERNAME'
                className='formGroup'
                contentClassName='formGroupContent'
              >
                <Label style={{ display: 'inline' }}>
                  Username <span className='requiredStar'>*</span>
                </Label>
                <InputGroup
                  id='connection-username'
                  disabled={isTesting || isSaving}
                  placeholder='Enter Username'
                  defaultValue={username}
                  onChange={(e) => onUsernameChange(e.target.value)}
                  className='input'
                  style={{ maxWidth: '300px' }}
                />
              </FormGroup>
            </div>
            <div className='formContainer'>
              <FormGroup
                disabled={isTesting || isSaving}
                label=''
                inline={true}
                labelFor='connection-password'
                helperText='PASSWORD'
                className='formGroup'
                contentClassName='formGroupContent'
              >
                <Label style={{ display: 'inline' }}>
                  Password <span className='requiredStar'>*</span>
                </Label>
                <InputGroup
                  id='connection-password'
                  type='password'
                  disabled={isTesting || isSaving}
                  placeholder='Enter Password'
                  defaultValue={password}
                  onChange={(e) => onPasswordChange(e.target.value)}
                  className='input'
                  style={{ maxWidth: '300px' }}
                />
              </FormGroup>
            </div>
          </>
        )}
        <div style={{ display: 'flex', marginTop: '30px', justifyContent: 'space-between', maxWidth: '50%' }}>
          <div>
            <Button
              icon={getConnectionStatusIcon()}
              text='Test Connection'
              onClick={onTest}
              loading={isTesting}
              disabled={isTesting || isSaving}
            />
          </div>
          <div>
            <Button icon='remove' text='Cancel' onClick={onCancel} disabled={isSaving || isTesting} />
            <Button
              icon='cloud-upload' intent='primary' text='Save Connection'
              loading={isSaving}
              disabled={isSaving || isTesting}
              onClick={onSave}
              style={{ marginLeft: '10px' }}
            />
          </div>
        </div>
      </form>
    </>
  )
}
