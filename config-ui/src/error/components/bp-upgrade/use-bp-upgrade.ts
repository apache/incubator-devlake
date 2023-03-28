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

import { useMemo, useState } from 'react';

import { operator } from '@/utils';

import * as API from './api';

export interface UseBPUpgradeProps {
  id?: ID;
  onResetError: () => void;
}

export const useBPUpgrade = ({ id, onResetError }: UseBPUpgradeProps) => {
  const [processing, setProcessing] = useState(false);

  const getScopeId = (plugin: string, scope: any) => {
    switch (plugin) {
      case 'github':
        return scope.githubId;
      case 'gitlab':
        return scope.gitlabId;
      case 'jira':
        return scope.boardId;
    }
  };

  const getScopeDetail = async (plugin: string, connectionId: ID, options: any) => {
    const prefix = `/plugins/${plugin}/connections/${connectionId}/proxy/rest`;

    if (plugin === 'github') {
      const res = await API.getGitHub(prefix, options.owner, options.repo);
      return {
        connectionId,
        githubId: res.id,
        name: `${res.owner.login}/${res.name}`,
        ownerId: res.owner.id,
        language: res.language,
        description: res.description,
        cloneUrl: res.clone_url,
        HTMLUrl: res.html_url,
      };
    }

    if (plugin === 'gitlab') {
      const res = await API.getGitLab(prefix, options.projectId);
      return {
        connectionId,
        gitlabId: res.id,
        name: res.path_with_namespace,
        pathWithNamespace: res.path_with_namespace,
        creatorId: res.creator_id,
        defaultBranch: res.default_branch,
        description: res.description,
        openIssuesCount: res.open_issues_count,
        starCount: res.star_count,
        visibility: res.visibility,
        webUrl: res.web_url,
        httpUrlToRepo: res.http_url_to_repo,
      };
    }

    if (plugin === 'jira') {
      const res = await API.getJira(prefix, options.boardId);
      return {
        connectionId,
        boardId: res.id,
        name: res.name,
        self: res.self,
        type: res.type,
        projectId: res?.location?.projectId,
      };
    }
  };

  const upgradeScope = async (plugin: string, connectionId: ID, scope: any) => {
    let transformationRule;

    if (scope.transformation) {
      // create transfromation template
      transformationRule = await API.createTransformation(plugin, connectionId, {
        ...scope.transformation,
        name: `upgrade-${plugin}-${connectionId}-${new Date().getTime()}`,
      });
    }

    // get data scope detail
    const scopeDetail = await getScopeDetail(plugin, connectionId, scope.options);

    // put data scope
    await API.updateDataScope(plugin, connectionId, getScopeId(plugin, scopeDetail), {
      ...scopeDetail,
      transformationRuleId: transformationRule?.id,
    });

    return {
      id: `${getScopeId(plugin, scopeDetail)}`,
      entities: scope.entities,
    };
  };

  const upgradeConnection = async (connection: any) => {
    const { plugin, connectionId } = connection;

    if (plugin == 'jenkins') {
      return {
        plugin,
        connectionId,
        scopes: [],
      };
    }
    const scope = await Promise.all((connection.scope ?? []).map((sc: any) => upgradeScope(plugin, connectionId, sc)));

    return {
      plugin,
      connectionId,
      scopes: scope,
    };
  };

  const handleUpgrade = async () => {
    if (!id) return;

    const bp = await API.getBlueprint(id);
    const connections = await Promise.all(bp.settings.connections.map((cs: any) => upgradeConnection(cs)));

    await API.updateBlueprint(id, {
      ...bp,
      settings: {
        version: '2.0.0',
        connections,
      },
    });
  };

  const handleSubmit = async () => {
    const [success] = await operator(handleUpgrade, {
      setOperating: setProcessing,
    });

    if (success) {
      onResetError();
    }
  };

  return useMemo(
    () => ({
      processing,
      onSubmit: handleSubmit,
    }),
    [processing],
  );
};
