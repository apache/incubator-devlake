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
import { Transformation } from '@/plugins'

import { useCreateBP } from '../bp-context'

import { useColumns } from './use-columns'

export const StepThree = () => {
  const [connection, setConnection] = useState<
    ConnectionItemType & { scope: any }
  >()

  const { connections } = useStore()
  const { uniqueList, scopeMap, onChangeShowDetail } = useCreateBP()

  const handleGoDetail = (c: ConnectionItemType & { scope: any }) => {
    setConnection(c)
    onChangeShowDetail(true)
  }

  const handleBack = () => {
    setConnection(undefined)
    onChangeShowDetail(false)
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
        <span>Cancel and Go Back</span>
      </div>
      <h2>Create/Select a Transformation</h2>
      <Divider />
      <Transformation
        plugin={connection.plugin}
        connectionId={connection.id}
        scopeIds={connection.scope.map((sc: any) => sc.id)}
        onCancel={handleBack}
        onSave={handleBack}
      />
    </Card>
  )
}
