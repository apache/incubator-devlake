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

import type { ItemType, ColumnType } from '../types'
import { RowStatus } from '../types'
import { CheckStatus } from '../components'

import { useItemMap } from './use-item-map'
import { useColumns } from './use-columns'

export interface UseMillerColumnsProps {
  items: ItemType[]
  activeItemId?: ItemType['id']
  onActiveItemId?: (id: ItemType['id']) => void
  selectedItemIds?: Array<ItemType['id']>
  onSelectedItemIds?: (ids: Array<ItemType['id']>) => void
}

export const useMillerColumns = ({
  items,
  onActiveItemId,
  onSelectedItemIds,
  ...props
}: UseMillerColumnsProps) => {
  const [activeItemId, setActiveItemId] = useState<ItemType['id']>()
  const [selectedItemIds, setSelectedItemIds] =
    useState<Array<ItemType['id']>>()

  const itemMap = useItemMap({ items, selectedItemIds })
  const columns = useColumns({ items, itemMap, activeItemId })

  useEffect(() => {
    setActiveItemId(props.activeItemId)
  }, [props.activeItemId])

  useEffect(() => {
    setSelectedItemIds(props.selectedItemIds)
  }, [props.selectedItemIds])

  return useMemo(
    () => ({
      columns,
      itemMap,
      activeItemId,
      selectedItemIds,
      getStatus(item: ItemType) {
        if (item.id === activeItemId) {
          return RowStatus.selected
        }
        return RowStatus.noselected
      },
      getChekecdStatus(item: ItemType) {
        if (!selectedItemIds?.length) {
          return CheckStatus.nochecked
        }
        if (selectedItemIds?.includes(item.id)) {
          return CheckStatus.checked
        }

        const hasChildCheckedIds = selectedItemIds.filter((id) =>
          item.items?.map((it) => it.id).includes(id)
        )

        if (!hasChildCheckedIds.length) {
          return CheckStatus.nochecked
        }

        if (hasChildCheckedIds.length === item.items?.length) {
          return CheckStatus.checked
        }

        return CheckStatus.indeterminate
      },
      getCheckedAllStatus(column: ColumnType) {
        const itemIds = column.items?.map((it) => it.id) ?? []
        const colSelectedIds = itemIds.filter((id) =>
          selectedItemIds?.includes(id)
        )
        switch (true) {
          case colSelectedIds.length === itemIds.length:
            return CheckStatus.checked
          case !!colSelectedIds.length:
            return CheckStatus.indeterminate
          default:
            return CheckStatus.nochecked
        }
      },
      getCheckedCount(item: ItemType) {
        return itemMap.getItemSelectedChildCount(item.id)
      },
      onExpandItem(item: ItemType) {
        if (!item.items?.length) {
          return
        }
        onActiveItemId ? onActiveItemId(item.id) : setActiveItemId(item.id)
      },
      onSelectItem(item: ItemType) {
        let newIds: Array<ItemType['id']>
        let targetIds: Array<ItemType['id']> = [item.id]
        const itemIds = item.items?.map((it) => it.id) ?? []

        const collect = (id: ItemType['id']) => {
          targetIds.push(id)
          const item = itemMap.getItem(id)
          if (item.items) {
            item.items.forEach((it) => collect(it.id))
          }
        }

        itemIds.forEach((id) => collect(id))

        const isRemoveExistedItem = !!selectedItemIds?.includes(item.id)

        if (isRemoveExistedItem) {
          const parentItem = itemMap.getItemParent(item.id)
          const deleteIds = [parentItem?.id, ...targetIds].filter(Boolean)
          newIds =
            selectedItemIds?.filter((id) => !deleteIds.includes(id)) ?? []
        } else {
          const parentItem = itemMap.getItemParent(item.id)
          const addIds = targetIds.filter(
            (id) => !selectedItemIds?.includes(id)
          )

          if (parentItem) {
            const parentChildIds = parentItem.items?.map((it) => it.id) ?? []
            const parentSelectedIds = parentChildIds.filter((id) =>
              [...(selectedItemIds ?? []), item.id].includes(id)
            )

            const isAllChildSelected =
              parentSelectedIds.length === parentItem?.items?.length

            if (isAllChildSelected) {
              addIds.push(parentItem.id)
            }
          }

          newIds = [...(selectedItemIds ?? []), ...addIds]
        }

        onSelectedItemIds
          ? onSelectedItemIds(newIds)
          : setSelectedItemIds(newIds)
      },
      onSelectAllItem(column: ColumnType) {
        let newIds: Array<ItemType['id']>
        let targetIds: Array<ItemType['id']> = []
        const itemIds = column.items?.map((it) => it.id) ?? []

        const collect = (id: ItemType['id']) => {
          targetIds.push(id)
          const item = itemMap.getItem(id)
          if (item.items) {
            item.items.forEach((it) => collect(it.id))
          }
        }

        itemIds.forEach((id) => collect(id))

        const isRemoveExistedItems =
          itemIds.filter((id) => selectedItemIds?.includes(id)).length ===
          itemIds.length

        if (isRemoveExistedItems) {
          const deleteIds = [...targetIds, column.parentId].filter(Boolean)
          newIds =
            selectedItemIds?.filter((id) => !deleteIds.includes(id)) ?? []
        } else {
          const addIds = targetIds.filter(
            (id) => !selectedItemIds?.includes(id)
          )

          if (column.parentId) {
            addIds.push(column.parentId)
          }

          newIds = [...(selectedItemIds ?? []), ...addIds]
        }

        onSelectedItemIds
          ? onSelectedItemIds(newIds)
          : setSelectedItemIds(newIds)
      }
    }),
    [columns, itemMap, activeItemId, selectedItemIds]
  )
}
