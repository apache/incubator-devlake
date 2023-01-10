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

import React, { useState, useMemo } from 'react';
import { Icon, ButtonGroup, Button, Position, Intent, IconName } from '@blueprintjs/core';
import { Tooltip2 } from '@blueprintjs/popover2';
import { pick } from 'lodash';
import { saveAs } from 'file-saver';

import { DEVLAKE_ENDPOINT } from '@/config';
import { Card, Loading, Table, ColumnType, Inspector } from '@/components';
import { useAutoRefresh } from '@/hooks';
import { formatTime } from '@/utils';

import type { PipelineType } from '../../types';
import { StatusEnum } from '../../types';
import { STATUS_ICON, STATUS_LABEL, STATUS_CLS } from '../../misc';
import * as API from '../../api';

import { PipelineDuration } from '../duration';

import * as S from './styled';

interface Props {
  blueprintId: ID;
}

export const PipelineHistorical = ({ blueprintId }: Props) => {
  const [isOpen, setIsOpen] = useState(false);
  const [json, setJson] = useState<any>({});

  const { loading, data } = useAutoRefresh<PipelineType[]>(
    async () => {
      const res = await API.getPipelineHistorical(blueprintId);
      return res.pipelines;
    },
    [],
    {
      cancel: (data) =>
        !!(
          data &&
          data.every((it) => [StatusEnum.COMPLETED, StatusEnum.CANCELLED, StatusEnum.FAILED].includes(it.status))
        ),
    },
  );

  const handleDownloadLog = async (id: ID) => {
    const res = await API.getPipelineLog(id);
    if (res) {
      saveAs(`${DEVLAKE_ENDPOINT}/pipelines/${id}/logging.tar.gz`, 'logging.tar.gz');
    }
  };

  const columns = useMemo(
    () =>
      [
        {
          title: 'Status',
          dataIndex: 'status',
          key: 'status',
          render: (val: StatusEnum) => (
            <S.StatusColumn className={STATUS_CLS(val)}>
              {STATUS_ICON[val] === 'loading' ? (
                <Loading style={{ marginRight: 4 }} size={14} />
              ) : (
                <Icon style={{ marginRight: 4 }} icon={STATUS_ICON[val] as IconName} />
              )}
              <span>{STATUS_LABEL[val]}</span>
            </S.StatusColumn>
          ),
        },
        {
          title: 'Started at',
          dataIndex: 'beganAt',
          key: 'beganAt',
          render: (val: string) => (val ? formatTime(val) : '-'),
        },
        {
          title: 'Completed at',
          dataIndex: 'finishedAt',
          key: 'finishedAt',
          render: (val: string) => (val ? formatTime(val) : '-'),
        },
        {
          title: 'Duration',
          dataIndex: ['status', 'beganAt', 'finishedAt'],
          key: 'duration',
          render: ({ status, beganAt, finishedAt }) => (
            <PipelineDuration status={status} beganAt={beganAt} finishedAt={finishedAt} />
          ),
        },
        {
          title: '',
          dataIndex: 'id',
          key: 'action',
          render: (id: ID, row) => (
            <ButtonGroup>
              <Tooltip2 position={Position.TOP} intent={Intent.PRIMARY} content="View JSON">
                <Button
                  minimal
                  intent={Intent.PRIMARY}
                  icon="code"
                  onClick={() => {
                    setIsOpen(true);
                    setJson(pick(row, ['id', 'name', 'plan', 'skipOnFail']));
                  }}
                />
              </Tooltip2>
              <Tooltip2 position={Position.TOP} intent={Intent.PRIMARY} content="Download Logs">
                <Button minimal intent={Intent.PRIMARY} icon="document" onClick={() => handleDownloadLog(id)} />
              </Tooltip2>
              {/* <Button minimal intent={Intent.PRIMARY} icon='chevron-right' /> */}
            </ButtonGroup>
          ),
        },
      ] as ColumnType<PipelineType>,
    [],
  );

  if (loading) {
    return (
      <Card>
        <Loading />
      </Card>
    );
  }

  if (!data) {
    return <Card>There are no historical runs associated with this blueprint.</Card>;
  }

  return (
    <div>
      <Table columns={columns} dataSource={data} />
      <Inspector isOpen={isOpen} title={`Pipeline ${json?.id}`} data={json} onClose={() => setIsOpen(false)} />
    </div>
  );
};
