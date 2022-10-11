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
  Tag
} from '@blueprintjs/core'
import { Providers, ProviderLabels } from '@/data/Providers'

const Deployment = (props) => {
  const {
    provider,
    transformation,
    isSaving = false,
    onSettingsChange = () => {}
  } = props

  const [selectValue, setSelectValue] = useState(1)

  useEffect(() => {
    setSelectValue(
      transformation?.deploymentPattern ||
        transformation?.deploymentPattern === '' ||
        transformation?.productionPattern ||
        transformation?.productionPattern === ''
        ? 1
        : 0
    )
  }, [transformation?.deploymentPattern, transformation?.productionPattern])

  const handleChangeSelectValue = (sv) => {
    if (sv === 0) {
      onSettingsChange({
        deploymentPattern: undefined,
        productionPattern: undefined
      })
    } else if (sv === 1) {
      onSettingsChange({ deploymentPattern: '', productionPattern: '' })
    }
    setSelectValue(sv)
  }

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
        tagHint = `A CI job/build with a name that matches the given regEx is considered as a Deployment.`
        break
    }
    return tagHint
  }

  const radioLabels = useMemo(() => {
    let radio1
    let radio2

    const providerId = provider?.id
    const providerName = ProviderLabels[provider?.id?.toUpperCase()]

    switch (providerId) {
      case Providers.JENKINS:
        radio1 = 'Detect Deployment from Builds in Jenkins'
        radio2 = 'Not using Jenkins Builds as Deployments'
        break
      case Providers.GITHUB:
        radio1 = `Detect Deployment from Jobs in GitHub Action`
        radio2 = 'Not using Jobs in GitHub Action as Deployments'
        break
      case Providers.GITLAB:
      default:
        radio1 = `Detect Deployment from Jobs in ${providerName} CI`
        radio2 = `Not using ${providerName} Builds as Deployments`
    }

    return [radio1, radio2]
  }, [provider])

  return (
    <>
      <h5>CI/CD</h5>
      <p style={{ color: '#292B3F' }}>
        <strong>What is a deployment?</strong>{' '}
        <Tag intent={Intent.PRIMARY} style={{ fontSize: '10px' }} minimal>
          DORA
        </Tag>
      </p>
      <p>Define Deployment using one of the following options.</p>
      <RadioGroup
        inline={false}
        label={false}
        name='deploy-tag'
        onChange={(e) => handleChangeSelectValue(+e.target.value)}
        selectedValue={selectValue}
        required
      >
        <Radio label={radioLabels[0]} value={1} />
        {selectValue === 1 && (
          <>
            <p>
              {getDeployTagHint(
                provider?.id,
                ProviderLabels[provider?.id?.toUpperCase()]
              )}
            </p>
            <div className='formContainer'>
              <FormGroup
                disabled={isSaving}
                inline={true}
                label={
                  <label
                    className='bp3-label'
                    style={{ minWidth: '150px', marginRight: '10px' }}
                  >
                    Deployment
                  </label>
                }
                labelFor='deploy-tag-production'
                className='formGroup'
                contentClassName='formGroupContent'
              >
                <InputGroup
                  id='deploy-tag-production'
                  placeholder='(?i)deploy'
                  value={transformation?.deploymentPattern}
                  onChange={(e) =>
                    onSettingsChange({ deploymentPattern: e.target.value })
                  }
                  disabled={isSaving}
                  className='input'
                  maxLength={255}
                />
              </FormGroup>
            </div>
            <p>
              The environment that matches the given regEx is considered as the
              Production environment. If you leave this field empty, all data
              will be tagged as in the Production environment.
            </p>
            <FormGroup
              disabled={isSaving}
              inline={true}
              label={
                <label
                  className='bp3-label'
                  style={{ minWidth: '150px', marginRight: '10px' }}
                >
                  Production
                </label>
              }
              labelFor='production'
              className='formGroup'
              contentClassName='formGroupContent'
            >
              <InputGroup
                id='deploy-tag-production'
                placeholder='(?i)production'
                value={transformation?.productionPattern}
                onChange={(e) =>
                  onSettingsChange({ productionPattern: e.target.value })
                }
                disabled={isSaving}
                className='input'
                maxLength={255}
              />
            </FormGroup>
          </>
        )}
        <Radio label={radioLabels[1]} value={0} />
      </RadioGroup>
    </>
  )
}

export default Deployment
