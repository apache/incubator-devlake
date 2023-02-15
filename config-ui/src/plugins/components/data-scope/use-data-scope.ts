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
import { omit } from 'lodash';

import { transformEntities } from '@/config';
import { operator } from '@/utils';

import * as API from './api';

export interface UseDataScope {
  plugin: string;
  connectionId: ID;
  entities: string[];
  initialValues?: {
    scope?: any;
    entites?: string[];
  };
  onSave?: (scope: any) => void;
}

export const useDataScope = ({ plugin, connectionId, entities, initialValues, onSave }: UseDataScope) => {
  const [saving, setSaving] = useState(false);
  const [selectedScope, setSelectedScope] = useState<any>([]);
  const [selectedEntities, setSelectedEntities] = useState<any>([]);

  useEffect(() => {
    (async () => {
      const scope = await Promise.all(
        (initialValues?.scope ?? []).map((sc: any) => API.getDataScope(plugin, connectionId, sc.id)),
      );
      setSelectedScope(scope);
    })();
  }, [initialValues?.scope]);

  useEffect(() => {
    setSelectedEntities(transformEntities(initialValues?.entites ?? entities));
  }, [entities, initialValues?.entites]);

  const getPluginId = (scope: any) => {
    switch (true) {
      case plugin === 'github':
        return scope.githubId;
      case plugin === 'jira':
        return scope.boardId;
      case plugin === 'gitlab':
        return scope.gitlabId;
      case plugin === 'jenkins':
        return scope.jobFullName;
      case plugin === 'bitbucket':
        return scope.bitbucketId;
    }
  };

  const getDataScope = async (scope: any) => {
    try {
      const res = await API.getDataScope(plugin, connectionId, getPluginId(scope));
      return {
        ...scope,
        transformationRuleId: res.transformationRuleId,
      };
    } catch {
      return scope;
    }
  };

  const handleSave = async () => {
    const scope = await Promise.all(selectedScope.map((sc: any) => getDataScope(sc)));

    const [success, res] = await operator(
      () =>
        API.updateDataScope(plugin, connectionId, {
          data: scope.map((sc: any) => omit(sc, 'from')),
        }),
      {
        setOperating: setSaving,
        hideToast: true,
      },
    );

    if (success) {
      onSave?.(
        res.map((it: any) => ({
          id: `${getPluginId(it)}`,
          entities: selectedEntities.map((it: any) => it.value),
        })),
      );
    }
  };

  return useMemo(
    () => ({
      saving,
      selectedScope,
      selectedEntities,
      onChangeScope: setSelectedScope,
      onChangeEntites: setSelectedEntities,
      onSave: handleSave,
    }),
    [saving, selectedScope, selectedEntities],
  );
};
