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
import React, { useEffect, useState, useCallback } from 'react'
import { FormGroup, Checkbox, InputGroup, NumericInput, Tag } from '@blueprintjs/core'
import { DataEntityTypes } from '@/data/DataEntities'

import '@/styles/integration.scss'
import '@/styles/connections.scss'

export default function GithubSettings (props) {
  const {
    // eslint-disable-next-line no-unused-vars
    connection,
    entities = [],
    transformation = {},
    isSaving,
    isSavingConnection,
    onSettingsChange = () => {},
    configuredProject
  } = props

  // eslint-disable-next-line no-unused-vars
  const [errors, setErrors] = useState([])
  const [enableAdditionalCalculations, setEnableAdditionalCalculations] = useState(false)

  // eslint-disable-next-line no-unused-vars
  const handleSettingsChange = useCallback((setting) => {
    onSettingsChange(setting, configuredProject)
  }, [onSettingsChange, configuredProject])

  const handleAdditionalSettings = useCallback((setting) => {
    setEnableAdditionalCalculations(setting)
    onSettingsChange({ refdiff: setting ? { tagsOrder: '', tagsPattern: '', tagsLimit: 10, } : null }, configuredProject)
  }, [setEnableAdditionalCalculations, configuredProject, onSettingsChange])

  useEffect(() => {
    console.log('>>>> TRANSFORMATION SETTINGS OBJECT....', transformation)
    setEnableAdditionalCalculations(!!transformation?.refdiff)
  }, [transformation])

  useEffect(() => {
    console.log('>>>> ENABLE GITHUB ADDITIONAL SETTINGS..?', enableAdditionalCalculations)
    if (enableAdditionalCalculations === 'disabled') {
      // onSettingsChange({gitextractorCalculation: ''}, configuredProject)
    }
  }, [enableAdditionalCalculations])

  return (
    <>
      {entities.some(e => e.value === DataEntityTypes.TICKET) && (
        <><h5>Issue Tracking{' '} <Tag className='bp3-form-helper-text'>RegExp</Tag></h5>
          <p className=''>Map your issue labels with each category
            to view corresponding metrics in the
            dashboard.
          </p>
          <div style={{ }}>
            <div className='formContainer'>
              <FormGroup
                disabled={isSaving || isSavingConnection}
                inline={true}
                label='Severity'
                labelFor='github-issue-severity'
                className='formGroup'
                contentClassName='formGroupContent'
              >
                <InputGroup
                  id='github-issue-severity'
                  placeholder='severity/(.*)$'
                // defaultValue={transformation?.issueSeverity}
                  value={transformation?.issueSeverity}
                  onChange={(e) => onSettingsChange({ issueSeverity: e.target.value }, configuredProject)}
                  disabled={isSaving || isSavingConnection}
                  className='input'
                  maxLength={255}
                  autoFocus={true}
                />
              </FormGroup>
            </div>
            <div className='formContainer'>
              <FormGroup
                disabled={isSaving || isSavingConnection}
                inline={true}
                label='Component'
                labelFor='github-issue-component'
                className='formGroup'
                contentClassName='formGroupContent'
              >
                <InputGroup
                  id='github-issue-component'
                  placeholder='component/(.*)$'
                  value={transformation?.issueComponent}
                  onChange={(e) => onSettingsChange({ issueComponent: e.target.value }, configuredProject)}
                  disabled={isSaving || isSavingConnection}
                  className='input'
                  maxLength={255}
                />
              </FormGroup>
            </div>
            <div className='formContainer'>
              <FormGroup
                disabled={isSaving || isSavingConnection}
                inline={true}
                label='Priority'
                labelFor='github-issue-priority'
                className='formGroup'
                contentClassName='formGroupContent'
              >
                <InputGroup
                  id='github-issue-priority'
                  placeholder='(highest|high|medium|low)$'
                  value={transformation?.issuePriority}
                  onChange={(e) => onSettingsChange({ issuePriority: e.target.value }, configuredProject)}
                  disabled={isSaving || isSavingConnection}
                  className='input'
                  maxLength={255}
                />
              </FormGroup>
            </div>
            <div className='formContainer'>
              <FormGroup
                disabled={isSaving || isSavingConnection}
                inline={true}
                label='Type/Requirement'
                labelFor='github-issue-requirement'
                className='formGroup'
                contentClassName='formGroupContent'
              >
                <InputGroup
                  id='github-issue-requirement'
                  placeholder='(feat|feature|proposal|requirement)$'
                  value={transformation?.issueTypeRequirement}
                  onChange={(e) => onSettingsChange({ issueTypeRequirement: e.target.value }, configuredProject)}
                  disabled={isSaving || isSavingConnection}
                  className='input'
                  maxLength={255}
                />
              </FormGroup>
            </div>
            <div className='formContainer'>
              <FormGroup
                disabled={isSaving || isSavingConnection}
                inline={true}
                label='Type/Bug'
                labelFor='github-issue-bug'
                className='formGroup'
                contentClassName='formGroupContent'
              >
                <InputGroup
                  id='github-issue-bug'
                  placeholder='(bug|broken)$'
                  value={transformation?.issueTypeBug}
                  onChange={(e) => onSettingsChange({ issueTypeBug: e.target.value }, configuredProject)}
                  disabled={isSaving || isSavingConnection}
                  className='input'
                  maxLength={255}
                />
              </FormGroup>
            </div>
            <div className='formContainer'>
              <FormGroup
                disabled={isSaving || isSavingConnection}
                inline={true}
                label='Type/Incident'
                labelFor='github-issue-incident'
                className='formGroup'
                contentClassName='formGroupContent'
              >
                <InputGroup
                  id='github-issue-incident'
                  placeholder='(incident|p0|p1|p2)$'
                  value={transformation?.issueTypeIncident}
                  onChange={(e) => onSettingsChange({ issueTypeIncident: e.target.value }, configuredProject)}
                  disabled={isSaving || isSavingConnection}
                  className='input'
                  maxLength={255}
                />
              </FormGroup>
            </div>
          </div>
        </>
      )}

      {entities.some(e => e.value === DataEntityTypes.CODE_REVIEW) && (
        <><h5>Code Review{' '} <Tag className='bp3-form-helper-text'>RegExp</Tag></h5>
          <p className=''>Map your pull requests labels with each category to view corresponding metrics in the dashboard.</p>

          <div style={{ }}>
            <div className='formContainer'>
              <FormGroup
                disabled={isSaving || isSavingConnection}
                inline={true}
                label='Type'
                labelFor='github-pr-type'
                className='formGroup'
                contentClassName='formGroupContent'
              >
                <InputGroup
                  id='github-pr-type'
                  placeholder='type/(.*)$'
                  value={transformation?.prType}
                  onChange={(e) => onSettingsChange({ prType: e.target.value }, configuredProject)}
                  disabled={isSaving || isSavingConnection}
                  className='input'
                  maxLength={255}
                />
              </FormGroup>
            </div>
            <div className='formContainer'>
              <FormGroup
                disabled={isSaving || isSavingConnection}
                inline={true}
                label='Component'
                labelFor='github-pr-component'
                className='formGroup'
                contentClassName='formGroupContent'
              >
                <InputGroup
                  id='github-pr-type'
                  placeholder='component/(.*)$'
                  value={transformation?.prComponent}
                  onChange={(e) => onSettingsChange({ prComponent: e.target.value }, configuredProject)}
                  disabled={isSaving || isSavingConnection}
                  className='input'
                  maxLength={255}
                />
              </FormGroup>
            </div>
          </div>

          <h5>Additional Settings</h5>
          <div>
            <Checkbox checked={enableAdditionalCalculations} label='Enable calculation of commit and issue difference' onChange={(e) => handleAdditionalSettings(!enableAdditionalCalculations)} />
            {enableAdditionalCalculations && (
            <>
              <div className='additional-settings-refdiff'>
                <FormGroup
                  disabled={isSaving || isSavingConnection}
                  inline={true}
                  label='Tags Limit'
                  className='formGroup'
                  contentClassName='formGroupContent'
                >
                  <NumericInput
                    id='refdiff-tags-limit'
                    fill={true}
                    placeholder='10'
                    allowNumericCharactersOnly={true}
                    onValueChange={(tagsLimitNumeric) => onSettingsChange({ refdiff: { ...transformation?.refdiff, tagsLimit: tagsLimitNumeric } }, configuredProject)}
                    value={transformation?.refdiff?.tagsLimit}
                  />
                </FormGroup>
                <FormGroup
                  disabled={isSaving || isSavingConnection}
                  inline={true}
                  label='Tags Pattern'
                  className='formGroup'
                  contentClassName='formGroupContent'
                >
                  <InputGroup
                    id='refdiff-tags-pattern'
                    placeholder='(regex)$'
                    value={transformation?.refdiff?.tagsPattern}
                    onChange={(e) => onSettingsChange({ refdiff: { ...transformation?.refdiff, tagsPattern: e.target.value } }, configuredProject)}
                    disabled={isSaving || isSavingConnection}
                    className='input'
                    maxLength={255}
                  />
                </FormGroup>
                <FormGroup
                  disabled={isSaving || isSavingConnection}
                  inline={true}
                  label='Tags Order'
                  className='formGroup'
                  contentClassName='formGroupContent'
                >
                  <InputGroup
                    id='refdiff-tags-order'
                    placeholder='reverse semver'
                    value={transformation?.refdiff?.tagsOrder}
                    onChange={(e) => onSettingsChange({ refdiff: { ...transformation?.refdiff, tagsOrder: e.target.value } }, configuredProject)}
                    disabled={isSaving || isSavingConnection}
                    className='input'
                    maxLength={255}
                  />
                </FormGroup>
              </div>

            </>
            )}
          </div>
        </>
      )}

      {(entities?.length === 0 || entities.some(e => e.value === DataEntityTypes.CROSSDOMAIN)) && (
        <div className='headlineContainer'>
          <h5>No Data Entities</h5>
          <p className='description'>
            You have not selected entities that require configuration.
          </p>
        </div>
      )}

    </>
  )
}
