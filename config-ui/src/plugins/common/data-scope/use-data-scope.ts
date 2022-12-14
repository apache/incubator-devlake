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
import { omit } from 'lodash'

import { Plugins } from '@/plugins'
import { operator } from '@/utils'

import * as API from './api'

export interface UseDataScope {
  plugin: string
  connectionId: ID
  entities: string[]
  initialValues?: any
  onSave?: (scope: any) => void
}

export const useDataScope = ({
  plugin,
  connectionId,
  entities,
  initialValues,
  onSave
}: UseDataScope) => {
  const [saving, setSaving] = useState(false)
  const [selectedScope, setSelectedScope] = useState<any>([])
  const [selectedEntities, setSelectedEntities] = useState<string[]>([])

  useEffect(() => {
    setSelectedScope(initialValues ?? [])
  }, [initialValues])

  useEffect(() => {
    setSelectedEntities(entities ?? [])
  }, [entities])

  const formatScope = (scope: any) => {
    return scope.map((sc: any) => {
      switch (true) {
        case plugin === Plugins.GitHub:
          return {
            ...sc,
            id: sc.githubId,
            entities: selectedEntities
          }
        case plugin === Plugins.JIRA:
          return {
            ...sc,
            id: sc.boardId,
            entities: selectedEntities
          }
        case plugin === Plugins.GitLab:
          return {
            ...sc,
            id: sc.gitlabId,
            entities: selectedEntities
          }
        case plugin === Plugins.Jenkins:
          return {
            ...sc,
            id: sc.fullName,
            entities: selectedEntities
          }
      }
    })
  }

  const handleSave = async () => {
    const [success, res] = await operator(
      () =>
        API.updateDataScope(plugin, connectionId, {
          data: selectedScope.map((sc: any) => omit(sc, 'from'))
        }),
      {
        setOperating: setSaving
      }
    )

    if (success) {
      onSave?.(formatScope(res))
    }
  }

  return useMemo(
    () => ({
      saving,
      selectedScope,
      selectedEntities,
      onChangeScope: setSelectedScope,
      onChangeEntites: setSelectedEntities,
      onSave: handleSave
    }),
    [saving, selectedScope, selectedEntities]
  )
}
