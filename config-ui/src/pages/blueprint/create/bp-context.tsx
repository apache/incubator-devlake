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

import React, { useState, useMemo, useContext } from 'react'
import { useHistory } from 'react-router-dom'

import type { ConnectionItemType } from '@/store'
import { useConnection, ConnectionStatusEnum } from '@/store'
import { operator } from '@/utils'

import type { BPContextType } from './types'
import { FromEnum, ModeEnum } from './types'
import * as API from './api'

export const BPContext = React.createContext<BPContextType>({
  step: 1,
  error: '',
  showInspector: false,
  showDetail: false,
  payload: {},

  name: 'MY BLUEPRINT',
  mode: ModeEnum.normal,
  rawPlan: JSON.stringify([[]], null, '  '),
  uniqueList: [],
  scopeMap: {},
  cronConfig: '0 0 * * *',
  isManual: false,
  skipOnFail: false,
  createdDateAfter: null,

  onChangeStep: () => {},
  onChangeShowInspector: () => {},
  onChangeShowDetail: () => {},

  onChangeMode: () => {},
  onChangeName: () => {},
  onChangeRawPlan: () => {},
  onChangeUniqueList: () => {},
  onChangeScopeMap: () => {},
  onChangeCronConfig: () => {},
  onChangeIsManual: () => {},
  onChangeSkipOnFail: () => {},
  onChangeCreatedDateAfter: () => {},

  onSave: () => {},
  onSaveAndRun: () => {}
})

interface Props {
  from: FromEnum
  projectName: string
  children: React.ReactNode
}

export const BPContextProvider = ({ from, projectName, children }: Props) => {
  const [step, setStep] = useState(1)
  const [showInspector, setShowInspector] = useState(false)
  const [showDetail, setShowDetail] = useState(false)

  const [name, setName] = useState(
    from === FromEnum.project ? `${projectName}-BLUEPRINT` : 'MY BLUEPRINT'
  )
  const [mode, setMode] = useState<ModeEnum>(ModeEnum.normal)
  const [rawPlan, setRawPlan] = useState(JSON.stringify([[]], null, '  '))
  const [uniqueList, setUniqueList] = useState<string[]>([])
  const [scopeMap, setScopeMap] = useState<Record<string, any>>({})
  const [cronConfig, setCronConfig] = useState('0 0 * * *')
  const [isManual, setIsManual] = useState(false)
  const [skipOnFail, setSkipOnFail] = useState(false)
  const [createdDateAfter, setCreatedDateAfter] = useState<string | null>(null)

  const history = useHistory()

  const { connections } = useConnection()

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
      case mode === ModeEnum.normal && !uniqueList.length:
        return 'Normal Mode: No Data Connections selected.'
      case mode === ModeEnum.normal &&
        !connections
          .filter((cs) => uniqueList.includes(cs.unique))
          .every((cs) => cs.status === ConnectionStatusEnum.ONLINE):
        return 'Normal Mode: Has some offline connections'
      case step === 2 && Object.keys(scopeMap).length !== uniqueList.length:
        return 'No Data Scope is Selected'
      default:
        return ''
    }
  }, [name, mode, rawPlan, uniqueList, connections, step, scopeMap])

  const payload = useMemo(() => {
    const params: any = {
      name,
      projectName,
      mode,
      enable: true,
      cronConfig,
      isManual,
      skipOnFail
    }

    if (mode === ModeEnum.normal) {
      params.settings = {
        version: '2.0.0',
        createdDateAfter,
        connections: uniqueList.map((unique) => {
          const connection = connections.find(
            (cs) => cs.unique === unique
          ) as ConnectionItemType
          const scope = scopeMap[unique] ?? []
          return {
            plugin: connection.plugin,
            connectionId: connection.id,
            scopes: scope.map((sc: any) => ({
              id: `${sc.id}`,
              entities: sc.entities
            }))
          }
        })
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
    cronConfig,
    isManual,
    skipOnFail,
    createdDateAfter,
    rawPlan,
    uniqueList,
    scopeMap,
    connections
  ])

  const handleSaveAfter = (id: ID) => {
    const path =
      from === FromEnum.blueprint
        ? `/blueprints/${id}`
        : `/projects/${projectName}`

    history.push(path)
  }

  const handleSave = async () => {
    const [success, res] = await operator(() => API.createBlueprint(payload))

    if (success) {
      handleSaveAfter(res.id)
    }
  }

  const hanldeSaveAndRun = async () => {
    const [success, res] = await operator(async () => {
      const res = await API.createBlueprint(payload)
      return await API.runBlueprint(res.id)
    })

    if (success) {
      handleSaveAfter(res.id)
    }
  }

  return (
    <BPContext.Provider
      value={{
        step,
        error,
        showInspector,
        showDetail,
        payload,

        name,
        mode,
        rawPlan,
        uniqueList,
        scopeMap,
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
        onChangeUniqueList: setUniqueList,
        onChangeScopeMap: setScopeMap,
        onChangeCronConfig: setCronConfig,
        onChangeIsManual: setIsManual,
        onChangeSkipOnFail: setSkipOnFail,
        onChangeCreatedDateAfter: setCreatedDateAfter,

        onSave: handleSave,
        onSaveAndRun: hanldeSaveAndRun
      }}
    >
      {children}
    </BPContext.Provider>
  )
}

export const useCreateBP = () => {
  return useContext(BPContext)
}
