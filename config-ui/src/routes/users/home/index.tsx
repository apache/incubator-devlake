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

import { useState, useMemo, useRef, useEffect } from 'react';
import { PlusOutlined, EditOutlined, DeleteOutlined, UsergroupAddOutlined } from '@ant-design/icons';
import { Flex, Table, Button, Modal, Input, Space, Popconfirm, Badge, Tooltip, Tag, Card, Descriptions } from 'antd';

import API from '@/api';
import { PageHeader, Block } from '@/components';
import { PATHS } from '@/config';
import { useAutoRefresh, useRefreshData } from '@/hooks';
import { PipelineDuration } from '@/routes/pipeline/components/duration';
import { PipelineStatus } from '@/routes/pipeline/components/status';
import { IPipelineStatus } from '@/types';
import { formatTime, operator } from '@/utils';
import type { IUser, ITeam, IPipeline } from '@/types';

const emptyForm = { name: '', email: '' };
const TEAM_OPTIONS_PAGE_SIZE = 10000;
const CONNECT_USER_ACCOUNTS_PIPELINE_NAME = 'Connect user accounts';

const PIPELINE_ACTIVE_STATUSES = new Set<IPipelineStatus>([
  IPipelineStatus.CREATED,
  IPipelineStatus.PENDING,
  IPipelineStatus.ACTIVE,
  IPipelineStatus.RUNNING,
  IPipelineStatus.RERUN,
]);

const PIPELINE_TERMINAL_STATUSES = new Set<IPipelineStatus>([
  IPipelineStatus.COMPLETED,
  IPipelineStatus.PARTIAL,
  IPipelineStatus.FAILED,
  IPipelineStatus.CANCELLED,
]);

const UserAccountSourcesRow = ({ userId, accountSources }: { userId: string; accountSources?: string[] }) => {
  const sources = accountSources ?? [];

  if (sources.length === 0) {
    return <span>No linked accounts</span>;
  }

  return (
    <Space wrap size={[8, 8]}>
      {sources.map((source) => (
        <Tag key={`${userId}-${source}`} color="blue">
          {source}
        </Tag>
      ))}
    </Space>
  );
};

const renderUserAccountSourcesRow = (record: IUser) => (
  <UserAccountSourcesRow userId={record.id} accountSources={record.accountSources} />
);

