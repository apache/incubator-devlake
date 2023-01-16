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

export enum StatusEnum {
  CREATED = 'TASK_CREATED',
  PENDING = 'TASK_PENDING',
  ACTIVE = 'TASK_ACTIVE',
  RUNNING = 'TASK_RUNNING',
  RERUN = 'TASK_RERUN',
  COMPLETED = 'TASK_COMPLETED',
  PARTIAL = 'TASK_PARTIAL',
  FAILED = 'TASK_FAILED',
  CANCELLED = 'TASK_CANCELLED',
}

export type PipelineType = {
  id: ID;
  status: StatusEnum;
  beganAt: string | null;
  finishedAt: string | null;
  stage: number;
  finishedTasks: number;
  totalTasks: number;
  message: string;
};

export type TaskType = {
  id: ID;
  plugin: string;
  status: StatusEnum;
  pipelineRow: number;
  pipelineCol: number;
  beganAt: string | null;
  finishedAt: string | null;
  options: string;
  message: string;
  progressDetail?: {
    finishedSubTasks: number;
    totalSubTasks: number;
  };
};
