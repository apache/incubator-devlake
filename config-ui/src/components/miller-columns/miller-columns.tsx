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

import React from 'react'

import { useMillerColumns, UseMillerColumnsProps } from './hooks'
import { Column, ColumnsProps, Item } from './components'

import * as S from './styled'

interface Props extends UseMillerColumnsProps {
  height?: number
  columnCount?: number
  firstColumnTitle?: React.ReactNode
  scrollProps?: ColumnsProps['scrollProps']
}

export const MillerColumns = ({
  height,
  columnCount,
  firstColumnTitle,
  scrollProps,
  ...props
}: Props) => {
  const { columns, getStatus, getChekecdStatus, onExpandItem, onSelectItem } =
    useMillerColumns(props)

  return (
    <S.Container>
      {columns.map((col, i) => (
        <Column
          key={col.parentId}
          items={col.items}
          renderItem={(item) => (
            <Item
              key={item.id}
              item={item}
              status={getStatus(item, col)}
              checkStatus={getChekecdStatus(item)}
              onExpandItem={onExpandItem}
              onSelectItem={onSelectItem}
            />
          )}
          height={height}
          title={i === 0 && firstColumnTitle}
          columnCount={columnCount}
          scrollProps={scrollProps}
        />
      ))}
    </S.Container>
  )
}
