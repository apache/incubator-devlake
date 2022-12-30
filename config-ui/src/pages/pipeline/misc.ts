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

import classNames from 'classnames';

import { StatusEnum } from './types';

export const STATUS_ICON = {
  [StatusEnum.CREATED]: 'stopwatch',
  [StatusEnum.PENDING]: 'stopwatch',
  [StatusEnum.ACTIVE]: 'loading',
  [StatusEnum.RUNNING]: 'loading',
  [StatusEnum.RERUN]: 'loading',
  [StatusEnum.COMPLETED]: 'tick-circle',
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
  [StatusEnum.FAILED]: 'Failed',
  [StatusEnum.CANCELLED]: 'Cancelled',
};

export const STATUS_CLS = (status: StatusEnum) =>
  classNames({
    ready: [StatusEnum.CREATED, StatusEnum.PENDING].includes(status),
    loading: [StatusEnum.ACTIVE, StatusEnum.RUNNING, StatusEnum.RERUN].includes(status),
    success: status === StatusEnum.COMPLETED,
    error: status === StatusEnum.FAILED,
    cancel: status === StatusEnum.CANCELLED,
  });
