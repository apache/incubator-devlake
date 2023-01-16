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

import React from 'react';

import { Loading } from '@/components';
import { useAutoRefresh } from '@/hooks';
import { formatTime } from '@/utils';

import type { PipelineType } from '../../types';
import { StatusEnum } from '../../types';
import * as API from '../../api';

import { usePipeline } from '../context';
import { PipelineStatus } from '../status';
import { PipelineDuration } from '../duration';
import { PipelineCancel } from '../cancel';
import { PipelineRerun } from '../rerun';

import * as S from './styled';

interface Props {
  id: ID;
  style?: React.CSSProperties;
}

export const PipelineInfo = ({ id, style }: Props) => {
  const { version } = usePipeline();

  const { loading, data } = useAutoRefresh<PipelineType>(() => API.getPipeline(id), [version], {
    cancel: (data) => {
      return !!(
        data &&
        [StatusEnum.COMPLETED, StatusEnum.PARTIAL, StatusEnum.FAILED, StatusEnum.CANCELLED].includes(data.status)
      );
    },
  });

  if (loading || !data) {
    return <Loading />;
  }

  const { status, beganAt, finishedAt, stage, finishedTasks, totalTasks, message } = data;

  return (
    <S.Wrapper style={style}>
      <ul>
        <li>
          <span>Status</span>
          <strong>
            <PipelineStatus status={status} />
          </strong>
        </li>
        <li>
          <span>Started at</span>
          <strong>{formatTime(beganAt)}</strong>
        </li>
        <li>
          <span>Duration</span>
          <strong>
            <PipelineDuration status={status} beganAt={beganAt} finishedAt={finishedAt} />
          </strong>
        </li>
        <li>
          <span>Current Stage</span>
          <strong>{stage}</strong>
        </li>
        <li>
          <span>Tasks Completed</span>
          <strong>
            {finishedTasks}/{totalTasks}
          </strong>
        </li>
        <li>
          <PipelineCancel id={id} status={status} />
          <PipelineRerun type="pipeline" id={id} status={status} />
        </li>
      </ul>
      {StatusEnum.FAILED === status && <p className="'message'">{message}</p>}
    </S.Wrapper>
  );
};
