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

type MapType = Record<
  MillerColumnsItem['id'],
  {
    page: number
    pageSize: number
    loaded: boolean
  }
>

export interface UseGitHubMillerColumnsProps {
  connectionId: string
}

export const useGitHubMillerColumns = ({
  connectionId
}: UseGitHubMillerColumnsProps) => {
  const [user, setUser] = useState<any>({})
  const [items, setItems] = useState<Array<MillerColumnsItem>>([])
  const [hasMore, setHasMore] = useState(true)
  const [map, setMap] = useState<MapType>({
    root: {
      page: 1,
      pageSize: 30,
      loaded: false
    }
  })

  const prefix = useMemo(
    () => getGitHubProxyApiPrefix(connectionId),
    [connectionId]
  )

  const getUserOrgs = (username: string, page: number, pageSize: number) => {
    return request(`${prefix}/users/${username}/orgs`, {
      data: { page, per_page: pageSize }
    })
  }

  const getUserRepos = (username: string, page: number, pageSize: number) => {
    return request(`${prefix}/users/${username}/repos`, {
      data: { page, per_page: pageSize }
    })
  }

  const getOrgRepos = (org: string, page: number, pageSize: number) => {
    return request(`${prefix}/orgs/${org}/repos`, {
      data: { page, per_page: pageSize }
    })
  }

  useEffect(() => {
    ;(async () => {
      const params = map.root

      const user = await request(`${prefix}/user`)
      const orgs = await getUserOrgs(user.login, params.page, params.pageSize)

      if (orgs.length < params.pageSize) {
        setHasMore(false)
        params.loaded = true
      } else {
        params.page += 1
      }

      setUser(user)
      setMap({
        ...map,
        root: params
      })
      setItems([
        {
          parentId: null,
          id: user.login,
          title: user.login,
          type: ItemTypeEnum.BRANCH,
          status: ItemStatusEnum.PENDING
        },
        ...orgs.map((it: any) => ({
          parentId: null,
          id: it.id,
          title: it.login,
          type: ItemTypeEnum.BRANCH,
          status: ItemStatusEnum.PENDING
        }))
      ])
    })()
  }, [prefix])

  const onExpandItem = useCallback(
    async (item: ItemType) => {
      if (map[item.id]) {
        return
      }

      let params = {
        page: 1,
        pageSize: 30,
        loaded: false
      }

      const isUser = item.id === user.login
      const repos = isUser
        ? await getUserRepos(item.id as string, params.page, params.pageSize)
        : await getOrgRepos(item.title, params.page, params.pageSize)

      if (repos.length < params.pageSize) {
        params.loaded = true
      } else {
        params.page += 1
      }

      setMap({
        ...map,
        [`${item.id}`]: params
      })
      setItems([
        ...items.map((it) =>
          it.id !== item.id
            ? it
            : !params.loaded
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
          type: ItemTypeEnum.LEAF,
          status: ItemStatusEnum.READY,
          owner: it.owner?.login,
          repo: it.name
        }))
      ])
    },
    [items, prefix, map]
  )

  const onScroll = async (parentId: MillerColumnsItem['parentId']) => {
    const params = map[parentId ?? 'root']

    if (params.loaded) {
      setItems(
        items.map((it) =>
          it.id !== parentId
            ? it
            : {
                ...it,
                status: ItemStatusEnum.READY
              }
        )
      )
    } else {
      const isUser = parentId === user.login
      const org = items.find((it) => it.id === parentId)
      const repos = !parentId
        ? await getUserOrgs(user.login, params.page, params.pageSize)
        : isUser
        ? await getUserRepos(org?.title as string, params.page, params.pageSize)
        : await getOrgRepos(org?.title as string, params.page, params.pageSize)

      if (!repos.length || repos.length < params.pageSize) {
        setHasMore(false)
        params.loaded = true
      } else {
        params.page += 1
      }

      setMap({
        ...map,
        [`${parentId ?? 'root'}`]: params
      })
      setItems([
        ...items.map((it) =>
          it.id !== parentId
            ? it
            : !params.loaded
            ? it
            : {
                ...it,
                status: ItemStatusEnum.READY
              }
        ),
        ...repos.map((it: any) => ({
          parentId,
          id: it.id,
          title: it.name,
          type: ItemTypeEnum.LEAF,
          status: ItemStatusEnum.READY,
          owner: it.owner.login,
          repo: it.name
        }))
      ])
    }
  }

  return useMemo(
    () => ({
      items,
      onExpandItem,
      hasMore,
      onScroll
    }),
    [items, hasMore]
  )
}
