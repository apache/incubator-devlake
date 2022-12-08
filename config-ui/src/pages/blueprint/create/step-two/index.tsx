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

import React, { useState, useEffect } from 'react'
import { Icon } from '@blueprintjs/core'

import { Table, Divider } from '@/components'
import { DataScope } from '@/plugins'

import type { BPConnectionItemType, BPScopeItemType } from '../types'
import { useBlueprint } from '../hooks'

import { useColumns } from './use-columns'
import * as S from './styled'

export const StepTwo = () => {
  const [connection, setConnection] = useState<BPConnectionItemType>()

  const {
    connections,
    onChangeConnections,
    onChangeShowDetail,
    onChangeError
  } = useBlueprint()

  const handleGoDetail = (c: BPConnectionItemType) => {
    setConnection(c)
    onChangeShowDetail(true)
  }

  const handleBack = () => {
    setConnection(undefined)
    onChangeShowDetail(false)
  }

  const handleSave = (scope: BPScopeItemType[]) => {
    if (!connection) return
    const newConnections = connections.map((cs) =>
      cs.unique !== connection.unique ? cs : { ...cs, scope }
    )
    onChangeConnections(newConnections)
    handleBack()
  }

  const columns = useColumns({ onDetail: handleGoDetail })

  useEffect(() => {
    switch (true) {
      case !connections.every((cs) => cs.scope.length):
        return onChangeError('No Data Scope is Selected')
      default:
        return onChangeError('')
    }
  }, [connections])

  return !connection ? (
    <S.Card style={{ padding: 0 }}>
      <Table columns={columns} dataSource={connections} />
    </S.Card>
  ) : (
    <S.Card>
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
        allEntities={connection.entities}
        onCancel={handleBack}
        onSaveAfter={handleSave}
      />
    </S.Card>
  )
}
