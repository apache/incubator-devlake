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

import type { WebhookItemType } from '../types'
import { WebhookCreateDialog } from '../create-dialog'
import { WebhookDeleteDialog } from '../delete-dialog'
import { WebhookViewOrEditDialog } from '../view-or-edit-dialog'

import type { UseConnectionProps } from './use-connection'
import { useConnection } from './use-connection'
import * as S from './styled'

type Type = 'add' | 'edit' | 'show' | 'delete'

interface Props extends UseConnectionProps {
  onCreateAfter?: (id: ID) => void
}

export const WebHookConnection = ({ onCreateAfter, ...props }: Props) => {
  const [type, setType] = useState<Type>()
  const [record, setRecord] = useState<WebhookItemType>()

  const { loading, connections, onRefresh } = useConnection({ ...props })

  const handleHideDialog = () => {
    setType(undefined)
    setRecord(undefined)
  }

  const handleShowDialog = (t: Type, r?: WebhookItemType) => {
    setType(t)
    setRecord(r)
  }

  const columns: ColumnType<WebhookItemType> = [
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
      <Table loading={loading} columns={columns} dataSource={connections} />
      {type === 'add' && (
        <WebhookCreateDialog
          isOpen
          onCancel={handleHideDialog}
          onSubmitAfter={(id) => {
            onRefresh()
            onCreateAfter?.(id)
          }}
        />
      )}
      {type === 'delete' && (
        <WebhookDeleteDialog
          isOpen
          initialValues={record}
          onCancel={handleHideDialog}
          onSubmitAfter={onRefresh}
        />
      )}
      {(type === 'edit' || type === 'show') && (
        <WebhookViewOrEditDialog
          type={type}
          isOpen
          initialValues={record}
          onCancel={handleHideDialog}
          onSubmitAfter={onRefresh}
        />
      )}
    </S.Wrapper>
  )
}
