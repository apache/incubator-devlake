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

import { useState, useMemo, useEffect } from 'react';
import { PlusOutlined } from '@ant-design/icons';
import { Flex, Tabs, Button, Alert } from 'antd';

import API from '@/api';
import { PageHeader } from '@/components';
import { useRefreshData } from '@/hooks';
import { operator, parseCSV, stringifyCSV } from '@/utils';
import type { CSVTable } from '@/utils';

import { EditableTable } from './editable-table';
import type { EditableRow, EditableField } from './editable-table';

type Resource = 'teams' | 'users' | 'accounts';

let keySeed = 0;
const nextKey = () => `row-${(keySeed += 1)}`;

const withKeys = (table: CSVTable): EditableRow[] =>
  table.rows.map((row) => ({ ...row, _key: nextKey() }) as EditableRow);

const toCSV = (headers: string[], rows: EditableRow[]) =>
  stringifyCSV({
    headers,
    rows: rows.map(({ _key, ...rest }) => rest),
  });

const required = (field: string, label: string) => (row: EditableRow) =>
  row[field]?.trim() ? undefined : `${label} is required`;

const uniqueId = (row: EditableRow, rows: EditableRow[]) => {
  const id = row.Id?.trim();
  if (!id) {
    return 'Id is required';
  }
  return rows.filter((it) => it.Id?.trim() === id).length > 1 ? 'Id must be unique' : undefined;
};