export const UsersHomePage = () => {
  const [version, setVersion] = useState(1);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [expandedRowKeys, setExpandedRowKeys] = useState<string[]>([]);
  const [open, setOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<IUser | null>(null);
  const [form, setForm] = useState(emptyForm);
  const [saving, setSaving] = useState(false);
  const [inputValue, setInputValue] = useState('');
  const [searchKeyword, setSearchKeyword] = useState('');
  const [teamsModalOpen, setTeamsModalOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<IUser | null>(null);
  const [selectedTeamIds, setSelectedTeamIds] = useState<string[]>([]);
  const [savingUserTeams, setSavingUserTeams] = useState(false);
  const [teamsSearchKeyword, setTeamsSearchKeyword] = useState('');
  const [triggeringPipeline, setTriggeringPipeline] = useState(false);
  const [triggeredPipelineId, setTriggeredPipelineId] = useState<ID | null>(null);

  const debounceRef = useRef<NodeJS.Timeout | null>(null);
  const refreshedPipelineIdRef = useRef<ID | null>(null);

  const { ready, data } = useRefreshData(
    () => API.user.list({ page, pageSize, ...(searchKeyword.trim() && { email: searchKeyword.trim() }) }),
    [version, page, pageSize, searchKeyword],
  );

  const { ready: teamsReady, data: teamsData } = useRefreshData(
    () =>
      teamsModalOpen
        ? API.team.list({ page: 1, pageSize: TEAM_OPTIONS_PAGE_SIZE, grouped: false })
        : Promise.resolve({ count: 0, teams: [] }),
    [teamsModalOpen],
  );

  const { ready: selectedUserTeamsReady, data: selectedUserTeamsData } = useRefreshData(() => {
    if (!teamsModalOpen || !selectedUser?.id) {
      return Promise.resolve({ userId: '', teamIds: [], teamNames: [], count: 0 });
    }
    return API.user.listTeams(selectedUser.id);
  }, [teamsModalOpen, selectedUser?.id]);

  const { loading: pipelineLoading, data: pipelineData } = useAutoRefresh<IPipeline | null>(
    () => {
      if (!triggeredPipelineId) {
        return Promise.resolve(null);
      }
      return API.pipeline.get(triggeredPipelineId);
    },
    [triggeredPipelineId],
    {
      cancel: (pipeline) => {
        if (!pipeline) {
          return true;
        }
        return PIPELINE_TERMINAL_STATUSES.has(pipeline.status);
      },
    },
  );

  const [dataSource, total] = useMemo(() => [data?.users ?? [], data?.count ?? 0], [data]);

  const teamsDataSource = useMemo<ITeam[]>(() => {
    const keyword = teamsSearchKeyword.trim().toLowerCase();
    const teams = teamsData?.teams ?? [];

    if (!keyword) {
      return teams;
    }

    return teams.filter((team) => {
      const name = team.name?.toLowerCase() ?? '';
      const alias = team.alias?.toLowerCase() ?? '';
      const id = team.id?.toLowerCase() ?? '';
      return name.includes(keyword) || alias.includes(keyword) || id.includes(keyword);
    });
  }, [teamsSearchKeyword, teamsData]);

  useEffect(() => {
    if (!teamsModalOpen) {
      return;
    }
    setSelectedTeamIds(selectedUserTeamsData?.teamIds ?? []);
  }, [teamsModalOpen, selectedUserTeamsData]);

  useEffect(() => {
    if (!pipelineData || !PIPELINE_TERMINAL_STATUSES.has(pipelineData.status)) {
      return;
    }
    if (refreshedPipelineIdRef.current === pipelineData.id) {
      return;
    }
    refreshedPipelineIdRef.current = pipelineData.id;
    setExpandedRowKeys([]);
    setVersion((v) => v + 1);
  }, [pipelineData]);

  const refresh = () => {
    setExpandedRowKeys([]);
    setVersion((v) => v + 1);
  };

  const handleShowCreate = () => {
    setEditingUser(null);
    setForm(emptyForm);
    setOpen(true);
  };

  const handleShowEdit = (user: IUser) => {
    setEditingUser(user);
    setForm({ name: user.name, email: user.email });
    setOpen(true);
  };

  const handleHideDialog = () => {
    setOpen(false);
    setEditingUser(null);
    setForm(emptyForm);
  };

  const handleShowTeamsModal = (user: IUser) => {
    setSelectedUser(user);
    setTeamsSearchKeyword('');
    setTeamsModalOpen(true);
  };

  const handleHideTeamsModal = () => {
    setTeamsModalOpen(false);
    setSelectedUser(null);
    setSelectedTeamIds([]);
    setTeamsSearchKeyword('');
  };

  const handleSave = async () => {
    const [success] = await operator(
      async () => {
        if (editingUser) {
          return API.user.update(editingUser.id, { ...form, teamIds: editingUser.teamIds ?? '' });
        }
        return API.user.create({ users: [{ ...form, teamIds: '' }] });
      },
      { setOperating: setSaving },
    );

    if (success) {
      handleHideDialog();
      refresh();
    }
  };

  const handleDelete = async (userId: string) => {
    const [success] = await operator(() => API.user.remove(userId));
    if (success) {
      refresh();
    }
  };

  const handleSaveUserTeams = async () => {
    if (!selectedUser) {
      return;
    }

    const [success] = await operator(() => API.user.updateTeams(selectedUser.id, { teamIds: selectedTeamIds }), {
      setOperating: setSavingUserTeams,
    });

    if (success) {
      handleHideTeamsModal();
      refresh();
    }
  };

  const handleConnectUserAccounts = async () => {
    const [success, pipeline] = await operator(
      () =>
        API.pipeline.create({
          name: CONNECT_USER_ACCOUNTS_PIPELINE_NAME,
          plan: [[{ plugin: 'org', subtasks: ['connectUserAccountsExact'] }]],
        }),
      {
        setOperating: setTriggeringPipeline,
        formatMessage: () => 'User account linking pipeline triggered.',
        formatReason: () => 'Failed to trigger user account linking pipeline.',
      },
    );

    if (success && pipeline?.id) {
      refreshedPipelineIdRef.current = null;
      setTriggeredPipelineId(pipeline.id);
    }
  };

  const handleSearch = (value: string) => {
    setInputValue(value);

    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
    }

    debounceRef.current = setTimeout(() => {
      setSearchKeyword(value.trim());
      setPage(1);
      refresh();
    }, 500);
  };

  const isPipelineActive =
    triggeredPipelineId !== null && (!pipelineData || PIPELINE_ACTIVE_STATUSES.has(pipelineData.status));

  const showPipelineInfo = triggeredPipelineId !== null;

  return (
    <PageHeader breadcrumbs={[{ name: 'Users', path: PATHS.USERS() }]}>
      <Flex style={{ marginBottom: 16, width: '100%' }} justify="space-between" align="center">
        <Space>
          <Button
            type="primary"
            onClick={handleConnectUserAccounts}
            loading={triggeringPipeline}
            disabled={isPipelineActive}
          >
            Connect User Accounts
          </Button>
          {showPipelineInfo && (
            <Button type="link" href={PATHS.PIPELINES()}>
              View Pipelines
            </Button>
          )}
        </Space>
        <Space>
          <Input
            placeholder="Search by email ..."
            style={{ width: 300 }}
            value={inputValue}
            onChange={(e) => handleSearch(e.target.value)}
          />
          <Button type="primary" icon={<PlusOutlined />} onClick={handleShowCreate}>
            New User
          </Button>
        </Space>
      </Flex>
      {showPipelineInfo && (
        <Card size="small" style={{ marginBottom: 16 }} loading={pipelineLoading && !pipelineData}>
          <Descriptions column={3} size="small">
            <Descriptions.Item label="Pipeline ID">{triggeredPipelineId}</Descriptions.Item>
            <Descriptions.Item label="Status">
              {pipelineData ? <PipelineStatus status={pipelineData.status} /> : 'Loading...'}
            </Descriptions.Item>
            <Descriptions.Item label="Tasks">
              {pipelineData ? `${pipelineData.finishedTasks}/${pipelineData.totalTasks}` : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="Started At">{formatTime(pipelineData?.beganAt ?? null)}</Descriptions.Item>
            <Descriptions.Item label="Finished At">{formatTime(pipelineData?.finishedAt ?? null)}</Descriptions.Item>
            <Descriptions.Item label="Duration">
              {pipelineData ? (
                <PipelineDuration
                  status={pipelineData.status}
                  beganAt={pipelineData.beganAt}
                  finishedAt={pipelineData.finishedAt}
                />
              ) : (
                '-'
              )}
            </Descriptions.Item>
            <Descriptions.Item label="Current Stage">{pipelineData?.stage ?? '-'}</Descriptions.Item>
            <Descriptions.Item label="Message" span={2}>
              {pipelineData?.message || '-'}
            </Descriptions.Item>
          </Descriptions>
        </Card>
      )}
      <Table
        rowKey="id"
        size="middle"
        loading={!ready}
        expandable={{
          expandedRowKeys,
          onExpand: (expanded, record: IUser) => setExpandedRowKeys(expanded ? [record.id] : []),
          rowExpandable: () => true,
          expandedRowRender: renderUserAccountSourcesRow,
        }}
        columns={[
          {
            title: 'Name',
            dataIndex: 'name',
            key: 'name',
          },
          {
            title: 'Email',
            dataIndex: 'email',
            key: 'email',
            render: (val: string) => val || '-',
          },
          {
            title: 'Accounts',
            dataIndex: 'accountCount',
            key: 'accountCount',
            width: 120,
            align: 'center',
            render: (val: number | undefined, record: IUser) => val ?? record.accountSources?.length ?? 0,
          },
          {
            title: 'Teams',
            dataIndex: 'teamCount',
            key: 'teamCount',
            width: 130,
            align: 'center',
            render: (val: number | undefined, record: IUser) => {
              const count = val ?? record.teamNames?.length ?? 0;
              const teams = record.teamNames ?? [];
              const tooltipTitle =
                teams.length > 0 ? (
                  <div>
                    {teams.map((teamName, index) => (
                      <div key={`${record.id}-${teamName}-${index}`}>{teamName}</div>
                    ))}
                  </div>
                ) : (
                  'No groups assigned'
                );

              return (
                <Tooltip title={tooltipTitle} placement="left">
                  <Badge count={count} showZero>
                    <Button
                      type="primary"
                      size="small"
                      icon={<UsergroupAddOutlined />}
                      onClick={() => handleShowTeamsModal(record)}
                    />
                  </Badge>
                </Tooltip>
              );
            },
          },
          {
            title: '',
            key: 'actions',
            width: 120,
            align: 'center',
            render: (_, record: IUser) => (
              <Space>
                <Button type="link" icon={<EditOutlined />} onClick={() => handleShowEdit(record)} />
                <Popconfirm
                  title="Are you sure you want to delete this user?"
                  onConfirm={() => handleDelete(record.id)}
                  okText="Yes"
                  cancelText="No"
                >
                  <Button type="link" danger icon={<DeleteOutlined />} />
                </Popconfirm>
              </Space>
            ),
          },
        ]}
        dataSource={dataSource}
        pagination={{
          current: page,
          pageSize,
          total,
          onChange: ((newPage: number, newPageSize: number) => {
            setPage(newPage);
            if (newPageSize !== pageSize) {
              setPageSize(newPageSize);
            }
          }) as (newPage: number) => void,
        }}
      />
      <Modal
        open={teamsModalOpen}
        width={820}
        centered
        title={selectedUser ? `Associate Teams to ${selectedUser.name}` : 'Associate Teams'}
        okText="Save"
        okButtonProps={{
          loading: savingUserTeams,
          disabled: !selectedUser || !teamsReady || !selectedUserTeamsReady,
        }}
        onOk={handleSaveUserTeams}
        onCancel={handleHideTeamsModal}
      >
        <Flex style={{ marginBottom: 16, width: '100%' }} justify="space-between" align="center">
          <Input
            placeholder="Search teams by name, alias or id ..."
            style={{ width: 360 }}
            value={teamsSearchKeyword}
            onChange={(e) => setTeamsSearchKeyword(e.target.value)}
          />
          <span>Selected teams: {selectedTeamIds.length}</span>
        </Flex>
        <Table
          rowKey="id"
          size="small"
          loading={!teamsReady || !selectedUserTeamsReady}
          columns={[
            {
              title: 'Name',
              dataIndex: 'name',
              key: 'name',
            },
            {
              title: 'Alias',
              dataIndex: 'alias',
              key: 'alias',
              render: (val: string) => val || '-',
            },
          ]}
          dataSource={teamsDataSource}
          rowSelection={{
            selectedRowKeys: selectedTeamIds,
            onChange: (selectedRowKeys) => setSelectedTeamIds(selectedRowKeys as string[]),
          }}
          pagination={{
            pageSize: 8,
            showSizeChanger: false,
          }}
        />
      </Modal>
      <Modal
        open={open}
        width={820}
        centered
        title={editingUser ? 'Edit User' : 'Create a New User'}
        okText="Save"
        okButtonProps={{
          disabled: !form.name,
          loading: saving,
        }}
        onOk={handleSave}
        onCancel={handleHideDialog}
      >
        <Block title="Name" description="Give the user a name" required>
          <Input
            style={{ width: 386 }}
            placeholder="Name"
            value={form.name}
            onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
          />
        </Block>
        <Block title="Email" description="The user's email address">
          <Input
            style={{ width: 386 }}
            placeholder="Email"
            value={form.email}
            onChange={(e) => setForm((f) => ({ ...f, email: e.target.value }))}
          />
        </Block>
      </Modal>
    </PageHeader>
  );
};
