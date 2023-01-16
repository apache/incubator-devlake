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

import { request } from '@/utils';

export const getPipeline = (id: ID) => request(`/pipelines/${id}`);

export const getPipelineTasks = (id: ID) => request(`/pipelines/${id}/tasks`);

export const deletePipeline = (id: ID) =>
  request(`/pipelines/${id}`, {
    method: 'delete',
  });

export const getPipelineHistorical = (id: ID) => request(`/blueprints/${id}/pipelines`);

export const pipelineRerun = (id: ID) =>
  request(`/pipelines/${id}/rerun`, {
    method: 'post',
  });

export const taskRerun = (id: ID) =>
  request(`/tasks/${id}/rerun`, {
    method: 'post',
  });

export const getPipelineLog = (id: ID) => request(`/pipelines/${id}/logging.tar.gz`);
