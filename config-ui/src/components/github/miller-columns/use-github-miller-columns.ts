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

import { useState, useEffect, useMemo, useCallback } from 'react'

import type { MillerColumnsItem, ItemType } from '@/components/miller-columns'
import { ItemTypeEnum, ItemStatusEnum } from '@/components/miller-columns'
import request from '@/components/utils/request'

import { getGitHubProxyApiPrefix } from '../config'

export interface UseGitHubMillerColumnsProps {
  connectionId: string
}

export const useGitHubMillerColumns = ({
  connectionId
}: UseGitHubMillerColumnsProps) => {
  const [items, setItems] = useState<Array<MillerColumnsItem>>([])

  const prefix = useMemo(
    () => getGitHubProxyApiPrefix(connectionId),
    [connectionId]
  )

  useEffect(() => {
    ;(async () => {
      const user = await request(`${prefix}/user`)
      const [repos, orgs] = await Promise.all([
        request(user.repos_url),
        request(user.organizations_url)
      ])

      setItems([
        ...orgs.map((it: any) => ({
          parentId: null,
          id: it.id,
          title: it.login,
          type: ItemTypeEnum.BRANCH,
          status: ItemStatusEnum.PENDING
        })),
        {
          parentId: null,
          id: 'owner',
          title: 'owner',
          type: ItemTypeEnum.BRANCH,
          status: ItemStatusEnum.READY
        },
        ...repos.map((it: any) => ({
          parentId: 'owner',
          id: it.id,
          title: it.name,
          owner: user.login,
          repo: it.name
        }))
      ])
    })()
  }, [prefix])

  const onExpandItem = useCallback(
    async (item: ItemType) => {
      if (item.status === ItemStatusEnum.READY) {
        return
      }

      const repos = await request(`${prefix}/orgs/${item.title}/repos`)
      setItems([
        ...items.map((it) =>
          it.id !== item.id
            ? it
            : {
                ...it,
                status: ItemStatusEnum.READY
              }
        ),
        ...repos.map((it: any) => ({
          parentId: item.id,
          id: it.id,
          title: it.name,
          owner: item.title,
          repo: it.name
        }))
      ])
    },
    [items, prefix]
  )

  return useMemo(
    () => ({
      items,
      onExpandItem
    }),
    [items, onExpandItem]
  )
}
