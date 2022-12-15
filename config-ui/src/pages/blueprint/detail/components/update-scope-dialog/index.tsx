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

import { Dialog } from '@/components'
import { DataScope } from '@/plugins'

import type { ConnectionItemType } from '../../types'

interface Props {
  connection?: ConnectionItemType
  onCancel: () => void
  onSubmit: (connection: any) => void
}

export const UpdateScopeDialog = ({
  connection,
  onCancel,
  onSubmit
}: Props) => {
  if (!connection) return null

  const { plugin, connectionId, entities } = connection

  const handleSaveScope = (sc: any) => {
    onSubmit({
      plugin,
      connectionId,
      scopes: sc.map((it: any) => ({
        ...it,
        id: `${it.id}`
      }))
    })
  }

  return (
    <Dialog
      isOpen
      title='Change Data Scope'
      footer={null}
      style={{ width: 900 }}
      onCancel={onCancel}
    >
      <DataScope
        plugin={plugin}
        connectionId={connectionId}
        entities={entities}
        onCancel={onCancel}
        onSave={handleSaveScope}
      />
    </Dialog>
  )
}
