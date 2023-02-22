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

import NoData from '@/images/no-data.svg';
import { Loading, Card } from '@/components';

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
        <Card>
          <S.NoData>
            <img src={NoData} alt="" />
            <p>{text ?? 'No Data'}</p>
            {onCreate && (
              <Button intent={Intent.PRIMARY} icon="plus" onClick={onCreate}>
                {btnText ?? 'Create'}
              </Button>
            )}
          </S.NoData>
        </Card>
      ) : (
        <Card style={{ padding: 0 }}>
          <S.Table loading={loading ? 1 : 0}>
            <S.Header>
              {columns.map(({ key, align = 'left', title }) => (
                <span key={key} style={{ textAlign: align }}>
                  {title}
                </span>
              ))}
            </S.Header>
            {dataSource.map((data, i) => (
              <S.Row key={i}>
                {columns.map(({ key, align = 'left', dataIndex, render }) => {
                  const value = Array.isArray(dataIndex)
                    ? dataIndex.reduce((acc, cur) => {
                        acc[cur] = data[cur];
                        return acc;
                      }, {} as any)
                    : data[dataIndex];
                  return (
                    <span key={key} style={{ textAlign: align }}>
                      {render ? render(value, data) : value}
                    </span>
                  );
                })}
              </S.Row>
            ))}
          </S.Table>
        </Card>
      )}
    </S.Container>
  );
};
