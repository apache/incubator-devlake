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

import { ItemType, ItemMapType, ColumnType } from '../types'

interface Props {
  items: ItemType[]
  itemMap: ItemMapType
  activeItemId?: ItemType['id']
}

export const useColumns = ({ items, itemMap, activeItemId }: Props) => {
  return useMemo(() => {
    const rootLeaf = { items, activeId: null, parentId: null }

    if (!activeItemId) {
      return [rootLeaf]
    }

    const activeItem = itemMap.getItem(activeItemId)

    const columns: ColumnType[] = [
      {
        parentId: activeItem.id,
        items: activeItem.items,
        activeId: null
      }
    ]

    const collect = (item: ItemType) => {
      const parent = itemMap.getItemParent(item.id)

      columns.unshift({
        parentId: parent?.id ?? null,
        items: parent?.items ?? items,
        activeId: item.id ?? null
      })

      if (parent) {
        collect(parent)
      }
    }

    collect(activeItem)

    return columns
  }, [items, itemMap, activeItemId])
}
