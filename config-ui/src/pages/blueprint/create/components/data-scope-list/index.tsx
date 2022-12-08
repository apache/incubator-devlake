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

import React, { useMemo } from 'react'
import { Button, Intent } from '@blueprintjs/core'

import { Table, ColumnType } from '@/components'

import type { BPConnectionItemType } from '../../types'
import { useBlueprint } from '../../hooks'

import { Scope } from './scope'
import * as S from './styled'

interface Props {
  loading?: boolean
  groupByTs?: boolean
  onDetail: (connection: BPConnectionItemType) => void
}

export const DataScopeList = ({
  loading = false,
  groupByTs = false,
  onDetail
}: Props) => {
  const { connections, scope } = useBlueprint()

  const columns = useMemo(
    () =>
      [
        {
          title: 'Data Connections',
          dataIndex: ['icon', 'name'],
          key: 'connection',
          render: ({
            icon,
            name
          }: Pick<BPConnectionItemType, 'icon' | 'name'>) => (
            <S.ConnectionColumn>
              <img src={icon} alt='' />
              <span>{name}</span>
            </S.ConnectionColumn>
          )
        },
        {
          title: `Data Scope${groupByTs ? ` and Transformation` : ''}`,
          dataIndex: 'unique',
          key: 'unique',
          render: ({ unique }: Pick<BPConnectionItemType, 'unique'>) => (
            <Scope
              groupByTs={groupByTs}
              scope={scope.filter((sc) => sc.unique === unique)}
            />
          )
        },
        {
          title: '',
          key: 'action',
          align: 'center',
          render: (_: any, connection: BPConnectionItemType) => (
            <Button
              small
              minimal
              intent={Intent.PRIMARY}
              icon='add'
              text={!groupByTs ? 'Add Scope' : 'Add Transformation'}
              onClick={() => onDetail(connection)}
            />
          )
        }
      ] as ColumnType<BPConnectionItemType>,
    [scope]
  )

  return (
    <S.Container>
      <Table loading={loading} columns={columns} dataSource={connections} />
    </S.Container>
  )
}
