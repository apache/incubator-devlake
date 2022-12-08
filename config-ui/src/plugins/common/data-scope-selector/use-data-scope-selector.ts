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

import * as API from './api'

type ScopeItem = {
  name: string
} & any

export interface UseDataScopeSelector {
  plugin: Plugins
  connectionId: ID
  scopeIds: ID[]
}

export const useDataScopeSelector = ({
  plugin,
  connectionId,
  scopeIds
}: UseDataScopeSelector) => {
  const [loading, setLoading] = useState(false)
  const [scope, setScope] = useState<ScopeItem[]>([])

  const getScopeDetail = async () => {
    setLoading(true)
    try {
      const res = await Promise.all(
        scopeIds.map((id) => API.getDataScopeRepo(plugin, connectionId, id))
      )
      setScope(res)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    getScopeDetail()
  }, [])

  return useMemo(
    () => ({
      loading,
      scope,
      getKey(sc: ScopeItem) {
        switch (true) {
          case plugin === Plugins.GitHub:
            return sc.githubId
        }
      },
      getName(sc: ScopeItem) {
        return sc.name
      }
    }),
    [loading, scope, plugin]
  )
}
