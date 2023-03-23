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

import React, { useMemo } from 'react';

import { DataScopeMillerColumns } from '@/plugins';

import type { ScopeItemType } from './types';

interface Props {
  connectionId: ID;
  selectedItems: ScopeItemType[];
  onChangeItems: (selectedItems: ScopeItemType[]) => void;
}

export const ZentaoDataScope = ({ connectionId, onChangeItems, ...props }: Props) => {
  const selectedItems = useMemo(
    () =>
      props.selectedItems.map((it) => ({
        id: it.type === 'project' ? `project/${it.id}` : `product/${it.id}`,
        name: it.name,
        data: it,
      })),
    [props.selectedItems],
  );

  return (
    <>
      <h3>Repositories *</h3>
      <p>Select the repositories you would like to sync.</p>
      <DataScopeMillerColumns
        plugin="zentao"
        connectionId={connectionId}
        selectedItems={selectedItems}
        onChangeItems={onChangeItems}
      />
    </>
  );
};
