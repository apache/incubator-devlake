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

import { operator } from '@/utils'
import { PluginConfig } from '@/plugins'

import type { BlueprintType, ConnectionItemType } from './types'
import * as API from './api'

export interface UseDetailProps {
  id: ID
}

export const useDetail = ({ id }: UseDetailProps) => {
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [blueprint, setBlueprint] = useState<BlueprintType>()
  const [connections, setConnections] = useState<ConnectionItemType[]>([])

  const transformConnection = (connections: any) => {
    return connections
      .map((cs: any) => {
        const plugin = PluginConfig.find((p) => p.plugin === cs.plugin)
        if (!plugin) return null
        return {
          icon: plugin.icon,
          name: plugin.name,
          connectionId: cs.connectionId,
          entities: cs.scopes[0].entities,
          plugin: cs.plugin,
          scopeIds: cs.scopes.map((sc: any) => sc.id)
        }
      })
      .filter(Boolean)
  }

  const getBlueprint = async () => {
    setLoading(true)
    try {
      const res = await API.getBlueprint(id)
      setBlueprint(res)
      setConnections(transformConnection(res.settings.connections))
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    getBlueprint()
  }, [])

  const handleUpdate = async (payload: any) => {
    const [success] = await operator(
      () =>
        API.updateBlueprint(id, {
          ...blueprint,
          ...payload
        }),
      {
        setOperating: setSaving
      }
    )

    if (success) {
      getBlueprint()
    }
  }

  return useMemo(
    () => ({
      loading,
      saving,
      blueprint,
      connections,
      onUpdate: handleUpdate
    }),
    [loading, saving, blueprint, connections]
  )
}
