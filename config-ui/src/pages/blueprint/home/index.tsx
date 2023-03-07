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
import { Link, useHistory } from 'react-router-dom';
import { ButtonGroup, Button, Tag, Intent } from '@blueprintjs/core';

import { PageLoading, PageHeader, Table, ColumnType, IconButton, TextTooltip } from '@/components';
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
          ellipsis: true,
        },
        {
          title: 'Data Connections',
          key: 'connections',
          align: 'center',
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
          width: 100,
          align: 'center',
          render: (_, row) => {
            const cron = getCron(row.isManual, row.cronConfig);
            return cron.label;
          },
        },
        {
          title: 'Next Run Time',
          key: 'nextRunTime',
          width: 200,
          align: 'center',
          render: (_, row) => {
            const cron = getCron(row.isManual, row.cronConfig);
            return formatTime(cron.nextTime);
          },
        },
        {
          title: 'Project',
          dataIndex: 'projectName',
          key: 'project',
          align: 'center',
          render: (val) =>
            val ? (
              <Link to={`/projects/${window.encodeURIComponent(val)}`}>
                <TextTooltip content={val}>{val}</TextTooltip>
              </Link>
            ) : (
              val
            ),
        },
        {
          title: 'Status',
          dataIndex: 'enable',
          key: 'enable',
          align: 'center',
          width: 100,
          render: (val) => (
            <Tag minimal intent={val ? Intent.SUCCESS : Intent.DANGER}>
              {val ? 'Enabled' : 'Disabled'}
            </Tag>
          ),
        },
        {
          title: '',
          dataIndex: 'id',
          key: 'action',
          width: 100,
          align: 'center',
          render: (val) => (
            <IconButton icon="cog" tooltip="Detail" onClick={() => history.push(`/blueprints/${val}`)} />
          ),
        },
      ] as ColumnType<BlueprintType>,
    [],
  );

  const handleCreate = () => history.push('/blueprints/create');

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
          <Button icon="plus" intent={Intent.PRIMARY} text="New Blueprint" onClick={handleCreate} />
        </div>
        <Table
          columns={columns}
          dataSource={dataSource}
          noData={{
            text: 'There is no Blueprint yet. Please add a new Blueprint here or from a Project.',
            btnText: 'New Blueprint',
            onCreate: handleCreate,
          }}
        />
      </S.Wrapper>
    </PageHeader>
  );
};
