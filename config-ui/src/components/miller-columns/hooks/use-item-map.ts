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

import { useState, useMemo, useEffect } from 'react'

import type { MillerColumnsItem, ItemType, ItemMapType } from '../types'
import { ItemTypeEnum, ItemStatusEnum } from '../types'

interface Props {
  items: Array<MillerColumnsItem>
}

export const useItemMap = ({ items }: Props) => {
  const [itemMap, setItemMap] = useState<ItemMapType>({})

  const checkChildLoaded = (item: MillerColumnsItem): boolean => {
    if (item.status === ItemStatusEnum.PENDING) {
      return false
    }

    return !items
      .filter((it) => it.parentId === item.id)
      .find((it) => it.status === ItemStatusEnum.PENDING)
  }

  const covertItem = (item: MillerColumnsItem): ItemType => {
    const type = item.type
      ? item.type
      : (item.items ?? []).length
      ? ItemTypeEnum.BRANCH
      : ItemTypeEnum.LEAF
    const status = item.status ? item.status : ItemStatusEnum.READY
    return {
      ...item,
      type,
      status,
      childLoaded: checkChildLoaded(item)
    } as ItemType
  }

  const collectChildItems = (
    items: Array<MillerColumnsItem>,
    item: MillerColumnsItem
  ): Array<ItemType> => {
    return items
      .filter((it) => {
        return it.parentId === item.id
      })
      .map((it) =>
        covertItem({
          ...it,
          items: collectChildItems(items, it)
        })
      )
  }

  const itemsToMap = (items: Array<MillerColumnsItem>): ItemMapType => {
    return items.reduce((acc, cur) => {
      acc[cur.id] = covertItem({
        ...cur,
        items: collectChildItems(items, cur)
      })
      return acc
    }, {} as any)
  }

  useEffect(() => {
    setItemMap(itemsToMap(items))
  }, [items])

  return useMemo(() => itemMap, [itemMap])
}
