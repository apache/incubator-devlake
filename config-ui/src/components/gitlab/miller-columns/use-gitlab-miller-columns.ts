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

import {
  MillerColumnsItem,
  ItemType,
  ItemStatusEnum
} from '@/components/miller-columns'
import { ItemTypeEnum } from '@/components/miller-columns'
import request from '@/components/utils/request'

import { getGitLabProxyApiPrefix } from '../config'

type MapType = Record<
  MillerColumnsItem['id'],
  {
    groupPage: number
    groupPageSize: number
    groupLoaded: boolean
    projectPage: number
    projectPageSize: number
    projectLoaded: boolean
  }
>

export interface UseGitLabMillerColumnsProps {
  connectionId: string
}

export const useGitLabMillerColumns = <T>({
  connectionId
}: UseGitLabMillerColumnsProps) => {
  const [userId, setUserId] = useState(0)
  const [items, setItems] = useState<Array<MillerColumnsItem>>([])
  const [hasMore, setHasMore] = useState(true)
  const [map, setMap] = useState<MapType>({
    root: {
      groupPage: 1,
      groupPageSize: 20,
      groupLoaded: false,
      projectPage: 1,
      projectPageSize: 20,
      projectLoaded: false
    }
  })

  const prefix = useMemo(
    () => getGitLabProxyApiPrefix(connectionId),
    [connectionId]
  )

  const getRootGroups = (page: number, pageSize: number) => {
    return request(`${prefix}/groups`, {
      data: { top_level_only: 1, page, per_page: pageSize }
    })
  }

  const getRootProjects = (id: number, page: number, pageSize: number) => {
    return request(`${prefix}/users/${id}/projects`, {
      data: { page, per_page: pageSize }
    })
  }

  const getChildGroups = (
    id: MillerColumnsItem['id'],
    page: number,
    pageSize: number
  ) => {
    return request(`${prefix}/groups/${id}/subgroups`, {
      data: { page, per_page: pageSize }
    })
  }

  const getChildProjects = (
    id: MillerColumnsItem['id'],
    page: number,
    pageSize: number
  ) => {
    return request(`${prefix}/groups/${id}/projects`, {
      data: { page, per_page: pageSize }
    })
  }

  const updateGroups = (
    arr: any,
    parentId: MillerColumnsItem['parentId'] = null
  ): Array<ItemType> =>
    arr.map((it: any) => ({
      parentId,
      id: it.id,
      title: it.name,
      type: ItemTypeEnum.BRANCH,
      status: ItemStatusEnum.PENDING
    }))

  const updateProjects = (
    arr: any,
    parentId: MillerColumnsItem['parentId'] = null
  ): Array<ItemType> =>
    arr.map((it: any) => ({
      parentId,
      id: it.id,
      title: it.name
    }))

  useEffect(() => {
    ;(async () => {
      const user = await request(`${prefix}/user`)
      setUserId(user.id)

      const target = map.root

      let projects = []

      const groups = await getRootGroups(target.groupPage, target.groupPageSize)

      if (groups.length < target.groupPageSize) {
        target.groupLoaded = true
        projects = await getRootProjects(
          user.id,
          target.projectPage,
          target.projectPageSize
        )

        if (projects.length < target.projectPageSize) {
          target.projectLoaded = true
          setHasMore(false)
        } else {
          target.projectPage += 1
        }
      } else {
        target.groupPage += 1
      }

      setItems([...updateGroups(groups), ...updateProjects(projects)])
      setMap({
        ...map,
        root: target
      })
    })()
  }, [prefix])

  const onExpandItem = async (item: ItemType) => {
    if (map[item.id]) {
      return
    }

    let target = {
      groupPage: 1,
      groupPageSize: 20,
      groupLoaded: false,
      projectPage: 1,
      projectPageSize: 20,
      projectLoaded: false
    }

    let groups = []
    let projects = []

    groups = await getChildGroups(
      item.id,
      target.groupPage,
      target.groupPageSize
    )

    if (groups.length < target.groupPageSize) {
      target.groupLoaded = true

      projects = await getChildProjects(
        item.id,
        target.projectPage,
        target.projectPageSize
      )

      if (projects.length < target.projectPageSize) {
        target.projectLoaded = true
      } else {
        target.projectPage += 1
      }
    } else {
      target.groupPage += 1
    }

    setItems([
      ...items,
      ...updateGroups(groups, item.id),
      ...updateProjects(projects, item.id)
    ])
    setMap({
      ...map,
      [`${item.id}`]: target
    })
  }

  const onScroll = async (parentId: ItemType['parentId']) => {
    const target = map[parentId ?? 'root']

    let groups = []
    let projects = []

    // All children ready
    if (target.groupLoaded && target.projectLoaded) {
      setItems(
        items.map((it) =>
          it.id !== parentId ? it : { ...it, status: ItemStatusEnum.READY }
        )
      )
      // groups ready
    } else if (target.groupLoaded) {
      projects = parentId
        ? await getChildProjects(
            parentId,
            target.projectPage,
            target.projectPageSize
          )
        : await getRootProjects(
            userId,
            target.projectPage,
            target.projectPageSize
          )

      if (projects.length < target.projectPageSize) {
        target.projectLoaded = true
      } else {
        target.projectPage += 1
      }
      // no group ready
    } else {
      groups = parentId
        ? await getChildGroups(parentId, target.groupPage, target.groupPageSize)
        : await getRootGroups(target.groupPage, target.groupPageSize)

      if (!groups.length) {
        target.groupLoaded = true
        projects = parentId
          ? await getChildProjects(
              parentId,
              target.projectPage,
              target.projectPageSize
            )
          : await getRootProjects(
              userId,
              target.projectPage,
              target.projectPageSize
            )

        if (projects.length < target.projectPageSize) {
          target.projectLoaded = true
        } else {
          target.projectPage += 1
        }
      } else if (groups.length < target.groupPageSize) {
        target.groupLoaded = true
      } else {
        target.groupPage += 1
      }
    }

    setItems([
      ...items.map((it) =>
        it.id !== parentId
          ? it
          : !(target.groupLoaded && target.projectLoaded)
          ? it
          : {
              ...it,
              status: ItemStatusEnum.READY
            }
      ),
      ...updateGroups(groups, parentId),
      ...updateProjects(projects, parentId)
    ])
    setMap({
      ...map,
      [`${parentId}`]: target
    })
  }

  return useMemo(
    () => ({
      items,
      onExpandItem,
      hasMore,
      onScroll
    }),
    [items, map, hasMore]
  )
}
