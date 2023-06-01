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

import { useMemo } from 'react';
import { Button, Intent, Position } from '@blueprintjs/core';
import { Tooltip2 } from '@blueprintjs/popover2';

import { Card } from '@/components';
import { getCron } from '@/config';
import { PipelineContextProvider, PipelineInfo, PipelineTasks, PipelineHistorical } from '@/pages';
import { formatTime } from '@/utils';

import type { BlueprintType } from '../types';

import * as S from './styled';

interface Props {
  blueprint: BlueprintType;
  pipelineId?: ID;
  operating: boolean;
  onRun: (skipCollectors: boolean) => void;
}

export const StatusPanel = ({ blueprint, pipelineId, operating, onRun }: Props) => {
  const cron = useMemo(() => getCron(blueprint.isManual, blueprint.cronConfig), [blueprint]);

  return (
    <S.StatusPanel>
      <div className="info">
        <span>{cron.value === 'manual' ? 'Manual' : `Next Run: ${formatTime(cron.nextTime, 'YYYY-MM-DD HH:mm')}`}</span>
        <Tooltip2
          position={Position.TOP}
          content="It is recommended to re-transform your data in this project if you have updated the transformation of the data scope in this project."
        >
          <Button
            disabled={!blueprint.enable}
            loading={operating}
            intent={Intent.PRIMARY}
            text="Re-transform Data"
            onClick={() => onRun(true)}
          />
        </Tooltip2>
        <Button
          disabled={!blueprint.enable}
          loading={operating}
          intent={Intent.PRIMARY}
          text="Collect All Data"
          onClick={() => onRun(false)}
        />
      </div>
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
    </S.StatusPanel>
  );
};
