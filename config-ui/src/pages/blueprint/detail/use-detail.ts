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

import { Error } from '@/error'
import { operator } from '@/utils'

import type { BlueprintType } from './types'
import * as API from './api'

export interface UseDetailProps {
  id: ID
}

export const useDetail = ({ id }: UseDetailProps) => {
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [blueprint, setBlueprint] = useState<BlueprintType>()
  const [, setError] = useState()

  const getBlueprint = async () => {
    setLoading(true)
    try {
      const res = await API.getBlueprint(id)

      // need to upgrade 2.0.0
      if (res.settings?.version === '1.0.0') {
        setError(() => {
          throw Error.BP_NEED_TO_UPGRADE
        })
      }

      setBlueprint(res)
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
      onUpdate: handleUpdate,
      onRefresh: getBlueprint
    }),
    [loading, saving, blueprint]
  )
}
