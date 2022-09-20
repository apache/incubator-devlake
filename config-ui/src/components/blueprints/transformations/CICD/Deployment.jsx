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
  Tooltip
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

  const [deployTag, setDeployTag] = useState(transformation?.productionPattern || '')
  const [enableDeployTag, setEnableDeployTag] = useState([
    transformation?.productionPattern,
    // transformation?.stagingPattern,
    // transformation?.testingPattern
  ].some(t => t && t !== '') ? 1 : 0)

  // @todo: check w/ product team about using standard message and avoid customized hints
  const getDeployTagHint = (providerId, providerName = 'Plugin') => {
    let tagHint = ''
    switch (providerId) {
      case Providers.JENKINS:
        // eslint-disable-next-line max-len
        tagHint = `The ${providerName} build with a name that matches the given regEx is considered as a deployment. You can define your Deployments for three environments: Production, Staging and Testing.`
        break
      case Providers.GITHUB:
      case Providers.GITLAB:
      case 'default':
        // eslint-disable-next-line max-len
        tagHint = 'A CI job/build with a name that matches the given regEx is considered as an deployment. You can define your Deployments for three environments: Production, Staging and Testing.'
        break
    }
    return tagHint
  }

  const getDeployOptionLabel = (providerId, providerName) => {
    let label = ''
    switch (providerId) {
      case Providers.JENKINS:
        // eslint-disable-next-line max-len
        label = `Detect Deployment from Builds in ${providerName}`
        break
      case Providers.GITHUB:
        label = `Detect Deployment from Jobs in ${providerName} Action`
        break
      case Providers.GITLAB:
      case 'default':
        // eslint-disable-next-line max-len
        label = `Detect Deployment from Jobs in ${providerName} CI`
        break
    }
    return label
  }

  useEffect(() => {
    console.log('>>> CI/CD Deployment Transform:', entityIdKey, deployTag)
    if (entityIdKey && enableDeployTag === 0) {
      onSettingsChange({ productionPattern: '' }, entityIdKey)
      // onSettingsChange({ stagingPattern: '' }, entityIdKey)
      // onSettingsChange({ testingPattern: '' }, entityIdKey)
    }
  }, [
    deployTag,
    enableDeployTag,
    entityIdKey,
    onSettingsChange
  ])

  useEffect(() => {
    console.log('>>> CI/CD Deployment: TRANSFORMATION OBJECT!', transformation)
  }, [transformation])

  return (
    <>
      <h5>CI/CD <Tag className='bp3-form-helper-text' minimal>RegExp</Tag></h5>
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
          label={getDeployOptionLabel(provider?.id, ProviderLabels[provider?.id?.toUpperCase()])}
          value={1}
        />
        {enableDeployTag === 1 && (
          <>
            <div
              className='bp3-form-helper-text'
              style={{ display: 'block', textAlign: 'left', color: '#94959F', marginBottom: '5px' }}
            >
              {getDeployTagHint(provider?.id, ProviderLabels[provider?.id?.toUpperCase()])}
            </div>
            <div className='formContainer'>
              <FormGroup
                disabled={isSaving}
                inline={true}
                label={<label className='bp3-label' style={{ minWidth: '150px', marginRight: '10px' }}>Deployment (Production)</label>}
                labelFor='deploy-tag-production'
                className='formGroup'
                contentClassName='formGroupContent'
              >
                <InputGroup
                  id='deploy-tag-production'
                  placeholder='(?i)deploy'
                  value={transformation?.productionPattern}
                  onChange={(e) => onSettingsChange({ productionPattern: e.target.value }, entityIdKey)}
                  disabled={isSaving}
                  className='input'
                  maxLength={255}
                  rightElement={
                    enableDeployTag && (transformation?.productionPattern === '' || !transformation?.productionPattern)
                      ? (
                        <Tooltip intent={Intent.PRIMARY} content='Deployment Tag RegEx required'>
                          <Icon icon='warning-sign' color={Colors.GRAY3} size={12} style={{ margin: '8px' }} />
                        </Tooltip>
                        )
                      : null
                  }
                  required
                />
              </FormGroup>
            </div>
            {/* <div className='formContainer'>
              <FormGroup
                disabled={isSaving}
                inline={true}
                label={<label className='bp3-label' style={{ minWidth: '150px', marginRight: '10px' }}>Deployment (Staging)</label>}
                labelFor='deploy-tag-staging'
                className='formGroup'
                contentClassName='formGroupContent'
              >
                <InputGroup
                  id='deploy-tag-staging'
                  placeholder='(?i)stag'
                  value={transformation?.stagingPattern}
                  onChange={(e) => onSettingsChange({ stagingPattern: e.target.value }, entityIdKey)}
                  disabled={isSaving}
                  className='input'
                  maxLength={255}
                  rightElement={
                    enableDeployTag && (transformation?.stagingPattern === '' || !transformation?.stagingPattern)
                      ? (
                        <Tooltip intent={Intent.PRIMARY} content='Deployment Tag RegEx required'>
                          <Icon icon='warning-sign' color={Colors.GRAY3} size={12} style={{ margin: '8px' }} />
                        </Tooltip>
                        )
                      : null
                  }
                  required
                />
              </FormGroup>
            </div>
            <div className='formContainer'>
              <FormGroup
                disabled={isSaving}
                inline={true}
                label={<label className='bp3-label' style={{ minWidth: '150px', marginRight: '10px' }}>Deployment (Testing)</label>}
                labelFor='deploy-tag-testing'
                className='formGroup'
                contentClassName='formGroupContent'
              >
                <InputGroup
                  id='deploy-tag-testing'
                  placeholder='(?i)test'
                  value={transformation?.testingPattern}
                  onChange={(e) => onSettingsChange({ testingPattern: e.target.value }, entityIdKey)}
                  disabled={isSaving}
                  className='input'
                  maxLength={255}
                  rightElement={
                    enableDeployTag && (transformation?.testingPattern === '' || !transformation?.testingPattern)
                      ? (
                        <Tooltip intent={Intent.PRIMARY} content='Deployment Tag RegEx required'>
                          <Icon icon='warning-sign' color={Colors.GRAY3} size={12} style={{ margin: '8px' }} />
                        </Tooltip>
                        )
                      : null
                  }
                  required
                />
              </FormGroup>
            </div> */}
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
