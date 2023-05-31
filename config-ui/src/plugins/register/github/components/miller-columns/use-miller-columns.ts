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

import { useState, useEffect, useMemo, useCallback } from 'react';
import { McsID, McsItem } from 'miller-columns-select';

import { useProxyPrefix } from '@/hooks';

import type { ScopeItemType } from '../../types';
import * as API from '../../api';
import { getConnection } from '@/pages/blueprint/connection-detail/api';

const DEFAULT_PAGE_SIZE = 30;

export type ExtraType = {
  type: 'org' | 'repo';
} & ScopeItemType;

type MapPageType = Record<McsID | 'root', number>;

export interface UseMillerColumnsProps {
  connectionId: string | number;
}

export const useMillerColumns = ({ connectionId }: UseMillerColumnsProps) => {
  const [user, setUser] = useState<any>({});
  const [items, setItems] = useState<McsItem<ExtraType>[]>([]);
  const [loadedIds, setLoadedIds] = useState<McsID[]>([]);
  const [mapPage, setMapPage] = useState<MapPageType>({});

  const prefix = useProxyPrefix({
    plugin: 'github',
    connectionId,
  });

  const formatOrgs = (orgs: any, parentId: McsID | null = null) =>
    orgs.map((it: any) => ({
      parentId,
      id: it.id,
      title: it.login,
      type: 'org',
    }));

  const formatRepos = (repos: any, parentId: McsID | null = null) =>
    repos.map((it: any) => ({
      parentId,
      id: it.id,
      title: it.name,
      type: 'repo',
      githubId: it.id,
      name: `${it.owner.login}/${it.name}`,
      ownerId: it.owner.id,
      language: it.language,
      description: it.description,
      cloneUrl: it.clone_url,
      HTMLUrl: it.html_url,
    }));

  const setLoaded = useCallback(
    (loaded: boolean, id: McsID, nextPage: number) => {
      if (loaded) {
        setLoadedIds([...loadedIds, id]);
      } else {
        setMapPage({ ...mapPage, [`${id}`]: nextPage });
      }
    },
    [loadedIds, mapPage],
  );

  useEffect(() => {
    (async () => {
      const connection = await getConnection('github', connectionId);

      if (connection.authMethod === 'AppKey') {
        const appInstallationRepos = await API.getInstallationRepos(prefix, {
          page: 1,
          per_page: 1,
        });

        setUser(null);
        setLoaded(true, 'root', 2);

        if (appInstallationRepos.total_count === 0) {
          setItems([]);
        } else {
          setItems([
            {
              parentId: null,
              id: appInstallationRepos.repositories[0].owner.login,
              title: appInstallationRepos.repositories[0].owner.login,
              type: 'org',
            } as any,
          ]);
        }
      } else {
        const user = await API.getUser(prefix);
        const orgs = await API.getUserOrgs(prefix, {
          page: 1,
          per_page: DEFAULT_PAGE_SIZE,
        });

        const loaded = !orgs.length || orgs.length < DEFAULT_PAGE_SIZE;

        setUser(user);
        setLoaded(loaded, 'root', 2);
        setItems([
          {
            parentId: null,
            id: user.login,
            title: user.login,
            type: 'org',
          },
          ...formatOrgs(orgs),
        ]);
      }
    })();
  }, [prefix]);

  const onExpand = useCallback(
    async (id: McsID) => {
      const item = items.find((it) => it.id === id) as McsItem<ExtraType>;
      let repos = [];

      if (user && id === user.login) {
        repos = await API.getUserRepos(prefix, {
          page: 1,
          per_page: DEFAULT_PAGE_SIZE,
        });
      } else if (user) {
        repos = await API.getOrgRepos(prefix, item.title, {
          page: 1,
          per_page: DEFAULT_PAGE_SIZE,
        });
      } else {
        const response = await API.getInstallationRepos(prefix, {
          page: 1,
          per_page: DEFAULT_PAGE_SIZE,
        });
        repos = response.repositories;
      }

      const loaded = !repos.length || repos.length < DEFAULT_PAGE_SIZE;
      setLoaded(loaded, id, 2);
      setItems([...items, ...formatRepos(repos, id)]);
    },
    [items, prefix],
  );

  const onScroll = async (id: McsID | null) => {
    const page = mapPage[id ?? 'root'];
    let orgs = [];
    let repos = [];
    let loaded = false;

    if (id) {
      const item = items.find((it) => it.id === id) as McsItem<ExtraType>;

      if (user && id === user.login) {
        repos = await API.getUserRepos(prefix, {
          page: 1,
          per_page: DEFAULT_PAGE_SIZE,
        });
      } else if (user) {
        repos = await API.getOrgRepos(prefix, item.title, {
          page: 1,
          per_page: DEFAULT_PAGE_SIZE,
        });
      } else {
        const response = await API.getInstallationRepos(prefix, {
          page: 1,
          per_page: DEFAULT_PAGE_SIZE,
        });
        repos = response.repositories;
      }

      loaded = !repos.length || repos.length < DEFAULT_PAGE_SIZE;
    } else {
      orgs = await API.getUserOrgs(prefix, {
        page,
        per_page: DEFAULT_PAGE_SIZE,
      });

      loaded = !orgs.length || orgs.length < DEFAULT_PAGE_SIZE;
    }

    setLoaded(loaded, id ?? 'root', page + 1);
    setItems([...items, ...formatOrgs(orgs), ...formatRepos(repos, id)]);
  };

  return useMemo(
    () => ({
      items,
      getHasMore(id: McsID | null) {
        return !loadedIds.includes(id ?? 'root');
      },
      onExpand,
      onScroll,
    }),
    [items, loadedIds],
  );
};
