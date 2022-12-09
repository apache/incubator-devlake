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

import type { BPConnectionItemType, BPScopeItemType } from '../types'
import { ModeEnum } from '../types'

export type BlueprintContextType = {
  step: number
  error: string
  showInspector: boolean
  showDetail: boolean
  payload: any

  name: string
  mode: ModeEnum
  rawPlan: string
  connections: BPConnectionItemType[]
  scope: BPScopeItemType[]
  cronConfig: string
  isManual: boolean
  skipOnFail: boolean

  onChangeStep: React.Dispatch<React.SetStateAction<number>>
  onChangeShowInspector: React.Dispatch<React.SetStateAction<boolean>>
  onChangeShowDetail: React.Dispatch<React.SetStateAction<boolean>>

  onChangeName: React.Dispatch<React.SetStateAction<string>>
  onChangeMode: (mode: ModeEnum) => void
  onChangeRawPlan: React.Dispatch<React.SetStateAction<string>>
  onChangeConnections: React.Dispatch<
    React.SetStateAction<BPConnectionItemType[]>
  >
  onChangeScope: React.Dispatch<React.SetStateAction<BPScopeItemType[]>>
  onChangeCronConfig: React.Dispatch<React.SetStateAction<string>>
  onChangeIsManual: React.Dispatch<React.SetStateAction<boolean>>
  onChangeSkipOnFail: React.Dispatch<React.SetStateAction<boolean>>

  onSave: () => void
  onSaveAndRun: () => void
}

export const BlueprintContext = React.createContext<BlueprintContextType>({
  step: 1,
  error: '',
  showInspector: false,
  showDetail: false,
  payload: {},

  name: 'MY BLUEPRINT',
  mode: ModeEnum.normal,
  rawPlan: JSON.stringify([[]], null, '  '),
  connections: [],
  scope: [],
  cronConfig: '0 0 * * *',
  isManual: false,
  skipOnFail: false,

  onChangeStep: () => {},
  onChangeShowInspector: () => {},
  onChangeShowDetail: () => {},

  onChangeMode: () => {},
  onChangeName: () => {},
  onChangeRawPlan: () => {},
  onChangeConnections: () => {},
  onChangeScope: () => {},
  onChangeCronConfig: () => {},
  onChangeIsManual: () => {},
  onChangeSkipOnFail: () => {},

  onSave: () => {},
  onSaveAndRun: () => {}
})
