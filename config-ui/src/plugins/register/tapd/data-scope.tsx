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

import { useEffect, useState } from 'react';
import { Button, ControlGroup, InputGroup, Intent } from '@blueprintjs/core';
import type { McsID, McsItem } from 'miller-columns-select';
import MillerColumnsSelect from 'miller-columns-select';

import API from '@/api';
import { ExternalLink, Loading } from '@/components';
import * as T from '@/plugins/components/data-scope-remote/types';
// import * as API from '@/plugins/components/data-scope-remote/api';

// import { prepareToken } from './api';

interface Props {
  connectionId: ID;
  disabledItems: T.ResItem[];
  selectedItems: T.ResItem[];
  onChangeSelectedItems: (items: T.ResItem[]) => void;
}

export const DataScope = ({ connectionId, disabledItems, selectedItems, onChangeSelectedItems }: Props) => {
  const [pageToken, setPageToken] = useState<string | undefined>(undefined);
  const [companyId, setCompanyId] = useState<string>(
    localStorage.getItem(`plugin/tapd/connections/${connectionId}/company_id`) || '',
  );

  const [miller, setMiller] = useState<{
    items: McsItem<T.ResItem>[];
    loadedIds: ID[];
  }>({
    items: [],
    loadedIds: [],
  });

  useEffect(() => {
    if (!pageToken) return;
    getItems(null, pageToken);
  }, [pageToken]);

  const getItems = async (groupId: ID | null, currentPageToken?: string) => {
    const res = await API.scope.remote('tapd', connectionId, {
      groupId,
      pageToken: currentPageToken,
    });

    const newItems = (res.children ?? []).map((it) => ({
      ...it,
      title: it.name,
    }));

    setMiller((m) => ({
      ...m,
      items: [...m.items, ...newItems],
      loadedIds: [...m.loadedIds, groupId ? groupId : 'root'],
    }));
  };

  const getPageToken = async (companyId: string | undefined) => {
    if (!companyId) {
      setPageToken(undefined);
      return;
    }
    const res = await API.plugin.tapd.remoteScopePrepareToken(connectionId, {
      companyId,
    });
    setPageToken(res.pageToken);
  };

  return (
    <>
      <h4>Workspaces *</h4>
      <p>Type in the company ID to list all the workspaces you want to sync. </p>
      <ExternalLink link="https://www.tapd.cn/help/show#1120003271001000103">
        Learn about how to get your company ID
      </ExternalLink>

      <ControlGroup fill={false} vertical={false} style={{ padding: '8px 0' }}>
        <InputGroup
          placeholder="Your company ID"
          value={companyId}
          style={{ width: 300 }}
          onChange={(e) => {
            setCompanyId(e.target.value);
            localStorage.setItem(`plugin/tapd/connections/${connectionId}/company_id`, e.target.value);
          }}
        />
        <Button intent={Intent.PRIMARY} onClick={() => getPageToken(companyId)}>
          Search
        </Button>
      </ControlGroup>

      {pageToken && (
        <MillerColumnsSelect
          items={miller.items}
          getCanExpand={(it) => it.type === 'group'}
          getHasMore={(id) => !miller.loadedIds.includes(id ?? 'root')}
          onExpand={(id: McsID) => getItems(id, pageToken)}
          onScroll={(id: McsID | null) => getItems(id, pageToken)}
          columnCount={2.5}
          columnHeight={300}
          renderLoading={() => <Loading size={20} style={{ padding: '4px 12px' }} />}
          disabledIds={(disabledItems ?? []).map((it) => it.id)}
          selectedIds={selectedItems.map((it) => it.id)}
          onSelectItemIds={(selectedIds: ID[]) =>
            onChangeSelectedItems(miller.items.filter((it) => selectedIds.includes(it.id)))
          }
        />
      )}
    </>
  );
};
