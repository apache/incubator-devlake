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

import { useState } from 'react';

import API from '@/api';
import { Loading, IconButton } from '@/components';
import { useAutoRefresh } from '@/hooks';
import { formatTime, operator } from '@/utils';

import * as T from '../types';
import * as S from '../styled';

import { PipelineStatus } from './status';
import { PipelineDuration } from './duration';

interface Props {
  id: ID;
}

export const PipelineInfo = ({ id }: Props) => {
  const [operating, setOperating] = useState(false);

  const { data } = useAutoRefresh<T.Pipeline>(() => API.pipeline.get(id), [], {
    cancel: (data) => {
      return !!(
        data &&
        [
          T.PipelineStatus.COMPLETED,
          T.PipelineStatus.PARTIAL,
          T.PipelineStatus.FAILED,
          T.PipelineStatus.CANCELLED,
        ].includes(data.status)
      );
    },
  });

  const handleCancel = async () => {
    const [success] = await operator(() => API.pipeline.remove(id), {
      setOperating,
    });

    if (success) {
      // setVersion((v) => v + 1);
    }
  };

  const handleRerun = async () => {
    const [success] = await operator(() => API.pipeline.rerun(id), {
      setOperating,
    });

    if (success) {
      // setVersion((v) => v + 1);
    }
  };

  if (!data) {
    return <Loading />;
  }

  const { status, beganAt, finishedAt, stage, finishedTasks, totalTasks, message } = data;

  return (
    <S.Info>
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
          {[T.PipelineStatus.ACTIVE, T.PipelineStatus.RUNNING, T.PipelineStatus.RERUN].includes(status) && (
            <IconButton loading={operating} icon="disable" tooltip="Cancel" onClick={handleCancel} />
          )}
          {[
            T.PipelineStatus.COMPLETED,
            T.PipelineStatus.PARTIAL,
            T.PipelineStatus.FAILED,
            T.PipelineStatus.CANCELLED,
          ].includes(status) && (
            <IconButton loading={operating} icon="repeat" tooltip="Rerun failed tasks" onClick={handleRerun} />
          )}
        </li>
      </ul>
      {T.PipelineStatus.FAILED === status && <p className="'message'">{message}</p>}
    </S.Info>
  );
};
