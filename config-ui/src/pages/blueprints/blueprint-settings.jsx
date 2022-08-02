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
import React, { useEffect, useState, useCallback, useRef } from 'react'
import { useParams, useHistory } from 'react-router-dom'
import dayjs from '@/utils/time'
import request from '@/utils/request'
import {
  Button,
  Elevation,
  Intent,
  Switch,
  Card,
  Tooltip,
  Icon,
  Colors,
  Divider,
  Spinner,
  Classes,
  Position,
  Popover,
  Collapse,
  Dialog,
} from '@blueprintjs/core'

import { integrationsData } from '@/data/integrations'
import { NullBlueprint } from '@/data/NullBlueprint'
import { NullPipelineRun } from '@/data/NullPipelineRun'
import { Providers, ProviderLabels, ProviderIcons } from '@/data/Providers'
import {
  StageStatus,
  TaskStatus,
  TaskStatusLabels,
  StatusColors,
  StatusBgColors,
} from '@/data/Task'

import Nav from '@/components/Nav'
import Sidebar from '@/components/Sidebar'
import Content from '@/components/Content'
import { ToastNotification } from '@/components/Toast'

import { DataEntities, DataEntityTypes } from '@/data/DataEntities'

import useBlueprintManager from '@/hooks/useBlueprintManager'

