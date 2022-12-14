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
import { Icon } from '@blueprintjs/core'

import type { ConnectionItemType } from '@/store'
import { Card, Table, Divider } from '@/components'
import { useStore } from '@/store'
import { DataScope } from '@/plugins'

import { useCreateBP } from '../bp-context'

import { useColumns } from './use-columns'

export const StepTwo = () => {
  const [connection, setConnection] = useState<
    ConnectionItemType & { scope: any }
  >()

  const { connections } = useStore()
  const { uniqueList, scopeMap, onChangeScopeMap, onChangeShowDetail } =
    useCreateBP()

  const handleGoDetail = (c: ConnectionItemType & { scope: any }) => {
    setConnection(c)
    onChangeShowDetail(true)
  }

  const handleBack = () => {
    setConnection(undefined)
    onChangeShowDetail(false)
  }

  const handleSave = (scope: any) => {
    if (!connection) return
    onChangeScopeMap({
      ...scopeMap,
      [`${connection.unique}`]: scope
    })
    handleBack()
  }

  const columns = useColumns({ onDetail: handleGoDetail })
  const dataSource = useMemo(
    () =>
      uniqueList.map((unique) => {
        const connection = connections.find(
          (cs) => cs.unique === unique
        ) as ConnectionItemType
        const scope = scopeMap[unique] ?? []
        return {
          ...connection,
          scope: scope.map((sc: any) => ({
            id: `${sc.id}`,
            entities: `${sc.entities}`
          }))
        }
      }),
    [uniqueList, connections, scopeMap]
  )

  return !connection ? (
    <Table columns={columns} dataSource={dataSource} />
  ) : (
    <Card>
      <div className='back' onClick={handleBack}>
        <Icon icon='arrow-left' size={14} />
        <span>Cancel and Go Back to the Data Scope List</span>
      </div>
      <h2>Add Data Scope</h2>
      <Divider />
      <div className='connection'>
        <img src={connection.icon} width={24} alt='' />
        <span>{connection.name}</span>
      </div>
      <DataScope
        plugin={connection.plugin}
        connectionId={connection.id}
        entities={connection.entities}
        onCancel={handleBack}
        onSave={handleSave}
      />
    </Card>
  )
}
