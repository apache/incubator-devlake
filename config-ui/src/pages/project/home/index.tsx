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
import { Link, useNavigate } from 'react-router-dom';
import { Button, InputGroup, Checkbox, Intent, FormGroup } from '@blueprintjs/core';
import dayjs from 'dayjs';

import { PageHeader, Table, Dialog, IconButton, toast } from '@/components';
import { getCron, cronPresets } from '@/config';
import { useConnections, useRefreshData } from '@/hooks';
import { formatTime, operator } from '@/utils';
import { PipelineStatus } from '@/pages/pipeline';

import { validName, encodeName } from '../utils';
import { BlueprintType, ModeEnum } from '../../blueprint';

import * as API from './api';
import * as S from './styled';

export const ProjectHomePage = () => {
  const [version, setVersion] = useState(1);
  const [isOpen, setIsOpen] = useState(false);
  const [name, setName] = useState('');
  const [enableDora, setEnableDora] = useState(true);
  const [saving, setSaving] = useState(false);

  const { ready, data } = useRefreshData(() => API.getProjects({ page: 1, pageSize: 200 }), [version]);
  const { onGet } = useConnections();

  const navigate = useNavigate();

  const presets = useMemo(() => cronPresets.map((preset) => preset.config), []);
  const dataSource = useMemo(
    () =>
      (data?.projects ?? []).map((it) => {
        return {
          name: it.name,
          connections: it.blueprint.settings.connections,
          isManual: it.blueprint.isManual,
          cronConfig: it.blueprint.cronConfig,
          createdAt: it.createdAt,
          lastRunCompletedAt: it.lastPipeline?.finishedAt,
          lastRunStatus: it.lastPipeline?.status,
        };
      }),
    [data],
  );

  const handleShowDialog = () => setIsOpen(true);
  const handleHideDialog = () => {
    setIsOpen(false);
    setName('');
    setEnableDora(true);
  };

  const handleCreate = async () => {
    if (!validName(name)) {
      toast.error('Please enter alphanumeric or underscore');
      return;
    }

    const [success] = await operator(
      async () => {
        await API.createProject({
          name: encodeName(name),
          description: '',
          metrics: [
            {
              pluginName: 'dora',
              pluginOption: '',
              enable: enableDora,
            },
          ],
        });
        return API.createBlueprint({
          name: `${name}-Blueprint`,
          projectName: name,
          mode: ModeEnum.normal,
          enable: true,
          cronConfig: presets[0],
          isManual: false,
          skipOnFail: true,
          settings: {
            version: '2.0.0',
            timeAfter: formatTime(dayjs().subtract(6, 'month').startOf('day').toDate(), 'YYYY-MM-DD[T]HH:mm:ssZ'),
            connections: [],
          },
        });
      },
      {
        setOperating: setSaving,
      },
    );

    if (success) {
      handleHideDialog();
      setVersion((v) => v + 1);
    }
  };

  return (
    <PageHeader
      breadcrumbs={[{ name: 'Projects', path: '/projects' }]}
      extra={<Button intent={Intent.PRIMARY} icon="plus" text="New Project" onClick={handleShowDialog} />}
    >
      <Table
        loading={!ready}
        columns={[
          {
            title: 'Project Name',
            dataIndex: 'name',
            key: 'name',
            render: (name: string) => (
              <Link to={`/projects/${encodeName(name)}?tab=configuration`} style={{ color: '#292b3f' }}>
                {name}
              </Link>
            ),
          },
          {
            title: 'Data Connections',
            dataIndex: 'connections',
            key: 'connections',
            render: (val: BlueprintType['settings']['connections']) =>
              !val.length
                ? 'N/A'
                : val
                    .map((it) => {
                      const cs = onGet(`${it.plugin}-${it.connectionId}`);
                      return cs.name;
                    })
                    .join(', '),
          },
          {
            title: 'Sync Frequency',
            dataIndex: ['isManual', 'cronConfig'],
            key: 'frequency',
            render: ({ isManual, cronConfig }) => {
              const cron = getCron(isManual, cronConfig);
              return cron.label;
            },
          },
          {
            title: 'Created at',
            dataIndex: 'createdAt',
            key: 'createdAt',
            render: (val) => formatTime(val),
          },
          {
            title: 'Last Run Completed at',
            dataIndex: 'lastRunCompletedAt',
            key: 'lastRunCompletedAt',
            render: (val) => (val ? formatTime(val) : '-'),
          },
          {
            title: 'Last Run Status',
            dataIndex: 'lastRunStatus',
            key: 'lastRunStatus',
            render: (val) => (val ? <PipelineStatus status={val} /> : '-'),
          },
          {
            title: '',
            dataIndex: 'name',
            key: 'action',
            width: 100,
            align: 'center',
            render: (name: any) => (
              <IconButton
                icon="cog"
                tooltip="Detail"
                onClick={() => navigate(`/projects/${encodeName(name)}?tab=configuration`)}
              />
            ),
          },
        ]}
        dataSource={dataSource}
        noData={{
          text: 'Add new projects to see engineering metrics based on projects.',
          btnText: 'New Project',
          onCreate: handleShowDialog,
        }}
      />
      <Dialog
        isOpen={isOpen}
        title="Create a New Project"
        style={{ width: 820 }}
        okText="Save"
        okDisabled={!name}
        okLoading={saving}
        onOk={handleCreate}
        onCancel={handleHideDialog}
      >
        <S.DialogWrapper>
          <FormGroup
            label={<S.Label>Project Name</S.Label>}
            subLabel={
              <S.LabelDescription>Give your project a unique name with letters, numbers, -, _ or /</S.LabelDescription>
            }
            labelInfo={<S.LabelInfo>*</S.LabelInfo>}
          >
            <InputGroup
              style={{ width: 386 }}
              placeholder="Your Project Name"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </FormGroup>

          <FormGroup
            label={<S.Label>Project Settings</S.Label>}
            subLabel={
              <S.LabelDescription>
                <a href="https://devlake.apache.org/docs/DORA/" rel="noreferrer" target="_blank">
                  DORA metrics
                </a>
                <span style={{ marginLeft: 4 }}>
                  are four widely-adopted metrics for measuring software delivery performance.
                </span>
              </S.LabelDescription>
            }
          >
            <Checkbox
              label="Enable DORA Metrics"
              checked={enableDora}
              onChange={(e) => setEnableDora((e.target as HTMLInputElement).checked)}
            />
          </FormGroup>
        </S.DialogWrapper>
      </Dialog>
    </PageHeader>
  );
};
