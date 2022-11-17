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

import type { ItemType } from '@/components/miller-columns'
import { MillerColumns, ItemTypeEnum } from '@/components/miller-columns'

import {
  useGitLabMillerColumns,
  UseGitLabMillerColumnsProps
} from './use-gitlab-miller-columns'

interface Props extends UseGitLabMillerColumnsProps {
  disabledItemIds?: Array<number>
  onChangeItems: (items: Array<Pick<ItemType, 'id' | 'title'>>) => void
}

export const GitLabMillerColumns = ({
  connectionId,
  disabledItemIds,
  onChangeItems
}: Props) => {
  const [seletedIds, setSelectedIds] = useState<Array<ItemType['id']>>([])

  const { items, onExpandItem, hasMore, onScroll } = useGitLabMillerColumns({
    connectionId
  })

  console.log(items)

  useEffect(() => {
    const curItems = items.filter((it) => seletedIds.includes(it.id))
    onChangeItems(curItems)
  }, [seletedIds])

  return (
    <MillerColumns
      height={300}
      firstColumnTitle='Subgroups/Projects'
      items={items}
      disabledItemIds={disabledItemIds}
      selectedItemIds={seletedIds}
      onSelectedItemIds={setSelectedIds}
      onExpandItem={onExpandItem}
      scrollProps={{
        hasMore,
        onScroll
      }}
    />
  )
}
