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

import React, { useState } from 'react'
import { ButtonGroup, Button, Icon, Intent } from '@blueprintjs/core'

import { Table, ColumnType } from '@/components'

import { CreateDialog, ViewOrEditDialog, DeleteDialog } from '../components'

import type { ConnectionItem } from './use-connection'
import { useConnection } from './use-connection'
import * as S from './styled'

type Type = 'add' | 'edit' | 'show' | 'delete'

export const WebHookConnection = () => {
  const [type, setType] = useState<Type>()
  const [record, setRecord] = useState<ConnectionItem>()

  const { loading, saving, connections, onCreate, onUpdate, onDelete } =
    useConnection()

  const handleHideDialog = () => {
    setType(undefined)
    setRecord(undefined)
  }

  const handleShowDialog = (t: Type, r?: ConnectionItem) => {
    setType(t)
    setRecord(r)
  }

  const columns: ColumnType<ConnectionItem> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id'
    },
    {
      title: 'Incoming Webhook Name',
      dataIndex: 'name',
      key: 'name',
      render: (name, row) => (
        <span onClick={() => handleShowDialog('show', row)}>{name}</span>
      )
    },
    {
      title: '',
      dataIndex: '',
      key: 'action',
      align: 'center',
      render: (_, row) => (
        <S.Action>
          <Icon icon='edit' onClick={() => handleShowDialog('edit', row)} />
          <Icon icon='trash' onClick={() => handleShowDialog('delete', row)} />
        </S.Action>
      )
    }
  ]

  return (
    <S.Wrapper>
      <ButtonGroup>
        <Button
          icon='plus'
          text='Add a Webhook'
          intent={Intent.PRIMARY}
          onClick={() => handleShowDialog('add')}
        />
      </ButtonGroup>
      <S.Inner>
        <Table loading={loading} columns={columns} dataSource={connections} />
      </S.Inner>
      {type === 'add' && (
        <CreateDialog
          isOpen
          saving={saving}
          onSubmit={onCreate}
          onCancel={handleHideDialog}
        />
      )}
      {(type === 'edit' || type === 'show') && (
        <ViewOrEditDialog
          type={type}
          initialValues={record}
          isOpen
          saving={saving}
          onSubmit={onUpdate}
          onCancel={handleHideDialog}
        />
      )}
      {type === 'delete' && (
        <DeleteDialog
          initialValues={record}
          isOpen
          saving={saving}
          onSubmit={onDelete}
          onCancel={handleHideDialog}
        />
      )}
    </S.Wrapper>
  )
}
