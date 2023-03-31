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

export const getBlueprint = (id: ID) => request(`/blueprints/${id}`);

export const updateBlueprint = (id: ID, payload: any) =>
  request(`/blueprints/${id}`, {
    method: 'patch',
    data: payload,
  });

export const createTransformation = (plugin: string, connectionId: ID, payload: any) =>
  request(`/plugins/${plugin}/connections/${connectionId}/transformation_rules`, {
    method: 'post',
    data: payload,
  });

export const getGitHub = (prefix: string, owner: string, repo: string) => request(`${prefix}/repos/${owner}/${repo}`);

export const getGitLab = (prefix: string, id: ID) => request(`${prefix}/projects/${id}`);

export const getJira = (prefix: string, id: ID) => request(`${prefix}/agile/1.0/board/${id}`);

export const updateDataScope = (plugin: string, connectionId: ID, repoId: ID, payload: any) =>
  request(`/plugins/${plugin}/connections/${connectionId}/scopes/${repoId}`, {
    method: 'patch',
    data: payload,
  });
