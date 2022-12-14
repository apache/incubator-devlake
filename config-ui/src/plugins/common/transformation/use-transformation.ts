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

import { useState, useEffect, useMemo } from 'react'

import { Plugins } from '@/plugins'
import { operator } from '@/utils'

import type { RuleItem, ScopeItem } from './types'
import { defaultConfig } from './config'
import * as API from './api'

export interface UseTransformationProps {
  plugin: Plugins
  connectionId: ID
  scopeIds: ID[]
  name: string
  selectedRule?: RuleItem
  selectedScope?: ScopeItem[]
  onSave?: () => void
}

export const useTransformation = ({
  plugin,
  connectionId,
  scopeIds,
  name,
  selectedRule,
  selectedScope,
  onSave
}: UseTransformationProps) => {
  const [loading, setLoading] = useState(false)
  const [rules, setRules] = useState<RuleItem[]>([])
  const [scope, setScope] = useState<ScopeItem[]>([])
  const [saving, setSaving] = useState(false)
  const [transformation, setTransformation] = useState({})

  useEffect(() => {
    setTransformation(selectedRule ? selectedRule : defaultConfig[plugin])
  }, [selectedRule])

  const getRules = async () => {
    const res = await API.getRules(plugin)
    setRules(res)
  }

  const getScope = async () => {
    const res = await Promise.all(
      scopeIds.map((id) => API.getDataScopeRepo(plugin, connectionId, id))
    )
    setScope(res)
  }

  const init = async () => {
    setLoading(true)
    try {
      await getRules()
      await getScope()
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    init()
  }, [])

  const handleUpdateScope = async (tid?: ID) => {
    if (!tid) return

    const payload = (selectedScope ?? []).map((sc: any) => ({
      ...sc,
      transformationRuleId: tid
    }))

    const [success] = await operator(() =>
      API.updateDataScope(plugin, connectionId, {
        data: payload
      })
    )

    if (success) {
      onSave?.()
    }
  }

  const handleSave = async () => {
    const [success, res] = await operator(
      () =>
        API.createTransformation(plugin, {
          ...transformation,
          name
        }),
      {
        setOperating: setSaving
      }
    )

    if (success) {
      const payload = (selectedScope ?? []).map((sc: any) => ({
        ...sc,
        transformationRuleId: res.id
      }))

      if (payload.length) {
        API.updateDataScope(plugin, connectionId, {
          data: payload
        })
      }

      onSave?.()
    }
  }

  return useMemo(
    () => ({
      loading,
      rules,
      scope,
      saving,
      transformation,
      getScopeKey(sc: any) {
        switch (true) {
          case plugin === Plugins.GitHub:
            return sc.githubId
          case plugin === Plugins.JIRA:
            return sc.boardId
          case plugin === Plugins.GitLab:
            return sc.gitlabId
          case plugin === Plugins.Jenkins:
            return sc.fullName
        }
      },
      onSave: handleSave,
      onUpdateScope: handleUpdateScope,
      onChangeTransformation: setTransformation
    }),
    [
      loading,
      rules,
      scope,
      saving,
      transformation,
      plugin,
      selectedScope,
      onSave
    ]
  )
}
