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
import { Item, ItemAll } from './components'

import { ColumnType } from './types'
import * as S from './styled'

interface Props extends UseMillerColumnsProps {
  height?: number
  firstColumnTitle?: React.ReactNode
  renderColumnBottom?: (col: ColumnType) => React.ReactNode
}

export const MillerColumns = ({
  firstColumnTitle,
  height,
  renderColumnBottom,
  ...props
}: Props) => {
  const {
    columns,
    getStatus,
    getChekecdStatus,
    onExpandItem,
    onSelectItem,
    getCheckedAllStatus,
    onSelectAllItem
  } = useMillerColumns(props)

  return (
    <S.Container>
      {columns.map((col, i) => {
        const bottom = renderColumnBottom?.(col)

        return (
          <S.Column key={col.parentId} height={height}>
            {i === 0 && firstColumnTitle && (
              <div className='title'>{firstColumnTitle}</div>
            )}
            {!!col.items.length && (
              <ItemAll
                column={col}
                checkStatus={getCheckedAllStatus(col)}
                onSelectAllItem={onSelectAllItem}
              />
            )}
            {col.items.map((it) => (
              <Item
                key={it.id}
                item={it}
                status={getStatus(it)}
                checkStatus={getChekecdStatus(it)}
                onExpandItem={onExpandItem}
                onSelectItem={onSelectItem}
              />
            ))}
            {bottom}
          </S.Column>
        )
      })}
    </S.Container>
  )
}
