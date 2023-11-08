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
import { Link } from 'react-router-dom';
import { ButtonGroup, Button, Tag, Intent, FormGroup, InputGroup, RadioGroup, Radio } from '@blueprintjs/core';
import dayjs from 'dayjs';

import API from '@/api';
import { PageHeader, Table, IconButton, TextTooltip, Dialog } from '@/components';
import { getCronOptions, cronPresets, getCron } from '@/config';
import { ConnectionName } from '@/features';
import { useRefreshData } from '@/hooks';
import { IBlueprint, IBPMode } from '@/types';
import { formatTime, operator } from '@/utils';

import * as S from './styled';

export const BlueprintHomePage = () => {
  const [version, setVersion] = useState(1);
  const [type, setType] = useState('all');
  const [page, setPage] = useState(1);
  const [pageSize] = useState(20);
  const [isOpen, setIsOpen] = useState(false);
  const [name, setName] = useState('');
  const [mode, setMode] = useState(IBPMode.NORMAL);
  const [saving, setSaving] = useState(false);

  const { ready, data } = useRefreshData(
    () => API.blueprint.list({ type: type.toLocaleUpperCase(), page, pageSize }),
    [version, type, page, pageSize],
  );

  const [options, presets] = useMemo(() => [getCronOptions(), cronPresets.map((preset) => preset.config)], []);
  const [dataSource, total] = useMemo(() => [data?.blueprints ?? [], data?.count ?? 0], [data]);

  const handleShowDialog = () => setIsOpen(true);
  const handleHideDialog = () => {
    setName('');
    setMode(IBPMode.NORMAL);
    setIsOpen(false);
  };

  const handleCreate = async () => {
    const payload: any = {
      name,
      mode,
      enable: true,
      cronConfig: presets[0],
      isManual: false,
      skipOnFail: true,
    };

    if (mode === IBPMode.NORMAL) {
      payload.timeAfter = formatTime(dayjs().subtract(6, 'month').startOf('day').toDate(), 'YYYY-MM-DD[T]HH:mm:ssZ');
      payload.connections = [];
    }

    if (mode === IBPMode.ADVANCED) {
      payload.timeAfter = undefined;
      payload.connections = undefined;
      payload.plan = [[]];
    }

    const [success] = await operator(() => API.blueprint.create(payload), {
      setOperating: setSaving,
    });

    if (success) {
      handleHideDialog();
      setVersion((v) => v + 1);
    }
  };

  return (
    <PageHeader
      breadcrumbs={[
        { name: 'Advanced', path: '/blueprints' },
        { name: 'Blueprints', path: '/blueprints' },
      ]}
    >
      <S.Wrapper>
        <p>This is a complete list of all Blueprints you have created, whether they belong to Projects or not.</p>
        <div className="action">
          <ButtonGroup>
            <Button intent={type === 'all' ? Intent.PRIMARY : Intent.NONE} text="All" onClick={() => setType('all')} />
            {options.map(({ label }) => (
              <Button
                key={label}
                intent={type === label ? Intent.PRIMARY : Intent.NONE}
                text={label}
                onClick={() => setType(label)}
              />
            ))}
          </ButtonGroup>
          <Button icon="plus" intent={Intent.PRIMARY} text="New Blueprint" onClick={handleShowDialog} />
        </div>
        <Table
          loading={!ready}
          columns={[
            {
              title: 'Blueprint Name',
              dataIndex: ['id', 'name'],
              key: 'name',
              render: ({ id, name }) => (
                <Link to={`/blueprints/${id}?tab=configuration`} style={{ color: '#292b3f' }}>
                  <TextTooltip content={name}>{name}</TextTooltip>
                </Link>
              ),
            },
            {
              title: 'Data Connections',
              dataIndex: ['mode', 'connections'],
              key: 'connections',
              render: ({ mode, connections }: Pick<IBlueprint, 'mode' | 'connections'>) => {
                if (mode === IBPMode.ADVANCED) {
                  return 'Advanced Mode';
                }
                return (
                  <S.ConnectionList>
                    {connections.map((it) => (
                      <li>
                        <ConnectionName
                          key={`${it.pluginName}-${it.connectionId}`}
                          plugin={it.pluginName}
                          connectionId={it.connectionId}
                        />
                      </li>
                    ))}
                  </S.ConnectionList>
                );
              },
            },
            {
              title: 'Frequency',
              dataIndex: ['isManual', 'cronConfig'],
              key: 'frequency',
              width: 100,
              render: ({ isManual, cronConfig }) => {
                const cron = getCron(isManual, cronConfig);
                return cron.label;
              },
            },
            {
              title: 'Next Run Time',
              dataIndex: ['isManual', 'cronConfig'],
              key: 'nextRunTime',
              width: 200,
              render: ({ isManual, cronConfig }) => {
                const cron = getCron(isManual, cronConfig);
                return formatTime(cron.nextTime);
              },
            },
            {
              title: 'Project',
              dataIndex: 'projectName',
              key: 'project',
              render: (val) =>
                val ? (
                  <Link to={`/projects/${window.encodeURIComponent(val)}`}>
                    <TextTooltip content={val}>{val}</TextTooltip>
                  </Link>
                ) : (
                  'N/A'
                ),
            },
            {
              title: 'Status',
              dataIndex: 'enable',
              key: 'enable',
              align: 'center',
              width: 100,
              render: (val) => (
                <Tag minimal intent={val ? Intent.SUCCESS : Intent.DANGER}>
                  {val ? 'Enabled' : 'Disabled'}
                </Tag>
              ),
            },
            {
              title: '',
              dataIndex: 'id',
              key: 'action',
              width: 100,
              align: 'center',
              render: (val) => (
                <Link to={`/blueprints/${val}?tab=configuration`}>
                  <IconButton icon="cog" tooltip="Detail" />
                </Link>
              ),
            },
          ]}
          dataSource={dataSource}
          pagination={{
            page,
            pageSize,
            total,
            onChange: setPage,
          }}
          noData={{
            text: 'There is no Blueprint yet. Please add a new Blueprint here or from a Project.',
            btnText: 'New Blueprint',
            onCreate: handleShowDialog,
          }}
        />
      </S.Wrapper>
      <Dialog
        style={{ width: 820 }}
        isOpen={isOpen}
        title="Create a New Blueprint"
        okText="Save"
        okDisabled={!name}
        okLoading={saving}
        onOk={handleCreate}
        onCancel={handleHideDialog}
      >
        <S.DialogWrapper>
          <FormGroup
            label={<S.Label>Blueprint Name</S.Label>}
            subLabel={
              <S.LabelDescription>
                Give your Blueprint a unique name to help you identify it in the future.
              </S.LabelDescription>
            }
            labelInfo={<S.LabelInfo>*</S.LabelInfo>}
          >
            <InputGroup
              style={{ width: 386 }}
              placeholder="Your Blueprint Name"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </FormGroup>
          <FormGroup
            label={<S.Label>Blueprint Mode</S.Label>}
            subLabel={
              <S.LabelDescription>
                Normal Mode is usually adequate for most usages. But if you need to customize how tasks are executed in
                the Blueprint, please use Advanced Mode to create a Blueprint.
              </S.LabelDescription>
            }
            labelInfo={<S.LabelInfo>*</S.LabelInfo>}
          >
            <RadioGroup
              inline
              selectedValue={mode}
              onChange={(e) => setMode((e.target as HTMLInputElement).value as IBPMode)}
            >
              <Radio value={IBPMode.NORMAL}>Normal Mode</Radio>
              <Radio value={IBPMode.ADVANCED}>Advanced Mode</Radio>
            </RadioGroup>
          </FormGroup>
        </S.DialogWrapper>
      </Dialog>
    </PageHeader>
  );
};
