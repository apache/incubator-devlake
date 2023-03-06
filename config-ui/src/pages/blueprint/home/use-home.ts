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

import { cronPresets } from '@/config';

import type { BlueprintType } from '../types';

import * as API from './api';

export const useHome = () => {
  const [loading, setLoading] = useState(false);
  const [blueprints, setBlueprints] = useState<BlueprintType[]>([]);
  const [dataSource, setDataSource] = useState<BlueprintType[]>([]);
  const [type, setType] = useState('all');

  const presets = cronPresets.map((preset) => preset.config);

  const getBlueprints = async () => {
    setLoading(true);
    try {
      const res = await API.getBlueprints({
        page: 1,
        pageSize: 200,
      });
      setBlueprints(res.blueprints);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    getBlueprints();
  }, []);

  useEffect(() => {
    setDataSource(
      blueprints.filter((bp) => {
        switch (type) {
          case 'all':
            return true;
          case 'manual':
            return bp.isManual;
          case 'custom':
            return !presets.includes(bp.cronConfig);
          default:
            return !bp.isManual && bp.cronConfig === type;
        }
      }),
    );
  }, [blueprints, type]);

  return useMemo(
    () => ({
      loading,
      blueprints,
      dataSource,
      type,
      onChangeType: setType,
    }),
    [loading, blueprints, dataSource, type],
  );
};
