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

import { Link } from 'react-router-dom';
import { Button, Intent } from '@blueprintjs/core';

import { Table } from '@/components';
import { ConnectionStatus, useConnection } from '@/store';

import { WebHookConnection } from '@/plugins/register/webook';

interface Props {
  plugin: string;
  onCreate: () => void;
}

export const ConnectionList = ({ plugin, onCreate }: Props) => {
  if (plugin === 'webhook') {
    return <WebHookConnection />;
  }

  return <BaseList plugin={plugin} onCreate={onCreate} />;
};

const BaseList = ({ plugin, onCreate }: Props) => {
  const { connections, onTest } = useConnection();

  return (
    <>
      <Table
        noShadow
        columns={[
          {
            title: 'Connection Name',
            dataIndex: 'name',
            key: 'name',
          },
          {
            title: 'Status',
            dataIndex: ['status', 'unique'],
            key: 'status',
            render: ({ status, unique }) => <ConnectionStatus status={status} unique={unique} onTest={onTest} />,
          },
          {
            title: '',
            dataIndex: ['plugin', 'id'],
            key: 'link',
            width: 100,
            render: ({ plugin, id }) => <Link to={`/connections/${plugin}/${id}`}>Details</Link>,
          },
        ]}
        dataSource={connections.filter((cs) => cs.plugin === plugin)}
        noData={{
          text: 'There is no data connection yet. Please add a new connection.',
        }}
      />
      <Button
        style={{ marginTop: 16 }}
        intent={Intent.PRIMARY}
        icon="add"
        text="Create a New Connection"
        onClick={onCreate}
      />
    </>
  );
};
