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
import React, { useState, useEffect, useMemo } from 'react'
import {
  Intent,
  FormGroup,
  RadioGroup,
  InputGroup,
  Radio,
  Tag,
  Icon,
  Colors,
} from '@blueprintjs/core'
import { Providers, ProviderLabels } from '@/data/Providers'
const Deployment = (props) => {
  const {
    connection,
    provider,
    transformation,
    configuredProject,
    configuredBoard,
    entityIdKey,
    entities = [],
    isSaving = false,
    onSettingsChange = () => {},
  } = props

  const [deployTag, setDeployTag] = useState(transformation?.deployTagPattern || '')
  const [enableDeployTag, setEnableDeployTag] = useState(transformation?.deployTagPattern !== '' ? 1 : 0)

  const getDeployTagHint = (providerId) => {
    let tagHint = ''
    switch (providerId) {
      case Providers.JENKINS:
        tagHint = 'The Jenkins build with a name that matches the given regEx is considered as a deployment.'
        break
      case Providers.GITLAB:
      case 'default':
        tagHint = 'A CI job/build with a name that matches the given regEx is considered as an deployment.'
        break
    }
    return tagHint
  }

  useEffect(() => {
    console.log('CI/CD Deployment Transform:', entityIdKey, deployTag)
    if (entityIdKey && enableDeployTag === 1) {
      onSettingsChange({ deployTagPattern: deployTag }, entityIdKey)
    } else if (entityIdKey && enableDeployTag === 0) {
      onSettingsChange({ deployTagPattern: '' }, entityIdKey)
    }
  }, [
    deployTag,
    enableDeployTag,
    entityIdKey,
    onSettingsChange
  ])

  return (
    <>
      <h5>CI/CD</h5>
      <p>Define deployment using one of the followng options</p>
      <p style={{ color: '#292B3F' }}>
        <strong>What is a deployment?</strong>{' '}
        <Tag intent={Intent.PRIMARY} style={{ fontSize: '10px' }} minimal>
          DORA
        </Tag>
      </p>

      <RadioGroup
        inline={false}
        label={false}
        name='deploy-tag'
        onChange={(e) => setEnableDeployTag(Number(e.target.value))}
        selectedValue={enableDeployTag}
        required
      >
        <Radio
          label={`Detect Deployment from Builds in ${
            ProviderLabels[provider?.id.toUpperCase()]
          }`}
          value={1}
        />
        {enableDeployTag === 1 && (
          <>
            <div
              className='bp3-form-helper-text'
              style={{ display: 'block', textAlign: 'left', color: '#94959F', marginBottom: '5px' }}
            >
              {getDeployTagHint(provider?.id)}
            </div>
            <div className='formContainer'>
              <FormGroup
                disabled={isSaving}
                inline={true}
                label='Deployment'
                labelFor='deploy-tag'
                className='formGroup'
                contentClassName='formGroupContent'
              >
                <InputGroup
                  id='github-pr-type'
                  placeholder='/deploy/'
                  value={transformation?.deployTagPattern}
                  onChange={(e) => setDeployTag(e.target.value)}
                  disabled={isSaving}
                  className='input'
                  maxLength={255}
                />
              </FormGroup>
            </div>
          </>
        )}
        <Radio
          label={`Not using ${
            ProviderLabels[provider?.id.toUpperCase()]
          } Builds as Deployments`}
          value={0}
        />
      </RadioGroup>
    </>
  )
}

export default Deployment
