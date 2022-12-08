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
import { Button, Icon, Intent, Classes } from '@blueprintjs/core'

import * as S from './styled'

interface Props {
  isOpen: boolean
  children: React.ReactNode
  style?: React.CSSProperties
  title?: string
  cancelText?: string
  okText?: string
  okDisabled?: boolean
  okLoading?: boolean
  onCancel?: () => void
  onOk?: () => void
}

export const Dialog = ({
  isOpen,
  children,
  style,
  title,
  cancelText = 'Cancel',
  okText = 'OK',
  okDisabled,
  okLoading,
  onCancel,
  onOk
}: Props) => {
  return (
    <S.Container isOpen={isOpen} style={style}>
      {title && (
        <S.Header className={Classes.DIALOG_HEADER}>
          <h2>{title}</h2>
          <Icon icon='cross' onClick={onCancel} />
        </S.Header>
      )}
      <S.Body className={Classes.DIALOG_BODY}>{children}</S.Body>
      <S.Footer className={Classes.DIALOG_FOOTER}>
        <div className={Classes.DIALOG_FOOTER_ACTIONS}>
          <Button
            outlined
            intent={Intent.PRIMARY}
            onClick={onCancel}
            text={cancelText}
          />
          <Button
            disabled={okDisabled}
            loading={okLoading}
            intent={Intent.PRIMARY}
            text={okText}
            onClick={onOk}
          />
        </div>
      </S.Footer>
    </S.Container>
  )
}
