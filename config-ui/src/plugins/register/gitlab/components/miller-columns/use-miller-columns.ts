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
import { McsID, McsItem } from 'miller-columns-select';

import { useProxyPrefix } from '@/hooks';

import type { ScopeItemType } from '../../types';
import * as API from '../../api';

const DEFAULT_PAGE_SIZE = 20;

export type ExtraType = {
  type: 'group' | 'project';
} & ScopeItemType;

type MapValueType = {
  groupPage: number;
  groupLoaded: boolean;
  projectPage: number;
  projectLoaded: boolean;
};

type MapType = Record<ID, MapValueType>;

export interface UseMillerColumnsProps {
  connectionId: ID;
}

export const useMillerColumns = ({ connectionId }: UseMillerColumnsProps) => {
  const [user, setUser] = useState<any>({});
  const [items, setItems] = useState<McsItem<ExtraType>[]>([]);
  const [map, setMap] = useState<MapType>({});

  const prefix = useProxyPrefix({
    plugin: 'gitlab',
    connectionId,
  });

  const formatGroups = (arr: any, parentId: ID | null = null): McsItem<ExtraType>[] =>
    arr.map((it: any) => ({
      parentId,
      id: it.id,
      title: it.name,
      type: 'group',
    }));

  const formatProjects = (arr: any, parentId: ID | null = null): McsItem<ExtraType>[] =>
    arr.map((it: any) => ({
      parentId,
      id: it.id,
      title: it.name,
      type: 'project',
      gitlabId: it.id,
      name: it.path_with_namespace,
      pathWithNamespace: it.path_with_namespace,
      creatorId: it.creator_id,
      defaultBranch: it.default_branch,
      description: it.description,
      openIssuesCount: it.open_issues_count,
      starCount: it.star_count,
      visibility: it.visibility,
      webUrl: it.web_url,
      httpUrlToRepo: it.http_url_to_repo,
    }));

  const setLoaded = (id: ID, params: MapValueType) => {
    setMap({
      ...map,
      [`${id}`]: params,
    });
  };

  useEffect(() => {
    (async () => {
      const user = await API.getUser(prefix);
      setUser(user);

      let groupLoaded = false;
      let projectLoaded = false;
      let groups = [];
      let projects = [];

      groups = await API.getUserGroups(prefix, {
        page: 1,
        per_page: DEFAULT_PAGE_SIZE,
      });

      groupLoaded = !groups.length || groups.length < DEFAULT_PAGE_SIZE;

      if (groupLoaded) {
        projects = await API.getUserProjects(prefix, user.id, {
          page: 1,
          per_page: DEFAULT_PAGE_SIZE,
        });

        projectLoaded = !projects.length || projects.length < DEFAULT_PAGE_SIZE;
      }

      setLoaded('root', {
        groupLoaded,
        groupPage: groupLoaded ? 1 : 2,
        projectLoaded,
        projectPage: projectLoaded ? 1 : 2,
      });
      setItems([...formatGroups(groups), ...formatProjects(projects)]);
    })();
  }, [prefix]);

  const onExpand = async (id: McsID) => {
    let groupLoaded = false;
    let projectLoaded = false;
    let groups = [];
    let projects = [];

    groups = await API.getGroupSubgroups(prefix, id, {
      page: 1,
      per_page: DEFAULT_PAGE_SIZE,
    });

    groupLoaded = !groups.length || groups.length < DEFAULT_PAGE_SIZE;

    if (groupLoaded) {
      projects = await API.getGroupProjects(prefix, id, {
        page: 1,
        per_page: DEFAULT_PAGE_SIZE,
      });

      projectLoaded = !projects.length || projects.length < DEFAULT_PAGE_SIZE;
    }

    setLoaded(id, {
      groupLoaded,
      groupPage: groupLoaded ? 1 : 2,
      projectLoaded,
      projectPage: projectLoaded ? 1 : 2,
    });
    setItems([...items, ...formatGroups(groups, id), ...formatProjects(projects, id)]);
  };

  const onScroll = async (id: McsID | null) => {
    const mapValue = map[id ?? 'root'];

    let groupLoaded = mapValue.groupLoaded;
    let projectLoaded = mapValue.projectLoaded;
    let groups = [];
    let projects = [];

    if (!groupLoaded) {
      groups = id
        ? await API.getGroupSubgroups(prefix, id, {
            page: mapValue.groupPage,
            per_page: DEFAULT_PAGE_SIZE,
          })
        : await API.getUserGroups(prefix, {
            page: mapValue.groupPage,
            per_page: DEFAULT_PAGE_SIZE,
          });

      groupLoaded = !groups.length || groups.length < DEFAULT_PAGE_SIZE;
    } else if (!projectLoaded) {
      projects = id
        ? await API.getGroupProjects(prefix, id, {
            page: mapValue.projectPage,
            per_page: DEFAULT_PAGE_SIZE,
          })
        : await API.getUserProjects(prefix, user.id, {
            page: mapValue.projectPage,
            per_page: DEFAULT_PAGE_SIZE,
          });

      projectLoaded = !projects.length || projects.length < DEFAULT_PAGE_SIZE;
    }

    setLoaded(id ?? 'root', {
      groupLoaded,
      groupPage: groupLoaded ? mapValue.groupPage : mapValue.groupPage + 1,
      projectLoaded,
      projectPage: projectLoaded ? mapValue.projectPage : mapValue.projectPage + 1,
    });
    setItems([...items, ...formatGroups(groups, id), ...formatProjects(projects, id)]);
  };

  return useMemo(
    () => ({
      items,
      getHasMore(id: McsID | null) {
        const mapValue = map[id ?? 'root'];
        return !(mapValue?.groupLoaded && mapValue?.projectLoaded);
      },
      onExpand,
      onScroll,
    }),
    [items, map],
  );
};
