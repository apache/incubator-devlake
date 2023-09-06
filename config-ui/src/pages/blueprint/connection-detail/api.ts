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

import { BlueprintType } from '../types';

export const getProject = (pname: string): Promise<{ blueprint: BlueprintType }> => request(`/projects/${pname}`);

export const getBlueprint = (id: ID): Promise<BlueprintType> => request(`/blueprints/${id}`);

export const updateBlueprint = (id: ID, payload: any) =>
  request(`/blueprints/${id}`, { method: 'patch', data: payload });

export const getConnection = (plugin: string, connectionId: ID) =>
  request(`/plugins/${plugin}/connections/${connectionId}`);

export const getDataScopes = (plugin: string, connectionId: ID) =>
  request(`/plugins/${plugin}/connections/${connectionId}/scopes`);

export const runBlueprint = (id: ID, skipCollectors: boolean) =>
  request(`/blueprints/${id}/trigger`, { method: 'post', data: { skipCollectors } });
