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

import { Button, Checkbox, InputGroup, Intent } from '@blueprintjs/core';
import { useMemo, useState } from 'react';
import { Link, useHistory } from 'react-router-dom';

import { ColumnType, Dialog, IconButton, PageHeader, Table } from '@/components';

import * as S from './styled';
import { useProject } from './use-project';

type ProjectItem = {
  name: string;
};

export const ProjectHomePage = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [name, setName] = useState('');
  const [enableDora, setEnableDora] = useState(true);

  const history = useHistory();

  const handleShowDialog = () => setIsOpen(true);
  const handleHideDialog = () => {
    setIsOpen(false);
    setName('');
  };

  const { loading, operating, projects, onSave } = useProject<ProjectItem>({
    name,
    enableDora,
    onHideDialog: handleHideDialog,
  });

  const columns = useMemo(
    () =>
      [
        {
          title: 'Project Name',
          dataIndex: 'name',
          key: 'name',
          render: (name: string) => (
            <Link to={`/projects/${window.encodeURIComponent(name)}`} style={{ color: '#292b3f' }}>
              {name}
            </Link>
          ),
        },
        {
          title: '',
          dataIndex: 'name',
          key: 'action',
          width: 100,
          align: 'right',
          render: (name: any) => (
            <IconButton
              icon="cog"
              tooltip="Detail"
              onClick={() => history.push(`/projects/${window.encodeURIComponent(name)}`)}
            />
          ),
        },
        {
          title: '',
          dataIndex: 'name',
          key: 'action',
          width: 100,
          align: 'center',
          render: (name: any) => (
            <IconButton
              icon="trash"
              tooltip="Delete"
              onClick={() => history.push(`/projects/${window.encodeURIComponent(name)}`)}
            />
          ),
        },
      ] as ColumnType<ProjectItem>,
    [],
  );

  return (
    <PageHeader
      breadcrumbs={[{ name: 'Projects', path: '/projects' }]}
      extra={
        projects.length ? (
          <Button intent={Intent.PRIMARY} icon="plus" text="New Project" onClick={handleShowDialog} />
        ) : null
      }
    >
      <Table
        loading={loading}
        columns={columns}
        dataSource={projects}
        noData={{
          text: 'Add new projects to see engineering metrics based on projects.',
          btnText: 'New Project',
          onCreate: handleShowDialog,
        }}
      />
      <Dialog
        isOpen={isOpen}
        title="Create a New Project"
        style={{
          top: -100,
          width: 820,
        }}
        okText="Save"
        okDisabled={!name}
        okLoading={operating}
        onCancel={handleHideDialog}
        onOk={onSave}
      >
        <S.DialogInner>
          <div className="block">
            <h3>Project Name *</h3>
            <p>Give your project a unique name with letters, numbers, -, _ or /</p>
            <InputGroup placeholder="Your Project Name" value={name} onChange={(e) => setName(e.target.value)} />
          </div>
          <div className="block">
            <h3>Project Settings</h3>
            <div className="checkbox">
              <Checkbox
                label="Enable DORA Metrics"
                checked={enableDora}
                onChange={(e) => setEnableDora((e.target as HTMLInputElement).checked)}
              />
              <p>
                <a href="https://devlake.apache.org/docs/DORA/" rel="noreferrer" target="_blank">
                  DORA metrics
                </a>
                <span style={{ marginLeft: 4 }}>
                  are four widely-adopted metrics for measuring software delivery performance.
                </span>
              </p>
            </div>
          </div>
        </S.DialogInner>
      </Dialog>
    </PageHeader>
  );
};
