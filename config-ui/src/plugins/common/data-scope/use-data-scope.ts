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
  allEntities: string[]
  onSaveAfter?: (scope: Array<{ id: ID; entities: string[] }>) => void
}

export const useDataScope = ({
  plugin,
  connectionId,
  allEntities,
  onSaveAfter
}: UseDataScope) => {
  const [saving, setSaving] = useState(false)
  const [scope, setScope] = useState<any>([])
  const [entities, setEntities] = useState<string[]>([])

  useEffect(() => {
    setEntities(allEntities ?? [])
  }, [allEntities])

  const formatScope = (scope: any) => {
    return scope.map((sc: any) => {
      switch (true) {
        case plugin === Plugins.GitHub:
          return {
            id: sc.githubId,
            entities
          }
        default:
          return {
            id: sc.id,
            entities
          }
      }
    })
  }

  const handleSave = async () => {
    const [success, res] = await operator(
      () =>
        API.updateDataScope(plugin, connectionId, {
          data: scope.map((sc: any) => omit(sc, 'from'))
        }),
      {
        setOperating: setSaving
      }
    )

    if (success) {
      onSaveAfter?.(formatScope(res))
    }
  }

  return useMemo(
    () => ({
      saving,
      scope,
      entities,
      onChangeScope: setScope,
      onChangeEntities: setEntities,
      onSave: handleSave
    }),
    [saving, scope]
  )
}
