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

const StageStatus = {
  PENDING: 'Pending',
  COMPLETE: 'Complete',
  FAILED: 'Failed',
  ACTIVE: 'In Progress',
}

const TaskStatus = {
  COMPLETE: 'TASK_COMPLETED',
  FAILED: 'TASK_FAILED',
  ACTIVE: 'TASK_RUNNING',
  RUNNING: 'TASK_RUNNING',
  CREATED: 'TASK_CREATED',
  PENDING: 'TASK_CREATED',
  CANCELLED: 'TASK_CANCELLED',
}

const TaskStatusLabels = {
  [TaskStatus.COMPLETE]: 'Succeeded',
  [TaskStatus.FAILED]: 'Failed',
  [TaskStatus.ACTIVE]: 'In Progress',
  [TaskStatus.RUNNING]: 'In Progress',
  [TaskStatus.CREATED]: 'Created (Pending)',
  [TaskStatus.PENDING]: 'Created (Pending)',
  [TaskStatus.CANCELLED]: 'Cancelled',
}

const StatusColors = {
  PENDING: '#292B3F',
  COMPLETE: '#4DB764',
  FAILED: '#E34040',
  ACTIVE: '#7497F7',
}

const StatusBgColors = {
  PENDING: 'transparent',
  COMPLETE: '#EDFBF0',
  FAILED: '#FEEFEF',
  ACTIVE: '#F0F4FE',
}

export {
  StageStatus,
  TaskStatus,
  TaskStatusLabels,
  StatusColors,
  StatusBgColors
}
