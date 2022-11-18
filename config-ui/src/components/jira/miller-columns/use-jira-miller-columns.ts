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

import type { MillerColumnsItem } from '@/components/miller-columns'
import { ItemTypeEnum, ItemStatusEnum } from '@/components/miller-columns'

import request from '@/components/utils/request'

import { getJIRAApiPrefix } from '../config'

export interface UseJIRAMillerColumnsProps {
  connectionId: string
}

export const useJIRAMillerColumns = ({
  connectionId
}: UseJIRAMillerColumnsProps) => {
  const [items, setItems] = useState<Array<MillerColumnsItem>>([])
  const [hasMore, setHasMore] = useState(true)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(50)

  const prefix = useMemo(() => getJIRAApiPrefix(connectionId), [connectionId])

  const updateItems = (arr: Array<{ id: number; name: string }>) =>
    arr.map((it) => ({
      parentId: null,
      id: it.id,
      title: it.name,
      type: ItemTypeEnum.LEAF,
      status: ItemStatusEnum.READY,
      items: []
    }))

  useEffect(() => {
    ;(async () => {
      const res = await request(`${prefix}/agile/1.0/board`, {
        data: { startAt: (page - 1) * pageSize, maxResults: pageSize }
      })
      setHasMore(!res.isLast)
      setItems([...items, ...updateItems(res.values)])
    })()
  }, [prefix, page, pageSize])

  return useMemo(
    () => ({
      items,
      hasMore,
      onScroll() {
        setPage(page + 1)
      }
    }),
    [items, hasMore]
  )
}
