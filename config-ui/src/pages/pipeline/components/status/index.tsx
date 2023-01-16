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
import { Icon, IconName } from '@blueprintjs/core';
import classNames from 'classnames';

import { Loading } from '@/components';

import { StatusEnum } from '../../types';

import * as S from './styled';

export const STATUS_ICON = {
  [StatusEnum.CREATED]: 'stopwatch',
  [StatusEnum.PENDING]: 'stopwatch',
  [StatusEnum.ACTIVE]: 'loading',
  [StatusEnum.RUNNING]: 'loading',
  [StatusEnum.RERUN]: 'loading',
  [StatusEnum.COMPLETED]: 'tick-circle',
  [StatusEnum.PARTIAL]: 'tick-circle',
  [StatusEnum.FAILED]: 'delete',
  [StatusEnum.CANCELLED]: 'undo',
};

export const STATUS_LABEL = {
  [StatusEnum.CREATED]: 'Created (Pending)',
  [StatusEnum.PENDING]: 'Created (Pending)',
  [StatusEnum.ACTIVE]: 'In Progress',
  [StatusEnum.RUNNING]: 'In Progress',
  [StatusEnum.RERUN]: 'In Progress',
  [StatusEnum.COMPLETED]: 'Succeeded',
  [StatusEnum.PARTIAL]: 'Partial Success',
  [StatusEnum.FAILED]: 'Failed',
  [StatusEnum.CANCELLED]: 'Cancelled',
};

interface Props {
  status: StatusEnum;
}

export const PipelineStatus = ({ status }: Props) => {
  const statusCls = classNames({
    ready: [StatusEnum.CREATED, StatusEnum.PENDING].includes(status),
    loading: [StatusEnum.ACTIVE, StatusEnum.RUNNING, StatusEnum.RERUN].includes(status),
    success: [StatusEnum.COMPLETED, StatusEnum.PARTIAL].includes(status),
    error: status === StatusEnum.FAILED,
    cancel: status === StatusEnum.CANCELLED,
  });

  return (
    <S.Wrapper className={statusCls}>
      {STATUS_ICON[status] === 'loading' ? (
        <Loading style={{ marginRight: 4 }} size={14} />
      ) : (
        <Icon style={{ marginRight: 4 }} icon={STATUS_ICON[status] as IconName} />
      )}
      <span>{STATUS_LABEL[status]}</span>
    </S.Wrapper>
  );
};