export const Team = () => {
  const [tab, setTab] = useState<Resource>('teams');
  const [operating, setOperating] = useState(false);
  const [version, setVersion] = useState(1);

  const [teams, setTeams] = useState<EditableRow[]>([]);
  const [users, setUsers] = useState<EditableRow[]>([]);
  const [accounts, setAccounts] = useState<EditableRow[]>([]);
  const [headers, setHeaders] = useState<Record<Resource, string[]>>({
    teams: ['Id', 'Name', 'Alias', 'ParentId', 'SortingIndex'],
    users: ['Id', 'Name', 'Email', 'TeamIds'],
    accounts: [],
  });

  const { data, ready } = useRefreshData(
    () =>
      Promise.all([API.org.getTeams(), API.org.getUsers(), API.org.getUserAccountMapping()]).then(
        ([teamsCSV, usersCSV, accountsCSV]) => ({
          teams: parseCSV(teamsCSV),
          users: parseCSV(usersCSV),
          accounts: parseCSV(accountsCSV),
        }),
      ),
    [version],
  );

  useEffect(() => {
    if (!data) {
      return;
    }
    setTeams(withKeys(data.teams));
    setUsers(withKeys(data.users));
    setAccounts(withKeys(data.accounts));
    setHeaders({
      teams: data.teams.headers.length ? data.teams.headers : headers.teams,
      users: data.users.headers.length ? data.users.headers : headers.users,
      accounts: data.accounts.headers,
    });
  }, [data]);

  const teamOptions = useMemo(
    () =>
      teams
        .filter((team) => team.Id?.trim())
        .map((team) => ({ label: team.Name?.trim() || team.Id, value: team.Id.trim() })),
    [teams],
  );

  const userOptions = useMemo(
    () =>
      users
        .filter((user) => user.Id?.trim())
        .map((user) => ({ label: user.Name?.trim() || user.Id, value: user.Id.trim() })),
    [users],
  );

  const teamFields: EditableField[] = [
    { key: 'Id', title: 'Id', width: 120, placeholder: 'e.g. 1', invalid: uniqueId },
    { key: 'Name', title: 'Name', width: 220, placeholder: 'e.g. Platform', invalid: required('Name', 'Name') },
    { key: 'Alias', title: 'Alias', width: 140, placeholder: 'e.g. PLT' },
    {
      key: 'ParentId',
      title: 'Parent team',
      width: 200,
      type: 'select',
      placeholder: 'No parent',
      options: teamOptions,
    },
    { key: 'SortingIndex', title: 'Sorting index', width: 130, placeholder: '0' },
  ];

  const userFields: EditableField[] = [
    { key: 'Id', title: 'Id', width: 140, placeholder: 'e.g. 1', invalid: uniqueId },
    { key: 'Name', title: 'Name', width: 200, placeholder: 'Full name', invalid: required('Name', 'Name') },
    { key: 'Email', title: 'Email', width: 240, placeholder: 'name@example.com' },
    { key: 'TeamIds', title: 'Teams', type: 'multiple', placeholder: 'Select teams', options: teamOptions },
  ];

  const accountFields: EditableField[] = [
    { key: 'Id', title: 'Account Id', width: 260, editable: false },
    { key: 'UserName', title: 'Username', width: 160, editable: false },
    { key: 'Email', title: 'Email', width: 220, editable: false },
    { key: 'Organization', title: 'Organization', width: 160, editable: false },
    { key: 'UserId', title: 'Mapped user', width: 220, type: 'select', placeholder: 'Unmapped', options: userOptions },
  ];

  const invalid = useMemo(() => {
    const rows = tab === 'teams' ? teams : tab === 'users' ? users : [];
    const fields = tab === 'teams' ? teamFields : tab === 'users' ? userFields : [];
    return rows.some((row) => fields.some((field) => field.invalid?.(row, rows)));
  }, [tab, teams, users]);

  const handleAdd = () => {
    if (tab === 'teams') {
      setTeams([...teams, { _key: nextKey(), Id: '', Name: '', Alias: '', ParentId: '', SortingIndex: '0' }]);
    } else if (tab === 'users') {
      setUsers([...users, { _key: nextKey(), Id: '', Name: '', Email: '', TeamIds: '' }]);
    }
  };

  const handleSave = async () => {
    const payload = {
      teams: () => API.org.updateTeams(toCSV(headers.teams, teams)),
      users: () => API.org.updateUsers(toCSV(headers.users, users)),
      accounts: () => API.org.updateUserAccountMapping(toCSV(headers.accounts, accounts)),
    }[tab];

    const [success] = await operator(payload, {
      setOperating,
      formatMessage: () => 'Saved. Re-run the project pipeline to refresh team metrics.',
    });

    if (success) {
      setVersion(version + 1);
    }
  };

  return (
    <PageHeader
      breadcrumbs={[{ name: 'Teams', path: '/teams' }]}
      description="Define your teams, their members, and which tool accounts belong to each member. Team dashboards in Grafana read this data."
    >
      <Alert
        style={{ marginBottom: 16 }}
        type="info"
        showIcon
        message="Saving replaces the entire table"
        description="The org plugin stores each of these as a whole file, so a save overwrites every row — including rows added by someone else since this page loaded."
      />
      <Tabs
        activeKey={tab}
        onChange={(key) => setTab(key as Resource)}
        items={[
          { key: 'teams', label: `Teams (${teams.length})` },
          { key: 'users', label: `Members (${users.length})` },
          { key: 'accounts', label: `Account Mapping (${accounts.length})` },
        ]}
      />
      {tab === 'teams' && (
        <EditableTable
          fields={teamFields}
          rows={teams}
          onChange={setTeams}
          emptyText={ready ? 'No teams yet. Add one to get started.' : 'Loading…'}
        />
      )}
      {tab === 'users' && (
        <EditableTable
          fields={userFields}
          rows={users}
          onChange={setUsers}
          emptyText={ready ? 'No members yet. Add one, then assign them to a team.' : 'Loading…'}
        />
      )}
      {tab === 'accounts' && (
        <EditableTable
          fields={accountFields}
          rows={accounts}
          onChange={setAccounts}
          removable={false}
          emptyText={
            ready
              ? 'No accounts collected yet. Run a data collection first — accounts come from your connected tools.'
              : 'Loading…'
          }
        />
      )}
      <Flex style={{ marginTop: 16 }} justify="space-between">
        {tab === 'accounts' ? (
          <span />
        ) : (
          <Button icon={<PlusOutlined />} onClick={handleAdd}>
            Add {tab === 'teams' ? 'team' : 'member'}
          </Button>
        )}
        <Button type="primary" loading={operating} disabled={!ready || invalid} onClick={handleSave}>
          Save changes
        </Button>
      </Flex>
    </PageHeader>
  );
};
