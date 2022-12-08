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
import { ItemType, ColumnType } from 'miller-columns-select'

import type { ScopeItemType } from '../../types'
import { useProxyPrefix } from '../../hooks'
import * as API from '../../api'

const DEFAULT_PAGE_SIZE = 20

export type GitLabItemType = ItemType<
  {
    type: 'group' | 'project'
  } & ScopeItemType
>

export type GitLabColumnType = ColumnType<
  {
    type: 'group' | 'project'
  } & ScopeItemType
>

type MapValueType = {
  groupPage: number
  groupLoaded: boolean
  projectPage: number
  projectLoaded: boolean
}

type MapType = Record<ID, MapValueType>

export interface UseMillerColumnsProps {
  connectionId: ID
}

export const useMillerColumns = <T>({
  connectionId
}: UseMillerColumnsProps) => {
  const [user, setUser] = useState<any>({})
  const [items, setItems] = useState<GitLabItemType[]>([])
  const [expandedIds, setExpandedIds] = useState<ID[]>([])
  const [map, setMap] = useState<MapType>({})

  const prefix = useProxyPrefix(connectionId)

  const formatGroups = (
    arr: any,
    parentId: ID | null = null
  ): GitLabItemType[] =>
    arr.map((it: any) => ({
      parentId,
      id: it.id,
      title: it.name,
      type: 'group'
    }))

  const formatProjects = (
    arr: any,
    parentId: ID | null = null
  ): GitLabItemType[] =>
    arr.map((it: any) => ({
      parentId,
      id: it.id,
      title: it.name,
      type: 'project',
      gitlabId: it.id,
      name: it.path_with_namespace
    }))

  const setLoaded = (id: ID, params: MapValueType) => {
    setMap({
      ...map,
      [`${id}`]: params
    })
  }

  useEffect(() => {
    ;(async () => {
      const user = await API.getUser(prefix)
      setUser(user)

      let groupLoaded = false
      let projectLoaded = false
      let groups = []
      let projects = []

      groups = await API.getUserGroups(prefix, {
        page: 1,
        per_page: DEFAULT_PAGE_SIZE
      })

      groupLoaded = !groups.length || groups.length < DEFAULT_PAGE_SIZE

      if (groupLoaded) {
        projects = await API.getUserProjects(prefix, user.id, {
          page: 1,
          per_page: DEFAULT_PAGE_SIZE
        })

        projectLoaded = !projects.length || projects.length < DEFAULT_PAGE_SIZE
      }

      setLoaded('root', {
        groupLoaded,
        groupPage: groupLoaded ? 1 : 2,
        projectLoaded,
        projectPage: projectLoaded ? 1 : 2
      })
      setItems([...formatGroups(groups), ...formatProjects(projects)])
    })()
  }, [prefix])

  const onExpandItem = async (item: GitLabItemType) => {
    if (expandedIds.includes(item.id)) {
      return
    }

    let groupLoaded = false
    let projectLoaded = false
    let groups = []
    let projects = []

    groups = await API.getGroupSubgroups(prefix, item.id, {
      page: 1,
      per_page: DEFAULT_PAGE_SIZE
    })

    groupLoaded = !groups.length || groups.length < DEFAULT_PAGE_SIZE

    if (groupLoaded) {
      projects = await API.getGroupProjects(prefix, item.id, {
        page: 1,
        per_page: DEFAULT_PAGE_SIZE
      })

      projectLoaded = !projects.length || projects.length < DEFAULT_PAGE_SIZE
    }

    setLoaded(item.id, {
      groupLoaded,
      groupPage: groupLoaded ? 1 : 2,
      projectLoaded,
      projectPage: projectLoaded ? 1 : 2
    })
    setExpandedIds([...expandedIds, item.id])
    setItems([
      ...items,
      ...formatGroups(groups, item.id),
      ...formatProjects(projects, item.id)
    ])
  }

  const onScrollColumn = async (column: GitLabColumnType) => {
    const mapValue = map[column.parentId ?? 'root']

    let groupLoaded = mapValue.groupLoaded
    let projectLoaded = mapValue.projectLoaded
    let groups = []
    let projects = []

    if (!groupLoaded) {
      groups = column.parentId
        ? await API.getGroupSubgroups(prefix, column.parentId, {
            page: mapValue.groupPage,
            per_page: DEFAULT_PAGE_SIZE
          })
        : await API.getUserGroups(prefix, {
            page: mapValue.groupPage,
            per_page: DEFAULT_PAGE_SIZE
          })

      groupLoaded = !groups.length || groups.length < DEFAULT_PAGE_SIZE
    } else if (!projectLoaded) {
      projects = column.parentId
        ? await API.getGroupProjects(prefix, column.parentId, {
            page: mapValue.projectPage,
            per_page: DEFAULT_PAGE_SIZE
          })
        : await API.getUserProjects(prefix, user.id, {
            page: mapValue.projectPage,
            per_page: DEFAULT_PAGE_SIZE
          })

      projectLoaded = !projects.length || projects.length < DEFAULT_PAGE_SIZE
    }

    setLoaded(column.parentId ?? 'root', {
      groupLoaded,
      groupPage: groupLoaded ? mapValue.groupPage : mapValue.groupPage + 1,
      projectLoaded,
      projectPage: projectLoaded
        ? mapValue.projectPage
        : mapValue.projectPage + 1
    })
    setItems([
      ...items,
      ...formatGroups(groups, column.parentId),
      ...formatProjects(projects, column.parentId)
    ])
  }

  return useMemo(
    () => ({
      items,
      getHasMore(column: GitLabColumnType) {
        const mapValue = map[column.parentId ?? 'root']
        if (mapValue?.groupLoaded && mapValue?.projectLoaded) {
          return false
        }
        return true
      },
      onExpandItem,
      onScrollColumn
    }),
    [items, map]
  )
}
