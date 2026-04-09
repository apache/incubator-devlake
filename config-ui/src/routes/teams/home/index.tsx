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
import { Flex, Table, Button, Modal, Input, InputNumber, Space, Popconfirm, Select, Badge } from 'antd';

import API from '@/api';
import { PageHeader, Block } from '@/components';
import { PATHS } from '@/config';
import { useRefreshData } from '@/hooks';
import { operator } from '@/utils';
import type { ITeam, IUser } from '@/types';

const emptyForm = { name: '', alias: '', parentId: '', sortingIndex: 0 };
const PARENT_TEAM_OPTIONS_PAGE_SIZE = 10000;
const TEAM_USERS_PAGE_SIZE = 10000;

export const TeamsHomePage = () => {
  const [version, setVersion] = useState(1);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [expandedRowKeys, setExpandedRowKeys] = useState<Array<string | number>>([]);
  const [open, setOpen] = useState(false);
  const [editingTeam, setEditingTeam] = useState<ITeam | null>(null);
  const [form, setForm] = useState(emptyForm);
  const [saving, setSaving] = useState(false);
  const [inputValue, setInputValue] = useState('');
  const [searchKeyword, setSearchKeyword] = useState('');
  const [usersModalOpen, setUsersModalOpen] = useState(false);
  const [selectedTeam, setSelectedTeam] = useState<ITeam | null>(null);
  const [selectedUserIds, setSelectedUserIds] = useState<string[]>([]);
  const [savingTeamUsers, setSavingTeamUsers] = useState(false);
  const [usersSearchKeyword, setUsersSearchKeyword] = useState('');

  const debounceRef = useRef<NodeJS.Timeout | null>(null);

  const { ready, data } = useRefreshData(
    () => API.team.list({ page, pageSize, grouped: true, ...(searchKeyword.trim() && { name: searchKeyword.trim() }) }),
    [version, page, pageSize, searchKeyword],
  );

  const { ready: parentOptionsReady, data: parentOptionsData } = useRefreshData(
    () => API.team.list({ page: 1, pageSize: PARENT_TEAM_OPTIONS_PAGE_SIZE, grouped: false }),
    [version],
  );

  const { ready: usersReady, data: usersData } = useRefreshData(
    () =>
      usersModalOpen
        ? API.user.list({ page: 1, pageSize: TEAM_USERS_PAGE_SIZE })
        : Promise.resolve({ count: 0, users: [] }),
    [usersModalOpen],
  );

  const { ready: selectedTeamUsersReady, data: selectedTeamUsersData } = useRefreshData(() => {
    if (!usersModalOpen || !selectedTeam?.id) {
      return Promise.resolve({ teamId: '', userIds: [], count: 0 });
    }
    return API.team.listUsers(selectedTeam.id);
  }, [usersModalOpen, selectedTeam?.id]);

  const [dataSource, total] = useMemo(() => [data?.teams ?? [], data?.count ?? 0], [data]);

  const parentTeamOptions = useMemo(
    () =>
      (parentOptionsData?.teams ?? [])
        .filter((team) => team.id !== editingTeam?.id)
        .sort((a, b) => {
          if (a.sortingIndex !== b.sortingIndex) {
            return a.sortingIndex - b.sortingIndex;
          }
          return a.name.localeCompare(b.name);
        })
        .map((team) => ({
          value: team.id,
          label: `${team.name} (${team.id})`,
        })),
    [parentOptionsData, editingTeam],
  );

  const usersDataSource = useMemo<IUser[]>(() => {
    const keyword = usersSearchKeyword.trim().toLowerCase();
    const users = usersData?.users ?? [];

    if (!keyword) {
      return users;
    }

    return users.filter((user) => {
      const name = user.name?.toLowerCase() ?? '';
      const email = user.email?.toLowerCase() ?? '';
      return name.includes(keyword) || email.includes(keyword);
    });
  }, [usersSearchKeyword, usersData]);

  useEffect(() => {
    if (!usersModalOpen) {
      return;
    }
    setSelectedUserIds(selectedTeamUsersData?.userIds ?? []);
  }, [usersModalOpen, selectedTeamUsersData]);

  const refresh = () => {
    setExpandedRowKeys([]);
    setVersion((v) => v + 1);
  };

  const handleShowCreate = () => {
    setEditingTeam(null);
    setForm(emptyForm);
    setOpen(true);
  };

  const handleShowEdit = (team: ITeam) => {
    setEditingTeam(team);
    setForm({ name: team.name, alias: team.alias, parentId: team.parentId, sortingIndex: team.sortingIndex });
    setOpen(true);
  };

  const handleHideDialog = () => {
    setOpen(false);
    setEditingTeam(null);
    setForm(emptyForm);
  };

  const handleShowUsersModal = (team: ITeam) => {
    setSelectedTeam(team);
    setUsersSearchKeyword('');
    setUsersModalOpen(true);
  };

  const handleHideUsersModal = () => {
    setUsersModalOpen(false);
    setSelectedTeam(null);
    setSelectedUserIds([]);
    setUsersSearchKeyword('');
  };

  const handleSave = async () => {
    const [success] = await operator(
      async () => {
        if (editingTeam) {
          return API.team.update(editingTeam.id, form);
        }
        return API.team.create({ teams: [form] });
      },
      { setOperating: setSaving },
    );

    if (success) {
      handleHideDialog();
      refresh();
    }
  };

  const handleDelete = async (teamId: string) => {
    const [success] = await operator(() => API.team.remove(teamId));
    if (success) {
      refresh();
    }
  };

  const handleSaveTeamUsers = async () => {
    if (!selectedTeam) {
      return;
    }

    const [success] = await operator(() => API.team.updateUsers(selectedTeam.id, { userIds: selectedUserIds }), {
      setOperating: setSavingTeamUsers,
    });

    if (success) {
      handleHideUsersModal();
      refresh();
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

  const handleExpand = (expanded: boolean, record: ITeam) => {
    if (!record.children?.length) {
      return;
    }
    setExpandedRowKeys(expanded ? [record.id] : []);
  };

  return (
    <PageHeader breadcrumbs={[{ name: 'Teams', path: PATHS.TEAMS() }]}>
      <Flex style={{ marginBottom: 16, width: '100%' }} justify="flex-end" align="center">
        <Input
          placeholder="Search team ..."
          style={{ width: 300, marginRight: 12 }}
          value={inputValue}
          onChange={(e) => handleSearch(e.target.value)}
        />
        <Button type="primary" icon={<PlusOutlined />} onClick={handleShowCreate}>
          New Team
        </Button>
      </Flex>
      <Table
        rowKey="id"
        size="middle"
        loading={!ready}
        expandable={{
          expandedRowKeys,
          onExpand: handleExpand,
          rowExpandable: (record: ITeam) => !!record.children?.length,
        }}
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
          {
            title: 'Users',
            dataIndex: 'userCount',
            key: 'userCount',
            width: 160,
            align: 'center',
            render: (val: number | undefined, record: ITeam) => (
              <Badge status="success" count={val ?? 0} showZero>
                <Button
                  type="primary"
                  size="small"
                  onClick={() => handleShowUsersModal(record)}
                  icon={<UsergroupAddOutlined />}
                />
              </Badge>
            ),
          },
          {
            title: '',
            key: 'actions',
            width: 120,
            align: 'center',
            render: (_, record: ITeam) => (
              <Space>
                <Button type="link" icon={<EditOutlined />} onClick={() => handleShowEdit(record)} />
                <Popconfirm
                  title="Are you sure you want to delete this team?"
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
        open={usersModalOpen}
        width={820}
        centered
        title={selectedTeam ? `Assign Users to ${selectedTeam.name}` : 'Assign Users'}
        okText="Save"
        okButtonProps={{
          loading: savingTeamUsers,
          disabled: !selectedTeam || !usersReady || !selectedTeamUsersReady,
        }}
        onOk={handleSaveTeamUsers}
        onCancel={handleHideUsersModal}
      >
        <Flex style={{ marginBottom: 16, width: '100%' }} justify="space-between" align="center">
          <Input
            placeholder="Search users by name or email ..."
            style={{ width: 360 }}
            value={usersSearchKeyword}
            onChange={(e) => setUsersSearchKeyword(e.target.value)}
          />
          <span>Selected users: {selectedUserIds.length}</span>
        </Flex>
        <Table
          rowKey="id"
          size="small"
          loading={!usersReady || !selectedTeamUsersReady}
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
          ]}
          dataSource={usersDataSource}
          rowSelection={{
            selectedRowKeys: selectedUserIds,
            onChange: (selectedRowKeys) => setSelectedUserIds(selectedRowKeys as string[]),
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
        title={editingTeam ? 'Edit Team' : 'Create a New Team'}
        okText="Save"
        okButtonProps={{
          disabled: !form.name,
          loading: saving,
        }}
        onOk={handleSave}
        onCancel={handleHideDialog}
      >
        <Block title="Team Name" description="Give your team a unique name" required>
          <Input
            style={{ width: 386 }}
            placeholder="Team Name"
            value={form.name}
            onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
          />
        </Block>
        <Block title="Alias" description="An optional short alias for the team">
          <Input
            style={{ width: 386 }}
            placeholder="Alias"
            value={form.alias}
            onChange={(e) => setForm((f) => ({ ...f, alias: e.target.value }))}
          />
        </Block>
        <Block title="Parent ID" description="The ID of the parent team, if any">
          <Select
            style={{ width: 386 }}
            placeholder="Select Parent Team"
            showSearch
            allowClear
            optionFilterProp="label"
            loading={!parentOptionsReady}
            options={parentTeamOptions}
            value={form.parentId || undefined}
            onChange={(val: string | undefined) => setForm((f) => ({ ...f, parentId: val ?? '' }))}
          />
        </Block>
        <Block title="Sorting Index" description="Numeric index used for ordering teams">
          <InputNumber
            style={{ width: 386 }}
            value={form.sortingIndex}
            onChange={(val) => setForm((f) => ({ ...f, sortingIndex: val ?? 0 }))}
          />
        </Block>
      </Modal>
    </PageHeader>
  );
};
