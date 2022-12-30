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

import { useState, useMemo, useEffect } from 'react';
import { useHistory } from 'react-router-dom';

import { Plugins } from '@/plugins';
import { operator } from '@/utils';

import * as API from './api';

interface Props {
  plugin: Plugins;
  id?: ID;
}

export const useForm = ({ plugin, id }: Props) => {
  const [loading, setLoading] = useState(false);
  const [operating, setOperating] = useState(false);
  const [connection, setConnection] = useState<any>({});

  const history = useHistory();

  const getConnection = async () => {
    if (!id) return;

    setLoading(true);
    try {
      const res = await API.getConnection(plugin, id);
      setConnection(res);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    getConnection();
  }, []);

  const handleTest = async (payload: any) => {
    const [success] = await operator(() => API.testConnection(plugin, payload), {
      setOperating,
      formatReason: (err) => (err as any)?.response?.data?.message,
    });

    if (success) {
    }
  };

  const handleCreate = async (payload: any) => {
    const [success] = await operator(() => API.createConnection(plugin, payload), {
      setOperating,
    });

    if (success) {
      history.push(`/connections/${plugin}`);
    }
  };

  const handleUpdate = async (id: ID, payload: any) => {
    const [success] = await operator(() => API.updateConnection(plugin, id, payload), {
      setOperating,
    });

    if (success) {
      history.push(`/connections/${plugin}`);
    }
  };

  return useMemo(
    () => ({
      loading,
      operating,
      connection,
      onTest: handleTest,
      onCreate: handleCreate,
      onUpdate: handleUpdate,
    }),
    [loading, operating, connection],
  );
};
