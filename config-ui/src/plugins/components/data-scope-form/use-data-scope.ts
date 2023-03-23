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

import { operator } from '@/utils';

import * as API from './api';

export interface UseDataScope {
  plugin: string;
  connectionId: ID;
  initialScope: any[];
  initialEntities: string[];
  onSubmit?: (scope: Array<{ id: string; entities: string[] }>, origin: any) => void;
}

export const useDataScope = ({ plugin, connectionId, initialScope, initialEntities, onSubmit }: UseDataScope) => {
  const [saving, setSaving] = useState(false);
  const [scope, setScope] = useState<any>([]);
  const [entities, setEntites] = useState<any>([]);

  useEffect(() => {
    (async () => {
      const scope = await Promise.all(initialScope.map((sc: any) => API.getDataScope(plugin, connectionId, sc.id)));
      setScope(scope);
    })();
  }, []);

  useEffect(() => {
    setEntites(initialEntities);
  }, []);

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
      case plugin === 'zentao':
        return scope.type === 'project' ? `project/${scope.id}` : `product/${scope.id}`;
      case plugin === 'sonarqube':
        return scope.projectKey;
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
    const data = await Promise.all(scope.map((sc: any) => getDataScope(sc)));

    let request: () => Promise<any>;
    if (plugin === 'zentao') {
      request = async () => {
        return [
          ...(await API.updateDataScopeWithType(plugin, connectionId, 'product', {
            data: data.filter((s) => s.type !== 'project'),
          })),
          ...(await API.updateDataScopeWithType(plugin, connectionId, 'project', {
            data: data.filter((s) => s.type === 'project'),
          })),
        ];
      };
    } else {
      request = () =>
        API.updateDataScope(plugin, connectionId, {
          data,
        });
    }

    const [success, res] = await operator(request, {
      setOperating: setSaving,
      hideToast: true,
    });

    if (success) {
      onSubmit?.(
        res.map((it: any) => ({
          id: `${getPluginId(it)}`,
          entities,
        })),
        res,
      );
    }
  };

  return useMemo(
    () => ({
      scope,
      setScope,
      entities,
      setEntites,
      saving,
      onSave: handleSave,
    }),
    [scope, entities, saving],
  );
};
