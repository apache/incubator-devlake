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

import type { WebhookItemType } from '../types'
import * as API from '../api'

export interface UseConnectionProps {
  filterIds?: ID[]
}

export const useConnection = ({ filterIds }: UseConnectionProps) => {
  const [loading, setLoading] = useState(false)
  const [connections, setConnections] = useState<WebhookItemType[]>([])

  const getConnections = async () => {
    setLoading(true)
    try {
      const res = await API.getConnections()
      setConnections(
        res.filter((cs: any) => (filterIds ? filterIds.includes(cs.id) : true))
      )
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    getConnections()
  }, [])

  return useMemo(
    () => ({
      loading,
      connections,
      onRefresh: getConnections
    }),
    [loading, connections]
  )
}
