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
import { useNavigate } from 'react-router-dom';
import { EyeOutlined, EditOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import { theme, Table, Button, Modal, message } from 'antd';

import { selectConnections, removeConnection } from '@/features/connections';
import { Message } from '@/components';
import { useAppDispatch, useAppSelector } from '@/hooks';
import { ConnectionStatus, ConnectionFormModal } from '@/plugins';
import { WebHookConnection } from '@/plugins/register/webhook';
import { operator } from '@/utils';

interface Props {
  plugin: string;
  onCreate: () => void;
}

export const ConnectionList = ({ plugin, onCreate }: Props) => {
  const [modalType, setModalType] = useState<'update' | 'delete' | 'deleteFailed'>();
  const [connectionId, setConnectionId] = useState<ID>();
  const [operating, setOperating] = useState(false);
  const [conflict, setConflict] = useState<string[]>([]);
  const [errorMsg, setErrorMsg] = useState('');

  const {
    token: { colorPrimary },
  } = theme.useToken();

  const dispatch = useAppDispatch();
  const connections = useAppSelector((state) => selectConnections(state, plugin));

  const navigate = useNavigate();

  const handleShowModal = (type: 'update' | 'delete', id: ID) => {
    setModalType(type);
    setConnectionId(id);
  };

  const hanldeHideModal = () => {
    setModalType(undefined);
    setConnectionId(undefined);
  };

  const handleDelete = async () => {
    const [, res] = await operator(
      async () => {
        try {
          await dispatch(removeConnection({ plugin, connectionId })).unwrap();
          return { status: 'success' };
        } catch (err: any) {
          const { status, data, message } = err;
          return {
            status: status === 409 ? 'conflict' : 'error',
            conflict: data ? [...data.projects, ...data.blueprints] : [],
            message,
          };
        }
      },
      {
        setOperating,
        hideToast: true,
      },
    );

    if (res.status === 'success') {
      message.success('Delete Connection Successful.');
      hanldeHideModal();
    } else if (res.status === 'conflict') {
      setModalType('deleteFailed');
      setConflict(res.conflict);
      setErrorMsg(res.message);
    } else {
      message.error('Operation failed.');
      hanldeHideModal();
    }
  };

  if (plugin === 'webhook') {
    return <WebHookConnection />;
  }

  return (
    <>
      <Table
        rowKey="id"
        size="small"
        columns={[
          {
            title: 'Connection Name',
            dataIndex: 'name',
            key: 'name',
          },
          {
            title: 'Status',
            key: 'status',
            width: 200,
            render: (_, row) => <ConnectionStatus connection={row} />,
          },
          {
            title: '',
            key: 'link',
            width: 300,
            render: (_, { plugin, id }) => (
              <>
                <Button type="link" icon={<EyeOutlined />} onClick={() => navigate(`/connections/${plugin}/${id}`)}>
                  Details
                </Button>
                <Button type="link" icon={<EditOutlined />} onClick={() => handleShowModal('update', id)}>
                  Edit
                </Button>
                <Button type="link" danger icon={<DeleteOutlined />} onClick={() => handleShowModal('delete', id)}>
                  Delete
                </Button>
              </>
            ),
          },
        ]}
        dataSource={connections}
        pagination={false}
      />
      <Button style={{ marginTop: 16 }} type="primary" icon={<PlusOutlined />} onClick={onCreate}>
        Create a New Connection
      </Button>
      {modalType === 'update' && (
        <ConnectionFormModal plugin={plugin} connectionId={connectionId} open onCancel={hanldeHideModal} />
      )}
      {modalType === 'delete' && (
        <Modal
          open
          width={820}
          centered
          title="Would you like to delete this Data Connection?"
          okText="Confirm"
          okButtonProps={{
            loading: operating,
          }}
          onCancel={hanldeHideModal}
          onOk={handleDelete}
        >
          <Message
            content=" This operation cannot be undone. Deleting a Data Connection will delete all data that have been collected
              in this Connection."
          />
        </Modal>
      )}
      {modalType === 'deleteFailed' && (
        <Modal
          open
          width={820}
          centered
          style={{ width: 820 }}
          title="This Data Connection can not be deleted."
          cancelButtonProps={{
            style: {
              display: 'none',
            },
          }}
          onCancel={hanldeHideModal}
          onOk={hanldeHideModal}
        >
          {!conflict.length ? (
            <Message content={errorMsg} />
          ) : (
            <>
              <Message
                content={`This Data Connection can not be deleted because it has been used in the following projects/blueprints:`}
              />
              <ul style={{ paddingLeft: 36 }}>
                {conflict.map((it) => (
                  <li key={it} style={{ color: colorPrimary }}>
                    {it}
                  </li>
                ))}
              </ul>
            </>
          )}
        </Modal>
      )}
    </>
  );
};
