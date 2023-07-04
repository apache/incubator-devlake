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
import { useNavigate } from 'react-router-dom';
import { Button, Switch, Icon, Intent, Position } from '@blueprintjs/core';
import { Tooltip2 } from '@blueprintjs/popover2';

import { Card, IconButton, Dialog } from '@/components';
import { getCron } from '@/config';
import { PipelineContextProvider, PipelineInfo, PipelineTasks, PipelineHistorical } from '@/pages';
import { formatTime, operator } from '@/utils';

import { BlueprintType, FromEnum } from '../types';

import * as API from './api';
import * as S from './styled';

interface Props {
  from: FromEnum;
  blueprint: BlueprintType;
  pipelineId?: ID;
  onRefresh: () => void;
}

export const StatusPanel = ({ from, blueprint, pipelineId, onRefresh }: Props) => {
  const [isOpen, setIsOpen] = useState(false);
  const [operating, setOperating] = useState(false);

  const navigate = useNavigate();

  const cron = useMemo(() => getCron(blueprint.isManual, blueprint.cronConfig), [blueprint]);

  const handleShowDeleteDialog = () => {
    setIsOpen(true);
  };

  const handleHideDeleteDialog = () => {
    setIsOpen(false);
  };

  const handleRun = async (skipCollectors: boolean) => {
    const [success] = await operator(() => API.runBlueprint(blueprint.id, skipCollectors), {
      setOperating,
      formatMessage: () => 'Trigger blueprint successful.',
    });

    if (success) {
      onRefresh();
    }
  };

  const handleUpdate = async (payload: any) => {
    const [success] = await operator(
      () =>
        API.updateBlueprint(blueprint.id, {
          ...blueprint,
          ...payload,
        }),
      {
        setOperating,
        formatMessage: () => 'Update blueprint successful.',
      },
    );

    if (success) {
      onRefresh();
    }
  };

  const handleDelete = async () => {
    const [success] = await operator(() => API.deleteBluprint(blueprint.id), {
      setOperating,
      formatMessage: () => 'Delete blueprint successful.',
    });

    if (success) {
      navigate('/blueprints');
    }
  };

  return (
    <S.StatusPanel>
      {from === FromEnum.project && (
        <S.ProjectACtion>
          <span>
            {cron.value === 'manual' ? 'Manual' : `Next Run: ${formatTime(cron.nextTime, 'YYYY-MM-DD HH:mm')}`}
          </span>
          <Tooltip2
            position={Position.TOP}
            content="It is recommended to re-transform your data in this project if you have updated the transformation of the data scope in this project."
          >
            <Button
              disabled={!blueprint.enable}
              loading={operating}
              intent={Intent.PRIMARY}
              text="Re-transform Data"
              onClick={() => handleRun(true)}
            />
          </Tooltip2>
          <Button
            disabled={!blueprint.enable}
            loading={operating}
            intent={Intent.PRIMARY}
            text="Collect All Data"
            onClick={() => handleRun(false)}
          />
        </S.ProjectACtion>
      )}

      {from === FromEnum.blueprint && (
        <S.BlueprintAction>
          <Button text="Run Now" onClick={() => handleRun(false)} />
          <Switch
            style={{ marginBottom: 0 }}
            label="Blueprint Enabled"
            disabled={!!blueprint.projectName}
            checked={blueprint.enable}
            onChange={(e) => handleUpdate({ enable: (e.target as HTMLInputElement).checked })}
          />
          <IconButton
            loading={operating}
            disabled={!!blueprint.projectName}
            icon="trash"
            tooltip="Delete Blueprint"
            onClick={handleShowDeleteDialog}
          />
        </S.BlueprintAction>
      )}

      <PipelineContextProvider>
        <div className="block">
          <h3>Current Pipeline</h3>
          {!pipelineId ? (
            <Card>There is no current run for this blueprint.</Card>
          ) : (
            <>
              <PipelineInfo id={pipelineId} />
              <Card style={{ marginTop: 16 }}>
                <PipelineTasks id={pipelineId} />
              </Card>
            </>
          )}
        </div>
        <div className="block">
          <h3>Historical Pipelines</h3>
          <PipelineHistorical blueprintId={blueprint.id} />
        </div>
      </PipelineContextProvider>

      <Dialog
        isOpen={isOpen}
        style={{ width: 820 }}
        title="Are you sure you want to delete this Blueprint?"
        okText="Confirm"
        okLoading={operating}
        onCancel={handleHideDeleteDialog}
        onOk={handleDelete}
      >
        <S.DialogBody>
          <Icon icon="warning-sign" />
          <span>
            Please note: deleting the Blueprint will not delete the historical data of the Data Scopes in this
            Blueprint. If you would like to delete the historical data of Data Scopes, please visit the Connection page
            and do so.
          </span>
        </S.DialogBody>
      </Dialog>
    </S.StatusPanel>
  );
};
