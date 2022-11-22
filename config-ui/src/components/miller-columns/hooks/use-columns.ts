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

import { useMemo } from 'react'

import type { ItemType, ItemMapType, ColumnType } from '../types'
import { ItemStatusEnum } from '../types'

interface Props {
  itemMap: ItemMapType
  activeItemId?: ItemType['id']
}

export const useColumns = ({ itemMap, activeItemId }: Props) => {
  return useMemo(() => {
    const rootLeaf = {
      parentId: null,
      activeId: null,
      items: Object.values(itemMap).filter((it) => it.parentId === null),
      hasMore: false
    }

    if (!activeItemId) {
      return [rootLeaf]
    }

    const activeItem = itemMap[activeItemId]

    const columns: ColumnType[] = [
      {
        parentId: activeItem.id,
        items: activeItem.items ?? [],
        activeId: null,
        hasMore: activeItem.status !== ItemStatusEnum.READY
      }
    ]

    const collect = (item: ItemType) => {
      const parent = itemMap[item.parentId ?? '']

      columns.unshift({
        parentId: item.parentId,
        items: parent
          ? parent.items
          : Object.values(itemMap).filter((it) => it.parentId === null),
        activeId: item.id ?? null,
        hasMore: false
      })

      if (parent) {
        collect(parent)
      }
    }

    collect(activeItem)

    return columns
  }, [itemMap, activeItemId])
}
