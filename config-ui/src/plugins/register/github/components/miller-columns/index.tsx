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

import React, { useState, useEffect } from 'react';
import type { McsID, McsItem, McsColumn } from 'miller-columns-select';
import MillerColumnsSelect from 'miller-columns-select';

import { Loading } from '@/components';

import type { ScopeItemType } from '../../types';

import type { UseMillerColumnsProps, ExtraType } from './use-miller-columns';
import { useMillerColumns } from './use-miller-columns';
import * as S from './styled';

interface Props extends UseMillerColumnsProps {
  selectedItems: ScopeItemType[];
  onChangeItems: (selectedItems: ScopeItemType[]) => void;
}

export const MillerColumns = ({ connectionId, selectedItems, onChangeItems }: Props) => {
  const [selectedIds, setSelectedIds] = useState<McsID[]>([]);

  const { items, getHasMore, onExpand, onScroll } = useMillerColumns({
    connectionId,
  });

  useEffect(() => {
    setSelectedIds(selectedItems.map((it) => it.githubId));
  }, [selectedItems]);

  const handleChangeItems = (selectedIds: McsID[]) => {
    const result = selectedIds.map((id) => {
      const selectedItem = selectedItems.find((it) => it.githubId === id);
      if (selectedItem) {
        return selectedItem;
      }

      const item = items.find((it) => it.id === id) as McsItem<ExtraType>;
      return {
        connectionId,
        githubId: item.githubId,
        name: item.name,
        ownerId: item.ownerId,
        language: item.language,
        description: item.description,
        cloneUrl: item.cloneUrl,
        HTMLUrl: item.HTMLUrl,
      };
    });

    onChangeItems(result);
  };

  const renderTitle = (column: McsColumn) => {
    return !column.parentId && <S.ColumnTitle>Organizations/Owners</S.ColumnTitle>;
  };

  const renderLoading = () => {
    return <Loading size={20} style={{ padding: '4px 12px' }} />;
  };

  return (
    <MillerColumnsSelect
      items={items}
      getCanExpand={(it) => it.type === 'org'}
      getHasMore={getHasMore}
      onExpand={onExpand}
      onScroll={onScroll}
      columnCount={2}
      columnHeight={300}
      renderTitle={renderTitle}
      renderLoading={renderLoading}
      selectedIds={selectedIds}
      onSelectItemIds={handleChangeItems}
    />
  );
};
