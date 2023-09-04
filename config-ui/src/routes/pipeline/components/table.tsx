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

import { useState, useMemo } from 'react';
import { ButtonGroup } from '@blueprintjs/core';
import { pick } from 'lodash';
import { saveAs } from 'file-saver';

import { DEVLAKE_ENDPOINT } from '@/config';
import { Table, ColumnType, IconButton, Inspector, Dialog } from '@/components';
import { formatTime } from '@/utils';

import * as T from '../types';
import * as API from '../api';

import { PipelineStatus } from './status';
import { PipelineDuration } from './duration';
import { PipelineTasks } from './tasks';

interface Props {
  loading: boolean;
  dataSource: T.Pipeline[];
  pagination?: {
    total: number;
    page: number;
    pageSize: number;
    onChange: (page: number) => void;
  };
  noData?: {
    text?: React.ReactNode;
    btnText?: string;
    onCreate?: () => void;
  };
}

export const PipelineTable = ({ dataSource, pagination, noData }: Props) => {
  const [JSON, setJSON] = useState<any>(null);
  const [id, setId] = useState<ID | null>(null);

  const handleShowJSON = (row: T.Pipeline) => {
    setJSON(pick(row, ['id', 'name', 'plan', 'skipOnFail']));
  };

  const handleDownloadLog = async (id: ID) => {
    const res = await API.getPipelineLog(id);
    if (res) {
      saveAs(`${DEVLAKE_ENDPOINT}/pipelines/${id}/logging.tar.gz`, 'logging.tar.gz');
    }
  };

  const handleShowDetails = (id: ID) => {
    setId(id);
  };

  const columns = useMemo(
    () =>
      [
        {
          title: 'ID',
          dataIndex: 'id',
          key: 'id',
        },
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
          render: (val) => formatTime(val),
        },
        {
          title: 'Completed at',
          dataIndex: 'finishedAt',
          key: 'finishedAt',
          align: 'center',
          render: (val) => formatTime(val),
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
          align: 'center',
          render: (id: ID, row) => (
            <ButtonGroup>
              <IconButton icon="code" tooltip="View JSON" onClick={() => handleShowJSON(row)} />
              <IconButton icon="document" tooltip="Download Logs" onClick={() => handleDownloadLog(id)} />
              <IconButton icon="chevron-right" tooltip="View Details" onClick={() => handleShowDetails(id)} />
            </ButtonGroup>
          ),
        },
      ] as ColumnType<T.Pipeline>,
    [],
  );

  return (
    <>
      <Table columns={columns} dataSource={dataSource} pagination={pagination} noData={noData} />
      {JSON && <Inspector isOpen title={`Pipeline ${JSON?.id}`} data={JSON} onClose={() => setJSON(null)} />}
      {id && (
        <Dialog style={{ width: 820 }} isOpen title={`Pipeline ${id}`} footer={null} onCancel={() => setId(null)}>
          <PipelineTasks id={id} />
        </Dialog>
      )}
    </>
  );
};
