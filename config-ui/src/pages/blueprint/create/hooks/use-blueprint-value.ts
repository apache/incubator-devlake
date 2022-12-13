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

import { useState, useMemo } from 'react'
import { useHistory } from 'react-router-dom'

import { operator } from '@/utils'

import type { BPConnectionItemType, BPScopeItemType } from '../types'
import { FromEnum, ModeEnum } from '../types'
import * as API from '../api'

interface Props {
  from: FromEnum
  projectName: string
}

export const useBlueprintValue = ({ from, projectName }: Props) => {
  const [step, setStep] = useState(1)
  const [showInspector, setShowInspector] = useState(false)
  const [showDetail, setShowDetail] = useState(false)

  const [name, setName] = useState(
    from === FromEnum.project ? `${projectName}-BLUEPRINT` : 'MY BLUEPRINT'
  )
  const [mode, setMode] = useState<ModeEnum>(ModeEnum.normal)
  const [rawPlan, setRawPlan] = useState(JSON.stringify([[]], null, '  '))
  const [connections, setConnections] = useState<BPConnectionItemType[]>([])
  const [scope, setScope] = useState<BPScopeItemType[]>([])
  const [cronConfig, setCronConfig] = useState('0 0 * * *')
  const [isManual, setIsManual] = useState(false)
  const [skipOnFail, setSkipOnFail] = useState(false)
  const [createdDateAfter, setCreatedDateAfter] = useState<string | null>(null)

  const history = useHistory()

  const validRawPlan = (rp: string) => {
    try {
      const p = JSON.parse(rp)
      if (p.flat().length === 0) {
        return true
      }
      return false
    } catch {
      return true
    }
  }

  const error = useMemo(() => {
    switch (true) {
      case !name:
        return 'Blueprint Name: Enter a valid Name'
      case name.length < 3:
        return 'Blueprint Name: Name too short, 3 chars minimum.'
      case mode === ModeEnum.advanced && validRawPlan(rawPlan):
        return 'Advanced Mode: Invalid/Empty Configuration'
      case mode === ModeEnum.normal && !connections.length:
        return 'Normal Mode: No Data Connections selected.'
      case step === 2 && !connections.every((cs) => cs.scope.length):
        return 'No Data Scope is Selected'
      default:
        return ''
    }
  }, [name, mode, rawPlan, connections, step])

  const payload = useMemo(() => {
    const params: any = {
      name,
      projectName,
      mode,
      enable: true,
      cronConfig,
      isManual,
      skipOnFail,
      createdDateAfter
    }

    if (mode === ModeEnum.normal) {
      params.settings = {
        version: '2.0.0',
        connections: connections.map((cs) => ({
          plugin: cs.plugin,
          connectionId: cs.id,
          scopes: cs.scope.map((sc) => ({
            ...sc,
            id: `${sc.id}`
          }))
        }))
      }
    }

    if (mode === ModeEnum.advanced) {
      params.plan = validRawPlan(rawPlan) ? JSON.parse(rawPlan) : [[]]
    }

    return params
  }, [
    name,
    projectName,
    mode,
    rawPlan,
    cronConfig,
    isManual,
    skipOnFail,
    createdDateAfter,
    connections,
    scope
  ])

  const handleSave = async () => {
    const [success] = await operator(() => API.createBlueprint(payload))
    const path =
      from === FromEnum.blueprint ? '/blueprints' : `/projects/${projectName}`
    if (success) {
      history.push(path)
    }
  }

  const hanldeSaveAndRun = () => {}

  return useMemo(
    () => ({
      step,
      error,
      showInspector,
      showDetail,
      payload,

      name,
      mode,
      rawPlan,
      connections,
      scope,
      cronConfig,
      isManual,
      skipOnFail,
      createdDateAfter,

      onChangeStep: setStep,
      onChangeShowInspector: setShowInspector,
      onChangeShowDetail: setShowDetail,

      onChangeName: setName,
      onChangeMode: setMode,
      onChangeRawPlan: setRawPlan,
      onChangeConnections: setConnections,
      onChangeScope: setScope,
      onChangeCronConfig: setCronConfig,
      onChangeIsManual: setIsManual,
      onChangeSkipOnFail: setSkipOnFail,
      onChangeCreatedDateAfter: setCreatedDateAfter,

      onSave: handleSave,
      onSaveAndRun: hanldeSaveAndRun
    }),
    [
      step,
      error,
      connections,
      showInspector,
      showDetail,
      name,
      mode,
      rawPlan,
      connections,
      scope,
      cronConfig,
      isManual,
      skipOnFail,
      createdDateAfter
    ]
  )
}
