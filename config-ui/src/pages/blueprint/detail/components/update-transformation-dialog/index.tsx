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
import { Transformation } from '@/plugins'

import type { ConnectionItemType } from '../../types'

interface Props {
  connection?: ConnectionItemType
  onCancel: () => void
}

export const UpdateTransformationDialog = ({ connection, onCancel }: Props) => {
  if (!connection) return null

  const { plugin, connectionId, scopeIds } = connection

  return (
    <Dialog
      isOpen
      title='Change Data Scope'
      footer={null}
      style={{ width: 900 }}
      onCancel={onCancel}
    >
      <Transformation
        plugin={plugin}
        connectionId={connectionId}
        scopeIds={scopeIds}
        onCancel={onCancel}
        onSave={onCancel}
      />
    </Dialog>
  )
}
