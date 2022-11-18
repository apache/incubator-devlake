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

import { getJenkinsProxyApiPrefix } from '../config'

export interface UseJenkinsMillerColumnsProps {
  connectionId: string
}

export const useJenkinsMillerColumns = ({
  connectionId
}: UseJenkinsMillerColumnsProps) => {
  const [items, setItems] = useState<Array<MillerColumnsItem>>([])
  const [hasMore, setHasMore] = useState(true)

  const prefix = useMemo(
    () => getJenkinsProxyApiPrefix(connectionId),
    [connectionId]
  )

  useEffect(() => {
    ;(async () => {
      const res = await request(
        `${prefix}/api/json?tree=jobs[name,jobs]{0,10000}`
      )
      setHasMore(false)
      setItems(
        res.jobs.map((it: any) => ({
          parentId: null,
          id: it.name,
          title: it.name,
          type: it.jobs ? ItemTypeEnum.BRANCH : ItemTypeEnum.LEAF,
          status: it.jobs ? ItemStatusEnum.PENDING : ItemStatusEnum.READY
        }))
      )
    })()
  }, [prefix])

  const getJobs = (
    item?: MillerColumnsItem
  ): Array<MillerColumnsItem['id']> => {
    let result = []

    if (item) {
      result.push(item.id)
      result.unshift(...getJobs(items.find((it) => it.id === item.parentId)))
    }
    return result
  }

  const onExpandItem = async (item: MillerColumnsItem) => {
    if (item.status === ItemStatusEnum.READY) {
      return
    }

    const jobs = getJobs(item)
    const res = await request(
      `${prefix}/job/${jobs.join(
        '/job/'
      )}/api/json?tree=jobs[name,jobs]{0,10000}`
    )
    setItems([
      ...items.map((it) =>
        it.id !== item.id
          ? it
          : {
              ...it,
              status: ItemStatusEnum.READY
            }
      ),
      ...res.jobs.map((it: any) => ({
        parentId: item.id,
        id: it.name,
        title: it.name,
        type: it.jobs ? ItemTypeEnum.BRANCH : ItemTypeEnum.LEAF,
        status: it.jobs ? ItemStatusEnum.PENDING : ItemStatusEnum.READY,
        jobPath: `${jobs.join('/')}/`
      }))
    ])
  }

  return useMemo(
    () => ({
      items,
      onExpandItem,
      hasMore
    }),
    [items]
  )
}
