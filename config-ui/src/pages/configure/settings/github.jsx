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
import {
  FormGroup,
  Checkbox,
  InputGroup,
  NumericInput,
  Tag,
  TextArea,
  Colors,
  Icon,
  Popover,
  Position,
  Intent
} from '@blueprintjs/core'
import { DataDomainTypes } from '@/data/DataDomains'
import Deployment from '@/components/blueprints/transformations/CICD/Deployment'

import '@/styles/integration.scss'
import '@/styles/connections.scss'

export default function GithubSettings(props) {
  const {
    Providers,
    ProviderLabels,
    provider,
    connection,
    dataDomains = [],
    transformation = {},
    isSaving,
    isSavingConnection,
    onSettingsChange = () => {}
  } = props
  const [enableAdditionalCalculations, setEnableAdditionalCalculations] =
    useState(false)

  const handleAdditionalEnable = useCallback(
    (enable) => {
      setEnableAdditionalCalculations(enable)
      onSettingsChange({
        refdiff: enable
          ? { tagsOrder: '', tagsPattern: '', tagsLimit: 10 }
          : null
      })
    },
    [setEnableAdditionalCalculations, onSettingsChange]
  )

  useEffect(() => {
    console.log(
      '>>>> GITHUB: TRANSFORMATION SETTINGS OBJECT....',
      transformation
    )
    setEnableAdditionalCalculations(!!transformation?.refdiff)
  }, [transformation])

  return (
    <>
      {dataDomains.some((e) => e.value === DataDomainTypes.TICKET) && (
        <>
          <h5>
            Issue Tracking{' '}
            <Tag className='bp3-form-helper-text' minimal>
              RegExp
            </Tag>
          </h5>
          <p className=''>
            Map your issue labels with each category to view corresponding
            metrics in the dashboard.
          </p>
          <div style={{}}>
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
                  onChange={(e) =>
                    onSettingsChange({ issueSeverity: e.target.value })
                  }
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
                  onChange={(e) =>
                    onSettingsChange({ issueComponent: e.target.value })
                  }
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
                  onChange={(e) =>
                    onSettingsChange({ issuePriority: e.target.value })
                  }
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
                  onChange={(e) =>
                    onSettingsChange({ issueTypeRequirement: e.target.value })
                  }
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
                  onChange={(e) =>
                    onSettingsChange({ issueTypeBug: e.target.value })
                  }
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
                label={
                  <>
                    Type/Incident
                    <Tag
                      intent={Intent.PRIMARY}
                      style={{ fontSize: '10px', marginLeft: '5px' }}
                      minimal
                    >
                      DORA
                    </Tag>
                  </>
                }
                labelFor='github-issue-incident'
                className='formGroup'
                contentClassName='formGroupContent'
              >
                <InputGroup
                  id='github-issue-incident'
                  placeholder='(incident|p0|p1|p2)$'
                  value={transformation?.issueTypeIncident}
                  onChange={(e) =>
                    onSettingsChange({ issueTypeIncident: e.target.value })
                  }
                  disabled={isSaving || isSavingConnection}
                  className='input'
                  maxLength={255}
                />
              </FormGroup>
            </div>
          </div>
        </>
      )}

      {dataDomains.some((e) => e.value === DataDomainTypes.DEVOPS) && (
        <Deployment
          provider={provider}
          transformation={transformation}
          onSettingsChange={onSettingsChange}
          isSaving={isSaving || isSavingConnection}
        />
      )}

      {dataDomains.some((e) => e.value === DataDomainTypes.CODE_REVIEW) && (
        <>
          <h5>
            Code Review{' '}
            <Tag className='bp3-form-helper-text' minimal>
              RegExp
            </Tag>
          </h5>
          <p className=''>
            Map your pull requests labels with each category to view
            corresponding metrics in the dashboard.
          </p>

          <div style={{}}>
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
                  onChange={(e) => onSettingsChange({ prType: e.target.value })}
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
                  onChange={(e) =>
                    onSettingsChange({ prComponent: e.target.value })
                  }
                  disabled={isSaving || isSavingConnection}
                  className='input'
                  maxLength={255}
                />
              </FormGroup>
            </div>
          </div>

          <h5>
            PR-Issue Mapping{' '}
            <Tag className='bp3-form-helper-text' minimal>
              RegExp
            </Tag>
          </h5>
          <p>
            Extract the issue numbers closed by pull requests. The issue numbers{' '}
            are parsed from PR bodies that meet the following RegEx.
          </p>

          <div className='formContainer'>
            <FormGroup
              disabled={isSaving || isSavingConnection}
              inline={true}
              label={
                <>
                  PR Body Pattern
                  <Popover
                    className='help-pr-body'
                    popoverClassName='popover-pr-body-help'
                    position={Position.TOP}
                    autoFocus={false}
                    enforceFocus={false}
                    usePortal={false}
                  >
                    <Icon
                      icon='help'
                      size={12}
                      color={Colors.GRAY3}
                      style={{ marginLeft: '4px', marginBottom: '4px' }}
                    />
                    <div
                      style={{
                        padding: '10px',
                        width: '300px',
                        maxWidth: '300px',
                        fontSize: '10px'
                      }}
                    >
                      <p style={{ margin: '0 0 10px 0', lineHeight: '110%' }}>
                        <Icon
                          icon='tick-circle'
                          size={10}
                          color={Colors.GREEN4}
                          style={{ marginRight: '4px' }}
                        />
                        Example 1: PR #321 body contains "
                        <strong>Closes #1234</strong>" (PR #321 and issue #1234
                        will be mapped by the following RegEx)
                      </p>
                      <p style={{ margin: 0, lineHeight: '110%' }}>
                        <Icon
                          icon='delete'
                          size={10}
                          color={Colors.RED4}
                          style={{ marginRight: '4px' }}
                        />
                        Example 2: PR #321 body contains "
                        <strong>Related to #1234</strong>" (PR #321 and issue
                        #1234 will NOT be mapped by the following RegEx)
                      </p>
                    </div>
                  </Popover>
                </>
              }
              labelFor='github-pr-body'
              className='formGroup'
              contentClassName='formGroupContent'
              style={{ alignItems: 'center' }}
            >
              <TextArea
                id='github-pr-body'
                className='textarea'
                value={transformation?.prBodyClosePattern}
                // eslint-disable-next-line max-len
                placeholder='(?mi)(fix|close|resolve|fixes|closes|resolves|fixed|closed|resolved)[\s]*.*(((and )?(#|https:\/\/github.com\/%s\/%s\/issues\/)\d+[ ]*)+)'
                onChange={(e) =>
                  onSettingsChange({ prBodyClosePattern: e.target.value })
                }
                disabled={isSaving || isSavingConnection}
                fill
                rows={2}
                growVertically={false}
              />
            </FormGroup>
          </div>

          <h5>Additional Settings</h5>
          <div>
            <Checkbox
              checked={enableAdditionalCalculations}
              label='Enable calculation of commit and issue difference'
              onChange={(e) =>
                handleAdditionalEnable(!enableAdditionalCalculations)
              }
            />
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
                      disabled={isSaving || isSavingConnection}
                      fill={true}
                      placeholder='10'
                      allowNumericCharactersOnly={true}
                      onValueChange={(tagsLimitNumeric) =>
                        onSettingsChange({
                          refdiff: {
                            ...transformation?.refdiff,
                            tagsLimit: tagsLimitNumeric
                          }
                        })
                      }
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
                      onChange={(e) =>
                        onSettingsChange({
                          refdiff: {
                            ...transformation?.refdiff,
                            tagsPattern: e.target.value
                          }
                        })
                      }
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
                      onChange={(e) =>
                        onSettingsChange({
                          refdiff: {
                            ...transformation?.refdiff,
                            tagsOrder: e.target.value
                          }
                        })
                      }
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

      {(dataDomains?.length === 0 ||
        dataDomains.every((e) => e.value === DataDomainTypes.CROSSDOMAIN)) && (
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
