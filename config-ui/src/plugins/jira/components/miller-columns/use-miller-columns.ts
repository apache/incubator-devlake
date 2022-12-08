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
import type { ItemType } from 'miller-columns-select'

import type { ScopeItemType } from '../../types'
import { useProxyPrefix } from '../../hooks'
import * as API from '../../api'

const DEFAULT_PAGE_SIZE = 50

type JIRAItemType = ItemType<ScopeItemType>

export interface UseMillerColumnsProps {
  connectionId: ID
}

export const useMillerColumns = ({ connectionId }: UseMillerColumnsProps) => {
  const [items, setItems] = useState<JIRAItemType[]>([])
  const [isLast, setIsLast] = useState(false)
  const [page, setPage] = useState(1)

  const prefix = useProxyPrefix(connectionId)

  const updateItems = (arr: any) =>
    arr.map((it: any) => ({
      parentId: null,
      id: it.id,
      title: it.name,
      connectionId,
      boardId: it.id,
      name: it.name,
      self: it.self,
      type: it.type,
      projectId: it?.location?.projectId
    }))

  useEffect(() => {
    ;(async () => {
      const res = await API.getBoards(prefix, {
        startAt: (page - 1) * DEFAULT_PAGE_SIZE,
        maxResults: DEFAULT_PAGE_SIZE
      })
      setIsLast(res.isLast)
      setItems([...items, ...updateItems(res.values)])
    })()
  }, [prefix, page])

  return useMemo(
    () => ({
      items,
      getHasMore() {
        return !isLast
      },
      onScrollColumn() {
        setPage(page + 1)
      }
    }),
    [items, isLast, page]
  )
}
