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
import { ItemType, ColumnType } from 'miller-columns-select'

import type { ScopeItemType } from '../../types'
import { useProxyPrefix } from '../../hooks'
import * as API from '../../api'

const DEFAULT_PAGE_SIZE = 30

type ExtraType = {
  type: 'org' | 'repo'
} & ScopeItemType

type GitHubItemType = ItemType<ExtraType>

type GitHubColumnType = ColumnType<ExtraType>

type MapPageType = Record<ID | 'root', number>

export interface UseMillerColumnsProps {
  connectionId: string | number
}

export const useMillerColumns = ({ connectionId }: UseMillerColumnsProps) => {
  const [user, setUser] = useState<any>({})
  const [items, setItems] = useState<GitHubItemType[]>([])
  const [expandedIds, setExpandedIds] = useState<ID[]>([])
  const [loadedIds, setLoadedIds] = useState<ID[]>([])
  const [mapPage, setMapPage] = useState<MapPageType>({})

  const prefix = useProxyPrefix(connectionId)

  const formatOrgs = (orgs: any, parentId: ID | null = null) =>
    orgs.map((it: any) => ({
      parentId,
      id: it.id,
      title: it.login,
      type: 'org'
    }))

  const formatRepos = (repos: any, parentId: ID | null = null) =>
    repos.map((it: any) => ({
      parentId,
      id: it.id,
      title: it.name,
      type: 'repo',
      githubId: it.id,
      name: `${it.owner.login}/${it.name}`
    }))

  const setLoaded = useCallback(
    (loaded: boolean, id: ID, nextPage: number) => {
      if (loaded) {
        setLoadedIds([...loadedIds, id])
      } else {
        setMapPage({ ...mapPage, [`${id}`]: nextPage })
      }
    },
    [loadedIds, mapPage]
  )

  useEffect(() => {
    ;(async () => {
      const user = await API.getUser(prefix)
      const orgs = await API.getUserOrgs(prefix, user.login, {
        page: 1,
        per_page: DEFAULT_PAGE_SIZE
      })

      const loaded = !orgs.length || orgs.length < DEFAULT_PAGE_SIZE

      setUser(user)
      setLoaded(loaded, 'root', 2)
      setItems([
        {
          parentId: null,
          id: user.login,
          title: user.login,
          type: 'org'
        },
        ...formatOrgs(orgs)
      ])
    })()
  }, [prefix])

  const onExpandItem = useCallback(
    async (item: GitHubItemType) => {
      if (expandedIds.includes(item.id)) {
        return
      }

      const isUser = item.id === user.login
      const repos = isUser
        ? await API.getUserRepos(prefix, user.login, {
            page: 1,
            per_page: DEFAULT_PAGE_SIZE
          })
        : await API.getOrgRepos(prefix, item.title, {
            page: 1,
            per_page: DEFAULT_PAGE_SIZE
          })

      const loaded = !repos.length || repos.length < DEFAULT_PAGE_SIZE
      setLoaded(loaded, item.id, 2)
      setExpandedIds([...expandedIds, item.id])
      setItems([...items, ...formatRepos(repos, item.id)])
    },
    [items, prefix]
  )

  const onScrollColumn = async (column: GitHubColumnType) => {
    const page = mapPage[column.parentId ?? 'root']
    const isUser = column.parentId === user.login
    const orgs = !column.parentId
      ? await API.getUserOrgs(prefix, user.login, {
          page,
          per_page: DEFAULT_PAGE_SIZE
        })
      : []

    const repos = column.parentId
      ? isUser
        ? await API.getUserRepos(prefix, user.login, {
            page,
            per_page: DEFAULT_PAGE_SIZE
          })
        : await API.getOrgRepos(prefix, column.parentTitle, {
            page,
            per_page: DEFAULT_PAGE_SIZE
          })
      : []

    const loaded = !column.parentId
      ? !orgs.length || orgs.length < DEFAULT_PAGE_SIZE
      : !repos.length || repos.length < DEFAULT_PAGE_SIZE

    setLoaded(loaded, column.parentId ?? 'root', page + 1)
    setItems([
      ...items,
      ...formatOrgs(orgs),
      ...formatRepos(repos, column.parentId)
    ])
  }

  return useMemo(
    () => ({
      items,
      getHasMore(column: GitHubColumnType) {
        if (loadedIds.includes(column.parentId ?? 'root')) {
          return false
        }
        return true
      },
      onExpandItem,
      onScrollColumn
    }),
    [items, loadedIds]
  )
}
