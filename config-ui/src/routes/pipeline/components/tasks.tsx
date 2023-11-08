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
import { Button, Collapse, Icon } from '@blueprintjs/core';
import { groupBy, sortBy } from 'lodash';

import API from '@/api';
import { Loading } from '@/components';
import { useAutoRefresh } from '@/hooks';
import { ITask, IPipelineStatus } from '@/types';

import * as S from '../styled';

import { PipelineTask } from './task';

interface Props {
  id: ID;
  style?: React.CSSProperties;
}

export const PipelineTasks = ({ id, style }: Props) => {
  const [isOpen, setIsOpen] = useState(true);

  // const { version } = usePipeline();

  const { data } = useAutoRefresh<ITask[]>(
    async () => {
      const taskRes = await API.pipeline.tasks(id);
      return taskRes.tasks;
    },
    [],
    {
      cancel: (data) => {
        return !!(
          data &&
          data.every((task) =>
            [IPipelineStatus.COMPLETED, IPipelineStatus.FAILED, IPipelineStatus.CANCELLED].includes(task.status),
          )
        );
      },
    },
  );

  const stages = groupBy(sortBy(data, 'id'), 'pipelineRow');

  const handleToggleOpen = () => setIsOpen(!isOpen);

  return (
    <S.Tasks>
      <div className="inner">
        <S.TasksHeader>
          {Object.keys(stages).map((key) => {
            let status;

            switch (true) {
              case !!stages[key].find((task) =>
                [IPipelineStatus.ACTIVE, IPipelineStatus.RUNNING].includes(task.status),
              ):
                status = 'loading';
                break;
              case stages[key].every((task) => task.status === IPipelineStatus.COMPLETED):
                status = 'success';
                break;
              case !!stages[key].find((task) => task.status === IPipelineStatus.FAILED):
                status = 'error';
                break;
              case !!stages[key].find((task) => task.status === IPipelineStatus.CANCELLED):
                status = 'cancel';
                break;
              default:
                status = 'ready';
                break;
            }

            return (
              <li key={key} className={status}>
                <strong>Stage {key}</strong>
                {status === 'loading' && <Loading size={14} />}
                {status === 'success' && <Icon icon="tick-circle" />}
                {status === 'error' && <Icon icon="cross-circle" />}
                {status === 'cancel' && <Icon icon="disable" />}
              </li>
            );
          })}
        </S.TasksHeader>
        <Collapse isOpen={isOpen}>
          <S.TasksList>
            {Object.keys(stages).map((key) => (
              <li key={key}>
                {stages[key].map((task) => (
                  <PipelineTask key={task.id} task={task} />
                ))}
              </li>
            ))}
          </S.TasksList>
        </Collapse>
      </div>
      <Button
        className="collapse-control"
        minimal
        icon={isOpen ? 'chevron-down' : 'chevron-up'}
        onClick={handleToggleOpen}
      />
    </S.Tasks>
  );
};
