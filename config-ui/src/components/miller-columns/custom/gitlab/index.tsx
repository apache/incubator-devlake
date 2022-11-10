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

import React, { useEffect, useState } from 'react'

import { MillerColumns } from '../../miller-columns'
import { ColumnType, ItemStatusEnum, ItemType, ItemTypeEnum } from '../../types'

import {
  useGitLabMillerColumns,
  UseGitLabMillerColumnsProps
} from './use-gitlab-miller-columns'

interface Props extends UseGitLabMillerColumnsProps {
  onChangeItems: (
    items: Array<Pick<ItemType, 'id' | 'title'> & { shortTitle: string }>
  ) => void
}

export const GitLabMillerColumns = ({ connectionId, onChangeItems }: Props) => {
  const [seletedIds, setSelectedIds] = useState<Array<ItemType['id']>>([])

  const { items, itemTree, onExpandItem } = useGitLabMillerColumns({
    connectionId
  })

  useEffect(() => {
    const curItems = seletedIds
      .filter((id) => itemTree[id].type === ItemTypeEnum.LEAF)
      .map((id) => ({
        id,
        title: itemTree[id].nameWithNameSpace ?? '',
        shortTitle: itemTree[id].title
      }))

    onChangeItems(curItems)
  }, [seletedIds])

  const renderColumnBottom = ({
    isLoading,
    isEmpty
  }: {
    isLoading: boolean
    isEmpty: boolean
  }) => {
    switch (true) {
      case isLoading:
        return <div>Loading...</div>
      case isEmpty:
        return <div>No Data.</div>
    }
  }

  return (
    <MillerColumns
      height={300}
      items={items}
      selectedItemIds={seletedIds}
      onSelectedItemIds={setSelectedIds}
      onExpandItem={onExpandItem}
      renderColumnBottom={(col: ColumnType) => {
        if (!col.parentId) {
          return renderColumnBottom({
            isLoading: !itemTree.root,
            isEmpty: !itemTree.root || !itemTree.root.items.length
          })
        } else {
          return renderColumnBottom({
            isLoading:
              !itemTree[col.parentId] ||
              itemTree[col.parentId].status === ItemStatusEnum.PENDING,
            isEmpty:
              !itemTree[col.parentId] || !itemTree[col.parentId].items.length
          })
        }
      }}
    />
  )
}
