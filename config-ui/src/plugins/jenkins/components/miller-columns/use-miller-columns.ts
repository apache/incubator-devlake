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
import type { ItemType, ColumnType } from 'miller-columns-select'

import type { ScopeItemType } from '../../types'
import { useProxyPrefix } from '../../hooks'
import * as API from '../../api'

export type JenkinsItemType = ItemType<
  {
    type: 'folder' | 'file'
  } & ScopeItemType
>

export type JenkinsColumnType = ColumnType<
  {
    type: 'folder' | 'file'
  } & ScopeItemType
>

export interface UseMillerColumnsProps {
  connectionId: ID
}

export const useMillerColumns = ({ connectionId }: UseMillerColumnsProps) => {
  const [items, setItems] = useState<JenkinsItemType[]>([])
  const [loadedIds, setLoadedIds] = useState<ID[]>([])
  const [expandedIds, setExpandedIds] = useState<ID[]>([])

  const prefix = useProxyPrefix(connectionId)

  const formatJobs = (jobs: any, parentId: ID | null = null) =>
    jobs.map((it: any) => ({
      parentId,
      id: it.name,
      title: it.name,
      type: it.jobs ? 'folder' : 'file'
    }))

  useEffect(() => {
    ;(async () => {
      const res = await API.getJobs(prefix)
      setItems(formatJobs(res.jobs))
      setLoadedIds(['root'])
    })()
  }, [prefix])

  const getJobs = (item?: JenkinsItemType): ID[] => {
    let result = []

    if (item) {
      result.push(item.id)
      result.unshift(...getJobs(items.find((it) => it.id === item.parentId)))
    }
    return result
  }

  const onExpandItem = async (item: JenkinsItemType) => {
    if (expandedIds.includes(item.id)) {
      return
    }

    const jobs = getJobs(item)
    const res = await API.getJobChildJobs(prefix, jobs.join('/job/'))

    setExpandedIds([...expandedIds, item.id])
    setLoadedIds([...loadedIds, item.id])
    setItems([...items, ...formatJobs(res.jobs, item.id)])
  }

  return useMemo(
    () => ({
      items,
      getHasMore(column: JenkinsColumnType) {
        if (loadedIds.includes(column.parentId ?? 'root')) {
          return false
        }
        return true
      },
      onExpandItem
    }),
    [items, loadedIds]
  )
}
