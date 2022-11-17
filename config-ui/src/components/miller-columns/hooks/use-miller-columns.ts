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

import type { MillerColumnsItem, ItemType, ColumnType } from '../types'
import { ItemTypeEnum, RowStatus } from '../types'
import { CheckStatus } from '../components'

import { useConvertItems } from './use-convert-items'
import { useItemMap } from './use-item-map'
import { useColumns } from './use-columns'

export interface UseMillerColumnsProps {
  items: Array<MillerColumnsItem>
  activeItemId?: ItemType['id']
  onActiveItemId?: (id: ItemType['id']) => void
  disabledItemIds?: Array<ItemType['id']>
  selectedItemIds?: Array<ItemType['id']>
  onSelectedItemIds?: (ids: Array<ItemType['id']>) => void
  onExpandItem?: (item: ItemType) => void
}

export const useMillerColumns = ({
  items,
  disabledItemIds,
  onActiveItemId,
  onSelectedItemIds,
  onExpandItem,
  ...props
}: UseMillerColumnsProps) => {
  const [activeItemId, setActiveItemId] = useState<ItemType['id']>()
  const [selectedItemIds, setSelectedItemIds] = useState<Array<ItemType['id']>>(
    []
  )

  const covertItems = useConvertItems({ items })
  const itemMap = useItemMap({ items: covertItems })
  const columns = useColumns({ itemMap, activeItemId })

  useEffect(() => {
    setActiveItemId(props.activeItemId)
  }, [props.activeItemId])

  useEffect(() => {
    setSelectedItemIds(props.selectedItemIds ?? [])
  }, [props.selectedItemIds])

  const collectAddParentIds = (item: ItemType) => {
    let result: Array<ItemType['id']> = []

    const parentItem = itemMap[item.parentId ?? '']

    if (parentItem) {
      const childSelectedIds = (parentItem.items ?? [])
        .map((it) => it.id)
        .filter((id) => [...selectedItemIds, item.id].includes(id))

      if (childSelectedIds.length === (parentItem.items ?? []).length) {
        result.push(parentItem.id)
        result.push(...collectAddParentIds(parentItem))
      }
    }

    return result
  }

  const collectRemoveParentIds = (item: ItemType) => {
    let result: Array<ItemType['id']> = []

    const parentItem = itemMap[item.parentId ?? '']

    if (parentItem) {
      result.push(parentItem.id)
      result.push(...collectRemoveParentIds(parentItem))
    }

    return result
  }

  return useMemo(
    () => ({
      columns,
      itemMap,
      activeItemId,
      selectedItemIds,
      getStatus(item: ItemType, column: ColumnType) {
        if (item.id === column.activeId) {
          return RowStatus.selected
        }
        return RowStatus.noselected
      },
      getChekecdStatus(item: ItemType) {
        const childSelectedIds = (item.items ?? [])
          .map((it) => it.id)
          .filter((id) => selectedItemIds.includes(id))

        switch (true) {
          case !itemMap[item.id].childLoaded:
          case (disabledItemIds ?? []).includes(item.id):
            return CheckStatus.disabled
          case selectedItemIds.includes(item.id):
            return CheckStatus.checked
          case !!childSelectedIds.length:
            return CheckStatus.indeterminate
          default:
            return CheckStatus.nochecked
        }
      },
      onExpandItem(item: ItemType) {
        if (item.type !== ItemTypeEnum.BRANCH) {
          return
        }
        onExpandItem?.(item)
        onActiveItemId ? onActiveItemId(item.id) : setActiveItemId(item.id)
      },
      onSelectItem(item: ItemType) {
        let newIds: Array<ItemType['id']> = [item.id]
        const isRemoveExistedItem = !!selectedItemIds.includes(item.id)

        const collectChildIds = (it: ItemType) => {
          newIds.push(it.id)
          it.items.forEach((it) => collectChildIds(it))
        }

        item.items.forEach((it) => collectChildIds(it))

        if (!isRemoveExistedItem) {
          newIds = [
            ...new Set([
              ...newIds,
              ...selectedItemIds,
              ...collectAddParentIds(item)
            ])
          ]
        } else {
          newIds = selectedItemIds.filter(
            (id) => ![...newIds, ...collectRemoveParentIds(item)].includes(id)
          )
        }

        onSelectedItemIds
          ? onSelectedItemIds(newIds)
          : setSelectedItemIds(newIds)
      }
    }),
    [columns, itemMap, activeItemId, disabledItemIds, selectedItemIds]
  )
}
