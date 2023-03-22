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

import { useState, useEffect, useMemo } from 'react';

import type { PluginConfigType } from '@/plugins';
import { PluginConfig } from '@/plugins';
import { operator } from '@/utils';

import * as API from './api';

export interface UseTransformationProps {
  plugin: string;
  connectionId: ID;
  id?: ID;
  onCancel?: () => void;
}

export const useTransformation = ({ plugin, connectionId, id, onCancel }: UseTransformationProps) => {
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [name, setName] = useState('');
  const [transformation, setTransformation] = useState({});

  const config = useMemo(() => PluginConfig.find((p) => p.plugin === plugin) as PluginConfigType, []);

  useEffect(() => {
    setLoading(true);
    (async () => {
      if (id) {
        const transformation = await API.getTransformation(plugin, connectionId, id);
        setTransformation(transformation);
        setName(transformation.name);
      } else {
        setTransformation(config.transformation);
      }
    })();
    setLoading(false);
  }, [id, config]);

  const handleSave = async () => {
    const request = id
      ? API.updateTransformation(plugin, connectionId, id, {
          ...transformation,
          name,
        })
      : API.createTransformation(plugin, connectionId, {
          ...transformation,
          name,
        });

    const [success] = await operator(() => request, {
      setOperating: setSaving,
      formatReason: (err) => (err as any).response?.data?.message,
    });

    if (success) {
      onCancel?.();
    }
  };

  return useMemo(
    () => ({
      loading,
      name,
      setName,
      transformation,
      setTransformation: setTransformation,
      saving,
      onSave: handleSave,
    }),
    [plugin, connectionId, id, name, transformation, saving],
  );
};
