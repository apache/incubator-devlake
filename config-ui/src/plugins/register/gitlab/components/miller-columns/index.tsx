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

import React, { useEffect, useState } from 'react';
import MillerColumnsSelect from 'miller-columns-select';

import { Loading } from '@/components';

import type { ScopeItemType } from '../../types';

import type { UseMillerColumnsProps, GitLabColumnType } from './use-miller-columns';
import { useMillerColumns } from './use-miller-columns';
import * as S from './styled';

interface Props extends UseMillerColumnsProps {
  selectedItems: ScopeItemType[];
  onChangeItems: (selectedItems: ScopeItemType[]) => void;
}

export const MillerColumns = ({ connectionId, selectedItems, onChangeItems }: Props) => {
  const [selectedIds, setSelectedIds] = useState<ID[]>([]);

  const { items, getHasMore, onExpandItem, onScrollColumn } = useMillerColumns({
    connectionId,
  });

  useEffect(() => {
    setSelectedIds(selectedItems.map((it) => it.gitlabId));
  }, [selectedItems]);

  const handleChangeItems = (selectedIds: ID[]) => {
    const result = items
      .filter((it) => (selectedIds.length ? selectedIds.includes(it.id) : false))
      .map((it) => ({
        connectionId,
        gitlabId: it.id,
        name: it.name,
        pathWithNamespace: it.pathWithNamespace,
        creatorId: it.creatorId,
        defaultBranch: it.defaultBranch,
        description: it.description,
        openIssuesCount: it.openIssuesCount,
        starCount: it.starCount,
        visibility: it.visibility,
        webUrl: it.webUrl,
        httpUrlToRepo: it.httpUrlToRepo,
      }));

    onChangeItems(result);
  };

  const renderTitle = (column: GitLabColumnType) => {
    return !column.parentId && <S.ColumnTitle>Subgroups/Projects</S.ColumnTitle>;
  };

  const renderLoading = () => {
    return <Loading size={20} style={{ padding: '4px 12px' }} />;
  };

  return (
    <MillerColumnsSelect
      columnCount={2.5}
      columnHeight={300}
      getCanExpand={(it) => it.type === 'group'}
      getHasMore={getHasMore}
      renderTitle={renderTitle}
      renderLoading={renderLoading}
      items={items}
      selectedIds={selectedIds}
      onSelectItemIds={handleChangeItems}
      onExpandItem={onExpandItem}
      onScrollColumn={onScrollColumn}
    />
  );
};
