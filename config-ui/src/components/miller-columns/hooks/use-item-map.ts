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

import type { ItemType, ItemInfoType } from '../types'

interface Props {
  items: ItemType[]
  selectedItemIds?: Array<ItemType['id']>
}

export const useItemMap = ({ items, selectedItemIds = [] }: Props) => {
  return useMemo(() => {
    const itemMap = new Map<ItemType['id'], ItemInfoType>()

    const collect = ({
      item,
      parent
    }: {
      item: ItemType
      parent?: ItemType
    }) => {
      if (!itemMap.has(item.id)) {
        itemMap.set(item.id, {
          item,
          parentId: parent?.id
        })
      }

      if (item.items) {
        item.items.forEach((it) => collect({ item: it, parent: item }))
      }
    }

    items.forEach((it) => collect({ item: it }))

    return {
      getItem(id: ItemType['id']) {
        return (itemMap.get(id) as ItemInfoType).item
      },
      getItemParent(id: ItemType['id']) {
        const parentId = itemMap.get(id)?.parentId
        return parentId ? (itemMap.get(parentId) as ItemInfoType).item : null
      },
      getItemMapSize() {
        return itemMap.size
      }
    }
  }, [items, selectedItemIds])
}
