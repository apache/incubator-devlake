/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
import React, { useEffect, useState, useCallback, useRef, useMemo } from 'react'
import {
  Button,
  Colors,
  FormGroup,
  InputGroup,
  Label,
  TextArea,
  Card,
  Icon,
  Tag,
  Elevation,
  Popover,
  // PopoverInteractionKind,
  Intent,
  PopoverInteractionKind,
  NumericInput,
  Switch
} from '@blueprintjs/core'
import { Providers } from '@/data/Providers'
import FormValidationErrors from '@/components/messages/FormValidationErrors'
import InputValidationError from '@/components/validation/InputValidationError'

import '@/styles/integration.scss'
import '@/styles/connections.scss'

export default function ConnectionForm(props) {
  const {
    isLocked = false,
    isValid = true,
    validationErrors = [],
    activeProvider,
    name,
    endpointUrl,
    token,
    initialTokenStore = {
      0: '',
      1: '',
      2: ''
    },
    username,
    password,
    proxy = '',
    rateLimitPerHour = 0,
    isSaving,
    isTesting,
    showError,
    errors,
    testStatus,
    testResponse,
    allTestResponses,
    onSave = () => {},
    onCancel = () => {},
    onTest = () => {},
    onNameChange = () => {},
    onEndpointChange = () => {},
    onTokenChange = () => {},
    onUsernameChange = () => {},
    onPasswordChange = () => {},
    onProxyChange = () => {},
    onRateLimitChange = () => {},
    onValidate = () => {},
    authType = 'token',
    sourceLimits = {},
    showLimitWarning = true,
    labels,
    placeholders,
    enableActions = true,
    formGroupClassName = 'formGroup',
    showHeadline = true
  } = props

  const connectionNameRef = useRef()
  const connectionEndpointRef = useRef()
  const connectionTokenRef = useRef()
  const connectionUsernameRef = useRef()
  const connectionPasswordRef = useRef()
  const connectionProxyRef = useRef()
  const connectionRateLimitRef = useRef()

  // const [isValidForm, setIsValidForm] = useState(true)
  const [allowedAuthTypes, setAllowedAuthTypes] = useState(['token', 'plain'])
  const [stateErrored, setStateErrored] = useState(false)
  const [tokenStore, setTokenStore] = useState(initialTokenStore)
  const [personalAccessTokens, setPersonalAccessTokens] = useState([])
  const [tokenTests, setTokenTests] = useState([])
  const [enableRateLimit, setEnableRateLimit] = useState(rateLimitPerHour > 0)

  const patTestPayload = useMemo(
    () => ({ endpoint: endpointUrl, proxy }),
    [endpointUrl, proxy]
  )

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

  const validate = useCallback(() => {
    onValidate({
      name,
      endpointUrl,
      token,
      username,
      password,
      rateLimitPerHour
    })
  }, [
    name,
    endpointUrl,
    token,
    username,
    password,
    rateLimitPerHour,
    onValidate
  ])
  const fieldHasError = (fieldId) => {
    return validationErrors.some((e) => e.includes(fieldId))
  }

  const getFieldError = (fieldId) => {
    return validationErrors.find((e) => e.includes(fieldId))
  }

  const getValidityStatus = useCallback(
    (token, testResponse, allTestResponses) => {
      let status = 'Invalid'
      // @todo: add duplicate user check
      if (testResponse?.success) {
        status = 'Valid'
      } else {
        status = 'Invalid'
      }
      return status
    },
    []
  )

  const activateErrorStates = (elementId) => {
    setStateErrored(elementId || false)
  }

  const addAnotherAccessToken = () => {
    const emptyToken = ''
    setTokenStore((tokens) => ({
      ...tokens,
      [Object.keys(tokens).length]: emptyToken
    }))
  }

  const setPersonalToken = (id, newToken) => {
    setTokenStore((tokens) => ({ ...tokens, [id]: newToken }))
  }

  const removePersonalToken = (id) => {
    setTokenStore((tokens) =>
      Object.values(tokens)
        .filter((t, tId) => tId !== id)
        .reduce((newStore, cT, tId) => ({ ...newStore, [tId]: cT }), {})
    )
  }

  const syncPersonalAcessTokens = useCallback(() => {
    onTokenChange(personalAccessTokens.join(',').trim())
    // eslint-disable-next-line no-unused-vars
    const tokenTestResponses = personalAccessTokens
      .filter((t) => t !== '')
      .map((t, tIdx) => {
        const withAlert = false
        const testPayload = {
          ...patTestPayload,
          token: t
        }
        if (
          ![
            connectionEndpointRef.current?.id,
            connectionProxyRef.current?.id
          ].includes(document.activeElement?.id)
        ) {
          onTest(withAlert, testPayload)
        }
        return t
      })
  }, [personalAccessTokens, patTestPayload, onTokenChange, onTest])

  useEffect(() => {
    if (!allowedAuthTypes.includes(authType)) {
      console.log('INVALID AUTH TYPE!')
    }
  }, [authType, allowedAuthTypes])

  useEffect(() => {
    setAllowedAuthTypes(['token', 'plain'])
  }, [])

  useEffect(() => {
    validate()
  }, [name, endpointUrl, token, username, password, validate])

  useEffect(() => {
    console.log('>> CONNECTION FORM VALIDATION STATUS CHANGED...', isValid)
  }, [isValid])

  useEffect(() => {
    console.log('>> PERSONAL TOKEN STORE UPDATED...', tokenStore)
    setPersonalAccessTokens(Object.values(tokenStore).filter((t) => t !== ''))
  }, [tokenStore])

  useEffect(() => {
    if (activeProvider?.id === Providers.GITHUB) {
      console.log('>> PERSONAL ACCESS TOKENS ENTERED...', personalAccessTokens)
      syncPersonalAcessTokens()
    }
  }, [activeProvider?.id, personalAccessTokens, syncPersonalAcessTokens])

  useEffect(() => {
    console.log(
      '>> INITIALIZED PERSONAL TOKEN STORE WITH...',
      initialTokenStore
    )
    setTokenStore(initialTokenStore)
  }, [initialTokenStore])

  useEffect(() => {
    console.log('>> PERSONAL ACCESS TOKENS TEST RESULTS...', tokenTests)
  }, [tokenTests])

  useEffect(() => {
    console.log('>> ALL TEST RESPONSES FROM CONN MGR...', allTestResponses)
    setTokenTests(
      personalAccessTokens
        .filter((t) => t !== '')
        .map((t, tIdx) => {
          return {
            id: tIdx,
            token: t,
            response: allTestResponses ? allTestResponses[t] : null,
            success: allTestResponses ? allTestResponses[t]?.success : false,
            message: allTestResponses
              ? allTestResponses[t]?.message
              : 'invalid token',
            username: allTestResponses ? allTestResponses[t]?.login : '',
            status: getValidityStatus(
              t,
              allTestResponses ? allTestResponses[t] : null,
              allTestResponses
            ),
            isDuplicate: !!(allTestResponses && allTestResponses[t])
          }
        })
    )
  }, [allTestResponses, personalAccessTokens, getValidityStatus])

  return (
    <>
      <form className='form form-add-connection'>
        {showHeadline && (
          <div className='headlineContainer'>
            <h2
              className='headline'
              style={{
                marginTop: 0,
                textDecoration: isLocked ? 'line-through' : 'none'
              }}
            >
              Configure Connection
            </h2>
            <p className='description'>
              Instance Account & Authentication settings
            </p>
          </div>
        )}

        {showError && (
          <Card
            className='app-error-card'
            interactive={false}
            elevation={showLimitWarning ? Elevation.TWO : Elevation.ZERO}
            style={{
              maxWidth: '480px',
              marginBottom: '20px',
              backgroundColor: showLimitWarning ? '#f0f0f0' : 'transparent',
              border: showLimitWarning ? 'inherit' : 0
            }}
          >
            <p className='warning-message' intent={Intent.WARNING}>
              <Icon
                icon='error'
                size='16'
                color={Colors.RED4}
                style={{ marginRight: '5px' }}
              />
              <strong>
                UNABLE TO SAVE CONNECTION ({name !== '' ? name : 'BLANK'})
              </strong>
              <br />
            </p>
            {errors.length > 0 && (
              <ul>
                {errors.map((errorMessage, idx) => (
                  <li key={`save-error-message-${idx}`}>{errorMessage}</li>
                ))}
              </ul>
            )}
          </Card>
        )}

        <div className='formContainer'>
          <FormGroup
            disabled={isTesting || isSaving || isLocked}
            label=''
            inline={true}
            labelFor='connection-name'
            className={formGroupClassName}
            contentClassName='formGroupContent'
          >
            <Label>
              {labels ? labels.name : <>Connection&nbsp;Name</>}
              <span className='requiredStar'>*</span>
            </Label>
            <InputGroup
              id='connection-name'
              autoComplete='false'
              autoFocus={true}
              inputRef={connectionNameRef}
              disabled={isTesting || isSaving || isLocked}
              // readOnly={[Providers.JENKINS].includes(activeProvider.id)}
              placeholder={
                placeholders ? placeholders.name : 'Enter Instance Name'
              }
              value={name}
              onChange={(e) => onNameChange(e.target.value)}
              className={`input connection-name-input ${
                stateErrored === 'connection-name' ? 'invalid-field' : ''
              }`}
              // leftIcon={[Providers.GITHUB, Providers.GITLAB, Providers.JENKINS].includes(activeProvider.id) ? 'lock' : null}
              inline={true}
              rightElement={
                <InputValidationError
                  error={getFieldError('Connection')}
                  elementRef={connectionNameRef}
                  onError={activateErrorStates}
                  onSuccess={() => setStateErrored(null)}
                  validateOnFocus
                />
              }
              // fill
            />
          </FormGroup>
        </div>

        <div className='formContainer'>
          <FormGroup
            disabled={isTesting || isSaving || isLocked}
            label=''
            inline={true}
            labelFor='connection-endpoint'
            className={formGroupClassName}
            contentClassName='formGroupContent'
          >
            <Label>
              {labels ? labels.endpoint : <>Endpoint&nbsp;URL</>}
              <span className='requiredStar'>*</span>
            </Label>
            <InputGroup
              id='connection-endpoint'
              autoComplete='false'
              inputRef={connectionEndpointRef}
              disabled={isTesting || isSaving || isLocked}
              placeholder={
                placeholders ? placeholders.endpoint : 'Enter Endpoint URL'
              }
              value={endpointUrl}
              onChange={(e) => onEndpointChange(e.target.value)}
              className={`input endpoint-url-input ${
                stateErrored === 'connection-endpoint' ? 'invalid-field' : ''
              }`}
              fill
              rightElement={
                <InputValidationError
                  error={getFieldError('Endpoint')}
                  elementRef={connectionEndpointRef}
                  onError={activateErrorStates}
                  onSuccess={() => setStateErrored(null)}
                  validateOnFocus
                />
              }
            />
            {/* <a href='#' style={{ margin: '5px 0 5px 5px' }}><Icon icon='info-sign' size='16' /></a> */}
          </FormGroup>
        </div>

        {authType === 'token' && (
          <div className='formContainer'>
            <FormGroup
              disabled={isTesting || isSaving || isLocked}
              label=''
              inline={true}
              labelFor='connection-token'
              className={formGroupClassName}
              contentClassName='formGroupContent'
            >
              <Label>
                {labels ? labels.token : <>Basic&nbsp;Auth&nbsp;Token</>}
                <span className='requiredStar'>*</span>
              </Label>
              {[Providers.GITHUB].includes(activeProvider?.id) ? (
                <>
                  {/* TEXTAREA Multi-line Token Input (Disabled) */}
                  {/* <div
                    className='bp3-input-group connection-token-group' style={{
                      boxSizing: 'border-box',
                      width: '99%',
                      position: 'relative',
                      display: 'flex'
                    }}
                  >
                    <TextArea
                      id='connection-token'
                      className={`input auth-input ${stateErrored === 'connection-token' ? 'invalid-field' : ''}`}
                      inputRef={connectionTokenRef}
                      disabled={isTesting || isSaving || isLocked}
                      placeholder={placeholders ? placeholders.token : 'Enter Auth Token eg. EJrLG8DNeXADQcGOaaaX4B47'}
                      growVertically={true}
                      large={true}
                      // intent={Intent.PRIMARY}
                      onChange={(e) => onTokenChange(e.target.value)}
                      value={token}
                      required
                      fill
                      style={{ maxWidth: '99%' }}
                    />
                    <span style={{ marginLeft: '-23px', zIndex: 1 }}>
                      <InputValidationError
                        error={getFieldError('Auth')}
                        elementRef={connectionTokenRef}
                        onError={activateErrorStates}
                        onSuccess={() => setStateErrored(null)}
                        validateOnFocus
                      />
                    </span>
                  </div> */}
                  <div className='connection-tokens-personal-group'>
                    <p>
                      Add one or more personal token(s) for authentication from
                      you and your organization members. Multiple tokens can
                      help speed up the data collection process.{' '}
                    </p>
                    <p>
                      <a
                        // eslint-disable-next-line max-len
                        href='https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token'
                        target='_blank'
                        rel='noreferrer'
                      >
                        Learn about how to create a personal access token
                      </a>
                    </p>
                    <label className='normal'>Personal Access Token(s)</label>
                    <div
                      className='personal-access-tokens'
                      style={{ margin: '5px 0' }}
                    >
                      <div
                        className='pats-inputgroup'
                        style={{ display: 'flex', flexDirection: 'column' }}
                      >
                        {Object.values(tokenStore).map((pat, patIdx) => (
                          <div
                            className='pat-input'
                            key={`pat-input-key-${patIdx}`}
                            style={{
                              display: 'flex',
                              flex: 1,
                              marginBottom: '8px'
                            }}
                          >
                            <div
                              className='token-input'
                              style={{ flex: 1, maxWidth: '55%' }}
                            >
                              <InputGroup
                                id={`pat-id-${patIdx}`}
                                type='password'
                                placeholder='Token'
                                value={pat}
                                onChange={(e) =>
                                  setPersonalToken(patIdx, e.target.value)
                                }
                                className='input personal-token-input'
                                fill
                                autoComplete='false'
                              />
                            </div>
                            {tokenTests[patIdx] && (
                              <div
                                className='token-info-status'
                                style={{ display: 'flex', padding: '0 10px' }}
                              >
                                {tokenTests[patIdx]?.success && pat !== '' ? (
                                  <>
                                    <span
                                      style={{
                                        color: Colors.GRAY4,
                                        marginRight: '10px'
                                      }}
                                    >
                                      From: {tokenTests[patIdx]?.username}
                                    </span>
                                    <span
                                      className='token-validation-status'
                                      style={{ color: Colors.GREEN4 }}
                                    >
                                      {tokenTests[patIdx].status}
                                    </span>
                                  </>
                                ) : (
                                  <>
                                    <span
                                      className='token-validation-status'
                                      style={{ color: Colors.RED4 }}
                                    >
                                      {pat === '' ? '' : 'Invalid'}
                                    </span>
                                  </>
                                )}
                              </div>
                            )}
                            <div
                              className='token-removal'
                              style={{
                                marginLeft: 'auto',
                                justifyContent: 'flex-end'
                              }}
                            >
                              <Button
                                icon='small-cross'
                                intent={Intent.PRIMARY}
                                minimal
                                small
                                onClick={() => removePersonalToken(patIdx)}
                              />
                            </div>
                          </div>
                        ))}
                      </div>
                      <div
                        className='pats-actions'
                        style={{ marginTop: '5px' }}
                      >
                        <Button
                          disabled={isSaving || isTesting}
                          text='Another Token'
                          icon='plus'
                          intent={Intent.PRIMARY}
                          small
                          outlined
                          onClick={() =>
                            addAnotherAccessToken(personalAccessTokens.length)
                          }
                        />
                      </div>
                    </div>
                  </div>
                </>
              ) : (
                <>
                  <InputGroup
                    id='connection-token'
                    type='password'
                    autoComplete='false'
                    inputRef={connectionTokenRef}
                    disabled={isTesting || isSaving || isLocked}
                    placeholder={
                      placeholders
                        ? placeholders.token
                        : 'Enter Auth Token eg. EJrLG8DNeXADQcGOaaaX4B47'
                    }
                    value={token}
                    onChange={(e) => onTokenChange(e.target.value)}
                    className={`input auth-input ${
                      stateErrored === 'connection-token' ? 'invalid-field' : ''
                    }`}
                    fill
                    required
                    rightElement={
                      <InputValidationError
                        error={getFieldError('Auth')}
                        elementRef={connectionTokenRef}
                        onError={activateErrorStates}
                        onSuccess={() => setStateErrored(null)}
                        validateOnFocus
                      />
                    }
                  />
                </>
              )}
              {/* activeProvider.id === Providers.JIRA &&
                  <Popover
                    className='popover-generate-token'
                    position={Position.RIGHT}
                    autoFocus={false}
                    enforceFocus={false}
                    isOpen={showTokenCreator}
                    onInteraction={handleTokenInteraction}
                    onClosed={() => setShowTokenCreator(false)}
                    usePortal={false}
                  >
                    <Button
                      disabled={isTesting || isSaving || isLocked}
                      type='button' icon='key' intent={Intent.PRIMARY} style={{ marginLeft: '5px' }}
                    />
                    <>
                      <div style={{ padding: '15px 20px 15px 15px' }}>
                        <GenerateTokenForm
                          isTesting={isTesting}
                          isSaving={isSaving}
                          isLocked={isLocked}
                          onTokenChange={onTokenChange}
                          setShowTokenCreator={setShowTokenCreator}
                        />
                      </div>
                    </>
                  </Popover> */}
              {/* <a href='#' style={{ margin: '5px 0 5px 5px' }}><Icon icon='info-sign' size='16' /></a> */}
            </FormGroup>
          </div>
        )}
        {authType === 'plain' && (
          <>
            {/* <div style={{ marginTop: '20px', marginBottom: '20px' }}>
              <h3 style={{ margin: 0 }}>Username & Password</h3>
              <span className='description' style={{ margin: 0, color: Colors.GRAY2 }}>
                If this connection uses login credentials to generate a token or uses PLAIN Auth, specify it here.
              </span>
            </div> */}
            <div className='formContainer'>
              <FormGroup
                label=''
                disabled={isTesting || isSaving || isLocked}
                inline={true}
                labelFor='connection-username'
                className={formGroupClassName}
                contentClassName='formGroupContent'
              >
                <Label>
                  {labels ? labels.username : <>Username</>}
                  <span className='requiredStar'>*</span>
                </Label>
                <InputGroup
                  id='connection-username'
                  tabIndex={0}
                  autoComplete='false'
                  inputRef={connectionUsernameRef}
                  disabled={isTesting || isSaving || isLocked}
                  placeholder='Enter Username'
                  defaultValue={username}
                  onChange={(e) => onUsernameChange(e.target.value)}
                  className={`input username-input ${
                    stateErrored === 'Username' ? 'invalid-field' : ''
                  }`}
                  // style={{ maxWidth: '300px' }}
                  rightElement={
                    <InputValidationError
                      error={getFieldError('Username')}
                      elementRef={connectionUsernameRef}
                      onError={activateErrorStates}
                      onSuccess={() => setStateErrored(null)}
                      validateOnFocus
                    />
                  }
                />
              </FormGroup>
            </div>
            <div className='formContainer'>
              <FormGroup
                disabled={isTesting || isSaving || isLocked}
                label=''
                inline={true}
                labelFor='connection-password'
                className={formGroupClassName}
                contentClassName='formGroupContent'
              >
                <Label>
                  {labels ? labels.password : <>Password</>}
                  <span className='requiredStar'>*</span>
                </Label>
                <InputGroup
                  id='connection-password'
                  tabIndex={1}
                  type='password'
                  autoComplete='false'
                  inputRef={connectionPasswordRef}
                  disabled={isTesting || isSaving || isLocked}
                  placeholder='Enter Password'
                  defaultValue={password}
                  onChange={(e) => onPasswordChange(e.target.value)}
                  className={`input password-input ${
                    stateErrored === 'Password' ? 'invalid-field' : ''
                  }`}
                  // style={{ maxWidth: '300px' }}
                  rightElement={
                    <InputValidationError
                      error={getFieldError('Password')}
                      elementRef={connectionPasswordRef}
                      onError={activateErrorStates}
                      onSuccess={() => setStateErrored(null)}
                      validateOnFocus
                    />
                  }
                />
              </FormGroup>
            </div>
          </>
        )}
        {[
          Providers.GITHUB,
          Providers.GITLAB,
          Providers.JIRA,
          Providers.JENKINS,
          Providers.TAPD
        ].includes(activeProvider?.id) && (
          <>
            <div className='formContainer'>
              <FormGroup
                disabled={isTesting || isSaving || isLocked}
                inline={true}
                labelFor='connection-proxy'
                className={formGroupClassName}
                contentClassName='formGroupContent'
              >
                <Label>{labels ? labels.proxy : <>Proxy&nbsp;URL</>}</Label>
                <InputGroup
                  id='connection-proxy'
                  inputRef={connectionProxyRef}
                  tabIndex={3}
                  placeholder={
                    placeholders.proxy
                      ? placeholders.proxy
                      : 'http://proxy.localhost:8080'
                  }
                  defaultValue={proxy}
                  onChange={(e) => onProxyChange(e.target.value)}
                  disabled={isTesting || isSaving || isLocked}
                  className={`input input-proxy ${
                    fieldHasError('Proxy') ? 'invalid-field' : ''
                  }`}
                  rightElement={
                    <InputValidationError error={getFieldError('Proxy')} />
                  }
                />
              </FormGroup>
            </div>
            <div className='formContainer'>
              <FormGroup
                disabled={isTesting || isSaving || isLocked}
                inline={true}
                labelFor='connection-ratelimit'
                className={formGroupClassName}
                contentClassName='formGroupContent'
              >
                <Label>
                  {labels ? labels.rateLimitPerHour : <>Rate&nbsp;Limit</>}
                </Label>
                <div
                  className='ratelimit-options'
                  style={{ display: 'flex', float: 'left' }}
                >
                  {enableRateLimit && (
                    <div style={{ marginRight: '10px' }}>
                      <NumericInput
                        id='connection-ratelimit'
                        ref={connectionRateLimitRef}
                        disabled={isTesting || isSaving || isLocked}
                        min={0}
                        max={1000000000}
                        clampValueOnBlur={true}
                        className={`input input-ratelimit ${
                          fieldHasError('RateLimit') ? 'invalid-field' : ''
                        }`}
                        fill={false}
                        placeholder={
                          placeholders.rateLimitPerHour
                            ? placeholders.rateLimitPerHour
                            : '1000'
                        }
                        allowNumericCharactersOnly={true}
                        onValueChange={(rateLimitPerHour) => {
                          onRateLimitChange(rateLimitPerHour)
                        }}
                        value={rateLimitPerHour}
                        rightElement={
                          <InputValidationError
                            error={getFieldError('RateLimit')}
                          />
                        }
                      />
                    </div>
                  )}
                  <Switch
                    checked={enableRateLimit}
                    label={
                      enableRateLimit ? (
                        <>Enabled &mdash; {rateLimitPerHour} Requests/hr</>
                      ) : (
                        'Disabled'
                      )
                    }
                    onChange={() => setEnableRateLimit(!enableRateLimit)}
                    style={{ marginBottom: '0' }}
                  />
                </div>
              </FormGroup>
            </div>
          </>
        )}
        {enableActions && (
          <div
            className='form-actions-block'
            style={{
              display: 'flex',
              marginTop: '30px',
              justifyContent: 'space-between'
            }}
          >
            <div style={{ display: 'flex' }}>
              <Button
                id='btn-test'
                className='btn-test-connection'
                icon={getConnectionStatusIcon()}
                text='Test Connection'
                onClick={onTest}
                loading={isTesting}
                disabled={isTesting || isSaving || isLocked}
              />
            </div>
            <div style={{ display: 'flex' }}>
              <div style={{ justifyContent: 'center', padding: '8px' }}>
                {validationErrors.length > 0 && (
                  <Popover
                    interactionKind={PopoverInteractionKind.HOVER_TARGET_ONLY}
                  >
                    <Icon
                      icon='warning-sign'
                      size={16}
                      color={Colors.RED5}
                      style={{ outline: 'none' }}
                    />
                    <div style={{ padding: '5px' }}>
                      <FormValidationErrors errors={validationErrors} />
                    </div>
                  </Popover>
                )}
              </div>
              <Button
                className='btn-cancel'
                icon='remove'
                text='Cancel'
                onClick={onCancel}
                disabled={isSaving || isTesting}
              />
              <Button
                id='btn-save'
                className='btn-save'
                icon='cloud-upload'
                intent='primary'
                text='Save Connection'
                loading={isSaving}
                disabled={isSaving || isTesting || isLocked || !isValid}
                onClick={onSave}
              />
            </div>
          </div>
        )}
      </form>
    </>
  )
}
