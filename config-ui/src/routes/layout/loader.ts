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

import { redirect } from 'react-router-dom';
import { intersection } from 'lodash';

import API from '@/api';
import { PATHS } from '@/config';
import { getRegisterPlugins } from '@/plugins';

type Props = {
  request: Request;
};

const normalizePath = (path: string) => (path.length > 1 ? path.replace(/\/+$/, '') : path);

const isOrgManagementPath = (path: string) => {
  const normalizedPath = normalizePath(path);
  return [normalizePath(PATHS.TEAMS()), normalizePath(PATHS.USERS())].includes(normalizedPath);
};

const parseEnabledPlugins = (plugins: string[]) => {
  const envPlugins = import.meta.env.DEVLAKE_PLUGINS;

  if (typeof envPlugins !== 'string') {
    return plugins;
  }

  const enabledPlugins = envPlugins
    .split(',')
    .map((plugin: string) => plugin.trim())
    .filter(Boolean);

  if (!enabledPlugins.length) {
    return plugins;
  }

  return plugins.filter((plugin) => enabledPlugins.includes(plugin));
};

const hasOrgApis = async () => {
  const [teamsResult, usersResult] = await Promise.allSettled([
    API.team.list({ page: 1, pageSize: 1, grouped: false }),
    API.user.list({ page: 1, pageSize: 1 }),
  ]);

  return teamsResult.status === 'fulfilled' && usersResult.status === 'fulfilled';
};

export const layoutLoader = async ({ request }: Props) => {
  const onboard = await API.store.get('onboard');

  if (!onboard) {
    return redirect('/onboard');
  }

  const fePlugins = parseEnabledPlugins(getRegisterPlugins());
  const bePlugins = await API.plugin.list();
  const bePluginNames = bePlugins.map((it) => it.plugin);
  const orgCapabilityAvailable = bePluginNames.includes('org') && (await hasOrgApis());

  if (!orgCapabilityAvailable && isOrgManagementPath(new URL(request.url).pathname)) {
    return redirect(PATHS.CONNECTIONS());
  }

  const res = await API.version(request.signal);

  return {
    version: res.version,
    plugins: intersection(fePlugins, bePluginNames),
    orgCapabilityAvailable,
  };
};
