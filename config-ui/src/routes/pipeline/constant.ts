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

import * as T from './types';

export const PipeLineStatusIcon = {
  [T.PipelineStatus.CREATED]: 'stopwatch',
  [T.PipelineStatus.PENDING]: 'stopwatch',
  [T.PipelineStatus.ACTIVE]: 'loading',
  [T.PipelineStatus.RUNNING]: 'loading',
  [T.PipelineStatus.RERUN]: 'loading',
  [T.PipelineStatus.COMPLETED]: 'tick-circle',
  [T.PipelineStatus.PARTIAL]: 'tick-circle',
  [T.PipelineStatus.FAILED]: 'delete',
  [T.PipelineStatus.CANCELLED]: 'undo',
};

export const PipeLineStatusLabel = {
  [T.PipelineStatus.CREATED]: 'Created (Pending)',
  [T.PipelineStatus.PENDING]: 'Created (Pending)',
  [T.PipelineStatus.ACTIVE]: 'In Progress',
  [T.PipelineStatus.RUNNING]: 'In Progress',
  [T.PipelineStatus.RERUN]: 'In Progress',
  [T.PipelineStatus.COMPLETED]: 'Succeeded',
  [T.PipelineStatus.PARTIAL]: 'Partial Success',
  [T.PipelineStatus.FAILED]: 'Failed',
  [T.PipelineStatus.CANCELLED]: 'Cancelled',
};
