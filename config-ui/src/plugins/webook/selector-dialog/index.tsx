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

import { Dialog } from '@/components'

import type { WebhookItemType } from '../types'

import { MillerColumns } from '../components'

import * as S from './styled'

interface Props {
  isOpen: boolean
  saving: boolean
  onCancel: () => void
  onSubmit: (items: WebhookItemType[]) => void
}

export const WebhookSelectorDialog = ({
  isOpen,
  saving,
  onCancel,
  onSubmit
}: Props) => {
  const [selectedItems, setSelectedItems] = useState<WebhookItemType[]>([])

  const handleSubmit = () => onSubmit(selectedItems)

  return (
    <Dialog
      isOpen={isOpen}
      title='Select Existing Webhooks'
      style={{
        width: 820
      }}
      okText='Confrim'
      okLoading={saving}
      okDisabled={!selectedItems.length}
      onCancel={onCancel}
      onOk={handleSubmit}
    >
      <S.Wrapper>
        <h3>Webhooks</h3>
        <p>Select an existing Webhook to import to the current project.</p>
        <MillerColumns
          selectedItems={selectedItems}
          onChangeItems={setSelectedItems}
        />
      </S.Wrapper>
    </Dialog>
  )
}
