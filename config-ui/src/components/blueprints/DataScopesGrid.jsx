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
import React from 'react'
import {
  Button,
  Intent,
  Card,
  Elevation,
  Colors,
  Icon,
  Classes,
  Tabs,
  Tab,
} from '@blueprintjs/core'
import { Providers, ProviderLabels, ProviderIcons } from '@/data/Providers'
import { NullBlueprint, BlueprintMode } from '@/data/NullBlueprint'

const DataScopesGrid = (props) => {
  const {
    connections = [],
    blueprint = NullBlueprint,
    mode = BlueprintMode.NORMAL,
    onModify = () => {},
    className='simplegrid',
    cardStyle = {
      padding: 0,
      minWidth: '878px'
    },
    gridStyle = {
      display: 'flex',
      flex: 1,
      width: '100%',
      flexDirection: 'column',
    },
    elevation = Elevation.TWO
  } = props

  return (

    <Card
      elevation={elevation}
      style={{ ...cardStyle }}
    >
      <div
        className={className}
        style={{
          ...gridStyle,
          backgroundColor: !blueprint?.enable ? '#f8f8f8' : 'inherit',
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
          <div
            className='cell-header entities'
            style={{ flex: 1 }}
          >
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
            <div className='cell connections' style={{ display: 'flex', flex: 1, alignItems: 'center' }}>
              <span style={{ marginBottom: '-5px', marginRight: '10px' }}>
                {c.icon}
              </span>
              <span>{c.name}</span>
            </div>
            <div className='cell entities' style={{ display: 'flex', flex: 1, alignItems: 'center' }}>
              <ul
                style={{
                  listStyle: 'none',
                  margin: 0,
                  padding: 0,
                }}
              >
                {c.entities.map((entityLabel, eIdx) => (
                  <li key={`list-item-key-${eIdx}`}>
                    {entityLabel}
                  </li>
                ))}
              </ul>
            </div>
            <div className='cell scope' style={{ display: 'flex', flex: 1, alignItems: 'center' }}>
              {[Providers.GITLAB, Providers.GITHUB].includes(
                c.provider?.id
              ) && (
                <ul
                  style={{
                    listStyle: 'none',
                    margin: 0,
                    padding: 0,
                  }}
                >
                  {c.projects.map((project, pIdx) => (
                    <li key={`list-item-key-${pIdx}`}>
                      {project}
                    </li>
                  ))}
                </ul>
              )}
              {[Providers.JIRA].includes(c.provider?.id) && (
                <ul
                  style={{
                    listStyle: 'none',
                    margin: 0,
                    padding: 0,
                  }}
                >
                  {c.boards.map((board, bIdx) => (
                    <li key={`list-item-key-${bIdx}`}>{board}</li>
                  ))}
                </ul>
              )}
            </div>
            <div
              className='cell transformation'
              style={{ display: 'flex', flex: 1, alignItems: 'center' }}
            >
              <ul
                style={{
                  listStyle: 'none',
                  margin: 0,
                  padding: 0,
                }}
              >
                {c.transformationStates.map((state, sIdx) => (
                  <li key={`list-item-key-${sIdx}`} style={{ minWidth: '80px' }}>
                    {state}
                  </li>
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
                // disabled={c.providerId === Providers.JENKINS}
                icon='annotation'
                intent={c.providerId === Providers.JENKINS ? Intent.NONE : Intent.PRIMARY}
                size={12}
                small
                minimal
                onClick={() => onModify(cIdx, c.connectionId, c.provider)}
              />
            </div>
          </div>
        ))}
      </div>
    </Card>

  )

}

export default DataScopesGrid
