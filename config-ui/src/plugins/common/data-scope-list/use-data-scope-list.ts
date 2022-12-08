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
import { groupBy } from 'lodash'

import { Plugins } from '@/plugins'

import * as API from './api'

type ScopeItem = {
  id: ID
  name: string
  transformationRuleName?: string
}

export interface UseDataScopeList {
  plugin: Plugins
  connectionId: ID
  scopeIds: ID[]
}

export const useDataScopeList = ({
  plugin,
  connectionId,
  scopeIds
}: UseDataScopeList) => {
  const [loading, setLoading] = useState(false)
  const [scope, setScope] = useState<ScopeItem[]>([])
  const [scopeTsMap, setScopeTsMap] = useState<Record<string, ScopeItem[]>>({})

  useEffect(() => {
    setScopeTsMap(
      groupBy(scope, (it) => it.transformationRuleName ?? 'No Transformation')
    )
  }, [scope])

  const formatScope = (scope: any) => {
    return scope.map((sc: any) => {
      switch (true) {
        case plugin === Plugins.GitHub:
          return {
            id: sc.githubId,
            name: sc.name,
            transformationRuleName: sc.transformationRuleName
          }
      }
    })
  }

  const getScopeDetail = async () => {
    setLoading(true)
    try {
      const res = await Promise.all(
        scopeIds.map((id) => API.getDataScopeRepo(plugin, connectionId, id))
      )
      setScope(formatScope(res))
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    getScopeDetail()
  }, [])

  return useMemo(
    () => ({ loading, scope, scopeTsMap }),
    [loading, scope, scopeTsMap]
  )
}
