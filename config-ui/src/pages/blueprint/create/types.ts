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

export enum FromEnum {
  project = 'project',
  blueprint = 'blueprint'
}

export enum ModeEnum {
  advanced = 'ADVANCED',
  normal = 'NORMAL'
}

export type BPContextType = {
  step: number
  error: string
  showInspector: boolean
  showDetail: boolean
  payload: any

  name: string
  mode: ModeEnum
  rawPlan: string
  uniqueList: string[]
  scopeMap: Record<string, any>
  cronConfig: string
  isManual: boolean
  skipOnFail: boolean
  createdDateAfter: string | null

  onChangeStep: React.Dispatch<React.SetStateAction<number>>
  onChangeShowInspector: React.Dispatch<React.SetStateAction<boolean>>
  onChangeShowDetail: React.Dispatch<React.SetStateAction<boolean>>

  onChangeName: React.Dispatch<React.SetStateAction<string>>
  onChangeMode: (mode: ModeEnum) => void
  onChangeRawPlan: React.Dispatch<React.SetStateAction<string>>
  onChangeUniqueList: React.Dispatch<React.SetStateAction<string[]>>
  onChangeScopeMap: React.Dispatch<React.SetStateAction<Record<string, any>>>
  onChangeCronConfig: React.Dispatch<React.SetStateAction<string>>
  onChangeIsManual: React.Dispatch<React.SetStateAction<boolean>>
  onChangeSkipOnFail: React.Dispatch<React.SetStateAction<boolean>>
  onChangeCreatedDateAfter: React.Dispatch<React.SetStateAction<string | null>>

  onSave: () => void
  onSaveAndRun: () => void
}
