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

import {
  FieldTimeOutlined,
  LoadingOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  UndoOutlined,
} from '@ant-design/icons';

import { IPipelineStatus } from '@/types';

export const PipeLineStatusIcon = {
  [IPipelineStatus.CREATED]: <FieldTimeOutlined rev={undefined} />,
  [IPipelineStatus.PENDING]: <FieldTimeOutlined rev={undefined} />,
  [IPipelineStatus.ACTIVE]: <LoadingOutlined rev={undefined} />,
  [IPipelineStatus.RUNNING]: <LoadingOutlined rev={undefined} />,
  [IPipelineStatus.RERUN]: <LoadingOutlined rev={undefined} />,
  [IPipelineStatus.COMPLETED]: <CheckCircleOutlined rev={undefined} />,
  [IPipelineStatus.PARTIAL]: <CheckCircleOutlined rev={undefined} />,
  [IPipelineStatus.FAILED]: <CloseCircleOutlined rev={undefined} />,
  [IPipelineStatus.CANCELLED]: <UndoOutlined rev={undefined} />,
};

export const PipeLineStatusLabel = {
  [IPipelineStatus.CREATED]: 'Created (Pending)',
  [IPipelineStatus.PENDING]: 'Created (Pending)',
  [IPipelineStatus.ACTIVE]: 'In Progress',
  [IPipelineStatus.RUNNING]: 'In Progress',
  [IPipelineStatus.RERUN]: 'In Progress',
  [IPipelineStatus.COMPLETED]: 'Succeeded',
  [IPipelineStatus.PARTIAL]: 'Partial Success',
  [IPipelineStatus.FAILED]: 'Failed',
  [IPipelineStatus.CANCELLED]: 'Cancelled',
};
