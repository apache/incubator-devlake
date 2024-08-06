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

import { encodeName } from '@/routes';

const PATH_PREFIX = import.meta.env.DEVLAKE_PATH_PREFIX ?? '';

export const PATHS = {
  ROOT: () => `${PATH_PREFIX}/`,
  CONNECTIONS: () => `${PATH_PREFIX}/connections`,
  CONNECTION: (plugin: string, connectionId: ID) => `${PATH_PREFIX}/connections/${plugin}/${connectionId}`,
  PROJECTS: () => `${PATH_PREFIX}/projects`,
  PROJECT: (pname: string) => `${PATH_PREFIX}/projects/${encodeName(pname)}`,
  PROJECT_CONNECTION: (pname: string, plugin: string, connectionId: ID) =>
    `${PATH_PREFIX}/projects/${encodeName(pname)}/${plugin}-${connectionId}`,
  BLUEPRINTS: () => `${PATH_PREFIX}/advanced/blueprints`,
  BLUEPRINT: (bid: ID) => `${PATH_PREFIX}/advanced/blueprints/${bid}`,
  BLUEPRINT_CONNECTION: (bid: ID, plugin: string, connectionId: ID) =>
    `${PATH_PREFIX}/advanced/blueprints/${bid}/${plugin}-${connectionId}`,
  PIPELINES: () => `${PATH_PREFIX}/advanced/pipelines`,
  APIKEYS: () => `${PATH_PREFIX}/keys`,
};
