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
import { Button, Intent } from '@blueprintjs/core';

import { Loading, Card, NoData, TextTooltip } from '@/components';

import { ColumnType } from './types';
import * as S from './styled';

interface Props<T> {
  loading?: boolean;
  columns: ColumnType<T>;
  dataSource: T[];
  noData?: {
    text?: React.ReactNode;
    btnText?: string;
    onCreate?: () => void;
  };
}

export const Table = <T extends Record<string, any>>({ loading, columns, dataSource, noData = {} }: Props<T>) => {
  const { text, btnText, onCreate } = noData;

  return (
    <S.Container>
      {loading ? (
        <Card>
          <S.Loading>
            <Loading />
          </S.Loading>
        </Card>
      ) : !dataSource.length ? (
        <NoData
          text={text}
          action={
            onCreate && (
              <Button intent={Intent.PRIMARY} icon="plus" onClick={onCreate}>
                {btnText ?? 'Create'}
              </Button>
            )
          }
        />
      ) : (
        <S.Table>
          <S.THeader>
            <S.TR>
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
                {columns.map(({ key, width, align = 'left', ellipsis, dataIndex, render }) => {
                  const value = Array.isArray(dataIndex)
                    ? dataIndex.reduce((acc, cur) => {
                        acc[cur] = data[cur];
                        return acc;
                      }, {} as any)
                    : data[dataIndex];
                  return (
                    <S.TD key={key} style={{ width, textAlign: align }}>
                      {render ? (
                        render(value, data)
                      ) : ellipsis ? (
                        <TextTooltip content={value}>{value}</TextTooltip>
                      ) : (
                        value
                      )}
                    </S.TD>
                  );
                })}
              </S.TR>
            ))}
          </S.TBody>
        </S.Table>
      )}
    </S.Container>
  );
};
