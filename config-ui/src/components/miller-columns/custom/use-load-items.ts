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
  ItemType,
  ItemTypeEnum,
  ItemHasStatusType,
  ItemStatusEnum
} from '../types'

interface Props {
  getInitItems: () => Promise<Array<ItemType>>
  loadMoreItems: (item: ItemType) => Promise<Array<ItemType>>
}

type TreeType = Record<
  ItemType['id'],
  ItemHasStatusType & {
    nameWithNameSpace?: string
  }
>

export const useLoadItems = ({ getInitItems, loadMoreItems }: Props) => {
  const [tree, setTree] = useState<TreeType>({})

  const itemsToTree = (items: Array<ItemType>) => {
    return items.reduce((acc, cur) => {
      acc[cur.id] = {
        ...cur,
        items: [],
        status: ItemStatusEnum.PENDING
      }
      return acc
    }, {} as TreeType)
  }

  const treeToItems = (t: TreeType) => {
    if (!t.root) {
      return []
    }

    const transform = (arr: Array<ItemType>): Array<ItemType> => {
      return arr.map((it) => ({
        ...it,
        items: transform(t[it.id].items)
      }))
    }

    return transform(t.root.items)
  }

  useEffect(() => {
    ;(async () => {
      const initItems = await getInitItems()
      setTree({
        root: {
          id: 'root',
          title: 'root',
          type: ItemTypeEnum.BRANCH,
          status: ItemStatusEnum.READY,
          items: initItems
        },
        ...itemsToTree(initItems)
      })
    })()
  }, [])

  return useMemo(() => {
    return {
      items: treeToItems(tree),
      itemTree: tree,
      async loadItems(item: ItemType) {
        if (tree[item.id].status === ItemStatusEnum.READY) {
          return
        }
        const items = await loadMoreItems(item)
        setTree({
          ...tree,
          [`${item.id}`]: {
            ...item,
            items,
            status: ItemStatusEnum.READY
          },
          ...itemsToTree(items)
        })
      }
    }
  }, [tree])
}