const BlueprintSettings = (props) => {
  // eslint-disable-next-line no-unused-vars
  const history = useHistory()
  const { bId } = useParams()

  const [blueprintId, setBlueprintId] = useState()
  const [activeBlueprint, setActiveBlueprint] = useState(NullBlueprint)
  const [currentRun, setCurrentRun] = useState(NullPipelineRun)

  const [connections, setConnections] = useState([
    {
      id: 0,
      provider: integrationsData.find((i) => i.id === Providers.GITHUB),
      providerLabel: ProviderLabels[Providers.GITHUB],
      name: 'Merico Github',
      entities: ['Source Code Management', 'Issue Tracking', 'Code Review'],
      projects: [
        'apache/incubator-devlake',
        'merico/devstream',
        'merico/another-project',
      ],
      boards: [],
      transformation: {},
      transformationStates: ['Added', '-', 'Added'],
      editable: true,
    },
    {
      id: 1,
      provider: integrationsData.find((i) => i.id === Providers.JIRA),
      providerLabel: ProviderLabels[Providers.JIRA],
      name: 'Merico JIRA',
      entities: ['Source Code Management', 'Issue Tracking', 'Code Review'],
      projects: [],
      boards: ['Board 1', 'Board 2', 'Board 3', 'Board 4'],
      transformation: {},
      transformationStates: ['Added', 'Added', '-', '-'],
      editable: true,
    },
  ])

  const [blueprintDialogIsOpen, setBlueprintDialogIsOpen] = useState(false)

  const {
    // eslint-disable-next-line no-unused-vars
    blueprint,
    isSaving,
    activateBlueprint,
    deactivateBlueprint,
    getNextRunDate,
    // eslint-disable-next-line no-unused-vars
    fetchBlueprint,
    fetchAllBlueprints,
    saveBlueprint,
    saveComplete,
  } = useBlueprintManager()

  const handleBlueprintActivation = useCallback(
    (blueprint) => {
      if (blueprint.enable) {
        deactivateBlueprint(blueprint)
      } else {
        activateBlueprint(blueprint)
      }
    },
    [activateBlueprint, deactivateBlueprint]
  )

  const viewBlueprintStatus = useCallback(() => {
    history.push(`/blueprints/detail/${blueprintId}`)
  }, [history, blueprintId])

  const viewBlueprintSettings = useCallback(() => {
    history.push(`/blueprints/settings/${blueprintId}`)
  }, [history])

  useEffect(() => {
    setBlueprintId(bId)
    console.log('>>> REQUESTED SETTINGS for BLUEPRINT ID ===', bId)
  }, [bId])

  useEffect(() => {
    if (!isNaN(blueprintId)) {
      console.log('>>>> FETCHING BLUEPRINT ID...', blueprintId)
      fetchBlueprint(blueprintId)
    }
  }, [blueprintId, fetchBlueprint])

  useEffect(() => {
    console.log('>>>> SETTING ACTIVE BLUEPRINT...', blueprint)
    if (blueprint?.id) {
      setActiveBlueprint((b) => ({
        ...b,
        ...blueprint,
      }))
    }
  }, [blueprint])

  useEffect(() => {
    setConnections(activeBlueprint?.settings?.connections.map((c, cIdx) => ({
      id: cIdx,
      provider: integrationsData.find((i) => i.id === c.plugin),
      providerLabel: ProviderLabels[c.plugin],
      name: `Connection ID #${c.connectionId}`,
      entities: ['Source Code Management', 'Issue Tracking', 'Code Review'],
      projects: [Providers.GITLAB, Providers.GITHUB].includes(c.plugin) ? c.scope.map(s => `${s.options.owner}/${s.options?.repo}`) : [],
      boards: [Providers.JIRA].includes(c.plugin) ? c.scope.map(s => `Board ${s.options?.boardId}`) : [],
      transformation: { ...c.transformation },
      transformationStates: c.scope.map(s => Object.values(s.transformation).some(v => v?.toString().length > 0) ? 'Added' : '-'),
      editable: true,
    })))
  }, [activeBlueprint, setConnections])

  return (
    <>
      <div className='container'>
        <Nav />
        <Sidebar />
        <Content>
          <main className='main'>
            <div
              className='blueprint-header'
              style={{
                display: 'flex',
                width: '100%',
                justifyContent: 'space-between',
                marginBottom: '10px',
                whiteSpace: 'nowrap'
              }}
            >
              <div className='blueprint-name' style={{}}>
                <h2 style={{ fontWeight: 'bold' }}>{activeBlueprint?.name}</h2>
              </div>
              <div
                className='blueprint-info'
                style={{ display: 'flex', alignItems: 'center' }}
              >
                <div className='blueprint-schedule'>
                  <span
                    className='blueprint-schedule-interval'
                    style={{ textTransform: 'capitalize', padding: '0 10px' }}
                  >
                    {activeBlueprint?.interval} (at{' '}
                    {dayjs(getNextRunDate(activeBlueprint?.cronConfig)).format(
                      'hh:mm A'
                    )}
                    )
                  </span>{' '}
                  &nbsp;{' '}
                </div>
                <div
                  className='blueprint-actions'
                  style={{ padding: '0 10px' }}
                >
                  {/* <Button
                    intent={Intent.PRIMARY}
                    small
                    text='Run Now'
                    onClick={runBlueprint}
                    disabled={!activeBlueprint?.enable || currentRun?.status === TaskStatus.RUNNING}
                  /> */}
                </div>
                <div className='blueprint-enabled'>
                  <Switch
                    id='blueprint-enable'
                    name='blueprint-enable'
                    checked={activeBlueprint?.enable}
                    label={
                      activeBlueprint?.enable
                        ? 'Blueprint Enabled'
                        : 'Blueprint Disabled'
                    }
                    onChange={() => handleBlueprintActivation(activeBlueprint)}
                    style={{
                      marginBottom: 0,
                      marginTop: 0,
                      color: !activeBlueprint?.enable
                        ? Colors.GRAY3
                        : 'inherit',
                    }}
                    disabled={currentRun?.status === TaskStatus.RUNNING}
                  />
                </div>
                <div style={{ padding: '0 10px' }}>
                  <Button
                    intent={Intent.PRIMARY}
                    icon='trash'
                    small
                    minimal
                    disabled
                  />
                </div>
              </div>
            </div>

            <div
              className='blueprint-navigation'
              style={{
                alignSelf: 'center',
                display: 'flex',
                margin: '20px auto',
              }}
            >
              <div style={{ marginRight: '10px' }}>
                <a
                  href='#'
                  className='blueprint-navigation-link'
                  onClick={viewBlueprintStatus}
                >
                  Status
                </a>
              </div>
              <div style={{ marginLeft: '10px' }}>
                <a
                  href='#'
                  className='blueprint-navigation-link active'
                  onClick={viewBlueprintSettings}
                >
                  Settings
                </a>
              </div>
            </div>

            <div
              className='blueprint-main-settings'
              style={{ display: 'flex', alignSelf: 'flex-start' }}
            >
              <div className='configure-settings-name'>
                <h3>Name</h3>
                <div style={{ display: 'flex', alignItems: 'center' }}>
                  <div className='blueprint-name'>{activeBlueprint?.name}</div>
                  <Button
                    icon='annotation'
                    intent={Intent.PRIMARY}
                    size={12}
                    small
                    minimal
                  />
                </div>
              </div>
              <div
                className='configure-settings-frequency'
                style={{ marginLeft: '40px' }}
              >
                <h3>Sync Frequency</h3>
                <div style={{ display: 'flex', alignItems: 'center' }}>
                  <div className='blueprint-frequency'>
                    {activeBlueprint?.interval} (at{' '}
                    {dayjs(getNextRunDate(activeBlueprint?.cronConfig)).format(
                      'hh:mm A'
                    )}
                    )
                  </div>
                  <Button
                    icon='annotation'
                    intent={Intent.PRIMARY}
                    size={12}
                    small
                    minimal
                  />
                </div>
              </div>
            </div>

            <div
              className='data-scopes-grid'
              style={{
                width: '100%',
                marginTop: '40px',
                alignSelf: 'flex-start',
              }}
            >
              <h2 style={{ fontWeight: 'bold' }}>
                Data Scope and Transformation
              </h2>

              <Card elevation={Elevation.TWO} style={{ padding: 0, minWidth: '878px' }}>
                <div
                  className='simplegrid'
                  style={{
                    display: 'flex',
                    flex: 1,
                    width: '100%',
                    flexDirection: 'column',
                  }}
                >
                  <div
                    className='simplegrid-header'
                    style={{
                      display: 'flex',
                      flex: 1,
                      width: '100%',
                      minHeight: '48px',
                      lineHeight: 'auto',
                      padding: '16px 20px',
                      fontWeight: 'bold',
                      borderBottom: '1px solid #BDCEFB',
                      justfiyContent: 'space-evenly',
                    }}
                  >
                    <div
                      className='cell-header connections'
                      style={{ flex: 1 }}
                    >
                      Data Connections
                    </div>
                    <div className='cell-header entities' style={{ flex: 1 }}>
                      Data Entities
                    </div>
                    <div className='cell-header scope' style={{ flex: 1 }}>
                      Data Scope
                    </div>
                    <div
                      className='cell-header transformation'
                      style={{ flex: 1 }}
                    >
                      Transformation
                    </div>
                    <div
                      className='cell-header actions'
                      style={{ minWidth: '100px' }}
                    >
                      &nbsp;
                    </div>
                  </div>

                  {connections.map((c, cIdx) => (
                    <div
                      key={`connection-row-key-${cIdx}`}
                      className='simplegrid-row'
                      style={{
                        display: 'flex',
                        flex: 1,
                        width: '100%',
                        minHeight: '48px',
                        lineHeight: 'auto',
                        padding: '10px 20px',
                        borderBottom: '1px solid #BDCEFB',
                        justfiyContent: 'space-evenly',
                      }}
                    >
                      <div className='cell connections' style={{ flex: 1 }}>
                        {c.name}
                      </div>
                      <div className='cell entities' style={{ flex: 1 }}>
                        <ul
                          style={{ listStyle: 'none', margin: 0, padding: 0 }}
                        >
                          {c.entities.map((entityLabel, eIdx) => (
                            <li key={`list-item-key-${eIdx}`}>{entityLabel}</li>
                          ))}
                        </ul>
                      </div>
                      <div className='cell scope' style={{ flex: 1 }}>
                        {[Providers.GITLAB, Providers.GITHUB].includes(c.provider?.id) && (
                          <ul
                            style={{ listStyle: 'none', margin: 0, padding: 0 }}
                          >
                            {c.projects.map((project, pIdx) => (
                              <li key={`list-item-key-${pIdx}`}>{project}</li>
                            ))}
                          </ul>
                        )}
                        {[Providers.JIRA].includes(c.provider?.id) && (
                          <ul
                            style={{ listStyle: 'none', margin: 0, padding: 0 }}
                          >
                            {c.boards.map((board, bIdx) => (
                              <li key={`list-item-key-${bIdx}`}>{board}</li>
                            ))}
                          </ul>
                        )}
                      </div>
                      <div className='cell transformation' style={{ flex: 1 }}>
                        <ul
                          style={{ listStyle: 'none', margin: 0, padding: 0 }}
                        >
                          {c.transformationStates.map((state, sIdx) => (
                            <li key={`list-item-key-${sIdx}`}>{state}</li>
                          ))}
                        </ul>
                      </div>
                      <div
                        className='cell actions'
                        style={{
                          display: 'flex',
                          minWidth: '100px',
                          textAlign: 'right',
                          alignItems: 'center',
                          justifyContent: 'flex-end',
                        }}
                      >
                        <Button
                          icon='annotation'
                          intent={Intent.PRIMARY}
                          size={12}
                          small
                          minimal
                        />
                      </div>
                    </div>
                  ))}
                </div>
              </Card>
            </div>
          </main>
        </Content>
      </div>
    </>
  )
}

export default BlueprintSettings
