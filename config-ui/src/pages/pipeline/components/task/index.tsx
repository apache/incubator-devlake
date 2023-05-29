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

import React, { useMemo } from 'react';
import { Intent } from '@blueprintjs/core';

import { TextTooltip } from '@/components';
import { getPluginConfig } from '@/plugins';

import type { TaskType } from '@/pages';
import { StatusEnum } from '@/pages';

import { PipelineDuration } from '../duration';
import { PipelineRerun } from '../rerun';

import * as S from './styled';

interface Props {
  task: TaskType;
}

export const PipelineTask = ({ task }: Props) => {
  const { id, beganAt, finishedAt, status, message, progressDetail } = task;

  const [icon, name] = useMemo(() => {
    const config = getPluginConfig(task.plugin);
    const options = JSON.parse(task.options);

    let name = config.name;

    switch (true) {
      case ['github', 'github_graphql'].includes(config.plugin):
        name = `${name}:${options.name}`;
        break;
      case ['gitextractor'].includes(config.plugin):
        name = `${name}:${options.repoId}`;
        break;
      case ['dora'].includes(config.plugin):
        name = `${name}:${options.projectName}`;
        break;
      case ['gitlab'].includes(config.plugin):
        name = `${name}:${options.projectId}`;
        break;
      case ['bitbucket'].includes(config.plugin):
        name = `${name}:${options.fullName}`;
        break;
      case ['tapd'].includes(config.plugin):
        name = `${name}:${options.workspaceId}`;
        break;
      case ['jira', 'jenkins'].includes(config.plugin):
        name = `${name}:${options.scopeId}`;
        break;
      case ['sonarqube'].includes(config.plugin):
        name = `${name}:${options.projectKey}`;
        break;
      case ['zentao'].includes(config.plugin):
        if (options.projectId) {
          name = `${name}:project/${options.projectId}`;
        } else {
          name = `${name}:product/${options.productId}`;
        }
        break;
      case ['kube_deployment'].includes(config.plugin):
        if (options.projectId) {
          name = `${name}:project/${options.id}`;
        } else {
          name = `${name}:product/${options.id}`;
        }
        break;
      case ['refdiff'].includes(config.plugin):
        name = `${name}:${options.repoId ?? options.projectName}`;
        break;
    }

    return [config.icon, name];
  }, [task]);

  return (
    <S.Wrapper>
      <S.Info>
        <div className="title">
          <img src={icon} alt="" />
          <strong>Task{task.id}</strong>
          <span>
            <TextTooltip content={name}>{name}</TextTooltip>
          </span>
        </div>
        {[status === StatusEnum.CREATED, StatusEnum.PENDING].includes(status) && <p>Subtasks pending</p>}

        {[StatusEnum.ACTIVE, StatusEnum.RUNNING].includes(status) && (
          <p>
            Subtasks running
            <strong style={{ marginLeft: 8 }}>
              {progressDetail?.finishedSubTasks}/{progressDetail?.totalSubTasks}
            </strong>
          </p>
        )}

        {status === StatusEnum.COMPLETED && <p>All Subtasks completed</p>}

        {status === StatusEnum.FAILED && (
          <TextTooltip intent={Intent.DANGER} content={message}>
            <p className="error">Task failed: hover to view the reason</p>
          </TextTooltip>
        )}

        {status === StatusEnum.CANCELLED && <p>Subtasks canceled</p>}
      </S.Info>
      <S.Duration>
        <PipelineDuration status={status} beganAt={beganAt} finishedAt={finishedAt} />
        <PipelineRerun type="task" id={id} status={status} />
      </S.Duration>
    </S.Wrapper>
  );
};
