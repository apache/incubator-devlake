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

import React, { useState } from 'react';
import { Button, Collapse, Icon } from '@blueprintjs/core';
import { groupBy, sortBy } from 'lodash';

import { Loading } from '@/components';
import { useAutoRefresh } from '@/hooks';

import type { TaskType } from '../../types';
import { StatusEnum } from '../../types';
import * as API from '../../api';

import { usePipeline } from '../context';
import { PipelineTask } from '../task';

import * as S from './styled';

interface Props {
  id: ID;
  style?: React.CSSProperties;
}

export const PipelineTasks = ({ id, style }: Props) => {
  const [isOpen, setIsOpen] = useState(true);

  const { version } = usePipeline();

  const { loading, data } = useAutoRefresh<TaskType[]>(
    async () => {
      const taskRes = await API.getPipelineTasks(id);
      return taskRes.tasks;
    },
    [version],
    {
      cancel: (data) => {
        return !!(
          data &&
          data.every((task) => [StatusEnum.COMPLETED, StatusEnum.FAILED, StatusEnum.CANCELLED].includes(task.status))
        );
      },
    },
  );

  if (loading) {
    return <Loading />;
  }

  const stages = groupBy(sortBy(data, 'id'), 'pipelineRow');

  const handleToggleOpen = () => setIsOpen(!isOpen);

  return (
    <S.Wrapper style={style}>
      <S.Inner>
        <S.Header>
          {Object.keys(stages).map((key) => {
            let status;

            switch (true) {
              case !!stages[key].find((task) => [StatusEnum.ACTIVE, StatusEnum.RUNNING].includes(task.status)):
                status = 'loading';
                break;
              case stages[key].every((task) => task.status === StatusEnum.COMPLETED):
                status = 'success';
                break;
              case !!stages[key].find((task) => task.status === StatusEnum.FAILED):
                status = 'error';
                break;
              case !!stages[key].find((task) => task.status === StatusEnum.CANCELLED):
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
        </S.Header>
        <Collapse isOpen={isOpen}>
          <S.Tasks>
            {Object.keys(stages).map((key) => (
              <li key={key}>
                {stages[key].map((task) => (
                  <PipelineTask key={task.id} task={task} />
                ))}
              </li>
            ))}
          </S.Tasks>
        </Collapse>
      </S.Inner>
      <Button
        className="collapse-control"
        minimal
        icon={isOpen ? 'chevron-down' : 'chevron-up'}
        onClick={handleToggleOpen}
      />
    </S.Wrapper>
  );
};
