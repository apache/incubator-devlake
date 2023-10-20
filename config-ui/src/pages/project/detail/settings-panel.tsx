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
import { InputGroup, Checkbox, Button, Icon, Intent } from '@blueprintjs/core';

import API from '@/api';
import { Card, FormItem, Buttons, toast, Dialog } from '@/components';
import { operator } from '@/utils';

import type { ProjectType } from '../types';
import { validName } from '../utils';

import * as S from './styled';

interface Props {
  project: ProjectType;
  onRefresh: () => void;
}

export const SettingsPanel = ({ project, onRefresh }: Props) => {
  const [name, setName] = useState('');
  const [enableDora, setEnableDora] = useState(false);
  const [operating, setOperating] = useState(false);
  const [isOpen, setIsOpen] = useState(false);

  const navigate = useNavigate();

  useEffect(() => {
    const doraMetrics = project.metrics.find((ms: any) => ms.pluginName === 'dora');

    setName(project.name);
    setEnableDora(doraMetrics?.enable ?? false);
  }, [project]);

  const handleUpdate = async () => {
    if (!validName(name)) {
      toast.error('Please enter alphanumeric or underscore');
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
              pluginOption: '',
              enable: enableDora,
            },
          ],
        }),
      {
        setOperating,
      },
    );

    if (success) {
      onRefresh();
      navigate(`/projects/${name}?tabId=settings`);
    }
  };

  const handleShowDeleteDialog = () => {
    setIsOpen(true);
  };

  const handleHideDeleteDialog = () => {
    setIsOpen(false);
  };

  const handleDelete = async () => {
    const [success] = await operator(() => API.project.remove(project.name), {
      setOperating,
      formatMessage: () => 'Delete project successful.',
    });

    if (success) {
      navigate(`/projects`);
    }
  };

  return (
    <>
      <Card>
        <FormItem label="Project Name" subLabel="Edit your project name with letters, numbers, -, _ or /" required>
          <InputGroup style={{ width: 386 }} value={name} onChange={(e) => setName(e.target.value)} />
        </FormItem>
        <FormItem subLabel="DORA metrics are four widely-adopted metrics for measuring software delivery performance.">
          <Checkbox
            label="Enable DORA Metrics"
            checked={enableDora}
            onChange={(e) => setEnableDora((e.target as HTMLInputElement).checked)}
          />
        </FormItem>
        <Buttons position="bottom">
          <Button text="Save" loading={operating} disabled={!name} intent={Intent.PRIMARY} onClick={handleUpdate} />
        </Buttons>
      </Card>
      <Buttons position="bottom" align="center">
        <Button intent={Intent.DANGER} text="Delete Project" onClick={handleShowDeleteDialog} />
      </Buttons>
      <Dialog
        isOpen={isOpen}
        style={{ width: 820 }}
        title="Are you sure you want to delete this Project?"
        okText="Confirm"
        okLoading={operating}
        onCancel={handleHideDeleteDialog}
        onOk={handleDelete}
      >
        <S.DialogBody>
          <Icon icon="warning-sign" />
          <span>
            This operation cannot be undone. Deleting a Data Connection will delete all data that have been collected in
            this Connection.
          </span>
        </S.DialogBody>
      </Dialog>
    </>
  );
};
