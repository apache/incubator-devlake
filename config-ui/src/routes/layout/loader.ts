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

import { AxiosError } from 'axios';
import { json } from 'react-router-dom';

import API from '@/api';
import { getRegisterPlugins } from '@/plugins';
import { ErrorEnum } from '@/routes/error';

type Props = {
  request: Request;
};

export const loader = async ({ request }: Props) => {
  let version = 'unknow';
  let plugins = [];

  try {
    const envPlugins = import.meta.env.DEVLAKE_PLUGINS.split(',').filter(Boolean);
    plugins = getRegisterPlugins().filter((plugin) => !envPlugins.length || envPlugins.includes(plugin));
  } catch (err) {
    plugins = getRegisterPlugins();
  }

  try {
    const res = await API.version(request.signal);
    version = res.version;
  } catch (err) {
    const status = (err as AxiosError).response?.status;
    if (status === 428) {
      throw json({ error: ErrorEnum.NEEDS_DB_MIRGATE }, { status: 428 });
    }
    throw json({ error: ErrorEnum.API_OFFLINE }, { status: 503 });
  }

  return {
    version,
    plugins,
  };
};
