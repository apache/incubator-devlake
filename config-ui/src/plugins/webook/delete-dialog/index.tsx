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

import type { UseDeleteProps } from './use-delete'
import { useDelete } from './use-delete'
import * as S from './styled'

interface Props extends UseDeleteProps {
  isOpen: boolean
  onCancel: () => void
}

export const WebhookDeleteDialog = ({ isOpen, onCancel, ...props }: Props) => {
  const { saving, onSubmit } = useDelete({ ...props })

  const handleSubmit = () => {
    onSubmit()
    onCancel()
  }

  return (
    <Dialog
      isOpen={isOpen}
      title='Delete this Incoming Webhook?'
      style={{ width: 600 }}
      okText='Confirm'
      okLoading={saving}
      onCancel={onCancel}
      onOk={handleSubmit}
    >
      <S.Wrapper>
        <div className='message'>
          <p>This Incoming Webhook cannot be recovered once it’s deleted.</p>
        </div>
      </S.Wrapper>
    </Dialog>
  )
}
