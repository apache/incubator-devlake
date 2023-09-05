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

import { Icon, IconName } from '@blueprintjs/core';
import classNames from 'classnames';

import { Loading } from '@/components';

import * as T from '../types';
import * as S from '../styled';
import * as C from '../constant';

interface Props {
  status: T.PipelineStatus;
}

export const PipelineStatus = ({ status }: Props) => {
  const cls = classNames({
    ready: [T.PipelineStatus.CREATED, T.PipelineStatus.PENDING].includes(status),
    loading: [T.PipelineStatus.ACTIVE, T.PipelineStatus.RUNNING, T.PipelineStatus.RERUN].includes(status),
    success: [T.PipelineStatus.COMPLETED, T.PipelineStatus.PARTIAL].includes(status),
    error: status === T.PipelineStatus.FAILED,
    cancel: status === T.PipelineStatus.CANCELLED,
  });

  return (
    <S.StatusWrapper className={cls}>
      {C.PipeLineStatusIcon[status] === 'loading' ? (
        <Loading style={{ marginRight: 4 }} size={14} />
      ) : (
        <Icon style={{ marginRight: 4 }} icon={C.PipeLineStatusIcon[status] as IconName} />
      )}
      <span>{C.PipeLineStatusLabel[status]}</span>
    </S.StatusWrapper>
  );
};
