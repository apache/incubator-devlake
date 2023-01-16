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

import { useProxyPrefix } from '@/hooks';

import type { ScopeItemType } from '../../types';
import * as API from '../../api';

export interface UseProjectSelectorProps {
  connectionId: ID;
}

export const useProjectSelector = ({ connectionId }: UseProjectSelectorProps) => {
  const [loading, setLoading] = useState(false);
  const [items, setItems] = useState<ScopeItemType[]>([]);
  const [search, setSearch] = useState('');
  const [membership, setMembership] = useState(true);

  const prefix = useProxyPrefix({
    plugin: 'gitlab',
    connectionId,
  });

  useEffect(() => {
    if (!search) return;
    setItems([]);
    setLoading(true);

    const timer = setTimeout(async () => {
      try {
        const res = await API.searchProject(prefix, {
          search,
          membership,
        });
        setItems(
          res.map((it: any) => ({
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
          })),
        );
      } finally {
        setLoading(false);
      }
    }, 1000);

    return () => clearTimeout(timer);
  }, [prefix, search, membership]);

  return useMemo(
    () => ({
      loading,
      items,
      membership,
      onSearch(s: string) {
        setSearch(s);
      },
      onChangeMembership(e: React.ChangeEvent<HTMLInputElement>) {
        setMembership(e.target.checked);
      },
    }),
    [loading, items, membership],
  );
};
