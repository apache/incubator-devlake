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

import { useState } from 'react';
import { CodeOutlined, FileZipOutlined, RightOutlined } from '@ant-design/icons';
import { Table, Space, Modal } from 'antd';
import { pick } from 'lodash';
import { saveAs } from 'file-saver';

import API from '@/api';
import { DEVLAKE_ENDPOINT } from '@/config';
import { IconButton, Inspector } from '@/components';
import { IPipeline } from '@/types';
import { formatTime } from '@/utils';

import { PipelineStatus } from './status';
import { PipelineDuration } from './duration';
import { PipelineTasks } from './tasks';

interface Props {
  loading: boolean;
  dataSource: IPipeline[];
  pagination?: {
    total: number;
    current: number;
    pageSize: number;
    onChange: (page: number) => void;
  };
}

export const PipelineTable = ({ loading, dataSource, pagination }: Props) => {
  const [JSON, setJSON] = useState<any>(null);
  const [id, setId] = useState<ID | null>(null);

  const handleShowJSON = (row: IPipeline) => {
    setJSON(pick(row, ['id', 'name', 'plan', 'skipOnFail']));
  };

  const handleDownloadLog = async (id: ID) => {
    const res = await API.pipeline.log(id);
    if (res) {
      saveAs(`${DEVLAKE_ENDPOINT}/pipelines/${id}/logging.tar.gz`, 'logging.tar.gz');
    }
  };

  const handleShowDetails = (id: ID) => {
    setId(id);
  };

  return (
    <>
      <Table
        rowKey="id"
        size="middle"
        loading={loading}
        columns={[
          {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            align: 'center',
          },
          {
            title: 'Blueprint Name',
            dataIndex: 'name',
            key: 'name',
            align: 'center',
          },
          {
            title: 'Status',
            dataIndex: 'status',
            key: 'status',
            align: 'center',
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
            key: 'duration',
            align: 'center',
            render: (_, { status, beganAt, finishedAt }) => (
              <PipelineDuration status={status} beganAt={beganAt} finishedAt={finishedAt} />
            ),
          },
          {
            title: '',
            dataIndex: 'id',
            key: 'action',
            align: 'center',
            render: (id: ID, row) => (
              <Space>
                <IconButton icon={<CodeOutlined />} helptip="Configuration" onClick={() => handleShowJSON(row)} />
                <IconButton icon={<FileZipOutlined />} helptip="Download Logs" onClick={() => handleDownloadLog(id)} />
                <IconButton icon={<RightOutlined />} helptip="Detail" onClick={() => handleShowDetails(id)} />
              </Space>
            ),
          },
        ]}
        dataSource={dataSource}
        pagination={pagination}
      />
      {JSON && <Inspector open title={`Pipeline ${JSON?.id}`} data={JSON} onClose={() => setJSON(null)} />}
      {id && (
        <Modal open width={820} centered title={`Pipeline ${id}`} footer={null} onCancel={() => setId(null)}>
          <PipelineTasks id={id} />
        </Modal>
      )}
    </>
  );
};
