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

import { Loading } from '@/components'

import { ColumnType } from './types'
import * as S from './styled'

interface Props<T> {
  loading?: boolean
  columns: ColumnType<T>
  dataSource: T[]
}

export const Table = <T extends Record<string, any>>({
  loading,
  columns,
  dataSource
}: Props<T>) => {
  return (
    <S.Container>
      <S.TableWrapper loading={loading}>
        <S.TableHeader>
          {columns.map(({ key, title }) => (
            <span key={key}>{title}</span>
          ))}
        </S.TableHeader>
        {dataSource.map((data, i) => (
          <S.TableRow key={i}>
            {columns.map(({ key, align = 'left', dataIndex, render }) => {
              const value = Array.isArray(dataIndex)
                ? dataIndex.reduce((acc, cur) => {
                    acc[cur] = data[cur]
                    return acc
                  }, {} as any)
                : data[dataIndex]
              return (
                <span key={key} style={{ textAlign: align }}>
                  {render ? render(value, data) : value}
                </span>
              )
            })}
          </S.TableRow>
        ))}
      </S.TableWrapper>
      {loading && (
        <S.TableMask>
          <Loading />
        </S.TableMask>
      )}
    </S.Container>
  )
}
