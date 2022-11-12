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
import classNames from 'classnames'

import { CheckStatus } from './types'
import * as S from './styled'

interface Props {
  status?: CheckStatus | Array<CheckStatus>
  children?: React.ReactNode
  onClick?: () => void
}

export const Checkbox = ({ children, status, onClick }: Props) => {
  const checkboxCls = classNames('checkbox', {
    'checkbox-checked':
      status === CheckStatus.checked ||
      (Array.isArray(status) && status.includes(CheckStatus.checked)),
    'checkbox-indeterminate':
      status === CheckStatus.indeterminate ||
      (Array.isArray(status) && status.includes(CheckStatus.indeterminate)),
    'checkbox-disabled':
      status === CheckStatus.disabled ||
      (Array.isArray(status) && status?.includes(CheckStatus.disabled))
  })

  const handleClick = (e: React.MouseEvent<HTMLDivElement>) => {
    e.stopPropagation()
    if (status === CheckStatus.disabled) return
    onClick?.()
  }

  return (
    <S.Wrapper>
      <span className={checkboxCls} onClick={handleClick}>
        <span className='checkbox-inner'></span>
      </span>
      {children && <span className='text'>{children}</span>}
    </S.Wrapper>
  )
}
