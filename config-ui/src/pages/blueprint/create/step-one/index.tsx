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

import React, { useState, useMemo, useEffect } from 'react'
import {
  InputGroup,
  TextArea,
  ButtonGroup,
  Button,
  Menu,
  MenuItem,
  Position
} from '@blueprintjs/core'
import { Popover2 } from '@blueprintjs/popover2'

import { useConnections, ConnectionStatusEnum } from '@/store'
import { Divider, MultiSelector, Loading } from '@/components'

import { ModeEnum } from '../types'
import { useBlueprint } from '../hooks'

import { DEFAULT_CONFIG } from './example'
import * as S from './styled'

export const StepOne = () => {
  const [isOpen, setIsOpen] = useState(false)

  const connectionsStore = useConnections()

  const {
    mode,
    name,
    rawPlan,
    connections,
    onChangeMode,
    onChangeName,
    onChangeRawPlan,
    onChangeConnections
  } = useBlueprint()

  const selectedConnections = useMemo(
    () =>
      connectionsStore.connections.filter((cs) =>
        connections.map((cs) => cs.unique).includes(cs.unique)
      ),
    [connectionsStore, connections]
  )

  return (
    <>
      <S.Card>
        <h2>Blueprint Name</h2>
        <Divider />
        <p>
          Give your Blueprint a unique name to help you identify it in the
          future.
        </p>
        <InputGroup
          placeholder='Enter Blueprint Name'
          value={name}
          onChange={(e) => onChangeName(e.target.value)}
        />
      </S.Card>

      {mode === ModeEnum.normal && (
        <>
          <S.Card>
            <h2>Add Data Connections</h2>
            <Divider />
            <h3>Select Connections</h3>
            <p>Select from existing or create new connections</p>
            <MultiSelector
              placeholder='Select Connections...'
              items={connectionsStore.connections}
              getKey={(it) => it.unique}
              getName={(it) => it.name}
              getIcon={(it) => it.icon}
              selectedItems={selectedConnections}
              onChangeItems={(selectedItems) => {
                connectionsStore.onTest(selectedItems)
                onChangeConnections(
                  selectedItems.map((it) => ({
                    ...it,
                    scope: []
                  }))
                )
              }}
            />
            <S.ConnectionList>
              {selectedConnections.map((cs) => (
                <li key={cs.unique}>
                  <span className='name'>{cs.name}</span>
                  <span className={`status ${cs.status}`}>
                    {cs.status === 'testing' && (
                      <Loading size={14} style={{ marginRight: 4 }} />
                    )}
                    {cs.status}
                  </span>
                </li>
              ))}
            </S.ConnectionList>
          </S.Card>
          <S.Tips>
            <span>
              To customize how tasks are executed in the blueprint, please use{' '}
            </span>
            <span onClick={() => onChangeMode(ModeEnum.advanced)}>
              Advanced Mode.
            </span>
          </S.Tips>
        </>
      )}

      {mode === ModeEnum.advanced && (
        <>
          <S.Card>
            <h2>JSON Configuration</h2>
            <Divider />
            <h3>Task Editor</h3>
            <p>Enter JSON Configuration or preload from a template</p>
            <TextArea
              fill
              value={rawPlan}
              onChange={(e) => onChangeRawPlan(e.target.value)}
            />
            <ButtonGroup minimal>
              <Button small text='Reset' icon='eraser' />
              <Popover2
                placement={Position.TOP}
                isOpen={isOpen}
                content={
                  <Menu>
                    {DEFAULT_CONFIG.map((it) => (
                      <MenuItem
                        key={it.id}
                        icon='code'
                        text={it.name}
                        onClick={() => {
                          setIsOpen(false)
                          onChangeRawPlan(JSON.stringify(it.config, null, '  '))
                        }}
                      />
                    ))}
                  </Menu>
                }
              >
                <Button
                  small
                  text='Load Templates'
                  rightIcon='caret-down'
                  onClick={() => setIsOpen(!isOpen)}
                />
              </Popover2>
            </ButtonGroup>
          </S.Card>
          <S.Tips>
            <span>To visually define blueprint tasks, please use </span>
            <span onClick={() => onChangeMode(ModeEnum.normal)}>
              Normal Mode.
            </span>
          </S.Tips>
        </>
      )}
    </>
  )
}
