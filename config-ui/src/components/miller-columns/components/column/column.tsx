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

import React, { useState, useEffect, useCallback } from 'react'
import InfiniteScroll from 'react-infinite-scroll-component'

import type { ItemType } from '../../types'

import * as S from './styled'

export interface ColumnsProps {
  items: Array<ItemType>
  renderItem: (item: ItemType) => React.ReactNode
  height?: number
  title?: string | React.ReactNode
  columnCount?: number
  scrollProps?: {
    hasMore: boolean
    onScroll: () => void
    renderLoader?: () => React.ReactNode
    renderBottom?: () => React.ReactNode
  }
}

export const Column = ({
  items,
  renderItem,
  height,
  title,
  columnCount = 3,
  scrollProps
}: ColumnsProps) => {
  const [hasMore, setHasMore] = useState(true)

  useEffect(() => {
    if (scrollProps) {
      setHasMore(scrollProps.hasMore)
    }
  }, [scrollProps])

  const handleNext = useCallback(() => {
    if (scrollProps) {
      scrollProps.onScroll()
    } else {
      setHasMore(false)
    }
  }, [scrollProps])

  const loader = scrollProps?.renderLoader?.() ?? (
    <S.StatusWrapper>Loading...</S.StatusWrapper>
  )

  const bottom = scrollProps?.renderBottom?.() ?? (
    <S.StatusWrapper>All Data Loaded.</S.StatusWrapper>
  )

  return (
    <S.Container
      id='miller-columns-column-container'
      height={height}
      columnCount={columnCount}
    >
      {title && <div className='title'>{title}</div>}
      <InfiniteScroll
        dataLength={items.length}
        hasMore={hasMore}
        next={handleNext}
        loader={loader}
        scrollableTarget='miller-columns-column-container'
        endMessage={bottom}
      >
        {items.map((it) => renderItem(it))}
      </InfiniteScroll>
    </S.Container>
  )
}
