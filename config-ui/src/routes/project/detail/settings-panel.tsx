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

import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { WarningOutlined } from '@ant-design/icons';
import { Flex, Space, Card, Modal, Input, Checkbox, Button, message } from 'antd';

import API from '@/api';
import { Block, HelpTooltip } from '@/components';
import { PATHS } from '@/config';
import { IProject } from '@/types';
import { operator } from '@/utils';

import { validName } from '../utils';

import * as S from './styled';

interface Props {
  project: IProject;
  onRefresh: () => void;
}

export const SettingsPanel = ({ project, onRefresh }: Props) => {
  const [name, setName] = useState('');
  const [dora, setDora] = useState({
    enable: false,
  });
  const [linker, setLinker] = useState({
    enable: false,
    prToIssueRegexp: '',
  });
  const [operating, setOperating] = useState(false);
  const [open, setOpen] = useState(false);

  const navigate = useNavigate();

  useEffect(() => {
    const dora = project.metrics.find((ms) => ms.pluginName === 'dora');
    const linker = project.metrics.find((ms) => ms.pluginName === 'linker');

    setName(project.name);
    setDora({
      enable: dora?.enable ?? false,
    });
    setLinker({
      enable: linker?.enable ?? false,
      prToIssueRegexp: linker?.pluginOption?.prToIssueRegexp ?? '',
    });
  }, [project]);

  const handleUpdate = async () => {
    if (!validName(name)) {
      message.error('Please enter alphanumeric or underscore');
      return;
    }

    const [success] = await operator(
      () =>
        API.project.update(project.name, {
          name,
          description: '',
          metrics: [
            {
              pluginName: 'dora',
              pluginOption: {},
              enable: dora.enable,
            },
            {
              pluginName: 'linker',
              pluginOption: {
                prToIssueRegexp: linker.prToIssueRegexp,
              },
              enable: linker.enable,
            },
          ],
        }),
      {
        setOperating,
      },
    );

    if (success) {
      onRefresh();
      navigate(PATHS.PROJECT(name, { tabId: 'settings' }));
    }
  };

  const handleShowDeleteDialog = () => {
    setOpen(true);
  };

  const handleHideDeleteDialog = () => {
    setOpen(false);
  };

  const handleDelete = async () => {
    const [success] = await operator(() => API.project.remove(project.name), {
      setOperating,
      formatMessage: () => 'Delete project successful.',
    });

    if (success) {
      navigate(PATHS.PROJECTS());
    }
  };

  return (
    <Flex vertical>
      <Space direction="vertical" size="large">
        <Card>
          <Block title="Project Name" description="Edit your project name with letters, numbers, -, _ or /" required>
            <Input style={{ width: 386 }} value={name} onChange={(e) => setName(e.target.value)} />
          </Block>
          <Block
            title={
              <Checkbox checked={dora.enable} onChange={(e) => setDora({ enable: e.target.checked })}>
                Enable DORA Metrics
              </Checkbox>
            }
            description="DORA metrics are four widely-adopted metrics for measuring software delivery performance."
          />
          <Block
            title={
              <Checkbox checked={linker.enable} onChange={(e) => setLinker({ ...linker, enable: e.target.checked })}>
                Associate pull requests with issues
              </Checkbox>
            }
            description={
              <span>
                Parse the issue key with the regex from the title and description of the pull requests in this project.
                <HelpTooltip
                  content={
                    <span>
                      The default regex will parse the issue key from the table.pull_requests where the description
                      contains "fix/close/.../resolved {'{'}issue_key{'}'}". The relationship between pull requests and
                      issues will be stored in the table.pull_request_issues
                    </span>
                  }
                />
              </span>
            }
          >
            {linker.enable && (
              <Input
                value={linker.prToIssueRegexp}
                onChange={(e) => setLinker({ ...linker, prToIssueRegexp: e.target.value })}
              />
            )}
          </Block>
          <Block>
            <Button type="primary" loading={operating} disabled={!name} onClick={handleUpdate}>
              Save
            </Button>
          </Block>
        </Card>
        <Flex justify="center">
          <Button type="primary" danger onClick={handleShowDeleteDialog}>
            Delete Project
          </Button>
        </Flex>
      </Space>
      <Modal
        open={open}
        width={820}
        centered
        title="Are you sure you want to delete this Project?"
        okText="Confirm"
        okButtonProps={{
          loading: operating,
        }}
        onCancel={handleHideDeleteDialog}
        onOk={handleDelete}
      >
        <S.DialogBody>
          <WarningOutlined />
          <span>
            This operation cannot be undone. Deleting a Data Connection will delete all data that have been collected in
            this Connection.
          </span>
        </S.DialogBody>
      </Modal>
    </Flex>
  );
};
