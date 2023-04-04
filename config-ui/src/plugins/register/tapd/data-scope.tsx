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

import React, { useMemo, useState } from 'react';

import { DataScopeMillerColumns } from '@/plugins';

import type { ScopeItemType } from './types';
import * as API from '@/plugins/components/data-scope-miller-columns/api';
import { Button, ControlGroup, InputGroup, Intent } from '@blueprintjs/core';
import { ExternalLink } from '@/components';

interface Props {
  connectionId: ID;
  selectedItems: ScopeItemType[];
  onChangeItems: (selectedItems: ScopeItemType[]) => void;
}

export const TapdDataScope = ({ connectionId, onChangeItems, ...props }: Props) => {
  const selectedItems = useMemo(
    () => props.selectedItems.map((it) => ({ id: `${it.id}`, name: it.name, data: it })),
    [props.selectedItems],
  );

  const [pageToken, setPageToken] = useState<string | undefined>(undefined);
  const [companyId, setCompanyId] = useState<string>(
    localStorage.getItem(`plugin/tapd/connections/${connectionId}/company_id`) || '',
  );

  const getPageToken = async (companyId: string | undefined) => {
    if (!companyId) {
      setPageToken(undefined);
      return;
    }
    const res = await API.prepareToken(`tapd`, connectionId, {
      companyId,
    });
    setPageToken(res.pageToken);
  };

  return (
    <>
      <h3>Workspaces *</h3>
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
        <DataScopeMillerColumns
          key={pageToken}
          plugin="tapd"
          connectionId={connectionId}
          selectedItems={selectedItems}
          onChangeItems={onChangeItems}
          pageToken={pageToken}
        />
      )}
    </>
  );
};
