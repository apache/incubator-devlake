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
  status?: CheckStatus
  children?: React.ReactNode
  onClick?: (e: React.MouseEvent<HTMLDivElement>) => void
}

export const Checkbox = ({ children, status, onClick }: Props) => {
  const checkboxCls = classNames('checkbox', {
    'checkbox-checked': status === CheckStatus.checked,
    'checkbox-indeterminate': status === CheckStatus.indeterminate
  })

  return (
    <S.Wrapper>
      <span className={checkboxCls} onClick={onClick}>
        <span className='checkbox-inner'></span>
      </span>
      {children && <span className='text'>{children}</span>}
    </S.Wrapper>
  )
}
