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

import React, { useMemo } from 'react';
import { useHistory } from 'react-router-dom';
import { ButtonGroup, Button, Intent } from '@blueprintjs/core';

import { PageLoading, PageHeader, Table, ColumnType } from '@/components';
import { getCron, getCronOptions } from '@/config';
import { formatTime } from '@/utils';

import type { BlueprintType } from '../types';
import { ModeEnum } from '../types';

import { useHome } from './use-home';
import * as S from './styled';

export const BlueprintHomePage = () => {
  const history = useHistory();

  const { loading, dataSource, type, onChangeType } = useHome();

  const options = useMemo(() => getCronOptions(), []);

  const columns = useMemo(
    () =>
      [
        {
          title: 'Blueprint Name',
          dataIndex: 'name',
          key: 'name',
        },
        {
          title: 'Data Connections',
          key: 'connections',
          render: (_, row) => {
            if (row.mode === ModeEnum.advanced) {
              return 'Advanced Mode';
            }
            return row.settings.connections.map((cs) => cs.plugin).join(',');
          },
        },
        {
          title: 'Frequency',
          key: 'frequency',
          render: (_, row) => {
            const cron = getCron(row.isManual, row.cronConfig);
            return cron.label;
          },
        },
        {
          title: 'Next Run Time',
          key: 'nextRunTime',
          render: (_, row) => {
            const cron = getCron(row.isManual, row.cronConfig);
            return formatTime(cron.nextTime);
          },
        },
        {
          title: '',
          dataIndex: 'id',
          key: 'action',
          align: 'center',
          render: (val) => (
            <Button minimal intent={Intent.PRIMARY} icon="cog" onClick={() => history.push(`/blueprints/${val}`)} />
          ),
        },
      ] as ColumnType<BlueprintType>,
    [],
  );

  if (loading) {
    return <PageLoading />;
  }

  return (
    <PageHeader breadcrumbs={[{ name: 'Blueprints', path: '/blueprints' }]}>
      <S.Wrapper>
        <div className="action">
          <ButtonGroup>
            <Button
              intent={type === 'all' ? Intent.PRIMARY : Intent.NONE}
              text="All"
              onClick={() => onChangeType('all')}
            />
            {options.map(({ label, value }) => (
              <Button
                key={value}
                intent={type === value ? Intent.PRIMARY : Intent.NONE}
                text={label}
                onClick={() => onChangeType(value)}
              />
            ))}
          </ButtonGroup>
          <Button intent={Intent.PRIMARY} text="Create Blueprint" onClick={() => history.push('/blueprints/create')} />
        </div>
        <Table columns={columns} dataSource={dataSource} />
      </S.Wrapper>
    </PageHeader>
  );
};
