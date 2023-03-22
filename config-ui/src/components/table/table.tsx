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

import React from 'react';
import { Checkbox, Radio } from '@blueprintjs/core';

import { TextTooltip } from '@/components';

import { ColumnType } from './types';
import { TableLoading, TableNoData } from './components';
import { useRowSelection, UseRowSelectionProps } from './hooks';
import * as S from './styled';

interface Props<T> extends UseRowSelectionProps<T> {
  loading?: boolean;
  columns: ColumnType<T>;
  dataSource: T[];
  noData?: {
    text?: React.ReactNode;
    btnText?: string;
    onCreate?: () => void;
  };
  noShadow?: boolean;
}

export const Table = <T extends Record<string, any>>({
  loading,
  columns,
  dataSource,
  noData = {},
  rowSelection,
  noShadow = false,
}: Props<T>) => {
  const { canSelection, selectionType, getCheckedAll, onCheckedAll, getChecked, onChecked } = useRowSelection<T>({
    dataSource,
    rowSelection,
  });

  if (loading) {
    return <TableLoading />;
  }

  if (!dataSource.length) {
    return <TableNoData {...noData} />;
  }

  return (
    <S.Table
      style={{
        boxShadow: noShadow ? 'none' : '0px 2.4px 4.8px -0.8px rgba(0, 0, 0, 0.1), 0px 1.6px 8px rgba(0, 0, 0, 0.07)',
      }}
    >
      <S.THeader>
        <S.TR>
          {canSelection && (
            <S.TH style={{ width: 40, textAlign: 'center' }}>
              {selectionType === 'checkbox' && <Checkbox checked={getCheckedAll()} onChange={() => onCheckedAll()} />}
            </S.TH>
          )}
          {columns.map(({ key, width, align = 'left', title }) => (
            <S.TH key={key} style={{ width, textAlign: align }}>
              {title}
            </S.TH>
          ))}
        </S.TR>
      </S.THeader>
      <S.TBody>
        {dataSource.map((data, i) => (
          <S.TR key={i}>
            {canSelection && (
              <S.TD style={{ width: 40, textAlign: 'center' }}>
                {selectionType === 'checkbox' && (
                  <Checkbox checked={getChecked(data)} onChange={() => onChecked(data)} />
                )}
                {selectionType === 'radio' && <Radio checked={getChecked(data)} onChange={() => onChecked(data)} />}
              </S.TD>
            )}
            {columns.map(({ key, width, align = 'left', ellipsis, dataIndex, render }) => {
              const value = Array.isArray(dataIndex)
                ? dataIndex.reduce((acc, cur) => {
                    acc[cur] = data[cur];
                    return acc;
                  }, {} as any)
                : data[dataIndex];
              return (
                <S.TD key={key} style={{ width, textAlign: align }}>
                  {render ? render(value, data) : ellipsis ? <TextTooltip content={value}>{value}</TextTooltip> : value}
                </S.TD>
              );
            })}
          </S.TR>
        ))}
      </S.TBody>
    </S.Table>
  );
};
