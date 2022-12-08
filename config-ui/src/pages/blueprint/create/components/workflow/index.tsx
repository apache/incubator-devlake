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

import React, { useMemo } from 'react'

import { ModeEnum } from '../../types'
import { useBlueprint } from '../../hooks'

import * as S from './styled'

export const WorkFlow = () => {
  const { step, mode } = useBlueprint()

  const steps = useMemo(
    () =>
      mode === ModeEnum.normal
        ? [
            'Add Data Connections',
            'Set Data Scope',
            'Add Transformation (Optional)',
            'Set Sync Frequency'
          ]
        : ['Create Advanced Configuration', 'Set Sync Frequency'],
    [mode]
  )

  return (
    <S.List>
      {steps.map((it, i) => (
        <S.Item key={it} active={i + 1 === step}>
          <span className='step'>{i + 1}</span>
          <span className='name'>{it}</span>
        </S.Item>
      ))}
    </S.List>
  )
}
