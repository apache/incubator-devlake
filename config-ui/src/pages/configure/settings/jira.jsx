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
import React, { useEffect, useState, useMemo } from 'react'
import {
  Button,
  ButtonGroup,
  Colors,
  FormGroup,
  Icon,
  InputGroup,
  Intent,
  Label,
  MenuItem,
  Position,
  Tag
} from '@blueprintjs/core'
import { MultiSelect, Select } from '@blueprintjs/select'
import { DataDomainTypes } from '@/data/DataDomains'

import '@/styles/integration.scss'
import '@/styles/connections.scss'

const MAPPING_TYPES = {
  Requirement: 'Requirement',
  Incident: 'Incident',
  Bug: 'Bug'
}

const createTypeMapObject = (customType, standardType) => {
  return customType && standardType
    ? {
        [customType]: {
          standardType
        }
      }
    : null
}

export default function JiraSettings(props) {
  const {
    Providers,
    ProviderLabels,
    provider,
    connection,
    issueTypes = [],
    fields: fieldsList = [],
    transformation = {},
    onSettingsChange = () => {},
    dataDomains = [],
    isSaving,
    isSavingConnection = false,
    jiraProxyError,
    isFetchingJIRA = false
  } = props

  // eslint-disable-next-line no-unused-vars
  const [statusMappings, setStatusMappings] = useState()
  const [jiraIssueEpicKeyField, setJiraIssueEpicKeyField] = useState('')
  const [jiraIssueStoryPointField, setJiraIssueStoryPointField] = useState('')
  const [remoteLinkCommitSha, setRemoteLinkCommitSha] = useState('')

  const [requirementTags, setRequirementTags] = useState([])
  const [bugTags, setBugTags] = useState([])
  const [incidentTags, setIncidentTags] = useState([])
  const allChosenTagsInThisBoard = useMemo(
    () => [
      ...(Array.isArray(requirementTags) ? requirementTags : []),
      ...(Array.isArray(bugTags) ? bugTags : []),
      ...(Array.isArray(incidentTags) ? incidentTags : [])
    ],
    [requirementTags, bugTags, incidentTags]
  )

  useEffect(() => {
    console.log('>>> JIRA SETTINGS :: FIELDS LIST DATA CHANGED!', fieldsList)
  }, [fieldsList])

  useEffect(() => {
    setBugTags(transformation?.bugTags || [])
  }, [transformation?.bugTags])

  useEffect(() => {
    setIncidentTags(transformation?.incidentTags || [])
  }, [transformation?.incidentTags])

  useEffect(() => {
    setRequirementTags(transformation?.requirementTags || [])
  }, [transformation?.requirementTags])

  useEffect(() => {
    setJiraIssueEpicKeyField(
      fieldsList.find((f) => f.value === transformation?.epicKeyField)
    )
  }, [fieldsList, transformation?.epicKeyField])

  useEffect(() => {
    setJiraIssueStoryPointField(
      fieldsList.find((f) => f.value === transformation?.storyPointField)
    )
  }, [fieldsList, transformation?.storyPointField])

  useEffect(() => {
    setRemoteLinkCommitSha(transformation?.remotelinkCommitShaPattern || '')
  }, [transformation?.remotelinkCommitShaPattern])

  useEffect(() => {
    console.log('>>> JIRA SETTINGS :: TRANSFORMATION OBJECT!', transformation)
  }, [transformation])

  const typeMappingAll = useMemo(
    () =>
      [
        ...bugTags.map((r) => createTypeMapObject(r.value, MAPPING_TYPES.Bug)),
        ...(incidentTags || []).map((r) =>
          createTypeMapObject(r.value, MAPPING_TYPES.Incident)
        ),
        ...requirementTags.map((r) =>
          createTypeMapObject(r.value, MAPPING_TYPES.Requirement)
        )
      ].reduce((c, p) => ({ ...c, ...p }), {}),
    [bugTags, incidentTags, requirementTags]
  )

  useEffect(() => {
    onSettingsChange({ typeMappings: typeMappingAll })
  }, [typeMappingAll, onSettingsChange])

  return (
    <>
      {dataDomains.some((e) => e.value === DataDomainTypes.TICKET) && (
        <>
          <h5>Issue Tracking</h5>
          <p className=''>
            Map your issue labels with each category to view corresponding
            metrics in the dashboard.
          </p>

          <div
            className='issue-type-multiselect'
            style={{ display: 'flex', marginBottom: '10px' }}
          >
            <div
              className='issue-type-label'
              style={{
                minWidth: '120px',
                paddingRight: '10px',
                paddingTop: '3px'
              }}
            >
              <label>Requirement</label>
            </div>
            <div
              className='issue-type-multiselect-selector'
              style={{ minWidth: '200px', width: '100%' }}
            >
              <MultiSelect
                disabled={isSaving}
                resetOnSelect={true}
                placeholder='Select...'
                popoverProps={{
                  usePortal: false,
                  popoverClassName: 'transformation-select-popover',
                  minimal: true,
                  fill: true,
                  style: { width: '100%' }
                }}
                className='multiselector-requirement-type'
                inline={true}
                fill={true}
                items={issueTypes}
                selectedItems={requirementTags}
                activeItem={null}
                itemPredicate={(query, item) =>
                  item?.title.toLowerCase().indexOf(query.toLowerCase()) >= 0
                }
                itemRenderer={(item, { handleClick, modifiers }) => (
                  <MenuItem
                    active={modifiers.active}
                    disabled={allChosenTagsInThisBoard?.some(
                      (t) => t.value === item.value
                    )}
                    label={
                      <span style={{ marginLeft: '20px' }}>
                        {item.description}
                      </span>
                    }
                    key={item.value}
                    onClick={handleClick}
                    text={
                      requirementTags.some((t) => t.value === item.value) ? (
                        <>
                          <img src={item.icon} width={12} height={12} />{' '}
                          {item.title}{' '}
                          <Icon icon='small-tick' color={Colors.GREEN5} />
                        </>
                      ) : (
                        <span style={{ fontWeight: 700 }}>
                          <img src={item.icon} width={12} height={12} />{' '}
                          {item.title}
                        </span>
                      )
                    }
                    style={{
                      marginBottom: '2px',
                      fontWeight: requirementTags.some(
                        (t) => t.value === item.value
                      )
                        ? 700
                        : 'normal'
                    }}
                  />
                )}
                tagRenderer={(item) => item.title}
                tagInputProps={{
                  tagProps: {
                    intent: Intent.NONE,
                    color: Colors.RED3,
                    minimal: true
                  }
                }}
                noResults={<MenuItem disabled={true} text='No results.' />}
                onRemove={(item) => {
                  const newValue = requirementTags.filter(
                    (t) => t.id !== item.id
                  )
                  setRequirementTags(newValue)
                  onSettingsChange({ requirementTags: newValue })
                }}
                onItemSelect={(item) => {
                  const newValue = !requirementTags.includes(item)
                    ? [...requirementTags, item]
                    : [...requirementTags]
                  setRequirementTags(newValue)
                  onSettingsChange({ requirementTags: newValue })
                }}
              />
            </div>
            <div
              className='multiselect-clear-action'
              style={{ marginLeft: '0' }}
            >
              <Button
                icon='eraser'
                disabled={requirementTags?.length === 0 || isSaving}
                intent={Intent.NONE}
                minimal={false}
                onClick={() => {
                  setRequirementTags([])
                  onSettingsChange({ requirementTags: [] })
                }}
                style={{
                  borderTopLeftRadius: 0,
                  borderBottomLeftRadius: 0,
                  marginLeft: '-2px'
                }}
              />
            </div>
          </div>

          <div
            className='issue-type-multiselect'
            style={{ display: 'flex', marginBottom: '10px' }}
          >
            <div
              className='issue-type-label'
              style={{
                minWidth: '120px',
                paddingRight: '10px',
                paddingTop: '3px'
              }}
            >
              <label>Bug</label>
            </div>
            <div
              className='issue-type-multiselect-selector'
              style={{ minWidth: '200px', width: '100%' }}
            >
              <MultiSelect
                disabled={isSaving}
                resetOnSelect={true}
                placeholder='Select...'
                popoverProps={{
                  usePortal: false,
                  popoverClassName: 'transformation-select-popover',
                  minimal: true
                }}
                className='multiselector-bug-type'
                inline={true}
                fill={true}
                items={issueTypes}
                selectedItems={bugTags}
                activeItem={null}
                itemPredicate={(query, item) =>
                  item?.title.toLowerCase().indexOf(query.toLowerCase()) >= 0
                }
                itemRenderer={(item, { handleClick, modifiers }) => (
                  <MenuItem
                    active={modifiers.active}
                    disabled={allChosenTagsInThisBoard?.some(
                      (t) => t.value === item.value
                    )}
                    label={
                      <span style={{ marginLeft: '20px' }}>
                        {item.description}
                      </span>
                    }
                    key={item.value}
                    onClick={handleClick}
                    text={
                      bugTags.some((t) => t.value === item.value) ? (
                        <>
                          <img src={item.icon} width={12} height={12} />{' '}
                          {item.title}{' '}
                          <Icon icon='small-tick' color={Colors.GREEN5} />
                        </>
                      ) : (
                        <span style={{ fontWeight: 700 }}>
                          <img src={item.icon} width={12} height={12} />{' '}
                          {item.title}
                        </span>
                      )
                    }
                    style={{
                      marginBottom: '2px',
                      fontWeight: bugTags.some((t) => t.value === item.value)
                        ? 700
                        : 'normal'
                    }}
                  />
                )}
                tagRenderer={(item) => item.title}
                tagInputProps={{
                  tagProps: {
                    intent: Intent.NONE,
                    color: Colors.RED3,
                    minimal: true
                  }
                }}
                noResults={<MenuItem disabled={true} text='No results.' />}
                onRemove={(item) => {
                  const newValue = bugTags.filter((t) => t.id !== item.id)
                  setBugTags(newValue)
                  onSettingsChange({ bugTags: newValue })
                }}
                onItemSelect={(item) => {
                  const newValue = !bugTags.includes(item)
                    ? [...bugTags, item]
                    : [...bugTags]
                  setBugTags(newValue)
                  onSettingsChange({ bugTags: newValue })
                }}
              />
            </div>
            <div
              className='multiselect-clear-action'
              style={{ marginLeft: '0' }}
            >
              <Button
                icon='eraser'
                disabled={bugTags.length === 0 || isSaving}
                intent={Intent.NONE}
                minimal={false}
                onClick={() => {
                  setBugTags([])
                  onSettingsChange({ bugTags: [] })
                }}
                style={{
                  borderTopLeftRadius: 0,
                  borderBottomLeftRadius: 0,
                  marginLeft: '-2px'
                }}
              />
            </div>
          </div>

          <div
            className='issue-type-multiselect'
            style={{ display: 'flex', marginBottom: '10px' }}
          >
            <div
              className='issue-type-label'
              style={{
                minWidth: '120px',
                paddingRight: '10px',
                paddingTop: '3px'
              }}
            >
              <label>
                Incident{' '}
                <Tag
                  intent={Intent.PRIMARY}
                  style={{ fontSize: '10px' }}
                  minimal
                >
                  DORA
                </Tag>
              </label>
            </div>
            <div
              className='issue-type-multiselect-selector'
              style={{ minWidth: '200px', width: '100%' }}
            >
              <MultiSelect
                disabled={isSaving}
                resetOnSelect={true}
                placeholder='Select...'
                popoverProps={{
                  usePortal: false,
                  popoverClassName: 'transformation-select-popover',
                  minimal: true
                }}
                className='multiselector-incident-type'
                inline={true}
                fill={true}
                items={issueTypes}
                selectedItems={incidentTags}
                activeItem={null}
                itemPredicate={(query, item) =>
                  item?.title.toLowerCase().indexOf(query.toLowerCase()) >= 0
                }
                itemRenderer={(item, { handleClick, modifiers }) => (
                  <MenuItem
                    active={modifiers.active}
                    disabled={allChosenTagsInThisBoard?.some(
                      (t) => t.value === item.value
                    )}
                    label={
                      <span style={{ marginLeft: '20px' }}>
                        {item.description}
                      </span>
                    }
                    key={item.value}
                    onClick={handleClick}
                    text={
                      incidentTags.some((t) => t.value === item.value) ? (
                        <>
                          <img src={item.icon} width={12} height={12} />{' '}
                          {item.title}{' '}
                          <Icon icon='small-tick' color={Colors.GREEN5} />
                        </>
                      ) : (
                        <span style={{ fontWeight: 700 }}>
                          <img src={item.icon} width={12} height={12} />{' '}
                          {item.title}
                        </span>
                      )
                    }
                    style={{
                      marginBottom: '2px',
                      fontWeight: incidentTags.some(
                        (t) => t.value === item.value
                      )
                        ? 700
                        : 'normal'
                    }}
                  />
                )}
                tagRenderer={(item) => item.title}
                tagInputProps={{
                  tagProps: {
                    intent: Intent.NONE,
                    color: Colors.RED3,
                    minimal: true
                  }
                }}
                noResults={<MenuItem disabled={true} text='No results.' />}
                onRemove={(item) => {
                  const newValue = incidentTags.filter((t) => t.id !== item.id)
                  setIncidentTags(newValue)
                  onSettingsChange({ incidentTags: newValue })
                }}
                onItemSelect={(item) => {
                  const newValue = !incidentTags.includes(item)
                    ? [...incidentTags, item]
                    : incidentTags
                  setIncidentTags(newValue)
                  onSettingsChange({ incidentTags: newValue })
                }}
              />
            </div>
            <div
              className='multiselect-clear-action'
              style={{ marginLeft: '0' }}
            >
              <Button
                icon='eraser'
                disabled={incidentTags.length === 0 || isSaving}
                intent={Intent.NONE}
                minimal={false}
                onClick={() => {
                  setIncidentTags([])
                  onSettingsChange({ incidentTags: [] })
                }}
                style={{
                  borderTopLeftRadius: 0,
                  borderBottomLeftRadius: 0,
                  marginLeft: '-2px'
                }}
              />
            </div>
          </div>

          <div
            className='epic-key-select'
            style={{ display: 'flex', marginBottom: '10px' }}
          >
            <div
              className='epick-key-label'
              style={{
                minWidth: '120px',
                paddingRight: '10px',
                paddingTop: '3px'
              }}
            >
              <label>Epic Link</label>
            </div>
            <div style={{ display: 'flex', minWidth: '260px' }}>
              <ButtonGroup disabled={isSaving} fill={true}>
                <Select
                  disabled={isSaving || fieldsList.length === 0}
                  className='select-epic-key'
                  inline={true}
                  fill={true}
                  items={fieldsList}
                  activeItem={jiraIssueEpicKeyField}
                  itemPredicate={(query, item) =>
                    item?.title.toLowerCase().indexOf(query.toLowerCase()) >= 0
                  }
                  itemRenderer={(item, { handleClick, modifiers }) => (
                    <MenuItem
                      disabled={jiraIssueStoryPointField?.value === item.value}
                      active={false}
                      intent={modifiers.active ? Intent.NONE : Intent.NONE}
                      key={item.value}
                      label={item.value}
                      onClick={handleClick}
                      text={
                        <>
                          <span>{item.title}</span>{' '}
                          <Tag
                            minimal
                            intent={Intent.PRIMARY}
                            style={{ fontSize: '9px' }}
                          >
                            {item.type}
                          </Tag>{' '}
                          {modifiers.active && (
                            <Icon
                              icon='small-tick'
                              color={Colors.GREEN5}
                              size={14}
                            />
                          )}
                        </>
                      }
                      style={{
                        fontSize: '11px',
                        fontWeight: modifiers.active ? 800 : 'normal',
                        backgroundColor: modifiers.active
                          ? Colors.LIGHT_GRAY4
                          : 'none'
                      }}
                    />
                  )}
                  noResults={
                    <MenuItem disabled={true} text='No epic results.' />
                  }
                  onItemSelect={(item) => {
                    setJiraIssueEpicKeyField(item)
                    onSettingsChange({ epicKeyField: item?.value })
                  }}
                  popoverProps={{
                    position: Position.TOP
                  }}
                >
                  <Button
                    intent={Intent.NONE}
                    disabled={isSaving || fieldsList.length === 0}
                    fill={true}
                    style={{
                      justifyContent: 'space-between',
                      display: 'flex',
                      minWidth: '260px',
                      maxWidth: '100%'
                    }}
                    text={
                      jiraIssueEpicKeyField
                        ? `${jiraIssueEpicKeyField.title}`
                        : 'Select...'
                    }
                    rightIcon='double-caret-vertical'
                    outlined
                  />
                </Select>
                <Button
                  disabled={!jiraIssueEpicKeyField || isSaving}
                  icon='eraser'
                  intent={jiraIssueEpicKeyField ? Intent.NONE : Intent.NONE}
                  minimal={false}
                  onClick={() => {
                    setJiraIssueEpicKeyField('')
                    onSettingsChange({ epicKeyField: '' })
                  }}
                />
              </ButtonGroup>
              <div style={{ marginLeft: '10px' }}>
                {jiraProxyError && (
                  <p style={{ color: Colors.GRAY4 }}>
                    <Icon
                      icon='warning-sign'
                      color={Colors.RED5}
                      size={12}
                      style={{ marginBottom: '2px' }}
                    />{' '}
                    {jiraProxyError.toString() || 'JIRA API not accessible.'}
                  </p>
                )}
              </div>
            </div>
          </div>

          <div
            className='story-point-select'
            style={{ display: 'flex', marginBottom: '10px' }}
          >
            <div
              className='story-point-label'
              style={{
                minWidth: '120px',
                paddingRight: '10px',
                paddingTop: '3px'
              }}
            >
              <label>Story Point</label>
            </div>
            <div style={{ display: 'flex', minWidth: '260px' }}>
              <ButtonGroup disabled={isSaving}>
                <Select
                  disabled={isSaving || fieldsList.length === 0}
                  className='select-story-key'
                  inline={true}
                  fill={true}
                  items={fieldsList}
                  activeItem={jiraIssueStoryPointField}
                  itemPredicate={(query, item) =>
                    item?.title.toLowerCase().indexOf(query.toLowerCase()) >= 0
                  }
                  itemRenderer={(item, { handleClick, modifiers }) => (
                    <MenuItem
                      disabled={jiraIssueEpicKeyField?.value === item.value}
                      active={false}
                      intent={modifiers.active ? Intent.NONE : Intent.NONE}
                      key={item.value}
                      label={item.value}
                      onClick={handleClick}
                      text={
                        <>
                          {item.title}{' '}
                          <Tag
                            minimal
                            intent={Intent.PRIMARY}
                            style={{ fontSize: '9px' }}
                          >
                            {item.type}
                          </Tag>{' '}
                          {modifiers.active && (
                            <Icon
                              icon='small-tick'
                              color={Colors.GREEN5}
                              size={14}
                            />
                          )}
                        </>
                      }
                      style={{
                        fontSize: '11px',
                        fontWeight: modifiers.active ? 800 : 'normal',
                        backgroundColor: modifiers.active
                          ? Colors.LIGHT_GRAY4
                          : 'none'
                      }}
                    />
                  )}
                  noResults={
                    <MenuItem disabled={true} text='No epic results.' />
                  }
                  onItemSelect={(item) => {
                    setJiraIssueStoryPointField(item)
                    onSettingsChange({ storyPointField: item?.value })
                  }}
                  popoverProps={{
                    position: Position.TOP
                  }}
                >
                  <Button
                    // loading={isFetchingJIRA}
                    disabled={isSaving || fieldsList.length === 0}
                    fill={true}
                    style={{
                      justifyContent: 'space-between',
                      display: 'flex',
                      minWidth: '260px',
                      maxWidth: '300px'
                    }}
                    text={
                      jiraIssueStoryPointField
                        ? `${jiraIssueStoryPointField.title}`
                        : 'Select...'
                    }
                    rightIcon='double-caret-vertical'
                    outlined
                  />
                </Select>
                <Button
                  loading={isFetchingJIRA}
                  disabled={!jiraIssueStoryPointField || isSaving}
                  icon='eraser'
                  intent={jiraIssueStoryPointField ? Intent.NONE : Intent.NONE}
                  minimal={false}
                  onClick={() => {
                    setJiraIssueStoryPointField('')
                    onSettingsChange({ storyPointField: '' })
                  }}
                />
              </ButtonGroup>
              <div style={{ marginLeft: '10px' }}>
                {jiraProxyError && (
                  <p style={{ color: Colors.GRAY4 }}>
                    <Icon
                      icon='warning-sign'
                      color={Colors.RED5}
                      size={12}
                      style={{ marginBottom: '2px' }}
                    />{' '}
                    {jiraProxyError.toString() || 'JIRA API not accessible.'}
                  </p>
                )}
              </div>
            </div>
          </div>
          <div className='headlineContainer'>
            <h5>Additional Settings</h5>
          </div>
          <div className='formContainer' style={{ maxWidth: '550px' }}>
            <FormGroup
              disabled={isSaving}
              labelFor='jira-remotelink-sha'
              className='formGroup'
              contentClassName='formGroupContent'
            >
              <label>Remotelink Commit SHA</label>
              <p>
                Issue Weblink <strong>Commit SHA Pattern</strong> (Add weblink
                to jira for gitlab.)
              </p>
              <InputGroup
                id='jira-remotelink-sha'
                fill={true}
                placeholder='/commit/([0-9a-f]{40})$'
                value={remoteLinkCommitSha}
                onChange={(e) =>
                  onSettingsChange({
                    remotelinkCommitShaPattern: e.target.value
                  })
                }
                disabled={isSaving}
                className='input'
              />
            </FormGroup>
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
