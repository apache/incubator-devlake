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
import React, { useEffect, useState } from 'react'
import {
  Button,
  FormGroup, InputGroup, Label,
  Tag,
  Intent
} from '@blueprintjs/core'
import { Buffer } from 'buffer'
import '@/styles/integration.scss'
import '@/styles/connections.scss'

export default function GenerateTokenForm (props) {
  const {
    isTesting,
    isSaving,
    isLocked,
    onTokenChange,
    setShowTokenCreator
  } = props
  const [generatorUsername, setGeneratorUsername] = useState('')
  const [generatorPassword, setGeneratorPassword] = useState('')
  const [newToken, setNewToken] = useState()

  const generateAuthToken = (username, password) => {
    const token = Buffer.from(`${username}:${password}`).toString('base64')
    onTokenChange(token)
    setNewToken(token)
    console.log('>> BASIC AUTH TOKEN ENCODED = ', token)
    setShowTokenCreator(false)
  }

  const resetTokenGenerator = () => {
    setGeneratorUsername('')
    setGeneratorPassword('')
  }

  useEffect(() => {
    // --token-data-changed
  }, [generatorUsername, generatorPassword, newToken])

  return (
    <>
      <h3 style={{ margin: 0 }}>GENERATE TOKEN <Tag>base64</Tag></h3>
      <p style={{ margin: '0 0 10px 0' }}>Enter <strong>Username</strong> (or E-mail) and <strong>Password</strong></p>
      <div className='formContainer' style={{ marginBottom: '0.2rem' }}>
        <FormGroup
          label=''
          disabled={isTesting || isSaving || isLocked}
          inline={true}
          labelFor='token-username'
          className='formGroup'
          contentClassName='formGroupContent'
        >
          <Label style={{ display: 'inline', minWidth: '50px', whiteSpace: 'nowrap' }}>
            Username <span className='requiredStar'>*</span>
          </Label>
          <InputGroup
            id='token-username'
            disabled={isTesting || isSaving || isLocked}
            placeholder='Enter Username'
            value={generatorUsername}
            onChange={(e) => setGeneratorUsername(e.target.value)}
            className='input'
            style={{ maxWidth: '300px' }}
          />
        </FormGroup>
      </div>
      <div className='formContainer' style={{ marginBottom: '0.2rem' }}>
        <FormGroup
          disabled={isTesting || isSaving || isLocked}
          label=''
          inline={true}
          labelFor='token-password'
          className='formGroup'
          contentClassName='formGroupContent'
        >
          <Label style={{ display: 'inline', minWidth: '50px', whiteSpace: 'nowrap' }}>
            Password <span className='requiredStar'>*</span>
          </Label>
          <InputGroup
            id='token-password'
            type='password'
            disabled={isTesting || isSaving || isLocked}
            placeholder='Enter Password'
            value={generatorPassword}
            onChange={(e) => setGeneratorPassword(e.target.value)}
            className='input'
            style={{ maxWidth: '300px' }}
          />
        </FormGroup>
      </div>
      <div style={{ display: 'flex' }}>
        <Button
          type='button' icon='eraser' text=''
          style={{ display: 'flex', marginLeft: 'auto', marginRight: '5px' }}
          onClick={resetTokenGenerator} small minimal
        />
        <Button
          type='button' intent={Intent.PRIMARY} icon='random' text='Generate'
          style={{ display: 'flex' }}
          disabled={!generatorUsername || !generatorPassword}
          onClick={() => generateAuthToken(generatorUsername, generatorPassword)} small
        />
      </div>
    </>
  )
}
