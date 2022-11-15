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

import request from '../request'
import { getGitLabProxyApiPrefix } from '../config'

export interface UseGitLabProjectSelectorProps {
  connectionId: string
  search: string
  membership: boolean
}

export const useGitLabProjectSelector = ({
  connectionId,
  search,
  membership
}: UseGitLabProjectSelectorProps) => {
  const [loading, setLoading] = useState(false)
  const [items, setItems] = useState([])

  const prefix = useMemo(
    () => getGitLabProxyApiPrefix(connectionId),
    [connectionId]
  )

  useEffect(() => {
    if (!search) return
    setLoading(true)

    const apiPath = `${prefix}/projects`

    const timer = setTimeout(async () => {
      const res = await request(apiPath, { data: { search, membership } })
      setItems(
        res.map((it: any) => ({
          id: it.id,
          key: it.id,
          title: it.name_with_namespace,
          shortTitle: it.name
        }))
      )
      setLoading(false)
    }, 1000)

    return () => clearTimeout(timer)
  }, [prefix, search, membership])

  return useMemo(() => {
    return {
      loading,
      items
    }
  }, [loading, items])
}
