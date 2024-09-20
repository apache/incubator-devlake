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

import dayjs from 'dayjs';

import { IPipelineStatus } from '@/types';

const duration = (minute: number) => {
  if (minute < 1) {
    return '< 1m';
  }

  if (minute < 60) {
    return `${Math.ceil(minute / 60)}m`;
  }

  if (minute < 60 * 24) {
    const hours = Math.floor(minute / 60);
    const minutes = minute - hours * 60;
    return `${hours}h${minutes}m`;
  }

  const days = Math.floor(minute / (60 * 24));
  const hours = Math.floor((minute - days * 60 * 24) / 60);
  const minutes = minute - days * 60 * 24 - hours * 60;

  return `${days}d${hours}h${minutes}m`;
};

interface Props {
  status: IPipelineStatus;
  beganAt: string | null;
  finishedAt: string | null;
}

export const PipelineDuration = ({ status, beganAt, finishedAt }: Props) => {
  if (!beganAt) {
    return <span>-</span>;
  }

  if (
    ![IPipelineStatus.CANCELLED, IPipelineStatus.COMPLETED, IPipelineStatus.PARTIAL, IPipelineStatus.FAILED].includes(
      status,
    )
  ) {
    return <span>{duration(dayjs(beganAt).diff(dayjs(), 'm'))}</span>;
  }

  if (!finishedAt) {
    return <span>-</span>;
  }

  return <span>{duration(dayjs(beganAt).diff(dayjs(finishedAt), 'm'))}</span>;
};
