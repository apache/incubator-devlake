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

import { Space } from 'antd';
import classNames from 'classnames';

import { IPipelineStatus } from '@/types';

import * as S from '../styled';
import * as C from '../constant';

interface Props {
  status: IPipelineStatus;
}

export const PipelineStatus = ({ status }: Props) => {
  const cls = classNames({
    ready: [IPipelineStatus.CREATED, IPipelineStatus.PENDING].includes(status),
    loading: [IPipelineStatus.ACTIVE, IPipelineStatus.RUNNING, IPipelineStatus.RERUN].includes(status),
    success: [IPipelineStatus.COMPLETED, IPipelineStatus.PARTIAL].includes(status),
    error: status === IPipelineStatus.FAILED,
    cancel: status === IPipelineStatus.CANCELLED,
  });

  return (
    <S.StatusWrapper className={cls}>
      <Space size="small">
        {C.PipeLineStatusIcon[status]}
        <span>{C.PipeLineStatusLabel[status]}</span>
      </Space>
    </S.StatusWrapper>
  );
};
