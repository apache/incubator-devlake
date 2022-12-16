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

import React, { useState, useMemo } from 'react'
import { Icon, Colors } from '@blueprintjs/core'
import dayjs from 'dayjs'

import { Table, ColumnType } from '@/components'
import { getCron } from '@/config'
import { DataScopeList } from '@/plugins'

import type { BlueprintType, ConnectionItemType } from '../types'
import {
  UpdateNameDialog,
  UpdatePolicyDialog,
  UpdateScopeDialog,
  UpdateTransformationDialog
} from '../components'
import * as S from '../styled'

type Type = 'name' | 'frequency' | 'scope' | 'transformation'

interface Props {
  blueprint: BlueprintType
  connections: ConnectionItemType[]
  saving: boolean
  onUpdate: (bp: any) => Promise<void>
}

export const Configuration = ({
  blueprint,
  connections,
  saving,
  onUpdate
}: Props) => {
  const [type, setType] = useState<Type>()
  const [curConnection, setCurConnection] = useState<ConnectionItemType>()

  const cron = useMemo(
    () => getCron(blueprint.isManual, blueprint.cronConfig),
    [blueprint]
  )

  const handleCancel = () => {
    setType(undefined)
  }

  const handleUpdateName = async (name: string) => {
    await onUpdate({ name })
    handleCancel()
  }

  const handleUpdatePolicy = async (policy: any) => {
    await onUpdate(policy)
    handleCancel
  }

  const handleUpdateConnection = async (updated: any) => {
    await onUpdate({
      settings: {
        version: '2.0.0',
        connections: blueprint.settings.connections.map((cs) =>
          cs.plugin === updated.plugin &&
          cs.connectionId === updated.connectionId
            ? updated
            : cs
        )
      }
    })
  }

  const columns = useMemo(
    () =>
      [
        {
          title: 'Data Connections',
          dataIndex: ['icon', 'name'],
          key: 'connection',
          render: ({
            icon,
            name
          }: Pick<ConnectionItemType, 'icon' | 'name'>) => (
            <S.ConnectionColumn>
              <img src={icon} alt='' />
              <span>{name}</span>
            </S.ConnectionColumn>
          )
        },
        {
          title: 'Data Entities',
          dataIndex: 'entities',
          key: 'entities',
          render: (val: string[]) => (
            <>
              {val.map((it) => (
                <div>{it}</div>
              ))}
            </>
          )
        },
        {
          title: 'Data Scope and Transformation',
          dataIndex: ['plugin', 'connectionId', 'scopeIds'],
          key: 'sopce',
          render: ({
            plugin,
            connectionId,
            scopeIds
          }: Pick<
            ConnectionItemType,
            'plugin' | 'connectionId' | 'scopeIds'
          >) => (
            <DataScopeList
              groupByTs
              plugin={plugin}
              connectionId={connectionId}
              scopeIds={scopeIds}
            />
          )
        },
        {
          title: '',
          key: 'action',
          align: 'center',
          render: (_, row: ConnectionItemType) => (
            <S.ActionColumn>
              <div
                className='item'
                onClick={() => {
                  setType('scope')
                  setCurConnection(row)
                }}
              >
                <Icon icon='annotation' color={Colors.BLUE2} />
                <span>Change Data Scope</span>
              </div>
              <div
                className='item'
                onClick={() => {
                  setType('transformation')
                  setCurConnection(row)
                }}
              >
                <Icon icon='annotation' color={Colors.BLUE2} />
                <span>Change Transformation</span>
              </div>
            </S.ActionColumn>
          )
        }
      ] as ColumnType<ConnectionItemType>,
    []
  )

  return (
    <S.ConfigurationPanel>
      <div className='top'>
        <div className='block'>
          <h3>Name</h3>
          <div className='detail'>
            <span>{blueprint.name}</span>
            <Icon
              icon='annotation'
              color={Colors.BLUE2}
              onClick={() => setType('name')}
            />
          </div>
        </div>
        <div className='block'>
          <h3>Sync Policy</h3>
          <div className='detail'>
            <span>
              {cron.label}
              {cron.value !== 'manual'
                ? dayjs(cron.nextTime).format('HH:mm A')
                : null}
            </span>
            <Icon
              icon='annotation'
              color={Colors.BLUE2}
              onClick={() => setType('frequency')}
            />
          </div>
        </div>
      </div>
      <div className='bottom'>
        <h3>Data Scope and Transformation</h3>
        <Table columns={columns} dataSource={connections} />
      </div>
      {type === 'name' && (
        <UpdateNameDialog
          name={blueprint.name}
          saving={saving}
          onCancel={handleCancel}
          onSubmit={handleUpdateName}
        />
      )}
      {type === 'frequency' && (
        <UpdatePolicyDialog
          blueprint={blueprint}
          isManual={blueprint.isManual}
          cronConfig={blueprint.cronConfig}
          skipOnFail={blueprint.skipOnFail}
          createdDateAfter={blueprint.settings.createdDateAfter}
          saving={saving}
          onCancel={handleCancel}
          onSubmit={handleUpdatePolicy}
        />
      )}
      {type === 'scope' && (
        <UpdateScopeDialog
          connection={curConnection}
          onCancel={handleCancel}
          onSubmit={handleUpdateConnection}
        />
      )}
      {type === 'transformation' && (
        <UpdateTransformationDialog
          connection={curConnection}
          onCancel={handleCancel}
        />
      )}
    </S.ConfigurationPanel>
  )
}
