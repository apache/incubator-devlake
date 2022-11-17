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

import React, { useState, useEffect } from 'react'

import type { ItemType } from '@/components/miller-columns'
import { MillerColumns } from '@/components/miller-columns'

import {
  useGitHubMillerColumns,
  UseGitHubMillerColumnsProps
} from './use-github-miller-columns'

interface Props extends UseGitHubMillerColumnsProps {
  onChangeItems: (items: Array<Pick<ItemType, 'id' | 'title'>>) => void
}

export const GitHubMillerColumns = ({ connectionId, onChangeItems }: Props) => {
  const [seletedIds, setSelectedIds] = useState<Array<ItemType['id']>>([])

  const { items, onExpandItem } = useGitHubMillerColumns({
    connectionId
  })

  useEffect(() => {
    onChangeItems(
      items
        .filter((it) => seletedIds.includes(it.id))
        .map((it: any) => {
          return {
            id: it.id,
            title: `${it.owner}/${it.repo}`,
            owner: it.owner,
            repo: it.repo,
            value: `${it.owner}/${it.repo}`,
            type: 'miller-columns'
          }
        })
    )
  }, [seletedIds])

  return (
    <MillerColumns
      height={300}
      columnCount={2}
      firstColumnTitle='Organizations/Owners'
      items={items}
      selectedItemIds={seletedIds}
      onSelectedItemIds={setSelectedIds}
      onExpandItem={onExpandItem}
    />
  )
}
