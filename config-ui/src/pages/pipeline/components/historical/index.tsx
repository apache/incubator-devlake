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
import { ButtonGroup } from '@blueprintjs/core';
import { pick } from 'lodash';
import { saveAs } from 'file-saver';

import { DEVLAKE_ENDPOINT } from '@/config';
import type { ColumnType } from '@/components';
import { Card, Loading, Table, Inspector, Dialog, IconButton } from '@/components';
import { useAutoRefresh } from '@/hooks';
import { formatTime } from '@/utils';

import type { PipelineType } from '../../types';
import { StatusEnum } from '../../types';
import * as API from '../../api';

import { usePipeline } from '../context';
import { PipelineStatus } from '../status';
import { PipelineDuration } from '../duration';
import { PipelineTasks } from '../tasks';

interface Props {
  blueprintId: ID;
}

export const PipelineHistorical = ({ blueprintId }: Props) => {
  const [JSON, setJSON] = useState<any>(null);
  const [ID, setID] = useState<ID | null>(null);

  const { version } = usePipeline();

  const { loading, data } = useAutoRefresh<PipelineType[]>(
    async () => {
      const res = await API.getPipelineHistorical(blueprintId);
      return res.pipelines;
    },
    [version],
    {
      cancel: (data) =>
        !!(
          data &&
          data.every((it) =>
            [StatusEnum.COMPLETED, StatusEnum.PARTIAL, StatusEnum.CANCELLED, StatusEnum.FAILED].includes(it.status),
          )
        ),
    },
  );

  const handleShowJSON = (row: PipelineType) => {
    setJSON(pick(row, ['id', 'name', 'plan', 'skipOnFail']));
  };

  const handleDownloadLog = async (id: ID) => {
    const res = await API.getPipelineLog(id);
    if (res) {
      saveAs(`${DEVLAKE_ENDPOINT}/pipelines/${id}/logging.tar.gz`, 'logging.tar.gz');
    }
  };

  const handleShowDetails = (id: ID) => {
    setID(id);
  };

  const columns = useMemo(
    () =>
      [
        {
          title: 'Status',
          dataIndex: 'status',
          key: 'status',
          render: (val) => <PipelineStatus status={val} />,
        },
        {
          title: 'Started at',
          dataIndex: 'beganAt',
          key: 'beganAt',
          align: 'center',
          render: (val: string | null) => (val ? formatTime(val) : '-'),
        },
        {
          title: 'Completed at',
          dataIndex: 'finishedAt',
          key: 'finishedAt',
          align: 'center',
          render: (val: string | null) => (val ? formatTime(val) : '-'),
        },
        {
          title: 'Duration',
          dataIndex: ['status', 'beganAt', 'finishedAt'],
          key: 'duration',
          align: 'center',
          render: ({ status, beganAt, finishedAt }) => (
            <PipelineDuration status={status} beganAt={beganAt} finishedAt={finishedAt} />
          ),
        },
        {
          title: '',
          dataIndex: 'id',
          key: 'action',
          align: 'center',
          render: (id: ID, row) => (
            <ButtonGroup>
              <IconButton icon="code" tooltip="View JSON" onClick={() => handleShowJSON(row)} />
              <IconButton icon="document" tooltip="Download Logs" onClick={() => handleDownloadLog(id)} />
              <IconButton icon="chevron-right" tooltip="View Details" onClick={() => handleShowDetails(id)} />
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
      {JSON && <Inspector isOpen title={`Pipeline ${JSON?.id}`} data={JSON} onClose={() => setJSON(null)} />}
      {ID && (
        <Dialog style={{ width: 720 }} isOpen title={`Pipeline ${ID}`} footer={null} onCancel={() => setID(null)}>
          <PipelineTasks id={ID} />
        </Dialog>
      )}
    </div>
  );
};
