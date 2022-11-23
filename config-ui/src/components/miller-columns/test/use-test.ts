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

import type { MillerColumnsItem, ItemType } from '..'
import { ItemStatusEnum } from '..'

import { flatDataFirst, flatDataSecond, flatDataThird } from './mock'

export const useTest = () => {
  const [items, setItems] = useState<Array<MillerColumnsItem>>([])
  const [hasMore, setHasMore] = useState(true)

  // Get the initial items data
  // And know whether the first column has completed all data loading
  useEffect(() => {
    setItems(flatDataFirst)
    setHasMore(false)
  }, [])

  // Load more data when expanding
  // And judge whether the data is loaded
  const onExpandItem = (item: ItemType) => {
    if (item.id === '2') {
      setItems([...items, ...flatDataSecond])
    }
  }

  const onScroll = (parentId: ItemType['parentId']) => {
    if (parentId === '2') {
      setItems([
        ...items.map((it) =>
          it.id !== '2' ? it : { ...it, status: ItemStatusEnum.READY }
        ),
        ...flatDataThird
      ])
    }
  }

  return useMemo(
    () => ({
      items,
      onExpandItem,
      hasMore,
      onScroll
    }),
    [items, hasMore]
  )
}
